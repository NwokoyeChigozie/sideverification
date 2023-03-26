# Build stage
FROM golang:1.20.1-alpine3.17 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .


RUN mv app-sample.env app.env && \
    go build -v -o /dist/vesicash-verification-ms

# Deployment stage
FROM alpine:3.17

WORKDIR /usr/src/app

COPY --from=build /usr/src/app ./

COPY --from=build /dist/vesicash-verification-ms /usr/local/bin/vesicash-verification-ms

CMD ["vesicash-verification-ms"]
