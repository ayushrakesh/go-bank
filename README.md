# GoPay: Dockerized Banking API with Transaction Safety

This project implements a simple banking API using Go, Gin-Gonic for the web framework, PostgreSQL as the database, and `sqlc` for type-safe SQL queries. It provides functionalities for managing user accounts, recording entries, and facilitating money transfers. Authentication is handled using PASETO tokens.

## Features

*   **User Management**: Register new users, authenticate existing users.
*   **Account Management**: Create, view, and list bank accounts for authenticated users.
*   **Money Transfers**: Securely transfer funds between user accounts.
*   **Transaction History**: Record all financial entries related to accounts.
*   **Authentication**: Secure API access using PASETO tokens.
*   **Data Validation**: Robust request validation using `gin-gonic/gin/binding` and `go-playground/validator/v10`.
*   **Database**: PostgreSQL with `sqlc` for generated, type-safe database access.
*   **Containerization**: Docker and Docker Compose for easy setup and deployment.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Go**: Version 1.18 or higher. [Download here](https://golang.org/dl/)
*   **PostgreSQL**: Database server. [Download here](https://www.postgresql.org/download/)
*   **Docker & Docker Compose**: (Optional, but recommended for easy setup). [Download here](https://www.docker.com/products/docker-desktop/)
*   **`golang-migrate`**: Database migration tool.
    ```bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```
*   **`sqlc`**: SQL code generation tool.
    ```bash
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    ```
*   **`curl`**: Command-line tool for making HTTP requests (usually pre-installed on Linux/macOS).

## Getting Started

Follow these steps to set up and run the GoPay project from scratch.

### 1. Clone the Repository

```bash
git clone https://github.com/ayushrakesh/gopay.git
cd gopay
```

### 2. Environment Variables (`app.env`)

Create an `app.env` file in the root directory of the project. This file will hold your application's configuration. A typical `app.env` file might look like this:

```env
SYMMETRIC_KEY=a8005714b619fd6efb5481145b61b6da # Keep this key secret and generate a strong one for production
DB_DRIVER=postgres
DB_SOURCE=postgres://postgres:postgres@localhost:5432/bank?sslmode=disable
ACCESS_TOKEN_DURATION=15m
SERVER_ADDRESS=0.0.0.0:8080
```

*   **`SYMMETRIC_KEY`**: A secret key used for signing authentication tokens. **Change this to a strong, random key in a production environment.**
*   **`DB_DRIVER`**: The database driver (e.g., `postgres`).
*   **`DB_SOURCE`**: The database connection string. Adjust `username`, `password`, `host`, `port`, and `database name` as per your PostgreSQL setup.
*   **`ACCESS_TOKEN_DURATION`**: The duration for which access tokens are valid (e.g., `15m` for 15 minutes).
*   **`SERVER_ADDRESS`**: The address and port where the API server will run.

### 3. Database Setup (PostgreSQL)

#### A. Install PostgreSQL

Ensure PostgreSQL is installed on your system. Refer to the [Prerequisites](#prerequisites) section for download links and basic installation instructions.

#### B. Create Database and User

By default, the project expects a database named `bank` and a user named `postgres` with the password `postgres`. You can create these using `psql` (the PostgreSQL command-line client):

```bash
# Connect to PostgreSQL as the default superuser
psql -U postgres

# Inside psql, create the user and database
CREATE USER postgres WITH PASSWORD 'postgres';
CREATE DATABASE bank OWNER postgres;
\q # Exit psql
```

**Important**: If you choose different credentials, make sure to update the `DB_SOURCE` in your `app.env` file accordingly.

#### C. Apply Database Migrations

This project uses `golang-migrate` to manage database schema changes. With your database set up, apply the initial migrations:

```bash
# From the project root directory
make migrateup1
# or, using the migrate tool directly:
migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable" -verbose up
```

### 4. Generate SQLC Code

This project uses `sqlc` to generate type-safe Go code from your SQL queries. Generate the necessary Go code for database interactions:

```bash
make sqlc
```

## Running the Application

You can run the application using Docker (recommended) or by running the Go server directly.

### Option 1: Using Docker (Recommended)

This method uses Docker Compose to set up both the PostgreSQL database and the Go API server in isolated containers, simplifying the setup process.

1.  **Ensure Docker is installed**: Make sure you have [Docker and Docker Compose](https://www.docker.com/products/docker-desktop/) installed on your system.

2.  **Navigate to Project Root**: Open your terminal in the root directory of the `gopay` project.

3.  **Build and Run**: Execute the following command to build the Docker images and start the services:
    ```bash
    docker-compose up -d --build
    ```
    *   `-d`: Runs containers in detached mode (in the background).
    *   `--build`: Ensures that Docker images are rebuilt, picking up any changes to `Dockerfile` or source code.

4.  **Verify Status**: You can check the status of your running containers with:
    ```bash
    docker-compose ps
    ```

5.  **Access API**: The API server will be accessible at `http://localhost:8080`.

6.  **Stop Containers**: To stop and remove the containers, run:
    ```bash
    docker-compose down
    ```

### Option 2: Running Manually

If you prefer to run the application directly on your host machine without Docker:

1.  **Start PostgreSQL**: Ensure your local PostgreSQL server is running.

2.  **Run the Go Server**: Navigate to the project's root directory in your terminal and start the API server:
    ```bash
    make server
    ```
    This command executes `go run main.go`.

3.  **Access API**: The API server will be running on the address specified in `SERVER_ADDRESS` in your `app.env` file (by default, `http://localhost:8080`).

## API Testing

There are two main ways to test the API:

### 1. Running Existing Go Tests

The project includes a comprehensive suite of unit and integration tests. It's recommended to run these to ensure all functionalities are working as expected.

```bash
# From the project root directory
make test
```

This command executes `go test -v -cover ./...`, which will run all tests, provide verbose output, and show code coverage.

### 2. Manual API Testing with `curl`

Ensure the API server is running (either via Docker Compose or `make server`) before proceeding with manual tests.

#### A. Unauthenticated Endpoints

These endpoints do not require an authentication token.

*   **Create User (`POST /users`)**

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{ "username": "testuser", "password": "password", "full_name": "Test User", "email": "test@example.com" }' http://localhost:8080/users
    ```

*   **Login User (`POST /users/login`)**

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{ "username": "testuser", "password": "password" }' http://localhost:8080/users/login
    ```

    **Important**: Copy the `access_token` from the response of the login request. You will need it for all authenticated requests. For convenience, you can store it in an environment variable:
    ```bash
    export AUTH_TOKEN="YOUR_ACCESS_TOKEN_HERE" # Replace with your actual token
    ```

#### B. Authenticated Endpoints

These endpoints require an `Authorization` header with a valid `Bearer` token.

*   **Create Account (`POST /accounts`)**

    ```bash
    curl -X POST -H "Content-Type: application/json" \
         -H "Authorization: Bearer $AUTH_TOKEN" \
         -d '{ "currency": "USD" }' http://localhost:8080/accounts
    ```

*   **Get Account (`GET /accounts/:id`)**

    ```bash
    curl -X GET -H "Authorization: Bearer $AUTH_TOKEN" \
         http://localhost:8080/accounts/[ACCOUNT_ID]
    ```

*   **List Accounts (`GET /accounts`)**

    ```bash
    curl -X GET -H "Authorization: Bearer $AUTH_TOKEN" \
         "http://localhost:8080/accounts?page_id=1&page_size=5"
    ```

*   **Create Transfer (`POST /transfers`)**

    ```bash
    curl -X POST -H "Content-Type: application/json" \
         -H "Authorization: Bearer $AUTH_TOKEN" \
         -d '{ "from_account_id": [FROM_ACCOUNT_ID], "to_account_id": [TO_ACCOUNT_ID], "amount": 100, "currency": "USD" }' http://localhost:8080/transfers
    ```
    *   Replace `[FROM_ACCOUNT_ID]` and `[TO_ACCOUNT_ID]` with valid account IDs.
    *   Ensure the `currency` matches the currency of the accounts involved in the transfer.

## Project Structure

The project is organized into the following key directories and files:

*   **`/api`**: Contains the HTTP server, API handlers for accounts, transfers, and users, and middleware for authentication.
*   **`/db`**:
    *   **`/db/migrations`**: SQL migration files for database schema management.
    *   **`/db/queries`**: Raw SQL query files used by `sqlc`.
    *   **`/db/sqlc`**: Generated Go code for interacting with the PostgreSQL database using `sqlc`.
    *   **`/db/mock`**: Mock implementations for testing database interactions.
*   **`/token`**: Handles PASETO token creation, verification, and payload management for authentication.
*   **`/util`**: Contains utility functions for configuration loading, currency validation, password hashing, and random data generation.
*   **`main.go`**: The application's entry point, responsible for initializing the server and database.
*   **`app.env`**: Environment variable configuration file.
*   **`Dockerfile`**: Defines the Docker image for the Go API server.
*   **`docker-compose.yaml`**: Orchestrates the Docker containers for the API and PostgreSQL database.
*   **`Makefile`**: Contains various commands for building, testing, database migrations, and SQLC generation.
*   **`sqlc.yaml`**: Configuration file for `sqlc`.

## Technologies Used

*   **Go**: Programming language.
*   **Gin-Gonic**: HTTP web framework.
*   **PostgreSQL**: Relational database.
*   **`sqlc`**: Generates type-safe Go code from SQL.
*   **PASETO**: Platform-Agnostic Security Tokens for authentication.
*   **`golang-migrate`**: Database migration management.
*   **`viper`**: For configuration management.
*   **`bcrypt`**: For password hashing.
*   **Docker & Docker Compose**: For containerization and orchestration.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
