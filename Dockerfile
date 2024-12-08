FROM php:8.3-cli-alpine

WORKDIR /app

COPY . .

EXPOSE 27420

ENTRYPOINT [ "php", "xrdebug", "-p", "27420" ]
