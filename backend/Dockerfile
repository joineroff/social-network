FROM golang:1.17-alpine AS builder

ENV APPDIR $GOPATH/src/app
ENV ARTIFACT /build/app

RUN mkdir -p ${APPDIR}
WORKDIR ${APPDIR}
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${ARTIFACT} main.go

FROM alpine

ENV ARTIFACT /build/app

RUN mkdir -p /app
WORKDIR /app

COPY --from=builder ${ARTIFACT} /app/bin
COPY ./config /app/config/
COPY ./templates /app/templates

ENTRYPOINT [ "/app/bin" ]
