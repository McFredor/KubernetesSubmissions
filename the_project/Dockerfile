FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go index.html ./
RUN CGO_ENABLED=0 go build -o /app .

FROM scratch
COPY --from=build /app /app
ENTRYPOINT ["/app"]
