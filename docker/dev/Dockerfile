FROM golang:1.20-alpine
RUN apk update && apk add --no-cache git
RUN go install github.com/cespare/reflex@latest
COPY run.sh /
COPY test.sh /

ENTRYPOINT reflex -r "(\.go$|go\.mod)" -s sh /run.sh