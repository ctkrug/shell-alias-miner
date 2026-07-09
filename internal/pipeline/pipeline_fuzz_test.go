package pipeline

import "testing"

// Run composes history.Parse, miner.Mine, and alias.Propose end to end —
// this is the exact call cmd/wasm makes with whatever a user drops on the
// page. Fuzz it directly rather than relying solely on the fixed random
// sample TestRunNeverPanicsOnRandomBinaryInput takes.
func FuzzRun(f *testing.F) {
	seeds := []string{
		"",
		"git status\ngit status\ngit status\n",
		": 1700000000:0;git commit -m \"a\"\n: 1700000010:0;git commit -m \"b\"\n",
		"mysql -phunter2 mydb\nmysql -phunter2 mydb\nmysql -phunter2 mydb\n",
		"git commit -m \"🎉\"\ngit commit -m \"🎉\"\ngit commit -m \"🎉\"\n",
		"\x00\x01binary\xff\xfe",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, historyText string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Run(%q) panicked: %v", historyText, r)
			}
		}()
		Run(historyText)
	})
}
