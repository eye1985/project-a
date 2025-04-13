FROM alpine:latest

WORKDIR /app

COPY migrations/ /app/migrations
COPY web/assets /app/web/assets
COPY web/templates/ /app/web/templates
COPY bin/myapp-linux-amd-64 .

# Ensure binary is executable
RUN chmod +x /app/myapp-linux-amd-64

EXPOSE 3000

CMD ["./myapp-linux-amd-64"]