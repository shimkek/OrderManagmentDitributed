## Overview
This project is a distributed system for managing orders, stock, payments, and kitchen operations. It uses microservices architecture with gRPC for communication, RabbitMQ for messaging, and MongoDB for data storage. The system is designed to handle order creation, stock validation, payment processing, and kitchen operations efficiently.

---

## Features
- **Order Service**: Handles order creation, validation, and updates.
- **Stock Service**: Manages inventory and validates stock availability.
- **Payment Service**: Processes payments with Stripe and updates order statuses.
- **Kitchen Service**: Handles cooking operations and updates order statuses. 
- **Gateway Service**: Acts as an API gateway for external clients.
- **Tracing and Logging**: OpenTelemetry with Jaeger UI for distributed tracing and structured logging with Zap.

---

## Architecture
The system consists of the following microservices:
1. **Orders Service**:
   - Handles order creation, validation, and updates.
   - Communicates with the Stock Service for stock validation.
   - Publishes events to RabbitMQ for other services.

2. **Stock Service**:
   - Manages inventory and validates stock availability.
   - Provides gRPC endpoints for stock validation and item retrieval.

3. **Payment Service**:
   - Processes payments and updates order statuses.
   - Publishes events to RabbitMQ for other services.

4. **Kitchen Service**:
   - Listens to RabbitMQ for order-paid events.
   - Updates order statuses after cooking.

5. **Gateway Service**:
   - Acts as an API gateway for external clients.
   - Provides HTTP endpoints for order creation and retrieval.

---

## Technologies Used
- **Programming Language**: Go
- **Service discovery**: Consul
- **Communication**: gRPC
- **Message Broker**: RabbitMQ
- **Database**: MongoDB
- **Tracing**: OpenTelemetry and Jaeger UI
- **Logging**: Zap
- **Containerization**: Docker

---

## Setup Instructions

### Prerequisites
- Install Docker and Docker Compose.
- Install Go (version 1.24.3 or higher).
- Install Go Air

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/shimkek/OrderManagmentDitributed.git
   cd OrderManagmentDistributed

2. Install dependencies for each service:
    ```bash
    go get ./...
    ```
3. Start the services using Docker Compose:
   ```bash
   docker-compose up
   ```

4. Run the services with live reload using Air:
    - Start Air in each service directory:
      ```bash
      cd orders && air
      cd stock && air
      cd payments && air
      cd kitchen && air
      cd gateway && air
      ```

5. Run Stripe webhook listener with stripe CLI:
    - Start the webhook listener:
      ```bash
      stripe listen --forward-to localhost:8081/webhook
      ```

6. Access the services:
   - Gateway API: `http://localhost:8080`
   - MongoExpress(database): `http://localhost:8082`
   - RabbitMQ Management UI(message broker): `http://localhost:15672` (guest:guest)
   - Consul UI(service discovery): `http://localhost:8500`
   - Jaeger UI(telemetry): `http://localhost:16686`


---

## API Endpoints

### Gateway Service
#### Create Order
- **URL**: `POST http://localhost:8080//api/customers/{customerID}/orders`
- **Request Body**:
```json
{
  "items": [
    {
      "productID": "1",
      "quantity": 2
    },
    {
      "productID": "2",
      "quantity": 5
    }
  ]
}
```

- **Response**:
```json
{
  "order": {
    "orderID": "6845be6c505134f8c02c2566",
    "customerID": "1",
    "items": [
      {
        "productID": "1",
        "productName": "Cheese",
        "quantity": 2,
        "priceID": "price_1RWg6C07x6KdXnTOR532hPIw"
      },
      {
        "productID": "2",
        "productName": "Chocolate",
        "quantity": 5,
        "priceID": "price_1RT6Lg07x6KdXnTOdZW6J5rc"
      }
    ],
    "status": "pending",
    "PaymentLink": "none"
  },
  "redirectToURL": "http://localhost:8080/success.html?customerID=1&orderID=6845be6c505134f8c02c2566"
}
```

#### Get Order
- **URL**: `http://localhost:8080/api/customers/{customerID}/orders/{orderID}`
- **Response**:
```json
{
  "orderID": "6845be6c505134f8c02c2566",
  "customerID": "1",
  "items": [
    {
      "productID": "1",
      "productName": "Cheese",
      "quantity": 2,
      "priceID": "price_1RWg6C07x6KdXnTOR532hPIw"
    },
    {
      "productID": "2",
      "productName": "Chocolate",
      "quantity": 5,
      "priceID": "price_1RT6Lg07x6KdXnTOdZW6J5rc"
    }
  ],
  "status": "waiting_payment",
  "PaymentLink": "https://checkout.stripe.com/c/pay/{link}"
}
```
## Project Structure
```
OrderManagementDistributed/
├── common/
│   ├── api/                # Protobuf definitions
│   ├── broker/             # RabbitMQ utilities
│   ├── discovery/          # Service discovery utilities
│   ├── tracer.go           # OpenTelemetry setup
├── gateway/
│   ├── gateway/            # Gateway logic
│   ├── public/             # Static files (HTML)
│   ├── main.go             # Gateway entry point
├── orders/
│   ├── gateway/            # Stock gateway logic
│   ├── service.go          # Order service logic
│   ├── store.go            # MongoDB store logic
│   ├── main.go             # Orders entry point
├── stock/
│   ├── store.go            # Stock store logic
│   ├── service.go          # Stock service logic
│   ├── main.go             # Stock entry point
├── payments/
│   ├── processor/          # Payment processors
│   ├── service.go          # Payment service logic
│   ├── main.go             # Payments entry point
├── kitchen/
│   ├── gateway/            # Orders gateway logic
│   ├── amqp_consumer.go    # RabbitMQ consumer logic
│   ├── main.go             # Kitchen entry point
├── docker-compose.yaml     # Docker Compose configuration
└── README.md               # Project documentation
```

---

## Development

### Generate Protobuf Files
Run the following command to regenerate gRPC code inside the "common" directory:
```bash
make gen
```

