FROM node:22.20.0-alpine3.22 AS node_modules_deps_builder
WORKDIR /web
COPY web/package.json .
RUN #npm config set registry http://registry.npmmirror.com
RUN npm install

FROM node:22.20.0-alpine3.22 AS web_builder
RUN apk add --no-cache git
WORKDIR /
COPY . .
COPY --from=node_modules_deps_builder /web/node_modules /web/node_modules
RUN cd web && npm install && npm rebuild && npm run build

FROM golang:1.25.1-alpine3.22 AS builder
WORKDIR /app
ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct
ARG VERSION=latest
ARG COMMIT=unknown
ARG BUILDTIME
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN BUILDTIME="${BUILDTIME:-$(date '+%Y-%m-%d %H:%M:%S %z')}" && \
	go build -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X 'main.BuildTime=${BUILDTIME}'" -o ./bin/go_web_template .

FROM alpine:3.22
WORKDIR /app
COPY --from=web_builder /web/dist/ /app/assets/web/v1/
COPY --from=builder /app/bin/go_web_template /app/go_web_template
RUN chmod +x /app/go_web_template && mkdir -p /app/data
EXPOSE 1423
CMD ["./go_web_template"]
