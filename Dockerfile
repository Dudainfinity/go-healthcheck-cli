# --- build: compila o binário estático em Go ---
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /healthcheck .

# --- runtime: imagem mínima só com o binário ---
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /healthcheck /usr/local/bin/healthcheck
COPY config.example.json /config.json
ENTRYPOINT ["healthcheck"]
CMD ["-config", "/config.json"]
