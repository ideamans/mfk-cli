# リソースとメソッド一覧

cobra のコマンド定義から `go generate ./...` で生成しています。
手編集しないこと。直すのはコマンド定義側です。

## Global flags

| flag | type | default | description |
| --- | --- | --- | --- |
| `--dry-run` | bool | `false` | Show request without executing |
| `--format` | string | `json` | Output format: json, table, csv |
| `--json` | string | — | Request body as JSON string |
| `--limit` | int | `0` | Number of items to fetch |
| `--page-all` | bool | `false` | Fetch all pages |
| `--sandbox` | bool | `false` | Use sandbox environment |

## `mfk authorizations`

Manage authorizations

### `mfk authorizations create`

Create an authorization

| flag | type | default | description |
| --- | --- | --- | --- |
| `--amount` | int | `0` | Amount |
| `--customer-id` | string | — | Customer ID |
| `--destination-id` | string | — | Destination ID |

### `mfk authorizations delete`

Cancel an authorization

```
mfk authorizations delete <id>
```

### `mfk authorizations get`

Get authorization details

```
mfk authorizations get <id>
```

### `mfk authorizations list`

List authorizations

## `mfk billings`

Manage billings (use `qualified` for invoice-system billings since Oct 2023)

Manage billings.

Money Forward Kessai splits billing retrieval into two endpoints around the
Japanese qualified-invoice system (インボイス制度, started Oct 2023):

  mfk billings qualified  -> GET /billings/qualified
      Qualified invoice system (適格請求書等保存方式). Billings from Oct 2023
      onward. THIS IS THE MAIN COMMAND for current billings.

  mfk billings list       -> GET /billings
      Legacy classification-style billings (区分記載請求書等保存方式), roughly
      before Oct 2023. Use this only to access older billings.

Add --page-all to fetch every page (otherwise only the latest 20 are returned).

### `mfk billings download-signed-url`

Get a signed URL to download one invoice PDF

Get a short-lived signed URL for downloading the PDF of one invoice.

The invoice ID is in the billing's invoice_ids (see "mfk billings get" or
"mfk billings qualified"). The response looks like
{"items":[{"signed_url": "...", "expired_at": "...", "type": "pdf"}]}.
Download the file from signed_url without the API key; the URL expires soon.

```
mfk billings download-signed-url <billing_id> <invoice_id>
```

Example:

```
mfk billings download-signed-url <billing_id> <invoice_id>
  curl -sSo invoice.pdf "$(mfk billings download-signed-url <billing_id> <invoice_id> | jq -r '.items[0].signed_url')"
```

### `mfk billings get`

Get billing details

```
mfk billings get <billing_id>
```

### `mfk billings list`

List legacy (pre-Oct-2023) billings — use `qualified` for current invoice-system billings

| flag | type | default | description |
| --- | --- | --- | --- |
| `--customer-id` | string | — | Filter by customer ID |
| `--destination-id` | string | — | Filter by destination ID |
| `--status` | string | — | Filter by status |

### `mfk billings qualified`

List invoice-system billings (適格請求書, since Oct 2023) — the main billings command

| flag | type | default | description |
| --- | --- | --- | --- |
| `--customer-id` | string | — | Filter by customer ID |
| `--destination-id` | string | — | Filter by destination ID |

### `mfk billings reissue`

Reissue invoice

```
mfk billings reissue <billing_id>
```

### `mfk billings upload-signed-url`

Get a signed URL to upload a file for one billing

Get a short-lived signed URL for uploading a file attached to one billing.

--content-type is required by the API and must be one of application/pdf,
application/json, text/csv or text/plain (default: application/pdf).

This endpoint is only available to sellers whose contract includes it; other
accounts get "forbidden_seller" (HTTP 403) even though the request is correct.
Ask Money Forward Kessai to enable it before using this command.

```
mfk billings upload-signed-url <billing_id>
```

Example:

