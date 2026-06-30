package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/PressStart2P-Regular.ttf
var fontBytes []byte

//go:embed assets/food.png
var foodPNG []byte

//go:embed assets/gameover.mp3
var gameoverMP3 []byte

//go:embed assets/eating.mp3
var eatingMP3 []byte

const (
	screenWidth  = 1110
	screenHeight = 720
	gameWidth    = 1080
	gameHeight   = 700
	gridSize     = 20
	sampleRate   = 44100
)

var (
	arcadeFont     font.Face
	smallArcade    font.Face
	foodImage      *ebiten.Image
	audioContext   *audio.Context
	eatPlayer      *audio.Player
	gameoverPlayer *audio.Player
)

type point struct {
	X int
	Y int
}

type game struct {
	snake        []point
	direction    point
	food         point
	timer        int
	score        int
	gameOver     bool
	speed        int
	audioEnabled bool
}

func (g *game) reset() {
	g.snake = []point{{X: gameWidth / 2 / gridSize, Y: gameHeight / 2 / gridSize}}
	g.direction = point{X: 1, Y: 0}
	g.food = g.newFood()
	g.timer = 0
	g.score = 0
	g.gameOver = false
	g.speed = 7
}

func loadAudio() error {
	audioContext = audio.NewContext(sampleRate)

	eatDecoded, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(eatingMP3))
	if err != nil {
		return err
	}
	gameoverDecoded, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(gameoverMP3))
	if err != nil {
		return err
	}

	eatPlayer, err = audio.NewPlayer(audioContext, eatDecoded)
	if err != nil {
		return err
	}
	gameoverPlayer, err = audio.NewPlayer(audioContext, gameoverDecoded)
	if err != nil {
		return err
	}

	return nil
}

func playSound(player *audio.Player, volume float64) {
	if player == nil || audioContext == nil || !audioContext.IsReady() {
		return
	}
	player.Rewind()
	player.SetVolume(volume)
	player.Play()
}

func (g *game) newFood() point {
	for {
		candidate := point{
			X: rand.Intn(gameWidth / gridSize),
			Y: rand.Intn(gameHeight / gridSize),
		}
		occupied := false
		for _, body := range g.snake {
			if body == candidate {
				occupied = true
				break
			}
		}
		if !occupied {
			return candidate
		}
	}
}

func (g *game) Update() error {
	if audioContext != nil {
		g.audioEnabled = audioContext.IsReady()
	}

	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.reset()
		}
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) && g.direction.Y == 0 {
		g.direction = point{X: 0, Y: -1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && g.direction.Y == 0 {
		g.direction = point{X: 0, Y: 1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.direction.X == 0 {
		g.direction = point{X: -1, Y: 0}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.direction.X == 0 {
		g.direction = point{X: 1, Y: 0}
	}

	g.timer++
	if g.timer%g.speed != 0 {
		return nil
	}

	head := g.snake[0]
	newHead := point{X: head.X + g.direction.X, Y: head.Y + g.direction.Y}

	if newHead.X < 0 || newHead.X >= gameWidth/gridSize || newHead.Y < 0 || newHead.Y >= gameHeight/gridSize {
		g.gameOver = true
		playSound(gameoverPlayer, 2)
		return nil
	}

	for _, body := range g.snake {
		if body == newHead {
			g.gameOver = true
			playSound(gameoverPlayer, 2)
			return nil
		}
	}

	g.snake = append([]point{newHead}, g.snake...)

	if newHead == g.food {
		g.food = g.newFood()
		g.score++
		if g.score%10 == 0 && g.speed > 1 {
			g.speed--
		}
		playSound(eatPlayer, 1)
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{85, 103, 120, 255})

	gameArea := ebiten.NewImage(gameWidth, gameHeight)
	gameArea.Fill(color.RGBA{0, 0, 0, 255})

	for _, body := range g.snake {
		ebitenutil.DrawRect(gameArea, float64(body.X*gridSize), float64(body.Y*gridSize), gridSize, gridSize, color.RGBA{0, 205, 50, 255})
	}

	foodOptions := &ebiten.DrawImageOptions{}
	foodOptions.GeoM.Translate(float64(g.food.X*gridSize), float64(g.food.Y*gridSize))
	gameArea.DrawImage(foodImage, foodOptions)

	gameAreaOptions := &ebiten.DrawImageOptions{}
	gameAreaOptions.GeoM.Translate(float64((screenWidth-gameWidth)/2), float64((screenHeight-gameHeight)/2))
	screen.DrawImage(gameArea, gameAreaOptions)

	scoreArea := ebiten.NewImage(220, 120)
	scoreArea.Fill(color.RGBA{0, 0, 0, 0})
	ebitenutil.DrawRect(scoreArea, 0, 20, 220, 120, color.RGBA{148, 73, 237, 128})

	borderColor := color.RGBA{255, 255, 255, 128}
	borderThickness := 6.0
	ebitenutil.DrawRect(scoreArea, 0, 20, 220, borderThickness, borderColor)
	ebitenutil.DrawRect(scoreArea, 0, 120-borderThickness, 220, borderThickness, borderColor)
	ebitenutil.DrawRect(scoreArea, 0, 20, borderThickness, 100, borderColor)
	ebitenutil.DrawRect(scoreArea, 220-borderThickness, 20, borderThickness, 100, borderColor)

	text.Draw(scoreArea, "SCORE", arcadeFont, 30, 65, color.White)
	text.Draw(scoreArea, fmt.Sprintf("%d", g.score), arcadeFont, 30, 105, color.White)

	scoreOptions := &ebiten.DrawImageOptions{}
	scoreOptions.GeoM.Translate(float64(screenWidth-310), 20)
	screen.DrawImage(scoreArea, scoreOptions)

	if g.gameOver {
		overlay := ebiten.NewImage(560, 110)
		overlay.Fill(color.White)
		ebitenutil.DrawRect(overlay, 10, 10, 540, 90, color.RGBA{255, 0, 0, 255})
		text.Draw(overlay, "GAME OVER", arcadeFont, 100, 50, color.White)
		text.Draw(overlay, "Press 'R' to Play", smallArcade, 40, 90, color.White)

		overlayOptions := &ebiten.DrawImageOptions{}
		overlayOptions.GeoM.Translate(float64((screenWidth-560)/2), float64((screenHeight-110)/2))
		screen.DrawImage(overlay, overlayOptions)
	}

	if !g.audioEnabled {
		text.Draw(screen, "Audio locked by browser: click the page and press any key.", smallArcade, 20, screenHeight-20, color.White)
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	arcadeFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{Size: 32, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}
	smallArcade, err = opentype.NewFace(ttf, &opentype.FaceOptions{Size: 28, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}

	foodImage, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(foodPNG))
	if err != nil {
		log.Fatal(err)
	}

	if err := loadAudio(); err != nil {
		log.Fatal(err)
	}

	g := &game{}
	g.reset()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake in Go")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
