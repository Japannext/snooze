FROM golang:1.20 AS deps

ARG PROJECT=github.com/japannext/snooze
ARG VERSION=0.0.0
ARG COMMIT
WORKDIR /code

# Local CA use
COPY .ca-bundle /usr/local/share/ca-certificates/
RUN chmod -R 644 /usr/local/share/ca-certificates/ && update-ca-certificates

# Copy and download dependencies.
# We do it before for cache reasons
COPY go.mod go.sum ./
RUN go mod download

# snooze-apiserver
FROM deps AS build-apiserver
WORKDIR /code
# Common
COPY common common
# App specific
COPY apiserver apiserver
# Version last since it changes often
COPY version version
RUN CGO_ENABLED=0 GOOS=linux go build \
  -o /build/snooze-apiserver \
  -ldflags "-X ${PROJECT}/server.Version=${VERSION} -X ${PROJECT}/server.Commit=${COMMIT} -w -s" \
  sources/apiserver/cmd
FROM scratch AS snooze-apiserver
USER 1000
COPY --from=build-apiserver --chown=1000 --chmod=755 /build/snooze-apiserver /
CMD ["/snooze-apiserver"]

# snooze-otel
FROM deps AS build-otel
WORKDIR /code
# Common
COPY common common
# App specific
COPY sources/otel sources/otel
# Version last since it changes often
COPY version version
RUN CGO_ENABLED=0 GOOS=linux go build \
  -o /build/snooze-otel \
  -ldflags "-X ${PROJECT}/version.Version=${VERSION} -X ${PROJECT}/version.Commit=${COMMIT} -w -s" \
  sources/otel/cmd
FROM scratch AS snooze-otel
USER 1000
COPY --from=build-otel --chown=1000 --chmod=755 /build/snooze-otel /
CMD ["/snooze-otel"]

# snooze-alertmanager
FROM deps AS build-alertmanager
WORKDIR /code
# Common
COPY common common
# App specific
COPY sources/alertmanager sources/alertmanager
# Version last since it changes often
COPY version version
RUN CGO_ENABLED=0 GOOS=linux go build \
  -o /build/snooze-alertmanager \
  -ldflags "-X ${PROJECT}/server.Version=${VERSION} -X ${PROJECT}/server.Commit=${COMMIT} -w -s" \
  sources/alertmanager/cmd
FROM scratch AS snooze-alertmanager
USER 1000
COPY --from=build-alertmanager --chown=1000 --chmod=755 /build/snooze-alertmanager /
CMD ["/snooze-alertmanager"]
