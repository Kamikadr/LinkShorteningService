FROM golang:1.22.1-alpine as builder
WORKDIR /build
COPY app/go.mod app/go.sum ./
RUN go mod download
COPY app .
RUN go build -o ./bin/app cmd/Rest-shortcut/main.go

FROM alpine:latest as runner
COPY --from=builder build/bin/app .
COPY config/local.yaml /local.yaml
EXPOSE 8080
ENTRYPOINT ["./app"]