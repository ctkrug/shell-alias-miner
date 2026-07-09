# Design

The visual and interaction direction for `site/`, decided once here so BUILD/QA
craft work has a fixed target instead of drifting toward generic defaults.

## 1. Aesthetic direction

**Terminal-mono.** Sift reads someone's shell history, so the page
itself looks like a terminal window: a dark, glowing command console with a
monospace data grid, a blinking cursor in the wordmark, and a faux window
titlebar around the mining workspace. It should feel like a tool built by and
for people who live in a shell, not a generic SaaS landing page wrapped around
a file input.

## 2. Tokens

| Token | Value | Use |
|---|---|---|
| `--bg` | `#0a0d12` | page background |
| `--surface-1` | `#11151d` | terminal panel body |
| `--surface-2` | `#171d29` | raised rows / inputs / hover |
| `--border` | `rgba(255,255,255,0.09)` | hairlines |
| `--text` | `#e7ebf3` | primary text |
| `--text-muted` | `#8a93a8` | secondary text, labels |
| `--accent` | `#3ee08c` | terminal green — primary actions, focus, cursor |
| `--accent-dim` | `rgba(62,224,140,0.14)` | accent fills / glows |
| `--support` | `#5ec8ff` | cyan — links, the Type=function tag |
| `--danger` | `#ff6b6b` | error state |

- Type pairing: **JetBrains Mono** (display — wordmark, headings, all tabular/
  code content) + **Inter** (UI — body copy, labels, buttons). Both from
  Google Fonts with system-monospace / system-sans fallbacks.
- Spacing unit: 4px scale (4/8/12/16/24/32/48/64).
- Corner radius: 10px for the terminal panel, 6px for controls.
- Shadow/glow: soft `--accent-dim` glow on focus and on the panel's top edge;
  a layered `0 20px 60px rgba(0,0,0,.45)` drop shadow under the panel for
  depth against the page background.
- Motion: UI transitions 150ms ease-out; button press 90ms; cursor blink
  1.1s step-end, disabled entirely under `prefers-reduced-motion`.

## 3. Layout intent

The hero **is** the terminal window: a titlebar (three dots, a path-style
label `~/sift`) over the drop zone / thresholds / results table,
composed as one continuous panel. On 1440×900 it's centered, ~1040px wide,
comfortably the dominant element on the page with room to breathe around it
(not full-bleed — a terminal window has edges). On 390×844 the panel goes
full-width edge-to-edge like a real mobile terminal app, the titlebar dots
shrink, and table rows collapse to a stacked card layout so nothing scrolls
horizontally.

## 4. Signature detail

The wordmark renders as a shell prompt: `sift█` with a block
cursor that blinks (CSS `steps()` animation, paused under reduced-motion).
The panel titlebar's three dots are real (red/amber/green) window-chrome
dots, reinforcing the terminal illusion without being a skeuomorphic photo.

## 5. Juice plan

Not a game; no SFX plan needed. The equivalent "feel" budget: the Copy
button's confirmation state, the min-occurrence/min-savings inputs
re-filtering instantly with a subtle row fade-in, and the accent glow that
sweeps in on focus — all under 250ms so the tool feels responsive to a
power-user audience that will judge it by feel.
