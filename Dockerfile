FROM golang:1.23.0-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o bet-sim main.go

ARG INTERNAL_PORT=3000
ENV INTERNAL_PORT=${INTERNAL_PORT}
EXPOSE ${INTERNAL_PORT}

CMD ["sh", "-c", "./bet-sim"]
