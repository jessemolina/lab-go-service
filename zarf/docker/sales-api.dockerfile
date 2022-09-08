# ================================================================
# BUILD GO BINARY
FROM golang:1.18 as build_sales-api
ENV CGO_ENABLE 0
ARG BUILD_REF

COPY . /service

WORKDIR /service
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# ================================================================
# RUN BINARY IN ALPINE

FROM alpine:3.16
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=build_sales-api /service /service/service
WORKDIR /service
CMD ["./service"]

# ================================================================
# LABEL

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="Jesse Molina <jesse.molina@lambdax.io>" \
      org.opencontainers.image.source="https://github.com/jessemolina/ultimate-service/app/sales-api" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="lambdax.io"
