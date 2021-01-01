FROM golang:alpine AS build

WORKDIR /app
COPY . .

RUN go get -v ./...
RUN go build

FROM alpine
COPY --from=build /app/backend /usr/local/bin

EXPOSE 8000

ENV SMTP_FROM="hello@acend.ch"
ENV SMTP_TO="hello@acend.ch"

# Furthermore, this image requires the following settings as env vars
# MAILCHIMP_API_KEY
# SMTP_HOST
# SMTP_PORT
# SMTP_USERNAME
# SMTP_PASSWORD

CMD ["backend"]
