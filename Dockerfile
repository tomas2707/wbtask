# Use a small base image
FROM alpine:latest

# Install ca-certificates in case the application makes HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a directory for the app
WORKDIR /app

# Copy the binary from your host machine to the Docker image
COPY main /app/

# Expose port 8080 for the application
EXPOSE 8080

# Command to run the binary
CMD ["/app/main"]
