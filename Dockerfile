FROM alpine:latest

ENV TIMEOUT 5s

COPY logicng-service-linux /opt/logicng-service

EXPOSE 8080
CMD /opt/logicng-service --timeout $TIMEOUT
