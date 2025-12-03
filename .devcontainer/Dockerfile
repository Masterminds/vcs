# Use the official Go + Node.js devcontainer base image
FROM mcr.microsoft.com/devcontainers/go:1.25

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install VCS tools and other dependencies
RUN apt-get update && \
    apt-get install -y \
    git \
    mercurial \
    subversion \
    bzr \
    make \
    curl \
    sudo && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install GitHub Copilot CLI (architecture-aware)
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then \
        COPILOT_ARCH="amd64"; \
    elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then \
        COPILOT_ARCH="arm64"; \
    else \
        echo "Unsupported architecture: $ARCH" && exit 1; \
    fi && \
    curl -fsSL -o /tmp/copilot.tar.gz "https://github.com/github/copilot-cli/releases/download/v0.0.366-9/copilot-linux-${COPILOT_ARCH}.tar.gz" && \
    tar -xzf /tmp/copilot.tar.gz -C /usr/local/bin && \
    chmod +x /usr/local/bin/copilot && \
    rm -f /tmp/copilot.tar.gz

# Set working directory
WORKDIR /workspace
