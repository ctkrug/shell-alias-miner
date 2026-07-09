// Unit tests for main.js's pure (DOM-free) logic. Run with:
//   node --test site/main.test.js

const test = require("node:test");
const assert = require("node:assert/strict");
const { filterProposals, explainKeystrokesSaved } = require("./main.js");

function proposal(overrides) {
  return Object.assign(
    {
      Name: "gs",
      Command: "git status --short",
      Definition: 'alias gs="git status --short"',
      Occurrences: 10,
      KeystrokesSaved: 100,
      Kind: "alias",
    },
    overrides
  );
}

test("filterProposals keeps candidates that clear both thresholds", () => {
  const p = proposal({ Occurrences: 5, KeystrokesSaved: 30 });
  assert.deepEqual(filterProposals([p], 3, 20), [p]);
});

test("filterProposals drops candidates below the occurrence threshold", () => {
  const p = proposal({ Occurrences: 2, KeystrokesSaved: 1000 });
  assert.deepEqual(filterProposals([p], 3, 0), []);
});

test("filterProposals drops candidates below the savings threshold", () => {
  const p = proposal({ Occurrences: 1000, KeystrokesSaved: 5 });
  assert.deepEqual(filterProposals([p], 0, 20), []);
});

test("filterProposals composes both thresholds (AND, not OR)", () => {
  const passesOccurrenceOnly = proposal({ Occurrences: 50, KeystrokesSaved: 5 });
  const passesSavingsOnly = proposal({ Occurrences: 1, KeystrokesSaved: 500 });
  const passesBoth = proposal({ Occurrences: 50, KeystrokesSaved: 500 });

  const got = filterProposals(
    [passesOccurrenceOnly, passesSavingsOnly, passesBoth],
    3,
    20
  );

  assert.deepEqual(got, [passesBoth]);
});

test("filterProposals on an empty list returns an empty list", () => {
  assert.deepEqual(filterProposals([], 3, 20), []);
});

test("explainKeystrokesSaved reports command length, alias length, and occurrences for an alias", () => {
  const p = proposal({
    Command: "git status --short",
    Name: "gs",
    Occurrences: 10,
    KeystrokesSaved: 170,
    Kind: "alias",
  });

  const got = explainKeystrokesSaved(p);

  assert.match(got, /18 chars/); // len("git status --short")
  assert.match(got, /2 chars/); // len("gs")
  assert.match(got, /10 uses/); // occurrences
  assert.match(got, /= 170/); // KeystrokesSaved
});

test("explainKeystrokesSaved recovers the fixed prefix for a function proposal", () => {
  const p = proposal({
    Command: 'git commit -m "fix bug"',
    Name: "gc",
    Definition: 'function gc() { git commit -m "$1"; }',
    Occurrences: 50,
    KeystrokesSaved: 550,
    Kind: "function",
  });

  const got = explainKeystrokesSaved(p);

  assert.match(got, /git commit -m/);
  assert.doesNotMatch(got, /fix bug/);
});
