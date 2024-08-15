FROM golang:1.22 AS build

ARG PROJECT=github.com/japannext/snooze
ARG VERSION=0.0.0
ARG COMMIT
WORKDIR /code

# Local CA use
COPY .ca-bundle /usr/local/share/ca-certificates/
RUN chmod -R 644 /usr/local/share/ca-certificates/ && update-ca-certificates

# Copy and download dependencies.
# We do it before for cache reasons
RUN \
  --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download

# Common
ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY main.go go.sum go.mod ./
COPY pkg/ ./
RUN \
  --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=bind,target=. \
  go build -o /build/snooze \
  -ldflags "-X ${PROJECT}/server.Version=${VERSION} -X ${PROJECT}/server.Commit=${COMMIT} -w -s" \
  .

FROM scratch AS snooze
USER 1000
COPY --from=build --chown=1000 --chmod=755 /build/snooze /snooze
ENTRYPOINT ["/snooze"]
