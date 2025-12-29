# FROM golang:1.25.5-alpine3.23

# WORKDIR /app

# COPY . .

# RUN go mod tidy

# RUN go test ./...

# RUN go build -o api .

# EXPOSE 8080

# CMD [ "./api" ]

# No need to run chmod in this case as the binary built by Go is already executable
# RUN chmod +x ./api

# FROM golang:1.25.5-alpine3.23 AS builder
# WORKDIR /app
# COPY . .
# RUN go build -o api .

# FROM alpine:3.23
# WORKDIR /app
# COPY --from=builder /app/api .
# EXPOSE 8080
# CMD ["./api"]

FROM golang:1.25.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 🔥 CI step
RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api .

FROM alpine:3.23

WORKDIR /app
COPY --from=builder /app/api .

EXPOSE 8090
CMD ["./api"]
