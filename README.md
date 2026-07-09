# Shell Alias Miner

[![CI](https://github.com/ctkrug/shell-alias-miner/actions/workflows/ci.yml/badge.svg)](https://github.com/ctkrug/shell-alias-miner/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Point it at your own shell history and it mines your most-repeated multi-word
commands, proposes ready-to-paste aliases and shell functions for them, and
tells you exactly how many keystrokes each one will save you.

History search tools help you find a command again. Shell Alias Miner closes
the loop: it finds the commands you keep retyping and hands you the alias so
you never have to.

Everything runs **locally in your browser** — the miner is a Go program
compiled to WebAssembly. Your shell history is never uploaded anywhere.

## Why

Most of us have a handful of long, awkward commands we type dozens or
hundreds of times a year — a `git log` incantation, a `docker run` with the
same six flags, an `ffmpeg` transcode. We know we *should* alias them. We
never get around to figuring out which ones are worth it, or what the alias
should even look like. Shell Alias Miner does that analysis for you in under
a second.

## The wow moment

Drop your `.zsh_history` (or `.bash_history`) onto the page. Within a
second it shows you the exact command you typed 340 times this year, a
ready-to-paste `alias` line, and the math on how many keystrokes that alias
will save you going forward.

## How it works

Shell Alias Miner doesn't just count identical lines — it tokenizes each
history entry and mines repeated **n-gram** patterns across commands, so it
catches a command you run constantly with varying arguments (for example
`git commit -m "..."` with a different message every time) as a single
candidate, not hundreds of "unique" lines.

For every candidate pattern it proposes either:

- a plain `alias` when the command has no varying arguments, or
- a shell `function` with a positional parameter when arguments vary.

Each proposal comes with a keystrokes-saved estimate: `(chars typed per use -
chars typed with the alias) x times seen`.

## Planned features

- Parse `.zsh_history` (with and without the `EXTENDED_HISTORY` timestamp
  format) and plain `.bash_history`.
- Frequency + n-gram mining over tokenized command lines, not raw string
  matching.
- Rank candidates by total keystrokes saved, not just raw frequency.
- Propose `alias` for fixed commands and `function` for commands with
  varying arguments.
- Adjustable minimum-occurrence and minimum-savings thresholds.
- One-click copy of the generated alias/function block.
- Everything runs client-side via WebAssembly — no server, no upload, no
  network request.

## Stack

- **Go**, compiled to `GOOS=js GOARCH=wasm` — the mining engine and parser.
- A minimal static HTML/CSS/JS shell that loads the wasm module and drives
  the UI (file picker, results table, copy-to-clipboard).
- No backend, no build server required at runtime — the wasm binary is a
  static asset.

## Building and running locally

```sh
make test    # run the Go test suite
make vet     # go vet, including the wasm entrypoint
make site    # build site/main.wasm and copy in wasm_exec.js
```

After `make site`, serve `site/` with any static file server (opening
`index.html` directly via `file://` won't work — `fetch()` requires
http(s)) and open it in a browser, e.g.:

```sh
cd site && python3 -m http.server 8080
```

## Status

Early scaffold. See [`docs/VISION.md`](docs/VISION.md) for the full design
and [`docs/BACKLOG.md`](docs/BACKLOG.md) for the build plan.

## License

MIT — see [`LICENSE`](LICENSE).
