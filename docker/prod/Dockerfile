FROM node:lts-bookworm AS frontend

COPY sdk/ /sdk
COPY frontend/ /frontend
COPY admin/ /admin

RUN cd /sdk && npm i && npx tsc

RUN cd /frontend && npm i && npm run build
RUN cd /admin && npm i && npm run build

FROM golang:1.23-bookworm AS backend

ARG PARTYHALL_VERSION="indev"
ARG PARTYHALL_COMMIT="XXXXXXX"

RUN apt-get update && apt-get install -y \
    gcc-aarch64-linux-gnu \
    g++-aarch64-linux-gnu \
    libc6-dev-arm64-cross \
    && rm -rf /var/lib/apt/lists/*

COPY ./backend /app
COPY --from=frontend /frontend/dist /app/frontend/appliance
COPY --from=frontend /admin/dist /app/frontend/app

RUN cd /app \
    && CGO_ENABLED=1 \
       GOOS=linux \
       GOARCH=amd64 \
       go build \
        -ldflags="-X github.com/partyhall/partyhall/utils.CURRENT_VERSION=${PARTYHALL_VERSION} -X github.com/partyhall/partyhall/utils.CURRENT_COMMIT=${PARTYHALL_COMMIT}" \
        -o /partyhall-${PARTYHALL_VERSION}-linux-x86_64 \
    && chmod +x /partyhall-${PARTYHALL_VERSION}-linux-x86_64

RUN cd /app \
    && CGO_ENABLED=1 \
       CC=aarch64-linux-gnu-gcc \
       GOOS=linux \
       GOARCH=arm64 \
       go build \
        -ldflags="-X github.com/partyhall/partyhall/utils.CURRENT_VERSION=${PARTYHALL_VERSION} -X github.com/partyhall/partyhall/utils.CURRENT_COMMIT=${PARTYHALL_COMMIT}" \
        -o /partyhall-${PARTYHALL_VERSION}-linux-aarch64 \
    && chmod +x /partyhall-${PARTYHALL_VERSION}-linux-aarch64
