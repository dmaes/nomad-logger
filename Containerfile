FROM golang:1.20

ENV APP_ROOT /srv
WORKDIR $APP_ROOT

COPY . $APP_ROOT

RUN go build -o nomad-logger

CMD ["/srv/nomad-logger"]
