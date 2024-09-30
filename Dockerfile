# Stage 1: Build stage
FROM golang:1.23-alpine AS build
RUN apk add --no-cache build-base=0.5-r3
# Set the working directory
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
# Copy the source code
COPY . .
# Build the Go application
# TODO: add cache for quicker builds locally ?
RUN CGO_ENABLED=1 GOOS=linux go build -o api -a -ldflags '-linkmode external -extldflags "-static"' .

# Stage 2: Final stage
FROM scratch
WORKDIR  /
COPY --from=build /src/api .
ENTRYPOINT ["/api"]
