# mfk CLI - マネーフォワードケッサイ CLI ツール設計計画

## 概要

マネーフォワードケッサイ API v2 を操作する Go 製 CLI ツール。
gws (Google Workspace CLI) のコマンド体系を参考に、`mfk <リソース> <メソッド> [フラグ]` 形式で操作する。

## コマンド体系

```
mfk <resource> <method> [flags]
```

### グローバルフラグ

| フラグ | 説明 |
|--------|------|
| `--format <json\|table\|csv>` | 出力形式（デフォルト: json） |
| `--sandbox` | サンドボックス環境を使用（`MFK_ENV=sandbox` と同等、フラグが優先） |
| `--page-all` | 全ページを自動取得 |
| `--limit <N>` | 取得件数 |
| `--dry-run` | リクエスト内容を表示するのみ（実行しない） |
| `-h, --help` | ヘルプ表示 |

### CLI フラグと `--json` の併用ルール

よく使うフィールドは個別の CLI フラグとしても指定可能。`--json` と CLI フラグを同時に指定した場合、**CLI フラグが優先**（`--json` のオブジェクトをベースに CLI フラグの値でマージ・上書き）。

```go
// 実装イメージ: --json をベースに CLI フラグで上書き
body := parseJSON(jsonFlag)        // --json があればパース（なければ空 map）
if nameFlag != "" {                // CLI フラグがあれば上書き
    body["name"] = nameFlag
}
```

### リソースとメソッド一覧

#### customers

```bash
# 一覧取得（フィルタ用フラグ）
mfk customers list [--number <S>] [--name <S>] [--has-alert]

# 詳細取得
mfk customers get <customer_id>

# 作成（--json でも CLI フラグでも可、併用も可）
mfk customers create [--json <JSON>] \
  [--name <名前>] [--number <番号>] \
  [--dest-name <請求先名>] [--dest-name-kana <カナ>] \
  [--dest-email <メール>] [--dest-tel <電話>] \
  [--dest-zip-code <郵便番号>] [--dest-address1 <住所1>] [--dest-address2 <住所2>] \
  [--dest-department <部署>] [--dest-title <肩書き>]

# 更新
mfk customers update <customer_id> [--json <JSON>] \
  [--name <名前>]

# 銀行振込口座の割り当て
mfk customers bank-transfer <customer_id> --json <JSON>

# 口座振替申込
mfk customers account-transfer <customer_id> --json <JSON>
```

#### destinations

```bash
# 一覧取得
mfk destinations list [--customer-id <ID>]

# 詳細取得
mfk destinations get <destination_id>

# 作成
mfk destinations create [--json <JSON>] \
  [--customer-id <ID>] \
  [--name <名前>] [--name-kana <カナ>] \
  [--email <メール>] [--tel <電話>] \
  [--zip-code <郵便番号>] [--address1 <住所1>] [--address2 <住所2>] \
  [--department <部署>] [--title <肩書き>]

# 更新
mfk destinations update <destination_id> [--json <JSON>] \
  [--name <名前>] [--name-kana <カナ>] \
  [--email <メール>] [--tel <電話>] \
  [--zip-code <郵便番号>] [--address1 <住所1>] [--address2 <住所2>] \
  [--department <部署>] [--title <肩書き>]
```

#### transactions

```bash
# 一覧取得
mfk transactions list [--customer-id <ID>] [--destination-id <ID>] [--status <S>]

# 詳細取得
mfk transactions get <id>

# 作成（amount, date 等の基本項目は CLI フラグ、明細は --json で）
mfk transactions create [--json <JSON>] \
  [--destination-id <ID>] [--amount <金額>] [--number <番号>] \
  [--date <YYYY-MM-DD>] [--due-date <YYYY-MM-DD>] [--issue-date <YYYY-MM-DD>] \
  [--delivery-method <posting|email>]

# 削除
mfk transactions delete <id>
```

#### billings

```bash
mfk billings list [--customer-id <ID>] [--destination-id <ID>] [--status <S>]
mfk billings qualified [--customer-id <ID>] [--destination-id <ID>]
mfk billings get <id>
mfk billings reissue <id> --json <JSON>

# アップロード用署名URL取得
mfk billings upload-signed-url --json <JSON>

# ダウンロード用署名URL取得
mfk billings download-signed-url --json <JSON>
```

#### issues

```bash
# 発行済み請求情報一覧
mfk issues list
```

#### customer-examinations

```bash
mfk customer-examinations list [--customer-id <ID>] [--status <S>]
mfk customer-examinations get <id>

# 作成
mfk customer-examinations create [--json <JSON>] \
  [--customer-id <ID>] [--amount <金額>] [--end-date <YYYY-MM-DD>] \
  [--address1 <住所>] [--email <メール>] [--tel <電話>] \
  [--zip-code <郵便番号>] [--business-type <業種>] \
  [--representative-name <代表者名>]
```

