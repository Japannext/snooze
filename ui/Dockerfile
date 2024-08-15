FROM node:18-alpine as build

WORKDIR /code

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

COPY package.json ./
COPY pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY . .
RUN pnpm run build

# Run in unprivileged nginx
FROM nginxinc/nginx-unprivileged:1.25-alpine
RUN mkdir -p /usr/share/nginx/html/ui
COPY --from=build /code/dist /usr/share/nginx/html/ui/
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
