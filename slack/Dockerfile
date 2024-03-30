FROM golang:alpine AS build

WORKDIR /src/
COPY ./slack /src/slack
COPY go.mod go.sum /src/
WORKDIR /src/slack
RUN apk add make
RUN make build OUTPUT_BIN=/bin/slack

FROM busybox

COPY --from=build /bin/slack /bin/slack
COPY --from=build /etc/ssl/certs /etc/ssl/certs
ENTRYPOINT [ "/bin/slack" ]
