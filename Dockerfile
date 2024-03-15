From golang:alpine as builder
WORKDIR /build/
COPY . .
RUN apk add nodejs npm git
RUN ./build.sh server
RUN ./build.sh frontend

FROM alpine
RUN apk --update add ca-certificates curl mailcap restic
WORKDIR /
COPY --from=builder /build/server /resticity-server
COPY --from=builder /build/frontend/.output/public /public

EXPOSE 11278

CMD ["/resticity-server", "--config", "/config.json"]