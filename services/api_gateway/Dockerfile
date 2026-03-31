# Берем официальный образ Nginx
FROM nginx:latest

# Удаляем стандартный конфиг Nginx, чтобы он не мешал
RUN rm /etc/nginx/conf.d/default.conf

# Копируем твой кастомный конфиг внутрь контейнера
COPY nginx.conf /etc/nginx/nginx.conf