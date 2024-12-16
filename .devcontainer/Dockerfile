FROM mcr.microsoft.com/devcontainers/go:1.23

# Install utilities
RUN go install golang.org/x/tools/gopls@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest
RUN go install golang.org/x/lint/golint@latest
RUN go install github.com/mgechev/revive@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/fatih/gomodifytags@latest
RUN go install github.com/haya14busa/goplay/cmd/goplay@latest
RUN go install github.com/cweill/gotests/gotests@latest
RUN go install github.com/josharian/impl@latest
