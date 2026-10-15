# Go-модуль и контейнерный toolchain — P01-01

Состояние на 15 октября 2026 года: созданы Go-модуль, три bootstrap entrypoints и контейнер инструментов. Локальная сборка временным Go на Windows прошла; контейнерные проверки Linux AMD64/ARM64 выполняются в CI, их фактические результаты сохраняются в контексте и PR. Это только первая задача P01; runtime сервисов и бизнес-API ещё не реализованы.

## Состав и версии

| Компонент | Источник закрепления | Версия |
| :--- | :--- | :--- |
| Go | [go.mod](../../go.mod), [Dockerfile](../../build/toolchain.Dockerfile) | 1.27.2 |
| goimports | tool directive в go.mod, go.sum | golang.org/x/tools v0.52.0 |
| govulncheck | tool directive в go.mod, go.sum | golang.org/x/vuln v1.9.0 |
| golangci-lint | Dockerfile, официальный образ с digest | v2.14.0 |
| База toolchain | Dockerfile, индекс с AMD64/ARM64 | golang:1.27.2-trixie, sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e |
| Образ линтера | Dockerfile, индекс с AMD64/ARM64 | golangci/golangci-lint:v2.14.0, sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f |

Версии проверены по официальным источникам 15 октября 2026 года: [Go releases](https://go.dev/dl/), [golangci-lint releases](https://github.com/golangci/golangci-lint/releases/tag/v2.14.0), Go module proxy и Docker Registry. Оба закреплённых digest относятся к multi-platform индексам, а не к одному варианту CPU.

Текущий код приложений использует стандартную библиотеку. Внешние модули в go.mod нужны для реально используемых инструментов, закреплённых современными [Go tool directives](https://go.dev/doc/modules/managing-dependencies#tools); ради go.sum в бизнес-код не добавлена фиктивная библиотека. Зависимости будущего gRPC/SQL/Kafka добавляются при их первом использовании.

## Основной путь Windows и macOS

Нужны Git и Docker с Linux-контейнерами и Compose v2. Установка Go и Bash на хосте не требуется. Из корня репозитория:

```text
docker compose build toolchain
docker compose run --rm toolchain versions
docker compose run --rm toolchain check
docker compose run --rm toolchain build
docker compose run --rm toolchain bootstrap-smoke
```

`versions` выводит фактические Go, golangci-lint, build metadata goimports и версию govulncheck. `check` проверяет gofmt/goimports без перезаписи исходников, tidy без diff, verify, build, согласованные линтеры, unit/race commands и govulncheck. Ошибка любого шага завершает команду с ненулевым кодом. Набор линтеров — govet, staticcheck, errcheck, ineffassign, unused.

`build` создаёт Linux-бинарники bin/orders, bin/inventory, bin/payments для архитектуры toolchain-контейнера и выводит build metadata. Они игнорируются Git. На Windows эти Linux-бинарники не запускаются напрямую. `bootstrap-smoke` собирает и запускает их внутри контейнера: каждый выводит сообщение об отсутствии service runtime и завершается. Это не проверка health или бизнес-процесса. Тестовых Go-файлов в этой задаче нет: вывод `[no test files]` не является доказательством покрытия будущих рисков.

Для осознанного исправления форматирования: `docker compose run --rm toolchain format`. Эта команда меняет Go-исходники в рабочей копии; check их не меняет. Неподдерживаемая команда и лишние аргументы дают exit code 2.

Toolchain — временный сервис профиля tools. Явный `compose run toolchain` активирует его; команды не поднимают PostgreSQL/Kafka/Redis или runtime сервисов. Shell-скрипт выполняется внутри Linux-контейнера. Исходники подключаются из относительного пути, кеши модулей и сборки — в отдельных named volumes. Контейнер не получает Docker socket и привилегии; Linux capabilities удалены.

Внутри образа Git доверяет только рабочему каталогу /workspace: bind mount может иметь другого владельца, а Go читает Git для настоящих VCS build metadata. Настройка действует в образе инструментов и не меняет Git-конфигурацию хоста.

## CI и границы доказательств

[Go workflow](../../.github/workflows/go.yml) использует нативные ubuntu-24.04 (AMD64) и ubuntu-24.04-arm (ARM64). На каждой архитектуре он собирает один и тот же Dockerfile, выполняет check, bootstrap-smoke и проверяет отказ на неизвестную команду. Затем проверяет отсутствие изменений отслеживаемых исходников. Успешная сборка ARM64 в CI не доказывает запуск Docker Desktop на личном Mac.

На Windows Docker пока не найден; локальный container smoke — NOT_RUN. Для дополнительной проверки скачан официальный переносимый Go 1.27.2 с проверкой SHA-256 в игнорируемую папку .local, без установки в PATH. Это вспомогательная проверка, а не основной путь разработки. Личная macOS-среда ещё не проверялась; измерений потребления и производительности нет.

Версии, директивы tools и образы обновляются совместно отдельным проверяемым изменением. GOTOOLCHAIN=local запрещает неявную загрузку другого Go, контейнер использует GOFLAGS=-mod=readonly. После изменения go.mod/go.sum контейнер пересобирается. Текущая защита main и required checks отражаются в [безопасности репозитория](REPOSITORY_SECURITY.md), а не выводятся из наличия workflow.
