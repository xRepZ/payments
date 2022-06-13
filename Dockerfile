FROM golang:1.18

RUN mkdir /app
COPY . /app

WORKDIR /app
RUN go mod download
RUN go build -o /bin/payments cmd/payments/main.go
CMD /bin/payments -c config/config.yaml