# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /repoear

##
## Deploy
##
FROM alpine:3.14

RUN apk --no-cache add \
    bash \
    curl>7.77.0-r0 \
    git \
    libc6-compat \
    openssh-client \
    py3-pip \
    py3-wheel \
    python3 && \
    pip install --upgrade pip==21.0.1

# Install a YAML Linter
ARG yamllint_version=1.26.1
LABEL yamllint_version=$yamllint_version
RUN pip install "yamllint==$yamllint_version"

# Install Yamale YAML schema validator
ARG yamale_version=3.0.7
LABEL yamale_version=$yamale_version
RUN pip install "yamale==$yamale_version"

# Install kubectl
ARG kubectl_version=v1.21.3
LABEL kubectl_version=$kubectl_version
RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/$kubectl_version/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Install Helm 3
ARG helm_version=v3.6.3
LABEL helm_version=$helm_version
RUN curl -LO "https://get.helm.sh/helm-$helm_version-linux-amd64.tar.gz" && \
    mkdir -p "/usr/local/helm-$helm_version" && \
    tar -xzf "helm-$helm_version-linux-amd64.tar.gz" -C "/usr/local/helm-$helm_version" && \
    ln -s "/usr/local/helm-$helm_version/linux-amd64/helm" /usr/local/bin/helm3 && \
    rm -f "helm-$helm_version-linux-amd64.tar.gz"

# Install Helm 2
ARG helm_version=v2.17.0
LABEL helm_version=$helm_version
RUN curl -LO "https://get.helm.sh/helm-$helm_version-linux-amd64.tar.gz" && \
    mkdir -p "/usr/local/helm-$helm_version" && \
    tar -xzf "helm-$helm_version-linux-amd64.tar.gz" -C "/usr/local/helm-$helm_version" && \
    ln -s "/usr/local/helm-$helm_version/linux-amd64/helm" /usr/local/bin/helm2 && \
    rm -f "helm-$helm_version-linux-amd64.tar.gz"

RUN mkdir -p ~/.ssh && touch ~/.ssh/known_hosts
RUN ssh-keyscan www.gitlab.com >> ~/.ssh/known_hosts
RUN ssh-keyscan www.github.com >> ~/.ssh/known_hosts
RUN ssh-keyscan gitlab.com >> ~/.ssh/known_hosts
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts


RUN echo "alias helm=helm3" >> ~/.bashrc
RUN echo "alias HELM_VERSION='f(){ alias helm=helm\$@; }; f'" >> ~/.bashrc

WORKDIR /

COPY --from=build /repoear /repoear
COPY config/* ./config/

EXPOSE 8000

#USER nonroot:nonroot

ENTRYPOINT ["/repoear"]