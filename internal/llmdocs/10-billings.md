# 請求書の取得（重要: インボイス制度で取得先が2つに分かれています）

マネーフォワードケッサイでは、インボイス制度（適格請求書等保存方式、2023年10月開始）を境に
請求オブジェクトの取得エンドポイントが2系統に分かれています。**取得したい期間で使い分けてください。**

| コマンド | エンドポイント | 内容 | 対象期間 |
|----------|----------------|------|----------|
| mfk billings qualified | GET /billings/qualified | インボイス制度（適格請求書等保存方式）対応の請求 | 2023年10月以降（現行・メイン） |
| mfk billings list | GET /billings | 区分記載請求書等保存方式の請求（旧方式） | おおよそ2023年10月より前 |

**原則: 現行の請求書は `mfk billings qualified` を使ってください。**
`mfk billings list` は旧方式（区分記載請求書）の請求しか返さないため、インボイス制度開始以降の
データは含まれません。「2023年までのデータしか取れない」場合は qualified を使っていないのが原因です。

両者はクエリパラメータ・レスポンス構造（Billing オブジェクトの一覧）は共通です。
全件取得する場合は --page-all を付けます。

```bash
# 現行（インボイス制度対応）の請求書を全件取得
mfk billings qualified --page-all

# 顧客で絞り込み
mfk billings qualified --customer-id cust_xxx --page-all

# 旧方式（〜2023年）の請求書にアクセスしたいときだけ list を使う
mfk billings list --page-all

# 請求の詳細取得（ID は qualified/list どちらの一覧から得た id でも可）
mfk billings get <billing_id>
```

### 請求関連のその他メソッド

```bash
mfk billings get <billing_id>                 # 請求の詳細取得
mfk billings reissue <billing_id> --json ...   # 請求書の再発行
mfk billings upload-signed-url --json ...       # アップロード用署名URL取得
mfk billings download-signed-url --json ...     # ダウンロード用署名URL取得
```
