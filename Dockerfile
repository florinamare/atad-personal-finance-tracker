FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILT_AT=unknown
RUN CGO_ENABLED=0 go build \
    -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.builtAt=${BUILT_AT}" \
    -o /app .

FROM alpine:3.20
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]