#### credit-facilities

```bash
mfk credit-facilities list [--customer-id <ID>] [--customer-examination-id <ID>]
mfk credit-facilities get <id>
```

#### payouts / payout-transactions / payout-refunds

```bash
mfk payouts list
mfk payouts get <id>

mfk payout-transactions list [--payout-id <ID>]
mfk payout-transactions get <id>

mfk payout-refunds list
```

#### customer-name-updates

```bash
mfk customer-name-updates list
mfk customer-name-updates get <id>

# 作成
mfk customer-name-updates create [--json <JSON>] \
  [--customer-id <ID>] [--name <新しい名前>]

mfk customer-name-updates delete <id>
```

#### authorizations

```bash
mfk authorizations list
mfk authorizations get <id>

# 作成
mfk authorizations create [--json <JSON>] \
  [--customer-id <ID>] [--amount <金額>] [--destination-id <ID>]

mfk authorizations delete <id>
```

### CLI フラグ対応方針

| 方針 | 説明 |
|------|------|
| **フラット項目を優先** | `name`, `amount`, `date` 等の単純な文字列・数値はフラグ化する |
| **ネスト項目は prefix** | `destination.name` → `--dest-name`、`destination.email` → `--dest-email` のように prefix を付ける |
| **配列・複雑構造は `--json` のみ** | `transaction_details`, `amounts_per_tax_rate_type`, `cc_emails` 等は `--json` で指定 |
| **list のフィルタは全てフラグ化** | API のクエリパラメータはすべて CLI フラグとして提供する |
| **`--json` は常に使える** | どのコマンドでも `--json` でフルボディを渡せる。フラグとの併用時はフラグが優先 |

## 環境変数

| 環境変数 | 必須 | 説明 |
|----------|------|------|
| `MFK_API_KEY` | 必須 | API キー。リクエストヘッダー `apikey` に設定。未設定時はエラー終了 |
| `MFK_SANDBOX_API_KEY` | 任意 | サンドボックス用 API キー。sandbox 環境時にこちらが設定されていれば優先使用 |
| `MFK_ENV` | 任意 | `production`（デフォルト）または `sandbox`。API のベースURLを切り替える |

## 環境の決定と API キーの解決

### 環境の決定（優先順位）

1. `--sandbox` フラグ → sandbox
2. `MFK_ENV` 環境変数 → `sandbox` or `production`
3. デフォルト → production

`MFK_ENV` に `production` / `sandbox` 以外の値が設定された場合はエラー終了。

### API キーの解決

| 環境 | 使用する API キー |
|------|-------------------|
| sandbox | `MFK_SANDBOX_API_KEY`（設定されていれば優先）→ フォールバック: `MFK_API_KEY` |
| production | `MFK_API_KEY` |

いずれも未設定の場合はエラー終了。

## API エンドポイント

| 環境 | ベースURL |
|------|-----------|
| production（デフォルト） | `https://api.mfkessai.co.jp/v2` |
| sandbox | `https://sandbox-api.mfkessai.co.jp/v2` |

## 技術スタック

