FROM golang:1.23-alpine AS build

ARG PROJECT=github.com/japannext/snooze
ARG VERSION=0.0.0
ARG COMMIT
WORKDIR /code

# Local CA use
COPY .ca-bundle /usr/local/share/ca-certificates/
RUN chmod -R 644 /usr/local/share/ca-certificates/ && update-ca-certificates

ENV GOCACHE=/root/.cache/go-build

# Copy and download dependencies.
# We do it before for cache reasons
COPY go.sum go.mod ./
RUN go mod download

# Common
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY main.go .
COPY pkg pkg
RUN \
    go build -o /build/snooze \
  -ldflags "-X ${PROJECT}/server.Version=${VERSION} -X ${PROJECT}/server.Commit=${COMMIT} -w -s" \
  .

FROM scratch AS snooze
USER 1000
COPY --from=build --chown=1000 --chmod=755 /build/snooze /snooze
ENTRYPOINT ["/snooze"]
