FROM golang:latest AS builder

WORKDIR /app

COPY repositories/go.mod repositories/go.sum ./

RUN go mod download

COPY repositories .

ARG MAIN_PATH

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ${MAIN_PATH}

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]