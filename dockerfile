FROM golang:1.15.2-alpine
ADD . /app
WORKDIR /app
COPY . .
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o main main.go
EXPOSE 1323
# CMD [ "/app/main" ]

CMD [ "go", "run", "main.go" ]
