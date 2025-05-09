package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// We're emulating the clipboard by writing to a file called clipboard

var CopyToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy a branch name to the clipboard using custom clipboard command template",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},

	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
	},

	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a").IsSelected(),
			).
			Press(keys.Universal.CopyToClipboard)

		t.ExpectToast(Equals("'branch-a' copied to clipboard"))

		t.Views().Files().
			Focus()

		t.GlobalPress(keys.Files.RefreshFiles)

		t.FileSystem().FileContent("clipboard", Equals("branch-a"))
	},
})
