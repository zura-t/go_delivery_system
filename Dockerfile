FROM alpine:latest

RUN mkdir /app

COPY apiGateway /app

CMD ["/app/apiGateway"]