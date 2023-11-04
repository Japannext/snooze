FROM golang:1.20 as build

ARG PROJECT=github.com/japannext/snooze
ARG VERSION=0.0.0
ARG COMMIT
WORKDIR /app

# Local CA use
COPY .ca-bundle /usr/local/share/ca-certificates/
RUN chmod -R 644 /usr/local/share/ca-certificates/ && update-ca-certificates

# Copy and download dependencies.
# We do it before for cache reasons
COPY go.mod go.sum ./
RUN go mod download

# Add all file/directories that contain code.
ADD server/ ./server
ADD cmd/ ./cmd
ADD api/ ./api
ADD main.go .

# Build snooze-process
RUN CGO_ENABLED=0 GOOS=linux go build \
  cmd/process/main.go \
  -o ./build/snooze-process \
  -ldflags "-X ${PROJECT}/server.Version=${VERSION} -X ${PROJECT}/server.Commit=${COMMIT} -w -s"

FROM scratch AS snooze-process
USER 1000
COPY --from=build --chown=1000 --chmod=755 /app/build/snooze-process /
CMD ["/snooze-process"]
