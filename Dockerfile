FROM golang:1.23-alpine

WORKDIR /app

COPY  go.mod go.sum ./

RUN apk add --no-cache git && \
    go mod download


COPY . .

RUN go build -o  main .

EXPOSE 8001

CMD [ "./main" ]