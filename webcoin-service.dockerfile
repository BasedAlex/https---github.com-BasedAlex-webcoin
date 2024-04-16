FROM alpine:latest

RUN mkdir /app

COPY webcoinApp /app

CMD [ "/app/webcoinApp" ] 