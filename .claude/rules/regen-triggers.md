---
paths:
  - "cmd/*.go"
  - "internal/llmdocs/0*.md"
  - "internal/llmdocs/1*.md"
---

# 埋め込みリファレンスの原本に触れました

コマンド・フラグ・ヘルプ文字列を変えたなら、終わる前に `/regen-ai` を実行して
`internal/llmdocs/90-commands.md` を一致させること。CI が再生成して差分が出たら
落ちます。

請求書の取得系統（`billings qualified` / `billings list`）やページネーションの
既定値に触れたなら、`10-billings.md` と `00-guide.md` を手で直すこと。
ここがズレると、エージェントは**エラーにならずに不完全なデータで集計します**。

`90-commands.md` は直接編集しないこと。
