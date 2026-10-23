# Карта документации

**Faultline Lab** · спецификации, решения и доказательства. Актуальное выполнение — в [GitHub Issues](https://github.com/Includeoyasi/faultline-lab/issues) и PR.

## Маршруты чтения

| Задача читателя | Маршрут |
| :--- | :--- |
| Быстро оценить проект | [README](../README.ru.md) → [Дизайн](architecture/DESIGN.md) → [Roadmap](ROADMAP.md) |
| Понять цели владельца | [Заметки об использовании](PRODUCT_INTENT.md) — цели, примеры и ограничения из обсуждений |
| Проверить распределённый процесс | [Атлас G0](architecture/PROCESSES.md) → [Принятие G0](adr/0001-approved-g0.md) |
| Разобрать подключение другого бизнеса | [Граница лаборатории](architecture/LAB_BOUNDARY.md) — согласованное направление и детализация для ревью |
| Подготовить среду Windows/macOS | [Контейнерная среда](development/ENVIRONMENT.md) — состав, профили и статус проверок |
| Собрать и проверить первые Go-бинарники | [Toolchain P01-01](development/TOOLCHAIN.md) — версии, Docker-команды и границы проверок |
| Запустить и проверить три процесса | [Запуск сервисов](development/SERVICE_RUNTIME.md) — настройки, JSON, health и границы shutdown |
| Поднять три runtime-контейнера | [Контейнерный каркас](development/CONTAINER_RUNTIME.md) — сборка, health, остановка и smoke с отчётом |
| Подготовить базы и исходные данные | [PostgreSQL](development/POSTGRES.md) — роли, миграции, seed, схемы взаимодействий и reset |
| Запустить инфраструктуру вместе | [Kafka и Redis](development/INFRASTRUCTURE.md) — общий smoke, ACL, TTL, сохранность Kafka и очистка |
| Прочитать платформы и ресурсы отчёта | [Метаданные запуска](development/RUN_METADATA.md) — происхождение значений, неизвестные поля и ограничения |
| Разобрать смену Kafka publisher и повторы | [Публикация Kafka](architecture/KAFKA_PUBLISHING.md) — проблема, три механизма и совместные сценарии на схемах |
| Проверить Windows smoke | [Отчёт P01](development/evidence/p01/WINDOWS_SMOKE.md) — исходные JSON/JSONL, бюджет и остановка среды |
| Проверить исходные ресурсы | [Профили сред](development/RESOURCE_PROFILES.md) — Windows, macOS и возможный будущий сервер |
| Посмотреть оформление | [Концепции баннера](design/BANNER_OPTIONS.md) — три варианта и выбранное решение |
| Принять участие в разработке | [CONTRIBUTING](../CONTRIBUTING.md) → [Правила ревью](development/CODE_REVIEW.md) → [Workflow](development/WORKFLOW.md) → [AGENTS](../AGENTS.md) |
| Разобрать принятый каркас | [P01](development/packets/P01.md) — границы, критерии и датированная приёмка |
| Выбрать следующий объём | [P02](development/packets/P02.md) — две задачи, порядок, границы PR и критерии |
| Следить за выполнением | [Учёт задач](development/TASK_TRACKING.md) — GitHub Issues/PR как единственный трекер |
| Проследить инженерные решения | [Решения и ревью](development/ENGINEERING_DECISIONS.md) — вопросы владельца, обоснования, изменения и доказательства |
| Проверить отдельный backlog | [Перенос smoke на Go](development/backlog/GO_SMOKE.md) — задача #9, границы и критерии |
| Проверить права и защиту main | [Безопасность репозитория](development/REPOSITORY_SECURITY.md) — действующая защита и переход на аккаунт исполнителя |
| Подготовить эксперимент | [Roadmap](ROADMAP.md) → [Шаблон отчёта](experiments/REPORT_TEMPLATE.md) |
| Ввести другого ИИ в задачу | [Переносимый контекст](../PROJECT_CONTEXT.md) — полный самостоятельный снимок |

## Источники решений

- [Дизайн](architecture/DESIGN.md): принятые границы сервисов, контракты, состояния и инварианты.
- [Процессы](architecture/PROCESSES.md): 15 подробных схем транзакций, отказов, гонок и восстановления.
- [Граница лаборатории](architecture/LAB_BOUNDARY.md): общий механизм, адаптеры, кандидат второго бизнеса и доказательство переносимости.
- [Среда](development/ENVIRONMENT.md): Docker для Windows/macOS, распределение контейнеров и ограничения измерений.
- [Ресурсные профили](development/RESOURCE_PROFILES.md): снимок Windows-хоста и текущий WSL-бюджет и отдельно ещё не измеренные среды.
- [Roadmap](ROADMAP.md): этапы 0–7, критерии готовности, риски и 12 сценариев экспериментов.
- [Workflow](development/WORKFLOW.md): роли человека и ИИ, состав ревью и контрольные точки G0–G4.
- [Правила ревью](development/CODE_REVIEW.md): стандарт Go, блокирующие замечания, проверки и критерии слияния.
- [ADR 0001](adr/0001-approved-g0.md): свидетельство явного принятия G0.
- [ADR 0002](adr/0002-record-review-decisions.md): согласованный порядок фиксации решений в PR, ADR и обзорной странице.
- [ADR 0003](adr/0003-go-container-smoke.md): принятое решение о Go smoke и его последствия.
- [Структурированное состояние](../CONTEXT_STATE.json): принятые решения, датированные проверки и история ревизий.

`PROJECT_CONTEXT.md` генерируется из источников и содержит их полные тексты. Источники редактируются напрямую; снимок обновляется через `tools/update_context.ps1`. Подробные документы пока ведутся на русском; публичное введение доступно на двух языках.
