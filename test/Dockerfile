ARG NODE_VERSION="19-slim"

ARG TERRAFORM_VERSION="1.3.6"

##

FROM node:${NODE_VERSION}

ARG TARGETOS
ARG TARGETARCH

ARG TERRAFORM_VERSION

# Install the system dependencies.
RUN apt-get -y update && \
    apt-get -y install --no-install-recommends git wget unzip ca-certificates jq && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install Terraform.
WORKDIR /tmp
RUN wget -qO terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${TARGETOS}_${TARGETARCH}.zip && \
    unzip terraform.zip && \
    rm terraform.zip && \
    mv terraform /usr/local/bin/

# Install Install the test framework.
WORKDIR /workspace/test
RUN chmod 777 /workspace/test
COPY package.json package-lock.json ./

# hadolint ignore=DL3016
RUN npm install --global npm && \
    npm install
