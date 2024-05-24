FROM alpine:latest
LABEL org.opencontainers.image.source=https://github.com/booleworks/logicng-service

ENV TIMEOUT 5s

ARG TARGETPLATFORM

COPY build/$TARGETPLATFORM/logicng-service /opt/logicng-service

EXPOSE 8080
CMD /opt/logicng-service --timeout $TIMEOUT
