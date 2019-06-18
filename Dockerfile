FROM golang:1.12.5-stretch as builder

RUN mkdir /rest-to-googlesheet
WORKDIR /rest-to-googlesheet

ADD go.mod .
ADD go.sum .

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/rest-to-googlesheet .

FROM alpine

RUN apk add --no-cache ca-certificates openssl

COPY --from=builder /go/bin/rest-to-googlesheet /app/
COPY client_secret.json /app/

WORKDIR /app
CMD ["./rest-to-googlesheet"]