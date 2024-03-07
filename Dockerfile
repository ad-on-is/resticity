From golang:alpine as builder
WORKDIR /build/
COPY . .
RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest
RUN apk add nodejs npm gcc musl-dev pkgconf gtk+3.0-dev webkit2gtk-dev
RUN npm i -g pnpm
RUN wails build


FROM alpine
RUN apk --update add ca-certificates curl mailcap
WORKDIR /
COPY --from=builder /build/build/bin/resticity /resticity

EXPOSE 11278

CMD ["/resticity", "--config", "/config.json", "--server"]