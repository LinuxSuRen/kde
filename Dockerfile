# Build the manager binary
ARG RUNTIME=ghcr.io/linuxsuren/distroless/static:nonroot
ARG BUILDER=ghcr.io/linuxsuren/library/golang:1.22
ARG NODE=ghcr.io/linuxsuren/library/node:22
FROM ${NODE} AS node
WORKDIR /workspace
COPY ui/kde-ui/ .
RUN npm install && npm run build

FROM ${BUILDER} AS builder
ARG TARGETOS
ARG TARGETARCH
ARG GOPROXY

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
ENV GOPROXY=${GOPROXY}
RUN go mod download

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY config/ config/
COPY internal/ internal/
COPY pkg/ pkg/
COPY main.go main.go
COPY --from=node /workspace/dist ui/kde-ui/dist
COPY ui/kde-ui/data.go ui/kde-ui

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o server main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM ${RUNTIME}
WORKDIR /
COPY --from=builder /workspace/manager .
COPY --from=builder /workspace/server .
USER 65532:65532

ENTRYPOINT ["/manager"]
