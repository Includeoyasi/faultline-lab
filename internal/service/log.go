package service

import (
	"io"
	"log/slog"
	"runtime/debug"
)

func newLogger(output io.Writer, name string, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: level})).With("service", name)
}

func buildAttributes(info *debug.BuildInfo, ok bool) []slog.Attr {
	version, revision, modified, goVersion := "unknown", "unknown", "unknown", "unknown"
	if ok && info != nil {
		if info.Main.Version != "" {
			version = info.Main.Version
		}
		if info.GoVersion != "" {
			goVersion = info.GoVersion
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.modified":
				modified = setting.Value
			}
		}
	}
	return []slog.Attr{
		slog.String("version", version), slog.String("revision", revision),
		slog.String("vcs_modified", modified), slog.String("go_version", goVersion),
	}
}
