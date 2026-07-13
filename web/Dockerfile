# syntax=docker/dockerfile:1

FROM oven/bun:1-alpine AS builder

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .

ARG VITE_API_URL=http://localhost:3000
ARG VITE_NODE_ENV=production
ARG VITE_APPLICATION_NAME=Fuse
ARG VITE_API_TIMEOUT=30000

ENV VITE_API_URL=$VITE_API_URL
ENV VITE_NODE_ENV=$VITE_NODE_ENV
ENV VITE_APPLICATION_NAME=$VITE_APPLICATION_NAME
ENV VITE_API_TIMEOUT=$VITE_API_TIMEOUT

RUN bun run build

FROM nginx:1.27-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
    CMD wget --quiet --tries=1 --spider http://127.0.0.1/ || exit 1

CMD ["nginx", "-g", "daemon off;"]