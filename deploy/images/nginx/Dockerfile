FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /usr/share/nginx/html /usr/share/nginx/static

RUN mkdir -p /var/cache/nginx

EXPOSE 443

CMD ["nginx", "-g", "daemon off;"] 