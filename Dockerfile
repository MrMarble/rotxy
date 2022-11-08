ARG GO_VERSION=1.18

## Build container
FROM docker.io/golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /rotxy /src/cmd/rotxy

## Final container
FROM gcr.io/distroless/static:nonroot AS final

COPY --from=builder /rotxy /

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/rotxy","-d"]
