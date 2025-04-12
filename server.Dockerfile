FROM alpine:latest

WORKDIR /app

COPY migrations/ /app/migrations
COPY web/assets /app/assets
COPY templates/ /app/templates
COPY bin/myapp-linux-amd-64 .

# Ensure binary is executable
RUN chmod +x /app/myapp-linux-amd-64

EXPOSE 3000

CMD ["./myapp-linux-amd-64"]