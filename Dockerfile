ARG GO_VERSION="1.19-alpine3.17"
ARG GORELEASER_VERSION="v1.13.1"

ARG USER="appuser"

##

FROM golang:${GO_VERSION} AS builder

ARG TARGETOS
ARG TARGETARCH

ARG USER

ARG GORELEASER_VERSION

# Install the system dependencies.
RUN apk add --no-cache upx git && \
    go install github.com/goreleaser/goreleaser@${GORELEASER_VERSION}

# Create a new user.
# This is done here as the `addgroup` and `adduser` commands are not available in the scratch image.
RUN addgroup --gid 1001 --system "${USER}" && \
    adduser  --uid 1001 --ingroup "${USER}" --shell /bin/false --disabled-password --no-create-home --system "${USER}"

# Install the project dependencies.
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# Build (and compress) the binary for the target architecture and operating system.
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} goreleaser build --snapshot --rm-dist  --single-target

##

FROM scratch AS base

ARG USER

# Import the user and group files from the builder.
# Set the user that will be used by default if not overriden with `docker run --user ...`.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER "${USER}":"${USER}"

##

FROM base AS binary_from_build_context

ARG TARGETOS
ARG TARGETARCH

# Copy the static executable from the build context.
COPY dist/terraform-graph-beautifier_${TARGETOS}_${TARGETARCH}*/terraform-graph-beautifier \
     /usr/local/bin/

# Set the container entrypoint.
ENTRYPOINT ["/usr/local/bin/terraform-graph-beautifier"]

##

FROM base AS binary_from_builder

ARG TARGETOS
ARG TARGETARCH

# Copy the static executable from the builder stage.
COPY --from=builder \
     /build/dist/terraform-graph-beautifier_${TARGETOS}_${TARGETARCH}*/terraform-graph-beautifier \
     /usr/local/bin/

# Set the container entrypoint.
ENTRYPOINT ["/usr/local/bin/terraform-graph-beautifier"]
