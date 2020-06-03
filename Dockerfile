FROM golang:alpine3.12

ENV APP_NAME musik
ENV PORT 5000

EXPOSE ${PORT}

WORKDIR /go/src/${APP_NAME}

COPY . /go/src/${APP_NAME}


RUN go mod download
RUN go build -o ${APP_NAME}

CMD ./${APP_NAME}


