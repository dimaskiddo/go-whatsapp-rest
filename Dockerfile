FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <dimas.restu@student.upi.edu>

ARG SERVICE_NAME="go-whatsapp-rest"

ENV CONFIG_ENV="PROD" \
    CONFIG_FILE_PATH="./configs" \
    CONFIG_LOG_LEVEL="INFO" \
    CONFIG_LOG_SERVICE="$SERVICE_NAME"

WORKDIR /usr/src/app
COPY build/ .
RUN chmod 777 stores uploads

EXPOSE 3000
HEALTHCHECK --interval=5s --timeout=3s CMD ["curl", "http://127.0.0.1:3000/health"] || exit 1

VOLUME ["/usr/src/app/stores","/usr/src/app/uploads"]
CMD ["./main"]
