# Subscriptions Service

A small Go service for tracking user subscriptions with PostgreSQL storage and REST API

---


## Running with Docker

### 1. Start services

```bash
docker compose up --build
```

This will start:

- PostgreSQL
- Go API service

DB migrations will run automatically

---

### 2. API access

Default:

```
http://localhost:3000
```

Swagger UI:

```
http://localhost:3000/swagger/index.html
```

---

## Environment Variables

Example `.env` provided in .env.example, simply rename it to .env

---

## Example API Requests

### Create Subscription

```bash
curl -X POST http://localhost:3000/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

---

### List Subscriptions

```bash
curl "http://localhost:3000/subscriptions?limit=10&offset=0"
```

---

### Sum of Subscription Prices

```bash
curl "http://localhost:3000/subscriptions/sum?period_start=05-2024&period_end=10-2026"
```