```
mfk billings upload-signed-url <billing_id>
  mfk billings upload-signed-url <billing_id> --content-type text/csv
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `--content-type` | string | `application/pdf` | Content type of the file to upload: application/pdf, application/json, text/csv or text/plain |

## `mfk credit-facilities`

Manage credit facilities

### `mfk credit-facilities get`

Get credit facility details

```
mfk credit-facilities get <id>
```

### `mfk credit-facilities list`

List credit facilities

| flag | type | default | description |
| --- | --- | --- | --- |
| `--customer-examination-id` | string | — | Filter by examination ID |
| `--customer-id` | string | — | Filter by customer ID |

## `mfk customer-examinations`

Manage customer examinations

### `mfk customer-examinations create`

Request a credit examination (also used to request a higher credit line)

Request a credit examination (与信枠審査) for a customer: POST /customer_examinations.

This is also how you request a higher credit line (増枠申請). The examination
finishes within two business days at most.

Required: customer_id, amount, end_date, address1, email, tel, zip_code.

  --amount     The requested credit line in JPY. While a credit line is in
               effect, an amount below it is rejected. The examination may
               grant less than requested.
  --end-date   Transactions can be registered against the line until this
               date. Only the last day of one of the next three months is
               accepted (e.g. 2026-12-31).

With automatic credit examination (自動与信), the raised amount carries over
to later automatic credit lines, and the customer's alert is cleared.

While an examination for the customer is still in progress, a new one is
rejected with HTTP 409.

--from-destination <destination_id> fetches an existing billing destination
(a read-only GET, sent even with --dry-run) and fills customer_id, address1,
address2, zip_code, tel and email from it. Explicit flags win over --json,
and --json wins over the destination. A credit line increase then needs
little more than:

  mfk customer-examinations create --from-destination <destination_id> \
    --amount 3000000 --end-date 2026-12-31 --dry-run

Always check the request with --dry-run first: this sends a real request
to Money Forward Kessai.

| flag | type | default | description |
| --- | --- | --- | --- |
| `--address1` | string | — | Address line 1 (required) |
| `--address2` | string | — | Address line 2 (building, floor) |
| `--amount` | int | `0` | Requested credit amount in JPY (required). Must not be below the credit line currently in effect |
| `--business-description` | string | — | Business description: the customer's main services or products |
| `--business-type` | string | — | Business type: corporate or individual. Leave empty if unknown |
| `--corporate-number` | string | — | Corporate number (法人番号). Only when --business-type corporate |
| `--customer-id` | string | — | Customer ID (required; filled from --from-destination if omitted) |
| `--email` | string | — | Email (required) |
| `--end-date` | string | — | Last day transactions may be registered against this credit line (required, YYYY-MM-DD). Only the last day of one of the next three months is accepted |
| `--from-destination` | string | — | Fill customer_id, address1, address2, zip_code, tel and email from this billing destination (GET /destinations/{id}). Explicit flags and --json win |
| `--remark` | string | — | Other information for the examination. In sandbox, also selects the examination result |
| `--representative-name` | string | — | Representative name |
| `--tel` | string | — | Phone (required) |
| `--website` | string | — | Customer's website URL |
| `--zip-code` | string | — | Zip code (required) |

### `mfk customer-examinations get`

Get customer examination details

```
mfk customer-examinations get <id>
```

### `mfk customer-examinations list`

List customer examinations

| flag | type | default | description |
| --- | --- | --- | --- |
| `--customer-id` | string | — | Filter by customer ID |
| `--status` | string | — | Filter by status |

## `mfk customer-name-updates`

Manage customer name updates

### `mfk customer-name-updates create`

Create a customer name update

| flag | type | default | description |
| --- | --- | --- | --- |
| `--customer-id` | string | — | Customer ID |
| `--name` | string | — | New customer name |

### `mfk customer-name-updates delete`

Cancel a customer name update

```
mfk customer-name-updates delete <id>
```

### `mfk customer-name-updates get`

Get customer name update details

```
mfk customer-name-updates get <id>
```

### `mfk customer-name-updates list`

List customer name updates

## `mfk customers`

Manage customers

### `mfk customers account-transfer`

Apply for account transfer

```
mfk customers account-transfer <customer_id>
```

### `mfk customers bank-transfer`

Assign bank transfer account

```
mfk customers bank-transfer <customer_id>
```

### `mfk customers create`

Create a customer

| flag | type | default | description |
| --- | --- | --- | --- |
| `--dest-address1` | string | — | Destination address1 |
| `--dest-address2` | string | — | Destination address2 |
| `--dest-department` | string | — | Destination department |
| `--dest-email` | string | — | Destination email |
| `--dest-name` | string | — | Destination name |
| `--dest-name-kana` | string | — | Destination name kana |
| `--dest-tel` | string | — | Destination phone |
| `--dest-title` | string | — | Destination title |
| `--dest-zip-code` | string | — | Destination zip code |
| `--name` | string | — | Customer name |
| `--number` | string | — | Customer number |

### `mfk customers get`

Get customer details

```
mfk customers get <customer_id>
```

### `mfk customers list`

List customers

| flag | type | default | description |
| --- | --- | --- | --- |
| `--has-alert` | bool | `false` | Filter by alert status |
| `--name` | string | — | Filter by customer name |
| `--number` | string | — | Filter by customer number |

### `mfk customers update`

Update a customer

```
mfk customers update <customer_id>
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `--name` | string | — | Customer name |

