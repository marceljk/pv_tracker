# docker build -t pv_tracker_go:v0.1.0-beta --file docker/Dockerfile .
# Build the application from source
FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /pv_tracker ./cmd/main.go

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /pv_tracker /pv_tracker

USER nonroot:nonroot

ENTRYPOINT ["/pv_tracker"]
