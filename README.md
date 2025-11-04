# MoneyForward Kessai CLI Tool (mfk)

マネーフォワードケッサイAPIを操作するコマンドラインツール

## 概要

`mfk`は、マネーフォワードケッサイのAPIを効率的に操作するためのCLIツールです。
現在は請求書のダウンロード機能を提供していますが、今後さらに機能を拡張予定です。

## 主な機能

### 現在利用可能なコマンド

#### `download-invoice` - 請求書ダウンロード

- 指定期間の請求書を検索
- 請求書PDFの署名付きURL取得と自動ダウンロード
- 並列ダウンロード（同時実行数の指定可能）
- 詳細な進捗ログ出力
- ホームディレクトリのチルダ（~）解決
- 出力ディレクトリの自動作成

### 将来追加予定の機能

- 顧客管理コマンド
- 取引管理コマンド
- レポート生成コマンド
- その他、マネーフォワードケッサイAPIの各種機能

## インストール

### 方法1: makeを使用したインストール（推奨）

```bash
# ビルドとインストール（/usr/local/binにインストール）
make install

# アンインストール
make uninstall

# ビルドのみ（インストールしない）
make build

# ビルド成果物の削除
make clean

# 使用可能なコマンド一覧
make help
```

**注意:** `make install`は`sudo`を使用して`/usr/local/bin`にインストールします。パスワードの入力を求められます。

### 方法2: 手動ビルド

```bash
go build -o mfk ./cmd/main.go

# 手動でインストールする場合
sudo cp mfk /usr/local/bin/
sudo chmod 755 /usr/local/bin/mfk
```

### インストール後の確認

```bash
# どこからでもmfkコマンドが使えるようになります
which mfk
# /usr/local/bin/mfk

mfk --help
```

## 設定

### APIキーの設定

APIキーは以下の3つの方法で指定できます（優先順位の高い順）：

#### 1. コマンドラインオプション（最優先）

```bash
./mfk --api-key YOUR_API_KEY download-invoice 2024-09-01 2024-09-30 ./invoices
# または短縮形
./mfk -k YOUR_API_KEY download-invoice 2024-09-01 2024-09-30 ./invoices
```

#### 2. 環境変数

```bash
export MFK_API_KEY=your_api_key_here
./mfk download-invoice 2024-09-01 2024-09-30 ./invoices
```

#### 3. .envファイル

プロジェクトルートに `.env` ファイルを作成：

```bash
MFK_API_KEY=your_api_key_here
```

**APIキーの取得方法:**
- サンドボックス: https://sandbox-s.mfk.jp/developers/apikey
- 本番: https://s.mfk.jp/developers/apikey

**注意**:
- デフォルトでは本番環境 (`https://api.mfkessai.co.jp`) に接続します
- APIキーが指定されていない場合はエラーになります

## 使い方

### グローバルオプション

すべてのコマンドで使用可能なオプション：

- `-k, --api-key <key>` - MF Kessai APIキー（環境変数より優先）
- `-h, --help` - ヘルプを表示

### コマンド一覧

#### `download-invoice` - 請求書ダウンロード

指定期間の請求書PDFをダウンロードします。

**書式:**
```bash
./mfk [グローバルオプション] download-invoice [オプション] <開始日> <終了日> <出力ディレクトリ>
```

**オプション:**
- `-c, --concurrency <num>` - 並列ダウンロード数（デフォルト: 5）
- `-h, --help` - このコマンドのヘルプを表示

**引数:**
- `開始日` - 請求日の開始日（YYYY-MM-DD形式）
- `終了日` - 請求日の終了日（YYYY-MM-DD形式）
- `出力ディレクトリ` - PDFの保存先ディレクトリ

### 使用例

#### 例1: 基本的な使い方（.envファイルからAPIキーを読み込み）

```bash
./mfk download-invoice 2024-09-01 2024-09-30 ~/Downloads/mfk
```

#### 例2: APIキーをコマンドラインで指定

```bash
./mfk -k YOUR_API_KEY download-invoice 2024-09-01 2024-09-30 ./invoices
```

