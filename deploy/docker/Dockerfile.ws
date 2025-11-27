FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod ./
COPY internal ./core
COPY cmd/ws ./cmd/ws
RUN go build -o /bin/ws ./cmd/ws

FROM gcr.io/distroless/base-debian12
USER nonroot:nonroot
COPY --from=builder /bin/ws /bin/ws
EXPOSE 8090 9101
ENTRYPOINT ["/bin/ws"]
