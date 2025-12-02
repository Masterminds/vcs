# Use SUSE BCI Golang as the base image
FROM registry.suse.com/bci/golang:1.25

# Install VCS tools and other dependencies
RUN zypper --non-interactive refresh && \
    zypper --non-interactive install -y \
    git \
    mercurial \
    subversion \
    bzr \
    make \
    curl \
    sudo \
    shadow && \
    zypper clean --all

# Create a non-root user for development
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

# Set up Go workspace
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Create Go workspace directories with correct permissions
RUN mkdir -p /go/src /go/bin /go/pkg && \
    chown -R $USERNAME:$USERNAME /go

# Switch to the non-root user
USER $USERNAME

# Set working directory
WORKDIR /workspace
