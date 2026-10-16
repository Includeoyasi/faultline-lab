# Go-модуль и контейнерный toolchain

Закреплённый toolchain используется CI и сборкой runtime-образов. Команды и границы проверок описаны ниже; фактические результаты конкретного изменения находятся в PR и датированных evidence.

## Состав и версии

| Компонент | Источник закрепления | Версия |
| :--- | :--- | :--- |
| Go | [go.mod](../../go.mod), [Dockerfile](../../build/Dockerfile) | 1.27.2 |
| goimports | tool directive в go.mod, go.sum | golang.org/x/tools v0.52.0 |
| govulncheck | tool directive в go.mod, go.sum | golang.org/x/vuln v1.9.0 |
| golangci-lint | Dockerfile, официальный образ с digest | v2.14.0 |
| База toolchain | Dockerfile, индекс с AMD64/ARM64 | golang:1.27.2-trixie, sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e |
| Образ линтера | Dockerfile, индекс с AMD64/ARM64 | golangci/golangci-lint:v2.14.0, sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f |

Версии проверены по официальным источникам 15 октября 2026 года: [Go releases](https://go.dev/dl/), [golangci-lint releases](https://github.com/golangci/golangci-lint/releases/tag/v2.14.0), Go module proxy и Docker Registry. Оба закреплённых digest относятся к multi-platform индексам, а не к одному варианту CPU.

gRPC v1.84.0 добавлен для реального health API в P01-04; версия проверена через официальный Go proxy. JSON-логи использует стандартный log/slog. Инструменты закреплены через [Go tool directives](https://go.dev/doc/modules/managing-dependencies#tools). PostgreSQL использует pgx v5.11.0; Kafka/Redis библиотеки добавляются при первом использовании.

## Основной путь Windows и macOS

Нужны Git и Docker с Linux-контейнерами и Compose v2. Установка Go и Bash на хосте не требуется. Из корня репозитория:

```text
docker compose build toolchain
docker compose run --rm toolchain versions
docker compose run --rm toolchain check
docker compose run --rm toolchain build
docker compose run --rm toolchain service-smoke
```

`versions` выводит фактические Go, golangci-lint, build metadata goimports и версию govulncheck. `check` проверяет gofmt/goimports без перезаписи исходников, tidy без diff, verify, build, согласованные линтеры, unit/race commands и govulncheck. Ошибка любого шага завершает команду с ненулевым кодом. Набор линтеров — govet, staticcheck, errcheck, ineffassign, unused.

`build` создаёт Linux-бинарники bin/orders, bin/inventory, bin/payments для архитектуры toolchain-контейнера и выводит build metadata. Они игнорируются Git. На Windows эти Linux-бинарники не запускаются напрямую. `service-smoke` собирает и проверяет три отдельных процесса: health, JSON-диагностику ошибок, выход и независимость. Unit/process tests, graceful drain и runtime-контейнеры реализованы и проверены; бизнес-процессы остаются следующими пакетами. Прежний bootstrap-smoke удалён: сервисы больше не завершаются сразу после запуска.

Для осознанного исправления форматирования: `docker compose run --rm toolchain format`. Эта команда меняет Go-исходники в рабочей копии; check их не меняет. Неподдерживаемая команда и лишние аргументы дают exit code 2.

`data-build` собирает host CLI для [подготовки PostgreSQL](POSTGRES.md); `data-cross-check` проверяет компиляцию для Windows AMD64, Linux AMD64/ARM64 и macOS AMD64/ARM64. Host OS/архитектура задаются через GOOS/GOARCH.

Toolchain — временный сервис профиля tools. Явный `compose run toolchain` активирует его; команды не поднимают PostgreSQL/Kafka/Redis или runtime сервисов. Shell-скрипт выполняется внутри Linux-контейнера. Исходники подключаются из относительного пути, кеши модулей и сборки — в отдельных named volumes. Используются стандартные capabilities Docker для записи в bind mount; privileged mode и Docker socket не включены, no-new-privileges включён. Ограничения runtime-контейнеров сервисов описаны отдельно в [контейнерном каркасе](CONTAINER_RUNTIME.md).

Внутри образа Git доверяет только рабочему каталогу /workspace: bind mount может иметь другого владельца, а Go читает Git для настоящих VCS build metadata. Настройка действует в образе инструментов и не меняет Git-конфигурацию хоста.

## CI и границы доказательств

[Go workflow](../../.github/workflows/go.yml) использует нативные ubuntu-24.04 (AMD64) и ubuntu-24.04-arm (ARM64). На каждой архитектуре он собирает один и тот же Dockerfile, выполняет check, service-smoke и проверяет отказ на неизвестную команду. Затем проверяет отсутствие изменений отслеживаемых исходников. Успешная сборка ARM64 в CI не доказывает запуск Docker Desktop на личном Mac.

На Windows установлены host Go1.27.2, WSL2.7.13 и Docker Desktop4.95. [Runtime smoke и shutdown](evidence/p01/WINDOWS_SMOKE.md) PASS; основной путь сборки остаётся контейнерным. Личная macOS не проверена, performance capacity не измерена. После работы Docker Desktop/WSL VM выключаются; установленные инструменты и полезные кеши сохраняются.

Версии, директивы tools и образы обновляются совместно отдельным проверяемым изменением. GOTOOLCHAIN=local запрещает неявную загрузку другого Go, контейнер использует GOFLAGS=-mod=readonly. После изменения go.mod/go.sum контейнер пересобирается. Текущая защита main и required checks отражаются в [безопасности репозитория](REPOSITORY_SECURITY.md), а не выводятся из наличия workflow.
