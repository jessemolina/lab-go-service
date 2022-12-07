# ================================================================
# BUILD GO BINARY
FROM golang:1.18 as build_service-api
ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /service

WORKDIR /service/app/services/service-api
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# ================================================================
# RUN BINARY IN ALPINE

FROM alpine:3.16
ARG BUILD_DATE
ARG BUILD_REF

RUN addgroup -g 1000 -S service && \
    adduser -u 1000 -h /service -G service -S service

COPY --from=build_service-api --chown=service:service /service/app/services/service-api/service-api /service/service-api

WORKDIR /service
USER service
CMD ["./service-api"]

# ================================================================
# LABEL

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="service-api" \
      org.opencontainers.image.authors="Jesse Molina <jesse.molina@lambdax.io>" \
      org.opencontainers.image.source="https://github.com/jessemolina/ultimate-service/app/services/service-api" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="lambdax.io"
