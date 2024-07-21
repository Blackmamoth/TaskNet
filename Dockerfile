FROM golang:1.22.5-alpine

WORKDIR /usr/src/app

RUN apk update && apk upgrade

RUN apk add --no-cache make

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN make build

EXPOSE 5500

CMD ["make", "run"]