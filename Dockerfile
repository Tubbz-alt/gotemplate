FROM semior/baseimage:latest
LABEL maintainer="Semior <ura2178@gmail.com>"

WORKDIR /srv

ENV GOFLAGS="-mod=vendor"

COPY ./app /srv/app
COPY ./vendor /srv/vendor
COPY ./go.mod /srv/go.mod
COPY ./go.sum /srv/go.sum

COPY ./.git/ /srv/.git

RUN \
    export version="$(git describe --tags --long)" && \
    echo $version && \
    go build -mod=vendor -o /go/build/app -ldflags "-X 'main.version=${version}' -s -w" /srv/app

COPY ./migrations /srv/migrations

# cd / is added due to unability to install user-wide packages in the application context
RUN \
    cd / && \
    go get -u github.com/pressly/goose/cmd/goose && \
    cd /srv/migrations

COPY ./scripts/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 2345

CMD ["/entrypoint.sh"]