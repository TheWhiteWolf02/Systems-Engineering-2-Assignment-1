FROM alpine:latest

RUN apk update

RUN apk add nginx

COPY ./nginx.conf /etc/nginx

EXPOSE 80

RUN mkdir -p /run/nginx

CMD [ "nginx" ]
