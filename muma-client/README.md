# muma-client 

A svelte implementation of a fun collaborative todolist.

## Development

1. Generate some local certificates to allow http2 to work on https.

```bash
mkcert -install -cert-file ./cert.pem -key-file ./key.pem localhost
```