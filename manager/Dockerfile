FROM golang:1.26-alpine AS builder

# Install git as it's often needed for fetching go modules
RUN apk add --no-cache git

WORKDIR /app
COPY . .
# Resolve dependencies and build
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o manager main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/manager .
EXPOSE 8080
ENTRYPOINT ["./manager"]
