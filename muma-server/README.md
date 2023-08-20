# muma-server

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

## Todo

- [ ] Look into idiomatic go error handling
- [ ] Error handling should include some sort of common response type for errors
- [ ] Consider using realtime as a way to create a `RealtimeApiResponse` struct or something
- [ ] Consider moving Realtime init to TodosApi init
- [ ] Look into ways to create a session for the todos. Some sort of unique identifier for a todolist session 