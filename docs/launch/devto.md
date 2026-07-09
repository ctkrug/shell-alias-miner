---
title: "I mined my shell history for the aliases I never got around to writing"
published: false
tags: go, webassembly, cli, productivity
---

I have typed `git status --short --branch` thousands of times. I know I should
have aliased it years ago. I never did, because the chore of aliasing is not
typing the alias line, it is sitting down to figure out *which* of my long
commands are worth it and what each one should be called. Shell history is not
something you read.

So I built a tool that reads it for me. It is called **Sift**. You drop your
`.zsh_history` or `.bash_history` onto the page and it hands you back
ready-to-paste aliases and shell functions, ranked by how many keystrokes each
one saves. It is live here: <https://apps.charliekrug.com/shell-alias-miner/>,
and the source is here: <https://github.com/ctkrug/shell-alias-miner>.

Here are the two decisions I found most interesting to make.

## History never leaves the machine, so there is no server

Shell history is personal. It has your paths, your hostnames, flags that hint
at the infrastructure you touch. I did not want a tool that asks you to upload
that, and I did not want to run a server that could. So the whole miner is a
Go program compiled to `GOOS=js GOARCH=wasm`, and the page is a static bundle.
Your file is read with the browser File API and passed straight into the wasm
module. There is no fetch of your data, because there is nowhere to fetch it
to. You can watch the network tab while you use it and see nothing move.

A nice side effect: the parsing and mining logic is plain Go with no
`syscall/js` dependency, so it is all unit-testable and fuzz-testable on the
host. The wasm entrypoint is a thin wrapper that marshals proposals to JSON.
The tests never touch a browser.

## The mining is n-gram templates, not `sort | uniq -c`

Counting identical lines is the easy 80%, and on its own it misses the
commands that actually dominate your history: the ones with a varying
argument. Fifty `git commit -m "..."` lines with fifty different messages look
like fifty unique commands to a naive counter.

Sift tokenizes each command (quote-aware, so a commit message does not get
split on its spaces), groups by every token except the last, and only treats a
group as a template when it sees at least two distinct trailing arguments. A
group with one observed suffix has no variation to catch, so it stays a plain
exact-match candidate. That one rule means history with no argument variation
mines identically to the simple counter, and history full of it collapses into
the handful of real patterns.

Templates become shell functions (`function gc() { git commit -m "$1"; }`)
because a plain alias cannot take an argument. Everything else becomes an
alias. And the ranking is by keystrokes saved, not frequency: a six-character
command you ran 500 times is not worth aliasing, but an eighty-character one
you ran twenty times is.

## What bit me, and what I would do differently

Two things I am glad I caught. Reading history with `bufio.Scanner` silently
stops at its buffer cap the moment one line is too long, dropping the rest of
the file; I switched to `bufio.Reader.ReadString`, which has no ceiling. And a
first pass at deriving alias names happily pulled the first byte off tokens
like `&&` or a quoted flag, leaking shell metacharacters into names; now it
takes the first alphanumeric rune or skips the token. A property-based fuzz
test guards both.

If I did it again I would add per-shell export helpers (a one-liner to grab
your history in the right place), and maybe a diff against your existing
`.zshrc` so it only proposes aliases you do not already have.

If you live in a terminal, point it at your history and see what falls out.
Feedback welcome, especially on commands it should catch but does not.
</content>
