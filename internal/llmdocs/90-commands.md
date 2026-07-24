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

Get download signed URL

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

Get upload signed URL

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

Create a customer examination

| flag | type | default | description |
| --- | --- | --- | --- |
| `--address1` | string | — | Address |
| `--amount` | int | `0` | Amount |
| `--business-type` | string | — | Business type |
| `--customer-id` | string | — | Customer ID |
| `--email` | string | — | Email |
| `--end-date` | string | — | End date (YYYY-MM-DD) |
| `--representative-name` | string | — | Representative name |
| `--tel` | string | — | Phone |
| `--zip-code` | string | — | Zip code |

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