#### 例3: 並列数を指定してダウンロード

```bash
./mfk -c 10 download-invoice 2024-09-01 2024-09-30 ./invoices
```

#### 例4: APIキーと並列数の両方を指定

```bash
./mfk -k YOUR_API_KEY download-invoice -c 3 2024-09-01 2024-09-30 ./invoices
```

#### 例5: 環境変数でAPIキーを設定

```bash
export MFK_API_KEY=your_api_key_here
./mfk download-invoice 2025-01-01 2025-12-31 ~/Documents/invoices-2025
```

#### 例6: ヘルプの表示

```bash
# 全体のヘルプ
./mfk --help

# download-invoiceコマンドのヘルプ
./mfk download-invoice --help
```

### 日付の形式

日付は `YYYY-MM-DD` 形式で指定してください。

- 正しい例: `2024-09-01`
- 間違い: `2024/09/01`, `9-1-2024`

### 出力ファイル名

ダウンロードされたPDFファイルは、請求書IDをファイル名として保存されます：

```
{請求書ID}.pdf
```

例: `BILL-XXXXX.pdf`

## ログ出力

ツールは以下の情報をログ出力します：

1. **開始時**
   - 日付範囲
   - 出力ディレクトリ
   - 並列数

2. **検索結果**
   - 見つかった請求書の件数

3. **ダウンロード進捗**
   - 各請求書の処理状況（進捗/全体）
   - 請求書ID、発行日、金額
   - ダウンロード成功/失敗

4. **完了時のサマリー**
   - 全体の請求書数
   - 成功数
   - 失敗数
   - エラー詳細（エラーがある場合）

### ログ出力例

```
2025/11/04 10:56:02 Starting invoice download
2025/11/04 10:56:02 Date range: 2024-09-01 to 2024-09-30
2025/11/04 10:56:02 Output directory: ./tmp
2025/11/04 10:56:02 Concurrency: 5
2025/11/04 10:56:02 Fetching billings from 2024-09-01 to 2024-09-30...
2025/11/04 10:56:03 Found 15 billing(s)
2025/11/04 10:56:03 [1/15] Processing billing ID: BILL-00001 (Issue Date: 2024-09-01, Amount: 100000)
2025/11/04 10:56:04 [1/15] SUCCESS: Saved ./tmp/BILL-00001.pdf (52341 bytes)
2025/11/04 10:56:04 [2/15] Processing billing ID: BILL-00002 (Issue Date: 2024-09-05, Amount: 250000)
...
2025/11/04 10:56:15
=== Download Summary ===
2025/11/04 10:56:15 Total billings: 15
2025/11/04 10:56:15 Successful downloads: 15
2025/11/04 10:56:15 Failed downloads: 0
2025/11/04 10:56:15 Download completed successfully
```

## エラーハンドリング

ツールは以下のエラーに対応しています：

- **APIキーが見つからない**: `.env` ファイルに `MFK_API_KEY` が設定されているか確認してください
- **API認証エラー**: APIキーが正しいか確認してください
- **ネットワークエラー**: インターネット接続を確認してください
- **ディレクトリ作成エラー**: 出力先ディレクトリのパーミッションを確認してください
- **ファイル書き込みエラー**: ディスク容量とパーミッションを確認してください

エラーが発生した場合、ツールは詳細なエラーメッセージを表示し、エラーの原因を特定できるようにします。

## プロジェクト構造

```
.
├── cmd/
│   └── main.go                 # エントリポイント
├── internal/
│   ├── cmd/
│   │   ├── root.go            # Cobraルートコマンド（グローバルフラグ定義）
│   │   └── download_invoice.go # download-invoiceサブコマンド実装
│   ├── mfkessai/
│   │   └── client.go          # MF Kessai APIクライアント
│   └── downloader/
│       └── downloader.go      # 並列ダウンロード処理
├── tmp/                        # テスト用出力ディレクトリ
├── .env                        # 環境変数（APIキー）※オプション
├── .gitignore                  # Git除外設定
├── Makefile                    # ビルド・インストール用Makefile
├── go.mod                      # Go依存関係管理
├── go.sum                      # 依存関係チェックサム
├── INVOICE.md                  # MF Kessai 請求APIドキュメント
└── README.md                   # このファイル
```

