# Simple Banking API

This project implements a simple banking API with basic functionalities such as credit and debit transactions, along with idempotency support to ensure that duplicate or replayed requests do not result in inconsistent transaction states. The application is built in Golang and uses thread-safe mechanisms for processing transactions to ensure consistency, even under high concurrency.

## Features

- **Credit Transaction**: Adds funds to a user's account.
- **Debit Transaction**: Deducts funds from a user's account.
- **Idempotency Support**: Ensures that duplicate requests do not result in multiple executions of the same transaction.
- **Thread-Safe Transactions**: Ensures transactions are atomic and consistent in multi-threaded environments.
- **External Integration**: Supports forwarding transaction information to third-party systems.
- **Error Handling and Logging**: Detailed error handling for different transaction failure scenarios.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/midedickson/simple-banking-app.git
   ```

2. Navigate into the project directory:

   ```bash
   cd simple-banking-app
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Set up your environment variables, which may include database connection strings, third-party service keys, etc.

## Running the Project

1. Start the server:

   ```bash
   go run main.go
   ```

2. The API will be available on `http://localhost:8080`.

## API Endpoints

### Health Check

- **GET** `/`
  - Simple endpoint to verify if the server is running.
  - Response: `"hello, you have reached simple banking api"`

### Create Credit Transaction

- **POST** `/transaction/credit`
  - Requires an `X-Idempotency-Key` in the request header.
  - Body:
    ```json
    {
      "account_id": "string",
      "amount": "float"
    }
    ```
  - Response: Success message or appropriate error (e.g., insufficient funds, duplicate idempotency key).

### Create Debit Transaction

- **POST** `/transaction/debit`
  - Similar to the credit endpoint, but it deducts funds.
  - Requires an `X-Idempotency-Key` in the request header.
  - Body:
    ```json
    {
      "account_id": "string",
      "amount": "float"
    }
    ```

### Request New Idempotency Key

- **GET** `/idempotency`
  - Generates a new idempotency key for making transactions.
  - Response includes the generated `X-Idempotency-Key` in the response header and body.

### Fetch Transaction Details

- **GET** `/transaction/{reference}`
  - Retrieves the details of a specific transaction using the transaction reference.

### Fetch User Account Details

- **GET** `/account/{id}`
  - Fetches the details of a specific user's account by `id`.

## Idempotency and Thread-Safe Transactions

### **Idempotency**

Idempotency is a key feature in ensuring that duplicate requests do not result in the duplication of transactions. This is particularly important in financial applications where retrying or replaying the same request due to network issues or user actions can cause unintended outcomes.

**How Idempotency Works:**

- A client sends a request to the server with an `X-Idempotency-Key` header.
- The server checks if the idempotency key has already been used for a transaction:
  - **If the key is new**: The server processes the request and stores the idempotency key with the transaction status.
  - **If the key is in progress (`PROCESSING`)**: The server rejects the request to avoid duplicate processing.
  - **If the key is successful (`SUCCESS`)**: The server returns the result of the original transaction.
  - **If the key failed (`FAILED`)**: The server allows the request to be retried.

This ensures that even if a client retries a request (e.g., due to a network timeout), the transaction will only be processed once.

### **Thread-Safe Transactions**

In a concurrent environment, multiple transactions might be processed simultaneously, potentially leading to inconsistent account states if proper precautions are not taken. To avoid this, thread safety is implemented to ensure that transactions are atomic and consistent.

- **Atomicity**: Each transaction (credit or debit) is processed in isolation, ensuring that no other transaction interferes with its execution.
- **Locks/Mutexes**: Appropriate locking mechanisms are applied when updating account balances or creating transactions to avoid race conditions.
- **External Transaction Handling**: The system ensures that all updates to accounts and communication with external systems (e.g., third-party transaction processors) are coordinated to prevent issues like double processing or lost updates.

### Benefits:

- **Concurrency**: Multiple requests can be processed at the same time without the risk of data corruption.
- **Consistency**: Account balances remain consistent even when several transactions are processed in parallel.

## Running Tests

Tests are available to ensure that the idempotency, thread safety, and transaction logic are functioning as expected.

To run the tests, use:

```bash
go test ./...
```

This command will run all the tests within the project.

## Contributing

1. Fork the repository.
2. Create a new feature branch (`git checkout -b feature/my-feature`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add new feature'`).
5. Push to your branch (`git push origin feature/my-feature`).
6. Open a pull request to the `main` branch.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.

---

Feel free to reach out with any questions or suggestions for improvements.
