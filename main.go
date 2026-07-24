package main

import (
	"fmt"
	"os"

	"github.com/ideamans/go-llm-cli-kit/llmcmd"

	"github.com/ideamans/mfk-cli/cmd"
)

func main() {
	// --llm はコマンドラインのどの位置にあってもリファレンスを表示して終了する。
	// `mfk llm` へ移行したが、既存の呼び出しを壊さないため互換を維持する。
	if handled, err := llmcmd.HandleLegacy(os.Args[1:], cmd.LLMConfig(), os.Stdout); handled {
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		return
	}

	cmd.Execute()
}
