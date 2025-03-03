# Key-Value Store Go Demo

This is a simple demo of a Key-Value store implemented in Go using a PostgreSQL database and the Gin framework for the web server.

## Features:

- Add key-value pairs via a POST request (`/items`).
- Retrieve all key-value pairs via a GET request (`/items`).
- Delete a key-value pair by key via a DELETE request (`/items?key=X`).

## Running the demo

1. Install Go and docker
2. Run the following commands:

```bash
make db-start
make
# visit http://localhost:8080/
```
