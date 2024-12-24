FROM node:18-alpine as build

WORKDIR /code

COPY .ca-bundle /usr/local/share/ca-certificates/
RUN find /usr/local/share/ca-certificates/ -type f -exec cat {} \; >> /etc/ssl/certs/ca-certificates.crt
ENV NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

COPY package.json ./
COPY pnpm-lock.yaml ./
# Installing pnpm
RUN corepack enable && corepack install
# Installing dependencies
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY . .
RUN pnpm run build

# Run in unprivileged nginx
FROM nginxinc/nginx-unprivileged:1.25-alpine
RUN mkdir -p /usr/share/nginx/html/ui
COPY --from=build /code/dist /usr/share/nginx/html/ui/
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
