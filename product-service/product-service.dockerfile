FROM alpine:lastest

RUN mkdir -p /app

COPY productApp /app

CMD ["/app/productApp"]
