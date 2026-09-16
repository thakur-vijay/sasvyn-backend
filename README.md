# sasvyn-backend

Minimal Go HTTP backend.

## Run

```sh
go run ./cmd/server
```

The server listens on `http://localhost:8080` by default. Set `PORT` to use a different port.

Check the service with:

```sh
curl http://localhost:8080/healthz
```

Run tests with:

```sh
go test ./...
```