FROM golang:alpine3.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./  

RUN go mod download

COPY . .

ENV PORT=80

RUN go build -o ./out/dist cmd/currencies/main.go

CMD ./out/dist

EXPOSE 80