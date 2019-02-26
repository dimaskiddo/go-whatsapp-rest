FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <dimas.restu@student.upi.edu>

ENV CONFIG_ENV=PROD \
    CONFIG_FILE_PATH=./configs

WORKDIR /usr/src/app
COPY build/ .
RUN chmod 777 stores uploads

EXPOSE 3000
HEALTHCHECK --interval=3s --timeout=3s CMD ["curl", "http://127.0.0.1:3000/health"] || exit 1

VOLUME ["/usr/src/app/stores","/usr/src/app/uploads"]
CMD ["./whatsapp-go"]
