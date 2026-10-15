FROM golangci/golangci-lint:v2.14.0@sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f AS lint

FROM golang:1.27.2-trixie@sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e AS toolchain
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
COPY . .
RUN case "$SERVICE" in orders|inventory|payments) ;; *) exit 2 ;; esac \
    && CGO_ENABLED=0 go build -trimpath -o /out/service "./cmd/$SERVICE"

FROM scratch AS runtime
COPY --from=service-build /out/service /app/service
USER 65532:65532
STOPSIGNAL SIGTERM
ENTRYPOINT ["/app/service"]
