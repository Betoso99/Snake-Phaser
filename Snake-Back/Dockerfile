FROM golang

WORKDIR /cmd

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 3000

CMD [ "go","run","cmd/main.go" ]

