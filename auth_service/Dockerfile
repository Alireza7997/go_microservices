FROM golang:1.21.5

WORKDIR /auth_service

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 6001

CMD ["/auth_service/main"]

