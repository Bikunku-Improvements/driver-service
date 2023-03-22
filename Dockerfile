FROM golang:1.19-alpine AS Builder
WORKDIR /app

# Download Depedency
COPY go.* ./
RUN go mod download

# Copy Source Code to Container
COPY . ./

# Build binary
RUN go build -v -o driver-service

FROM alpine:3.9

# Copy binary from build stage
COPY --from=Builder /app/driver-service /driver-service

# Run app
CMD ["./driver-service"]
