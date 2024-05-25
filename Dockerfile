FROM golang:1.22.3
LABEL org.opencontainers.image.source=https://github.com/booleworks/logicng-service

ENV TIMEOUT 30s

WORKDIR /app

COPY go.mod go.sum main.go ./
RUN mkdir ./computation
COPY computation/ ./computation/
RUN mkdir ./config
COPY config/ ./config/
RUN mkdir ./middleware/
COPY middleware/ ./middleware/
RUN mkdir ./sio/
COPY sio/ ./sio/
RUN mkdir ./slog/
COPY slog/ ./slog/
RUN mkdir ./srv/
COPY srv/ ./srv/

RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN CGO_ENABLED=0 GOOS=linux go build -o /logicng-service main.go

EXPOSE 8080

CMD /logicng-service --timeout $TIMEOUT
