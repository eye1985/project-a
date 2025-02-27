FROM alpine:latest

WORKDIR /app

COPY bin/myapp-linux-amd-64 .
COPY templates/ /app/templates
COPY assets/ /app/assets
COPY migrations/ /app/migrations

# Ensure binary is executable
RUN chmod +x /app/myapp-linux-amd-64

EXPOSE 8080

CMD ["./myapp-linux-amd-64"]