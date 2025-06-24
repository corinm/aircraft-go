FROM golang:1.23.10-alpine3.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./discoverer

CMD [ "./discoverer" ]
