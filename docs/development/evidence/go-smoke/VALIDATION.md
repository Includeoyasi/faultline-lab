# Проверка Go smoke — 16 октября 2026 года

Датированные результаты для [#9](https://github.com/Includeoyasi/faultline-lab/issues/9) и F01/F02/F04/F07/F08 из [#20](https://github.com/Includeoyasi/faultline-lab/issues/20). Это доказательства конкретного изменения; актуальное выполнение и приёмка находятся в GitHub.

## Среда и происхождение

Windows AMD64, Go 1.27.2; Docker Desktop 4.95, Linux backend 29.8.2/x86_64, Compose 5.6.0. Toolchain работает в контейнере. Ресурсы WSL ограничены двумя CPU и 2 GiB памяти; эти проверки не измеряют производительность бизнеса.

Исходники заключительного запуска: `8e49367c868af33023a6c1f7a8f94db447a828eb`, чистая рабочая копия. [Итоговый отчёт](windows-final-report.json) содержит `source_dirty=false`, реальные container/image IDs, health/security, изоляцию Orders, graceful stop с exit 0, OOM=false и cleanup=PASS. Рядом сохранены JSONL [Orders](windows-final-orders.jsonl), [Inventory](windows-final-inventory.jsonl), [Payments](windows-final-payments.jsonl).

## Фактические проверки

| Проверка | Результат и границы |
|---|---|
| `docker compose run --rm toolchain check` | PASS: format/imports, tidy/verify, build, lint, unit, race, govulncheck. Linux AMD64 внутри Docker Desktop. |
| `docker compose run --rm toolchain service-smoke` | PASS: процессы трёх сервисов, health, изоляция и остановка. Существующие Go-тесты подтверждают завершение активного RPC и ограниченный forced stop зависшего RPC/Watch. |
| `docker compose run --rm toolchain smoke-cross-check` | PASS: Windows AMD64, Linux AMD64/ARM64, macOS AMD64/ARM64. Это проверка сборки, а не runtime всех платформ. |
| Нативный Windows `go test -count=1 -timeout=90s ./internal/smoke ./cmd/smoke` | PASS, включая завершение дерева процессов через Job Object. |
| `bin/smoke-windows-amd64.exe` | PASS на чистом исходном коммите; все собственные контейнеры/сеть удалены. |
| Windows `CTRL_BREAK_EVENT` на этапе check | Ожидаемый FAIL, `cancelled=true`, cleanup=PASS, собственных ресурсов не осталось; независимый контрольный контейнер продолжил работу. [Наблюдения](go-smoke-windows-interrupt.json), [отчёт](windows-interrupt-report.json). |
| Отмена build/start/check/stop и ошибки команд/cleanup | PASS в Go fixtures с явной моделью Docker. Они не выдаются за четыре реальных OS-interrupt прогона. Проверены независимый cleanup context, deadline, позднее появление собственного ресурса и сохранение чужого. |
| Эффективный Compose-профиль | 20s/15s отклонён до запуска; 20s/25s принят; намеренный 20s/15s допустим только через явный unsafe preflight. [Наблюдения](go-smoke-preflight.json). |
| Backend выключен | Ожидаемый NOT_RUN / exit 2 / cleanup=NOT_NEEDED, [отчёт](no-backend-report.json). Новые сервисы не запускались. |
| Гигиена путей | `node --test tools/doc_safety.test.mjs`: PASS. Windows, JSON-escaped Windows, macOS/Linux homes отклоняются; контейнерные/относительные пути и HTTP URL допускаются. |
| Переносимый контекст | Генерация и `-CheckOnly` в PowerShell 7 и Windows PowerShell 5.1; изменяемый `snapshot_ref=main` отклонён. SHA исходных материалов проверяется через Git. |

Тест Unix process group отдельно проверяет потомка, удерживающего TCP-соединение после успешного выхода родителя с перенаправленным stdio. Executor завершает остаток своей группы. Ограничение буфера проверяется через обычный Write и `io.Copy`; причина ошибки executor сохраняется с редактированием домашних путей, sensitive config stdout/stderr не публикуются.

## Сравнение с прежним runner

На одной Windows/Docker среде последовательно выполнены прежний Node runner и новый Go runner: оба PASS, cleanup PASS, одинаковые health/stop/exit/OOM/isolation результаты. [Сводка](go-smoke-equivalence.json), [Node report](node-baseline-report.json), [Go report](go-equivalence-report.json).

Эти два запуска и прерывание выполнены при HEAD `511b9a85c7006cc7042263e0b87e7eb7fb8b116c` с незакоммиченными изменениями (`source_dirty=true`), до заключительного чистого запуска. Image IDs различаются после пересборки; идентичность образов или сравнительная производительность не заявляются. Старый runner удалён после проверки замены; исходная версия доступна в истории Git.

## Ограничения и завершение локальной работы

Linux host smoke и настоящий SIGINT запускаются отдельными CI jobs для нативных AMD64/ARM64; их результаты и артефакты находятся в PR. Проверки выше не подменяют эти прогоны. macOS runtime — **NOT_RUN**, доступа к Mac в этом запуске не было. Полный контракт платформ F03 остаётся отдельным решением; поле backend не обозначает архитектуру каждого image.

После локальных проверок собственных контейнеров осталось 0, Docker Desktop выключен, Docker WSL VM и процессы разработки завершены. Автозапуск отключён; программы, образы и полезные кеши сохранены. [Подтверждение](windows-cleanup.json), время наблюдения указано внутри. Проверка NOT_RUN после выключения backend не поднимает его снова.

Полный разрешённый Compose config и inspect Env не сохранялись. Публикуемые JSON/JSONL не содержат персональных путей; произвольные сторонние логи требуют отдельной проверки перед публикацией. Для SIGKILL/аварии хоста восстановление выполняется адресно по `recovery_command`, без обещания запуска signal handler.
