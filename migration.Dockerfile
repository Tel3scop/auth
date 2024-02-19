FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /root

# Create the migrations directory to ensure it exists
RUN mkdir -p migrations

# Copy the entire migrations directory (even if it's empty)
COPY migrations migrations/

ADD migration.sh .
ADD migration.env .

RUN chmod +x migration.sh

ENTRYPOINT ["bash", "migration.sh"]
