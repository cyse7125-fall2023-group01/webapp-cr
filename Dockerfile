############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
#ENV DATABASE_URL "user=postgres password=postgres dbname=app host=127.0.0.2 port=5432 sslmode=disable"
# Build the binary.
RUN go build -o /go/bin/http-check
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
#ENV DATABASE_URL "user=postgres password=postgres dbname=app host=127.0.0.1 port=5432 sslmode=disable"
COPY --from=builder /go/bin/http-check /go/bin/http-check
# Run the hello binary.
ENTRYPOINT ["/go/bin/http-check"]