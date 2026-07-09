# Architecture

A concise map of the codebase for anyone (human or model) picking this up
fresh. See `VISION.md` for *why*, `BACKLOG.md` for *what's left*.

## Data flow

```
history text
    │
    ▼
internal/history.Parse        strip zsh EXTENDED_HISTORY timestamps,
                               join line continuations, skip blanks
    │  []string (one command per line)
    ▼
internal/miner.Mine           dedupe + count; group by every token except
                               the last; a group with ≥2 distinct trailing
                               arguments becomes one KindTemplate candidate,
                               everything else is KindExact
    │  []Candidate
    ▼
internal/alias.Propose        drop candidates matching containsSecret;
                               turn each survivor into a ready-to-paste
                               alias (KindExact) or function (KindTemplate)
                               with a keystrokes-saved estimate; resolve
                               name collisions
    │  []Proposal
    ▼
cmd/wasm.mineHistory          JSON-marshal []Proposal, exposed to JS as
                               the global mineHistory(text) function
    │  JSON string
    ▼
site/main.js                  fetch/instantiate the wasm module, wire the
                               file picker, filter by the threshold
                               controls, render the results table
```

`internal/pipeline.Run(historyText string) []alias.Proposal` composes the
first three stages and has no `syscall/js` dependency, so it's what
`internal/pipeline/*_test.go` exercises directly — `cmd/wasm` is a thin
wrapper that only exists because `internal/miner`/`internal/alias` can't be
unit-tested under `GOOS=js`.

## Modules

- **`internal/history`** — `Parse(io.Reader) []string`. Handles plain
  `bash_history`, zsh `EXTENDED_HISTORY` (`: <ts>:<dur>;cmd`), and trailing
  `\` line continuations. Never errors; best-effort on malformed input.
- **`internal/miner`** — `Candidate{Command, Prefix, Count, Kind}`.
  - `tokenize(string) []string`: quote-aware word splitter (a fast path via
    `strings.Fields` when there are no quotes, a rune-by-rune quote tracker
    otherwise). Shared by `Mine`.
  - `CountFrequencies`: the original exact-match-only miner. Still used
    directly by its own tests as the "no n-gram variation" baseline that
    `Mine` must match when history has no argument variation.
  - `Mine`: the n-gram template miner described above. See the comment on
    `Mine` in `template.go` for the covered-set optimization (a flat
    `(command, prefix)` slice instead of a set keyed by full command text —
    matters at 100k+-line scale, see `bench_test.go`).
- **`internal/alias`** — `Propose([]miner.Candidate) []Proposal`.
  - `containsSecret` (`secrets.go`): regex match on common long-form
    credential flags (`--password`, `--token`, `--api-key`, `--secret`) plus
    an allowlist-gated check for `mysql`/`mongo`'s inline `-p<password>`
    convention. Candidates matching it are dropped before proposing.
  - `proposeAlias` / `proposeFunction`: build the `alias name="..."` or
    `function name() { prefix "$1"; }` snippet and the keystrokes-saved
    number. `uniqueName` resolves name collisions in O(1) amortized via a
    per-base "next suffix" counter (not a re-probe-from-2 scan).
- **`internal/pipeline`** — `Run(string) []alias.Proposal`, the composed
  pipeline; also home to `bench_test.go` (10k/100k realistic-mix
  benchmarks).
- **`cmd/wasm`** — `GOOS=js GOARCH=wasm` entrypoint. Registers
  `mineHistory(text) -> JSON string` on `js.Global()` and blocks forever
  (`select{}`) so the wasm runtime doesn't exit and unregister the function.
- **`site/`** — static HTML/CSS/JS, no build step beyond `make site`
  (compiles the wasm binary and copies in `wasm_exec.js`). All paths are
  relative so it can be hosted at any subpath.
  - `main.js`: loads the wasm module, wires the file picker, and holds two
    small pure functions — `filterProposals` (AND-composes the
    min-occurrences/min-savings thresholds) and `explainKeystrokesSaved`
    (renders the per-row formula text, recovering a function proposal's
    fixed prefix from its `Definition` string since `Proposal` doesn't
    carry `Prefix`). Both are exported under `module.exports` when
    `require`d from Node so `main.test.js` can unit-test them without a
    DOM; that export is a no-op in the browser.

## Tests

- Go: `make test` (`go test ./internal/...`). Every package has unit tests;
  `internal/alias` and `internal/pipeline` also include real-execution/
  integration-style tests (sourcing a generated function in an actual
  `bash` shell; running the full pipeline against the `testdata/` fixtures).
- JS: `make test-js` (`node --test site/main.test.js`) — the DOM-free logic
  in `main.js` only; DOM-driving behavior (drag-drop, copy button, explain
  toggle) is verified by hand in a real browser (headless Chromium via
  Playwright) during development, not by a checked-in browser test suite.
- Both suites run in CI (`.github/workflows/ci.yml`).

## Running things

```sh
make test              # Go unit tests
make test-js           # site/main.js unit tests (needs Node)
make vet                # go vet, including the wasm entrypoint
make fmt                # gofmt check (non-zero exit if unformatted)
make site               # builds site/main.wasm + copies in wasm_exec.js
cd site && python3 -m http.server 8080   # serve it (file:// won't work)
```
