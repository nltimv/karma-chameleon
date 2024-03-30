FROM golang:alpine AS build

WORKDIR /src/
COPY ./slack /src/slack
COPY go.mod go.sum /src/
WORKDIR /src/slack
RUN go build -o /bin/slack

FROM busybox

COPY --from=build /bin/slack /bin/slack
COPY --from=build /etc/ssl/certs /etc/ssl/certs
ENTRYPOINT [ "/bin/slack" ]
