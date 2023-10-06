FROM golang:1.21.2-alpine as build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /main

FROM scratch
COPY --from=build /main /main
EXPOSE 80
CMD [ "/main" ]