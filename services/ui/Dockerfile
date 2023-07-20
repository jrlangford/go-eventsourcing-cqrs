FROM golang:1.20 as builder

WORKDIR /services/service/
COPY go.mod go.sum ./
COPY --from=lib . /lib
RUN go mod download
COPY cmd cmd
COPY internal internal
COPY web web
RUN go build -ldflags "-linkmode 'external' -extldflags '-static'" cmd/main.go

FROM scratch as dev
COPY --from=builder /services/service/main /main
COPY --from=builder /services/service/web /web
EXPOSE 8088
CMD ["/main"]
