# Builder Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:go-1.12 AS go-builder

WORKDIR /usr/src/app

COPY . ./

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -a -o dist/go-whatsapp *.go


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <dimas.restu@student.upi.edu>

ARG SERVICE_NAME="go-whatsapp-rest"
ENV PATH $PATH:/opt/${SERVICE_NAME}

ENV CONFIG_ENV="production" \
    PRODUCTION_ROUTER_BASE_PATH="/api/v1/whatsapp"

WORKDIR /opt/${SERVICE_NAME}

COPY --from=go-builder /usr/src/app/dist/go-whatsapp ./go-whatsapp
COPY share/ ./share

RUN chmod 777 share/store share/upload

EXPOSE 3000
HEALTHCHECK --interval=5s --timeout=3s CMD curl --fail http://127.0.0.1:3000${PRODUCTION_ROUTER_BASE_PATH}/health || exit 1

VOLUME ["/usr/src/app/share/store","/usr/src/app/share/upload"]
CMD ["go-whatsapp"]
