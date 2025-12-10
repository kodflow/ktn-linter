FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

# Detect architecture automatically for multi-platform builds
ARG TARGETARCH

# Install all system tools in a single layer (rarely changes - good for cache)
# hadolint ignore=DL3008
RUN apt-get update \
    # Install base dependencies for adding repositories
    && apt-get install -y --no-install-recommends \
        curl \
        wget \
        ca-certificates \
        gnupg \
        lsb-release \
        cmake \
        libpam0g-dev \
    # Add GitHub CLI repository
    && mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | gpg --dearmor -o /etc/apt/keyrings/githubcli-archive-keyring.gpg \
    && chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list \
    # Add Node.js repository
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list \
    # Add HashiCorp repository
    && curl -fsSL https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /etc/apt/keyrings/hashicorp-archive-keyring.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/hashicorp.list \
    # Update with new repositories and install packages
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        g++ \
        zlib1g-dev \
        git \
        jq \
        zsh \
        unzip \
        nodejs \
        gh \
        terraform \
        vault \
        consul \
        nomad \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    # Install yq binary (not available in Ubuntu repos)
    && YQ_ARCH=$(dpkg --print-architecture) \
    && curl -fsSL "https://github.com/mikefarah/yq/releases/latest/download/yq_linux_${YQ_ARCH}" -o /usr/local/bin/yq \
    && chmod +x /usr/local/bin/yq \
    # Install 1Password CLI separately (may fail on some architectures)
    && curl -fsSL https://downloads.1password.com/linux/keys/1password.asc | gpg --dearmor -o /etc/apt/keyrings/1password-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/$(dpkg --print-architecture) stable main" > /etc/apt/sources.list.d/1password.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends 1password-cli || true \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install AWS CLI v2 and golangci-lint v2 in parallel (rarely changes - good for cache)
# AWS CLI uses different architecture naming: arm64 -> aarch64, amd64 -> x86_64
RUN AWS_ARCH=$([ "$TARGETARCH" = "arm64" ] && echo "aarch64" || echo "x86_64") && \
    curl -fsSL "https://awscli.amazonaws.com/awscli-exe-linux-${AWS_ARCH}.zip" -o "/tmp/awscliv2.zip" & \
    curl -fsSL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" -o "/tmp/golangci-install.sh" & \
    wait && \
    unzip /tmp/awscliv2.zip -d /tmp && \
    /tmp/aws/install && \
    sh /tmp/golangci-install.sh -b /usr/local/bin latest && \
    rm -rf /tmp/awscliv2.zip /tmp/aws /tmp/golangci-install.sh

# Switch to vscode user
USER vscode
WORKDIR /home/vscode

# Create directories for persistent volumes and clean old golangci-lint
RUN mkdir -p /home/vscode/.cache \
    /home/vscode/.config \
    /home/vscode/.config/op \
    /home/vscode/.zsh_history_dir \
    /home/vscode/.local/bin \
    /home/vscode/.claude \
    /home/vscode/.local/share/@anthropic \
    # Set strict permissions for 1Password config directory
    && chmod 700 /home/vscode/.config/op \
    # Remove any old golangci-lint installations from user cache
    && rm -rf /home/vscode/.cache/go/bin/golangci-lint

# Install Oh My Zsh and themes in parallel (rarely changes - good for cache)
RUN rm -rf "$HOME/.oh-my-zsh" && \
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended && \
    git clone --depth=1 https://github.com/romkatv/powerlevel10k.git \
        "${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/powerlevel10k" & \
    git clone --depth=1 https://github.com/dracula/zsh.git \
        "${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/dracula" & \
    wait

# Copy p10k configuration (changes rarely)
COPY --chown=vscode:vscode p10k.sh /home/vscode/.p10k.zsh

# Configure zshrc (changes rarely)
RUN sed -i 's/^ZSH_THEME=.*/ZSH_THEME="powerlevel10k\/powerlevel10k"/' "$HOME/.zshrc" && \
    { \
        echo ''; \
        echo '# To customize prompt, run "p10k configure" or edit ~/.p10k.zsh'; \
        echo '[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh'; \
        echo ''; \
        echo '# Persistent zsh history'; \
        echo 'export HISTFILE=~/.zsh_history_dir/.zsh_history'; \
        echo 'export HISTSIZE=10000'; \
        echo 'export SAVEHIST=10000'; \
        echo ''; \
        echo '# Ensure /usr/local/bin takes priority over user cache'; \
        echo 'export PATH="/usr/local/bin:$PATH"'; \
        echo ''; \
        echo '# Super Claude alias'; \
        echo 'alias super-claude="claude --dangerously-skip-permissions --mcp-config /home/vscode/.devcontainer/mcp.json"'; \
    } >> "$HOME/.zshrc"

# Install Go 1.25.2 as root (placed near end for cache optimization)
# Go uses standard architecture naming: arm64, amd64
USER root
RUN GO_ARCH=$([ "$TARGETARCH" = "amd64" ] && echo "amd64" || echo "arm64") && \
    curl -fsSL https://go.dev/dl/go1.25.2.linux-${GO_ARCH}.tar.gz -o /tmp/go.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz && \
    ln -sf /usr/local/go/bin/go /usr/local/bin/go && \
    ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt

# Configure npm to use user-local directory for global packages
USER vscode
ENV PATH="/home/vscode/.npm-global/bin:$PATH"
RUN mkdir -p /home/vscode/.npm-global && \
    npm config set prefix '/home/vscode/.npm-global' && \
    echo 'export PATH=/home/vscode/.npm-global/bin:$PATH' >> /home/vscode/.zshrc && \
    npm install -g @anthropic-ai/claude-code@latest @bazel/bazelisk

# Create .devcontainer directory with empty mcp.json (fallback if setup-mcp.sh fails)
RUN mkdir -p /home/vscode/.devcontainer && \
    echo '{"mcpServers":{}}' > /home/vscode/.devcontainer/mcp.json

# Set working directory
WORKDIR /workspace

# Copy and set permissions for setup-mcp.sh
COPY --chown=vscode:vscode --chmod=755 setup-mcp.sh /workspace/.devcontainer/setup-mcp.sh
