# Builder
#########
FROM golang:latest as builder
WORKDIR /go/src/github.com/scriptsmith/depot

# Copy src
COPY *.go ./

# Build project
RUN go get ./...
RUN go build

# App
#####
FROM golang:latest
WORKDIR /app

# Copy app
COPY --from=builder /go/src/github.com/scriptsmith/depot/depot .
COPY templates ./templates

# Setup environment
RUN mkdir /data
ENV DEPOT_ROOT /data
ENV DEPOT_PORT 8080

# Run depot
CMD ["./depot"]