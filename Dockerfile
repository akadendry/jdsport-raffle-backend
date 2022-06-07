FROM golang:1.18-alpine AS builder

# Time zone
RUN apk add --no-cache tzdata

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY *.go ./

# Build the application
RUN go build -o main main.go

# Move to /dist directory as the place for resulting binary folder
# WORKDIR /dist

# Copy binary from build to main folder
# RUN cp /build/main .

# Build a small image
# FROM scratch
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# COPY --from=builder /build/database/migrationSchema/ /database/migrationSchema/
# COPY --from=builder /dist/main /
# COPY ./.env /.env
# COPY ./gcp-cardential.json /gcp-cardential.json
ENV TZ=Asia/Jakarta
# ENV MODE=dev
# ENV MIGRATE_DATABASE=true

# Run executable
CMD ["./jdsport-raffle-backend"]