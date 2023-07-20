FROM golang:1.19 AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build -o main cmd/api/main.go

FROM alpine:3.13

RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime
RUN echo "Asia/Bangkok" >  /etc/timezone

WORKDIR /usr/src/app

COPY --from=builder /src/main /usr/src/app/main

EXPOSE 8080
CMD ["./main"]
