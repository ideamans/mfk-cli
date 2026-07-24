package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ideamans/go-llm-cli-kit/llmcmd"

	"github.com/ideamans/mfk-cli/internal/llmdocs"
)

func TestEmbeddedReference(t *testing.T) {
	g, err := llmdocs.Docs().Markdown()
	if err != nil {
		t.Fatalf("Markdown: %v", err)
	}
	for _, want := range []string{
		"AIエージェント向けリファレンス",
		"MFK_API_KEY", // 認証
		"--page-all",  // ページネーション
		"インボイス制度",     // billings 章
		"mfk billings qualified",
		"リソースとメソッド一覧", // 生成物
	} {
		if !strings.Contains(g, want) {
			t.Errorf("埋め込みリファレンスに %q がありません", want)
		}
	}
}

func TestChapterOrder(t *testing.T) {
	sections, err := llmdocs.Docs().Sections()
	if err != nil {
		t.Fatalf("Sections: %v", err)
	}
	var files []string
	for _, s := range sections {
		files = append(files, s.File)
	}
	want := "00-guide.md,10-billings.md,90-commands.md"
	if got := strings.Join(files, ","); got != want {
		t.Errorf("chapters = %s, want %s", got, want)
	}
}

// TestLegacyLLMFlag は互換の約束を守る: --llm は従来どおりコマンドラインの
// どの位置でも効く。
func TestLegacyLLMFlag(t *testing.T) {
	for _, args := range [][]string{{"--llm"}, {"billings", "list", "--llm"}} {
		var out bytes.Buffer
		handled, err := llmcmd.HandleLegacy(args, LLMConfig(), &out)
		if err != nil {
			t.Fatalf("HandleLegacy(%v): %v", args, err)
		}
		if !handled {
			t.Errorf("HandleLegacy(%v) が --llm を処理しませんでした", args)
		}
		if !strings.Contains(out.String(), "MFK_API_KEY") {
			t.Errorf("HandleLegacy(%v) の出力が想定と違います", args)
		}
	}
}
