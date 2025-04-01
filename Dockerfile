FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /pdkm_project2
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o ./pdkm_project2 .

ENTRYPOINT [ "./pdkm_project2" ]

