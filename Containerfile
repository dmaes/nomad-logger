FROM golang:1.16

ENV APP_ROOT /srv
WORKDIR $APP_ROOT

COPY . $APP_ROOT

RUN go build

CMD ["/srv/nomad-logger"]
