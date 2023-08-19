# mumago

A try at realtime sync with json patch in golang.

## Development

1. Generate some local certificates to allow http2 to work on https.

```bash
mkcert -install -cert-file ./cert.pem -key-file ./key.pem localhost
```

2. Create a `.env` file with the following env vars set:

```bash
DATABASE_DSN=""
```

3. Run the application:

```bash
go run main.go
```