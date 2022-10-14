FROM golang:1.18.7-alpine3.16 AS builder

RUN go version
RUN apk add git

COPY ./ /github.com/indikator/aggregator_lets_go
WORKDIR /github.com/indikator/aggregator_lets_go

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/webservice/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root

COPY --from=0 /github.com/indikator/aggregator_lets_go/.bin/app .

CMD [ "./app"]