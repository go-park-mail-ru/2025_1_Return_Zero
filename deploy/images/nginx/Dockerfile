FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf

RUN mkdir -p /usr/share/nginx/html /usr/share/nginx/static

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"] 