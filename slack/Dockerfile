FROM golang:alpine as build

WORKDIR /src/
COPY ./slack /src/slack
COPY go.mod go.sum /src/
WORKDIR /src/slack
RUN go build -o /bin/slack

FROM alpine:latest

COPY --from=build /bin/slack /bin/slack
ENTRYPOINT [ "/bin/slack" ]
