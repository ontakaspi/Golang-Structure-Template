FROM golang:1.18-stretch AS builder
RUN go env -w GO111MODULE=on
WORKDIR /app
COPY ./.env ./
COPY ./_public_key.pem ./
COPY ./    ./
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.15 as sca-app
USER root
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash
RUN apk add git
RUN apk add --no-cache tzdata
WORKDIR /root/
COPY --from=builder /app/main ./
COPY --from=builder /app/.env ./
CMD ["./main"]
