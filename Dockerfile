FROM golang:1.22.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY tests ./tests
COPY web ./web
COPY cmd ./cmd
COPY internal ./internal
COPY config ./config

EXPOSE 7540

RUN GOOS=linux go build -o app ./cmd/yet-another-todo-list/main.go

CMD ["./app"]