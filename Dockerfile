ARG GO_VERSION=1.11
ARG APP_NAME="keeprice"
ARG PORT=8976

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git bzr

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

FROM scratch AS final

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app /app

EXPOSE ${PORT}

ENTRYPOINT ["/app"]