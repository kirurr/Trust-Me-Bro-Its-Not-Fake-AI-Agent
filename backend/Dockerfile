FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.work go.work.sum ./
COPY tui/go.mod tui/go.sum ./tui/
COPY shared/go.mod shared/go.sum ./shared/
COPY backend/go.mod backend/go.sum ./backend/

WORKDIR /src
RUN go mod download

COPY shared/ ./shared/
COPY backend/ ./backend/

WORKDIR /src/backend
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/backend ./cmd/server

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

COPY --from=build /out/backend /usr/local/bin/backend

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/backend"]
