---
name: regen-ai
description: 埋め込みLLMリファレンスを再生成して検証する。コマンド・フラグ・ヘルプ文・手書き章を変更した後に使う。
allowed-tools: Bash(go generate:*) Bash(go test:*) Bash(go build:*) Bash(git status:*) Bash(git diff:*) Read
---

# regen-ai

`internal/llmdocs/` をコードに合わせ直します。

1. `git status --short` — 既存の差分を把握する。
2. `go generate ./...` — `90-commands.md` を書き直す。
3. `go build ./... && go test ./...`。
4. どのリソース・メソッドが増減したかを報告する。

生成器が見ているのはコマンドツリーだけです。請求書の取得系統やページネーションの
挙動を変えた場合は、`10-billings.md` / `00-guide.md` を手で直す必要があります。

このスキルは Claude Code ローカル用で、配布プラグインには含みません。
