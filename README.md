# go-native-api

## Requirements

- Go (go1.24.4 linux/amd64)
- PostgreSQL 15.13

## Installation

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd go-native-api
    ```

2.  **Set up the .env file:**

    Copy the contents of `.env-example` to `.env` and fill in the necessary values.

    ```bash
    cp .env-example .env
    ```

    Edit the `.env` file with your PostgreSQL credentials:

    ```
    DB_HOST=localhost
    DB_USERNAME=zaropheus
    DB_PASSWORD=your_db_password
    DB_NAME=go-native-api
    DB_PORT=5432
    ```

3.  **Run the application:**

    ```bash
    go run cmd/main.go
    ```

**Important:** Ensure you have Go installed and a PostgreSQL database named `go-native-api` created.

Make sure to replace `<repository_url>` with the actual URL of your repository and `your_db_password` with your actual
