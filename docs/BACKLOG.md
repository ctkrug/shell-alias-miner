# Backlog

Epics and stories for the v1 build. Every story lists concrete,
verifiable acceptance criteria — no "works well" vibes. Story 1.1 is the
wow moment and must land before anything else in this backlog.

## Epic 1 — The wow moment: real n-gram mining, end to end

- [x] **1.1 Load a real `.zsh_history` and see the top repeated command
  with a ready alias and its keystrokes-saved count, in under a second**
  - Given a `.zsh_history` file (10k+ lines) containing one multi-word
    command repeated 340 times with no argument variation, dropping the
    file on the page shows that command in the top 3 result rows, ranked
    by keystrokes saved, within 1 second wall-clock on a typical laptop.
  - The result row's Definition cell contains a syntactically valid
    `alias name="the exact command"` line, pasteable into `.zshrc`
    without edits.
  - The Keystrokes Saved cell shows a number greater than 0 equal to
    `(command length - alias name length) * occurrences`.

- [x] **1.2 Mine n-gram templates so commands with varying trailing
  arguments are recognized as one repeated pattern**
  - History containing `git commit -m "..."` 50 times with 50 distinct
    messages produces a single candidate for the `git commit -m`
    template, not 50 separate exact-match rows.
  - History with only exact-duplicate commands (no argument variation)
    returns the same candidates as the baseline exact-match miner — no
    regression.
  - A varying-argument candidate reports a `PatternKind` (or equivalent
    field) distinguishing it from an exact-match candidate.

- [x] **1.3 Propose shell functions (not aliases) for commands with a
  varying trailing argument**
  - A parameterized candidate's proposed definition is a
    `function name() { ... "$1" ...; }`-shaped snippet, not a plain
    `alias` (a plain alias can't take an argument, so proposing one
    would hand back a proposal that silently drops user input).
  - Sourcing the generated function definition in a real `zsh` or `bash`
    shell and calling it with an argument reproduces the original
    command with that argument substituted in the correct position.

- [ ] **1.4 Surface alias-vs-function classification in the results
  table**
  - The results table has a "Type" column showing `alias` or `function`
    per row.
  - Rows remain sorted by keystrokes saved regardless of type, so the
    single highest-value proposal — alias or function — is always first.

## Epic 2 — Trust and control over the results

- [ ] **2.1 Minimum-occurrence threshold**
  - Setting the "min occurrences" control to N hides every candidate
    seen fewer than N times in the loaded history; default is 3.
  - Changing the threshold re-filters the already-mined results
    instantly, without re-parsing or re-uploading the file.

- [ ] **2.2 Minimum-keystrokes-saved threshold**
  - Setting the "min savings" control to N hides every candidate whose
    total keystrokes saved is below N; default is 20.
  - The two thresholds (2.1, 2.2) compose: a candidate must clear both
    to appear.

- [ ] **2.3 One-click copy of a proposal**
  - Clicking "Copy" on a row copies exactly that row's Definition text
    to the clipboard (verified via the Clipboard API's write call
    receiving that row's string).
  - The button shows a brief "Copied" confirmation state after a
    successful copy, then reverts.

- [ ] **2.4 Explain the keystrokes-saved math per row**
  - An info affordance on the Keystrokes Saved column, when hovered or
    tapped, shows the formula and the specific numbers used for that
    row (command length, alias name length, occurrence count).

## Epic 3 — Robustness on real-world history files

- [ ] **3.1 Validate mixed-format parsing end to end in the UI**
  - Dropping a plain `bash_history` file (no `: <ts>:<dur>;` prefixes)
    produces commands with no leftover timestamp artifacts in any row.
  - Dropping a zsh `EXTENDED_HISTORY` file produces commands with
    timestamps correctly stripped, verified against a fixture file
    checked into `testdata/`.

- [ ] **3.2 Handle large history files without freezing the tab**
  - A synthetic 150k-line history fixture completes mining in under 3
    seconds in the browser.
  - The browser tab shows no "page unresponsive" warning during mining
    (mining runs without blocking the main thread for longer than the
    browser's hang-detection window).

- [ ] **3.3 Handle malformed or non-text input without crashing**
  - Feeding a file that isn't valid shell history (random binary bytes)
    returns zero candidates and a friendly status message in the UI
    instead of an uncaught JS exception in the console.

- [ ] **3.4 Never propose an alias that bakes in a secret**
  - Commands matching common secret-bearing patterns (e.g. `--password`,
    `--token`, a bare `-p<value>` on tools known to accept inline
    credentials) are excluded from proposals entirely.
  - A history fixture containing one such command and several unrelated
    repeated commands produces proposals for the unrelated commands but
    zero rows referencing the credential-bearing one.
