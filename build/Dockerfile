FROM golangci/golangci-lint:v2.14.0@sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f AS lint

FROM golang:1.27.3-trixie@sha256:604106102d33b1084e37ca08adf98ac0813624ae8595bd0d7f3965ecebded09c AS toolchain
ARG TARGETARCH
RUN apt-get update && apt-get install -y --no-install-recommends unzip \
    && case "$TARGETARCH" in \
       amd64) archive=linux-x86_64; checksum=e9c129c176bb7df02546c4cd6185126ca53c89e7d2f09511e209319704b5dd7e ;; \
       arm64) archive=linux-aarch_64; checksum=4a802ed23d70f7bad7eb19e5a3e724b3aa967250d572cadfd537c1ba939aee6a ;; \
       *) exit 2 ;; esac \
    && curl -fsSL "https://github.com/protocolbuffers/protobuf/releases/download/v32.1/protoc-32.1-$archive.zip" -o /tmp/protoc.zip \
    && echo "$checksum  /tmp/protoc.zip" | sha256sum -c - \
    && unzip -q /tmp/protoc.zip -d /opt/protoc \
    && ln -s /opt/protoc/bin/protoc /usr/local/bin/protoc \
    && rm /tmp/protoc.zip && rm -rf /var/lib/apt/lists/*
ENV GOTOOLCHAIN=local \
    GOFLAGS=-mod=readonly \
    GOBIN=/usr/local/bin
WORKDIR /workspace
COPY --from=lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint
COPY go.mod go.sum ./
RUN git config --system --add safe.directory /workspace \
    && go mod download && go install tool
COPY --chmod=0755 tools/go.sh /usr/local/bin/lab-go
ENTRYPOINT ["/usr/local/bin/lab-go"]
CMD ["versions"]

FROM toolchain AS service-build
ARG SERVICE
COPY cmd ./cmd
COPY api ./api
COPY internal ./internal
COPY migrations ./migrations
RUN case "$SERVICE" in orders|inventory|payments) ;; *) exit 2 ;; esac \
    && CGO_ENABLED=0 go build -trimpath -o /out/service "./cmd/$SERVICE"

FROM scratch AS runtime
COPY --from=service-build /out/service /app/service
USER 65532:65532
STOPSIGNAL SIGTERM
ENTRYPOINT ["/app/service"]

FROM toolchain AS labdata-build
COPY cmd ./cmd
COPY api ./api
COPY internal ./internal
COPY migrations ./migrations
RUN CGO_ENABLED=0 go build -trimpath -o /out/labdata ./cmd/labdata \
    && CGO_ENABLED=0 go test -c -o /out/database.test ./internal/database \
    && CGO_ENABLED=0 go test -c -o /out/inventory.test ./internal/inventory/store \
    && CGO_ENABLED=0 go build -trimpath -o /out/inventory ./cmd/inventory

FROM scratch AS labdata
COPY --from=labdata-build /out/labdata /app/labdata
COPY --from=labdata-build /out/database.test /app/database.test
COPY --from=labdata-build /out/inventory.test /app/inventory.test
COPY --from=labdata-build /out/inventory /app/inventory
USER 65532:65532
ENTRYPOINT ["/app/labdata"]
