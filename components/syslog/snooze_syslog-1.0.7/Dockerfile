FROM ghcr.io/astral-sh/uv:python3.8-alpine AS build

COPY .ca-bundle /usr/local/share/ca-certificates/
RUN find /usr/local/share/ca-certificates/ -type f -exec cat {} \; >> /etc/ssl/certs/ca-certificates.crt

ENV CA_REQUESTS_BUNDLE /etc/ssl/certs/ca-certificates.crt
ENV UV_NATIVE_TLS=true
ENV UV_LINK_MODE=copy
ENV UV_CACHE_DIR=/opt/uv-cache

#RUN apk add --no-cache \
#      gcc==13.2.1_git20240309-r1 \
#      libldap==2.6.8-r0 \
#      musl-dev==1.2.5-r1 \
#      libffi-dev==3.4.6-r0

WORKDIR /code/syslog
COPY pyproject.toml uv.lock README.md /code/syslog/

RUN --mount=type=cache,target=/opt/uv-cache \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    --mount=type=bind,source=README.md,target=README.md \
    uv sync --locked --no-install-project --compile-bytecode --no-editable

ADD src/ .
RUN --mount=type=cache,target=/opt/uv-cache \
    uv sync --locked --compile-bytecode --no-editable


FROM python:3.8-alpine

USER 1001
WORKDIR /app
COPY --from=build --chown=1001:1001 /code/syslog/.venv /app/.venv
ENV PATH="/app/.venv/bin:$PATH"

CMD ["snooze-syslog"]
