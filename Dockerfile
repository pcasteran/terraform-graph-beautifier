ARG GO_VERSION="1.19-alpine3.17"

ARG USER="appuser"

FROM golang:${GO_VERSION} AS builder

# Install the system dependencies.
RUN apk add --no-cache upx

# Create a new user.
# This is done here as the `addgroup` and `adduser` commands are not available in the scratch image.
ARG USER
RUN addgroup --gid 1001 --system "${USER}" && \
    adduser  --uid 1001 --ingroup "${USER}" --shell /bin/false --disabled-password --no-create-home --system "${USER}"

# Install the project dependencies.
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# Build (and compress) the binary for the target architecture and operating system.
COPY . .
RUN CGO_ENABLED=0 GOARCH=${TARGETARCH} GOOS=${TARGETOS} go build -ldflags="-w -s" -o terraform-graph-beautifier && \
    upx --best --lzma terraform-graph-beautifier

##

FROM scratch

ARG USER

# Copy the static executable from the builder.
COPY --from=builder /build/terraform-graph-beautifier /usr/local/bin/

# Import the user and group files from the builder.
# Set the user that will be used by default if not overriden with `docker run --user ...`.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER "${USER}":"${USER}"

# Set the container entrypoint.
ENTRYPOINT ["/usr/local/bin/terraform-graph-beautifier"]
