# Stage 1: Build stage
FROM golang:1.23-alpine AS build
# Set the working directory
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code
COPY . .
# Build the Go application
# RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .
RUN go build -o api .

# Stage 2: Final stage
FROM scratch
COPY --from=build /src/api .
ENTRYPOINT ["/api"]
