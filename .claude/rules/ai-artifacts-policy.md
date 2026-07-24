# 生成物 — 手編集しないこと

| 生成物 | 原本 |
| --- | --- |
| `internal/llmdocs/90-commands.md` | `cmd/*.go` の cobra コマンド定義（`internal/gen-llmdocs` が描画） |

手書き（編集してよい）:

- `internal/llmdocs/00-guide.md` — 認証・環境・グローバルフラグ・出力・書き込み時の注意
- `internal/llmdocs/10-billings.md` — インボイス制度による請求取得系統の分岐
- `plugins/mfk-cli/skills/*/SKILL.md`
- `context7.json`

生成物を直接編集しても次の `go generate ./...` で消え、それまでは CI が古い差分で
落ちる。カタログを良くしたいときはコマンドの `Short` / `Long` / フラグ説明を直すこと。

再生成は `/regen-ai`、または `go generate ./... && go test ./...`。
