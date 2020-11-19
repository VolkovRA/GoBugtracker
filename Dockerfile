FROM alpine:latest

# Информация:
LABEL maintainer="VolkovRA"
LABEL description="The bugtracker service for runtime errors"

# Создание директории приложения
WORKDIR /usr/src/app

# Порты:
EXPOSE 80/tcp
EXPOSE 443/tcp

# Копируем файлы приложения:
COPY ./bin ./bin
COPY ./ssl ./ssl

# SSL Сертификаты для замены:
VOLUME $PWD/ssl

# Точка входа:
CMD ["bin/bugtracker.bin"]