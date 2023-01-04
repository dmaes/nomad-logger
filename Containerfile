FROM golang:1.19

ENV APP_ROOT /srv
WORKDIR $APP_ROOT

COPY . $APP_ROOT

RUN go build

CMD ["/srv/nomad-logger"]
