FROM node:18.16.0 as webui
ENV WEBUI_DIR /src/webui

RUN mkdir -p $WEBUI_DIR
COPY views/ $WEBUI_DIR/
WORKDIR $WEBUI_DIR

RUN yarn install --silent
RUN yarn build


FROM golang:1.20-alpine as builder
RUN apk update && apk add --no-cache git
RUN go install github.com/cespare/reflex@latest

RUN rm -rf /ostraka/views/dist/
COPY --from=webui /src/webui/dist/ /ostraka/webui/dist/

COPY scripts/run.sh /
COPY scripts/test.sh /

WORKDIR /ostraka
COPY go.* /
RUN go mod download

EXPOSE 4000
ENTRYPOINT reflex -r "(\.go$|go\.mod)" -s sh /run.sh