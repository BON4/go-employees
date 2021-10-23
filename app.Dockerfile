FROM golang:latest as build
ENV APP_HOME /go/go-employees

COPY ./ $APP_HOME

WORKDIR $APP_HOME/cmd/app

RUN CGO_ENABLED=0 GOOS=linux go build -o service_app

FROM alpine as image

COPY --from=build /go/go-employees/cmd/app/service_app service_app
COPY --from=build /go/go-employees/config/docker_debug_conf.yaml config/

CMD ["./service_app"]