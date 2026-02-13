# Wallet Service

A simple digital wallet service implemented in Go.

This service allows users to:

- Withdraw funds from their wallet
- Check wallet balance

The system ensures:

- Strongly consistent withdrawals
- Atomic balance updates
- Idempotent request handling
- Ledger-based audit trail

---

## 🚀 Tech Stack

- Go
- PostgreSQL
- sqlx
- Gin
- Goose (database migrations)

---

## 📦 Installation Requirements

Before running the project, ensure you have:

- Go 1.22+
- PostgreSQL 14+
- Goose

---

## 🛠 Setup & Run

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Create PostgreSQL database

If you don't have postgres in your local machine you can:

```bash
docker run --name postgres -e POSTGRES_PASSWORD=12345678 -p 5432:5432 -d postgres    
```

Then run this on sql:

```sql
CREATE DATABASE digital_wallet;
```

### 3. Copy env

Run this script:

```bash
cp .env.example .env
```

and change the db username and password on DB_URL

### 4. Run database migrations

If you don't have goose in your local machine you can install it by:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Then run this on the project

```bash
goose up
```

### 5. Run API

```bash
go run cmd/app/main.go
```

The server will start on the configured port.

---

## 🧪 Running Tests

Integration tests require PostgreSQL.

Then run:

```bash

go test ./internal/service -v
```

Tests include:

- Withdraw success
- Withdraw insufficient funds
- Concurrent withdrawals
- Idempotency replay

---

## 📌 API Routes

### 1. Withdraw

```http

POST /v1/wallets/withdraw
```

#### Idempotency Replay

The Withdraw API uses a unique idempotency_key to prevent double deduction.
If the same key is sent again, the system does not execute the withdrawal twice.
Instead, it returns the previously stored result.

#### CURL

```curl
curl --request POST \
  --url http://localhost:8000/v1/wallets/withdraw \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/12.3.0' \
  --header 'X-Idempotency-Key: test-5' \
  --header 'X-User-ID: 1' \
  --data '{
 "amount": 20000
}'

```

#### Headers

- Idempotency-Key: it's client generated and ideally it would be uuid but it can be anything.
- X-User-ID: 1

#### Request Body

```json
{
  "amount": 30000
}

```

#### Success Response

```json
{
  "userId": 1,
  "amount": 30000,
  "balance": 70000
}
```

#### Error Response

| HTTP | Code                   | Description                                   |
| ---- | ---------------------- | --------------------------------------------- |
| 400  | INVALID_AMOUNT         | Amount must be greater than 0                 |
| 404  | WALLET_NOT_FOUND       | Wallet does not exist                         |
| 409  | INSUFFICIENT_FUND      | Not enough balance                            |
| 409  | IDEMPOTENCY_KEY_REUSED | Idempotency key reused with different payload |
| 409  | REQUEST_IN_PROGRESS    | Previous request still being processed        |
| 500  | UNKNOWN_ERROR          | Unexpected server error                       |

### 2. Balance Inquiry

```http
POST /v1/wallets/balance
```

#### CURL

```curl
curl --request GET \
  --url http://localhost:8000/v1/wallets/balance \
  --header 'User-Agent: insomnia/12.3.0' \
  --header 'X-User-ID: 1'
```

#### Headers

- X-User-ID: 1

#### Success Response

```json
{
  "balance": 70000
}
```

### 3. Create User

```http
POST /v1/users
```

#### CURL

```curl
curl --request POST \
  --url http://localhost:8000/v1/users \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/12.3.0' \
  --data '{
 "name": "Bolang",
 "balance": 200000
}'
```

#### Request Body

```json
{
 "name": "Bolang",
 "balance": 200000
}
```

#### Response Body

```json
{
 "data": {
  "id": 1 // User id for easy testing
 }
}

```

---

## 🏗 Design Decisions

### 1. Atomic Balance Update

Withdrawals use a single SQL statement:

```sql
UPDATE wallets SET balance = balance - $1 WHERE user_id = $2 AND balance >= $1 RETURNING balance;
```

This ensures:

- No race condition
- No negative balance
- Safe concurrent withdrawals
Only one concurrent withdraw will succeed when funds are insufficient.

### 2. Ledger-Based Idempotency

Each withdraw requires an Idempotency-Key.

The system:

- Stores withdraw attempts in ledgers
- Returns the same result if the same key is reused
- Rejects reused keys with different payloads

This guarantees safe retry behavior.

### 3. Money Representation

All monetary values use `int64`.

Reason:

- Avoid floating-point precision issues
- Represent currency in smallest unit

Example:
100000 = Rp 100.000

---

## 📂 Folder Structure

```csharp

cmd/
  └── app/
        └── main.go           # Application entry point
internal/
├── domain/                   # Domain models and business errors
├── repository/               # Database access layer
├── service/                  # Business logic layer
├── handler/                  # HTTP handlers
└── utils/
      └── response/           # Standardized API response helpers
```

### domain/

Contains:

- Wallet model
- Ledger model
- Business errors (ErrInsufficientFund, etc.)
No database or HTTP logic.

### repository/

Responsible for:

- SQL queries
- Atomic balance updates
- Idempotency enforcement
- Transaction integration

Uses sqlx for database interaction.

### service/

Implements:

- Withdraw business flow
- Idempotency conflict handling
- Validation rules

Returns domain-level errors only.

### handler/

Responsible for:

- Parsing HTTP requests
- Extracting parameters and headers
- Mapping domain errors to HTTP responses
- Returning JSON responses
No business logic inside handlers.

### utils/response/

Contains helper utilities for:

- Standardized JSON success responses
- Standardized error responses
- Keeps handlers clean and consistent.

---

## 📝 Assumptions

- Authentication is out of scope
- Only withdrawal operation implemented
- Single currency system
- No overdraft allowed
