FROM golang:alpine
COPY . /app
WORKDIR /app/src/rest_module

RUN go build -o app main.go api.go

# Открытие порта 8081
EXPOSE 8081

CMD ["./app"]