FROM golang:alpine AS build

WORKDIR /src/
COPY ./ /src/
WORKDIR /src/cmd/slack
RUN apk add make git
RUN make build OUTPUT_BIN=/bin/slack

FROM busybox

COPY --from=build /bin/slack /bin/slack
COPY --from=build /etc/ssl/certs /etc/ssl/certs
ENTRYPOINT [ "/bin/slack" ]
