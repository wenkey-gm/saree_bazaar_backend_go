# Saree Bazaar Backend

## Overview

This is a backend application for a e-commerce website called Saree Bazaar. The application provides the following features:

## Features

- User authentication & registration
- Saree Bazaar API
- Manage storage (MongoDB)
- Dockerized for deployment
- Unit tests included(Not Dockerized)

## Prerequisites

Before you start, ensure you have the following installed on your local machine:

- **Go**: The backend is written in Go. [Download Go](https://golang.org/dl/)
- **Docker**: For containerization. [Download Docker](https://www.docker.com/products/docker-desktop)
- **Air**: For live reloading. [Install Air](https://github.com/air-verse/air)

## Run Application(Docker)

1. Clone the repository:
   ```bash
   git clone https://github.com/wenkey-gm/saree_bazaar_backend_go
   cd saree_bazaar_backend_go
   ```

2. Build and run the project using Docker:
   ```bash
    docker-compose up --build
    ```
3. Access the application:

   Once the containers are up, the backend will be accessible at `http://localhost:8080`.

## Run Application(Local Machine Setup)

1. Clone the repository:
   ```bash
   git clone https://github.com/wenkey-gm/saree_bazaar_backend_go
   cd saree_bazaar_backend_go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
    go run cmd/app/main.go
   ```

## Environment Variables
The backend requires certain environment variables for proper functioning. These should be set in a .env file.
   ```
    MONGO_URI=mongodb://mongo:27017
    REFRESH_SECRET=refreshsecret
    PRIV_KEY_FILE=private.pem
    PUB_KEY_FILE=public.pem
    ID_TOKEN_EXP=900
    REFRESH_TOKEN_EXP=259200
   ```
These variables are used for:

- **MONGO_URI**: The connection string for the MongoDB database.
- **REFRESH_SECRET**: The secret key for generating refresh tokens.
- **PRIV_KEY_FILE**: The private key file for generating JWT tokens.
- **PUB_KEY_FILE**: The public key file for generating JWT tokens.
- **ID_TOKEN_EXP**: The expiry time for ID tokens.
- **REFRESH_TOKEN_EXP**: The expiry time for refresh tokens.

## Folder Structure

The project is structured as follows:

- **cmd**: Contains the main application code.
- **config**: Contains the configuration files.
- **controllers**: Contains the controller functions.
- **database**: Contains the database configuration.
- **middleware**: Contains the middleware functions.
- **models**: Contains the data models.
- **routes**: Contains the route definitions.
- **utils**: Contains utility functions.
- **tests**: Contains the unit tests.

