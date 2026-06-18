package cmd

// llmGuide is a comprehensive, self-contained guide printed by `mfk --llm`.
// It is intended for LLM agents driving this CLI, so it favors explicit,
// copy-pasteable commands and clearly documents non-obvious behavior such as
// the split between legacy and qualified (invoice-system) billings.
const llmGuide = `# mfk CLI — LLM エージェント向けガイド

mfk はマネーフォワードケッサイ API v2 を操作する CLI です。
すべて ` + "`mfk <リソース> <メソッド> [フラグ]`" + ` の形式で実行します。
デフォルト出力は JSON なので、` + "`jq`" + ` と組み合わせて処理できます。

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

## 請求書の取得（重要: インボイス制度で取得先が2つに分かれています）

マネーフォワードケッサイでは、インボイス制度（適格請求書等保存方式、2023年10月開始）を境に
請求オブジェクトの取得エンドポイントが2系統に分かれています。**取得したい期間で使い分けてください。**

| コマンド | エンドポイント | 内容 | 対象期間 |
|----------|----------------|------|----------|
| mfk billings qualified | GET /billings/qualified | インボイス制度（適格請求書等保存方式）対応の請求 | 2023年10月以降（現行・メイン） |
| mfk billings list | GET /billings | 区分記載請求書等保存方式の請求（旧方式） | おおよそ2023年10月より前 |

**原則: 現行の請求書は ` + "`mfk billings qualified`" + ` を使ってください。**
` + "`mfk billings list`" + ` は旧方式（区分記載請求書）の請求しか返さないため、インボイス制度開始以降の
データは含まれません。「2023年までのデータしか取れない」場合は qualified を使っていないのが原因です。

両者はクエリパラメータ・レスポンス構造（Billing オブジェクトの一覧）は共通です。
全件取得する場合は --page-all を付けます。

` + "```bash" + `
# 現行（インボイス制度対応）の請求書を全件取得
mfk billings qualified --page-all

# 顧客で絞り込み
mfk billings qualified --customer-id cust_xxx --page-all

# 旧方式（〜2023年）の請求書にアクセスしたいときだけ list を使う
mfk billings list --page-all

# 請求の詳細取得（ID は qualified/list どちらの一覧から得た id でも可）
mfk billings get <billing_id>
` + "```" + `

### 請求関連のその他メソッド

` + "```bash" + `
mfk billings get <billing_id>                 # 請求の詳細取得
mfk billings reissue <billing_id> --json ...   # 請求書の再発行
mfk billings upload-signed-url --json ...       # アップロード用署名URL取得
mfk billings download-signed-url --json ...     # ダウンロード用署名URL取得
` + "```" + `

## リソースとメソッド一覧

### customers（取引先）
` + "```bash" + `
mfk customers list [--number <S>] [--name <S>] [--has-alert]
mfk customers get <customer_id>
mfk customers create [--json <JSON>] [--name ...] [--number ...] [--dest-name ...] ...
mfk customers update <customer_id> [--json <JSON>] [--name ...]
mfk customers bank-transfer <customer_id> --json <JSON>
mfk customers account-transfer <customer_id> --json <JSON>
` + "```" + `

### destinations（請求先）
` + "```bash" + `
mfk destinations list [--customer-id <ID>]
mfk destinations get <destination_id>
mfk destinations create [--json <JSON>] [--customer-id ...] [--name ...] ...
mfk destinations update <destination_id> [--json <JSON>] [--name ...] ...
` + "```" + `

### transactions（取引）
` + "```bash" + `
mfk transactions list [--customer-id <ID>] [--destination-id <ID>] [--status <S>]
mfk transactions get <id>
mfk transactions create [--json <JSON>] [--destination-id ...] [--amount ...] [--date ...] ...
mfk transactions delete <id>
` + "```" + `

### billings（請求） … 上記「請求書の取得」を参照
### customer-examinations（与信枠審査）
` + "```bash" + `
mfk customer-examinations list [--customer-id <ID>] [--status <S>]
mfk customer-examinations get <id>
mfk customer-examinations create [--json <JSON>] [--customer-id ...] [--amount ...] ...
` + "```" + `

### credit-facilities（与信枠）
` + "```bash" + `
mfk credit-facilities list [--customer-id <ID>] [--customer-examination-id <ID>]
mfk credit-facilities get <id>
` + "```" + `

### payouts / payout-transactions / payout-refunds（振込・振込明細・返金）
` + "```bash" + `
mfk payouts list
mfk payouts get <id>
mfk payout-transactions list [--payout-id <ID>]
mfk payout-transactions get <id>
mfk payout-refunds list
` + "```" + `

### customer-name-updates（取引先名変更）
` + "```bash" + `
mfk customer-name-updates list
mfk customer-name-updates get <id>
mfk customer-name-updates create [--json <JSON>] [--customer-id ...] [--name ...]
mfk customer-name-updates delete <id>
` + "```" + `

### authorizations（オーソリ）
` + "```bash" + `
mfk authorizations list
mfk authorizations get <id>
mfk authorizations create [--json <JSON>] [--customer-id ...] [--amount ...] [--destination-id ...]
mfk authorizations delete <id>
` + "```" + `

## 出力の扱い

- デフォルトは JSON。一覧は ` + "`{\"object\":\"list\",\"items\":[...],\"pagination\":{...}}`" + ` 形式。
- --page-all を付けると全ページの items を結合した ` + "`{\"items\":[...]}`" + ` を返します。
- jq での絞り込み例: ` + "`mfk billings qualified --page-all | jq '.items[] | select(.status==\"invoice_issued\")'`" + `

## よくあるタスクの手順

- 「今年の請求書を全部取得」→ ` + "`mfk billings qualified --page-all`" + `（list ではなく qualified）
- 「特定顧客の現行請求」→ ` + "`mfk billings qualified --customer-id <ID> --page-all`" + `
- 「2023年より前の古い請求」→ ` + "`mfk billings list --page-all`" + `
`
