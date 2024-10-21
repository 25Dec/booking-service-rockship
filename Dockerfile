# syntax=docker/dockerfile:1

FROM golang:latest as builder


# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o docker-gs-ping

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose

# Run

# Step 3: Final
FROM scratch
ENV ENV PROD
COPY --from=builder /app/docker-gs-ping docker-gs-ping
COPY --from=builder /app/config/config.yml config/config.yml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080

CMD ["/docker-gs-ping"]