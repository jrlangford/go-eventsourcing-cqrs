FROM golang:1.20 as builder

WORKDIR /services/service/
COPY go.mod go.sum ./
COPY --from=lib . /lib
RUN go mod download
COPY cmd cmd
COPY internal internal
RUN go build -ldflags "-linkmode 'external' -extldflags '-static'" cmd/main.go

FROM scratch as dev
COPY --from=builder /services/service/main /main
EXPOSE 5051
CMD ["/main"]
