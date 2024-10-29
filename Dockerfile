# Stage 1: Build stage
FROM golang:1.23-alpine AS builder
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificatesRUN apk add --no-cache build-base=0.5-r3
EXPOSE 3000
# Set the working directory
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code
COPY . .
# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o api -a -ldflags '-linkmode external -extldflags "-static"' .
# Stage 2: Final stage
FROM scratch
WORKDIR  /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/api .
ENTRYPOINT ["/api"]
