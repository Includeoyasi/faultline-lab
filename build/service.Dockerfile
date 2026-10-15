# The named context is the pinned toolchain service in compose.yaml.
FROM toolchain AS build
ARG SERVICE
COPY . .
RUN case "$SERVICE" in orders|inventory|payments) ;; *) exit 2 ;; esac \
    && CGO_ENABLED=0 go build -trimpath -o /out/service "./cmd/$SERVICE"

FROM scratch AS runtime
COPY --from=build /out/service /app/service
USER 65532:65532
STOPSIGNAL SIGTERM
ENTRYPOINT ["/app/service"]
