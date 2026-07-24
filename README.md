# mfk - マネーフォワードケッサイ CLI

マネーフォワードケッサイ API v2 をコマンドラインから操作するための CLI ツールです。

## インストール

### Homebrew (macOS / Linux)

```bash
brew install ideamans/tap/mfk
```

### Go

```bash
go install github.com/ideamans/mfk-cli@latest
```

### GitHub Releases

[Releases](https://github.com/ideamans/mfk-cli/releases) からバイナリをダウンロードしてください。

## セットアップ

環境変数に API キーを設定してください。

```bash
export MFK_API_KEY="your-api-key"
```

サンドボックス環境を使用する場合:

```bash
export MFK_ENV=sandbox
export MFK_SANDBOX_API_KEY="your-sandbox-api-key"
```

## 使い方

```
mfk <リソース> <メソッド> [フラグ]
```

### グローバルフラグ

| フラグ | 説明 |
|--------|------|
| `--format <json\|table\|csv>` | 出力形式（デフォルト: json） |
| `--sandbox` | サンドボックス環境を使用 |
| `--page-all` | 全ページを自動取得 |
| `--limit <N>` | 取得件数 |
| `--dry-run` | リクエスト内容を表示するのみ |
| `--json <JSON>` | リクエストボディを JSON で指定 |
| `--llm` | 非推奨。`mfk llm` を使ってください（互換のため動作します） |

### リソース一覧

| リソース | 説明 |
|----------|------|
| `customers` | 取引先 |
| `destinations` | 請求先 |
| `transactions` | 取引 |
| `billings` | 請求（`qualified` = インボイス制度対応・現行 / `list` = 旧請求） |
| `issues` | 発行済み請求情報 |
| `customer-examinations` | 与信枠審査 |
| `credit-facilities` | 与信枠 |
| `payouts` | 振込 |
| `payout-transactions` | 振込明細 |
| `payout-refunds` | 返金 |
| `customer-name-updates` | 取引先名変更 |
| `authorizations` | オーソリ |

### 使用例

```bash
# 取引先一覧を取得
mfk customers list

# テーブル形式で表示
mfk customers list --format table

# 取引先名でフィルタ
mfk customers list --name "テスト株式会社"

# CLI フラグで取引先を作成
mfk customers create --name "テスト株式会社" --number "CUST-001" \
  --dest-name "経理部" --dest-name-kana "ケイリブ" \
  --dest-email "keiri@example.com" --dest-tel "03-1234-5678" \
  --dest-zip-code "100-0001" --dest-address1 "東京都千代田区1-1"

# JSON で取引先を作成
mfk customers create --json '{"name":"テスト株式会社","number":"CUST-001","destination":{...}}'

# JSON をベースに CLI フラグで一部上書き
mfk customers create --json "$(cat template.json)" --name "別の会社名"

# 取引先名変更
mfk customer-name-updates create --customer-id "cust_abc" --name "新しい名前"

# オーソリ作成
mfk authorizations create --customer-id "cust_abc" --amount 500000 --destination-id "dest_abc"

# 全ページを自動取得して jq でフィルタ
mfk transactions list --page-all | jq '.items[] | select(.status == "passed")'

# dry-run でリクエスト内容を確認
mfk customers create --dry-run --name "テスト" --number "TEST-001"

# サンドボックス環境を使用
mfk customers list --sandbox
```

### 請求書の取得（インボイス制度対応）

マネーフォワードケッサイでは、**インボイス制度（適格請求書等保存方式、2023年10月開始）を境に請求の取得エンドポイントが2系統に分かれています。** 取得したい期間に応じて使い分けてください。

| コマンド | API | 内容 | 対象期間 |
|----------|-----|------|----------|
| `mfk billings qualified` | `GET /billings/qualified` | インボイス制度（適格請求書等保存方式）対応の請求 | **2023年10月以降（現行・メイン）** |
| `mfk billings list` | `GET /billings` | 区分記載請求書等保存方式の請求（旧方式） | おおよそ2023年10月より前 |

**現行の請求書は `mfk billings qualified` を使ってください。** `mfk billings list` は旧方式（区分記載請求書）の請求しか返さないため、インボイス制度開始以降のデータは含まれません（「2023年までのデータしか取得できない」場合は qualified を使っていないのが原因です）。

全件取得するには `--page-all` を付けます（付けないと最新20件のみ）。

```bash
# 現行（インボイス制度対応）の請求書を全件取得
mfk billings qualified --page-all

# 特定顧客の現行請求
mfk billings qualified --customer-id cust_xxx --page-all

# 旧方式（〜2023年）の請求書にアクセスしたいときだけ list を使う
mfk billings list --page-all
```

### AIエージェントから使う

`mfk llm` で、全コマンドと請求書（qualified / list）の使い分けを含む詳細リファレンスを Markdown で出力します。コマンドカタログは cobra の定義から生成されるため実装と乖離せず、バイナリに埋め込まれているのでオフラインでも動作します。

```bash
mfk llm                  # Markdown
mfk llm --format json    # 章ごとのJSON配列
mfk --llm                # 非推奨エイリアス。従来どおりどの位置でも動作
```

Claude Code ではプラグインを導入すると `/mfk-usage` と `/mfk-install` が使えます。

```
/plugin marketplace add ideamans/claude-public-plugins
/plugin install mfk-cli@ideamans-plugins
```

同じスキルは Copilot や Cursor など Agent Skills 対応ホストでも利用できます。

```bash
gh skill install ideamans/mfk-cli/plugins/mfk-cli/skills/mfk-usage --agent copilot
```

### CLI フラグと `--json` の併用

- `--json` でリクエストボディ全体を指定できます
- よく使うフィールドは個別の CLI フラグでも指定可能です
- 両方を指定した場合、**CLI フラグが優先**されます（`--json` をベースにマージ・上書き）
- 配列やネストの深い構造は `--json` で指定してください

## 環境変数

| 変数 | 必須 | 説明 |
|------|------|------|
| `MFK_API_KEY` | 必須 | API キー |
| `MFK_SANDBOX_API_KEY` | 任意 | サンドボックス用 API キー（sandbox 時に優先使用） |
| `MFK_ENV` | 任意 | `production`（デフォルト）または `sandbox` |

### 環境の決定順序

1. `--sandbox` フラグ
2. `MFK_ENV` 環境変数
3. デフォルト: `production`

## 開発

```bash
# ビルド
make build

# 開発用（ビルド + 実行）
make run ARGS="customers list --dry-run"

# テスト
make test

# リント
make lint

# クリーンアップ
make clean
```

## ライセンス

MIT
