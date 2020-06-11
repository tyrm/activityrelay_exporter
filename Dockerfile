FROM golang:1.14 AS builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o activitypub_exporter .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/activitypub_exporter /activitypub_exporter
CMD ["/activitypub_exporter"]