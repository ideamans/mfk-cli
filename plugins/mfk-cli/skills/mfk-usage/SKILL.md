---
name: mfk-usage
description: Work with Money Forward Kessai (マネーフォワードケッサイ) through the mfk CLI — look up customers and credit examinations, fetch billings and invoices, check transactions, payouts and destinations. Use when the user asks about their MFK account, an invoice or 請求書 issued through Money Forward Kessai, a customer's credit line, or a payout.
license: MIT
compatibility: Requires the `mfk` binary on PATH — run the mfk-install skill if it is missing. Needs MFK_API_KEY (or MFK_SANDBOX_API_KEY with --sandbox). Operates on live billing and credit data for real customers; create and update methods have real financial consequences.
allowed-tools: Bash(mfk:*) Bash(jq:*) Bash(command:*) Read Write
---

# mfk-usage

Drive the Money Forward Kessai API v2 through the `mfk` CLI. Every command has
the shape `mfk <resource> <method> [flags]` and prints JSON by default.

## 1. Confirm the tool and the environment

```bash
command -v mfk && mfk --version
```

Missing? Run the `mfk-install` skill.

`MFK_API_KEY` must be set for production; `--sandbox` switches to the sandbox
base URL and prefers `MFK_SANDBOX_API_KEY`. The environment is decided by
`--sandbox` → `MFK_ENV` → production. **Check which one you are about to hit
before running anything that writes** — the same command against production
touches real customers.

## 2. The mistake that wastes the most time: billings

Money Forward Kessai split invoice retrieval in two when the Japanese invoice
system (インボイス制度, 適格請求書等保存方式) started in October 2023:

| Command | Returns | Period |
| --- | --- | --- |
| `mfk billings qualified` | invoice-system billings | **October 2023 onwards — current** |
| `mfk billings list` | old-format (区分記載請求書) billings | roughly before October 2023 |

**Use `mfk billings qualified` for anything current.** If a user says "I can only
get data up to 2023", the cause is almost always `billings list` being used
instead of `qualified`. Both take the same query parameters and return the same
Billing shape.

## 3. Paginate deliberately

Lists return the **20 most recent** records by default, newest first. Without
`--page-all` you are looking at a truncated view — which silently produces wrong
totals.

```bash
mfk billings qualified --page-all | jq '.items | length'
mfk billings qualified --customer-id cust_xxx --page-all
```

`--limit` sets the page size (1–200). `--page-all` follows
`pagination.has_next` internally and merges every page into `{"items":[...]}`.

## 4. Read the reference for anything else

```bash
mfk llm | head -60
mfk llm | grep -i 'customers'
```

Embedded in the binary (~520 lines), so it matches the installed version.

## 5. Treat writes as financial actions

Creating a customer, requesting a credit examination, issuing or reissuing a
billing, and anything touching payouts **affect a real business relationship and
real money**. Before any non-read method:

```bash
mfk <resource> <method> --json '{...}' --dry-run
```

`--dry-run` prints the request without sending it. Show the user what will be
sent, to which customer, in which environment, and get their agreement. Then run
it for real.

## 6. Report

Give back the identifiers (billing id, customer id) and the amounts and statuses
that the user needs. When you paginated, say how many records the number covers
and over what period.

## Failure modes

| Symptom | Cause | Fix |
| --- | --- | --- |
| `command not found: mfk` | not installed | run the `mfk-install` skill |
| authentication errors | `MFK_API_KEY` unset or wrong environment | check the variable and whether `--sandbox` is intended |
| only pre-2023 invoices come back | used `billings list` | use `mfk billings qualified` |
| totals look too small | default 20-record page | add `--page-all` |
| unexpected write against production | `--sandbox` omitted | verify the environment before every write; use `--dry-run` first |
