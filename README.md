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
    ALLOWED_ORIGINS=https://munichandrasarees.com,http://localhost:3001
    ADMIN_EMAILS=owner@example.com
    UPLOAD_DIR=./uploads
    PUBLIC_BASE_URL=https://api.example.com
   ```
These variables are used for:

- **MONGO_URI**: The connection string for the MongoDB database.
- **REFRESH_SECRET**: The secret key for generating refresh tokens.
- **PRIV_KEY_FILE**: The private key file for generating JWT tokens.
- **PUB_KEY_FILE**: The public key file for generating JWT tokens.
- **ID_TOKEN_EXP**: The expiry time for ID tokens.
- **REFRESH_TOKEN_EXP**: The expiry time for refresh tokens.
- **ADMIN_EMAILS**: Comma-separated emails that get the admin role and can manage the catalog. Checked at every login, so adding or removing an email takes effect the next time that person signs in.
- **UPLOAD_DIR**: Where photos uploaded through the admin panel are stored (default `./uploads`). Keep it on persistent storage — docker-compose mounts a volume for it.
- **PUBLIC_BASE_URL**: Public URL of this server, used to build photo URLs (e.g. `https://api.example.com`). If unset, it's taken from the request's host, which is fine locally but should be set in production behind a proxy.
- **ALLOWED_ORIGINS**: Comma-separated browser origins allowed to call the API (CORS). Server-to-server calls, such as the storefront's server-side catalog fetch, don't need to be listed.

## Adding sarees (admin panel)

Open **`/admin`** on the backend (e.g. `http://localhost:8080/admin`) to add, edit, hide and delete sarees from a phone or computer — no API calls needed.

1. Put your email in `ADMIN_EMAILS` and restart the server.
2. Create your account once (only needed the first time):
   ```bash
   curl -X POST http://localhost:8080/signup -H 'Content-Type: application/json' \
     -d '{"email":"owner@example.com","password":"choose-a-strong-password"}'
   ```
3. Sign in at `/admin`. Tap **Add saree**, take or choose photos (the first is the cover; large photos are shrunk to 1600 px before upload), fill in the name, price and stock, and save. Untick **Show on the website** to hide a saree without deleting it.

The storefront picks up changes within 5 minutes. Sessions last 15 minutes; if one expires while you're editing, sign in again and your form is still there.

## Frontend integration

The storefront ([munichandra-sarees-web](https://github.com/wenkey-gm/munichandra-sarees-web)) reads its catalog from this API. Set `SAREE_API_URL` there to this server's base URL (e.g. `http://localhost:8080`). The `/collections` page then lists sarees from `GET /sarees`, and falls back to its built-in product list if the API is unreachable.

| Method | Path | Auth |
|---|---|---|
| `GET` | `/sarees`, `/sarees/:id` | Public |
| `POST` | `/sarees` | Admin (`Authorization: Bearer <access_token>`) |
| `PUT`, `DELETE` | `/sarees/:id` | Admin |
| `POST` | `/uploads` | Admin; multipart `file` field (JPEG, PNG, WebP or GIF, max 15 MB); returns `{"url": ...}` |
| `GET` | `/uploads/:name` | Public |
| `POST` | `/signup`, `/login` | Public; both return `tokens` and `user` (with `role`) |
| `DELETE` | `/signout` | `Authorization: Bearer <access_token>` |

A saree's `name`, `description`, `category`, `color` and `images_url` are what the storefront displays. `name` is required; `id` is generated on create. "Admin" means a signed-in user whose email is in `ADMIN_EMAILS`; other users get `403`.

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

