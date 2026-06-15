# syntax=docker/dockerfile:1.7

# ---- build stage ----
FROM golang:1.23-alpine AS build
WORKDIR /src

# Use a mirror reachable from Iran. Override at build time with --build-arg GOPROXY=...
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY} \
    GOSUMDB=off \
    CGO_ENABLED=0

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /out/server ./cmd/server

# ---- runtime stage ----
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /out/server /app/server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/server"]
