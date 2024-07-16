FROM golang:1.22-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags="-s -w" -o /main

FROM alpine

RUN apk update --no-cache

WORKDIR /app

COPY --from=builder /main /main

CMD ["/main"]

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux go build -o /main

# EXPOSE 80

# CMD ["/main"]