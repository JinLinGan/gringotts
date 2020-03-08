FROM golang:1.13.5-buster as builder

RUN mkdir /gringotts
WORKDIR /gringotts

COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

RUN go build -o ./bin/gringotts-agent github.com/jinlingan/gringotts/gringotts-agent

FROM golang:1.13.5-buster

COPY --from=builder /gringotts/bin/gringotts-agent /usr/local/bin/gringotts-agent

CMD ["gringotts-agent","start","-s","gringotts-server:6667"]