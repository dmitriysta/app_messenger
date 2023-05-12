
FROM golang:latest

WORKDIR /app/users

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /messenger/message

EXPOSE 8080

CMD [ "/messenger/users" ]