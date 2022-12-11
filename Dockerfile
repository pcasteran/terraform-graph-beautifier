ARG GO_VERSION="1.19-alpine3.17"
ARG ALPINE_VERSION="3.17"

FROM golang:${GO_VERSION} AS builder

##

FROM alpine:${ALPINE_VERSION}

# Create the workspace directory.
ENV WORKSPACE_DIR="/workspace"
RUN mkdir -p ${WORKSPACE_DIR} && chmod 777 ${WORKSPACE_DIR}
WORKDIR ${WORKSPACE_DIR}
VOLUME ["${WORKSPACE_DIR}"]

# Create a new user that will be used by default if not overriden with `docker run --user ...`.
RUN addgroup --gid 1001 --system tf_beautifier && \
    adduser  --uid 1001 --ingroup tf_beautifier --shell /bin/false --disabled-password --no-create-home --system tf_beautifier
USER tf_beautifier

# Set the container entrypoint.
# TODO
