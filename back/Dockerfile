FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.26.1 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -a -installsuffix cgo -ldflags="-w -s" -o server ./main.go

FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine:3 AS certs
RUN apk --no-cache add ca-certificates

FROM --platform=${BUILDPLATFORM:-linux/amd64} scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/server /server

EXPOSE 8888

ENTRYPOINT ["/server"]
