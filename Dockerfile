# Builder Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:go-1.15 AS go-builder

WORKDIR /usr/src/app

COPY . ./

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -o main cmd/main/main.go


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <dimas.restu@student.upi.edu>

ARG SERVICE_NAME="go-whatsapp-rest"
ENV PATH="$PATH:/usr/app/${SERVICE_NAME}" \
    CONFIG_ENV="production" \
    PRODUCTION_ROUTER_BASE_PATH="/api/v1/whatsapp"

WORKDIR /usr/app/${SERVICE_NAME}

COPY --from=go-builder /usr/src/app/config/ ./config
COPY --from=go-builder /usr/src/app/main ./go-whatsapp-rest

RUN chmod 777 ./config/stores ./config/uploads

EXPOSE 3000
HEALTHCHECK --interval=5s --timeout=3s CMD ["sh", "-c", "curl http://127.0.0.1:3000${PRODUCTION_ROUTER_BASE_PATH}/health || exit 1"]

VOLUME ["/usr/app/${SERVICE_NAME}/config/stores","/usr/app/${SERVICE_NAME}/config/uploads"]
CMD ["go-whatsapp-rest"]
