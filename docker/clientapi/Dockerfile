FROM golang:1.14.4-alpine3.12 as build

WORKDIR /ports

COPY ./vendor ./vendor
COPY ./pkg ./pkg
COPY ./cmd/clientapi/ ./cmd/clientapi
COPY ./go.* ./

WORKDIR /ports/cmd/clientapi

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o clientapi .

######################################################################

FROM alpine:3.12 as clientapi

RUN apk update
RUN apk add ca-certificates

COPY --from=build /ports/cmd/clientapi/clientapi .

CMD ["./clientapi"]