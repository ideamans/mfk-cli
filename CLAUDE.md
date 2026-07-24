# CLAUDE.md — mfk-cli

マネーフォワードケッサイ API v2 の CLI。**バイナリ名は `mfk`**（リポジトリ名と
リリースアーカイブ名は `mfk-cli`、`go install` した場合も `mfk-cli`）。

実際の請求・与信・支払データを扱うため、書き込み系の変更は特に慎重に。

## 変更時の必須手順

**機能を追加した、フラグを増やした、既存の挙動を変えた — このいずれかをしたら、
3か所すべてを更新してから終わること。**

| 更新先 | 対象 | やり方 |
| --- | --- | --- |
| ① ドキュメント | `README.md` | 使い方が変わったときのみ |
| ② ヘルプ | cobra の `Short` / `Long` / フラグ説明 | コード内。**カタログはここから生成される** |
| ③ **LLMナレッジ** | `internal/llmdocs/00-guide.md` | 認証・環境・グローバルフラグ・出力形式が変わったら |
| | `internal/llmdocs/10-billings.md` | **請求書の取得系統（qualified / list）を変えたら必ず** |
| | `internal/llmdocs/90-commands.md` | **生成物。手編集しない** → `go generate ./...` |
| | `plugins/mfk-cli/skills/*/SKILL.md` | 手順や前提が変わったとき |
| | `context7.json` の `rules` | 新しい落とし穴が生まれたとき |

③ を忘れやすい。ドキュメントとヘルプは人間が読んで気づくが、**LLMナレッジが
古いことには誰も気づかない**（エージェントが黙って間違えるだけ）。

判断に迷ったときの目安:

- **請求書の取得系統に関わる変更** → `10-billings.md` は必須。インボイス制度で
  `qualified` と `list` に分かれている件はこの CLI 最大の落とし穴で、
  「2023年までのデータしか取れない」という問い合わせの原因が常にこれ
- ページネーションの既定値を変えた → `00-guide.md` と `context7.json`。
  既定20件のまま `--page-all` を付け忘れると集計が静かに間違う
- 書き込み系メソッドを追加した → `mfk-usage` の SKILL.md。`--dry-run` での
  事前確認とユーザー同意の手順を保つこと
- 新しいリソースを足した → ②を書いてから `go generate ./...`

## リリース

`PluginVersion`（`cmd/root.go`）と `plugin.json` の `version` と git タグの3つを
揃える。テストとリリースワークフローが不一致を検出する。手順は
`plugins/mfk-cli/PUBLISH.md`。

## 確認

```bash
go generate ./...     # 生成物を作り直す
git diff --exit-code  # 差分が出たらコミット漏れ
go test ./...         # SKILL.md 検証とバージョン整合を含む
go run . llm | head
```

## 参照

- 標準: <https://github.com/ideamans/go-llm-cli-kit/blob/main/LLM.md>
- 生成物と原本の対応: `.claude/rules/ai-artifacts-policy.md`
- 再生成: `/regen-ai`
