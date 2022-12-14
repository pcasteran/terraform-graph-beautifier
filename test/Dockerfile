ARG NODE_VERSION="19-slim"

##

FROM node:${NODE_VERSION}

# Install the system dependencies.
RUN apt-get -y update && \
    apt-get -y install --no-install-recommends git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install Install the test framework.
WORKDIR /workspace
COPY package.json package-lock.json ./

# hadolint ignore=DL3016
RUN npm install --global npm && \
    npm install
