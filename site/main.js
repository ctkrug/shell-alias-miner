// Loads the compiled wasm module and wires the file picker to the
// mineHistory() function it exposes. All relative paths so this works
// mounted at any base path.

const statusEl = document.getElementById("status");
const fileInput = document.getElementById("history-file");
const resultsTable = document.getElementById("results");
const resultsBody = document.getElementById("results-body");

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
    resultsTable.hidden = true;
    return;
  }

  renderResults(parsed);
});

function renderResults(proposals) {
  resultsBody.innerHTML = "";

  if (proposals.length === 0) {
    statusEl.textContent = "No repeated commands found in that history file.";
    resultsTable.hidden = true;
    return;
  }

  for (const p of proposals) {
    const row = document.createElement("tr");

    const alias = document.createElement("td");
    alias.textContent = p.Name;

    const definition = document.createElement("td");
    const code = document.createElement("code");
    code.textContent = p.Definition;
    definition.appendChild(code);

    const seen = document.createElement("td");
    seen.textContent = p.Occurrences;

    const saved = document.createElement("td");
    saved.textContent = p.KeystrokesSaved;

    row.append(alias, definition, seen, saved);
    resultsBody.appendChild(row);
  }

  resultsTable.hidden = false;
  statusEl.textContent =
    "Found " + proposals.length + " candidate" + (proposals.length === 1 ? "" : "s") + ".";
}
