From golang:alpine as builder
WORKDIR /build/
COPY . .
RUN apk add nodejs npm
RUN npm i -g pnpm
RUN CGO_ENABLED=0 go build ./cmd/server
RUN cd frontend && pnpm install && pnpm build

FROM alpine
RUN apk --update add ca-certificates curl mailcap restic
WORKDIR /
COPY --from=builder /build/server /resticity-server
COPY --from=builder /build/frontend/.output/public /public

EXPOSE 11278

CMD ["/resticity-server", "--config", "/config.json"]