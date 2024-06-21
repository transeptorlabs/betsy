# Build 4337-in-a-box in a container
FROM golang:1.22-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers git

# Get dependencies
COPY go.mod /4337-in-a-box/
COPY go.sum /4337-in-a-box/
RUN cd /4337-in-a-box && go mod download

ADD . /4337-in-a-box
RUN cd /4337-in-a-box && go build -o ./bin/4337-in-a-box ./cmd/4337-in-a-box

# Pull the binary into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /4337-in-a-box/bin/4337-in-a-box /usr/local/bin/

EXPOSE 4337
ENTRYPOINT ["4337-in-a-box"]
