FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <dimas.restu@student.upi.edu>

ARG SERVICE_NAME="go-whatsapp-rest"
ENV CONFIG_ENV="production"

WORKDIR /usr/src/app

COPY share/ ./share
COPY dist/${SERVICE_NAME}_linux_amd64/go-whatsapp ./go-whatsapp

RUN chmod 777 share/store share/upload

EXPOSE 3000
HEALTHCHECK --interval=5s --timeout=3s CMD ["curl", "http://127.0.0.1:3000/api/v1/whatsapp/health"] || exit 1

VOLUME ["/usr/src/app/share/store","/usr/src/app/share/upload"]
CMD ["./go-whatsapp"]