| 項目 | 選定 |
|------|------|
| 言語 | Go 1.22+ |
| CLI フレームワーク | [cobra](https://github.com/spf13/cobra) |
| HTTP クライアント | 標準 `net/http` |
| JSON 処理 | 標準 `encoding/json` |
| テーブル出力 | [tablewriter](https://github.com/olekukonez/tablewriter) |

## プロジェクト構造

```
mfk-cli/
├── main.go                  # エントリポイント
├── go.mod
├── go.sum
├── plan.md
├── cmd/
│   ├── root.go              # ルートコマンド（グローバルフラグ）
│   ├── customers.go          # customers サブコマンド
│   ├── destinations.go       # destinations サブコマンド
│   ├── customer_examinations.go
│   ├── credit_facilities.go
│   ├── transactions.go
│   ├── billings.go
│   ├── issues.go
│   ├── payouts.go
│   ├── payout_transactions.go
│   ├── payout_refunds.go
│   ├── customer_name_updates.go
│   └── authorizations.go
├── pkg/
│   ├── api/
│   │   ├── client.go         # HTTP クライアント（認証・ベースURL管理）
│   │   └── client_test.go
│   └── output/
│       ├── formatter.go      # json / table / csv 出力
│       └── formatter_test.go
```

## 実装フェーズ

### Phase 1: 基盤（MVP）

1. **プロジェクト初期化** - `go mod init`, cobra セットアップ
2. **API クライアント** - `pkg/api/client.go`
   - `MFK_API_KEY` / `MFK_SANDBOX_API_KEY` / `MFK_ENV` 読み取り
   - GET / POST / PATCH / DELETE メソッド
   - `--sandbox` フラグ / `MFK_ENV` による環境切り替え
   - sandbox 時は `MFK_SANDBOX_API_KEY` を優先使用
   - `--dry-run` 対応
   - ページネーション（`--page-all`）
3. **出力フォーマッター** - `pkg/output/formatter.go`
   - JSON（デフォルト、`jq` 互換のインデント付き）
   - テーブル形式
   - CSV 形式
4. **ルートコマンド** - `cmd/root.go`

### Phase 2: 主要リソース

5. **customers** コマンド - CRUD + bank-transfer + account-transfer
6. **destinations** コマンド - CRUD
7. **transactions** コマンド - CRUD
8. **billings** コマンド - list / qualified / get / reissue / upload-signed-url / download-signed-url
9. **issues** コマンド - list

### Phase 3: 残りのリソース

10. **customer-examinations** コマンド
11. **credit-facilities** コマンド
12. **payouts** / **payout-transactions** / **payout-refunds** コマンド
13. **customer-name-updates** コマンド
14. **authorizations** コマンド

### Phase 4: 品質向上

14. エラーハンドリング改善（API エラーレスポンスの整形表示）
15. テスト追加
16. README 作成

## 使用例

```bash
# 取引先一覧を取得
mfk customers list

# 取引先一覧をテーブル形式で取得
mfk customers list --format table

# 取引先名でフィルタ
mfk customers list --name "テスト株式会社"

# CLI フラグだけで取引先を作成（シンプルなケース）
mfk customers create --name "テスト株式会社" --number "CUST-001" \
  --dest-name "経理部" --dest-name-kana "ケイリブ" \
  --dest-email "keiri@example.com" --dest-tel "03-1234-5678" \
  --dest-zip-code "100-0001" --dest-address1 "東京都千代田区1-1"

# --json でフルボディを渡す（複雑なケース）
mfk customers create --json '{"name":"テスト株式会社","number":"CUST-001","destination":{...}}'

# --json をベースに CLI フラグで一部上書き
mfk customers create --json "$(cat customer-template.json)" --name "別の会社名"

# 取引の作成（基本項目は CLI フラグ、明細は --json）
mfk transactions create \
  --destination-id "dest_abc123" --number "INV-2026-001" \
  --amount 110000 --date 2026-03-19 --due-date 2026-04-30 --issue-date 2026-03-20 \
  --delivery-method email \
  --json '{"transaction_details":[{"description":"サービス利用料","amount":110000,"quantity":1,"unit_price":100000,"tax_rate_type":"normal_10","tax_included_type":"excluded"}],"amounts_per_tax_rate_type":[{"tax_rate_type":"normal_10","amount":110000}]}'

# 取引を JSON ファイルから作成
mfk transactions create --json "$(cat transaction.json)"

# 与信枠審査を申請
mfk customer-examinations create \
  --customer-id "cust_abc123" --amount 1000000 --end-date 2026-12-31 \
  --address1 "東京都千代田区1-1" --email "info@example.com" --tel "03-1234-5678" \
  --zip-code "100-0001" --business-type "IT" --representative-name "山田太郎"

# 取引先名変更
mfk customer-name-updates create --customer-id "cust_abc123" --name "新しい会社名"

# オーソリ作成
mfk authorizations create --customer-id "cust_abc123" --amount 500000 --destination-id "dest_abc123"

# 全ページ自動取得して jq でフィルタ
mfk transactions list --page-all | jq '.items[] | select(.status == "passed")'

# サンドボックス環境で dry-run（--sandbox フラグ）
mfk customers create --sandbox --dry-run --name "テスト" --number "TEST-001" \
  --dest-name "テスト" --dest-name-kana "テスト" \
  --dest-email "test@example.com" --dest-tel "000-0000-0000" \
  --dest-zip-code "000-0000" --dest-address1 "テスト住所"

# または環境変数で切り替え（.envrc 等で設定する想定）
# export MFK_ENV=sandbox
# export MFK_SANDBOX_API_KEY=your-sandbox-api-key
```

## 設計方針

- **シンプルさ優先**: gws のように動的にスキーマを取得する仕組みは不要。APIは固定的なのでコマンドを静的に定義する
- **JSON ファースト**: デフォルト出力は JSON。パイプで `jq` と組み合わせやすくする
- **最小依存**: cobra + tablewriter のみ。HTTP は標準ライブラリ
- **`--json` + CLI フラグ併用**: `--json` でフルボディも渡せるが、よく使うフラット項目は個別フラグでも指定可能。併用時は CLI フラグが優先（`--json` をベースにマージ上書き）。配列やネスト構造は `--json` のみ
