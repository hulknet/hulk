FROM golang:latest as builder
RUN mkdir /build
COPY go.mod /build/
COPY go.sum /build/
ADD . /build/
WORKDIR /build
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  -o main ./cmd/srv/
FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
EXPOSE 80
CMD ["./main"]