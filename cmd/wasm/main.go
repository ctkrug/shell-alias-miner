// Command wasm compiles the mining pipeline to WebAssembly and exposes it
// as a single JavaScript-callable function, mineHistory(text) -> JSON.
//
// Everything runs inside the browser's wasm sandbox: history text goes in,
// a JSON array of alias proposals comes out, and nothing ever leaves the
// machine.
package main

import (
	"encoding/json"
	"strings"
	"syscall/js"

	"github.com/ctkrug/shell-alias-miner/internal/alias"
	"github.com/ctkrug/shell-alias-miner/internal/history"
	"github.com/ctkrug/shell-alias-miner/internal/miner"
)

func main() {
	js.Global().Set("mineHistory", js.FuncOf(mineHistory))

	// Block forever: a wasm program that returns exits the runtime and
	// mineHistory would stop being callable.
	select {}
}

func mineHistory(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return errorJSON("mineHistory requires the history text as its argument")
	}

	text := args[0].String()
	commands := history.Parse(strings.NewReader(text))
	candidates := miner.CountFrequencies(commands)
	proposals := alias.Propose(candidates)

	out, err := json.Marshal(proposals)
	if err != nil {
		return errorJSON(err.Error())
	}
	return string(out)
}

func errorJSON(message string) string {
	out, _ := json.Marshal(map[string]string{"error": message})
	return string(out)
}
