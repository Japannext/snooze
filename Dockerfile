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

# snooze-processor
FROM deps AS build
WORKDIR /code
# Common
COPY pkg pkg
COPY go.sum go.sum main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build \
  -o /build/snooze \
  -ldflags "-X ${PROJECT}/server.Version=${VERSION} -X ${PROJECT}/server.Commit=${COMMIT} -w -s" \
  .

FROM scratch AS snooze
USER 1000
COPY --from=build --chown=1000 --chmod=755 /build/snooze /bin
ENTRYPOINT ["/bin/snooze"]
