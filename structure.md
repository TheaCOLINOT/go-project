ecommerce/
│
├── cmd/
│   ├── server/
│   │   └── main.go
│   ├── client/
│   │   └── main.go
│   ├── admin/
│   │   └── main.go
│   ├── ssh-client/
│   │   └── main.go
│   └── ssh-admin/
│       └── main.go
│
├── internal/
│   │
│   ├── domain/
│   │   ├── user/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── errors.go
│   │   ├── product/
│   │   ├── cart/
│   │   ├── order/
│   │   └── payment/
│   │
│   ├── application/
│   │   ├── auth/
│   │   ├── user/
│   │   ├── product/
│   │   ├── cart/
│   │   ├── order/
│   │   └── payment/
│   │
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   └── postgres/
│   │   │       ├── user_repository.go
│   │   │       ├── product_repository.go
│   │   │       ├── cart_repository.go
│   │   │       └── order_repository.go
│   │   ├── security/
│   │   └── email/
│   │
│   ├── delivery/
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   ├── routes.go
│   │   │   └── dto/
│   │   ├── cli/
│   │   │   ├── client/
│   │   │   ├── admin/
│   │   │   └── styles/
│   │   └── ssh/
│   │
│   └── config/
│
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_products.sql
│   ├── 003_create_carts.sql
│   ├── 004_create_orders.sql
│   └── 005_create_payments.sql
│
├── tests/
│   ├── integration/
│   └── e2e/
│
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── .env.example
└── README.md