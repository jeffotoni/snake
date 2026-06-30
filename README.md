# Snake v7

[![GoDoc](https://godoc.org/github.com/jeffotoni/snake?status.svg)](https://godoc.org/github.com/jeffotoni/snake) [![Go Report](https://goreportcard.com/badge/github.com/jeffotoni/snake)](https://goreportcard.com/report/github.com/jeffotoni/snake) [![License](https://img.shields.io/github/license/jeffotoni/snake)](https://github.com/jeffotoni/snake/blob/main/LICENSE) ![GitHub last commit](https://img.shields.io/github/last-commit/jeffotoni/snake) ![GitHub contributors](https://img.shields.io/github/contributors/jeffotoni/snake) ![GitHub forks](https://img.shields.io/github/forks/jeffotoni/snake?style=social) ![GitHub stars](https://img.shields.io/github/stars/jeffotoni/snake?style=social)

A classic Snake game built with **Go** and **Ebitengine**, designed to run as a native desktop game or directly in the browser through **WebAssembly**.

This repository keeps only the final, clean version of the project while preserving the development journey as context.

---

## About

Snake v7 is a compact game project focused on the core mechanics that make Snake fun:

- Grid-based movement
- Food collection and score tracking
- Snake growth
- Speed progression
- Wall and self-collision detection
- Game Over state with restart
- Embedded textures, font, and audio
- Desktop and browser execution

---

## Snake in Action

![Snake gameplay preview](./snake.jpeg)

> The current preview uses `snake.jpeg`. It can be replaced later with an animated GIF using the same filename or by updating this README image path.

---

## Requirements

- Go 1.22 or newer
- A modern browser with WebAssembly support for the web version

---

## Quick Start

```bash
git clone https://github.com/jeffotoni/snake.git
cd snake
go run .
```

---

## Controls

| Action | Key |
|---|---|
| Move up | Arrow Up |
| Move down | Arrow Down |
| Move left | Arrow Left |
| Move right | Arrow Right |
| Restart after Game Over | R |

---

## Run on Desktop

```bash
go run .
```

Build a local binary:

```bash
go build -o snake .
./snake
```

---

## Run in the Browser with WebAssembly

### Option A: Build with Go and serve static files

```bash
GOOS=js GOARCH=wasm go build -o web/snake.wasm .
cd web
python3 -m http.server 8080
```

Open:

```text
http://localhost:8080
```

### Option B: Use wasmserve

```bash
go run github.com/hajimehoshi/wasmserve@latest .
```

Open:

```text
http://localhost:8080/web/
```

Important:

- Run `wasmserve` from the project root.
- Do not run `wasmserve ./web`, because `web/` does not contain the Go entrypoint.
- `web/index.html` can load both `snake.wasm` and `/main.wasm`, so it works with both options above.

---

## Browser Audio

Browsers can block game audio until the page receives user interaction.

If the game starts without sound:

- Click the game page or canvas once.
- Press any key, such as an arrow key.
- The message `Audio locked by browser...` disappears when audio is ready.

---

## Project Layout

```text
.
├── assets/
│   ├── PressStart2P-Regular.ttf
│   ├── eating.mp3
│   ├── food.png
│   └── gameover.mp3
├── web/
│   ├── index.html
│   └── wasm_exec.js
├── main.go
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

---

## Troubleshooting

### `expected magic word ... found 34 30 34 20`

This means the browser received a `404` HTML response instead of a real `.wasm` file.

Fixes:

- With Option A, make sure `web/snake.wasm` exists after running the build command.
- With Option B, run `go run github.com/hajimehoshi/wasmserve@latest .` from the project root and open `/web/`.

### No sound in the browser

This is usually caused by browser autoplay policy.

Fix:

- Click the page and press a key once to unlock audio.

---

## Development Journey

| Version | Focus |
|---|---|
| Version 0 | Basic game board and snake rendering |
| Version 1 | Keyboard movement and direction control |
| Version 2 | Collision rules for walls and snake body |
| Version 3 | WebAssembly build support |
| Version 4 | Texture support for food |
| Version 5 | Audio effects integration |
| Version 6 | Game Over flow and restart |
| Version 7 | Current clean version with desktop and web structure |

---

## Contributing

Contributions are welcome.

Suggested ways to improve the project:

- Add tests for game rules and movement behavior.
- Improve responsive behavior in the web wrapper.
- Add new visual assets or polish the game screen.
- Open an issue with ideas or bugs.

---

## License

This project is open source under the **MIT License**.

Copyright (c) 2026 Jefferson Otoni Lima.

See [LICENSE](./LICENSE) for the full license text.
