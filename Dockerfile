FROM golang:1.20

WORKDIR /bin

COPY . .

RUN go build -o /bin/apiGateway ./cmd/app/main.go

EXPOSE 8080

CMD ["/bin/apiGateway"]