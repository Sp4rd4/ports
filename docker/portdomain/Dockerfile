FROM golang:1.14.4-alpine3.12 as build

WORKDIR /go/src/github.com/sp4rd4/ports

COPY ./vendor ./vendor
COPY ./pkg ./pkg
COPY ./cmd/portdomain/ ./cmd/portdomain
WORKDIR ./cmd/portdomain

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o portdomain .

######################################################################

FROM alpine:3.12 as portdomain

RUN apk update
RUN apk add ca-certificates

COPY --from=build /go/src/github.com/sp4rd4/ports/cmd/portdomain/portdomain .
COPY --from=build /go/src/github.com/sp4rd4/ports/pkg/portdomain/store/postgres/migrations ./migrations

CMD ["./portdomain"]