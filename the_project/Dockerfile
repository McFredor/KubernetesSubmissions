FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go index.html ./
RUN CGO_ENABLED=0 go build -o /app .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app /app
ENTRYPOINT ["/app"]
