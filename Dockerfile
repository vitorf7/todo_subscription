FROM golang:1.22-alpine AS build

# Dependencies
RUN apk update && apk add make git gcc musl-dev

# Allow access to private dependencies
ARG GITHUB_TOKEN
ENV GOPRIVATE=github.com/utilitywarehouse/*
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# Add project
ADD . /go/src/github.com/utilitywarehouse/contact-channels-graphql-go-service
WORKDIR /go/src/github.com/utilitywarehouse/contact-channels-graphql-go-service

# Download dependencies
ENV GOPRIVATE=github.com/utilitywarehouse/*
RUN make install

# Build binary
ARG SERVICE
RUN make $SERVICE
RUN mv $SERVICE /$SERVICE

FROM alpine:latest

ARG SERVICE
ENV SERVICE=$SERVICE

# Install ca certs
RUN apk add --no-cache ca-certificates && mkdir /app

# Copy executables from build
COPY --from=build /$SERVICE /app/$SERVICE

ENTRYPOINT /app/$SERVICE
