FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /notifier

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 3000

CMD ["/notifier"]