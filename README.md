# mumago

A try at realtime sync with json patch in golang.

## Development

Generate some local certificates to allow http2 to work on https.

```bash
mkcert -install -cert-file ./cert.pem -key-file ./key.pem localhost
```

Run the application:

```bash
go run main.go
```