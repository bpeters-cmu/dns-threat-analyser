# Multistage Docker build
FROM golang:1.17.2-alpine as builder
COPY . /dns-threat-analyser
WORKDIR /dns-threat-analyser
RUN apk add --no-cache gcc musl-dev
ENV GO111MODULE=on
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/DnsThreatAnalyser cmd/server/server.go &&\
    CGO_ENABLED=1 GOOS=linux go build -o bin/SmokeTestClient cmd/client/smoketest.go

FROM alpine:3.13.6
RUN apk update \
    && apk upgrade \
    && apk add --no-cache sqlite \
    && apk add bind-tools
WORKDIR /root/
COPY --from=builder /dns-threat-analyser/bin /dns-threat-analyser/credentials.json ./

CMD ["./DnsThreatAnalyser"]