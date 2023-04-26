FROM golang:1.20

WORKDIR /app
COPY main.go go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /sample-scheduler-extender

EXPOSE 8080

CMD ["/sample-scheduler-extender"]
