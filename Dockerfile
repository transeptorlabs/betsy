# Build 4337-in-a-box in a container
FROM golang:1.22-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers git

# Get dependencies
COPY go.mod /betsy/
COPY go.sum /betsy/
RUN cd /betsy && go mod download

ADD . /betsy
RUN cd /betsy && go build -o ./bin/betsy ./cmd/betsy

# Pull the binary into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /betsy/bin/betsy /usr/local/bin/

EXPOSE 4337
ENTRYPOINT ["betsy"]