### 拡張性について

このツールは拡張可能な設計になっています：

**新しいサブコマンドの追加方法:**
1. `internal/cmd/`に新しいコマンドファイルを作成（例: `list_customers.go`）
2. `root.go`の`init()`関数でコマンドを登録
3. グローバルフラグ（`-k/--api-key`等）は自動的に全サブコマンドで利用可能

**例: 顧客一覧コマンドの追加**
```go
// internal/cmd/list_customers.go
var listCustomersCmd = &cobra.Command{
    Use:   "list-customers",
    Short: "List all customers",
    RunE:  runListCustomers,
}

func init() {
    rootCmd.AddCommand(listCustomersCmd)  // root.goで登録
}
```

## API仕様

このツールは以下のMF Kessai API v2エンドポイントを使用します：

1. `GET /v2/billings` - 請求書一覧取得
2. `POST /v2/billings/{billing_id}/download_signed_url` - 署名付きダウンロードURL取得

詳細は [INVOICE.md](./INVOICE.md) を参照してください。

## 依存関係

- [cobra](https://github.com/spf13/cobra) - CLIフレームワーク
- [godotenv](https://github.com/joho/godotenv) - .envファイル読み込み

## ライセンス

MIT

## 開発

### Makefileコマンド

```bash
# ビルド
make build

# テスト実行
make test

# 依存関係の整理
make tidy

# クリーンアップ
make clean

# ヘルプ表示
make help
```

### 手動でのテスト実行

```bash
# ビルドしたバイナリでテスト
make build
./mfk download-invoice 2023-06-28 2023-06-28 ./tmp

# 並列数を変えてテスト
./mfk -c 3 download-invoice 2023-06-01 2023-08-31 ./tmp

# APIキーを指定してテスト
./mfk -k YOUR_API_KEY download-invoice 2023-06-28 2023-06-28 ./tmp
```

### 開発ワークフロー

```bash
# 1. コードを編集
# 2. ビルド
make build

# 3. テスト実行
./mfk download-invoice 2023-06-28 2023-06-28 ./tmp

# 4. 問題なければインストール
make install

# 5. システム全体でテスト
mfk download-invoice 2023-06-28 2023-06-28 ~/Downloads/test
```

## 技術詳細

### ダウンロードエンドポイントについて

請求書のダウンロードには、以下のエンドポイントを使用します：

```
POST /v2/billings/{billing_id}/issues/i/{invoice_id}/download_signed_url
```

**重要なポイント:**
- `i`は請求書（invoice）を示すリテラルパスセグメントです
- 口座振替通知書の場合は`a`を使用します
- `invoice_id`は`GET /v2/billings`のレスポンスの`invoice_ids`配列から取得できます

**レスポンス例:**
```json
{
  "items": [
    {
      "signed_url": "https://download.mfk.jp/api/billings/...",
      "expired_at": "2023-03-08T10:36:43+09:00",
      "type": "pdf"
    }
  ],
  "object": "list"
}
```

署名付きURLには有効期限があるため、取得後速やかにダウンロードを実行します。

## トラブルシューティング

### "No billings found" と表示される

指定した期間に請求書が存在しない可能性があります。日付範囲を変更して再試行してください。

### "API error (status 401)" と表示される

APIキーが正しくない、または期限切れの可能性があります。`.env` ファイルの `MFK_API_KEY` を確認してください。

### ダウンロードが遅い

並列数を増やすことで高速化できます。ただし、API側のレート制限に注意してください：

```bash
./mfk -c 10 download-invoice 2024-09-01 2024-09-30 ./invoices
```

### "API rate limit exceeded" と表示される

並列数が多すぎます。`-c`オプションで並列数を減らしてください：

```bash
./mfk -c 1 download-invoice 2024-09-01 2024-09-30 ./invoices
```
