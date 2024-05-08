
FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY tests ./tests
COPY web ./web
COPY *.go ./

ENV TODO_PORT="7540"
ENV TODO_DBFILE="/app/scheduler.db"

EXPOSE 7540

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /yet_another_todo_list

CMD ["/yet_another_todo_list"]
