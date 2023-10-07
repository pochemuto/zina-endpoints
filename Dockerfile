FROM golang:1.21.2-alpine as build
ARG APP_VERSION=<unset>

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN echo "Building version ${APP_VERSION}";\
    go build -o /main -ldflags="-X 'main.app_version=${APP_VERSION}'"

FROM alpine
COPY --from=build /main /main
EXPOSE 80
CMD [ "/main" ]