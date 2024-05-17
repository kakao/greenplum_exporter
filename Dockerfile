FROM golang:1.16-alpine

RUN mkdir -p /app
WORKDIR /app
ENV TZ=Asia/Seoul

COPY go.sum main.go go.mod .
COPY collector collector
RUN mkdir /app/bin

ARG ARCH="amd64"
ARG OS="linux"

RUN go build -o /app/bin/greenplum_exporter
ENTRYPOINT ["sh"]
CMD ["entrypoint.sh"]
