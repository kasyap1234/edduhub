ARG GO_VERSION=1.24.2
FROM golang:${GO_VERSION}

# Create a non-root user
RUN adduser --disabled-password --gecos '' vscode && \
    apt-get update && apt-get install -y git curl unzip

# Set Go env
ENV PATH="/go/bin:$PATH"
ENV GO111MODULE=on

# Set working dir
WORKDIR /workspace

# Use non-root user
USER vscode
