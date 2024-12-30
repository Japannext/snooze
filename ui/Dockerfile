FROM oven/bun:1-alpine AS build
WORKDIR /code

COPY .ca-bundle /usr/local/share/ca-certificates/
RUN find /usr/local/share/ca-certificates/ -type f -exec cat {} \; >> /etc/ssl/certs/ca-certificates.crt
ENV NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

# Installing dependencies
COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile

# Copy code and build
COPY . .
RUN bun run build

# Run in unprivileged nginx
FROM nginxinc/nginx-unprivileged:1.25-alpine
RUN mkdir -p /usr/share/nginx/html/ui
COPY --from=build /code/dist /usr/share/nginx/html/ui/
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
