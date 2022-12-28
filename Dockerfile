ARG ALPINE_VERSION="3.17"
ARG GO_VERSION="1.19-alpine3.17"
ARG GORELEASER_VERSION="v1.14.0"

ARG USER="appuser"

##

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS binary_builder

ARG TARGETOS
ARG TARGETARCH

ARG GORELEASER_VERSION

# Install the system dependencies.
RUN apk add --no-cache upx git && \
    go install github.com/goreleaser/goreleaser@${GORELEASER_VERSION}

# Install the project dependencies.
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# Build the binary for the target architecture and operating system.
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} goreleaser build --snapshot --rm-dist  --single-target

##

FROM --platform=$BUILDPLATFORM alpine:${ALPINE_VERSION} AS user_builder

ARG USER

# Create a new user.
# This is done in a dedicated step as the `addgroup` and `adduser` commands are not available in the scratch image.
RUN addgroup --gid 1001 --system "${USER}" && \
    adduser  --uid 1001 --ingroup "${USER}" --shell /bin/false --disabled-password --no-create-home --system "${USER}"

##

FROM scratch AS base

ARG USER

# Import the user and group files from the builder stage.
# Set the user that will be used by default if not overriden with `docker run --user ...`.
COPY --from=user_builder /etc/passwd /etc/passwd
COPY --from=user_builder /etc/group /etc/group
USER "${USER}":"${USER}"

##

FROM base AS binary_from_build_context

ARG TARGETOS
ARG TARGETARCH

# Copy the static executable from the build context.
COPY dist/terraform-graph-beautifier_${TARGETOS}_${TARGETARCH}*/terraform-graph-beautifier* \
     /usr/local/bin/

# Set the container entrypoint.
ENTRYPOINT ["/usr/local/bin/terraform-graph-beautifier"]

##

FROM base AS binary_from_builder

ARG TARGETOS
ARG TARGETARCH

# Copy the static executable from the builder stage.
COPY --from=binary_builder \
     /build/dist/terraform-graph-beautifier_${TARGETOS}_${TARGETARCH}*/terraform-graph-beautifier* \
     /usr/local/bin/

# Set the container entrypoint.
ENTRYPOINT ["/usr/local/bin/terraform-graph-beautifier"]
