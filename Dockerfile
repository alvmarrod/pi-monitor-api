FROM alpine:3.20.2

ARG BINARY_NAME=pi_monitor_api
ARG BINARY_PORT=8080
ARG VERSION="1.0.0"

COPY src/${BINARY_NAME}_${VERSION} /pi_monitor_api
WORKDIR /

EXPOSE ${BINARY_PORT}

CMD [ "/pi_monitor_api" ]