## `mfk destinations`

Manage billing destinations

### `mfk destinations create`

Create a destination

| flag | type | default | description |
| --- | --- | --- | --- |
| `--address1` | string | — | Destination address1 |
| `--address2` | string | — | Destination address2 |
| `--customer-id` | string | — | Customer ID |
| `--department` | string | — | Destination department |
| `--email` | string | — | Destination email |
| `--name` | string | — | Destination name |
| `--name-kana` | string | — | Destination name_kana |
| `--tel` | string | — | Destination tel |
| `--title` | string | — | Destination title |
| `--zip-code` | string | — | Destination zip_code |

### `mfk destinations get`

Get destination details

```
mfk destinations get <destination_id>
```

### `mfk destinations list`

List destinations

| flag | type | default | description |
| --- | --- | --- | --- |
| `--customer-id` | string | — | Filter by customer ID |

### `mfk destinations update`

Update a destination

```
mfk destinations update <destination_id>
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `--address1` | string | — | Destination address1 |
| `--address2` | string | — | Destination address2 |
| `--department` | string | — | Destination department |
| `--email` | string | — | Destination email |
| `--name` | string | — | Destination name |
| `--name-kana` | string | — | Destination name_kana |
| `--tel` | string | — | Destination tel |
| `--title` | string | — | Destination title |
| `--zip-code` | string | — | Destination zip_code |

## `mfk issues`

Manage issued billing information

### `mfk issues list`

List issued billing information

## `mfk payout-refunds`

Manage payout refunds

### `mfk payout-refunds list`

List payout refunds

## `mfk payout-transactions`

Manage payout transactions

### `mfk payout-transactions get`

Get payout transaction details

```
mfk payout-transactions get <id>
```

### `mfk payout-transactions list`

List payout transactions

| flag | type | default | description |
| --- | --- | --- | --- |
| `--payout-id` | string | — | Filter by payout ID |

## `mfk payouts`

Manage payouts

### `mfk payouts get`

Get payout details

```
mfk payouts get <payout_id>
```

### `mfk payouts list`

List payouts

## `mfk transactions`

Manage transactions

### `mfk transactions create`

Create a transaction

| flag | type | default | description |
| --- | --- | --- | --- |
| `--amount` | int | `0` | Transaction amount |
| `--date` | string | — | Transaction date (YYYY-MM-DD) |
| `--delivery-method` | string | — | Invoice delivery method: posting or email |
| `--destination-id` | string | — | Destination ID |
| `--due-date` | string | — | Due date (YYYY-MM-DD) |
| `--issue-date` | string | — | Issue date (YYYY-MM-DD) |
| `--number` | string | — | Transaction number |

### `mfk transactions delete`

Cancel a transaction

```
mfk transactions delete <transaction_id>
```

### `mfk transactions get`

Get transaction details

```
mfk transactions get <transaction_id>
```

### `mfk transactions list`

List transactions

| flag | type | default | description |
| --- | --- | --- | --- |
| `--customer-id` | string | — | Filter by customer ID |
| `--destination-id` | string | — | Filter by destination ID |
| `--status` | string | — | Filter by status |
