# acend.ch website backend

Backend for acend.ch website: Mailchimp newsletter subscription and form submissions

## Configuration

The following env vars need to be set, most of them are sensitive:

```
MAILCHIMP_API_KEY
SMTP_HOST
SMTP_PORT
SMTP_USERNAME
SMTP_PASSWORD
SMTP_FROM
SMTP_TO
```

## Development

* `go get -v ./...`
* `go run ./cmd/backend`
* `go test -v ./...` (no tests yet, sorry...)
* `golangci-lint run`

## Docker

See `Dockerfile` and make sure to set all expected env vars.
