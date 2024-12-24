FROM golang:1.23 as build

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go vet -v ./...

RUN go test -v ./...

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/xrdebug

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/xrdebug /

EXPOSE 27420

CMD ["/xrdebug"]

