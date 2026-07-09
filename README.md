# Sift

**▶ Live demo — [apps.charliekrug.com/shell-alias-miner](https://apps.charliekrug.com/shell-alias-miner/)**

[![CI](https://github.com/ctkrug/shell-alias-miner/actions/workflows/ci.yml/badge.svg)](https://github.com/ctkrug/shell-alias-miner/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> The aliases hiding in your shell history.

Point Sift at your own shell history and it mines your most-repeated
multi-word commands, proposes ready-to-paste aliases and shell functions for
them, and tells you exactly how many keystrokes each one will save you.

History search tools help you find a command again. Sift closes the loop: it
finds the commands you keep retyping and hands you the alias so you never have
to.

Everything runs **locally in your browser**. The miner is a Go program
compiled to WebAssembly, so your shell history is never uploaded anywhere.

## Why

Most of us have a handful of long, awkward commands we type dozens or hundreds
of times a year: a `git log` incantation, a `docker run` with the same six
flags, an `ffmpeg` transcode. We know we *should* alias them. We never get
around to figuring out which ones are worth it, or what the alias should even
look like. Sift does that analysis for you in under a second.

## The wow moment

Drop your `.zsh_history` (or `.bash_history`) onto the page. Within a second
it shows you the exact command you typed 340 times this year, a ready-to-paste
`alias` line, and the math on how many keystrokes that alias will save you
going forward.

## How it works

Sift doesn't just count identical lines. It tokenizes each history entry and
mines repeated **n-gram** patterns across commands, so it catches a command
you run constantly with varying arguments (for example `git commit -m "..."`
with a different message every time) as a single candidate, not hundreds of
"unique" lines.

For every candidate pattern it proposes either:

- a plain `alias` when the command has no varying arguments, or
- a shell `function` with a positional parameter when arguments vary.

Each proposal comes with a keystrokes-saved estimate:
`(chars typed per use - chars typed with the alias) x times seen`.

## Sample output

| Alias | Type | Definition | Seen | Keystrokes saved |
|---|---|---|---|---|
| `gs` | alias | `alias gs="git status --short --branch"` | 340 | 9,860 |
| `gc` | function | `function gc() { git commit -m "$1"; }` | 128 | 1,664 |
| `dcu` | alias | `alias dcu="docker compose up -d"` | 58 | 1,392 |
| `gl` | alias | `alias gl="git log --oneline --graph --all"` | 21 | 861 |

## Features

- Parses `.zsh_history` (with and without the `EXTENDED_HISTORY` timestamp
  format) and plain `.bash_history`.
- Drag-and-drop the file onto the page, or use the file picker.
- Frequency and n-gram mining over tokenized command lines, not raw string
  matching. Commands that vary only in a trailing argument (a commit message,
  a filename) are recognized as one repeated pattern.
- Ranks candidates by total keystrokes saved, not just raw frequency.
- Proposes `alias` for fixed commands and `function` for commands with varying
  arguments; a "Type" column shows which is which.
- Adjustable minimum-occurrence and minimum-savings thresholds, applied
  instantly without re-mining.
- One-click copy of the generated alias/function block, with an info
  affordance that explains the keystrokes-saved math per row.
- Never proposes an alias/function that would bake in a password, token, or
  other inline credential.
- Everything runs client-side via WebAssembly. No server, no upload, no
  network request.

## Stack

- **Go**, compiled to `GOOS=js GOARCH=wasm`, is the mining engine and parser.
- A minimal static HTML/CSS/JS shell loads the wasm module and drives the UI
  (file picker, results table, copy-to-clipboard).
- No backend and no build server required at runtime; the wasm binary is a
  static asset.

## Building and running locally

```sh
make test    # run the Go test suite
make test-js # run site/main.js's unit tests (needs Node)
make vet     # go vet, including the wasm entrypoint
make site    # build site/main.wasm and copy in wasm_exec.js
```

After `make site`, serve `site/` with any static file server (opening
`index.html` directly via `file://` won't work, because `fetch()` requires
http(s)) and open it in a browser, for example:

```sh
cd site && python3 -m http.server 8080
```

## Documentation

- [`docs/VISION.md`](docs/VISION.md) for the full design rationale.
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) for how the pieces fit
  together.
- [`docs/DESIGN.md`](docs/DESIGN.md) for the site's visual direction.
- [`docs/BACKLOG.md`](docs/BACKLOG.md) for what's shipped and what's left.

## License

MIT, see [`LICENSE`](LICENSE).

---

More of Charlie's projects → [apps.charliekrug.com](https://apps.charliekrug.com)
</content>
