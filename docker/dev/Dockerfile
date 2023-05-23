FROM golang:1.20-alpine
RUN apk update && apk add --no-cache git

RUN go install github.com/cespare/reflex@latest
COPY reflex.conf /
COPY run.sh /

ENTRYPOINT ["reflex", "-c", "/reflex.conf"]