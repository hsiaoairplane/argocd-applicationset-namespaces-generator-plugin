FROM --platform=$BUILDPLATFORM golang:1.24 AS build

ARG BUILDPLATFORM
ARG TARGETARCH
ARG VERSION

COPY . /src
RUN cd /src && GOOS=linux GOARCH=$TARGETARCH go build -a -ldflags="-X 'github.com/hsiaoairplane/argocd-applicationset-namespaces-generator-plugin/cmd/version.Version=$VERSION'" -o /bin/argocd-applicationset-namespaces-generator-plugin

FROM golang:1.24

COPY --from=build /bin/argocd-applicationset-namespaces-generator-plugin /usr/local/bin/argocd-applicationset-namespaces-generator-plugin

RUN useradd -s /bin/bash -u 999 argocd
WORKDIR /home/argocd
USER argocd

ENTRYPOINT ["/usr/local/bin/argocd-applicationset-namespaces-generator-plugin", "server"]
