# Vision

## The problem

Everyone who lives in a terminal has a handful of commands they type
dozens or hundreds of times: a `git log` incantation with the same six
flags, a `docker run` with the same mount and port mapping, an `ffmpeg`
transcode, a `kubectl` command against the same namespace. Everyone knows,
in the abstract, that they *should* alias these. Almost nobody does,
because figuring out *which* commands are worth it and *what the alias
should look like* is its own small chore, and shell history is not
something you sit down and read.

History search tools (`Ctrl+R`, `fzf`, `atuin`, `mcfly`) solve a different
problem: finding a command you already know you ran, right now, so you can
run it again. They don't look backward across your whole history and tell
you "you've paid a 40-character tax on this exact shape of command 340
times this year — here's the alias that ends it."

## Who it's for

Anyone with a shell history file worth mining: developers, sysadmins,
data folks — anyone who has been living in the same terminal long enough
to have built up repeated habits without automating them. No install
beyond opening a web page; no account; no dependency on which shell
plugin manager they use.

## The core idea

Treat shell history as a corpus to mine, not a list to search:

1. **Parse** the history file into individual command lines, handling the
   real-world formats (plain `bash_history`, zsh's `EXTENDED_HISTORY`
   timestamps, line continuations).
2. **Mine** it for repeated patterns — starting with exact-line frequency,
   extending to n-gram template matching so a command that varies only in
   its trailing argument (a commit message, a filename, a container tag)
   is still recognized as "the same command" rather than counted as
   hundreds of unique lines.
3. **Propose** a ready-to-paste `alias` (for fixed commands) or shell
   `function` (for commands with a varying argument), plus a concrete
   keystrokes-saved estimate so the value is never a hand-wave.

All of this must run **locally, in the browser, in well under a second**
for a realistic history file (tens of thousands of lines). History is
personal and often sensitive (paths, hostnames, flags that hint at
infrastructure) — it should never leave the machine, which is why the
whole pipeline is a Go program compiled to WebAssembly rather than a
service with an upload step.

## Key design decisions

- **Go compiled to WASM, no backend.** History never crosses the network.
  This is also just simpler to host: a static bundle, deployable to any
  subpath, no server to run or pay for.
- **N-gram mining, not exact-line counting, is the differentiator.**
  Exact-line counting is the easy 80% and ships first (it's still useful
  on its own), but the wow moment — and the reason this beats "just grep
  your history and sort | uniq -c" — is catching commands that vary by
  argument. That's tracked as its own epic, not bolted on later.
  Everything else routes on top of `internal/miner`, so upgrading exact
  matching to template mining doesn't touch the parser or the proposer.
- **Aliases for fixed commands, functions for parameterized ones.** A
  plain `alias` can't take an argument; proposing an `alias` for a command
  that actually varies would hand the user something broken. The miner's
  job is to correctly classify which candidates are which before the
  proposer ever runs.
- **Keystrokes saved is the ranking metric, not raw frequency.** A command
  typed 500 times that's already only 6 characters isn't worth aliasing.
  A command typed 20 times that's 80 characters is. Ranking by savings
  keeps the results useful instead of just a list of "your most common
  short commands."
- **No config file, no account, no persistence.** Point it at a file, get
  an answer. State lives in the browser tab for the session and nowhere
  else.

## What "v1 done" looks like

- Drop a real `.zsh_history` or `.bash_history` file onto the page and,
  within about a second, see a ranked table of alias/function proposals
  with occurrence counts and keystrokes-saved numbers.
- The miner catches both exact-repeat commands and commands that vary by
  a trailing argument (the n-gram template case) — not just `sort | uniq
  -c`-equivalent exact matching.
- Every proposal is copy-pasteable as-is into `.zshrc`/`.bashrc` and is
  syntactically valid shell.
- Adjustable minimum-occurrence and minimum-savings thresholds so a huge
  history file doesn't drown the user in noise.
- The whole thing is a static site: open `index.html` locally or host it
  at any subpath, no build step or server required at runtime.
