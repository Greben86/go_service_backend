FROM golang:alpine
COPY . /app
WORKDIR /app/src/rest_module

RUN go build -o app *.go

# Открытие порта 8080
EXPOSE 8080

CMD ["./app"]