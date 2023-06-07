ARG GO_VER=1.20
ARG GO_ARCH=alpine
ARG INSTALL_KUBECTL=true
ARG OPERATOR_VER=1.29.0
ARG INSTALL_OPERATOR_SDK=true
ARG KUBECTL_VER=1.24.7
ARG INSTALL_KUBEBUILDER=true
ARG KUBEBUILDER_VER=3.10.0

FROM golang:${GO_VER}-${GO_ARCH}

# https://blog.scottchayaa.com/post/2018/11/04/docker-arg-env-variable/#arg
# ARG need to be declared again after FROM

ARG OPERATOR_VER
ARG GO_VER
ARG GO_ARCH
ARG INSTALL_KUBECTL
ARG INSTALL_OPERATOR_SDK
ARG KUBECTL_VER
ARG INSTALL_KUBEBUILDER
ARG KUBEBUILDER_VER

# Install Golang extension package tools
RUN go install github.com/ramya-rao-a/go-outline@latest         && \
    go install github.com/cweill/gotests/gotests@latest         && \
    go install github.com/fatih/gomodifytags@latest             && \
    go install github.com/josharian/impl@latest                 && \
    go install github.com/haya14busa/goplay/cmd/goplay@latest   && \
    go install github.com/go-delve/delve/cmd/dlv@latest         && \
    go install honnef.co/go/tools/cmd/staticcheck@latest        && \
    go install golang.org/x/tools/gopls@latest


RUN if [ "${GO_ARCH}" == "alpine" ]; then \
        apk add git curl bash alpine-sdk; \
    fi



# Install kubectl
RUN if [ "${INSTALL_KUBECTL}" == "true" ]; then \
        export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac) && \
        export OS=$(uname | awk '{print tolower($0)}') && \
        curl -LO "https://dl.k8s.io/release/v${KUBECTL_VER}/bin/${OS}/${ARCH}/kubectl" && \
        chmod +x kubectl && mv kubectl /usr/local/bin/kubectl; \
    fi

# Install kubebuilder
RUN if [ "${INSTALL_KUBEBUILDER}" == "true" ]; then \
        export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac) && \
        export OS=$(uname | awk '{print tolower($0)}') && \
        export KUBEBUILDER_DL_URL=https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VER} && \
        curl -LO ${KUBEBUILDER_DL_URL}/kubebuilder_${OS}_${ARCH} && \
        chmod +x kubebuilder_${OS}_${ARCH} && \
        mv kubebuilder_${OS}_${ARCH} /usr/local/bin/kubebuilder; \
    fi

# Install operator-sdk 
RUN if [ "${INSTALL_OPERATOR_SDK}" == "true" ]; then \
        export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac) && \
        export OS=$(uname | awk '{print tolower($0)}') && \
        export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/download/v${OPERATOR_VER} && \        
        curl -LO ${OPERATOR_SDK_DL_URL}/operator-sdk_${OS}_${ARCH} && \
        chmod +x operator-sdk_${OS}_${ARCH} && \
        mv operator-sdk_${OS}_${ARCH} /usr/local/bin/operator-sdk; \
    fi

