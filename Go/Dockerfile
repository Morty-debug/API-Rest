FROM golang:1.12-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh mysql-client

LABEL maintainer="Ricardo Antonio Valladares Renderos <r_a_v_r_@hotmail.com>"

WORKDIR /docker

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 5002

CMD ["./main"]
