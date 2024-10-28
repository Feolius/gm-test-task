FROM golang:alpine as builder
WORKDIR /app
COPY go.mod go.sum  /app/
RUN go mod download
COPY . .
RUN ls
RUN go build -buildvcs=false -v -o app ./cmd/gateway-service

FROM alpine as prod
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/app /app
CMD [ "/app/app" ]