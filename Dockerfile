FROM golang:1.20-alpine AS builder

RUN apk add --no-cache make

WORKDIR /src

COPY go.mod go.sum ./
ENV GOPROXY='https://goproxy.cn'
RUN go mod download

COPY . .

ARG VERSION=unknown
RUN make build

FROM alpine

RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=builder /src/bailian2openai /app/


ENTRYPOINT ["./bailian2openai"]
