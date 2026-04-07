FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.work go.work.sum ./
COPY tui/go.mod tui/go.sum ./tui/

WORKDIR /src/tui
RUN go mod download

COPY tui/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/tui ./cmd/tui

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /out/tui /usr/local/bin/tui

ENTRYPOINT ["/usr/local/bin/tui"]
