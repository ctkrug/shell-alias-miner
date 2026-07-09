// Loads the compiled wasm module and wires the file picker to the
// mineHistory() function it exposes. All relative paths so this works
// mounted at any base path.

// filterProposals hides every proposal seen fewer than minOccurrences times
// or whose KeystrokesSaved is below minSavings. Both thresholds must be
// cleared for a proposal to remain; it never re-mines, just re-filters the
// already-mined list.
function filterProposals(proposals, minOccurrences, minSavings) {
  return proposals.filter(
    (p) => p.Occurrences >= minOccurrences && p.KeystrokesSaved >= minSavings
  );
}

// explainKeystrokesSaved renders the row-specific formula text for the
// Keystrokes Saved info affordance. For a function proposal the fixed
// portion is the prefix baked into Definition (the "$1" argument itself
// costs the same to type either way), recovered from the definition text
// since Proposal doesn't carry Prefix directly.
function explainKeystrokesSaved(p) {
  let fixedPortion = p.Command;
  if (p.Kind === "function") {
    const match = p.Definition.match(/^function \S+\(\) \{ (.+) "\$1"; \}$/);
    if (match) {
      fixedPortion = match[1];
    }
  }
  return (
    "(" +
    fixedPortion.length +
    ' chars in "' +
    fixedPortion +
    '" − ' +
    p.Name.length +
    ' chars in "' +
    p.Name +
    '") × ' +
    p.Occurrences +
    " uses = " +
    p.KeystrokesSaved
  );
}

// Node's `require`-based test runner can load this file without a DOM; the
// browser never sets `module`, so this export is a no-op there.
if (typeof module !== "undefined") {
  module.exports = { filterProposals, explainKeystrokesSaved };
}

// Everything below drives the live page and needs a DOM, so it's skipped
// when this file is `require`d from a non-browser test runner.
if (typeof document !== "undefined") {
  const statusEl = document.getElementById("status");
  const fileInput = document.getElementById("history-file");
  const thresholdsSection = document.getElementById("thresholds");
  const minOccurrencesInput = document.getElementById("min-occurrences");
  const minSavingsInput = document.getElementById("min-savings");
  const resultsTable = document.getElementById("results");
  const resultsBody = document.getElementById("results-body");

  let minedProposals = [];

  const go = new Go();
  let wasmReady = WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
    .then((result) => {
      go.run(result.instance);
      statusEl.textContent = "Ready. Choose a history file to get started.";
    })
    .catch((err) => {
      statusEl.textContent = "Failed to load the wasm module: " + err.message;
    });

  fileInput.addEventListener("change", async () => {
    const file = fileInput.files[0];
    if (!file) {
      return;
    }

    statusEl.textContent = "Mining " + file.name + "...";
    await wasmReady;

    const text = await file.text();
    const raw = window.mineHistory(text);
    const parsed = JSON.parse(raw);

    if (parsed.error) {
      statusEl.textContent = "Error: " + parsed.error;
      thresholdsSection.hidden = true;
      resultsTable.hidden = true;
      return;
    }

    minedProposals = parsed;
    applyFilterAndRender();
  });

  minOccurrencesInput.addEventListener("input", applyFilterAndRender);
  minSavingsInput.addEventListener("input", applyFilterAndRender);

  function applyFilterAndRender() {
    if (minedProposals.length === 0) {
      thresholdsSection.hidden = true;
      statusEl.textContent = "No repeated commands found in that history file.";
      resultsTable.hidden = true;
      return;
    }

    thresholdsSection.hidden = false;
    const minOccurrences = Number(minOccurrencesInput.value) || 0;
    const minSavings = Number(minSavingsInput.value) || 0;
    const filtered = filterProposals(minedProposals, minOccurrences, minSavings);
    renderResults(filtered, minedProposals.length);
  }

  function renderResults(proposals, totalMined) {
    resultsBody.innerHTML = "";

    if (proposals.length === 0) {
      resultsTable.hidden = true;
      statusEl.textContent = "No candidates clear the current thresholds.";
      return;
    }

    for (const p of proposals) {
      resultsBody.appendChild(buildRow(p));
    }

    resultsTable.hidden = false;
    statusEl.textContent =
      "Showing " + proposals.length + " of " + totalMined + " candidate" +
      (totalMined === 1 ? "" : "s") + ".";
  }

  function buildRow(p) {
    const row = document.createElement("tr");

    const alias = document.createElement("td");
    alias.dataset.label = "Alias";
    alias.textContent = p.Name;

    const type = document.createElement("td");
    type.dataset.label = "Type";
    type.textContent = p.Kind;

    const definition = document.createElement("td");
    definition.dataset.label = "Definition";
    const code = document.createElement("code");
    code.textContent = p.Definition;
    definition.appendChild(code);

    const seen = document.createElement("td");
    seen.dataset.label = "Seen";
    seen.textContent = p.Occurrences;

    const saved = document.createElement("td");
    saved.dataset.label = "Keystrokes saved";
    saved.appendChild(buildSavedExplainer(p));

    const actions = document.createElement("td");
    actions.dataset.label = "Actions";
    actions.appendChild(buildCopyButton(p.Definition));

    row.append(alias, type, definition, seen, saved, actions);
    return row;
  }

  function buildSavedExplainer(p) {
    const details = document.createElement("details");
    details.className = "explain";

    const summary = document.createElement("summary");
    summary.textContent = p.KeystrokesSaved;
    summary.title = explainKeystrokesSaved(p);

    const formula = document.createElement("span");
    formula.textContent = explainKeystrokesSaved(p);

    details.append(summary, formula);
    return details;
  }

  function buildCopyButton(definition) {
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = "Copy";
    button.addEventListener("click", () => {
      const revert = () => {
        button.textContent = "Copy";
        button.classList.remove("copied");
      };
      navigator.clipboard
        .writeText(definition)
        .then(() => {
          button.textContent = "Copied";
          button.classList.add("copied");
          setTimeout(revert, 1500);
        })
        .catch(() => {
          button.textContent = "Copy failed";
          setTimeout(revert, 1500);
        });
    });
    return button;
  }
}
