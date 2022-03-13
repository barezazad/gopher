# Build stage
FROM golang:1.18-rc-alpine3.15 as builder

RUN apk --no-cache add tzdata

WORKDIR /backend

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/app/main.go

FROM scratch as final

COPY --from=builder /backend/server .
COPY --from=builder /backend/storage/ /storage/
COPY --from=builder /backend/logs/   /logs/

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Baghdad

ENTRYPOINT ["/server"]