# syntax=docker/dockerfile:1

FROM golang:1.21.1 as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download


COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /notise ./cmd/ 

FROM scratch
WORKDIR /app
COPY --from=builder /notise /app/

# Run
CMD ["/app/notise"]