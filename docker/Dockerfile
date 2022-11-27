# Build stage
FROM golang:1.19-alpine3.16 AS builder
ARG SOURCE_LOCATION=/app
WORKDIR ${SOURCE_LOCATION}
ADD go.mod go.mod
COPY . .
ENV GO111MODULE on
RUN go mod download
RUN go build -o app main.go

# Run stage
FROM alpine:3.16
ARG BUILDER_SOURCE=/app
LABEL maintainer="Ludin Nento <ludyyn@gmail.com>"
WORKDIR ${BUILDER_SOURCE}
COPY --from=builder ${BUILDER_SOURCE}/app .
EXPOSE 8080
ENTRYPOINT ["./app"]