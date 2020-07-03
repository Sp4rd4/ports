FROM golang:1.14.4-alpine3.12 as build

WORKDIR /ports

COPY ./go.* ./
COPY ./vendor ./vendor
COPY ./pkg ./pkg
COPY ./cmd/portdomain/ ./cmd/portdomain

WORKDIR /ports/cmd/portdomain

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o portdomain .

######################################################################

FROM alpine:3.12 as portdomain

RUN apk update
RUN apk add ca-certificates

COPY --from=build /ports/cmd/portdomain/portdomain .
COPY --from=build /ports/pkg/storage/postgres/migrations ./migrations

CMD ["./portdomain"]