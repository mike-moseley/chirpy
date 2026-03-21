# API Reference

## Authentication

Request body:
``` json

```
Response body:
``` json

```

### `POST /api/login`
Request body:
``` json
{
    "password": "<password>",
    "email": "<email>"
}
```
Response body:
``` json
{
    "id": "<user_db_uuid>",
    "email": "<email>",
    "token": "<access_token>",
    "refresh_token": "<refresh_token>",
    "is_chirpy_red": "<premium_status>",
}
```

### `POST /api/refresh`
Request headers:
``` json

```
Response body:
``` json

```

<!-- Request headers, response body, error codes -->

### `POST /api/revoke`
Request body:
``` json

```
Response body:
``` json

```

<!-- Request headers, response, error codes -->

## Users

### `POST /api/users`
Request body:
``` json

```
Response body:
``` json

```

<!-- Request body, response body, error codes -->

## Chirps

### `GET /api/chirps`
Request body:
``` json

```
Response body:
``` json

```

<!-- Query params, response body, error codes -->

### `GET /api/chirps/{chirpID}`
Request body:
``` json

```
Response body:
``` json

```

<!-- Path params, response body, error codes -->

### `POST /api/chirps`
Request body:
``` json

```
Response body:
``` json

```

<!-- Request headers, request body, response body, error codes -->

## Admin

### `GET /admin/metrics`

<!-- Response format -->

### `POST /admin/reset`

<!-- What it resets, response -->

## Health

#### `GET /api/healthz`

<!-- Response -->
