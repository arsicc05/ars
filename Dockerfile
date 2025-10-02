FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM alpine:latest

COPY --from=build /main .


EXPOSE 8000

CMD ["/main"]