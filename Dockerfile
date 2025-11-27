FROM golang:1.25.4 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
COPY . . 
RUN swag init -g ./api/main.go -o ./docs

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o shopping-kart ./api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=builder /app/shopping-kart ./
COPY --from=builder /app/data ./data
COPY --from=builder /app/docs ./docs
RUN chmod +x shopping-kart
EXPOSE 8080
CMD ["./shopping-kart"]