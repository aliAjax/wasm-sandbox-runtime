FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/sandboxd ./cmd/sandboxd
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/sandboxd /sandboxd
EXPOSE 8087
USER nonroot:nonroot
ENTRYPOINT ["/sandboxd"]
