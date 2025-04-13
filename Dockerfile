FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /order-management
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o ./order-management .

# Final image
FROM alpine:3.19
WORKDIR /order-management

COPY --from=builder /order-management/order-management .
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]