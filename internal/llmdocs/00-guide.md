# mfk CLI — AIエージェント向けリファレンス

mfk はマネーフォワードケッサイ API v2 を操作する CLI です。
すべて `mfk <リソース> <メソッド> [フラグ]` の形式で実行します。
デフォルト出力は JSON なので、`jq` と組み合わせて処理できます。

このリファレンスはバイナリに埋め込まれています。`mfk llm` は常に実行中の
バージョンそのものを説明します。

## 認証・環境

| 環境変数 | 必須 | 説明 |
|----------|------|------|
| MFK_API_KEY | 必須 | 本番 API キー（リクエストヘッダー apikey に設定） |
| MFK_SANDBOX_API_KEY | 任意 | サンドボックス用 API キー（sandbox 時に優先使用） |
| MFK_ENV | 任意 | production（デフォルト）または sandbox |

環境の決定順: (1) --sandbox フラグ → (2) MFK_ENV → (3) production。
ベースURL: production=https://api.mfkessai.co.jp/v2, sandbox=https://sandbox-api.mfkessai.co.jp/v2

## グローバルフラグ

| フラグ | 説明 |
|--------|------|
| --format <json\|table\|csv> | 出力形式（デフォルト json） |
| --sandbox | サンドボックス環境を使用 |
| --page-all | 全ページを自動取得（ページネーションを内部で辿る） |
| --limit <N> | 1ページの取得件数（1〜200、デフォルト 20） |
| --dry-run | リクエスト内容を表示するのみ（実行しない） |
| --json <JSON> | リクエストボディを JSON で指定（CLI フラグと併用時はフラグが優先） |
| --llm | このガイドを表示して終了 |

ページネーション: 一覧は「追加日時の降順」で返り、デフォルト 20 件。
全件取得するには --page-all を付けます（内部で pagination.has_next / pagination.end を辿って全ページを結合）。
件数が多い場合は --page-all を必ず付けてください。付けないと最新 20 件しか取得できません。

## 出力の扱い

- デフォルトは JSON。一覧は `{"object":"list","items":[...],"pagination":{...}}` 形式。
- --page-all を付けると全ページの items を結合した `{"items":[...]}` を返します。
- jq での絞り込み例: `mfk billings qualified --page-all | jq '.items[] | select(.status=="invoice_issued")'`

## よくあるタスクの手順

- 「今年の請求書を全部取得」→ `mfk billings qualified --page-all`（list ではなく qualified）
- 「特定顧客の現行請求」→ `mfk billings qualified --customer-id <ID> --page-all`
- 「2023年より前の古い請求」→ `mfk billings list --page-all`

## 課金・与信に関わる操作について

このCLIは実際の請求・与信・支払データを操作します。作成・更新・取り消し系の
メソッドは取引先に影響するため、実行前に `--dry-run` で内容を確認し、
何をどの顧客に対して行うのかをユーザーに提示して同意を得てください。
