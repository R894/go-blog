FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./app/app

VOLUME [ "/static" ]

EXPOSE 4000

CMD [ "./app/app" ]