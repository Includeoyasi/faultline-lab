# Inventory: проверка на Windows с Linux-контейнерами

20 октября 2026 года проверен код рабочего дерева перед исходным коммитом. [Нормализованные SHA256 исходников](source-hashes.json) связывают проверку с конкретными файлами; исходные полные отчёты остаются в игнорируемом `artifacts/`. Здесь сохранены выбранные события с SHA256 исходных отчётов, без секретов и машинных путей.

| Проверка | Фактический результат |
|---|---|
| Закреплённый Linux toolchain: генерация API, форматирование, tidy/verify, build, lint | PASS; lint: 0 замечаний |
| `go test -count=1 ./...`, `go test -race -count=1 ./...`, govulncheck | PASS; уязвимостей не обнаружено |
| Процессы каркаса, health, диагностика и остановка | PASS; Inventory в явном health-only режиме |
| [PostgreSQL и Inventory](windows-postgres.json) | PASS; проект `faultline-data-test-4591c9fe80d3aecb` |
| [Общий PostgreSQL/Kafka/Redis и Inventory](windows-infrastructure.json) | PASS; проект `faultline-data-test-0e9512f319183416` |
| Runtime smoke трёх контейнеров каркаса | PASS; запуск `db74e00c-b6a6-41fa-9729-2fb381727ad7` |
| Очистка каждого disposable проекта | PASS; контейнеры, сети, тома и временные image tags отсутствуют |

## Что проверено в PostgreSQL и через gRPC

- 100 конкурентных попыток купить 10 единиц дали ровно 10 резервов; остаток неотрицателен и баланс сохранён.
- Многопозиционный отказ не списывает доступную первую позицию; сохранённый отказ не меняется после пополнения.
- Конкурентный повтор ключа даёт один эффект; изменённый payload отклоняется, новый ключ не создаёт второй резерв заказа.
- Release до Reserve и их гонка сохраняют барьер; Commit/Release выбирают один терминальный переход. Повтор Commit не продаёт дважды, Release после продажи отклоняется.
- Контролируемый отказ вставки outbox откатывает остаток, резерв, операцию и ledger. Отмена запроса, ожидающего DB-блокировку, не оставляет эффект и не ломает пул.
- Отдельный настоящий процесс Inventory переживает перезапуск с прежним результатом операции. После рестарта PostgreSQL сохраняются барьер, ключ, ledger и событие.
- Через отдельный процесс проверены Reserve, CommitReservation и GetAvailability; через обычный scratch-контейнер с `compose.inventory.yaml` и одним прикладным секретом — Reserve, ReleaseReservation и GetReservation. Так все пять RPC проходят транспортную проверку. Повтор через контейнер после рестарта PostgreSQL сохраняет исторический результат и текущее RELEASED.

## Воспроизведение

Из корня checkout с доступным Docker:

```text
docker compose build toolchain
docker compose run --rm toolchain check
docker compose run --rm toolchain service-smoke
docker compose run --rm -e GOOS=windows -e GOARCH=amd64 toolchain data-build
```

Затем на Windows:

```powershell
.\bin\labdata-windows-amd64.exe smoke
.\bin\labdata-windows-amd64.exe smoke -infrastructure
```

Host Go использовался для сборки локальных CLI; PostgreSQL и бизнес-процесс проверялись в Linux-контейнерах. Для связанного Windows worktree локальный toolchain получал явное read-only подключение Git-метаданных. Эти особенности не меняют бизнес-конфигурацию стенда.

После проверки Docker Desktop остановлен, его WSL VM завершена; процессов Docker/vmmemWSL не осталось, `autoStart=false`. Полезные кеши и установленные инструменты сохранены. Результаты AMD64/ARM64 CI относятся к конкретному HEAD и доступны в PR; этот документ не подменяет их. Личный macOS runtime, native/emulation-профиль, capacity и полный F10 трёх сервисов здесь не подтверждаются.
