# Orders: датированные доказательства — 23 октября 2026 года

Исходный код: `d09f1fe6652cc18cf88cd4ce30f604be379f2ef0`; проверенный HEAD вместе с контекстом: `1aec8ab098bccb9b7b501285fffda06360310d19`. Это сохранённый результат конкретного прогона; актуальные проверки и решение владельца находятся в [PR #29](https://github.com/Includeoyasi/faultline-lab/pull/29).

| Среда / проверка | Фактический результат |
|---|---|
| Windows AMD64, portable Go 1.27.3 | Orders/config/model/worker, datalab/smokecli, vet и build PASS; интеграционные случаи без FL_DATA_INTEGRATION пропускаются |
| Windows, отдельные health-only процессы | PASS после разрешённого запуска вне ограничения loopback; первый sandbox-прогон получил запрет сокета |
| Native Linux AMD64 и ARM64, Go 1.27.3 | [Go CI](https://github.com/Includeoyasi/faultline-lab/actions/runs/39989398467) PASS: генерация, форматирование, tidy/verify, build, lint, unit/race, govulncheck, сервисный и контейнерный smoke |
| Реальные PostgreSQL + Inventory, оба native профиля | PostgreSQL-only и совместный PostgreSQL/Kafka/Redis PASS |
| Документы и контекст | [Documentation](https://github.com/Includeoyasi/faultline-lab/actions/runs/39989398444) PASS |

В обоих режимах каждый native job выполнял [TestOrdersFoundation и TestOrdersProcessCrash](../../../../internal/orders/acceptance/foundation_test.go) с явным FL_DATA_INTEGRATION=1. Проверены конкурентный повтор Create и конфликт payload, откат при ошибке INSERT job/outbox, поколения worker и истечение lease во время ожидания строки заказа. Сценарий G1 принимает 1000 запросов на 100 единиц и проверяет 100 PAYMENT_PENDING, 900 сохранённых INSUFFICIENT_STOCK, доступность 0 и сумму резервов 100. Настоящий процесс Orders получает SIGKILL до первого Reserve либо после commit Inventory; после запуска прежняя операция даёт один складской эффект и одно локальное событие. Исторический успешный Reserve с текущим RELEASED оставляет диагностируемый незавершённый заказ.

Scratch-контейнер Orders проверен до и после перезапуска PostgreSQL: неизменный ответ Create, текущее чтение обоими API, конфликт изменённого payload и подтверждённый резерв. Последующие host CLI tests проверили фактические образы, ресурсные снимки и метаданные Orders/Inventory во всех отчётах; адресная очистка собственных ресурсов завершилась успешно.

Отрисовка документов, 23 октября: Mermaid 11.4.1 воспроизвёл parse error на неэкранированной точке с запятой в подписи потерянного ответа. После замены разделителей обе схемы — описание PR и руководство Orders — успешно отрисованы в SVG/PNG и визуально проверены. Структурная проверка Documentation сама по себе не подтверждает отрисовку Mermaid; эта проверка выполнена отдельно и не меняет результаты бизнес-сценариев выше.

Границы: takeover моделируется реальными SQL Claim/Resolve разных поколений, а не двумя одновременно работающими процессами Orders. Потеря именно ответа Create не инъецируется; проверяется неизменный повтор после прогресса и рестарта. Отсутствие SQL-транзакции во время RPC следует из кода, отдельного измерительного теста нет. Docker/WSL на личном Windows ПК для этого среза не запускались; Windows Docker и macOS runtime — NOT_RUN. Эти функциональные результаты не являются capacity/A/B-измерением или приёмкой G1 владельцем.
