FROM oven/bun:1.2.21-alpine AS build

WORKDIR /app

ARG VITE_BACKEND_URL=http://localhost:8080
ARG VITE_BACKEND_WS_URL=ws://localhost:8080/ws

ENV VITE_BACKEND_URL=$VITE_BACKEND_URL
ENV VITE_BACKEND_WS_URL=$VITE_BACKEND_WS_URL

COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile

COPY web/ ./

RUN bun run build

FROM nginx:1.29-alpine

COPY web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
