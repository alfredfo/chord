FROM golang:1.21-alpine
WORKDIR /app
COPY . /app
RUN go build /app/cmd/chord
ENTRYPOINT ["/app/chord"]