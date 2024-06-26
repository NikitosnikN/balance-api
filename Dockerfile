FROM --platform=linux/amd64 golang:1.22-alpine3.18 as builder

WORKDIR /home/app

COPY . /home/app

RUN go mod download
RUN go build -o /balance-api ./cmd/cli


FROM --platform=linux/amd64 alpine:3.18.4

RUN mkdir /home/app
WORKDIR /home/app

RUN apk --no-cache add gcompat tini
COPY --from=builder /balance-api /home/app/balance-api

# copy default settings
COPY config.example.yaml /home/app/config.example.yaml

# expose default port
EXPOSE 8000

ENTRYPOINT ["/sbin/tini", "--"]
CMD [ "/home/app/headers-api" ]