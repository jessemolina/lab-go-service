# ================================================================
# BUILD GO BINARY
FROM golang:1.18 as build_sales-api
ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /service

WORKDIR /service/app/services/sales-api
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# ================================================================
# RUN BINARY IN ALPINE

FROM alpine:3.16
ARG BUILD_DATE
ARG BUILD_REF

RUN addgroup -g 1000 -S sales && \
    adduser -u 1000 -h /service -G sales -S sales

COPY --from=build_sales-api --chown=sales:sales /service/app/services/sales-api/sales-api /service/sales-api

WORKDIR /service
USER sales
CMD ["./sales-api"]

# ================================================================
# LABEL

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="Jesse Molina <jesse.molina@lambdax.io>" \
      org.opencontainers.image.source="https://github.com/jessemolina/ultimate-service/app/sales-api" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="lambdax.io"