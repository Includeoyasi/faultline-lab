# PROJECT_CONTEXT — Faultline Lab

Самостоятельный переносимый контекст для ИИ-агентов и чатов ChatGPT. Полные источники включены ниже; локальный компьютер для чтения не требуется.

Ревизия: **51**. Дата снимка: **2026-10-16**, Europe/Moscow. Исходный коммит внешних материалов: `1c54d9c6aee2be4da7c38b7775b6d862402caac1`.

## Выполнение задач

Актуальные статусы, блокеры и ревью ведутся только в GitHub: [Issues](https://github.com/Includeoyasi/faultline-lab/issues), [Pull requests](https://github.com/Includeoyasi/faultline-lab/pulls), [Milestones](https://github.com/Includeoyasi/faultline-lab/milestones). Этот файл хранит спецификации, принятые решения и датированные доказательства; он не копирует текущую доску задач.

## Назначение переносимого контекста

Это самостоятельный пакет спецификаций, решений и датированных доказательств Faultline Lab. Его можно читать без исходного компьютера, передавать пользователю или ИИ-тестировщику. Полные зарегистрированные документы и исходники включены в приложения; SHA-256 позволяет проверить соответствие исходным файлам.

**Актуальное выполнение ведётся только в GitHub Issues, PR и milestones.** Этот снимок не содержит отдельной доски задач, счётчиков прогресса или назначения «следующей задачи». Успешная генерация подтверждает соответствие локальным источникам, а не актуальность GitHub или работоспособность сервисов. Датированные результаты проверок применимы только к указанным версиям, среде и commit.

## Проект и роли

Faultline Lab — лаборатория масштабирования, отказов и восстановления вокруг небольшой бизнес-системы. Публичный репозиторий: [Includeoyasi/faultline-lab](https://github.com/Includeoyasi/faultline-lab). Владелец выбрал MIT 13 октября 2026 года. Прежнее имя Microservices Lab обозначает тот же проект.

Владелец проекта определяет направление, проводит ревью и принимает ключевые решения. Исполнитель — Codex, публикующий изменения от отдельного аккаунта **AI-Your-Bro**. Другие ИИ-агенты рассматриваются как потенциальные пользователи, подписчики или тестировщики. Передача контекста не даёт им поручения разрабатывать код.

Исполнитель не одобряет и не сливает собственные PR. Автор замечания из чата — владелец; исполнитель может пересказать решение со ссылкой на чат как источник, но не выдаёт его за GitHub approval. Требования защиты main и CODEOWNERS описаны в приложениях. Автор Git-коммита не заменяет проверку реально авторизованного API actor.

## Как читать и проверять сведения

1. Для понимания проекта прочитайте README, принятый дизайн и атлас процессов.
2. Для поведения системы используйте полные спецификации и соответствующие исходники, а не предположение по дереву каталогов.
3. Для выполнения задачи, блокеров и ревью откройте GitHub. При отсутствии доступа явно укажите это ограничение; не восстанавливайте актуальный прогресс из старой истории.
4. Для анализа измерений проверьте дату, исходный commit, бюджет среды, команды и ограничения доказательства.
5. Различайте принятое направление, текст спецификации, применённый код, выполненную проверку и приёмку человеком.

Приоритет внутри проектных решений: последнее прямое указание владельца, датированные согласования/ADR, принятые спецификации. Историческое предложение не отменяет более позднего решения. Если документы расходятся, назовите конкретное противоречие и основание выбора.

Владелец полностью принял **G0 13 октября 2026 года**, прямое сообщение: «G0 полностью принят». Это принятие архитектуры и правил, не доказательство реализации. G0–G4 — контрольные точки приёмки; P01–P05 — пакеты изменений между контрольными точками. Первый кодовый P01 расположен после G0 и перед G1. Приёмка PR сама по себе не означает приёмки следующей контрольной точки.

## Архитектурные границы

Ровно три бизнес-сервиса: Orders, Inventory и Payments. Сохранённые шаги оркестрации находятся внутри Orders. Инфраструктура и инструменты лаборатории не становятся дополнительными бизнес-сервисами. Go/gRPC, PostgreSQL, Kafka, Redis, Prometheus/Grafana и Loki/Alloy составляют планируемый стек; компоненты вводятся по этапам roadmap.

Общая лаборатория отделена от тестируемой системы и адаптера бизнес-операций/проверок. Магазин подключается через ту же границу. Обработка изображений — кандидат другого бизнеса; автоматический мастер подключения — перспектива. Детальные контракты и доказательство переносимости находятся в приложениях; пример не объявляется принятой реализацией.

Бизнес-изменение и outbox фиксируются атомарно в БД владельца. Повторная доставка и неизвестный результат RPC ожидаемы. Конечные состояния подтверждаются данными оплаты и резерва; задержанный запрос не должен отменять сохранённый запрет операции. Redis не разрешает продажу. Полные автоматы и инварианты имеют приоритет над этим кратким введением.

## Среда и безопасность выполнения

Согласован контейнерный путь сборки и запуска на Windows/macOS, без обязательного host Go SDK. Применённые конфигурации, версии и результаты измерений сохраняются с датами. Возможный арендуемый сервер — будущее расширение, не поручение на аренду. Фиксированный бюджет Docker не является измерением capacity бизнеса.

Адреса listener, peer DNS и публикация портов задаются раздельно в deployment-конфигурации. `FL_GRPC_ADDR` обязателен, скрытого fallback в прикладном коде нет. Loopback в тестовом fixture не является production-настройкой. Секреты, полный inspect Env и личные пути не публикуются.

На личном Windows-компьютере после проверок нужно удалить только ресурсы собственного теста, выключить Docker Desktop и его WSL VM, проверить остановку и сохранить выключенный автозапуск. Установленные инструменты, образы и полезные кеши остаются. Чужие workloads нельзя останавливать общей очисткой.

## Аудит и принятые направления

16 октября владелец передал внешний аудит v2, затем выбрал направления F01–F12. Их причины, критерии и отложенные вопросы зафиксированы в [Issue #20](https://github.com/Includeoyasi/faultline-lab/issues/20). Перенос smoke на Go связан с [#9](https://github.com/Includeoyasi/faultline-lab/issues/9); вопросы Kafka — с [#16](https://github.com/Includeoyasi/faultline-lab/issues/16). Статусы этих работ берутся только из трекера.

В структурированных решениях сохранены выбранные варианты и первоначальные предложения с явным различием. Отложенный вопрос не считается решённым из-за наличия описания. Чужие сохранённые тесты не превращаются в свежий результат; приложенные аудиторские скрипты сами по себе не являются поручением на запуск.

## Обновление снимка

Редактируются CONTEXT_STATE.json, это введение и канонические приложения; PROJECT_CONTEXT.md генерируется. После содержательного решения или проверки увеличивается revision, фиксируются источник и дата Europe/Moscow, выполняются генерация и CheckOnly. Изменение стадии Issue не требует синхронно редактировать копии прогресса в документах.

Внешние материалы снимка закрепляются за существующим исходным SHA; ссылки на актуальный проект отделены. Снимок не может содержать hash собственного будущего коммита: сначала фиксируются исходные материалы, затем их SHA записывается в метаданные и создаётся обновлённый контекст. Локальный git HEAD и статус ревью проверяются отдельно.

## Решения, доказательства и история

JSON ниже включён для автономного чтения решений и истории. Даты и исходные commit обязательны при использовании сохранённых наблюдений; они не означают актуальный статус задачи. Отсутствие измерения не является нулевым значением или успешной проверкой.

```json
{
  "schema_version": 2,
  "project_id": "faultline-lab",
  "context_revision": 51,
  "context_created_on": "2026-10-13",
  "context_updated_on": "2026-10-16",
  "user_timezone": "Europe/Moscow",
  "language": "ru",
  "repository": {
    "name": "Faultline Lab",
    "slug": "faultline-lab",
    "owner": "Includeoyasi",
    "url": "https://github.com/Includeoyasi/faultline-lab",
    "visibility": "PUBLIC",
    "default_branch": "main",
    "description": "A Go microservices lab for testing scale, failure and recovery — built through human-led AI collaboration.",
    "license": "MIT",
    "license_approval": "Прямое сообщение владельца проекта 2026-10-13: возьми MIT",
    "first_commit": "6897f4ae1596c7e9e2b44e64ca8974e78fad54b8",
    "snapshot_ref": "1c54d9c6aee2be4da7c38b7775b6d862402caac1",
    "snapshot_ref_semantics": "Existing source commit for immutable external materials; not a claim about current branch HEAD or review status"
  },
  "roles": {
    "human": "Владелец проекта; ревьюер при проверке кода",
    "ai": "Исполнитель разработки, проверок, экспериментов и документации",
    "terminology_confirmed_on": "2026-10-16",
    "authority_changed": false,
    "github_owner": "Includeoyasi",
    "github_executor": "AI-Your-Bro",
    "github_identity_transition_confirmed_on": "2026-10-13",
    "other_ai_agents": "Potential users, subscribers or testers; implementation is assigned to the executor",
    "clarified_on": "2026-10-16",
    "development_executor": "AI-Your-Bro (Codex)"
  },
  "tracker_navigation": {
    "authority": "GitHub is the only source of live issue, blocker and review status; no local progress cache",
    "issues": "https://github.com/Includeoyasi/faultline-lab/issues",
    "pull_requests": "https://github.com/Includeoyasi/faultline-lab/pulls",
    "milestones": "https://github.com/Includeoyasi/faultline-lab/milestones",
    "related_issues": [
      9,
      12,
      16,
      20
    ]
  },
  "accepted_architecture": [
    {
      "id": "G0",
      "status": "APPROVED",
      "approved_on": "2026-10-13",
      "approval_evidence": "Прямое сообщение пользователя: G0 полностью принят.",
      "conditions": []
    }
  ],
  "authorization_history": [
    {
      "date": "2026-10-16",
      "source": "Direct owner message in chat: take issue9 and small F02/F04/F07/F08 fixes next",
      "scope": [
        "Issue9 / F01",
        "F02",
        "F04",
        "F07",
        "F08"
      ],
      "implementation_and_review_pr_authorized": true,
      "P02_implementation_authorized": false,
      "other_audit_findings": "Follow their agreed future checkpoints in issue20"
    },
    {
      "date": "2026-10-16",
      "source": "Direct owner chat after PR21 merge: take issue12; explain additions and seed with interaction diagrams; remove obsolete Docker data",
      "scope": [
        "Issue12 PostgreSQL foundation",
        "Detailed operational guide and diagrams",
        "Scoped disposable-image cleanup"
      ],
      "implementation_and_review_pr_authorized": true,
      "other_work": "Issue16 and remaining audit findings keep their agreed checkpoints"
    }
  ],
  "review_history": [
    {
      "number": 1,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/1",
      "author": null,
      "head_ref": "docs/go-review-policy",
      "base_ref": "main",
      "head_at_merge": "214eb34a160a2e331ecb45bca3f5f11fdeb5a73b",
      "merge_commit": "c1e55fdfca98e12e518f722300bd85ed85c1d85e",
      "merged_at_utc": null,
      "merged": true,
      "human_review": "ACCEPTED_BY_PROJECT_OWNER_ACCOUNT_MERGE",
      "scope": null,
      "observation_note": null
    },
    {
      "number": 2,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/2",
      "author": null,
      "head_ref": "docs/review-policy-handoff",
      "base_ref": "main",
      "head_at_merge": "681660a2a659d68c0336d05c0ac11db70d0c76d2",
      "merge_commit": "1a0b3f356529e262fe902a1e419e16169e8c6b08",
      "merged_at_utc": null,
      "merged": true,
      "human_review": "MERGED_REPOSITORY_STATE_OBSERVED",
      "scope": "Отражение принятия PR #1 и обновление переносимого контекста",
      "observation_note": null
    },
    {
      "number": 4,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/4",
      "author": null,
      "head_ref": "docs/portable-lab-plan",
      "base_ref": "main",
      "head_at_merge": "83e00760c67f5aa79940bda05444c2c6838e84a0",
      "merge_commit": "8c392b8ed4145d8b2b9efa93c5be14d9fdf113d1",
      "merged_at_utc": "2026-10-13T13:50:50Z",
      "merged": true,
      "human_review": "MERGED_REPOSITORY_STATE_OBSERVED",
      "scope": "Согласованные документы, граница лаборатории, Docker-план, ресурсные профили, выбранный баннер и P01 backlog",
      "observation_note": "GitHub REST подтвердил merge; ИИ не выполнял approve/merge. Переоформлять этот PR не требуется."
    },
    {
      "number": 5,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/5",
      "author": "AI-Your-Bro",
      "head_ref": "chore/executor-security",
      "base_ref": "main",
      "head_at_merge": "704a7ecae42372dbbe89de1d20fa7afcfd274491",
      "merge_commit": "d616a221566be4f7f7be89f6f3223501a539c45c",
      "merged_at_utc": null,
      "merged": true,
      "human_review": "USER_CONFIRMED_AND_GITHUB_MERGE_OBSERVED",
      "scope": "Защита main, CODEOWNERS, разделение GitHub identity и обновление контекста",
      "observation_note": "Владелец сообщил о merge; API подтвердил. ИИ не approve/merge."
    },
    {
      "number": 6,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/6",
      "author": "AI-Your-Bro",
      "head_ref": "feat/p01-01-toolchain",
      "base_ref": "main",
      "head_at_merge": "001d333132613f642b9f2415865bcd1255cb397a",
      "merge_commit": "e4a9b0b80f97b2fb20ca973ef57030d32f8ba878",
      "merged_at_utc": null,
      "merged": true,
      "human_review": "USER_CONFIRMED_AND_GITHUB_MERGE_OBSERVED",
      "scope": "P01-01: Go module, bootstrap entrypoints, pinned container toolchain",
      "observation_note": "Пользователь сообщил о merge; GitHub GET /pulls/6 подтвердил merged=true. ИИ не выполнял approve/merge."
    },
    {
      "number": 7,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/7",
      "author": "AI-Your-Bro",
      "head_ref": "feat/p01-02-04-service-runtime",
      "base_ref": "main",
      "head_at_merge": "2e8c7f8149f27c8fd9820401bd1e5d7e476ca0c9",
      "merge_commit": "bc09d08ce7867f3109d996f71015e3ce73986452",
      "merged_at_utc": null,
      "merged": true,
      "human_review": "USER_ACCEPTED_AND_GITHUB_MERGE_OBSERVED",
      "scope": "P01-02..P01-04: config, JSON logs, gRPC health",
      "observation_note": "Владелец сообщил о принятии; GitHub API merged=true. ИИ не approve/merge."
    },
    {
      "number": 8,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/8",
      "author": "AI-Your-Bro",
      "head_ref": "feat/p01-05-08-runtime",
      "base_ref": "main",
      "head_at_merge": "56f367278a0fabdd621425b38e13fcec323fb5c9",
      "merge_commit": "3ca290d818a98fc986fe42ae51dcfc0546673fb8",
      "merged_at_utc": "2026-10-15T21:10:58Z",
      "merged": true,
      "human_review": "ACCEPTED_BY_OWNER_ACCOUNT_MERGE",
      "scope": null,
      "observation_note": null
    },
    {
      "number": 11,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/11",
      "author": "AI-Your-Bro",
      "head_ref": "docs/p02-plan",
      "base_ref": "main",
      "head_at_merge": null,
      "merge_commit": null,
      "merged_at_utc": null,
      "merged": false,
      "human_review": "NO_CODE_REVIEW; OWNER_OBJECTED_TO_PREMATURE_PR_CREATION",
      "scope": "P02 planning and observed P01 acceptance only",
      "observation_note": null
    },
    {
      "number": 21,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/21",
      "author": "AI-Your-Bro",
      "head_ref": "feat/go-smoke-audit",
      "base_ref": "main",
      "head_at_merge": "4f283b78dae1edf62fe53ada39d6c3350c0c7883",
      "merge_commit": "843e95ed2b65a44729fbf297b986363c588d51ca",
      "merged_at_utc": "2026-10-16T15:14:01Z",
      "merged": true,
      "merged_by": "Includeoyasi",
      "human_review": "ACCEPTED_BY_PROJECT_OWNER_ACCOUNT_MERGE",
      "scope": "Issue9/F01: Go smoke; audit F02/F04/F07/F08; review wording and agreed process documentation",
      "observation_note": "Owner reported merge in chat; executor verified merged=true and actual merged_by via GitHub API. No acceptance of G1 or authorization of P02 implementation is inferred."
    }
  ],
  "review_policy": {
    "document": "docs/development/CODE_REVIEW.md",
    "status": "ACCEPTED_VIA_PROJECT_OWNER_REPOSITORY_MERGE",
    "requested_on": "2026-10-13",
    "request_evidence": "Прямое сообщение владельца проекта: Давай теперь установим классические для современной разработки на Go правила код-ревью",
    "human_approval": {
      "kind": "OBSERVED_REPOSITORY_MERGE",
      "actor": "Includeoyasi",
      "date_local": "2026-10-13",
      "merged_at_utc": "2026-10-12T22:47:27Z",
      "accepted_head": "214eb34a160a2e331ecb45bca3f5f11fdeb5a73b",
      "merge_commit": "c1e55fdfca98e12e518f722300bd85ed85c1d85e",
      "pull_request": "https://github.com/Includeoyasi/faultline-lab/pull/1",
      "detail": "GitHub API подтвердил слияние аккаунтом владельца проекта. ИИ не вызывал merge/approve; текстовое approval не приписывается пользователю."
    },
    "normal_change_delivery": "FOCUSED_BRANCH_AND_PR",
    "merge": "EXPLICIT_HUMAN_DECISION_AND_ACTION_AUTHORIZATION",
    "go_checks_activation": "ACTIVE_WITH_P01_01_FIRST_GO_MODULE; REQUIRED_NATIVE_AMD64_ARM64",
    "github_branch_protection": "ENABLED_AND_VERIFIED_2026_10_13",
    "pull_request": "https://github.com/Includeoyasi/faultline-lab/pull/1",
    "pull_request_author": "AI-Your-Bro",
    "github_approval_required": true,
    "chat_approval_does_not_bypass_github": true,
    "current_identity_transition_pr": "https://github.com/Includeoyasi/faultline-lab/pull/5"
  },
  "review_decision_tracking": {
    "requested_on": "2026-10-15",
    "policy_status": "ACCEPTED_IN_CHAT",
    "policy_evidence": "Owner replied Хорошо to the proposal: working discussion in PR, significant decisions in ADR, one engineering overview linked from README",
    "overview": "docs/development/ENGINEERING_DECISIONS.md",
    "policy_adr": "docs/adr/0002-record-review-decisions.md",
    "publication_pr": "https://github.com/Includeoyasi/faultline-lab/pull/8",
    "publication_review": "ACCEPTED_BY_INCLUDEOYASI_MERGE_PR8_ON_2026_10_16_LOCAL",
    "record_fields": [
      "Question and source",
      "Alternatives and rationale",
      "Change",
      "Verification evidence",
      "Status and acceptance source"
    ],
    "attribution_rule": "AI-Your-Bro publishes executor summaries; owner requirements from chat are identified as such. Do not fabricate GitHub comments or approvals by the project owner.",
    "routine_fixes": "PR only; ADR reserved for significant decisions; approved scope remains autonomous",
    "cases": [
      {
        "topic": "Explicit environment configuration instead of application network defaults",
        "status": "ACCEPTED_PR7_OWNER_MERGE",
        "source": "Owner review in chat",
        "evidence": "https://github.com/Includeoyasi/faultline-lab/pull/7#issuecomment-6192239609"
      },
      {
        "topic": "Go container smoke runner",
        "source": "Owner question in chat and explicit separate-backlog instruction",
        "evidence": "https://github.com/Includeoyasi/faultline-lab/issues/9",
        "decision_on_2026_10_16": "Owner selected Go migration and assigned issue9 in chat; live stage only in GitHub"
      },
      {
        "topic": "Personal PC runtime shutdown and bounded VM resources",
        "status": "ACCEPTED_PR8_OWNER_MERGE",
        "source": "Owner personal-computer instruction in chat; executor selected VM limits",
        "evidence": "docs/development/evidence/p01/WINDOWS_SMOKE.md"
      }
    ]
  },
  "repository_security": {
    "observed_on": "2026-10-15",
    "document": "docs/development/REPOSITORY_SECURITY.md",
    "main_protected": true,
    "source": "GitHub REST PUT branch protection followed by independent GET; branches/main confirms protected=true",
    "enforce_admins": true,
    "required_status_checks": {
      "strict": true,
      "checks": [
        {
          "context": "Links and portable context",
          "app_id": 15368
        },
        {
          "context": "Go toolchain (amd64)",
          "app_id": 15368
        },
        {
          "context": "Go toolchain (arm64)",
          "app_id": 15368
        }
      ]
    },
    "required_pull_request_reviews": {
      "required_approving_review_count": 1,
      "dismiss_stale_reviews": true,
      "require_last_push_approval": true,
      "require_code_owner_reviews": true
    },
    "required_conversation_resolution": true,
    "allow_force_pushes": false,
    "allow_deletions": false,
    "codeowners": {
      "path": ".github/CODEOWNERS",
      "owner": "Includeoyasi",
      "status": "ACTIVE_IN_MAIN",
      "effective_after": "HUMAN_REVIEW_AND_MERGE_INTO_MAIN",
      "pull_request_url": "https://github.com/Includeoyasi/faultline-lab/pull/5",
      "merge_commit": "d616a221566be4f7f7be89f6f3223501a539c45c",
      "verified_on": "2026-10-15"
    },
    "go_checks": "REQUIRED_NATIVE_TOOLCHAIN_AMD64_AND_ARM64_FROM_2026_10_15",
    "bypasses_added": false,
    "destructive_probes": "NOT_PERFORMED_CONFIGURATION_READ_BACK_ONLY",
    "executor_has_admin": false
  },
  "executor_identity": {
    "login": "AI-Your-Bro",
    "github_user_id": 342030453,
    "profile_url": "https://github.com/AI-Your-Bro",
    "profile_status": "VERIFIED_EXISTS_USER_PROVIDED_LOGIN",
    "profile_created_at_utc": "2026-10-13T13:58:33Z",
    "registered_by_ai": false,
    "machine_account_control": "PROJECT_OWNER",
    "terms": "GitHub requires human registration of machine accounts; automatic registration was not performed",
    "invitation": {
      "id": 337527166,
      "status": "ACCEPTED_WRITE_ACCESS_VERIFIED",
      "permissions": "write",
      "url": "https://github.com/Includeoyasi/faultline-lab/invitations"
    },
    "observed_repository_permission_before_invitation": "read",
    "authentication": "VERIFIED_GET_USER_AI_YOUR_BRO",
    "two_factor_authentication": {
      "status": "USER_DECLINED_ENABLEMENT",
      "api_status": "NOT_RETURNED_FOR_CURRENT_OAUTH_TOKEN",
      "date": "2026-10-13",
      "evidence": "Прямой ответ владельца: я не хочу ее включать"
    },
    "credentials_storage": "Windows Credential Manager; OAuth device flow completed; credentials not written to Git/context",
    "connector_actor": "Includeoyasi; do not create new PRs through this connector until executor authentication is available",
    "git_identity": {
      "scope": "LOCAL_REPOSITORY_ONLY",
      "name": "AI-Your-Bro",
      "email": "342030453+AI-Your-Bro@users.noreply.github.com",
      "credential_username": "AI-Your-Bro",
      "global_identity_changed": false
    },
    "new_pr_policy": "AUTHOR_MUST_BE_AI_YOUR_BRO; OWNER_REVIEWS; AI_NEVER_APPROVES_OWN_WORK",
    "verified_access": {
      "permission": "write",
      "push": true,
      "admin": false,
      "source": "GET repo collaborator permission authenticated as AI-Your-Bro"
    },
    "api_write_transport": "Executor credential requested by username; GET /user actor checked before every API mutation. Owner-connected GitHub integration not used to create new PRs.",
    "first_pull_request": {
      "number": 5,
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/5",
      "author": "AI-Your-Bro",
      "review_requested_from": "Includeoyasi",
      "verification": "Creation response user.login and requested_reviewers checked"
    },
    "first_commit": {
      "sha": "fbbaef0a94e5d7477f4be5aa6f3f114ce03855f5",
      "github_author": "AI-Your-Bro",
      "github_committer": "AI-Your-Bro",
      "verified": true
    }
  },
  "personal_computer_policy": {
    "authorized_on": "2026-10-15",
    "owner_instruction": "После работы выключать рабочие компоненты; это личный компьютер. Полезные программы и библиотеки можно оставить.",
    "scope": "Windows personal computer",
    "rules": [
      "Stop this task containers and remove its network after checks",
      "Quit Docker Desktop and stop its WSL VM after work",
      "Keep Docker autostart disabled",
      "Keep installed tools, current reusable images and useful caches; remove disposable run-specific images and verify cleanup",
      "Stop only task-owned resources; preserve unrelated workloads"
    ],
    "configured_wsl_budget": {
      "memory": "2GB",
      "processors": 2,
      "swap": "1GB"
    },
    "verified_cleanup": {
      "containers_remaining": 0,
      "wsl_shutdown_exit": 0,
      "images_and_build_cache_preserved": true,
      "desktop_stop_exit": 0,
      "remaining_development_processes": [],
      "installed_tools_preserved": true,
      "docker_autostart": false,
      "docker_run_entry_present": false,
      "observed_at": "2026-10-15T22:49:53.4008625+03:00",
      "reboot_performed": false
    }
  },
  "lab_evolution_notes": {
    "status": "USER_CONFIRMED_DIRECTION_WITH_REVIEWABLE_DETAIL",
    "discussed_on": "2026-10-13",
    "confirmed_direction": [
      "Лаборатория должна иметь границу с тестируемой бизнес-системой для последующего подключения других проектов",
      "Эталонный магазин подключается через заранее настроенный адаптер с бизнес-операциями и проверками, используя ту же границу, что и будущие проекты"
    ],
    "confirmation_evidence": "Владелец проекта подтвердил границу лаборатории, встроенный адаптер магазина, затем попросил детализацию и сохранение обсуждения",
    "proposals_pending_design": [
      "Точные Go-интерфейсы, capabilities, типы исходов, доказательства и формат профиля по первому работающему срезу",
      "Обработка изображений через HTTP как кандидат второго бизнеса; конкретный проект не выбран",
      "Полный эксперимент на втором бизнесе как предлагаемое доказательство переносимости; дополнение к G4 ещё не принято",
      "Мастер подготовки операций, данных и проверок с подтверждением владельцем целевой системы; перспектива после проверенного профиля"
    ],
    "scope_note": "Детализация сохранена для ревью, не является работающей реализацией или одобрением точных интерфейсов. Три бизнес-сервиса магазина и инварианты G0 сохраняются; второй бизнес подключается как другая тестируемая система.",
    "detail_document": "docs/architecture/LAB_BOUNDARY.md"
  },
  "environment_strategy": {
    "status": "VERIFIED_WINDOWS_X64_AND_NATIVE_LINUX_AMD64_ARM64; MACOS_NOT_RUN",
    "approved_on": "2026-10-13",
    "confirmation_evidence": "Владелец проекта описал работу на Windows и macOS, согласовал Docker-подход сообщением Да, давай так и сделаем",
    "documents": [
      "docs/development/ENVIRONMENT.md",
      "docs/development/RESOURCE_PROFILES.md"
    ],
    "primary_hosts": [
      "Windows desktop",
      "macOS laptop"
    ],
    "runtime": "Docker Desktop, Linux containers, Docker Compose",
    "toolchain": "Контейнерная сборка Go, тесты и закреплённые инструменты; host Go необязателен",
    "planned_image_platforms": [
      "linux/amd64",
      "linux/arm64"
    ],
    "container_layout_recommendation": {
      "base": "3 бизнес-сервиса + 1 PostgreSQL с 3 базами + 1 Kafka + 1 Redis",
      "observability": "Prometheus, Grafana, Loki, Alloy — 4 дополнительных компонента",
      "tools": "Временные запуски toolchain/lab; адаптер — модуль без отдельного обязательного контейнера",
      "count_scope": "6 постоянных core, 10 с базовой наблюдаемостью при одной реплике; дополнительные exporters/реплики увеличат состав",
      "measured_resource_cost": "NOT_MEASURED"
    },
    "comparison_rule": "По умолчанию A/B на одной машине с одинаковыми лимитами и нативной архитектурой; разные хосты/ресурсы — отдельный фактор",
    "future_server": {
      "status": "FUTURE_OPTION",
      "user_intent": "Владелец проекта может арендовать сервер позже, если понадобятся дополнительные ресурсы и масштабирование",
      "server_selected": false,
      "rental_or_deployment_authorized": false
    },
    "runtime_created": true,
    "windows_container_start_verified": true,
    "macos_container_start_verified": false,
    "toolchain_status": "VERIFIED_NATIVE_LINUX_CI; RUNTIME_STAGE_BUILT_ON_WINDOWS",
    "runtime_scope": "THREE_HEALTH_SERVICES_WITH_BOUNDED_SHUTDOWN; NO_BUSINESS_OR_INFRA",
    "runtime_containers_created": true,
    "runtime_container_linux_verification": "PASS_NATIVE_AMD64_AND_ARM64",
    "runtime_container_windows_verification": "PASS"
  },
  "environment_facts": {
    "local_host": "Windows",
    "shell": "PowerShell",
    "user_timezone": "Europe/Moscow",
    "project_path_at_snapshot": null,
    "path_required_to_read_context": false,
    "go_available": "HOST_MSI_GO_1_27_2_INSTALLED_MACHINE_PATH_REGISTERED",
    "docker_available": "DESKTOP_4_95_0_INSTALLED_PER_USER_CURRENTLY_STOPPED",
    "tool_commands_checked_on": "2026-10-15",
    "tool_check_scope": "Установлены host Go1.27.2 через официальный MSI и Docker Desktop4.95 per-user/WSL2.7.13; machine PATH Go зарегистрирован, установленный go.exe version проверен. Windows container smoke PASS. Docker Desktop и WSL VM после проверки остановлены. Существующему терминалу может потребоваться переоткрытие для обновления PATH.",
    "powershell_version": "7.6.5",
    "cpu_ram_disk_budget": "HOST_SNAPSHOT_COLLECTED_DOCKER_BUDGET_NOT_MEASURED",
    "local_directory_name": "faultline-lab",
    "cpu_architecture": "X64",
    "wsl_command_available": true,
    "resource_snapshot": {
      "observed_on": "2026-10-13",
      "source": "Windows CIM: CPU/OS; PowerShell volume capacity/free space",
      "os": "Windows 11 Pro, build 26200, x64",
      "cpu_model": "Intel Core i5-12400F",
      "physical_cores": 6,
      "logical_processors": 12,
      "os_visible_memory_gib": 15.72,
      "available_memory_gib_at_observation": 3.38,
      "workspace_volume_capacity_gib": 930.51,
      "workspace_volume_free_gib_at_observation": 134.32,
      "docker_budget": "NOT_MEASURED",
      "application_consumption": "NOT_MEASURED_NO_RUNTIME",
      "macos_resources": "NOT_MEASURED",
      "future_server_resources": "NOT_SELECTED",
      "note": "Свободные ресурсы изменяются; снимок хоста не равен бюджету Docker и не подтверждает capacity приложения"
    },
    "local_go_health_check": "EARLIER_NATIVE_WINDOWS_GO_TCP_UNIT_AND_PROCESS_TESTS_PASS; WINDOWS_CONTAINER_SMOKE_NOW_PASS",
    "wsl_status": "INSTALLED_2_7_13_0_VM_SHUT_DOWN",
    "docker_service_status": "PER_USER_INSTALL_NO_PRIVILEGED_DOCKER_SERVICE_REQUIRED",
    "local_container_smoke": {
      "date": "2026-10-15",
      "result": "PASS",
      "exit_code": 0,
      "source_revision": "d81b51ce2eb46f56d698fd8bf5f9966b91a16036",
      "source_dirty": false,
      "run_id": "8119b50d-f910-412a-9bd2-ab50fccd1b71",
      "evidence": "docs/development/evidence/p01/WINDOWS_SMOKE.md"
    },
    "os_verified_on": "2026-10-15",
    "os_verified": "Windows 11 Pro 10.0.26200 x64",
    "hypervisor_present": true,
    "lanmanserver": "RUNNING_AUTOMATIC",
    "installation_status": "COMPLETE",
    "host_versions": {
      "go": "1.27.2",
      "wsl": "2.7.13.0",
      "docker_desktop": "4.95.0",
      "docker_engine": "29.8.2",
      "compose": "5.6.0",
      "node": "v24.19.0"
    },
    "docker_budget": {
      "wsl_config": {
        "memory": "2GB",
        "processors": 2,
        "swap": "1GB"
      },
      "engine_cpu_count": 2,
      "engine_memory_bytes": 1996603392,
      "performance_capacity": "NOT_MEASURED"
    },
    "post_task_cleanup": {
      "containers_remaining": 0,
      "wsl_shutdown_exit": 0,
      "images_and_build_cache_preserved": true,
      "desktop_stop_exit": 0,
      "remaining_development_processes": [],
      "installed_tools_preserved": true,
      "docker_autostart": false,
      "docker_run_entry_present": false,
      "observed_at": "2026-10-15T22:49:53.4008625+03:00",
      "reboot_performed": false
    }
  },
  "toolchain": {
    "document": "docs/development/TOOLCHAIN.md",
    "go": "1.27.2",
    "goimports": "golang.org/x/tools v0.52.0",
    "govulncheck": "golang.org/x/vuln v1.9.0",
    "golangci_lint": "v2.14.0",
    "go_image": "golang:1.27.2-trixie@sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e",
    "lint_image": "golangci/golangci-lint:v2.14.0@sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f",
    "image_indexes": "VERIFIED_LINUX_AMD64_AND_ARM64",
    "application_dependencies": "google.golang.org/grpc v1.84.0; STANDARD_LOG_SLOG; TOOL_DIRECTIVES_PINNED",
    "native_ci_runners": [
      "ubuntu-24.04",
      "ubuntu-24.04-arm"
    ],
    "bootstrap_behavior": "THREE_SERVERS_WITH_BOUNDED_GRACEFUL_SHUTDOWN; HEALTHCHECK_SUBCOMMAND",
    "smoke_command": "docker compose run --rm toolchain service-smoke",
    "service_runtime_document": "docs/development/SERVICE_RUNTIME.md",
    "runtime_image": "scratch; static Go executable; USER 65532:65532; no shell/compiler/CA bundle",
    "container_runtime_document": "docs/development/CONTAINER_RUNTIME.md",
    "upload_artifact": {
      "tag": "v7.0.2",
      "sha": "cf430e030ddbb5b0abf93d22962f4752f3646cd9",
      "source": "Official GitHub release and git tag ref checked 2026-10-15"
    }
  },
  "network_configuration_review": {
    "requested_on": "2026-10-15",
    "pull_request_url": "https://github.com/Includeoyasi/faultline-lab/pull/7",
    "source": "Direct project owner review in chat",
    "request": "Мне не нравится, что в коде присутствуют локальные IP адреса типа 127.0.0.1:50053. Давай сразу оформим красиво и правильно (с точки зрения бизнес разработки и безопасности).",
    "status": "ACCEPTED_AS_PART_OF_OWNER_MERGED_PR7",
    "decisions": [
      "FL_GRPC_ADDR is mandatory: no IP/port literals or implicit network fallback in application Go source. Missing/empty address fails before bind with sanitized JSON and exit 1.",
      "Local and container profiles are explicit nonsecret configs/env/*.env.example files; application itself does not load dotenv.",
      "Separate listener configuration, future peer DNS targets and host publication. Container wildcard is explicit, not an automatic fallback or authentication policy.",
      "Default toolchain has no published ports. Manual development instructions publish only on host loopback; unqualified host publication removed.",
      "Personal *.env files are ignored by Git and excluded from Docker build context, including nested dotenv files.",
      "Loopback in isolated test fixtures and observed logs is intentional and not a deployment default. TLS/mTLS, runtime images and client service discovery remain future work."
    ],
    "local_checks": "PASS: build/vet/unit/real process TCP smoke; three missing and three empty address process cases; Git ignore verified; runtime source IP/port search returned no matches",
    "native_ci": {
      "name": "Go toolchain",
      "head_sha": "2922b4a00b80d4e83dd63767111c8ee2d60a179c",
      "status": "completed",
      "conclusion": "success",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38737981164",
      "environments": [
        "ubuntu-24.04 amd64",
        "ubuntu-24.04-arm arm64"
      ]
    },
    "owner_acceptance": "ACCEPTED_PR7_MERGE",
    "documentation_ci": {
      "name": "Documentation",
      "head_sha": "2922b4a00b80d4e83dd63767111c8ee2d60a179c",
      "status": "completed",
      "conclusion": "success",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38737981056"
    }
  },
  "review_communication": {
    "updated_on": "2026-10-15",
    "authority": "Latest direct owner instruction",
    "rule": "Use clear descriptive headings: Что было сделано / Проверки / На что обратить внимание, or another plain heading that describes the content. Do not use the old question-style guide title.",
    "applies_to": "All new PR descriptions and comments; old requests in change history are superseded."
  },
  "branding_notes": {
    "status": "USER_SELECTED_MERGED_TO_MAIN",
    "selected_on": "2026-10-13",
    "brief": "Три варианта: технологичность, современность, доступность, мягкие пастельные тона, исследование и будущее",
    "document": "docs/design/BANNER_OPTIONS.md",
    "alternatives": [
      "01-system-map",
      "02-experiment-flow",
      "03-research-horizon"
    ],
    "selected": "02-experiment-flow",
    "selection_evidence": "Прямое сообщение владельца проекта: Мне нравится вариант 2",
    "primary_asset": "docs/assets/banner.svg",
    "source_asset": "docs/assets/banner-concepts/02-experiment-flow.svg",
    "preview": "docs/assets/banner-concepts/02-experiment-flow.png",
    "comparison": "docs/assets/banner-concepts/comparison.png",
    "svg_sha256": "7eb3a6e1ec1176769be56d05c959dd75893d17e7dd27fea0d9767701d6a7d119",
    "published": true,
    "note": "Линии декоративные, не изображают измерения; оба README используют общий SVG",
    "published_branch": "docs/portable-lab-plan",
    "published_to_main": true,
    "pull_request_url": "https://github.com/Includeoyasi/faultline-lab/pull/4",
    "merge_commit": "8c392b8ed4145d8b2b9efa93c5be14d9fdf113d1"
  },
  "validation_evidence": [
    {
      "date": "2026-10-13",
      "check": "Inventory of the project files and Git repository status",
      "result": "PASS",
      "detail": "До создания контекста обнаружены только шесть проектных Markdown-файлов. Git rev-parse подтвердил отсутствие репозитория в этом пути и родителях."
    },
    {
      "date": "2026-10-11",
      "check": "Markdown structure and encoding for G0 visuals",
      "result": "PASS",
      "detail": "Проверены кодировка, парные fences и типы Mermaid: 15 схем в docs/architecture/PROCESSES.md, 3 в docs/architecture/DESIGN.md, 1 в docs/ROADMAP.md. Это не проверка рендеринга Mermaid и не тесты сервисов."
    },
    {
      "date": "2026-10-13",
      "check": "Generate portable context and check source freshness",
      "result": "PASS",
      "commands": [
        ".\\tools\\update_context.ps1",
        ".\\tools\\update_context.ps1 -CheckOnly"
      ],
      "detail": "Собраны 7 полных проектных документов и реестр 10 SHA-256 источников. CheckOnly подтвердил соответствие переносимого снимка исходным файлам."
    },
    {
      "date": "2026-10-13",
      "check": "Portable context structure and approval consistency",
      "result": "PASS",
      "detail": "Встроены 19 Mermaid-схем; fences сбалансированы; все 7 внутренних anchors существуют; UTF-8 без replacement characters. Устаревшие утверждения, что принятие G0 ещё не получено, в документах не обнаружены. Рендеринг схем и бизнес-тесты не запускались."
    },
    {
      "date": "2026-10-13",
      "check": "GitHub repository creation and local Git bootstrap",
      "result": "PASS",
      "detail": "GitHub API подтвердил владельца Includeoyasi, публичность, URL и создание репозитория. Git main и origin созданы. Авторство настроено локально с GitHub noreply email. Коммит и push ещё не выполнены."
    },
    {
      "date": "2026-10-13",
      "check": "Repository documentation and portable context checks",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "powershell.exe -NoProfile -File tools/update_context.ps1 -CheckOnly"
      ],
      "detail": "Проверены 15 Markdown-файлов, 115 локальных ссылок и 50 Mermaid-блоков вместе с повторениями внутри снимка. Снимок включает 12 документов и 15 SHA-256. CheckOnly проходит в PowerShell 7 и Windows PowerShell 5.1. SVG-обложка отрендерена и визуально проверена. Рендеринг всех Mermaid-схем и тесты сервисов не выполнялись."
    },
    {
      "date": "2026-10-13",
      "check": "Initial commit publication and GitHub Actions",
      "result": "PASS",
      "commands": [
        "git commit",
        "git push -u origin main",
        "GitHub API: workflow runs and job steps"
      ],
      "detail": "Git push создал удалённый main и tracking origin/main. Documentation run 38317062261 завершён success на коммите 6897f4a; оба шага проверки документов и контекста проходят на Ubuntu. Публичный README, MIT, темы и файловая структура доступны в GitHub.",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38317062261"
    },
    {
      "date": "2026-10-13",
      "check": "Go review policy documentation and portable snapshot",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "powershell.exe -NoProfile -File tools/update_context.ps1 -CheckOnly",
        "git diff --check"
      ],
      "detail": "Проверены 16 Markdown-файлов, 144 локальные ссылки, fences/кодировка и отсутствие машинных путей. Снимок содержит 13 документов и 16 fingerprints; актуальность проходит в PowerShell 7 и Windows PowerShell 5.1. Go-проверки NOT_RUN: модуль и сервисы ещё отсутствуют. Человеческая приёмка правил не заявляется."
    },
    {
      "date": "2026-10-13",
      "check": "Publish Go review policy branch and PR",
      "result": "PASS",
      "detail": "Ветка docs/go-review-policy опубликована, GitHub PR #1 открыт в main с HEAD 214eb34. Documentation push run 38318256118 завершился success. PR не одобрен ИИ и не слит; решение владельца проекта не получено.",
      "pull_request_url": "https://github.com/Includeoyasi/faultline-lab/pull/1",
      "ci_url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38318256118"
    },
    {
      "date": "2026-10-13",
      "check": "Observe project owner merge and reconcile current main",
      "result": "PASS",
      "detail": "GitHub API подтвердил PR #1 merged=true, merged_by=Includeoyasi, accepted HEAD 214eb34, merge c1e55fd. Main CI run 38318650575 success. Origin/main получен; последующее обновление контекста перенесено на актуальный main без конфликтов. ИИ не выполнял слияние.",
      "pull_request_url": "https://github.com/Includeoyasi/faultline-lab/pull/1",
      "ci_url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38318650575"
    },
    {
      "date": "2026-10-13",
      "check": "Verify the merged context handoff before preparing local terminology changes",
      "result": "PASS",
      "detail": "GitHub metadata подтвердили слияние PR #2 с HEAD 681660a и merge commit 1a0b3f3. Git fetch получил origin/main; новые правки подготовлены на этой основе. Слияние не выполнялось в этой задаче.",
      "pull_request_url": "https://github.com/Includeoyasi/faultline-lab/pull/2"
    },
    {
      "date": "2026-10-13",
      "check": "Project owner terminology and refreshed portable context",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "powershell.exe -NoProfile -ExecutionPolicy Bypass -File tools/update_context.ps1 -CheckOnly",
        "git diff --check",
        "Search canonical documents, templates and generated context for obsolete role terminology"
      ],
      "detail": "Проверены 16 Markdown-файлов, 144 локальные ссылки и структура 50 Mermaid-блоков; снимок ревизии 8 содержит 13 документов и 16 fingerprints. CheckOnly прошёл в PowerShell 7 и Windows PowerShell 5.1. Прежнее название роли и дублированные фразы не обнаружены. Коммит и push ещё не выполнены; CI нового изменения и Go-проверки не запускались."
    },
    {
      "date": "2026-10-13",
      "check": "Initial local tool availability check before G1-P01",
      "result": "OBSERVED",
      "detail": "Get-Command не обнаружил go и docker в текущем PATH; PowerShell 7.6.5 доступен. Go/Docker версии и Docker daemon не проверены, поскольку команды не найдены. Ничего не устанавливалось; ресурсы и другие расположения инструментов ещё не исследованы."
    },
    {
      "date": "2026-10-13",
      "check": "Read Windows host resource snapshot",
      "result": "PASS",
      "detail": "CPU/OS CIM чтение сначала ограничено песочницей; разрешённый повтор дал CPU 6/12, 15.72 GiB видимой памяти, 3.38 GiB свободной и 134.32 GiB свободного места. Нулевые значения неудачного чтения не используются. Лимиты Docker и потребление приложения не измерены."
    },
    {
      "date": "2026-10-13",
      "check": "Render and visually inspect three banner concepts",
      "result": "PASS",
      "detail": "Три SVG отрендерены в PNG через локальный sharp, просмотрен общий лист. Текст не обрезан; владелец выбрал вариант 2. Исходник применён к общему баннеру README. Публикация и просмотр изменённого GitHub README ещё не выполнялись."
    },
    {
      "date": "2026-10-13",
      "check": "Validate saved lab/environment decisions, chosen banner and portable context revision 10",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "powershell.exe -NoProfile -ExecutionPolicy Bypass -File tools/update_context.ps1 -CheckOnly",
        "git diff --check",
        "git diff --cached --check",
        "Compare selected SVG with common README banner; render SVG previews"
      ],
      "detail": "Проверены 20 Markdown-файлов, 194 локальные ссылки, UTF-8/LF/fences и структура 56 Mermaid-блоков. Контекст ревизии 10 содержит 17 встроенных документов и 20 fingerprints; CheckOnly проходит в PowerShell 7 и Windows PowerShell 5.1. PNG previews созданы, баннер совпадает с выбранным SVG. Полный Mermaid-render и сервисные проверки не выполнялись; runtime отсутствует. Изменения не закоммичены и не опубликованы."
    },
    {
      "date": "2026-10-13",
      "check": "Inspect existing task tracking",
      "result": "OBSERVED",
      "detail": "GitHub search repo:Includeoyasi/faultline-lab is:issue вернул пустой список. До этой задачи G1-P01 имел общий объём в design/context без явной декомпозиции и отдельных критериев."
    },
    {
      "date": "2026-10-13",
      "check": "G1-P01 decomposition and tracking proposal documentation",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "git diff --check"
      ],
      "detail": "Проверены 22 Markdown-файла, 230 локальных ссылок и структура 56 Mermaid-блоков. Снимок ревизии 11 содержит 19 документов и 22 fingerprints. Критерии описаны, но не исполнены: код 0/8 и Go/Docker проверки NOT_RUN. Внешние Issues/Projects не создавались."
    },
    {
      "date": "2026-10-13",
      "check": "Separate gate and packet naming; first code PR P01 between G0 and G1",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "git diff --check",
        "Search active documents for obsolete packet labels"
      ],
      "detail": "Активные документы, пути и ID задач используют P01–P05. Старое обозначение сохранено только в истории/таблице соответствия; восемь задач и критерии не изменены. Проверены 22 Markdown-файла, 230 локальных ссылок и актуальность контекста ревизии 12 с 19 документами и 22 fingerprints. Код 0/8, runtime/Go проверки не выполнялись."
    },
    {
      "date": "2026-10-13",
      "check": "GitHub P01 issue creation and body completeness",
      "result": "PASS",
      "detail": "GitHub REST вернул открытый Issue #3; в сохранённом ответе проверены 8 пустых task checkboxes и 8 разделов критериев, текст не усечён. Это создание backlog, а не выполнение задач.",
      "url": "https://github.com/Includeoyasi/faultline-lab/issues/3"
    },
    {
      "date": "2026-10-13",
      "check": "Documentation and context for created P01 backlog",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "git diff --check",
        "git diff --cached --check"
      ],
      "detail": "Проверены 22 Markdown-файла, 234 локальные ссылки, структура 56 Mermaid-блоков и актуальность ревизии 13 с 19 документами/22 fingerprints. Активные документы не содержат старых имён пакетов или роли тимлида. Проверки Go/runtime не выполнялись: реализации нет."
    },
    {
      "date": "2026-10-13",
      "check": "Initial documentation PR CI",
      "result": "PASS",
      "head_sha": "2c28d127764210f1019c114a4b3d5b1cd59b529c",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38402553170",
      "detail": "GitHub Actions Documentation pull_request run завершён успешно на 2c28d12. Результат не распространяется автоматически на последующий handoff commit и не является тестом Go/runtime."
    },
    {
      "date": "2026-10-13",
      "check": "Published PR and issue handoff documentation",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "git diff --check"
      ],
      "detail": "Проверены 22 Markdown-файла, 234 локальные ссылки, структура 56 Mermaid-блоков и актуальность ревизии 14 с 19 документами/22 fingerprints. Снимок содержит созданные Issue #3 и PR #4; code/runtime checks NOT_RUN. Наблюдённый CI указан с точным SHA; новый handoff HEAD проверяется отдельно в PR."
    },
    {
      "date": "2026-10-13",
      "check": "GitHub main protection and executor profile",
      "result": "PASS",
      "detail": "Независимый GET подтвердил administrator enforcement, 1 approval, stale dismissal, last-push approval, CODEOWNERS flag, strict CI с GitHub Actions app 15368, conversation resolution, запрет force push/deletion. branches/main protected=true. Профиль AI-Your-Bro ID 342030453 существует. Приглашение write ID 337527166 создано. Это не проверка принятия приглашения, OAuth или 2FA."
    },
    {
      "date": "2026-10-13",
      "check": "Main protection and executor onboarding documentation",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "git diff --check"
      ],
      "detail": "23 Markdown-файла, 244 локальные ссылки, структура 56 Mermaid-блоков, 21 встроенный документ и 24 fingerprints ревизии 15. CODEOWNERS проверен как локально подготовленный; публикация/авторизация исполнителя не заявляются. Go/runtime tests NOT_RUN."
    },
    {
      "date": "2026-10-13",
      "check": "Executor authenticated identity and least privilege",
      "result": "PASS",
      "detail": "GCM login завершён с exit 0. Авторизованный GET /user: AI-Your-Bro ID 342030453. GET collaborators/AI-Your-Bro/permission от исполнителя: write, push=true, admin=false. Локальные user.name, user.email и credential username прочитаны обратно. 2FA не заявляется включённой."
    },
    {
      "date": "2026-10-13",
      "check": "Authenticated executor and protected-main documentation",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "git diff --check",
        "git diff --cached --check"
      ],
      "detail": "23 Markdown-файла, 244 локальные ссылки, структура 56 Mermaid-блоков; ревизия 16, 21 встроенный документ и 24 fingerprints. CODEOWNERS назначает все файлы Includeoyasi; авторизация исполнителя и write без admin наблюдены. 2FA не включена этим заданием по решению владельца. Go/runtime tests не запускались."
    },
    {
      "date": "2026-10-13",
      "check": "Executor-authored security PR publication",
      "result": "PASS",
      "head_sha": "fbbaef0a94e5d7477f4be5aa6f3f114ce03855f5",
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/5",
      "detail": "PR #5 user.login=AI-Your-Bro, requested_reviewers содержит Includeoyasi. GitHub commit author и committer AI-Your-Bro. Исходный Documentation pull_request run SUCCESS. ИИ не approve/merge; проверка финального handoff HEAD отдельно в PR."
    },
    {
      "date": "2026-10-13",
      "check": "Published executor PR portable handoff",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly",
        "git diff --check"
      ],
      "detail": "23 Markdown-файла, 244 локальные ссылки, структура 56 Mermaid-блоков; ревизия 17, 21 встроенный документ и 24 fingerprints. Публикация PR #5 и авторство исполнителя отражены, CODEOWNERS в main не заявлен до merge, решение владельца о 2FA сохранено. Текущий HEAD/CI проверяются в PR; Go/runtime N/A."
    },
    {
      "date": "2026-10-15",
      "check": "Local P01-01 Go checks",
      "result": "PASS",
      "commands": [
        "go version",
        "go mod tidy -diff",
        "go mod verify",
        "go build ./...",
        "go test -count=1 ./...",
        "go tool goimports -l cmd",
        "go tool govulncheck ./..."
      ],
      "environment": "Temporary checksum-verified Go 1.27.2 windows/amd64; caches in ignored .local; no system Go installation",
      "detail": "Сборка трёх пакетов проходит; tidy не требует изменений, verify PASS, goimports не вывел файлов, govulncheck: No vulnerabilities found. Go tests: no test files. Docker, Linux race/lint и bootstrap execution ожидают CI; macOS NOT_RUN."
    },
    {
      "date": "2026-10-15",
      "check": "Initial Linux ARM64 container check",
      "result": "FAIL",
      "head_sha": "c1e7c7699660fedfb997d5bcaac38514b8723853",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38722835503",
      "detail": "Образ собрался, версии инструментов/format/tidy/verify прошли. go build не смог получить Git VCS status из-за ownership bind mount. В образ добавлен точный safe.directory /workspace; Git metadata сохраняется. Повторный CI ожидается."
    },
    {
      "date": "2026-10-15",
      "check": "Portable source embedding and initial documentation",
      "result": "PASS",
      "detail": "24 Markdown-файла, 275 локальных ссылок, структура 56 Mermaid; ревизия 18, 32 документа и 35 fingerprints. CheckOnly прошёл PS7 и Windows PS5.1. Documentation CI исходного коммита успешен."
    },
    {
      "date": "2026-10-15",
      "check": "Linux AMD64/ARM64 Go checks and build-output permissions",
      "result": "PARTIAL_PASS",
      "head_sha": "1cc5d8c63a7263d677f277c989e6e2fcb051020a",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38723121526",
      "detail": "На обеих архитектурах прошли версии, format/imports, tidy/verify, build, golangci-lint (0 issues), unit/race commands (no test files), govulncheck (no vulnerabilities). Bootstrap build не смог создать bin: cap_drop ALL лишил root возможности писать в bind mount runner с другим UID. Для временной среды сборки восстановлены стандартные Docker capabilities; privileged mode не включён, no-new-privileges сохранён. Повторный CI ожидается."
    },
    {
      "date": "2026-10-15",
      "check": "P01-01 native Linux container checks",
      "result": "PASS",
      "head_sha": "7ea4b89c1ce5c916da1d4293501a3f8f24b66c63",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38723552125",
      "commands": [
        "docker compose build toolchain",
        "docker compose run --rm toolchain check",
        "docker compose run --rm toolchain bootstrap-smoke",
        "invalid-command exits exactly 2",
        "git diff --exit-code"
      ],
      "environments": [
        "ubuntu-24.04 amd64",
        "ubuntu-24.04-arm arm64"
      ],
      "detail": "Оба нативных job SUCCESS. Построен pinned Dockerfile, Go 1.27.2, линтер 2.14.0 (0 issues), format/imports/tidy/verify/build, unit/race commands (no test files), govulncheck (no vulnerabilities). Все три бинарника собраны и выполнены. VCS ownership и вывод bin исправлены и проверены. Windows Docker/macOS NOT_RUN; бизнес-проверок нет."
    },
    {
      "date": "2026-10-15",
      "check": "Required Go checks activated",
      "result": "PASS",
      "detail": "После успешного Go CI PATCH добавил Go toolchain (amd64)/(arm64), GitHub Actions app 15368. Независимый GET подтвердил три required checks со strict=true и сохранённые owner review/CODEOWNERS/admin enforcement, запреты force push/deletion."
    },
    {
      "date": "2026-10-15",
      "check": "P01-01 final portable review packet",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly (PowerShell 7)",
        "tools/update_context.ps1 -CheckOnly (Windows PowerShell 5.1)",
        "git diff --check"
      ],
      "detail": "24 Markdown-файла, 276 локальных ссылок, структура 56 Mermaid-блоков; ревизия 21, 33 встроенных источника и 36 fingerprints. Исходники P01-01 включены с language-tagged fences. Проверено 1/8, принято владельцем 0/8; оставшиеся семь задач не начаты."
    },
    {
      "date": "2026-10-15",
      "check": "PR #6 owner merge",
      "result": "PASS",
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/6",
      "head_sha": "001d333132613f642b9f2415865bcd1255cb397a",
      "merge_commit": "e4a9b0b80f97b2fb20ca973ef57030d32f8ba878",
      "detail": "Владелец сообщил о слиянии; API merged=true. ИИ не approve/merge."
    },
    {
      "date": "2026-10-15",
      "check": "P01-02–04 local Windows native Go",
      "result": "PASS",
      "commands": [
        "go mod tidy -diff",
        "go mod verify",
        "go build ./...",
        "go vet ./...",
        "go test -count=1 -v -timeout=3m ./...",
        "go tool govulncheck ./..."
      ],
      "detail": "Go 1.27.2 portable, gRPC v1.84.0. Unit tests: defaults/overrides/invalid inputs, safe JSON/levels/build data, bind errors, context cancellation and real TCP health; process tests: all 3 health endpoints, error exit 1, isolation. Sandbox initially blocked loopback; rerun with permission passed. Windows child processes use Kill, not Unix signals; graceful signal exit remains Linux CI. No vulnerabilities found. Local Docker/macOS/race NOT_RUN; Linux container CI pending."
    },
    {
      "date": "2026-10-15",
      "check": "P01-02–04 first native Linux CI",
      "result": "FAIL_FIXED_AWAITING_RERUN",
      "head_sha": "e8d9102acca633883d3b01ae3488b113111c0c92",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38734508988",
      "detail": "Оба image build и module verify прошли; staticcheck QF1001 потребовал упростить инвертированную проверку символа DNS. Переписано в явный switch допустимых символов, без подавления анализатора. Tests/race/smoke этой попытки ещё не выполнялись. Documentation CI PASS."
    },
    {
      "date": "2026-10-15",
      "check": "P01-02–04 native Linux container checks",
      "result": "PASS",
      "head_sha": "4d3b2356b472f275063efb576d8226de4a35f594",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38734815221",
      "documentation_ci": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38734815235",
      "commands": [
        "docker compose build toolchain",
        "docker compose run --rm toolchain check",
        "docker compose run --rm toolchain service-smoke",
        "invalid-command exits exactly 2",
        "git diff --exit-code"
      ],
      "environments": [
        "ubuntu-24.04 amd64",
        "ubuntu-24.04-arm arm64"
      ],
      "detail": "Both jobs SUCCESS. Go 1.27.2; gRPC 1.84.0; lint 0 issues; unit and process tests PASS; -race PASS, including in-process server paths; govulncheck No vulnerabilities found. All three binaries SERVING, unknown health NotFound, invalid config/bind exit 1 with JSON, orders SIGINT exit 0 and final JSON, other two stay SERVING. Immediate stop only; no drain deadline or runtime images. Child smoke binaries are ordinary builds; race instruments the server unit test. Windows Docker/macOS/capacity NOT_RUN."
    },
    {
      "date": "2026-10-15",
      "check": "PR #7 review handoff and tracker",
      "result": "PASS",
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/7",
      "detail": "Executor GraphQL confirmed draft=false and author AI-Your-Bro; review requested from Includeoyasi. Issue PATCH verified 4 completed implementation checkboxes; 1 owner-accepted task, 3 in review, 4 not started. Issue remains open. AI did not approve or merge."
    },
    {
      "date": "2026-10-15",
      "check": "PR #7 network configuration review fix, local Windows",
      "result": "PASS",
      "commands": [
        "go build ./...",
        "go vet ./...",
        "go test -count=1 -timeout=3m ./...",
        "git check-ignore configs/env/local.env configs/env/container.env",
        "rg runtime Go sources excluding *_test.go: no literal loopback/wildcard IP or service ports"
      ],
      "detail": "FL_GRPC_ADDR required; missing/empty address rejected for each real binary. Unit defaults are only log level and timeout; valid explicit listen addresses still work. All three real health endpoints and process isolation pass. Linux CI for this change is pending; previous green CI is historical evidence only."
    },
    {
      "date": "2026-10-15",
      "check": "PR #7 explicit network configuration, native CI",
      "result": "PASS",
      "head_sha": "2922b4a00b80d4e83dd63767111c8ee2d60a179c",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38737981164",
      "documentation_ci": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38737981056",
      "commands": [
        "docker compose build toolchain",
        "docker compose run --rm toolchain check",
        "docker compose run --rm toolchain service-smoke",
        "invalid-command exits exactly 2",
        "git diff --exit-code"
      ],
      "environments": [
        "ubuntu-24.04 amd64",
        "ubuntu-24.04-arm arm64"
      ],
      "detail": "Both jobs SUCCESS. Format, tidy/verify/build, lint 0 issues, unit/race tests and govulncheck passed. Three real binaries reject missing or empty FL_GRPC_ADDR with one sanitized JSON record and exit 1; explicit addresses still SERVING. Bind/config failures, independent processes and Linux signal exit pass. No environment-specific addresses/ports in application Go sources; test fixtures use loopback deliberately. Runtime images, TLS/mTLS and local Docker Desktop/macOS remain unimplemented/unverified as documented."
    },
    {
      "date": "2026-10-15",
      "check": "PR #7 owner acceptance and merge",
      "result": "PASS",
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/7",
      "head_sha": "2e8c7f8149f27c8fd9820401bd1e5d7e476ca0c9",
      "merge_commit": "bc09d08ce7867f3109d996f71015e3ce73986452",
      "detail": "Owner accepted in chat; API merged=true. AI did not approve/merge."
    },
    {
      "date": "2026-10-15",
      "check": "P01-05–08 local native Go preparation",
      "result": "PASS",
      "commands": [
        "go mod tidy -diff",
        "go build ./...",
        "go vet ./...",
        "go test -count=1 -timeout=3m ./...",
        "node --check tools/smoke.mjs"
      ],
      "detail": "Portable Go 1.27.2 on Windows. Controlled real RPC tests verify NOT_SERVING before listener closure, active call completion and bounded forced return for a handler ignoring cancellation; fixtures release/join their work afterward. Health probe and process tests pass. Unix process signals and real runtime containers await Linux CI."
    },
    {
      "date": "2026-10-15",
      "check": "Windows runtime container smoke prerequisites",
      "result": "NOT_RUN",
      "command": "node tools/smoke.mjs",
      "exit_code": 2,
      "detail": "Docker executable ENOENT; report.json written. Standard Docker Desktop path and service absent; wsl --list --verbose reports subsystem not installed. Local host install/configuration not changed. Criterion P01-06 remains unmet; draft required."
    },
    {
      "date": "2026-10-15",
      "check": "First P01-05–08 container CI",
      "result": "FAIL_FIXED_AWAITING_RERUN",
      "head_sha": "afdce9351ed7d9616954526f3e8ac444bdaef9b4",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38744565647",
      "detail": "Both architectures failed Compose model validation before build: inactive profile service toolchain was referenced as an additional context. Replaced cross-service build dependency with one multi-stage build/Dockerfile and explicit toolchain/runtime targets. Checks remain enabled; rerun required."
    },
    {
      "date": "2026-10-15",
      "check": "P01-05–08 native Linux Go and runtime CI",
      "result": "PASS",
      "head_sha": "d5e13a9ece471c6ab8410a53a973b96ad9345aa1",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38744910340",
      "documentation_ci": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38744910192",
      "environments": [
        "ubuntu-24.04 amd64",
        "ubuntu-24.04-arm arm64"
      ],
      "commands": [
        "docker compose build toolchain",
        "docker compose run --rm toolchain check",
        "docker compose run --rm toolchain service-smoke",
        "node tools/smoke.mjs",
        "git diff --exit-code"
      ],
      "detail": "Both jobs SUCCESS. Format/tidy/verify/build/lint (0 issues), unit/race, govulncheck (No vulnerabilities found). Controlled RPC drain/deadline checks pass. All 3 real processes verify SIGTERM/SIGINT normal exit and held Watch forced exit 1. Both architectures build scratch runtime images, start 3 containers, verify health/security/internal network/no host ports or mounts, stop orders without stopping peers, verify all normal exits/log order and cleanup. Artifacts retained 14 days and independently read; image IDs and runtime outcomes recorded in runtime_smoke_evidence. Windows/macOS Docker are NOT_RUN; no capacity claim."
    },
    {
      "date": "2026-10-15",
      "check": "P01 final implementation CI and retained reports",
      "result": "PASS_LINUX; WINDOWS_CONTAINER_NOT_RUN",
      "head_sha": "59684070fa85a0aea8561b48828088e574a06436",
      "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38747476656",
      "documentation_ci": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38747476687",
      "detail": "Both native architecture jobs passed all Go/process/runtime checks. Both retained JSON/JSONL artifacts independently read, merge checkout cf9462c321e939279bd926f367ef4fde3015ef71, clean tree, all 3 containers healthy, unprivileged/readonly/internal/no host ports or mounts, normal exit0, peer isolation and cleanup PASS. No source edits after this proof; later snapshot publication requires its own HEAD checks."
    },
    {
      "date": "2026-10-15",
      "check": "Windows Docker Desktop/WSL2 runtime smoke and shutdown",
      "result": "PASS",
      "head_sha": "d81b51ce2eb46f56d698fd8bf5f9966b91a16036",
      "command": "node tools/smoke.mjs",
      "run_id": "8119b50d-f910-412a-9bd2-ab50fccd1b71",
      "source_dirty": false,
      "evidence": "docs/development/evidence/p01/WINDOWS_SMOKE.md",
      "detail": "All 3 real Linux containers healthy/security constraints verified; orders isolation PASS; all graceful exit0; scoped cleanup PASS. Go MSI 1.27.2 installed/version verified; WSL2.7.13.0; Desktop4.95.0; Engine29.8.2; Compose5.6.0; Node24.19.0. Desktop and WSL VM stopped afterward, no Docker/vmmemWSL processes, autostart disabled. No reboot."
    },
    {
      "date": "2026-10-15",
      "check": "Engineering decision records and separate Go smoke backlog",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly (PowerShell 7)",
        "tools/update_context.ps1 -CheckOnly (Windows PowerShell 5.1)",
        "git diff --check"
      ],
      "detail": "31 Markdown files, 437 local links, 62 Mermaid blocks checked structurally; context revision36, 60 embedded documents, 63 fingerprints. Executor GET confirms Issue9 OPEN/backlog/enhancement, AI-Your-Bro author, unassigned. No implementation or local runtime started; P01 owner acceptance4/8, G1 NOT_READY. Final publication HEAD CI is checked in PR8, not inferred from baseline CI."
    },
    {
      "date": "2026-10-16",
      "check": "P01 owner merge and next-packet tracker",
      "result": "PASS",
      "url": "https://github.com/Includeoyasi/faultline-lab/pull/8",
      "head_sha": "56f367278a0fabdd621425b38e13fcec323fb5c9",
      "merge_commit": "3ca290d818a98fc986fe42ae51dcfc0546673fb8",
      "detail": "GitHub GET merged=true, merged_by Includeoyasi, merged_at UTC2026-10-15T21:10:58Z = 2026-10-16 Europe/Moscow. Executor PATCH closed Issue3 after all criteria/owner merge; Issue10 created BACKLOG unassigned with 8 unchecked tasks. No code/runtime checks or implementation performed while planning."
    },
    {
      "date": "2026-10-16",
      "check": "P02 planning documentation and portable context",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly (PowerShell7)",
        "tools/update_context.ps1 -CheckOnly (Windows PowerShell5.1)",
        "git diff --check"
      ],
      "detail": "32 Markdown files; 464 local links; 62 Mermaid blocks structure only. Revision37 embeds61 documents with64 source fingerprints. Changes only documentation/state; P02 and Go smoke9 NOT_STARTED, priority/PR grouping PROPOSED. No Go/service/container tests or runtime startup needed for the planning change; automatic final HEAD CI tracked in draftPR11."
    },
    {
      "date": "2026-10-16",
      "check": "Approved P02 task publication and portable context",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly (PowerShell7)",
        "tools/update_context.ps1 -CheckOnly (Windows PowerShell5.1)",
        "git diff --check"
      ],
      "detail": "Executor create responses verify milestoneP02 and eight open Issues12–19, authorAI-Your-Bro, unassigned, all criteria unchecked. Local documentation/context checks passed; no application/runtime source changes and no local runtime started. Dependency links and final branch publication verified separately after publication."
    },
    {
      "date": "2026-10-16",
      "check": "P02 consolidation and documentation",
      "result": "PASS",
      "commands": [
        "node tools/check_docs.mjs",
        "tools/update_context.ps1 -CheckOnly (PowerShell7)",
        "tools/update_context.ps1 -CheckOnly (Windows PowerShell5.1)",
        "git diff --check",
        "Executor GET milestone issue collection and PATCH response verification"
      ],
      "detail": "Milestone-filtered issue collection contains only open12/16; titles/bodies match canonical descriptions and criteria, all checkboxes unchecked and no assignees. Six replaced cards are closed not_planned with milestone=null. Local docs/context checks passed. Milestone cached aggregate counters are checked separately; no development/runtime startup or PR creation."
    },
    {
      "date": "2026-10-16",
      "check": "External auditv2 integrity and source/evidence review",
      "result": "PASS_READ_ONLY_REVIEW",
      "detail": "144 manifest SHA256 entries and79 Git source blobs match. Current executable sources/configs unchanged. Saved test JSON independently recounts15 top-level+59 subtest pass events; coverage116/128 internal and116/140 total. Main3ca and openIssues9/12/16 checked via read-only API. No fresh runtime/Go tests and no execution of attached scripts; not a new PASS for their behavior."
    },
    {
      "date": "2026-10-16",
      "check": "Audit remediation issue publication",
      "result": "PUBLISHED_AND_PAYLOAD_MATCHED",
      "detail": "Executor helper verified fresh /user actor AI-Your-Bro before every API request. Issue20 matches reviewed Markdown with12 F sections and25 unchecked acceptance criteria. Existing9/16 bodies updated and matched while states/milestones preserved. No Go/runtime checks or audit scripts executed."
    }
  ],
  "runtime_smoke_evidence": {
    "source_head": "59684070fa85a0aea8561b48828088e574a06436",
    "workflow": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38747476656",
    "artifacts": [
      {
        "architecture": "amd64",
        "artifact_id": 11979987375,
        "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38747476656/artifacts/11979987375",
        "retention_days": 14,
        "result": "PASS",
        "built_merge_revision": "cf9462c321e939279bd926f367ef4fde3015ef71",
        "source_dirty": false,
        "docker_platform": "linux/x86_64",
        "cleanup": "PASS",
        "isolation_after_orders_stop": "PASS",
        "services": {
          "orders": {
            "image_id": "sha256:618542af6f61d5dafabc413a5ff5f814c6bf1a5768adb01e60398401b85d5621",
            "healthy": true,
            "exit_code": 0,
            "stop_mode": "graceful",
            "security": {
              "user": "65532:65532",
              "read_only": true,
              "capabilities_dropped": true,
              "no_new_privileges": true,
              "host_ports": [],
              "mounts": [],
              "internal_network": true
            }
          },
          "inventory": {
            "image_id": "sha256:968ea2314c56405c938ed9d5835d5ebb7015711aa4960856531e65e6289449c8",
            "healthy": true,
            "exit_code": 0,
            "stop_mode": "graceful",
            "security": {
              "user": "65532:65532",
              "read_only": true,
              "capabilities_dropped": true,
              "no_new_privileges": true,
              "host_ports": [],
              "mounts": [],
              "internal_network": true
            }
          },
          "payments": {
            "image_id": "sha256:ca8d88c079db4033b6c2a85c084433f09e1b19accbff7fdeb3478d755a08d2bc",
            "healthy": true,
            "exit_code": 0,
            "stop_mode": "graceful",
            "security": {
              "user": "65532:65532",
              "read_only": true,
              "capabilities_dropped": true,
              "no_new_privileges": true,
              "host_ports": [],
              "mounts": [],
              "internal_network": true
            }
          }
        }
      },
      {
        "architecture": "arm64",
        "artifact_id": 11979852311,
        "url": "https://github.com/Includeoyasi/faultline-lab/actions/runs/38747476656/artifacts/11979852311",
        "retention_days": 14,
        "result": "PASS",
        "built_merge_revision": "cf9462c321e939279bd926f367ef4fde3015ef71",
        "source_dirty": false,
        "docker_platform": "linux/aarch64",
        "cleanup": "PASS",
        "isolation_after_orders_stop": "PASS",
        "services": {
          "orders": {
            "image_id": "sha256:9c5d84c9c9ed97e51ad63a76ad3e4ab89497f17fce29b814a9441b27968d2f2b",
            "healthy": true,
            "exit_code": 0,
            "stop_mode": "graceful",
            "security": {
              "user": "65532:65532",
              "read_only": true,
              "capabilities_dropped": true,
              "no_new_privileges": true,
              "host_ports": [],
              "mounts": [],
              "internal_network": true
            }
          },
          "inventory": {
            "image_id": "sha256:a997c2b03757626902f16851d3a8661784c3af348f1ab97c564966297fd72d32",
            "healthy": true,
            "exit_code": 0,
            "stop_mode": "graceful",
            "security": {
              "user": "65532:65532",
              "read_only": true,
              "capabilities_dropped": true,
              "no_new_privileges": true,
              "host_ports": [],
              "mounts": [],
              "internal_network": true
            }
          },
          "payments": {
            "image_id": "sha256:efa7621e0f31b72adf6d7dd2ef61a7d97e01ef20a0ca4834e37afc6bd61b9aef",
            "healthy": true,
            "exit_code": 0,
            "stop_mode": "graceful",
            "security": {
              "user": "65532:65532",
              "read_only": true,
              "capabilities_dropped": true,
              "no_new_privileges": true,
              "host_ports": [],
              "mounts": [],
              "internal_network": true
            }
          }
        }
      }
    ],
    "local_windows": "PASS",
    "macos": "NOT_RUN",
    "performance_capacity": "NOT_MEASURED",
    "windows": {
      "source_head": "d81b51ce2eb46f56d698fd8bf5f9966b91a16036",
      "source_dirty": false,
      "run_id": "8119b50d-f910-412a-9bd2-ab50fccd1b71",
      "result": "PASS",
      "report": "docs/development/evidence/p01/windows-2026-10-15/report.json",
      "report_sha256": "b08d9deabadc65c765d9173482e9171782aeff28affb5e7fd9bf0b201d5add65",
      "environment": {
        "desktop": "4.95.0",
        "docker": "29.8.2",
        "remaining_containers_before_stop": [],
        "cpus": 2,
        "memory_bytes": 1996603392,
        "wsl_config": {
          "memory": "2GB",
          "processors": 2,
          "swap": "1GB"
        },
        "platform": "linux/x86_64",
        "observed_at": "2026-10-15T22:49:35.706453+03:00"
      },
      "cleanup": {
        "containers_remaining": 0,
        "wsl_shutdown_exit": 0,
        "images_and_build_cache_preserved": true,
        "desktop_stop_exit": 0,
        "remaining_development_processes": [],
        "installed_tools_preserved": true,
        "docker_autostart": false,
        "docker_run_entry_present": false,
        "observed_at": "2026-10-15T22:49:53.4008625+03:00",
        "reboot_performed": false
      },
      "services": {
        "orders": {
          "id": "0e35127ad00aa13be4f551edf6b032c00ffc4894081efb4e4c33cc723fff9b04",
          "image_id": "sha256:c15eb12d4ab1369a2c49fd8da5f6ca178c530d39cf0f19c37d4ad513eac4574f",
          "architecture": "linux/x86_64",
          "healthy": true,
          "security": {
            "user": "65532:65532",
            "read_only": true,
            "capabilities_dropped": true,
            "no_new_privileges": true,
            "host_ports": [],
            "mounts": [],
            "internal_network": true
          },
          "stopped": true,
          "exit_code": 0,
          "oom_killed": false,
          "stop_mode": "graceful"
        },
        "inventory": {
          "id": "926b6baf437c4b0d0a01302d158872112689ddf89472c58a56f0d66c8f621d84",
          "image_id": "sha256:1fb289e7407f8559f671e9ec276faac6715f7c6e364220f70b30839624c7f778",
          "architecture": "linux/x86_64",
          "healthy": true,
          "security": {
            "user": "65532:65532",
            "read_only": true,
            "capabilities_dropped": true,
            "no_new_privileges": true,
            "host_ports": [],
            "mounts": [],
            "internal_network": true
          },
          "stopped": true,
          "exit_code": 0,
          "oom_killed": false,
          "stop_mode": "graceful"
        },
        "payments": {
          "id": "8a8c51c0112100bbdf7b981f9ebc1f9c4e099baaa24654a85eb85e956f55b868",
          "image_id": "sha256:3481c4f6aea54924e0f16b1da34fb2082bef2e937f5477fe077f5bce8b4eec35",
          "architecture": "linux/x86_64",
          "healthy": true,
          "security": {
            "user": "65532:65532",
            "read_only": true,
            "capabilities_dropped": true,
            "no_new_privileges": true,
            "host_ports": [],
            "mounts": [],
            "internal_network": true
          },
          "stopped": true,
          "exit_code": 0,
          "oom_killed": false,
          "stop_mode": "graceful"
        }
      },
      "limitations": [
        "No business operations or capacity measurement",
        "Personal macOS NOT_RUN"
      ]
    }
  },
  "audit_review": {
    "date": "2026-10-16",
    "request": "Read attached audit carefully, give concise summary and correction plans; agree further action afterward",
    "archive_name": "faultline-lab-agent-report-v2.zip",
    "archive_sha256": "224da358bfe06ab2cb4d39ed8fce9627b5ff7cd42b5111966e29be69e9859bb0",
    "audited_commit": "3ca290d818a98fc986fe42ae51dcfc0546673fb8",
    "remote_main_observed_commit": "3ca290d818a98fc986fe42ae51dcfc0546673fb8",
    "remote_observation": "Read-only GitHub branch API on2026-10-16; open issues9/12/16 verified; no server administration performed",
    "bundle_verification": {
      "manifest_files_checked": 144,
      "manifest_mismatches": [],
      "snapshot_files_checked_against_git": 79,
      "snapshot_mismatches": [],
      "changed_since_audit": [
        "AGENTS.md",
        "CONTEXT_STATE.json",
        "PROJECT_CONTEXT.md",
        "README.md",
        "README.ru.md",
        "docs/README.md",
        "docs/ROADMAP.md",
        "docs/adr/0002-record-review-decisions.md",
        "docs/adr/0003-go-container-smoke.md",
        "docs/architecture/DESIGN.md",
        "docs/development/ENGINEERING_DECISIONS.md",
        "docs/development/ENVIRONMENT.md",
        "docs/development/TASK_TRACKING.md",
        "docs/development/WORKFLOW.md",
        "docs/development/backlog/GO_SMOKE.md",
        "docs/development/packets/P01.md",
        "tools/context_intro.md"
      ],
      "implementation_files_changed": []
    },
    "archived_test_evidence": {
      "origin": "Auditor machine: native macOS ARM64 and Docker Linux ARM64, original audit session",
      "runtime_rerun_in_report_v2": false,
      "runtime_rerun_in_this_review": false,
      "independently_recounted_top_level_pass": 15,
      "independently_recounted_subtest_pass": 59,
      "pass_events": 74,
      "fail_events": 0,
      "coverage_internal_service": {
        "statements": 128,
        "covered": 116,
        "percent": 90.6
      },
      "total_instrumented_coverage": {
        "statements": 140,
        "covered": 116,
        "percent": 82.9
      },
      "note": "Process binaries lack coverage instrumentation; 0% main is not absence of process tests. These results do not establish business correctness or capacity. Auditor Mac evidence does not silently change personal Mac profile."
    },
    "findings": [
      {
        "id": "F01",
        "report_kind": "confirmed_defect",
        "report_priority": "P1",
        "priority_is_not_security_severity": true,
        "assessment": "SIGINT завершает Node smoke до итогового report/cleanup; в evidence остались три контейнера, затем удалённые аудитором.",
        "owner_decision": "Перенести smoke на Go по #9; исправление отмены, очистки и отчёта входит в перенос. Возможно объединение с небольшими исправлениями; точный состав PR ещё не выбран.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "Срочно исправить поддерживаемую отмену и регрессию; не ждать полного переноса runner на Go.",
        "assessment_at_review": "CURRENT_CONFIRMED_BY_ARCHIVED_REPRO_AND_UNCHANGED_SOURCE"
      },
      {
        "id": "F02",
        "report_kind": "state_drift",
        "report_priority": "P1",
        "priority_is_not_security_severity": true,
        "assessment": "P01 DONE отражён в нашей ветке документации; main по свежему API всё ещё 3ca290d с прежним IN_REVIEW/4 из8.",
        "owner_decision": "GitHub — единственный источник актуальных статусов; контекст хранит решения и датированную историю. Исполнитель — AI-Your-Bro; остальные агенты — потенциальные пользователи, подписчики или тестировщики.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "Довести уже подготовленное обновление статуса до main в согласованной публикации; не считать приёмкой G1.",
        "assessment_at_review": "ADDRESSED_IN_DOCS_BRANCH_MAIN_STILL_STALE"
      },
      {
        "id": "F03",
        "report_kind": "report_provenance_ambiguity",
        "report_priority": "P2",
        "priority_is_not_security_severity": true,
        "assessment": "services.architecture копирует Docker backend platform; это не измерение платформы конкретного image ID.",
        "owner_decision": "Отложить выбор платформенных метаданных; открытый пункт #20, вернуться на ближайшем планировании P02 до performance-сравнений.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "Разделить backend/requested/actual image platform до использования отчётов для сравнений.",
        "assessment_at_review": "CURRENT_METADATA_AMBIGUITY"
      },
      {
        "id": "F04",
        "report_kind": "documented_configuration_risk",
        "report_priority": "P2",
        "priority_is_not_security_severity": true,
        "assessment": "Для app20s/container15s с активным Watch сохранён exit137; default10s/15s не сломан.",
        "owner_decision": "Сохранить независимые app/container timeout; проверять совместимость и запас в профиле запуска.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "Предложить deployment preflight таймаутов или единую конфигурацию с запасом, проверить валидную и несовместимую пары.",
        "assessment_at_review": "DOCUMENTED_OVERRIDE_RISK_DEFAULTS_OK"
      },
      {
        "id": "F05",
        "report_kind": "documented_lifecycle_tradeoff",
        "report_priority": "P3",
        "priority_is_not_security_severity": true,
        "assessment": "Долгоживущий health Watch удерживает drain до forced exit; штатный Compose использует unary Check.",
        "owner_decision": "До постоянных Watch-клиентов завершать health-подписки отдельно, сохраняя graceful drain бизнес-RPC.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "Выбрать политику инфраструктурных streams перед внедрением Watch-клиентов; сейчас не считать регрессией P01.",
        "assessment_at_review": "DOCUMENTED_BEHAVIOR_DEFER_POLICY"
      },
      {
        "id": "F06",
        "report_kind": "future_acceptance_gap",
        "report_priority": "P2",
        "priority_is_not_security_severity": true,
        "assessment": "Smoke не является полноценным manifest для A/B; лимиты контейнеров не заданы отдельно от бюджета VM.",
        "owner_decision": "Вводить описание запуска постепенно: фактическая среда и ресурсы в P02, полный контракт перед сравнительными экспериментами.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "В P02 фиксировать фактические ресурсы; полную schema/правила INCONCLUSIVE подготовить перед экспериментами.",
        "assessment_at_review": "FUTURE_EXPERIMENT_ACCEPTANCE_WORK"
      },
      {
        "id": "F07",
        "report_kind": "artifact_reproducibility_risk",
        "report_priority": "P2",
        "priority_is_not_security_severity": true,
        "assessment": "Генератор всё ещё использует изменяемый snapshot_ref для внешних дополнений; полные встроенные тексты сохраняются.",
        "owner_decision": "Материалы снимка закреплять за существующим SHA коммита; актуальную навигацию отделять.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "Фиксировать внешние доказательные ссылки по source SHA, отделив актуальную навигацию.",
        "assessment_at_review": "CURRENT_REPRODUCIBILITY_RISK"
      },
      {
        "id": "F08",
        "report_kind": "checker_contract_gap",
        "report_priority": "P3",
        "priority_is_not_security_severity": true,
        "assessment": "Regex checker покрывает Windows Users, но пропускает Unix home paths. Реальная утечка не обнаружена.",
        "owner_decision": "Добавить распознавание домашних путей macOS/Linux и положительные/отрицательные синтетические примеры.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "Добавить синтетические positive/negative fixtures и Unix patterns либо честно сузить описание проверки.",
        "assessment_at_review": "CURRENT_CHECKER_CONTRACT_GAP"
      },
      {
        "id": "F09",
        "report_kind": "design_risk",
        "report_priority": "P2",
        "priority_is_not_security_severity": true,
        "assessment": "Нет исполняемого протокола защиты от старого publisher после lease takeover; outbox ещё отсутствует.",
        "owner_decision": "Разобрать протокол порядка и старого publisher в Kafka/P02 #16; outbox и его takeover-проверку реализовать на соответствующем этапе.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "До outbox определить ordering/fencing или безопасную обработку старых версий и проверить takeover-сценарий; дельту G0 согласовать при необходимости.",
        "assessment_at_review": "FUTURE_DESIGN_QUESTION_NOT_RUNTIME_BUG"
      },
      {
        "id": "F10",
        "report_kind": "design_risk",
        "report_priority": "P2",
        "priority_is_not_security_severity": true,
        "assessment": "Не определены сроки хранения ключей идемпотентности и cancellation barriers относительно replay/GC.",
        "owner_decision": "Хранить ключи и защитные записи всю жизнь изолированного набора; согласованный reset не допускает влияния старых сообщений на новый запуск.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "До операций/очистки определить lifetime namespace и гарантии replay; допустимый MVP — без GC защитных записей.",
        "assessment_at_review": "FUTURE_DESIGN_QUESTION_NOT_RUNTIME_BUG"
      },
      {
        "id": "F11",
        "report_kind": "contract_clarification",
        "report_priority": "P2",
        "priority_is_not_security_severity": true,
        "assessment": "Исключительный RELEASED при неизвестном Commit описан неоднозначно; формального противоречия автомата не доказано.",
        "owner_decision": "Исключительный RELEASED после неизвестного Commit учитывать диагностируемым незавершённым восстановлением, например NEEDS_REVIEW; не объявлять успешным terminal outcome.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "Перед соответствующим контрактом задать состояние, ошибку, durable step и критерий завершения; новые переходы G0 только после решения владельца.",
        "assessment_at_review": "FUTURE_CONTRACT_CLARIFICATION"
      },
      {
        "id": "F12",
        "report_kind": "experiment_design_risk",
        "report_priority": "P2",
        "priority_is_not_security_severity": true,
        "assessment": "Одинаковый seed не гарантирует одинаковые плановые исходы при разных logical IDs/RNG schedule; loadgen ещё отсутствует.",
        "owner_decision": "Плановый исход зависит от logical operation ID, seed и версии алгоритма; namespace запуска отдельный; сохраняется hash канонического плана.",
        "tracking_issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
        "initial_proposal": "Перед A/B закрепить logical operations, RNG/outcome version и workload hash; проверить детерминизм при разной конкурентности.",
        "assessment_at_review": "FUTURE_EXPERIMENT_DESIGN_QUESTION"
      },
      {
        "id": "F13",
        "report_kind": "maintenance_debt",
        "report_priority": "P3",
        "priority_is_not_security_severity": true,
        "assessment": "Ручные статусы повторяются в README/документах/intro; полный PROJECT_CONTEXT полезен как архив, но тяжёл для первого входа.",
        "assessment_at_review": "CURRENT_MAINTENANCE_DEBT",
        "initial_report_proposal": "Генерировать краткий handoff/статусы из canonical state, сохранив полный переносимый архив и подробные инварианты."
      }
    ],
    "local_notes": "Detailed review saved under ignored .local; audit originals include host-specific paths and are not published",
    "scope_note": "Dated review of supplied archive; remediation directions and live execution are tracked in issue20. Saved runtime evidence is not a rerun.",
    "owner_review": {
      "date": "2026-10-16",
      "source": "Direct owner chat message reviewing F01–F12 and requesting one consolidated issue",
      "issue_url": "https://github.com/Includeoyasi/faultline-lab/issues/20",
      "author_of_summary": "AI-Your-Bro; not an owner GitHub approval",
      "selected_directions": {
        "F01": "Перенести smoke на Go по #9; исправление отмены, очистки и отчёта входит в перенос. Возможно объединение с небольшими исправлениями; точный состав PR ещё не выбран.",
        "F02": "GitHub — единственный источник актуальных статусов; контекст хранит решения и датированную историю. Исполнитель — AI-Your-Bro; остальные агенты — потенциальные пользователи, подписчики или тестировщики.",
        "F03": "Отложить выбор платформенных метаданных; открытый пункт #20, вернуться на ближайшем планировании P02 до performance-сравнений.",
        "F04": "Сохранить независимые app/container timeout; проверять совместимость и запас в профиле запуска.",
        "F05": "До постоянных Watch-клиентов завершать health-подписки отдельно, сохраняя graceful drain бизнес-RPC.",
        "F06": "Вводить описание запуска постепенно: фактическая среда и ресурсы в P02, полный контракт перед сравнительными экспериментами.",
        "F07": "Материалы снимка закреплять за существующим SHA коммита; актуальную навигацию отделять.",
        "F08": "Добавить распознавание домашних путей macOS/Linux и положительные/отрицательные синтетические примеры.",
        "F09": "Разобрать протокол порядка и старого publisher в Kafka/P02 #16; outbox и его takeover-проверку реализовать на соответствующем этапе.",
        "F10": "Хранить ключи и защитные записи всю жизнь изолированного набора; согласованный reset не допускает влияния старых сообщений на новый запуск.",
        "F11": "Исключительный RELEASED после неизвестного Commit учитывать диагностируемым незавершённым восстановлением, например NEEDS_REVIEW; не объявлять успешным terminal outcome.",
        "F12": "Плановый исход зависит от logical operation ID, seed и версии алгоритма; namespace запуска отдельный; сохраняется hash канонического плана."
      },
      "F13": "Outside the twelve requested sections; original audit observation retained as history"
    },
    "local_head_at_review": "511b9a85c7006cc7042263e0b87e7eb7fb8b116c"
  },
  "source_documents": [
    {
      "path": "README.md",
      "anchor": "appendix-readme-en",
      "title": "Публичное введение на английском"
    },
    {
      "path": "README.ru.md",
      "anchor": "appendix-readme-ru",
      "title": "Публичное введение на русском"
    },
    {
      "path": "docs/README.md",
      "anchor": "appendix-docs-index",
      "title": "Карта документации"
    },
    {
      "path": "docs/ROADMAP.md",
      "anchor": "appendix-roadmap",
      "title": "Полный roadmap, архитектура, эксперименты и методика"
    },
    {
      "path": "docs/architecture/DESIGN.md",
      "anchor": "appendix-design-review",
      "title": "Принятый архитектурный пакет G0"
    },
    {
      "path": "docs/architecture/PROCESSES.md",
      "anchor": "appendix-g0-processes",
      "title": "Подробные визуальные схемы G0"
    },
    {
      "path": "docs/development/WORKFLOW.md",
      "anchor": "appendix-development-process",
      "title": "Роли, ревью, контрольные точки и обновление контекста"
    },
    {
      "path": "docs/experiments/REPORT_TEMPLATE.md",
      "anchor": "appendix-report-template",
      "title": "Полный шаблон отчёта"
    },
    {
      "path": "AGENTS.md",
      "anchor": "appendix-agents",
      "title": "Локальные инструкции ИИ-исполнителю"
    },
    {
      "path": "docs/adr/0001-approved-g0.md",
      "anchor": "appendix-adr-0001",
      "title": "Запись принятия архитектуры G0"
    },
    {
      "path": "CONTRIBUTING.md",
      "anchor": "appendix-contributing",
      "title": "Участие, ревью и проверки"
    },
    {
      "path": "LICENSE",
      "anchor": "appendix-license",
      "title": "MIT License"
    },
    {
      "path": "docs/development/CODE_REVIEW.md",
      "anchor": "appendix-code-review",
      "title": "Правила Go code review и критерии принятия изменений"
    },
    {
      "path": "docs/architecture/LAB_BOUNDARY.md",
      "anchor": "appendix-lab-boundary",
      "title": "Общая лаборатория, адаптеры и кандидат второго бизнеса"
    },
    {
      "path": "docs/development/ENVIRONMENT.md",
      "anchor": "appendix-environment",
      "title": "Переносимая Docker-среда и распределение контейнеров"
    },
    {
      "path": "docs/development/RESOURCE_PROFILES.md",
      "anchor": "appendix-resource-profiles",
      "title": "Ресурсный снимок Windows и будущие профили macOS/сервера"
    },
    {
      "path": "docs/design/BANNER_OPTIONS.md",
      "anchor": "appendix-banner-options",
      "title": "Концепции баннера и выбранное оформление"
    },
    {
      "path": "docs/development/packets/P01.md",
      "anchor": "appendix-p01",
      "title": "Первый PR с кодом: задачи, зависимости и критерии приёмки"
    },
    {
      "path": "docs/development/TASK_TRACKING.md",
      "anchor": "appendix-task-tracking",
      "title": "GitHub Issues, статусы и переносимые снимки"
    },
    {
      "path": "docs/development/REPOSITORY_SECURITY.md",
      "anchor": "appendix-repository-security",
      "title": "Защита main, отдельный исполнитель и состояние подключения"
    },
    {
      "path": ".github/CODEOWNERS",
      "anchor": "appendix-codeowners",
      "title": "Действующий CODEOWNERS для человеческого ревью"
    },
    {
      "path": "docs/development/TOOLCHAIN.md",
      "anchor": "appendix-toolchain",
      "title": "Контейнерный toolchain и команды"
    },
    {
      "path": "go.mod",
      "anchor": "appendix-go-mod",
      "title": "Go-модуль и версии инструментов",
      "language": "go"
    },
    {
      "path": "go.sum",
      "anchor": "appendix-go-sum",
      "title": "Контрольные суммы модулей",
      "language": "text"
    },
    {
      "path": "cmd/orders/main.go",
      "anchor": "appendix-orders-main",
      "title": "Orders gRPC health entrypoint",
      "language": "go"
    },
    {
      "path": "cmd/inventory/main.go",
      "anchor": "appendix-inventory-main",
      "title": "Inventory gRPC health entrypoint",
      "language": "go"
    },
    {
      "path": "cmd/payments/main.go",
      "anchor": "appendix-payments-main",
      "title": "Payments gRPC health entrypoint",
      "language": "go"
    },
    {
      "path": "build/Dockerfile",
      "anchor": "appendix-toolchain-dockerfile",
      "title": "Единый multi-stage Dockerfile: toolchain и runtime",
      "language": "dockerfile"
    },
    {
      "path": "compose.yaml",
      "anchor": "appendix-compose-toolchain",
      "title": "Compose для временных инструментов",
      "language": "yaml"
    },
    {
      "path": "tools/go.sh",
      "anchor": "appendix-go-commands",
      "title": "Команды внутри Linux-контейнера",
      "language": "sh"
    },
    {
      "path": ".golangci.yml",
      "anchor": "appendix-golangci-config",
      "title": "Согласованные Go-линтеры",
      "language": "yaml"
    },
    {
      "path": ".github/workflows/go.yml",
      "anchor": "appendix-go-workflow",
      "title": "Проверки Go в нативном Linux CI",
      "language": "yaml"
    },
    {
      "path": ".dockerignore",
      "anchor": "appendix-dockerignore",
      "title": "Исключения Docker build context",
      "language": "text"
    },
    {
      "path": "docs/development/SERVICE_RUNTIME.md",
      "anchor": "appendix-service-runtime",
      "title": "Конфигурация, JSON-логи и gRPC health: поведение и проверки"
    },
    {
      "path": "internal/service/config.go",
      "anchor": "appendix-service-config",
      "title": "Валидация конфигурации",
      "language": "go"
    },
    {
      "path": "internal/service/log.go",
      "anchor": "appendix-service-logging",
      "title": "JSON-логи и build metadata",
      "language": "go"
    },
    {
      "path": "internal/service/run.go",
      "anchor": "appendix-service-run",
      "title": "Запуск gRPC health и немедленная остановка",
      "language": "go"
    },
    {
      "path": "internal/service/config_test.go",
      "anchor": "appendix-service-config-tests",
      "title": "Проверки конфигурации",
      "language": "go"
    },
    {
      "path": "internal/service/log_test.go",
      "anchor": "appendix-service-log-tests",
      "title": "Проверки JSON и build metadata",
      "language": "go"
    },
    {
      "path": "internal/service/run_test.go",
      "anchor": "appendix-service-run-tests",
      "title": "Проверки bind, TCP health и отмены контекста",
      "language": "go"
    },
    {
      "path": "tests/smoke/services_test.go",
      "anchor": "appendix-service-process-tests",
      "title": "Проверки трёх реальных процессов",
      "language": "go"
    },
    {
      "path": "configs/env/local.env.example",
      "anchor": "appendix-local-env-example",
      "title": "Явное локальное окружение: образец",
      "language": "dotenv"
    },
    {
      "path": "configs/env/container.env.example",
      "anchor": "appendix-container-env-example",
      "title": "Явное контейнерное окружение: образец",
      "language": "dotenv"
    },
    {
      "path": ".gitignore",
      "anchor": "appendix-gitignore",
      "title": "Исключение личных настроек и артефактов из Git",
      "language": "text"
    },
    {
      "path": "internal/service/command.go",
      "anchor": "appendix-service-command",
      "title": "Entrypoint command and bounded health probe",
      "language": "go"
    },
    {
      "path": "internal/service/command_test.go",
      "anchor": "appendix-service-command-tests",
      "title": "Health probe and command tests",
      "language": "go"
    },
    {
      "path": "internal/service/shutdown_test.go",
      "anchor": "appendix-shutdown-tests",
      "title": "Controlled active RPC and deadline tests",
      "language": "go"
    },
    {
      "path": "docs/development/CONTAINER_RUNTIME.md",
      "anchor": "appendix-container-runtime",
      "title": "Runtime build, health, shutdown, smoke and review limitations"
    },
    {
      "path": "docs/development/evidence/p01/WINDOWS_SMOKE.md",
      "anchor": "windows-smoke-proof",
      "title": "Windows smoke: итог и ограничения"
    },
    {
      "path": "docs/development/evidence/p01/windows-2026-10-15/report.json",
      "anchor": "windows-smoke-report",
      "title": "Полный Windows smoke report",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p01/windows-2026-10-15/orders.jsonl",
      "anchor": "windows-smoke-orders",
      "title": "Windows: Orders JSONL",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p01/windows-2026-10-15/inventory.jsonl",
      "anchor": "windows-smoke-inventory",
      "title": "Windows: Inventory JSONL",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p01/windows-2026-10-15/payments.jsonl",
      "anchor": "windows-smoke-payments",
      "title": "Windows: Payments JSONL",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p01/windows-2026-10-15/environment.json",
      "anchor": "windows-smoke-environment",
      "title": "Windows: фактический бюджет",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p01/windows-2026-10-15/cleanup.json",
      "anchor": "windows-smoke-cleanup",
      "title": "Windows: подтверждение остановки",
      "language": "json"
    },
    {
      "path": "docs/development/ENGINEERING_DECISIONS.md",
      "anchor": "appendix-engineering-decisions",
      "title": "Инженерные решения, вклад владельца и доказательства ревью"
    },
    {
      "path": "docs/adr/0002-record-review-decisions.md",
      "anchor": "appendix-adr-0002",
      "title": "ADR 0002: согласованный порядок фиксации решений"
    },
    {
      "path": "docs/adr/0003-go-container-smoke.md",
      "anchor": "appendix-adr-0003",
      "title": "ADR 0003: proposed Go container smoke"
    },
    {
      "path": "docs/development/backlog/GO_SMOKE.md",
      "anchor": "appendix-go-smoke-backlog",
      "title": "Backlog #9: перенос smoke на Go, объём и критерии"
    },
    {
      "path": "docs/development/packets/P02.md",
      "anchor": "appendix-p02",
      "title": "Предложенный P02: задачи, зависимости, критерии и порядок PR"
    },
    {
      "path": "docs/development/SMOKE.md",
      "title": "Go container smoke guide",
      "anchor": "go-smoke-guide"
    },
    {
      "path": "cmd/smoke/main.go",
      "title": "Go smoke entrypoint",
      "anchor": "go-smoke-main",
      "language": "go"
    },
    {
      "path": "cmd/smoke/main_test.go",
      "title": "Smoke CLI tests",
      "anchor": "go-smoke-main-tests",
      "language": "go"
    },
    {
      "path": "internal/smoke/run.go",
      "title": "Smoke lifecycle",
      "anchor": "go-smoke-run",
      "language": "go"
    },
    {
      "path": "internal/smoke/report.go",
      "title": "Smoke report",
      "anchor": "go-smoke-report",
      "language": "go"
    },
    {
      "path": "internal/smoke/checks.go",
      "title": "Smoke container checks",
      "anchor": "go-smoke-checks",
      "language": "go"
    },
    {
      "path": "internal/smoke/config.go",
      "title": "Compose shutdown preflight",
      "anchor": "go-smoke-config",
      "language": "go"
    },
    {
      "path": "internal/command/exec.go",
      "title": "Bounded command output",
      "anchor": "go-smoke-exec",
      "language": "go"
    },
    {
      "path": "internal/command/exec_unix.go",
      "title": "Unix command process groups",
      "anchor": "go-smoke-unix",
      "language": "go"
    },
    {
      "path": "internal/command/exec_windows.go",
      "title": "Windows command jobs",
      "anchor": "go-smoke-windows",
      "language": "go"
    },
    {
      "path": "internal/smoke/run_test.go",
      "title": "Smoke failure and cancellation tests",
      "anchor": "go-smoke-run-tests",
      "language": "go"
    },
    {
      "path": "internal/smoke/config_test.go",
      "title": "Shutdown budget tests",
      "anchor": "go-smoke-config-tests",
      "language": "go"
    },
    {
      "path": "internal/command/exec_test.go",
      "title": "Command tree cancellation tests",
      "anchor": "go-smoke-exec-tests",
      "language": "go"
    },
    {
      "path": "tools/doc_safety.mjs",
      "title": "Documentation path hygiene",
      "anchor": "doc-path-hygiene",
      "language": "javascript"
    },
    {
      "path": "tools/doc_safety.test.mjs",
      "title": "Path hygiene examples",
      "anchor": "doc-path-hygiene-tests",
      "language": "javascript"
    },
    {
      "path": "tests/smokecli/interrupt_test.go",
      "title": "Host CLI SIGINT regression",
      "anchor": "go-smoke-cli-interrupt",
      "language": "go"
    },
    {
      "path": "internal/command/exec_unix_test.go",
      "title": "Unix orphan plugin lifecycle regression",
      "anchor": "source-smoke-exec-unix-test",
      "language": "go"
    },
    {
      "path": "docs/development/evidence/go-smoke/go-equivalence-report.json",
      "title": "Go smoke evidence: go-equivalence-report.json",
      "anchor": "go-smoke-evidence-go-equivalence-report",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/go-smoke/go-smoke-equivalence.json",
      "title": "Go smoke evidence: go-smoke-equivalence.json",
      "anchor": "go-smoke-evidence-go-smoke-equivalence",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/go-smoke/go-smoke-preflight.json",
      "title": "Go smoke evidence: go-smoke-preflight.json",
      "anchor": "go-smoke-evidence-go-smoke-preflight",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/go-smoke/go-smoke-windows-interrupt.json",
      "title": "Go smoke evidence: go-smoke-windows-interrupt.json",
      "anchor": "go-smoke-evidence-go-smoke-windows-interrupt",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/go-smoke/no-backend-report.json",
      "title": "Go smoke evidence: no-backend-report.json",
      "anchor": "go-smoke-evidence-no-backend-report",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/go-smoke/node-baseline-report.json",
      "title": "Go smoke evidence: node-baseline-report.json",
      "anchor": "go-smoke-evidence-node-baseline-report",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/go-smoke/VALIDATION.md",
      "title": "Go smoke evidence: VALIDATION.md",
      "anchor": "go-smoke-evidence-VALIDATION"
    },
    {
      "path": "docs/development/evidence/go-smoke/windows-cleanup.json",
      "title": "Go smoke evidence: windows-cleanup.json",
      "anchor": "go-smoke-evidence-windows-cleanup",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/go-smoke/windows-final-inventory.jsonl",
      "title": "Go smoke evidence: windows-final-inventory.jsonl",
      "anchor": "go-smoke-evidence-windows-final-inventory",
      "language": "jsonl"
    },
    {
      "path": "docs/development/evidence/go-smoke/windows-final-orders.jsonl",
      "title": "Go smoke evidence: windows-final-orders.jsonl",
      "anchor": "go-smoke-evidence-windows-final-orders",
      "language": "jsonl"
    },
    {
      "path": "docs/development/evidence/go-smoke/windows-final-payments.jsonl",
      "title": "Go smoke evidence: windows-final-payments.jsonl",
      "anchor": "go-smoke-evidence-windows-final-payments",
      "language": "jsonl"
    },
    {
      "path": "docs/development/evidence/go-smoke/windows-final-report.json",
      "title": "Go smoke evidence: windows-final-report.json",
      "anchor": "go-smoke-evidence-windows-final-report",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/go-smoke/windows-interrupt-report.json",
      "title": "Go smoke evidence: windows-interrupt-report.json",
      "anchor": "go-smoke-evidence-windows-interrupt-report",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p02/windows-runtime-report.json",
      "title": "docs/development/evidence/p02/windows-runtime-report.json",
      "anchor": "p02-docs-development-evidence-p02-windows-runtime-report-json",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p02/windows-cleanup.json",
      "title": "docs/development/evidence/p02/windows-cleanup.json",
      "anchor": "p02-docs-development-evidence-p02-windows-cleanup-json",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p02/windows-postgres-report.json",
      "title": "docs/development/evidence/p02/windows-postgres-report.json",
      "anchor": "p02-docs-development-evidence-p02-windows-postgres-report-json",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p02/postgres-image.json",
      "title": "docs/development/evidence/p02/postgres-image.json",
      "anchor": "p02-docs-development-evidence-p02-postgres-image-json",
      "language": "json"
    },
    {
      "path": "docs/development/evidence/p02/POSTGRES.md",
      "title": "docs/development/evidence/p02/POSTGRES.md",
      "anchor": "p02-docs-development-evidence-p02-postgres-md"
    },
    {
      "path": "internal/command/command.go",
      "title": "internal/command/command.go",
      "anchor": "p02-internal-command-command-go",
      "language": "go"
    },
    {
      "path": "internal/database/migrate.go",
      "title": "internal/database/migrate.go",
      "anchor": "p02-internal-database-migrate-go",
      "language": "go"
    },
    {
      "path": "internal/database/integration_test.go",
      "title": "internal/database/integration_test.go",
      "anchor": "p02-internal-database-integration-test-go",
      "language": "go"
    },
    {
      "path": "internal/database/fingerprint.go",
      "title": "internal/database/fingerprint.go",
      "anchor": "p02-internal-database-fingerprint-go",
      "language": "go"
    },
    {
      "path": "internal/database/config_test.go",
      "title": "internal/database/config_test.go",
      "anchor": "p02-internal-database-config-test-go",
      "language": "go"
    },
    {
      "path": "internal/database/config.go",
      "title": "internal/database/config.go",
      "anchor": "p02-internal-database-config-go",
      "language": "go"
    },
    {
      "path": "internal/database/bootstrap.go",
      "title": "internal/database/bootstrap.go",
      "anchor": "p02-internal-database-bootstrap-go",
      "language": "go"
    },
    {
      "path": "configs/postgres/local.json",
      "title": "configs/postgres/local.json",
      "anchor": "p02-configs-postgres-local-json",
      "language": "json"
    },
    {
      "path": "cmd/labdata/main_test.go",
      "title": "cmd/labdata/main_test.go",
      "anchor": "p02-cmd-labdata-main-test-go",
      "language": "go"
    },
    {
      "path": "cmd/labdata/main.go",
      "title": "cmd/labdata/main.go",
      "anchor": "p02-cmd-labdata-main-go",
      "language": "go"
    },
    {
      "path": "migrations/inventory/0002_catalog.sql",
      "title": "migrations/inventory/0002_catalog.sql",
      "anchor": "p02-migrations-inventory-0002-catalog-sql",
      "language": "sql"
    },
    {
      "path": "migrations/inventory/0001_schema.sql",
      "title": "migrations/inventory/0001_schema.sql",
      "anchor": "p02-migrations-inventory-0001-schema-sql",
      "language": "sql"
    },
    {
      "path": "migrations/embed.go",
      "title": "migrations/embed.go",
      "anchor": "p02-migrations-embed-go",
      "language": "go"
    },
    {
      "path": "migrations/seed/inventory.sql",
      "title": "migrations/seed/inventory.sql",
      "anchor": "p02-migrations-seed-inventory-sql",
      "language": "sql"
    },
    {
      "path": "migrations/orders/0001_schema.sql",
      "title": "migrations/orders/0001_schema.sql",
      "anchor": "p02-migrations-orders-0001-schema-sql",
      "language": "sql"
    },
    {
      "path": "migrations/payments/0001_schema.sql",
      "title": "migrations/payments/0001_schema.sql",
      "anchor": "p02-migrations-payments-0001-schema-sql",
      "language": "sql"
    },
    {
      "path": "internal/datalab/run_test.go",
      "title": "internal/datalab/run_test.go",
      "anchor": "p02-internal-datalab-run-test-go",
      "language": "go"
    },
    {
      "path": "internal/datalab/run.go",
      "title": "internal/datalab/run.go",
      "anchor": "p02-internal-datalab-run-go",
      "language": "go"
    },
    {
      "path": "compose.postgres.yaml",
      "title": "compose.postgres.yaml",
      "anchor": "p02-compose-postgres-yaml",
      "language": "yaml"
    },
    {
      "path": "internal/smoke/exec_adapter.go",
      "title": "internal/smoke/exec_adapter.go",
      "anchor": "p02-internal-smoke-exec-adapter-go",
      "language": "go"
    },
    {
      "path": "docs/development/POSTGRES.md",
      "title": "docs/development/POSTGRES.md",
      "anchor": "p02-docs-development-postgres-md"
    }
  ],
  "change_history": [
    {
      "revision": 1,
      "date": "2026-10-13",
      "summary": "Первый переносимый контекст; явное принятие G0; полные проектные приложения; механизм обновления и проверки свежести",
      "trigger": "Пользователь: G0 полностью принят. Создай файл для хранения контекста задачи..."
    },
    {
      "revision": 2,
      "date": "2026-10-13",
      "summary": "Проект получил имя Faultline Lab; создан GitHub-репозиторий, оформлены README/обложка/структура, выбрана MIT. Подготовка первого коммита.",
      "trigger": "Пользователь поручил инициализировать Git, создать публичный репозиторий и первый коммит; затем выбрал MIT."
    },
    {
      "revision": 3,
      "date": "2026-10-13",
      "summary": "Первый коммит опубликован; CI успешно проверен; репозиторий и модель human-led AI development оформлены; переносимый контекст фиксирует завершение bootstrap.",
      "trigger": "Фактические результаты первого push, GitHub Actions и проверки публичного README."
    },
    {
      "revision": 4,
      "date": "2026-10-13",
      "summary": "Подготовлен стандарт Go code review и обновлены инструкции/шаблон PR; создана ветка docs/go-review-policy. Go CI включается с первым модулем; человеческое одобрение не заявляется.",
      "trigger": "Прямое поручение владельца проекта установить современные правила Go code review."
    },
    {
      "revision": 5,
      "date": "2026-10-13",
      "summary": "Зафиксированы PR #1, опубликованная ветка, успешные локальные проверки и CI 214eb34; правила готовы к ревью владельца проекта, слияние не выполнялось.",
      "trigger": "Фактические результаты Git push, GitHub PR creation и Actions."
    },
    {
      "revision": 6,
      "date": "2026-10-13",
      "summary": "Замечено и проверено слияние PR #1 аккаунтом владельца проекта. Правила действуют в main; статусы reconciled с merge c1e55fd и успешным CI. Handoff обновляется отдельной веткой.",
      "trigger": "GitHub PR metadata и фактическое обновление origin/main во время текущей задачи."
    },
    {
      "revision": 7,
      "date": "2026-10-13",
      "summary": "В переносимый снимок добавлен PR #2 для отражения принятия Go review standard; правила уже действуют через PR #1 в main.",
      "trigger": "Создание отдельного handoff PR после наблюдаемого слияния владельцем проекта."
    },
    {
      "revision": 8,
      "date": "2026-10-13",
      "summary": "Унифицировано название роли человека, обновлены исходные документы и переносимый контекст для следующего коммита; отражено слияние PR #2 и сохранено согласованное направление развития лаборатории.",
      "trigger": "Прямое поручение владельца проекта подготовить изменения терминологии к следующему коммиту; предшествующее голосовое обсуждение границы лаборатории и бизнес-системы."
    },
    {
      "revision": 9,
      "date": "2026-10-13",
      "summary": "Зафиксирована первая проверка среды перед G1-P01: Go и Docker не обнаружены как команды текущего терминала; подготовленные документы по-прежнему ожидают следующего коммита.",
      "trigger": "Владелец проекта спросил о необходимых приготовлениях перед следующим пакетом реализации."
    },
    {
      "revision": 10,
      "date": "2026-10-13",
      "summary": "Сохранены граница лаборатории и кандидат второго бизнеса, согласованный Docker-подход, анализ контейнеров и ресурсы Windows; применён выбранный пастельный баннер. Все изменения подготовлены к следующему коммиту, без runtime и публикации.",
      "trigger": "Голосовое обсуждение и прямые поручения сохранить решения/измерить ресурсы/подготовить три обложки; затем выбор владельцем варианта 2"
    },
    {
      "revision": 11,
      "date": "2026-10-13",
      "summary": "Первый кодовый PR G1-P01 декомпозирован на 8 задач с критериями и зависимостями; предложен GitHub Issues и уточнены статусы. Реализация 0/8, внешние карточки и PR с Go-кодом не созданы.",
      "trigger": "Владелец проекта запросил список задач ближайшего MR, статус декомпозиции, критерии приёмки и место учёта выполнения"
    },
    {
      "revision": 12,
      "date": "2026-10-13",
      "summary": "Первый кодовый пакет переименован из G1-P01 в P01. В документах и контексте явно зафиксирован порядок G0 → кодовые PR с ревью → G1; состав первого PR и статусы не изменились.",
      "trigger": "Владелец проекта указал на неоднозначность обозначения и подтвердил место первого PR между G0 и G1"
    },
    {
      "revision": 13,
      "date": "2026-10-13",
      "summary": "Принят GitHub Issues и создан P01 #3. Согласованные документы, новая обложка и разделение G0/P01/G1 подготовлены к документационному PR; первый PR с кодом остаётся будущим P01.",
      "trigger": "Прямое поручение владельца проекта: молодец, фиксируй изменения и заводи задачи"
    },
    {
      "revision": 14,
      "date": "2026-10-13",
      "summary": "Сохранены фактические публикация коммита, PR #4 и успешный CI его исходного HEAD. Переносимый контекст и документы теперь содержат реальные ссылки на tracker и review. Main не изменён; код P01 не начат.",
      "trigger": "Завершение поручения владельца: фиксация изменений и создание задач"
    },
    {
      "revision": 15,
      "date": "2026-10-13",
      "summary": "Main защищён и настройки прочитаны обратно. Владелец предоставил AI-Your-Bro; профиль существует, приглашение write отправлено. Новые PR от Includeoyasi запрещены; публикация локального CODEOWNERS/контекста ждёт авторизации исполнителя. Наблюдено слияние PR #4.",
      "trigger": "Владелец поручил отдельный аккаунт исполнителя для всех новых PR и защиту main; затем предоставил логин AI-Your-Bro"
    },
    {
      "revision": 16,
      "date": "2026-10-13",
      "summary": "Подтверждены OAuth, принятие приглашения и write без admin у AI-Your-Bro. Настроена локальная Git identity, учтён прямой отказ владельца от 2FA. Защита main сохранена; подготовка первого security PR от исполнителя.",
      "trigger": "Владелец завершил device flow и сообщил готово; отдельным ответом отказался включать 2FA"
    },
    {
      "revision": 17,
      "date": "2026-10-13",
      "summary": "Первый PR нового порядка #5 создан от AI-Your-Bro к владельцу проекта. Сохранены реальная публикация, автор/коммиттер, requested review и исходный успешный CI. Main защищён; CODEOWNERS ещё не слит.",
      "trigger": "Завершение поручения по подключению исполнителя и защите main; handoff для человеческого ревью"
    },
    {
      "revision": 18,
      "date": "2026-10-15",
      "summary": "Слияние security PR #5 учтено. Реализуется первая задача P01-01: модуль, bootstrap entrypoints, toolchain, Go CI и команды; остальные семь задач не начаты. Контекст расширен полными исходниками первой задачи.",
      "trigger": "Владелец сообщил о merge PR #5 и поручил взять первую утверждённую задачу и создать PR для ревью"
    },
    {
      "revision": 19,
      "date": "2026-10-15",
      "summary": "Создан draft PR #6 от исполнителя; первый ARM64 CI выявил Git ownership bind mount. Образ доверяет точному /workspace для VCS metadata; повторные Go-проверки ожидаются.",
      "trigger": "Фактическая диагностика контейнерного CI первой задачи"
    },
    {
      "revision": 20,
      "date": "2026-10-15",
      "summary": "Go quality checks прошли на AMD64/ARM64. Исправлена запись bootstrap-бинарников в bind mount: toolchain использует стандартные capabilities Docker. Добавлен полный .dockerignore в переносимый снимок.",
      "trigger": "Фактическая проверка bootstrap build после исправления VCS"
    },
    {
      "revision": 21,
      "date": "2026-10-15",
      "summary": "P01-01 проверена на обеих Linux-архитектурах и предъявлена в ready PR #6. Оба Go checks стали обязательными в main. Issue отражает 1/8 проверено, 0/8 принято. Контекст содержит точные источники и CI evidence; следующие задачи не выполнялись.",
      "trigger": "Завершение первой задачи и передача владельцу на ревью"
    },
    {
      "revision": 22,
      "date": "2026-10-15",
      "summary": "В итоговом self-review устранена устаревшая строка 0/8 реализации в конце спецификации; все активные сводки различают 1/8 проверено и 0/8 принято. Уточнены локальный runtime prerequisite и дата security snapshot. Go-код, инструменты и критерии не менялись.",
      "trigger": "Сверка итоговых статусов перед передачей PR #6 владельцу"
    },
    {
      "revision": 23,
      "date": "2026-10-15",
      "summary": "Наблюдено принятие PR #6. По прямому поручению владельца P01-02–04 реализуются одним PR: config, JSON, health, unit/process tests и документация. Локальные Windows Go проверки прошли; Linux CI ожидается, приёмка новых задач не заявляется.",
      "trigger": "Merge P01-01 и поручение следующей группы задач"
    },
    {
      "revision": 24,
      "date": "2026-10-15",
      "summary": "Создан draft PR #7 от AI-Your-Bro, прикреплён к чату. Первый Linux CI остановлен staticcheck QF1001; проверка допустимых символов переписана без подавления линтера, ожидается повторный полный прогон. Новые задачи ещё не IN_REVIEW.",
      "trigger": "CI feedback и исправление перед готовностью к ревью"
    },
    {
      "revision": 25,
      "date": "2026-10-15",
      "summary": "Задачи P01-02–04 проверены на обеих нативных Linux-архитектурах и предъявлены в ready PR #7. Issue #3 отражает 4/8 проверено и 1/8 принято владельцем. Добавлены реальные JSON-записи и полный переносимый исходный код. P01-05–08 остаются будущими задачами; immediate stop не выдаётся за graceful shutdown.",
      "trigger": "Завершение согласованной группы задач и передача на ревью"
    },
    {
      "revision": 26,
      "date": "2026-10-15",
      "summary": "Финальная сверка переносимого контекста: repository.snapshot_ref и current_publication переключены с уже слитого PR #6 на текущий ready PR #7. Код и критерии не менялись; доказательства слияния #6 сохраняются в истории.",
      "trigger": "Проверка согласованности активных полей перед передачей владельцу"
    },
    {
      "revision": 27,
      "date": "2026-10-15",
      "summary": "Замечание владельца о сетевых адресах в коде реализовано в текущем PR #7: явный обязательный FL_GRPC_ADDR, образцы окружения, отказ без настройки, разделение listener/DNS/publication, игнорирование личных env. Локальные проверки PASS; новый Linux CI ожидается. Приёмка человеком не заявляется.",
      "trigger": "Прямое замечание владельца во время ревью PR #7"
    },
    {
      "revision": 28,
      "date": "2026-10-15",
      "summary": "Полный CI сетевого исправления PR #7 прошёл на AMD64/ARM64. Зафиксирован обязательный адрес, профили и безопасная локальная публикация; замечание ADDRESSED_VERIFIED_AWAITING_OWNER_REVIEW. Человеческое принятие новых задач и исправления ещё ожидается.",
      "trigger": "Проверка изменения по замечанию владельца и повторная передача на ревью"
    },
    {
      "revision": 29,
      "date": "2026-10-15",
      "summary": "Принятие/merge PR #7 подтверждены. По прямому поручению владельца P01-05–08 реализуются одним PR с понятными заголовками. Локальные Go lifecycle tests PASS; runtime/CI подготовлены, Docker/WSL Windows отсутствуют. Новый PR будет draft до обязательного Windows smoke.",
      "trigger": "Принятый предыдущий PR и поручение оставшихся задач"
    },
    {
      "revision": 30,
      "date": "2026-10-15",
      "summary": "Создан draft PR #8 от AI-Your-Bro и прикреплён к чату. Первый CI выявил зависимость сборки от неактивного toolchain-профиля; введён единый multi-stage Dockerfile с target toolchain/runtime. Повторный CI ожидается; Windows smoke по-прежнему NOT_RUN.",
      "trigger": "Публикация и исправление модели Compose по результату CI"
    },
    {
      "revision": 31,
      "date": "2026-10-15",
      "summary": "Полный P01-05–08 CI и runtime smoke PASS на AMD64/ARM64, оба отчёта/логи прочитаны. P01-05/07/08 завершили реализацию и проверки; P01-06 остаётся с невыполненным Windows container-критерием. Draft PR #8 сохранён; установка Docker/WSL запрошена у владельца, автоматическая перезагрузка не разрешена.",
      "trigger": "Нативные контейнерные доказательства и подготовка итогового draft review packet"
    },
    {
      "revision": 32,
      "date": "2026-10-15",
      "summary": "Зафиксировано разрешение владельца на Docker/WSL и реальные результаты установки: WSL требует admin, UAC отменён, Docker prerequisite VirtualMachinePlatform отсутствует. Windows smoke остаётся NOT_RUN, PR #8 draft. Повторный полный CI HEAD5968407 и оба отчёта PASS; понятное описание PR подготовлено.",
      "trigger": "Ответ владельца об установке среды и проверка финального кода"
    },
    {
      "trigger": "Ответ владельца о недоступности компьютера",
      "date": "2026-10-15",
      "summary": "Владелец сообщил об отсутствии физического доступа к Windows. Установка уже разрешена, но отложена до возможности подтвердить UAC. Повторные UAC не запускать без готовности владельца; перезагрузка отдельно. PR #8 опубликован с понятным комментарием; обязательный Windows smoke остаётся NOT_RUN.",
      "revision": 33
    },
    {
      "revision": 34,
      "date": "2026-10-15",
      "summary": "Завершена установка Go/WSL/Docker и реальный Windows smoke P01-06 PASS. Опубликованы воспроизводимые JSON/JSONL доказательства. По поручению владельца Docker Desktop и WSL VM выключены, автозапуск отключён, полезные кеши сохранены. PR #8 ready; P01 реализация 8/8, owner acceptance4/8, G1 NOT_READY.",
      "trigger": "Владелец вернулся к компьютеру, подтвердил UAC и поручил освобождать ресурсы личного ПК после работы"
    },
    {
      "revision": 35,
      "summary": "Согласованы текущие tracker_stage/current_packet_stage с фактическим IN_REVIEW PR #8 и заменены прежние сведения об отсутствующем host Go/Docker на проверенные результаты установки. Доказательства Windows/cleanup не менялись; владелец принял 4/8, G1 NOT_READY.",
      "date": "2026-10-15",
      "trigger": "Проверка согласованности переносимого контекста перед окончательной передачей на ревью"
    },
    {
      "revision": 36,
      "date": "2026-10-15",
      "summary": "Зафиксирован согласованный порядок PR/ADR/обзор с явным авторством замечаний из чата. По прямому поручению создана отдельная Issue #9 BACKLOG для Go smoke с критериями, без начала реализации; все четыре новых документа включены в переносимый контекст. P01 остаётся IN_REVIEW, принято 4/8, G1 NOT_READY; runtime не запускался.",
      "trigger": "Ответ владельца Хорошо на предложение фиксации инженерных решений; затем поручение оформить Go smoke отдельной backlog-задачей"
    },
    {
      "revision": 37,
      "date": "2026-10-16",
      "summary": "P01 принят 8/8 через наблюдаемый owner merge PR8, карточка3 закрыта. Предложены восемь задач P02 в Issue10/draftPR11: PostgreSQL и права, миграции/seed, Kafka/Redis, проверки и доказательства. Приоритет и группы PR ещё не приняты; P02/Go smoke9 NOT_STARTED, G1 NOT_READY, Docker не запускался.",
      "trigger": "Владелец попросил собрать ближайшие задачи, определить очередность и состав P02"
    },
    {
      "revision": 38,
      "date": "2026-10-16",
      "summary": "Владелец уточнил: P02 пока обсуждается; создание PR11 было преждевременным. PR11 закрыт без merge, никакого принятия P02 или старта реализации нет. Сохранено правило discussion-only; эта корректировка локальная, без нового commit/push.",
      "trigger": "Почему ты уже создал PR11? Не торопись, мы пока что обсуждаем планы"
    },
    {
      "revision": 39,
      "date": "2026-10-16",
      "summary": "План P02 согласован в чате и опубликован по отдельному прямому поручению: milestoneP02, восемь Issues12–19, два будущих PR (01–04/05–08), короткие тексты и критерии. Все TODO/NOT_STARTED; разработку владелец пока запретил. Issue10 закрыт not_planned, PR11 закрыт без merge. Обновлены источники и контекст; публикация документов в docs/p02-tasks, новый PR не создаётся.",
      "trigger": "публикуй на гит, разработку не начинай"
    },
    {
      "revision": 40,
      "date": "2026-10-16",
      "summary": "Вместо восьми мелких Issues P02 теперь две цельные задачи #12/#16, каждая под один будущий PR. Полные описания и критерии заменяют организационную воду и дублирование. #13–15/#17–19 закрыты not_planned, исключены из milestone, не считаются выполненными. Объём G0/P02 сохранён; реализация0/2, G1 NOT_READY, новые PR не созданы.",
      "trigger": "Владелец попросил объединить01–04 в одну задачу и05–08 в другую, с хорошими полными описаниями и качественными критериями"
    },
    {
      "revision": 41,
      "date": "2026-10-16",
      "summary": "Прочитан аудитv2 и проверены доказательства/источники. F01 актуален по архивному repro и неизменному коду; F02 исправлен вdocs-ветке, main ещё старый; F04/F05 документированные ограничения; F09–12 будущие контракты/эксперименты. План исправлений только предложен. Контекст обновлён локально, безcommit/push, GitHub-мутаций и запускаruntime.",
      "trigger": "Владелец принёс faultline-lab-agent-report-v2.zip и попросил краткую выжимку и план для последующего согласования"
    },
    {
      "revision": 42,
      "date": "2026-10-16",
      "summary": "Владелец выбрал пути F01–F12; AI-Your-Bro создал единый Issue #20 с причинами и25 критериями, связал F01 с#9 и F09 с#16. GitHub выбран единственным актуальным трекером; роли внешних агентов уточнены. Направления согласованы, точный PR и старт реализации не поручены; локальный контекст обновлён безcommit/push иruntime.",
      "trigger": "Прямое ревью F01–F12 и поручение оформить все12 задач в одном Issue"
    },
    {
      "revision": 43,
      "date": "2026-10-16",
      "summary": "Владелец поручил issue9 и F02/F04/F07/F08. Контекст schema2 хранит решения, датированные доказательства и навигацию; live task board и счётчики удалены. Внешние материалы фиксируются существующим SHA. Статус выполнения и ревью — только в GitHub.",
      "trigger": "Direct owner implementation instruction; accepted audit directions"
    },
    {
      "revision": 44,
      "date": "2026-10-16",
      "summary": "Реализованы Go smoke/F01 и F02/F04/F07/F08. Заключительный Windows smoke на чистом исходном коммите PASS; старый/новый runner функционально эквивалентны, отмена Windows проверена отдельно, cross-build пяти целей и Go quality PASS. Сохранены ограничения Mac/F03 и адресные доказательства выключения личной среды. Внешние материалы закреплены за существующим source SHA; live review/status только в GitHub.",
      "trigger": "Выполнение выбранного владельцем пакета и проверка результата"
    },
    {
      "revision": 45,
      "date": "2026-10-16",
      "summary": "По прямому поручению владельца закреплено предварительное краткое описание объёма, проверок и возможных трудностей каждой следующей задачи, с обоснованным ориентиром времени и уточнением при существенных изменениях. Правило внесено в AGENTS и WORKFLOW; новый этап согласования не вводится.",
      "trigger": "Обсуждение длительности задачи #9 и согласованное правило предварительной оценки"
    },
    {
      "revision": 46,
      "date": "2026-10-16",
      "summary": "По поручению владельца в Codex создан ежедневный анализ процесса разработки в этом чате: основная попытка07:00 МСК, резерв19:00 при отсутствии готового отчёта, не чаще одного результата за24часа. Анализ опирается на факты, оценивает самостоятельность и достаточность живого общения; запуск стенда и внедрение рекомендаций не разрешены. Сохранено описание настройки, актуальные запуски управляются в Scheduled.",
      "trigger": "Прямое поручение владельца о ежедневной ретроспективе и повторной попытке через12часов"
    },
    {
      "revision": 47,
      "date": "2026-10-16",
      "summary": "По замечаниям владельца в PR21 унифицированы названия ролей во всех исходных документах и контексте. Реальные GitHub identifiers сохранены для авторизации и доказательств. Из инструкции smoke удалены отсылки к прежнему runtime; оставлены назначение, prerequisites и команды. Включены ранее согласованные правила предварительной оценки задач и описание ежедневного анализа процесса.",
      "trigger": "PR21 review comments4295745781/4295759572 and browser note on redundant smoke instructions"
    },
    {
      "revision": 48,
      "date": "2026-10-16",
      "summary": "Внешние материалы переносимого снимка закреплены за существующим исходным коммитом после правок по ревью формулировок. Проверки генерации и документации выполнены локально; актуальные CI/review остаются в GitHub.",
      "trigger": "Обновление snapshot_ref после фиксации исходных материалов ревью"
    },
    {
      "revision": 49,
      "date": "2026-10-16",
      "summary": "Владелец сообщил о слиянии PR21; GitHub API подтвердил owner merge принятого HEAD4f283b7 в843e95e. Сохранено датированное основание приёмки #9/F01/F02/F04/F07/F08, локальная main синхронизирована. Snapshot закреплён за merge-коммитом main, доступным после squash и удаления feature-ветки. Текущее выполнение остаётся только в GitHub; новые работы и приёмка G1 не заявлены.",
      "trigger": "Прямое сообщение владельца: PR смержил; проверка GitHub merge evidence"
    },
    {
      "revision": 50,
      "date": "2026-10-16",
      "summary": "PostgreSQL foundation: isolated DB roles, serialized checksummed migrations, repeatable Inventory seed, host smoke and CI, detailed diagrams/guide; scoped run-image cleanup",
      "trigger": "Owner assigned issue12, requested detailed explanation and cleanup of obsolete Docker images"
    },
    {
      "revision": 51,
      "date": "2026-10-16",
      "summary": "Закреплены внешние материалы PostgreSQL foundation за существующим исходным коммитом 1c54d9c6aee2be4da7c38b7775b6d862402caac1; переносимый контекст синхронизирован",
      "trigger": "Immutable snapshot contract F07 after source commit"
    }
  ],
  "go_smoke_validation_2026_10_16": {
    "date": "2026-10-16",
    "source_commit": "8e49367c868af33023a6c1f7a8f94db447a828eb",
    "source_dirty": false,
    "windows_run": "609bf33e-2a05-f1d0-0f0c-2326b75f5f9f",
    "windows_result": "PASS",
    "cleanup": "PASS",
    "linux_container_quality": "PASS: format/build/lint/unit/race/govulncheck",
    "cross_compilation": [
      "windows/amd64",
      "linux/amd64",
      "linux/arm64",
      "darwin/amd64",
      "darwin/arm64"
    ],
    "mac_runtime": "NOT_RUN",
    "evidence": "docs/development/evidence/go-smoke/VALIDATION.md",
    "personal_pc_cleanup": "Docker Desktop and Docker WSL VM stopped; autostart disabled; tools/images/caches preserved",
    "live_review_and_host_CI": "GitHub PR and check runs, not cached in this snapshot"
  },
  "task_start_communication": {
    "confirmed_on": "2026-10-16",
    "source": "Direct owner instruction in chat: в следующие задачи заранее обозначь ожидаемый объем работы и возможные трудности",
    "before_each_task": [
      "Expected result and scope boundaries",
      "Main implementation and validation work",
      "Likely difficulties, dependencies and uncertainty",
      "Indicative elapsed-time range when supportable; distinguish build/CI waits"
    ],
    "update_on_material_change": true,
    "additional_start_approval_required": false
  },
  "process_review_automation": {
    "configured_on": "2026-10-16",
    "source": "Direct owner request in chat for a daily development-process analysis with an evening retry",
    "automation_id": "faultline-lab",
    "management": "Codex Scheduled; this records configuration intent, not live run status",
    "purpose": [
      "What works well and where development loses time or quality",
      "Executor autonomy and decisions needing more discussion",
      "Whether live conversation with clear responsibilities is sufficient without elaborate prompts"
    ],
    "schedule": {
      "timezone": "Europe/Moscow",
      "primary": "07:00 daily",
      "fallback": "19:00 if no eligible successful report",
      "first_primary": "2026-10-17 07:00",
      "maximum_report_frequency": "One successful report per 24-hour scheduled-slot interval; skip the next morning after evening success"
    },
    "delivery": "This project chat, Russian, normally 200–350 words with evidence and one or two improvements",
    "sources": "Available project conversation, Git, GitHub PR/Issues/reviews, CI and saved evidence; state gaps explicitly",
    "allowed_writes": "Only its own reports and period journal in ignored .local/process-analysis; no project/source/tracker changes",
    "execution_limits": "Analysis only; no Docker/WSL/tests/builds, development or external publication; no automatic implementation of proposed process improvements",
    "runtime_dependency": "Local scheduled runs require the host computer and Codex app to be running"
  },
  "role_wording_review": {
    "date": "2026-10-16",
    "source": "Owner review comments in PR21 and follow-up instruction in chat",
    "owner_term": "владелец проекта / project owner",
    "executor_term": "исполнитель / executor",
    "authority_changed": false,
    "review_comments": [
      "https://github.com/Includeoyasi/faultline-lab/pull/21#discussion_r4295745781",
      "https://github.com/Includeoyasi/faultline-lab/pull/21#discussion_r4295759572"
    ],
    "identity_exception": "Keep real GitHub logins in URLs, CODEOWNERS, account metadata, authentication checks and API evidence",
    "usage_docs": "Describe prerequisites and actions directly; migration rationale belongs in ADR and task materials. Keep role definitions in the workflow; omit repeated role paragraphs from the task-tracking guide."
  },
  "postgres_foundation": {
    "recorded_on": "2026-10-16",
    "scope": "Issue12; no business RPC or Kafka/Redis implementation",
    "postgres": "18.6-bookworm",
    "digest": "sha256:afc7e2d441324c0388fa80c3d24f733b4194a4eb7f47dd8ee2b08eb1a24a647c",
    "client": "github.com/jackc/pgx/v5 v5.11.0",
    "topology": "Separate Compose overlay; one private PostgreSQL, 3 service DBs, 6 distinct non-superuser roles; temporary dbtool",
    "migration_contract": "Per-service transaction, advisory lock before journal, version/name/SHA-256 verification; errors roll back pending SQL and journal; no transaction control in migration files",
    "seed_contract": "Three fixed Inventory SKU/catalog and stock fixtures; ON CONFLICT DO NOTHING preserves modified prices/stocks; repeat is not reset",
    "reset_contract": "Explicit faultline-data-* project plus identical confirmation; disposable smoke deletes own volume and run image; ordinary down preserves data",
    "configuration": "configs/postgres/local.json and ignored file-backed credentials; SCRAM auth, internal network, no host port; local plaintext profile only",
    "guide": "docs/development/POSTGRES.md",
    "validation": "docs/development/evidence/p02/POSTGRES.md",
    "architecture_limits": "Single DB instance shares failures/resources; no service business DB clients yet; macOS runtime is NOT_RUN",
    "personal_cleanup_refinement": "Owner observed historical image accumulation; remove run tags after success/error/cancel, including partial build; retain shared image layers/caches; no global prune",
    "local_evidence": {
      "date": "2026-10-16",
      "windows_smoke": "PASS; project faultline-data-test-90eda236ac29dee2; 19:12:11–19:13:14 MSK",
      "runtime_services_smoke": "PASS; run 5562c957-f2db-dcae-2c62-0ec73c2e8a91",
      "migration_seed_rights": "PASS on real PostgreSQL",
      "scoped_reset": "PASS; independent control volume preserved",
      "cross_compile": "PASS windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; macOS runtime NOT_RUN",
      "quality": "Container quality checks PASS; final HEAD CI records live validation",
      "resource_sample": "23.43MiB of512MiB, CPU5.12%; single sample, not capacity",
      "cleanup": "0 containers, 0 run images, Docker Desktop/WSL stopped, autostart=false; preserve2sharedimages and2Go cache volumes",
      "source": "Working-tree validation on main base843e95e, not a claim of clean final-HEAD validation",
      "ci": "Final checks and artifacts in GitHub Actions; no duplicated live status"
    }
  }
}
```

## Реестр включённых источников

SHA-256 относится к исходному файлу на момент сборки. Генерируемый PROJECT_CONTEXT.md исключён из реестра, чтобы не создавать циклическую контрольную сумму. Проверка CheckOnly сравнивает весь снимок с результатом повторной сборки из актуальных источников.

| Источник | SHA-256 |
|---|---|
| `CONTEXT_STATE.json` | `30f24e6054049c2ca0123f1fa08b9ba4dd41308246db1d7218336b62a90f6f3e` |
| `tools/context_intro.md` | `efa653b38275069b42953862f69f20c351002d21221d8c83123404a8411a56b9` |
| `tools/update_context.ps1` | `73dff2e2da1f7ea7e8fc8fe2dc70bde722f4d697c1ca7c685d297214cd133e33` |
| `README.md` | `2251b4268c6b3fd5f3c03ac37d47d01ab53f1b0f2b121f4e4e99ee3b4887d5fa` |
| `README.ru.md` | `98032f7bda43cb6cbd1f8bd3c7ada83a07e334ee7e7ed190db0d455d44dbccc7` |
| `docs/README.md` | `cee97c537a2162efc9a2da7e4035e01ffb41341170f156aa6e558309f930e451` |
| `docs/ROADMAP.md` | `3881f7898db08cd56221d60fe1578b8642060c99bc53f3574bc59e8defcff637` |
| `docs/architecture/DESIGN.md` | `74a141c2ce11d07cb568f01a6eb72b849f9556d7be85705544bb9e8a15806f2e` |
| `docs/architecture/PROCESSES.md` | `dd6d00ee9d169236fe0f0d14de4b6f54d778a7cf79489f6b6a3354357e2e9ec6` |
| `docs/development/WORKFLOW.md` | `a32d68ca2ae9085b7284101bd6d3bfc4cee2d3b50b4adefbea9c15c4fc78f97a` |
| `docs/experiments/REPORT_TEMPLATE.md` | `4f569aab52ba34a68bb2e8b4551c6a3bf4932277c26c3b4206dbb505781a2cab` |
| `AGENTS.md` | `7e601464d6bd8330e7b1df6e730e7ee7ceb4ce39a4b94731d1141e6e77e3b9a0` |
| `docs/adr/0001-approved-g0.md` | `b8518c5c49cdefea3dbf0bf2a79ab40459dba2b565e71a6c534a555cbafa20df` |
| `CONTRIBUTING.md` | `5d7fe456ad6e2967e3e324bb227f711dd23a735e4080db9ceb4d01760c94327e` |
| `LICENSE` | `26a8432e65ef8fc4a2fbb6f4170ee2c31f0cacd3d2696bb678fb0ab7fd14ad1e` |
| `docs/development/CODE_REVIEW.md` | `ede879352cb1f36c935970fae3be8886132e5080c3327c422e43ddf2da1db5d9` |
| `docs/architecture/LAB_BOUNDARY.md` | `94030566542a3e5fbda07593ed06cea7f95ec98985483de8ac0d66f4001ae326` |
| `docs/development/ENVIRONMENT.md` | `097735b23b23cc8b7606491682f17df2f7230c95d5166743bb234238ddf7df8c` |
| `docs/development/RESOURCE_PROFILES.md` | `60212954c54085f2c8b74c4d3cf5c50f8ccf90192e8e9a502da357208645d44f` |
| `docs/design/BANNER_OPTIONS.md` | `0b880a782a766f85c33fb5fba1a8c933a3d05f949d93b539d03b9a37b71efd0f` |
| `docs/development/packets/P01.md` | `82d62963046bde7289b9c22dde1a0fcc26adb6911e4bbf5bd792e068a6a22067` |
| `docs/development/TASK_TRACKING.md` | `46baf95126af2e7e362d6751e957057e915a98f4b7606739178c57427aec61d8` |
| `docs/development/REPOSITORY_SECURITY.md` | `e15d0de41447fcf81b76f4264a76483895940dc1c03f6ab809e3a626b9f8278a` |
| `.github/CODEOWNERS` | `95de6cc3de1ba640d0ab9e2b9bfbdea79b7ce6641a8bcdb0bcef1f7607666415` |
| `docs/development/TOOLCHAIN.md` | `e5325b03c3caee54ad63252f8bc24e608caf33b269c60e9129da584b9decac16` |
| `go.mod` | `32aa5667a8f832ad567ec8af3dc68f7952b1ef14ec40f8176d7e19633feba4a1` |
| `go.sum` | `fdcd4f5c8dac3b2c90899b2da40284659a1503f5ce01faf3b432576cd1708502` |
| `cmd/orders/main.go` | `131626df0bae0819d49c3d46793809c839cacb6dbd71a93b4bc1740a04ea3e6c` |
| `cmd/inventory/main.go` | `3a8d1959cf1530e56c786e7e26ac0a0d04c145e9cf5266f95f6a766e57fb5f39` |
| `cmd/payments/main.go` | `3ceb9fe70d25948079c3682ad409c4970918d30036f74afb6451858725c351c1` |
| `build/Dockerfile` | `05266b15d55807239775e1895c5456896126670f8b6def76989523e1a1edae1f` |
| `compose.yaml` | `76cabf52d3c42c963c0ed56bc486b8e68686159ad32e680f23c2574bfa2fb6ed` |
| `tools/go.sh` | `40fdab41afd4022715a929a0c9e1252037e9896f8abeb288fa286d591a0a2697` |
| `.golangci.yml` | `44651a36c889767cd425bc12b46e07c8bffc924024de23af1176057f78f3a74a` |
| `.github/workflows/go.yml` | `269203aa9d9aec1b0c6c4d4df68ee6b704cc44dffc6debccde73849c374c28bc` |
| `.dockerignore` | `e75b314d79f67eec6dafc763b1e4c118597186f106f91a73749750d4ab0cf045` |
| `docs/development/SERVICE_RUNTIME.md` | `85d5e127812af0a6ccfdd29c3fd54dbbe4d2203cfa5c3e08cc4923f9fa3fbbd6` |
| `internal/service/config.go` | `a5cd706cc41d0fb88ea96b0ffaeff4f7ecf32bf31394ec5398efb39fba240a89` |
| `internal/service/log.go` | `034f82245011a1a905514e9f21d16953dd4ce325cdfc720d8832d1d924a5773f` |
| `internal/service/run.go` | `5a1160b916ac892c1deffc38a5fda57dffb1cf352785d75cab5d960a57556a13` |
| `internal/service/config_test.go` | `c628504d6e9ad8432832b6ceb98796589ef86140dc9cf76b8dd6b2018250014f` |
| `internal/service/log_test.go` | `1b8f3fc837128930d5ef7a99d75a183bc7037035d5e756ea2faabe800d61f732` |
| `internal/service/run_test.go` | `8a70e4e7ed779fd77f3823a03ab3d1fda80edd6cd7da5ca315472815c096694d` |
| `tests/smoke/services_test.go` | `7f011324f860dc3521decfdbbf850c8724bbea8edfd98f4600ec08a16fef4ed1` |
| `configs/env/local.env.example` | `e1ed433b87891f1416c0eee2143c8c5be7131002169fb6061e86ce00fb97ba12` |
| `configs/env/container.env.example` | `968bffa2aa651eb168d4cf58fb6935ec8eeb9d636ad5e20dfb5c836be61b2c2b` |
| `.gitignore` | `f2b608989f7ce4e52f1dd1bed1015febf699d7593df371b99ee957b308a03546` |
| `internal/service/command.go` | `e571f49daf1664aa64c3c9db53a8e0b0d7fb5f2be18b6b2ba8b10cdce64d41b4` |
| `internal/service/command_test.go` | `56eb2c958131551d61097e00fd66a34068332c141b81c3eaa5ec8358072a7051` |
| `internal/service/shutdown_test.go` | `9387689f5183003faabddb0c5f6d6d5fd5c89e410fec66cf76aa94eebf25be58` |
| `docs/development/CONTAINER_RUNTIME.md` | `0807c58d99e01fb3ab6c7a4d1163b1c5965fe877e869aca7a021062025dcd13e` |
| `docs/development/evidence/p01/WINDOWS_SMOKE.md` | `8b17000ed9fb0c10083508e57a3d6da29ca36dc624e0d5d27be8c1f9ff837b53` |
| `docs/development/evidence/p01/windows-2026-10-15/report.json` | `b08d9deabadc65c765d9173482e9171782aeff28affb5e7fd9bf0b201d5add65` |
| `docs/development/evidence/p01/windows-2026-10-15/orders.jsonl` | `84ff81f5e1e0d74fc37e092972581d0243e46ef5ff6f88d7c75e8bfa4a2fe8c2` |
| `docs/development/evidence/p01/windows-2026-10-15/inventory.jsonl` | `c74c87b48339c263f785ceb169c48415448bef963eaca5e804be4ad05ff6d3f3` |
| `docs/development/evidence/p01/windows-2026-10-15/payments.jsonl` | `d3466907750bc5c137738ba0045f8d37f6c83e7fb4e3690a56896cbf9ad83c58` |
| `docs/development/evidence/p01/windows-2026-10-15/environment.json` | `e32d9a5c383c6fcc426179fc716ab8ad39bece710670215b208405ee099bc90b` |
| `docs/development/evidence/p01/windows-2026-10-15/cleanup.json` | `6c78a1ca41714a3ccb50e4777a2532ffeaca7d1cfb5fb561127e76e91fe0ad79` |
| `docs/development/ENGINEERING_DECISIONS.md` | `c553ac968846968b3c2109d7938ca80c47464b82393f82fceb2182bef66cb335` |
| `docs/adr/0002-record-review-decisions.md` | `00ad82591791da40b82d0c98c4be2c098a069fee70a5ba495a658a3657ba305e` |
| `docs/adr/0003-go-container-smoke.md` | `228110462151044ffc87a6048293b8ad12e2673bee74bb627c8675ec62347181` |
| `docs/development/backlog/GO_SMOKE.md` | `158ba7ad889fa7562db4136fc4fa475e607c0485f6d80f7bc59ab504bf80a4d4` |
| `docs/development/packets/P02.md` | `c58b3cc3a1288c61cc90ac47cbad752409230551c1de69366166d26529195364` |
| `docs/development/SMOKE.md` | `f9952d1b7033d7e2ffecec19e23f118b6c154194a86e8162d0bbc8d4aba0e043` |
| `cmd/smoke/main.go` | `16ff87a9b0691547b5d4c79663efeab797ebccee1d8beb5b48df01ea0e13b0cd` |
| `cmd/smoke/main_test.go` | `2970b74307c907e153f3eeeffb4b309c0ecc99652b4ed9bfe5f1e869c34e4154` |
| `internal/smoke/run.go` | `c190e71b9c47c9c950be8334e5302a41f63ddcb592bcfaa9186ea00b792231db` |
| `internal/smoke/report.go` | `0f16b8801b83b1d61142f180776de2e38c8830fa137d9d0bf76b6c28a96f428c` |
| `internal/smoke/checks.go` | `cb2ad73b3a258531f21c15b99f853bff819a3fc58feccad725c7be969dc32253` |
| `internal/smoke/config.go` | `ac4fffa46895f9d2108866a9614b217f72c637dc55bc7eefcdf40c8471e1f557` |
| `internal/command/exec.go` | `69aa0f4eef3cb8cef8518b528980751c5ed62104563e6778f8adeae5dbfc7d4f` |
| `internal/command/exec_unix.go` | `46f04c817063ac9b643ce92fe793059f3b285ce19fc005d68e9142733f9fae93` |
| `internal/command/exec_windows.go` | `4c54a5dcf35035f37514f92253a8af670c0c40db2a6512ba2c44ecffbe1328a5` |
| `internal/smoke/run_test.go` | `2a71a2eab90c3b36a510ed12acaf85f9adf5672bf418242d0ed66b89f304d37f` |
| `internal/smoke/config_test.go` | `770d8d3161ac428cc70c141585ab6925c7780daca683dea8b1db2b17729adf64` |
| `internal/command/exec_test.go` | `310e08549e9027c2fb6a0efb4be346567c848aabe42e0c766527ea63aff0f25e` |
| `tools/doc_safety.mjs` | `785ca466148638927a73eaeb3ca55dfbc6af65069739e5820d264dc646da1dc9` |
| `tools/doc_safety.test.mjs` | `14e2666906c021e8ae739aaa4d6c164f63b4b931f5a1726e2d42180d9b973ab7` |
| `tests/smokecli/interrupt_test.go` | `250cde8480adedd046523adc112070b8a33b6838321b30a2b4de8519da0c8f2d` |
| `internal/command/exec_unix_test.go` | `a6ab1fe31b20421193d944d1445ef15329d6d21096ba1e56e582b1788801e513` |
| `docs/development/evidence/go-smoke/go-equivalence-report.json` | `e7a68953b8101343ba23868864e02a6cb293d2a7679edcb56748fae954936af5` |
| `docs/development/evidence/go-smoke/go-smoke-equivalence.json` | `da1916fc24b5629dbb3b820797e4544a0e641b7a38536077293220b99e4bff3a` |
| `docs/development/evidence/go-smoke/go-smoke-preflight.json` | `812ecdf057b99863dee5dc0e20be998535a4fdb1d1fc81369cb8b15cb3c19da2` |
| `docs/development/evidence/go-smoke/go-smoke-windows-interrupt.json` | `f870d1471e99666768e9587429ae4b510e50b930a0397c196beae93b688ea029` |
| `docs/development/evidence/go-smoke/no-backend-report.json` | `4b5fd74d98d26f445926def471d5b76c47ffb5f306dc9196917c9a22142443a2` |
| `docs/development/evidence/go-smoke/node-baseline-report.json` | `71dfd871132333e1762454d402f30db81ac9d1bfbf2a40769831333ab5803c8b` |
| `docs/development/evidence/go-smoke/VALIDATION.md` | `80ad7a0431fe460595acb2869c8e924d6d0129c2a1f1dc6bb8f0e34c149246b9` |
| `docs/development/evidence/go-smoke/windows-cleanup.json` | `5b2ce61f636f584fab05807ba46c3602e51f2b588f730751327cb3031b8595e8` |
| `docs/development/evidence/go-smoke/windows-final-inventory.jsonl` | `40af5264a5ab207a05cc1f8c1ed296ff53ed40352f90b1431129cfbc48750ad4` |
| `docs/development/evidence/go-smoke/windows-final-orders.jsonl` | `bfa10fea28be71ae40ca822e38d28a3e79a6d898c539180383d898b09d86c07b` |
| `docs/development/evidence/go-smoke/windows-final-payments.jsonl` | `5d1c49d32c8eb80ec390e15736245290181c7a95a8178e1acc5002056494971d` |
| `docs/development/evidence/go-smoke/windows-final-report.json` | `d528361131d3eb0a2d553afd9e801de7f69caf999a7e1f6d95caffc77c1211ac` |
| `docs/development/evidence/go-smoke/windows-interrupt-report.json` | `e20bdf2620358046854ed57c8b4f831a42854c38688d9d92aa257a26feba12ae` |
| `docs/development/evidence/p02/windows-runtime-report.json` | `46f24a541852560d584ced05e0d7b8032bdf01f5c2553f1fc38251652f4fec2c` |
| `docs/development/evidence/p02/windows-cleanup.json` | `82000dcc8920ec8742fea39b1dc91b691acf056d61d335cabea759dac710db6f` |
| `docs/development/evidence/p02/windows-postgres-report.json` | `d50cfc21f7e9259c5c6e71a178d792482ab93ce8e34c8b77939b59445b0847bb` |
| `docs/development/evidence/p02/postgres-image.json` | `8298b566f5b86090978a88d7e1734fe8576e315529661b58a5f38a9b66a18c8f` |
| `docs/development/evidence/p02/POSTGRES.md` | `50a36e94eeb9ff6bb4e1a8afa5fccec7804bd5e35ef4c643aeaf58b8216a9c70` |
| `internal/command/command.go` | `3eb0aa422d5b40ba355479ca534c2f7c22ec9c118212beafbbf477ee5da28af3` |
| `internal/database/migrate.go` | `0c42150be719f33bd2f743dfc1220864f1ee215a92f4b4138c7ae549948c7373` |
| `internal/database/integration_test.go` | `d6301720392b4442d0ad87e739d2367579eaf7ca56763847c73a2b7f717c236a` |
| `internal/database/fingerprint.go` | `35544189197803be0ddd88fae70fceff8c83776ad8c31655545ecfe00eaa9ef5` |
| `internal/database/config_test.go` | `99b8d5b78e2de98e1448952f470cbce68fd7e4ab05f3747e50529828d8d732f5` |
| `internal/database/config.go` | `2e73f3f8ab4787426abe4c15cc1b74b4b9b96bb6c188a51a7458a4a52756d4c5` |
| `internal/database/bootstrap.go` | `48eefda101b01a2eaf99e44f35a0645ed9fdddc02768cc7b2b148affccb36e51` |
| `configs/postgres/local.json` | `8a935f46a670b33b9b0c66d6617f87b356506c5449cb8835cbb9f185b4b9bcc8` |
| `cmd/labdata/main_test.go` | `2d8392be252b0d370b98358cee5f57b55b114cbcd8fbf972ad78bf854d955436` |
| `cmd/labdata/main.go` | `879941426bc549ec6b4579aca23f887dde83652f5ef779162c11cab40615e8fe` |
| `migrations/inventory/0002_catalog.sql` | `cecccb789616eede86cd6cea63162264a24f023c62a7f6ddfb182e2aed387d97` |
| `migrations/inventory/0001_schema.sql` | `60a1389cfff22a0a79566c1264ed338cb6049b761398af8dc6224c3d87fe32ac` |
| `migrations/embed.go` | `439b7ae21a933585e4a5e21c8c14360bc4e2298fd6e3c39330bcfd5d17c8b258` |
| `migrations/seed/inventory.sql` | `61b9efdd980a8508fd1f36e4c4e2d0a72f907cb2c5ae0a8c9afe8c75de0deb2d` |
| `migrations/orders/0001_schema.sql` | `7a2abc0f03d1dac75ea31f86a6abfadae9532eaee42416167be45aee5d5ab16f` |
| `migrations/payments/0001_schema.sql` | `2856144143bedd41e0380818034ae85025f18b3f40baa295c4eace2d48519af4` |
| `internal/datalab/run_test.go` | `3a1019ca0ce7efe6cb871ee6278f7f5e5d4954ada079812408cda830fdcbf254` |
| `internal/datalab/run.go` | `bfeb8312d3ff80e81f6eb3805793283e4765ddd72c5e78651656d4a9a94362da` |
| `compose.postgres.yaml` | `f8e28373a62b0135d30f9b568f1b548eb4f7d8322d9296dbd909708a52a8ef4d` |
| `internal/smoke/exec_adapter.go` | `5f903f8c590d90655b333d6988139d9276ca5e6dbb400bfff10b0af0dda0c4ba` |
| `docs/development/POSTGRES.md` | `7d9c8f56acf84056287efd527c6e1a1dbbb47343a75202d6d15d100257799041` |

## Полные самостоятельные приложения

- [Публичное введение на английском](#appendix-readme-en) — `README.md`
- [Публичное введение на русском](#appendix-readme-ru) — `README.ru.md`
- [Карта документации](#appendix-docs-index) — `docs/README.md`
- [Полный roadmap, архитектура, эксперименты и методика](#appendix-roadmap) — `docs/ROADMAP.md`
- [Принятый архитектурный пакет G0](#appendix-design-review) — `docs/architecture/DESIGN.md`
- [Подробные визуальные схемы G0](#appendix-g0-processes) — `docs/architecture/PROCESSES.md`
- [Роли, ревью, контрольные точки и обновление контекста](#appendix-development-process) — `docs/development/WORKFLOW.md`
- [Полный шаблон отчёта](#appendix-report-template) — `docs/experiments/REPORT_TEMPLATE.md`
- [Локальные инструкции ИИ-исполнителю](#appendix-agents) — `AGENTS.md`
- [Запись принятия архитектуры G0](#appendix-adr-0001) — `docs/adr/0001-approved-g0.md`
- [Участие, ревью и проверки](#appendix-contributing) — `CONTRIBUTING.md`
- [MIT License](#appendix-license) — `LICENSE`
- [Правила Go code review и критерии принятия изменений](#appendix-code-review) — `docs/development/CODE_REVIEW.md`
- [Общая лаборатория, адаптеры и кандидат второго бизнеса](#appendix-lab-boundary) — `docs/architecture/LAB_BOUNDARY.md`
- [Переносимая Docker-среда и распределение контейнеров](#appendix-environment) — `docs/development/ENVIRONMENT.md`
- [Ресурсный снимок Windows и будущие профили macOS/сервера](#appendix-resource-profiles) — `docs/development/RESOURCE_PROFILES.md`
- [Концепции баннера и выбранное оформление](#appendix-banner-options) — `docs/design/BANNER_OPTIONS.md`
- [Первый PR с кодом: задачи, зависимости и критерии приёмки](#appendix-p01) — `docs/development/packets/P01.md`
- [GitHub Issues, статусы и переносимые снимки](#appendix-task-tracking) — `docs/development/TASK_TRACKING.md`
- [Защита main, отдельный исполнитель и состояние подключения](#appendix-repository-security) — `docs/development/REPOSITORY_SECURITY.md`
- [Действующий CODEOWNERS для человеческого ревью](#appendix-codeowners) — `.github/CODEOWNERS`
- [Контейнерный toolchain и команды](#appendix-toolchain) — `docs/development/TOOLCHAIN.md`
- [Go-модуль и версии инструментов](#appendix-go-mod) — `go.mod`
- [Контрольные суммы модулей](#appendix-go-sum) — `go.sum`
- [Orders gRPC health entrypoint](#appendix-orders-main) — `cmd/orders/main.go`
- [Inventory gRPC health entrypoint](#appendix-inventory-main) — `cmd/inventory/main.go`
- [Payments gRPC health entrypoint](#appendix-payments-main) — `cmd/payments/main.go`
- [Единый multi-stage Dockerfile: toolchain и runtime](#appendix-toolchain-dockerfile) — `build/Dockerfile`
- [Compose для временных инструментов](#appendix-compose-toolchain) — `compose.yaml`
- [Команды внутри Linux-контейнера](#appendix-go-commands) — `tools/go.sh`
- [Согласованные Go-линтеры](#appendix-golangci-config) — `.golangci.yml`
- [Проверки Go в нативном Linux CI](#appendix-go-workflow) — `.github/workflows/go.yml`
- [Исключения Docker build context](#appendix-dockerignore) — `.dockerignore`
- [Конфигурация, JSON-логи и gRPC health: поведение и проверки](#appendix-service-runtime) — `docs/development/SERVICE_RUNTIME.md`
- [Валидация конфигурации](#appendix-service-config) — `internal/service/config.go`
- [JSON-логи и build metadata](#appendix-service-logging) — `internal/service/log.go`
- [Запуск gRPC health и немедленная остановка](#appendix-service-run) — `internal/service/run.go`
- [Проверки конфигурации](#appendix-service-config-tests) — `internal/service/config_test.go`
- [Проверки JSON и build metadata](#appendix-service-log-tests) — `internal/service/log_test.go`
- [Проверки bind, TCP health и отмены контекста](#appendix-service-run-tests) — `internal/service/run_test.go`
- [Проверки трёх реальных процессов](#appendix-service-process-tests) — `tests/smoke/services_test.go`
- [Явное локальное окружение: образец](#appendix-local-env-example) — `configs/env/local.env.example`
- [Явное контейнерное окружение: образец](#appendix-container-env-example) — `configs/env/container.env.example`
- [Исключение личных настроек и артефактов из Git](#appendix-gitignore) — `.gitignore`
- [Entrypoint command and bounded health probe](#appendix-service-command) — `internal/service/command.go`
- [Health probe and command tests](#appendix-service-command-tests) — `internal/service/command_test.go`
- [Controlled active RPC and deadline tests](#appendix-shutdown-tests) — `internal/service/shutdown_test.go`
- [Runtime build, health, shutdown, smoke and review limitations](#appendix-container-runtime) — `docs/development/CONTAINER_RUNTIME.md`
- [Windows smoke: итог и ограничения](#windows-smoke-proof) — `docs/development/evidence/p01/WINDOWS_SMOKE.md`
- [Полный Windows smoke report](#windows-smoke-report) — `docs/development/evidence/p01/windows-2026-10-15/report.json`
- [Windows: Orders JSONL](#windows-smoke-orders) — `docs/development/evidence/p01/windows-2026-10-15/orders.jsonl`
- [Windows: Inventory JSONL](#windows-smoke-inventory) — `docs/development/evidence/p01/windows-2026-10-15/inventory.jsonl`
- [Windows: Payments JSONL](#windows-smoke-payments) — `docs/development/evidence/p01/windows-2026-10-15/payments.jsonl`
- [Windows: фактический бюджет](#windows-smoke-environment) — `docs/development/evidence/p01/windows-2026-10-15/environment.json`
- [Windows: подтверждение остановки](#windows-smoke-cleanup) — `docs/development/evidence/p01/windows-2026-10-15/cleanup.json`
- [Инженерные решения, вклад владельца и доказательства ревью](#appendix-engineering-decisions) — `docs/development/ENGINEERING_DECISIONS.md`
- [ADR 0002: согласованный порядок фиксации решений](#appendix-adr-0002) — `docs/adr/0002-record-review-decisions.md`
- [ADR 0003: proposed Go container smoke](#appendix-adr-0003) — `docs/adr/0003-go-container-smoke.md`
- [Backlog #9: перенос smoke на Go, объём и критерии](#appendix-go-smoke-backlog) — `docs/development/backlog/GO_SMOKE.md`
- [Предложенный P02: задачи, зависимости, критерии и порядок PR](#appendix-p02) — `docs/development/packets/P02.md`
- [Go container smoke guide](#go-smoke-guide) — `docs/development/SMOKE.md`
- [Go smoke entrypoint](#go-smoke-main) — `cmd/smoke/main.go`
- [Smoke CLI tests](#go-smoke-main-tests) — `cmd/smoke/main_test.go`
- [Smoke lifecycle](#go-smoke-run) — `internal/smoke/run.go`
- [Smoke report](#go-smoke-report) — `internal/smoke/report.go`
- [Smoke container checks](#go-smoke-checks) — `internal/smoke/checks.go`
- [Compose shutdown preflight](#go-smoke-config) — `internal/smoke/config.go`
- [Bounded command output](#go-smoke-exec) — `internal/command/exec.go`
- [Unix command process groups](#go-smoke-unix) — `internal/command/exec_unix.go`
- [Windows command jobs](#go-smoke-windows) — `internal/command/exec_windows.go`
- [Smoke failure and cancellation tests](#go-smoke-run-tests) — `internal/smoke/run_test.go`
- [Shutdown budget tests](#go-smoke-config-tests) — `internal/smoke/config_test.go`
- [Command tree cancellation tests](#go-smoke-exec-tests) — `internal/command/exec_test.go`
- [Documentation path hygiene](#doc-path-hygiene) — `tools/doc_safety.mjs`
- [Path hygiene examples](#doc-path-hygiene-tests) — `tools/doc_safety.test.mjs`
- [Host CLI SIGINT regression](#go-smoke-cli-interrupt) — `tests/smokecli/interrupt_test.go`
- [Unix orphan plugin lifecycle regression](#source-smoke-exec-unix-test) — `internal/command/exec_unix_test.go`
- [Go smoke evidence: go-equivalence-report.json](#go-smoke-evidence-go-equivalence-report) — `docs/development/evidence/go-smoke/go-equivalence-report.json`
- [Go smoke evidence: go-smoke-equivalence.json](#go-smoke-evidence-go-smoke-equivalence) — `docs/development/evidence/go-smoke/go-smoke-equivalence.json`
- [Go smoke evidence: go-smoke-preflight.json](#go-smoke-evidence-go-smoke-preflight) — `docs/development/evidence/go-smoke/go-smoke-preflight.json`
- [Go smoke evidence: go-smoke-windows-interrupt.json](#go-smoke-evidence-go-smoke-windows-interrupt) — `docs/development/evidence/go-smoke/go-smoke-windows-interrupt.json`
- [Go smoke evidence: no-backend-report.json](#go-smoke-evidence-no-backend-report) — `docs/development/evidence/go-smoke/no-backend-report.json`
- [Go smoke evidence: node-baseline-report.json](#go-smoke-evidence-node-baseline-report) — `docs/development/evidence/go-smoke/node-baseline-report.json`
- [Go smoke evidence: VALIDATION.md](#go-smoke-evidence-VALIDATION) — `docs/development/evidence/go-smoke/VALIDATION.md`
- [Go smoke evidence: windows-cleanup.json](#go-smoke-evidence-windows-cleanup) — `docs/development/evidence/go-smoke/windows-cleanup.json`
- [Go smoke evidence: windows-final-inventory.jsonl](#go-smoke-evidence-windows-final-inventory) — `docs/development/evidence/go-smoke/windows-final-inventory.jsonl`
- [Go smoke evidence: windows-final-orders.jsonl](#go-smoke-evidence-windows-final-orders) — `docs/development/evidence/go-smoke/windows-final-orders.jsonl`
- [Go smoke evidence: windows-final-payments.jsonl](#go-smoke-evidence-windows-final-payments) — `docs/development/evidence/go-smoke/windows-final-payments.jsonl`
- [Go smoke evidence: windows-final-report.json](#go-smoke-evidence-windows-final-report) — `docs/development/evidence/go-smoke/windows-final-report.json`
- [Go smoke evidence: windows-interrupt-report.json](#go-smoke-evidence-windows-interrupt-report) — `docs/development/evidence/go-smoke/windows-interrupt-report.json`
- [docs/development/evidence/p02/windows-runtime-report.json](#p02-docs-development-evidence-p02-windows-runtime-report-json) — `docs/development/evidence/p02/windows-runtime-report.json`
- [docs/development/evidence/p02/windows-cleanup.json](#p02-docs-development-evidence-p02-windows-cleanup-json) — `docs/development/evidence/p02/windows-cleanup.json`
- [docs/development/evidence/p02/windows-postgres-report.json](#p02-docs-development-evidence-p02-windows-postgres-report-json) — `docs/development/evidence/p02/windows-postgres-report.json`
- [docs/development/evidence/p02/postgres-image.json](#p02-docs-development-evidence-p02-postgres-image-json) — `docs/development/evidence/p02/postgres-image.json`
- [docs/development/evidence/p02/POSTGRES.md](#p02-docs-development-evidence-p02-postgres-md) — `docs/development/evidence/p02/POSTGRES.md`
- [internal/command/command.go](#p02-internal-command-command-go) — `internal/command/command.go`
- [internal/database/migrate.go](#p02-internal-database-migrate-go) — `internal/database/migrate.go`
- [internal/database/integration_test.go](#p02-internal-database-integration-test-go) — `internal/database/integration_test.go`
- [internal/database/fingerprint.go](#p02-internal-database-fingerprint-go) — `internal/database/fingerprint.go`
- [internal/database/config_test.go](#p02-internal-database-config-test-go) — `internal/database/config_test.go`
- [internal/database/config.go](#p02-internal-database-config-go) — `internal/database/config.go`
- [internal/database/bootstrap.go](#p02-internal-database-bootstrap-go) — `internal/database/bootstrap.go`
- [configs/postgres/local.json](#p02-configs-postgres-local-json) — `configs/postgres/local.json`
- [cmd/labdata/main_test.go](#p02-cmd-labdata-main-test-go) — `cmd/labdata/main_test.go`
- [cmd/labdata/main.go](#p02-cmd-labdata-main-go) — `cmd/labdata/main.go`
- [migrations/inventory/0002_catalog.sql](#p02-migrations-inventory-0002-catalog-sql) — `migrations/inventory/0002_catalog.sql`
- [migrations/inventory/0001_schema.sql](#p02-migrations-inventory-0001-schema-sql) — `migrations/inventory/0001_schema.sql`
- [migrations/embed.go](#p02-migrations-embed-go) — `migrations/embed.go`
- [migrations/seed/inventory.sql](#p02-migrations-seed-inventory-sql) — `migrations/seed/inventory.sql`
- [migrations/orders/0001_schema.sql](#p02-migrations-orders-0001-schema-sql) — `migrations/orders/0001_schema.sql`
- [migrations/payments/0001_schema.sql](#p02-migrations-payments-0001-schema-sql) — `migrations/payments/0001_schema.sql`
- [internal/datalab/run_test.go](#p02-internal-datalab-run-test-go) — `internal/datalab/run_test.go`
- [internal/datalab/run.go](#p02-internal-datalab-run-go) — `internal/datalab/run.go`
- [compose.postgres.yaml](#p02-compose-postgres-yaml) — `compose.postgres.yaml`
- [internal/smoke/exec_adapter.go](#p02-internal-smoke-exec-adapter-go) — `internal/smoke/exec_adapter.go`
- [docs/development/POSTGRES.md](#p02-docs-development-postgres-md) — `docs/development/POSTGRES.md`

<a id="appendix-readme-en"></a>

## Приложение: Публичное введение на английском

Полный исходный документ: `README.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

<p align="center">
  <img src="https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/docs/assets/banner.svg" alt="Faultline Lab — a research lab for hypotheses, experiments and evidence" width="100%">
</p>

<p align="center">
  <strong>A small business system. A deliberate breaking point. Evidence for every fix.</strong><br>
  A Go microservices lab for studying scale, failure and recovery through human-led AI development.
</p>

<p align="center">
  <a href="https://github.com/Includeoyasi/faultline-lab/actions/workflows/docs.yml"><img src="https://github.com/Includeoyasi/faultline-lab/actions/workflows/docs.yml/badge.svg" alt="Documentation checks"></a>
  · <a href="#appendix-license">MIT License</a>
</p>

<p align="center">
  <a href="#appendix-readme-ru">Русский</a> ·
  <a href="#appendix-docs-index">Documentation</a> ·
  <a href="#appendix-g0-processes">Process atlas</a> ·
  <a href="#appendix-roadmap">Roadmap</a> ·
  <a href="https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/PROJECT_CONTEXT.md">Agent handoff</a>
</p>

#### Project guide

The repository contains three Go service processes with explicit configuration, gRPC health, JSON logs and bounded shutdown. The shop workflow, infrastructure and experiments follow the [roadmap](#appendix-roadmap); the architecture diagram below describes their intended boundaries.

[Build tools](#appendix-toolchain) · [Runtime guide](#appendix-container-runtime) · [Go container smoke](#go-smoke-guide).

**Live work and review status:** [Issues](https://github.com/Includeoyasi/faultline-lab/issues), [pull requests](https://github.com/Includeoyasi/faultline-lab/pulls) and [milestones](https://github.com/Includeoyasi/faultline-lab/milestones). The repository keeps specifications and dated evidence, including the [Windows baseline](#windows-smoke-proof), rather than a second progress board. Performance claims require their own measurements.

The owner chose container-based development for Windows/macOS and a reusable lab boundary with business adapters. The [environment plan](#appendix-environment) and [boundary design](#appendix-lab-boundary) explain those decisions.

#### What this project investigates

Faultline Lab models a mini-shop: **order → reservation → simulated payment → completion or cancellation**. Around that business process, the planned lab will generate load, inject faults, verify invariants and compare engineering decisions.

The useful result of an experiment is more than a latency chart: it explains whether the business stayed correct, what failed, which change helped, how much it cost and whether the result can be reproduced.

| Business service | Owns | Critical responsibility |
| :--- | :--- | :--- |
| **Orders** | Order state and durable workflow jobs | Coordinate completion and compensation; account for accepted requests |
| **Inventory** | Stock, reservations and stock ledger | Prevent overselling; commit or release a reservation safely |
| **Payments** | Simulated provider, payment jobs and payment ledger | Prevent duplicate debits; cancel or refund safely |

Exactly three business services. Workflow coordination lives inside Orders; load generators and verifiers are lab tools.

#### Architecture

The diagram shows the approved boundaries. The full transaction sequence, cancellation races and recovery paths are documented in the [process atlas](#appendix-g0-processes).

```mermaid
flowchart LR
    Client["Client / load generator"] -->|gRPC| Orders["Orders<br/>Workflow + durable jobs"]
    Orders -->|"gRPC: reserve / commit / release"| Inventory["Inventory<br/>Stock + reservation ledger"]
    Orders -->|"gRPC: query / cancel / refund"| Payments["Payments<br/>Simulated provider + ledger"]
    Orders -->|outbox| Kafka["Kafka<br/>Domain events"]
    Inventory -->|outbox| Kafka
    Payments -->|outbox| Kafka
    Kafka -->|"payment requests"| Payments
    Kafka -->|"payment + inventory results"| Orders
    Orders --- ODB[("Orders database")]
    Inventory --- IDB[("Inventory database")]
    Payments --- PDB[("Payments database")]
    Inventory -.->|"availability cache"| Redis[(Redis)]
```

**Planned stack:** Go · gRPC · PostgreSQL · Kafka · Redis · Prometheus · Grafana · Loki · Alloy. Local Docker Compose comes first; distributed tracing is a later extension. A single local PostgreSQL instance hosts three databases with separate roles. Services never read one another's tables. Redis caches availability; PostgreSQL authorizes reservations.

The accepted design makes recovery explicit:

- Business state and its outbox event commit together. Duplicate delivery is expected; inbox deduplication and operation keys preserve meaning.
- An RPC timeout means the outcome is unknown. Recovery reuses the operation key or reconciles authoritative state.
- `COMPLETED` requires confirmed payment and committed inventory. `CANCELLED` requires no active reservation and no unrefunded debit.
- Cancellation and finalization compete for one versioned transition. Persisted cancellation barriers block late reservation or debit.
- Consumer lag and unfinished business work are measured separately. A committed Kafka offset does not prove a workflow has finished.

#### The experiment loop

```mermaid
flowchart LR
    H["Hypothesis"] --> B["Baseline<br/>fixed seed + resource budget"]
    B --> F["Load / fault scenario"]
    F --> V["Drain + business verifier"]
    V --> C["Change one decision"]
    C --> R["Repeat under equal conditions"]
    R --> E["Report<br/>correctness + performance + cost"]
    E -.-> H
```

The [experiment catalog](#appendix-roadmap) includes load ramps, bursts, hot products, database contention, Kafka backlog and partition skew, slow dependencies, retry storms, crash windows, Redis failures, late payments and poison events.

Reports retain offered versus actual load, unresolved outcomes, resource limits and repeat variability. The [report template](#appendix-report-template) defines the evidence to collect. There are no completed experiment reports yet.

#### Human-led AI development

The **project owner** directs development and reviews changes. The **executor, AI-Your-Bro (Codex)**, prepares focused changes, runs available checks, documents evidence and addresses review comments. The project owner approves architecture and milestone transitions.

```mermaid
flowchart LR
    Task["Scoped task"] --> AI["AI implementation<br/>+ checks + evidence"]
    AI --> Review["Human code review"]
    Review -->|changes requested| AI
    Review -->|accepted| Context["Update portable context"]
    Context --> Next["Next task / approved gate"]
```

| Gate | Result |
| :--- | :--- |
| **G0** | Architecture, invariants and development plan |
| **G1** | Orders and Inventory foundation |
| **G2** | End-to-end MVP and first experiment |
| **G3** | Failure recovery and resilience |
| **G4** | Reproducible engineering lab |

The [code review standard](#appendix-code-review) covers Go correctness, concurrency, data integrity and evidence. Review expectations and milestone decisions are in the [development workflow](#appendix-development-process). [AGENTS.md](#appendix-agents) gives implementers their working rules; [PROJECT_CONTEXT.md](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/PROJECT_CONTEXT.md) embeds specifications, decisions and dated evidence for users and testers without repository access.

[Engineering decisions and review](#appendix-engineering-decisions) traces owner feedback through the rationale, changes and verification. It also links the separate [Go smoke task](https://github.com/Includeoyasi/faultline-lab/issues/9); its execution and review are tracked in GitHub.

#### Repository map

```text
faultline-lab/
├── README.md                    # Public overview in English
├── README.ru.md                 # Public overview in Russian
├── CONTRIBUTING.md              # Review and contribution expectations
├── AGENTS.md                    # Instructions for AI implementers
├── CONTEXT_STATE.json           # Status, approvals and validation evidence
├── PROJECT_CONTEXT.md           # Generated, self-contained agent handoff
├── docs/
│   ├── README.md                # Guided reading paths
│   ├── ROADMAP.md               # Stages, experiments and acceptance criteria
│   ├── architecture/           # Approved design and process atlas
│   ├── development/            # Human review and AI implementation workflow
│   ├── experiments/            # Experiment report template
│   ├── adr/                    # Recorded decisions and approvals
│   └── assets/                 # Repository visuals
├── cmd/                        # Business service and host-tool entrypoints
├── internal/service/           # Shared configuration, logging and gRPC startup
├── configs/                    # Explicit service and database configuration
├── migrations/                 # Versioned SQL and repeatable Inventory seed
├── tests/smoke/                # Real process health and isolation checks
├── build/                      # Pinned toolchain image
├── compose.yaml                # Three services and temporary toolchain
├── compose.postgres.yaml       # PostgreSQL and temporary dbtool overlay
├── tools/                      # Go commands, context and documentation checks
└── .github/                    # Documentation CI and review templates
```

Compose now defines three runtime services and the temporary toolchain profile; see the [runtime guide](#appendix-container-runtime). The [PostgreSQL guide](#p02-docs-development-postgres-md) explains roles, migrations, seed and reset with interaction diagrams. The [roadmap](#appendix-roadmap) contains the planned code layout.

#### Working with this snapshot

Clone the repository and start with the [documentation index](#appendix-docs-index):

```bash
git clone https://github.com/Includeoyasi/faultline-lab.git
cd faultline-lab
```

Documentation checks require Node.js 20+ and PowerShell 7+; the context generator also runs on Windows PowerShell 5.1:

```powershell
node tools/check_docs.mjs
pwsh -NoProfile -File tools/update_context.ps1 -CheckOnly
```

After changing canonical documents, update `CONTEXT_STATE.json` and its revision, regenerate with `pwsh -NoProfile -File tools/update_context.ps1`, then run the checks again. These commands validate documentation and snapshot freshness. Service settings and process checks are described in the [runtime guide](#appendix-service-runtime).

#### License

[MIT](#appendix-license) · Copyright © 2026 Dmitrii Morozov.

<a id="appendix-readme-ru"></a>

## Приложение: Публичное введение на русском

Полный исходный документ: `README.ru.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

<p align="center">
  <img src="https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/docs/assets/banner.svg" alt="Faultline Lab — исследовательская лаборатория: гипотеза, эксперимент, доказательства" width="100%">
</p>

<p align="center">
  <strong>Небольшой бизнес. Управляемые сбои. Доказательства эффективности решений.</strong><br>
  Лаборатория микросервисов на Go, разрабатываемая ИИ под руководством человека.
</p>

<p align="center">
  <a href="https://github.com/Includeoyasi/faultline-lab/actions/workflows/docs.yml"><img src="https://github.com/Includeoyasi/faultline-lab/actions/workflows/docs.yml/badge.svg" alt="Проверки документации"></a>
  · <a href="#appendix-license">MIT License</a>
</p>

<p align="center">
  <a href="#appendix-readme-en">English</a> ·
  <a href="#appendix-docs-index">Документация</a> ·
  <a href="#appendix-g0-processes">Атлас процессов</a> ·
  <a href="#appendix-roadmap">Roadmap</a> ·
  <a href="https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/PROJECT_CONTEXT.md">Контекст для агентов</a>
</p>

#### Навигация по проекту

Репозиторий содержит три Go-процесса с явной конфигурацией, gRPC health, JSON-логами и ограниченным завершением. Бизнес-процесс, инфраструктура и эксперименты вводятся по [roadmap](#appendix-roadmap); схема ниже показывает принятые границы.

[Команды сборки](#appendix-toolchain) · [Запуск контейнеров](#appendix-container-runtime) · [Go smoke](#go-smoke-guide).

**Актуальные задачи и ревью:** [Issues](https://github.com/Includeoyasi/faultline-lab/issues), [PR](https://github.com/Includeoyasi/faultline-lab/pulls) и [milestones](https://github.com/Includeoyasi/faultline-lab/milestones). Здесь хранятся спецификации и датированные доказательства, включая [Windows baseline](#windows-smoke-proof). Выводы об эффективности требуют отдельных измерений.

#### Идея

Faultline Lab — мини-магазин и лаборатория инженерных экспериментов вокруг него. Бизнес-процесс: **заказ → резервирование товара → симулированная оплата → завершение или отмена**.

Мы будем воспроизводить перегрузку, сбои зависимостей, задержки событий и ошибки восстановления; затем менять инженерное решение и сравнивать результаты в одинаковых условиях. Отчёт должен отвечать, сохранилась ли корректность бизнеса, что изменилось в производительности и ресурсах и можно ли повторить результат.

| Сервис | Ответственность | Что защищает |
| :--- | :--- | :--- |
| **Orders** | Заказ, сохранённые шаги процесса, отмена и компенсации | Учёт принятых запросов и корректные конечные состояния |
| **Inventory** | Остатки, резервы и складской ledger | Отсутствие перепродажи и безопасное освобождение резерва |
| **Payments** | Симуляция провайдера, платежи и возвраты | Отсутствие повторного списания и поздней оплаты после отмены |

Ровно три бизнес-сервиса. Управление процессом остаётся в Orders; генератор нагрузки и verifier — инструменты лаборатории.

#### Принятая архитектура

```mermaid
flowchart LR
    Client["Клиент / генератор нагрузки"] -->|gRPC| Orders["Orders<br/>Заказ + durable jobs"]
    Orders -->|"gRPC: резерв / продажа / освобождение"| Inventory["Inventory<br/>Остатки + ledger"]
    Orders -->|"gRPC: запрос / отмена / возврат"| Payments["Payments<br/>Симулятор + ledger"]
    Orders -->|outbox| Kafka["Kafka<br/>Бизнес-события"]
    Inventory -->|outbox| Kafka
    Payments -->|outbox| Kafka
    Kafka -->|"запрос оплаты"| Payments
    Kafka -->|"результаты оплаты и склада"| Orders
    Inventory -.->|"кэш доступности"| Redis[(Redis)]
```

Планируемый стек: **Go, gRPC, PostgreSQL, Kafka, Redis, Prometheus, Grafana, Loki и Alloy**. Первая среда — локальный Docker Compose. У каждого сервиса своя база и роль PostgreSQL; чтение чужих таблиц исключено. Redis помогает читать доступность, а резерв разрешает PostgreSQL.

Принятые опорные правила:

- Бизнес-изменение и outbox-событие сохраняются одной транзакцией. Дубликаты доставки учитываются через inbox и ключи операций.
- Timeout RPC означает неизвестный результат. Восстановление повторяет тот же ключ или сверяет состояние владельца данных.
- `COMPLETED` требует подтверждённой оплаты и продажи товара. `CANCELLED` требует отсутствия активного резерва и невозвращённого списания.
- Отмена и финализация конкурируют за переход версии Orders. Сохранённые запреты отмены блокируют поздний резерв или списание.
- Kafka lag измеряется отдельно от незавершённых бизнес-задач.

В [атласе G0](#appendix-g0-processes) разобраны локальные транзакции, успешная покупка, отказы, гонки, поздние события и восстановление.

#### Цикл эксперимента

```mermaid
flowchart LR
    H["Гипотеза"] --> B["Базовый прогон<br/>seed + лимиты ресурсов"]
    B --> F["Нагрузка / сбой"]
    F --> V["Drain + проверка бизнес-инвариантов"]
    V --> C["Одно изменение решения"]
    C --> R["Повтор в тех же условиях"]
    R --> E["Отчёт<br/>корректность + скорость + ресурсы"]
    E -.-> H
```

В каталоге roadmap — постепенная нагрузка, всплески, горячий SKU, блокировки PostgreSQL, backlog и перекос партиций Kafka, медленные зависимости, лавина повторов, рестарты, сбои Redis, поздняя оплата и poison events.

Мы учитываем предложенную и фактическую нагрузку, неизвестные исходы и разброс повторов. Цифры эффективности появятся после измерений. [Шаблон отчёта](#appendix-report-template) определяет состав доказательств; готовых результатов пока нет.

#### Человек и ИИ

**Владелец проекта** задаёт направление разработки, проводит код-ревью, принимает архитектуру и ключевые этапы. **Исполнитель — AI-Your-Bro (Codex)**: готовит изменения, выполняет проверки, фиксирует доказательства и исправляет замечания.

```mermaid
flowchart LR
    T["Конкретная задача"] --> AI["ИИ: реализация<br/>проверки и доказательства"]
    AI --> R["Ревью владельца проекта"]
    R -->|замечания| AI
    R -->|принято| C["Обновление контекста"]
    C --> N["Следующая задача / этап"]
```

G0 — архитектура; G1 — основа заказа и склада; G2 — сквозной MVP; G3 — устойчивость; G4 — воспроизводимая лаборатория. Порядок и критерии приёмки описаны в workflow; факты принятия сохраняются с датой и источником.

[Правила код-ревью](#appendix-code-review) задают стандарт Go, проверки корректности и критерии слияния. [Порядок разработки](#appendix-development-process) описывает процесс и приёмку этапов. [AGENTS.md](#appendix-agents) задаёт правила исполнителю. [PROJECT_CONTEXT.md](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/PROJECT_CONTEXT.md) включает полные документы, решения и датированные доказательства: его можно передать агенту без доступа к репозиторию.

[Инженерные решения и ревью](#appendix-engineering-decisions) связывают замечания владельца с обоснованием, изменениями и проверками. Выполнение переноса контейнерного smoke на Go и его ревью ведутся в [задаче #9](https://github.com/Includeoyasi/faultline-lab/issues/9).

#### С чего начать

1. [Карта документации](#appendix-docs-index) — выбрать маршрут чтения.
2. [Принятый дизайн](#appendix-design-review) — границы и инварианты.
3. [Процессы G0](#appendix-g0-processes) — подробные визуальные сценарии.
4. [Roadmap](#appendix-roadmap) — этапы, критерии и эксперименты.
5. [CONTRIBUTING.md](#appendix-contributing) — как подготовить проверяемое изменение.

Структура репозитория показана в [English README](#appendix-readme-en). Каталоги реализации будут добавляться вместе с работающим содержимым.

```bash
git clone https://github.com/Includeoyasi/faultline-lab.git
cd faultline-lab
```

Проверки документации: Node.js 20+ и PowerShell 7+. Генератор контекста также совместим с Windows PowerShell 5.1:

```powershell
node tools/check_docs.mjs
pwsh -NoProfile -File tools/update_context.ps1 -CheckOnly
```

После изменения источников обнови `CONTEXT_STATE.json`, увеличь ревизию и выполни `pwsh -NoProfile -File tools/update_context.ps1`, затем проверки. Настройки сервисов и process smoke описаны в [инструкции запуска](#appendix-service-runtime).

Подготовка баз, права, миграции, повторяемый seed и явный сброс описаны в [инструкции PostgreSQL со схемами](#p02-docs-development-postgres-md).

#### Лицензия

[MIT](#appendix-license) · Copyright © 2026 Dmitrii Morozov.

<a id="appendix-docs-index"></a>

## Приложение: Карта документации

Полный исходный документ: `docs/README.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Карта документации

**Faultline Lab** · спецификации, решения и доказательства. Актуальное выполнение — в [GitHub Issues](https://github.com/Includeoyasi/faultline-lab/issues) и PR.

#### Маршруты чтения

| Задача читателя | Маршрут |
| :--- | :--- |
| Быстро оценить проект | [README](#appendix-readme-ru) → [Дизайн](#appendix-design-review) → [Roadmap](#appendix-roadmap) |
| Проверить распределённый процесс | [Атлас G0](#appendix-g0-processes) → [Принятие G0](#appendix-adr-0001) |
| Разобрать подключение другого бизнеса | [Граница лаборатории](#appendix-lab-boundary) — согласованное направление и детализация для ревью |
| Подготовить среду Windows/macOS | [Контейнерная среда](#appendix-environment) — состав, профили и статус проверок |
| Собрать и проверить первые Go-бинарники | [Toolchain P01-01](#appendix-toolchain) — версии, Docker-команды и границы проверок |
| Запустить и проверить три процесса | [Запуск сервисов](#appendix-service-runtime) — настройки, JSON, health и границы shutdown |
| Поднять три runtime-контейнера | [Контейнерный каркас](#appendix-container-runtime) — сборка, health, остановка и smoke с отчётом |
| Подготовить базы и исходные данные | [PostgreSQL](#p02-docs-development-postgres-md) — роли, миграции, seed, схемы взаимодействий и reset |
| Проверить Windows smoke | [Отчёт P01](#windows-smoke-proof) — исходные JSON/JSONL, бюджет и остановка среды |
| Проверить исходные ресурсы | [Профили сред](#appendix-resource-profiles) — Windows, macOS и возможный будущий сервер |
| Посмотреть оформление | [Концепции баннера](#appendix-banner-options) — три варианта и выбранное решение |
| Принять участие в разработке | [CONTRIBUTING](#appendix-contributing) → [Правила ревью](#appendix-code-review) → [Workflow](#appendix-development-process) → [AGENTS](#appendix-agents) |
| Разобрать принятый каркас | [P01](#appendix-p01) — границы, критерии и датированная приёмка |
| Выбрать следующий объём | [P02](#appendix-p02) — две задачи, порядок, границы PR и критерии |
| Следить за выполнением | [Учёт задач](#appendix-task-tracking) — GitHub Issues/PR как единственный трекер |
| Проследить инженерные решения | [Решения и ревью](#appendix-engineering-decisions) — вопросы владельца, обоснования, изменения и доказательства |
| Проверить отдельный backlog | [Перенос smoke на Go](#appendix-go-smoke-backlog) — задача #9, границы и критерии |
| Проверить права и защиту main | [Безопасность репозитория](#appendix-repository-security) — действующая защита и переход на аккаунт исполнителя |
| Подготовить эксперимент | [Roadmap](#appendix-roadmap) → [Шаблон отчёта](#appendix-report-template) |
| Ввести другого ИИ в задачу | [Переносимый контекст](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/PROJECT_CONTEXT.md) — полный самостоятельный снимок |

#### Источники решений

- [Дизайн](#appendix-design-review): принятые границы сервисов, контракты, состояния и инварианты.
- [Процессы](#appendix-g0-processes): 15 подробных схем транзакций, отказов, гонок и восстановления.
- [Граница лаборатории](#appendix-lab-boundary): общий механизм, адаптеры, кандидат второго бизнеса и доказательство переносимости.
- [Среда](#appendix-environment): Docker для Windows/macOS, распределение контейнеров и ограничения измерений.
- [Ресурсные профили](#appendix-resource-profiles): снимок Windows-хоста и текущий WSL-бюджет и отдельно ещё не измеренные среды.
- [Roadmap](#appendix-roadmap): этапы 0–7, критерии готовности, риски и 12 сценариев экспериментов.
- [Workflow](#appendix-development-process): роли человека и ИИ, состав ревью и контрольные точки G0–G4.
- [Правила ревью](#appendix-code-review): стандарт Go, блокирующие замечания, проверки и критерии слияния.
- [ADR 0001](#appendix-adr-0001): свидетельство явного принятия G0.
- [ADR 0002](#appendix-adr-0002): согласованный порядок фиксации решений в PR, ADR и обзорной странице.
- [ADR 0003](#appendix-adr-0003): принятое решение о Go smoke и его последствия.
- [Структурированное состояние](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/CONTEXT_STATE.json): принятые решения, датированные проверки и история ревизий.

`PROJECT_CONTEXT.md` генерируется из источников и содержит их полные тексты. Источники редактируются напрямую; снимок обновляется через `tools/update_context.ps1`. Подробные документы пока ведутся на русском; публичное введение доступно на двух языках.

<a id="appendix-roadmap"></a>

## Приложение: Полный roadmap, архитектура, эксперименты и методика

Полный исходный документ: `docs/ROADMAP.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Faultline Lab: roadmap лаборатории инженерных экспериментов

Цель проекта — воспроизводимо исследовать, как небольшой бизнес-процесс ведёт себя под нагрузкой и при сбоях, а затем проверять эффективность изменений. Итоговый продукт состоит из трёх бизнес-сервисов и инструментов вокруг них: генератора нагрузки, сценариев отказов, наблюдаемости, проверок корректности и отчётов.

Это план разработки, а не результаты уже проведённых испытаний. Все численные пороги ниже — начальные гипотезы, которые нужно откалибровать на своём компьютере.

Разработку выполняет Codex; пользователь — владелец проекта и ревьюер. [Workflow](#appendix-development-process) определяет приёмку. G0 принят 13 октября 2026 года. Выполнение этапов, задачи и ревью ведутся в [GitHub](https://github.com/Includeoyasi/faultline-lab/issues); переносимый контекст хранит спецификации и датированную историю.

#### 1. Допущения и границы

Обозначения G0–G4 относятся к контрольным точкам приёмки; P01–P05 — к небольшим пакетам кода с отдельным ревью. Первый PR с каркасом — P01, между G0 и G1. G1 принимается позже, после готовности основы заказа и склада; следующий PR не обозначает прохождения этой точки.

- Один исполнитель разработки — Codex, владелец проекта и ревьюер — пользователь. Первоначальные оценки объёма рассчитаны для ручной работы разработчика, знакомого с Go, SQL и Docker: один день равен примерно шести продуктивным часам. Эти оценки не являются сроками выполнения Codex; прогресс оценивается по готовности результатов и согласованным контрольным точкам.
- Первая среда — Linux-контейнеры в Docker Desktop/Compose на Windows и macOS; Go и инструменты сборки/проверок выполняются в контейнерной среде. Этот подход подтверждён владельцем проекта. Планируются нативные образы для x64 и ARM64; ресурсы и ограничения Docker фиксируются в каждом эксперименте. [Состав среды и статус проверок](#appendix-environment).
- Ровно три бизнес-сервиса: Orders, Inventory и Payments. Воркеры, оркестратор процесса и outbox publisher работают внутри соответствующего сервиса. Утилиты лаборатории и инфраструктура не становятся дополнительными бизнес-сервисами.
- Одна валюта, целые суммы в минимальных денежных единицах, один склад, небольшой фиксированный каталог и заранее загруженные остатки.
- В MVP нет авторизации, доставки, налогов, скидок, UI магазина и возвратов уже доставленных товаров. Возврат платежа нужен для компенсации отменённого заказа.
- Платёжный провайдер — модель внутри Payments. Она умеет детерминированно отвечать успехом или отказом, задерживать ответ и моделировать неизвестный результат. Реальные деньги не используются.
- Сначала один PostgreSQL с тремя базами и отдельными ролями. У сервисов нет прав читать чужие базы. Позже базы можно разнести, если именно разделение ресурсов становится предметом эксперимента.
- Один Kafka broker в первом профиле. Этот профиль проверяет обработку событий, backlog и перезапуски, но не отказоустойчивость реплицированного кластера. Для последней нужен отдельный профиль с несколькими брокерами.
- Версии Go, зависимостей и образы контейнеров закрепляются явно. В отчёт попадают теги и digest образов; плавающий `latest` не является описанием воспроизводимой среды.

Достаточный конечный результат: новый сценарий задаётся конфигурацией; один запуск создаёт нагрузку, применяет сбой, собирает данные и проверяет инварианты; сравнение двух вариантов формирует отчёт с выводом, ограничениями и исходными данными.

#### 2. Архитектура

```mermaid
flowchart LR
    L[Loadgen и labctl] -->|gRPC CreateOrder / GetOrder| O[Orders]
    L -->|gRPC GetAvailability| I[Inventory]

    subgraph B[Три бизнес-сервиса]
        O -->|gRPC Reserve / Commit / Release| I
        O -->|gRPC GetPayment / Cancel / Refund| P[Payments]
    end

    O -->|OrderAwaitingPayment| K[Kafka]
    K -->|Событие заказа| P
    P -->|PaymentSucceeded / Declined / Refunded| K
    I -->|InventoryReserved / Committed / Released| K
    K -->|Результаты платежа и резерва| O

    subgraph DB[Один PostgreSQL, три отдельные базы]
        DO[(orders_db)]
        DI[(inventory_db)]
        DP[(payments_db)]
    end
    O --> DO
    I --> DI
    P --> DP
    I --> R[(Redis: кеш доступности)]

    O -.->|Метрики| M[Prometheus]
    I -.->|Метрики| M
    P -.->|Метрики| M
    O -.->|JSON stdout| A[Grafana Alloy]
    I -.->|JSON stdout| A
    P -.->|JSON stdout| A
    A --> LK[Loki]
    M --> G[Grafana и локальные алерты]
    LK --> G
```

Инфраструктура для MVP: PostgreSQL, Kafka, Redis, Prometheus, Grafana, Loki и Alloy. Prometheus хранит метрики; Loki хранит логи; Grafana показывает оба источника. Сбор Docker-логов через Alloy соответствует официальному примеру Grafana. [Документация Loki и Alloy](https://grafana.com/docs/loki/latest/get-started/quick-start/tutorial/).

На этапе расширенной диагностики добавить OpenTelemetry и Tempo. Они нужны для ответа на конкретный вопрос: где заказ провёл время — в gRPC, очереди, SQL или обработчике событий. Сбор трасс должен иметь одинаковую конфигурацию в сравниваемых прогонах.

##### Владение данными

| Сервис | Ответственность | Основные данные | Что не делает |
|---|---|---|---|
| Orders | Приём заказа, долговечное управление шагами, отмена и компенсации | orders, order_items, workflow_jobs, transition_history, inbox, outbox, request_keys | Не меняет складские и платёжные таблицы |
| Inventory | Остатки, цены каталога для стенда, резервирование, подтверждение продажи, освобождение | products, stock, reservations, reservation_items, operations, inventory_ledger, outbox | Не решает, успешен ли заказ в целом |
| Payments | Симуляция платежа, запрет позднего списания после отмены, возврат | payments, payment_jobs, payment_ledger, inbox, outbox, operation_keys | Не резервирует товар и не завершает заказ |

Цены фиксируются Inventory в момент резервирования; Orders сохраняет этот снимок. Payments получает сумму из события Orders, а не произвольную сумму клиента. Общими могут быть protobuf-контракты и инструменты телеметрии; бизнес-логика и SQL каждого сервиса остаются внутри него.

##### Назначение транспорта

**gRPC** используется там, где нужен непосредственный результат операции: принять заказ, узнать состояние, атомарно зарезервировать товар, подтвердить или освободить резерв, запросить статус платежа и выполнить компенсацию.

**Kafka** переносит факты и отделяет приём заказа от платежной обработки. Payments запускается по `OrderAwaitingPayment`. Orders принимает результаты Payments и Inventory. Inventory подтверждает фиксацию продажи событием `InventoryCommitted`, после которого Orders завершает заказ.

Таким образом, все три сервиса используют и gRPC, и Kafka. Долговечный управляющий процесс находится в Orders; отдельный сервис для него не нужен.

**Redis** — кеш чтения `GetAvailability`, принадлежащий Inventory. Сначала TTL около двух секунд с небольшим случайным разбросом и best-effort invalidation после изменения склада. Кеш может временно показывать устаревшую доступность. `Reserve` всегда проверяет PostgreSQL. Redis не хранит единственную копию идемпотентности, платежей или остатков.

Пользу кеша проверять на смешанном профиле просмотра и заказа. Проверять только создающие записи запросы и ожидать ускорения от кеша чтений было бы неверной постановкой эксперимента.

#### 3. Бизнес-процесс и контракты

##### Успешный путь

1. `Orders.CreateOrder(client_request_id, items)` в одной SQL-транзакции создаёт заказ, ключ запроса и долговечную задачу резервирования. После commit клиент получает `order_id` и `ACCEPTED`. Ответ означает принятие, а не завершение покупки.
2. Внутренний воркер Orders вызывает `Inventory.Reserve(order_id, operation_id, items)`.
3. Inventory атомарно создаёт резерв, обновляет остатки и пишет `InventoryReserved` в outbox. Повтор с тем же ключом возвращает прежний результат. Orders сохраняет резерв, сумму и `OrderAwaitingPayment` в своей транзакции.
4. Publisher Orders доставляет событие в Kafka. Payments записывает inbox и долговечную payment job в своей транзакции, затем подтверждает обработанный Kafka offset.
5. Воркер Payments выполняет симуляцию задержки вне SQL-транзакции. Короткая транзакция проверяет текущее состояние, записывает единственное списание и `PaymentSucceeded` в outbox либо фиксирует отказ и `PaymentDeclined`.
6. Orders атомарно принимает `PaymentSucceeded`, меняет состояние на `FINALIZING` и создаёт задачу `CommitReservation`. Inventory фиксирует продажу и `InventoryCommitted` в одной транзакции.
7. Orders обрабатывает `InventoryCommitted`, проверяет подтверждённую оплату, переводит заказ в `COMPLETED` и публикует `OrderCompleted`.

Если Orders упал после ответа Inventory, повторяется тот же шаг с тем же `operation_id`. Если Kafka временно недоступна, событие остаётся в outbox. Если воркер умер после получения задачи, истечение lease позволяет продолжить её другим экземпляром.

##### Состояния

| Сущность | Основные состояния |
|---|---|
| Заказ | ACCEPTED → RESERVING → PAYMENT_PENDING → FINALIZING → COMPLETED |
| Ветка отказа заказа | REJECTED при нехватке товара; CANCELLING → CANCELLED после завершения компенсаций |
| Резерв | RESERVED → COMMITTED или RELEASED; RELEASED также служит запретом позднего Reserve для отменённого order_id |
| Платёж | PROCESSING → SUCCEEDED или DECLINED; отмена до списания → CANCELLED; возврат списанного платежа → REFUND_PENDING → REFUNDED |

`NEEDS_REVIEW` можно добавить для исчерпанного автоматического восстановления. Это нерешённый бизнес-процесс, который ухудшает результат эксперимента, а не успешный терминальный исход.

В MVP отмена клиентом разрешена до FINALIZING. Переход PAYMENT_PENDING → FINALIZING и принятие CancelOrder конкурируют за один условный переход версии в Orders: выигрывает одна ветка. После FINALIZING отмена отклоняется, процесс доводится до COMPLETED. Это позволяет не моделировать возврат уже проданного товара. Если намерение отмены победило раньше, поздний PaymentSucceeded запускает возврат и не создаёт Commit job. Старый воркер проверяет актуальную версию и разрешённый шаг перед выполнением мутации.

##### Отказы и компенсации

| Ситуация | Действие |
|---|---|
| Недостаточно товара | Атомарно отказать в резервировании всего заказа, завершить как REJECTED, не запускать оплату |
| Платёж отклонён | Зафиксировать отказ, освободить резерв, завершить как CANCELLED |
| gRPC timeout после мутации | Считать результат неизвестным; повторить тот же operation_id или запросить состояние |
| Kafka недоступна | Держать событие в outbox, контролировать возраст и объём накопления; при исчерпании заданного лимита ограничивать новый приём |
| Оплата успешна, Inventory временно недоступен | Оставить FINALIZING, повторять CommitReservation и сверять состояние; временный отказ не является поводом немедленно вернуть деньги |
| Резерв уже освобождён, оплата успешна | Перейти к возврату; CANCELLED выставить после подтверждения возврата |
| Позднее платёжное событие после начала отмены | Продолжить отмену: запретить ещё не выполненное списание либо вернуть выполненное |
| Результат Commit неизвестен | Сначала сверить резерв. Если COMMITTED, завершить заказ; если RELEASED, компенсировать оплату |

`ReleaseReservation(order_id)` должен сохранять терминальный запрет даже если Reserve ещё не успел создать запись. `CancelPayment(order_id)` аналогично сохраняет запрет будущего списания; если списание уже произошло, возвращает состояние для последующего `RefundPayment`. Иначе задержанный запрос или событие могут воскресить отменённый процесс.

В MVP резерв не исчезает по независимому TTL. Срок процесса контролирует Orders, который сохраняет намерение отмены и доводит компенсации до конца. Автоматическое истечение резерва добавить позже как отдельный протокол и проверить гонки с оплатой и Commit.

##### gRPC API

| Сервис | Методы |
|---|---|
| Orders | CreateOrder, GetOrder, GetOrderByRequestKey, CancelOrder |
| Inventory | GetAvailability, Reserve, GetReservation, CommitReservation, ReleaseReservation |
| Payments | GetPayment, CancelPayment, RefundPayment |

Mutating RPC принимают ключ операции и проверяют hash существенных параметров. Повтор ключа с другим payload — конфликт, а не новый заказ. Платёжная симуляция запускается событиями; отдельный публичный `Charge` в MVP не нужен.

Идемпотентность включает подтверждённые бизнес-отказы. Например, недостаток товара сохраняется как результат operation_id без частичных складских изменений. Повтор того же ключа после пополнения склада не превращает завершённый отказ в новый успешный резерв.

##### Kafka

| Топик | События | Consumer group |
|---|---|---|
| orders.events.v1 | OrderAwaitingPayment, OrderCompleted, OrderCancelled | payments-v1, фильтрует только нужные события |
| payments.events.v1 | PaymentSucceeded, PaymentDeclined, PaymentRefunded, PaymentCancelled | orders-payments-v1 |
| inventory.events.v1 | InventoryReserved, InventoryCommitted, InventoryReleased | orders-inventory-v1 |

Envelope: `event_id`, `event_type`, `schema_version`, `order_id`, `aggregate_version`, `occurred_at`, `correlation_id`, `causation_id`, `traceparent`, `payload`. `aggregate_version` относится к владельцу события; нельзя сравнивать версии Inventory и Payments как одну последовательность.

Начать с трёх партиций на топик и ключа `order_id`, затем экспериментировать с другим числом. Использовать обычные consumer groups. Kafka упорядочивает записи внутри партиции; общего порядка между топиками нет. Параллельная обработка в приложении также не должна нарушать требуемый порядок. Kafka-транзакции сами по себе не создают общую атомарную транзакцию с PostgreSQL. [Официальное описание Kafka](https://kafka.apache.org/43/design/design/).

Publisher outbox сохраняет порядок версий одного агрегата: несколько publishers не должны одновременно отправлять разные версии одного order_id в случайном порядке. В MVP допустим один активный publisher на сервис с DB lease; другие экземпляры могут принять lease после его потери. Позже проверить шардирование publishers. State machine обязана безопасно переживать дубликаты и опоздавшие факты даже при корректном порядке публикации.

#### 4. Гарантии корректности

##### Обязательные инварианты

1. `available >= 0`, `reserved >= 0`, `on_hand >= reserved` в Inventory.
2. При отсутствии пополнений `initial_stock = available + reserved + sold`. При пополнениях они учитываются в ledger; одна операция Commit или Release меняет баланс только один раз.
3. На order_id не более одного успешного списания. Сумма возвратов не превышает сумму списания. В MVP один полный возврат.
4. `COMPLETED` означает подтверждённую оплату на ожидаемую сумму и `COMMITTED` резерв.
5. `CANCELLED` означает отсутствие активного резерва и отсутствие невозвращённого списания; поздние операции не могут нарушить этот результат.
6. На client_request_id существует максимум один заказ с неизменным смыслом запроса.
7. Каждый принятый заказ существует в Orders и учтён итоговым verifier как терминальный, незавершённый или неизвестный из-за недоступности данных. Неучтённые заказы недопустимы.

Состояние PAYMENT_PENDING при резерве RESERVED — нормальная временная рассогласованность. Успешная оплата и FINALIZING во время краткого отказа Inventory также допустимы, но их длительность учитывается. Завершённый заказ без оплаты, отрицательные остатки и двойной debit — нарушение безопасности при любой длительности.

Глобальный verifier после прекращения нагрузки читает экспорт данных владельцев или использует отдельные read-only учётные записи трёх баз. Это инструмент лаборатории; сами сервисы не получают доступ к чужим данным. Во время активных изменений глобальная сверка должна учитывать версии и повторно подтверждать подозрительные пары, чтобы разные моменты чтения не выглядели ложным нарушением. После drain выполняется обязательная полная сверка.

##### Механизмы и момент внедрения

| Механизм | Когда | Назначение и проверка |
|---|---|---|
| SQL constraints, уникальные operation keys, условные переходы состояния | Этап 1 | Конкурентные Reserve и повторные запросы не создают лишний резерв или эффект |
| Долговечные шаги Orders, lease, version check | Этапы 1–2 | Перезапуск продолжает сохранённый шаг; старый воркер не перезаписывает новое состояние |
| Transactional outbox у всех трёх сервисов | Этап 2 | Падение между SQL commit и Kafka publish не теряет событие |
| Inbox и durable jobs у Kafka-потребителей | Этап 2 | Дубликат события не создаёт второй job или бизнес-эффект |
| Reconciliation и компенсации | Минимум в этапе 2, расширение в этапе 4 | Неизвестный результат RPC и поздние результаты не оставляют ложный terminal state |
| DLQ/quarantine и управляемый replay | Этап 4 | Некорректное событие диагностируется и воспроизводится без тихой потери процесса |
| Политика circuit breaker и admission control | Этап 5 после измерения | Ограничить каскадное накопление работы; оценить успешный goodput, отказы и восстановление |

Inbox, бизнес-переход и создаваемые jobs/outbox записываются одной транзакцией. Kafka offset подтверждается после неё. Если дальше нужна долгая работа, она уже представлена durable job. При параллельной обработке offset можно продвигать только через непрерывный диапазон завершённых записей, а не сразу до самого высокого готового offset.

Outbox помечается отправленным после подтверждения Kafka; падение перед отметкой приведёт к повторной публикации. Поэтому повторная доставка ожидаема, и защита строится на бизнес-ключах и SQL-транзакциях. Обещание проекта — определённые идемпотентные эффекты и проверяемые инварианты в заданной модели отказов, а не универсальная сквозная exactly-once.

Для склада использовать короткую транзакцию: условное увеличение `reserved` при достаточном `on_hand - reserved`, запись резерва и ledger. Несколько SKU блокировать в стабильном порядке; при недостатке хотя бы одной позиции откатывать весь резерв. PostgreSQL при конкурентном UPDATE в Read Committed повторно проверяет WHERE для изменённой строки. Это полезно для такого условного обновления, но не заменяет проектирование всей транзакции. [Transaction isolation](https://www.postgresql.org/docs/current/transaction-iso.html). Единый порядок блокировок и короткие транзакции уменьшают риск deadlock. [Explicit locking](https://www.postgresql.org/docs/current/explicit-locking.html).

На RPC всегда задавать deadline и передавать отмену в SQL и последующие вызовы. gRPC по умолчанию не устанавливает deadline; приложение должно прекращать запущенную работу после отмены. [gRPC deadlines](https://grpc.io/docs/guides/deadlines/).

Это относится к работе текущего запроса. Уже сохранённая и принятая durable job имеет собственный контекст процесса и бюджет бизнес-шага; уход клиента после CreateOrder не удаляет заказ и не отменяет его обработку. Повтор клиентского запроса разрешает неизвестный результат через тот же request key.

Начальная политика: одна повторная попытка для идемпотентной мутации, только для классифицированной временной ошибки и внутри оставшегося deadline; backoff с jitter. Дальнейшее восстановление делает durable job с сохранённым next_attempt_at. Не включать независимые длинные retry loops одновременно у клиента, gRPC и воркера. Встроенную retry policy gRPC нужно задавать явно и проверять счётчик реальных попыток. [gRPC retry](https://grpc.io/docs/guides/retry/).

Длительные задержки провайдера, ожидание gRPC и Kafka не должны удерживать открытую SQL-транзакцию. Ограниченные пулы и worker queues — исходная управляемая конфигурация; изменение их размеров станет экспериментом.

#### 5. Наблюдаемость

##### Логи и корреляция

JSON-логи в stdout: `timestamp`, `level`, `service`, `instance`, `version`, `run_id`, `order_id`, `operation_id`, `event_id`, `correlation_id`, `trace_id`, `step`, `attempt`, `duration_ms`, `outcome`, `error_code`.

Логировать переход состояния, неизвестный результат операции, повтор, компенсацию, конфликт ключа и помещение в quarantine. Содержимое платёжных данных для стенда — синтетическое. Успешные однотипные чтения можно семплировать; журналы бизнес-операций, на которые опирается проверка, находятся в PostgreSQL, а не только в логах.

`order_id`, `trace_id`, `event_id` и `run_id` — поля логов, не Prometheus labels и не индексные labels Loki. Labels метрик ограничить `service`, `method`, `code`, `step`, `event_type`, `outcome` с заранее ограниченными наборами. Высокая кардинальность увеличивает стоимость хранения и обработки метрик. [Prometheus instrumentation](https://prometheus.io/docs/practices/instrumentation/).

##### Метрики и дашборды

| Дашборд | Показатели | Какой вопрос решает |
|---|---|---|
| Бизнес-процесс | Offered/sent/accepted orders, completed orders/s, business rejects, technical failures, отмены, время до терминального состояния, oldest pending age, незавершённые заказы | Завершает ли система покупки при заданном спросе? |
| gRPC и Go | p50/p95/p99, timeout/error/retry rate, inflight, goroutines, GC, CPU/RSS | Где растёт задержка и сколько работы удерживается? |
| PostgreSQL | Время acquire из пула, занятые/свободные соединения, SQL latency, lock waits, deadlocks, долгие транзакции | База медленно работает или приложение ждёт её допуска? |
| Kafka и jobs | Lag по group/topic/partition, возраст старейшего необработанного события, event completion rate, длительность rebalance, job/outbox depth и oldest age | Где накапливается работа и успевает ли она завершаться? |
| Redis | Hit/miss, ошибки и latency, частота fallback, загрузка SQL при fallback, возраст возвращённых данных | Кеш снижает нагрузку без опасного доверия устаревшим данным? |
| Корректность | Нарушения каждого инварианта, dangling reservations, unmatched payments, failed reconciliation | Не куплено ли ускорение ценой неверных результатов? |

Lag по committed offsets показывает прогресс передачи в durable job, если consumer подтверждает событие на этом этапе. Поэтому низкий Kafka lag не доказывает, что платежи уже обработаны: отдельно нужны возраст и размер payment_jobs.

Метрики сервисов снимает Prometheus. На ранних этапах runner сохраняет `docker stats`, доступные SQL-показатели и Kafka offsets. На этапе 3 добавить инфраструктурные exporters с проверкой совместимости своего Docker окружения. Накладные расходы verifier и exporters также фиксировать и держать одинаковыми в A/B.

##### Начальные SLI/SLO

Профиль для SLO: достаточно товара, симулятор всегда одобряет платёж, нет внесённых отказов, нагрузка R0. R0 — не обещанное число: её определить первым ступенчатым тестом.

| SLI | Начальная гипотеза |
|---|---|
| Приём валидного заказа | Не менее 99% уникальных отправленных заказов получают подтверждение принятия за 300 мс |
| Завершение бизнеса | Не менее 99% уникальных отправленных валидных заказов переходят в COMPLETED за 2 секунды с момента отправки |
| Безопасность данных | Ноль отрицательных остатков, двойных списаний и ложных COMPLETED |
| Полнота обработки после drain | Каждый принятый заказ учтён; в профиле без отказов нет незавершённых процессов |
| Восстановление | После 30 секунд недоступности зависимости backlog и задержки возвращаются в исходный диапазон за 120 секунд; уточнить по baseline |

Доля завершённых среди принятых — полезная диагностическая метрика. Для оценки способности обслужить спрос нужен также знаменатель всех отправленных валидных заказов: иначе отклонение большей части нагрузки искусственно улучшит результат.

Ожидаемые sold-out и payment declines учитывать отдельно от технических ошибок. Правила ожидаемого поведения задаются до запуска. Для fault-сценариев заранее задаются отдельные критерии допустимой деградации и восстановления; нельзя после теста исключать из знаменателя неудачные запросы.

Начальные локальные алерты в Grafana:

- Верификатор подтвердил нарушение инварианта — сразу критический сигнал.
- При положительном входящем потоке нет бизнес-завершений 30 секунд — сигнал остановки процесса.
- Старейший заказ старше 10 секунд либо объём backlog растёт несколько окон подряд — сигнал деградации.
- gRPC p99 превышает 300 мс или доля технических ошибок превышает 1% на достаточной выборке несколько минут — сигнал нарушения начального бюджета.
- Долгое ожидание соединения, рост outbox age, недоступность telemetry collectors — отдельные диагностические сигналы.

Runner фиксирует момент начала проблемы, переход алерта в firing и восстановление. Порог и `for` подбирать по baseline; выбранный `for` ограничивает минимальную скорость обнаружения. Во время запланированных сбоев сигналы не скрывать: они являются частью результата.

#### 6. Генератор и управление экспериментами

Начать с маленького Go CLI `loadgen`, использующего те же protobuf-контракты. Поддержать open-loop запуск заказов с заданной частотой, фиксированный seed, список SKU и независимую сверку результатов. Не создавать для него отдельную постоянную бизнес-платформу.

Open-loop нужен, чтобы рост latency не снижал автоматически предлагаемую системе нагрузку. Такое снижение в закрытой модели способно скрывать перегрузку; различие подробно описано в документации k6. Реализуется оно и в собственном Go-генераторе. [Open and closed load models](https://grafana.com/docs/k6/latest/using-k6/scenarios/concepts/open-vs-closed/).

Для каждой операции сохранять плановое время отправки, фактическую отправку, ответ, ключ запроса и итоговый бизнес-статус. Если генератор исчерпал свой лимит inflight, записать skipped/late launches: не тормозить незаметно расписание. Измерять CPU генератора и расхождение target/sent rate. Результат при насыщении самого генератора маркировать как непригодный для оценки capacity стенда.

Параметры профиля:

- orders/s, ступени роста и burst;
- доля чтений каталога и статуса заказа;
- число товаров в заказе, размер каталога и распределение popular SKU;
- исходный склад: избыточный для capacity-тестов, ограниченный для hot-stock;
- детерминированная доля payment declines и задержек, keyed по order_id и seed;
- лимиты inflight генератора, число воркеров сервисов и pool sizes;
- таймлайн сбоя, длительность, направление сетевого воздействия;
- прогрев, основное окно, drain, повторы и ресурсный бюджет.

Способы внесения отказов: test-only delay/error hooks, управляемая остановка и запуск контейнера, сетевой proxy по нужной зависимости. Задержка SQL в тестовой ветке должна быть отличима от реальной блокировки/плохого плана SQL. Все hooks отключаются обычной конфигурацией, а их активация попадает в manifest.

#### 7. Roadmap

Этапы 0–2 дают MVP. Он содержит весь обязательный стек, сквозной процесс и первый отчёт, собранный вручную. Этапы 3–7 делают лабораторию удобной для более глубоких экспериментов.

Контрольные точки владельца проекта: G0 — архитектура перед кодом; G1 — базовая работа с заказом и складом после этапа 1; G2 — сквозной MVP после этапа 2; G3 — устойчивость после этапа 4; G4 — итоговая лаборатория после этапа 7. G0 полностью принят 2026-10-13; G1–G4 ещё не достигнуты. Внутри этапов изменения готовятся небольшими пакетами для код-ревью. Переход через контрольную точку выполняется после явного согласования пользователя. Существенные изменения ранее согласованной архитектуры выносятся на отдельное согласование.

| Этап | Цель и задачи | Зависит от | Результат и критерий готовности | Демонстрация | Дни | Основной риск |
|---|---|---|---|---|---|---|
| 0. Контракт эксперимента | Описать состояние и инварианты; определить API/события; выбрать baseline data и ресурсный профиль; подготовить monorepo и Compose | — | Контракты и initial manifest сохранены; инфраструктура поднимается; у каждой БД отдельная роль | Проверка доступности и невозможности читать чужую БД | 2–3 | Долгое проектирование вместо первого работающего среза |
| 1. Создание и резерв | Go-сервисы, миграции; Create/Get; атомарный Reserve/Commit/Release; модель платежа; request keys; durable step; JSON-логи и базовые метрики | 0 | Заказ переживает перезапуск в RESERVING; повтор запроса не дублирует заказ; конкуренция не создаёт отрицательный склад | 1000 попыток купить 100 единиц одного SKU | 4–6 | Неверная транзакция или слишком долгие блокировки |
| 2. Сквозной MVP | Kafka workflow; outbox/inbox; payment jobs; минимальные компенсации; кеш чтения Redis; Loki/Alloy/Grafana; loadgen и verifier | 1 | 1000 заказов завершаются; отказы оплаты освобождают товар; повторы и контрольные crash windows не нарушают инварианты; первый baseline отчёт сохранён | Ступени нагрузки и сравнение одного изменения worker limit | 6–9 | Попытка отложить восстановление и потерять доверие к baseline |
| 3. Диагностика | Полные дашборды и алерты; SQL/пул/Kafka/job telemetry; OpenTelemetry и Tempo; учёт overhead | 2 | Для медленного заказа виден участок ожидания; искусственная задержка даёт ожидаемый сигнал и алерт; конфигурация dashboards хранится в репозитории | 500 мс задержки Payments и сравнение таймлайна событий, метрик и trace | 3–5 | Высокая кардинальность и влияние самой телеметрии |
| 4. Устойчивость | Late events, cancel races, timeout-after-commit, DLQ/quarantine и replay; проверки reconciliation; schema compatibility | 2; наблюдаемость из 3 помогает | Сценарии имеют утверждённые результаты; после исправления и replay нет потерь и двойных эффектов; нерешённые заказы явно учтены | Убить publisher после publish до отметки; задержать успех оплаты до отмены | 6–9 | Объявить компенсацию завершённой при неизвестном результате |
| 5. Нагрузка и решения | Прогон каталога узких мест; SQL/pool/workers; retries; cache; backpressure; сравнение гипотез по одной | 3–4 | Минимум четыре сценария имеют baseline, изменение и повторное измерение; улучшение не ухудшает корректность | Short transaction против долгой; ограниченные retries против retry storm | 5–8 | Сменить несколько факторов и потерять причинную связь |
| 6. Масштабирование | gRPC routing, несколько реплик, Kafka partitions/publishers; connection budgets; при необходимости разнести базы; optional multi-broker profile | 5 | Сравнение 1/2/4 реплик при том же суммарном CPU и при увеличенном бюджете; показано ограничение масштабирования | Replicas vs partitions; hot SKU vs равномерная нагрузка | 4–6 | Пиннинг gRPC на одном backend, горячая строка склада или рост суммарного SQL-пула |
| 7. Отчёты и regression suite | labctl автоматизирует setup/run/drain/verify/export/compare; JSON manifest, графики, Markdown report; короткий набор сценариев | 2, шаблон стабилизирован в 5–6 | Одна команда воспроизводит сравнение и выдаёт артефакты; smoke/regression suite находит намеренно внесённую ошибку | Повтор A/B с одинаковым seed и сравнение отчётов | 4–6 | Отчёт скрывает пропуски телеметрии, backlog или различия среды |

Базовая сумма: **34–52 рабочих дня, приблизительно 204–312 часов**. При 20 часах в неделю — ориентировочно 10–16 недель; при 10 часах — 20–31 неделя. Это оценка при принятых допущениях, без времени на изучение Go/Docker с нуля и без Kubernetes. MVP: **12–18 рабочих дней, 72–108 часов**. После MVP уточнить остальные оценки по фактической скорости.

##### Первый полезный эксперимент после MVP

1. При достаточном складе и всегда успешной оплате выполнить ступени 10 → 25 → 50 → 100 orders/s, по 3–5 минут. При необходимости скорректировать диапазон. Снять accepted/completed throughput, e2e latency, oldest job age, pool waits, Kafka и outbox backlog.
2. Найти первый устойчивый узкий участок. Если это последовательная payment processing, проверить изменение `payment_workers: 1 → 4` при тех же CPU/памяти. Если ограничение другое, выбрать изменение по наблюдаемому сигналу.
3. Повторить одинаковый профиль и сравнить goodput, p95/p99, незавершённые заказы, CPU и инварианты. Не менять одновременно pool size, число реплик и профиль отказов.
4. Сохранить первый отчёт. Это диагностический эксперимент; подтверждение предельной capacity требует более долгих и повторных прогонов на этапе 5.

##### Когда нужен Kubernetes

После этапа 7 возможен отдельный этап на 5–10 дней, если предмет исследования — rolling update, readiness, распределение по узлам, лимиты/requests или autoscaling. Он не нужен для первых экспериментов с Kafka, gRPC, SQL, кешем и нагрузкой. Автомасштабирование должно опираться на проверенный сигнал, например возраст работы, а не только CPU. Одноузловый кластер всё ещё не доказывает устойчивость к потере физического узла.

#### 8. Каталог экспериментов

Нагрузочные сценарии растят спрос или меняют его распределение. Fault-сценарии намеренно ломают зависимость. Они имеют разные критерии успеха: максимальный устойчивый goodput против ограниченной деградации и восстановления.

| ID / тип | Проблема и воспроизведение | Сигналы | Гипотеза / изменение | Оценка улучшения | Компромисс |
|---|---|---|---|---|---|
| L1: рост нагрузки | Ступени до выхода из SLO; товар не заканчивается, оплата всегда успешна | offered/accepted/completed, p99, oldest pending, CPU, pool/job waits | Изменить ровно один обнаруженный bottleneck: индекс, короткую транзакцию, число воркеров | Максимальный устойчивый orders/s при одном SLO и бюджете; backlog не имеет устойчивого роста | Латентное накопление работы может выглядеть как высокая скорость приёма |
| L2: burst | После спокойного окна подать 5×R0 на 30 секунд, затем R0 | Peak pending, heap/goroutines, rejects, drain time | Ограниченный inflight, admission control, раздельные worker budgets | Меньше технических ошибок и faster recovery; goodput и доля отказов показаны вместе | Shed requests нельзя исключать из общей оценки спроса |
| L3: горячий товар | 80–90% заказов на один SKU; отдельные тесты с обильным и ограниченным остатком | Lock wait, Reserve latency, business sold-out, inventory ledger | Сократить locked section; стабильный lock order; сравнить сериализацию по SKU как альтернативу | Короткие ожидания и корректный баланс; после terminal/drain для 100 единиц нет более 100 продаж и нет dangling reserve | Больше реплик не устраняет конкурентную запись в одну строку; serialization может снизить параллелизм |
| L4: SQL bottleneck | Отдельно плохой план чтения, долгий lock и маленький pool; не смешивать их | SQL latency/plan, acquire wait, lock wait, deadlocks, query count | Индекс для своего запроса, уменьшение времени транзакции, изменение pool cap | Улучшение e2e и SQL при том же профиле; pool не перегружает общий PostgreSQL | Большой pool способен ухудшить базу; индекс увеличивает цену записи |
| L5: Kafka backlog | Снизить throughput consumer либо processing worker; сохранять поток producer | Lag по партициям, возраст events/jobs, oldest outbox, rebalance | Batching, workers, число обычных consumers и partitions в отдельных тестах | Возраст backlog и время догонки, бизнес goodput | Низкий Kafka lag при большом durable-job backlog — ложный признак здоровья |
| L6: partition skew | Для отдельного профиля большинство ключей направить в одну партицию; держать stock contention отдельным фактором | Lag/throughput по партициям, idle consumers | Изменить ключ или распределение, если бизнес-порядок это допускает | Снижение перекоса и e2e latency при сохранении семантики | Новое число партиций меняет распределение ключей; не менять его на ходу без плана порядка и миграции |
| F1: зависимость медленная/недоступна | 500 мс → 2 с задержки или остановка Inventory/Payments на 30 секунд | Deadlines, job age, accepted/completed, inflight, recovery | Deadline budgets, bounded retries, durable reconciliation; затем breaker | Ограниченная память/очередь, ноль нарушений, recovery target | Слишком короткий deadline увеличивает неизвестные результаты и повторную работу |
| F2: retry storm | Временно возвращать UNAVAILABLE; сравнить агрессивную и ограниченную политику | Attempts/logical request, load amplification, CPU, queue age | Одна retry policy, backoff+jitter и budget | Меньше лишних попыток, выше completed goodput, короче восстановление | Более длинное ожидание единичного клиента; breaker/admission отклоняют часть спроса |
| F3: crash windows | Убить процесс после DB commit до publish, после publish до sent marker, после consumer commit до offset | Replay count, duplicate suppression, outbox/jobs, ledger | Outbox/inbox и устойчивые operation keys; проверяется действующая реализация | Сверка всех заказов после рестарта; ноль лишних списаний и lost events в модели сохранённых дисков | Проверка не включает уничтожение данных БД или исчерпание retention |
| F4: Redis down/stale | Остановить Redis, одновременно истечь ключам, задержать invalidation | Hit rate, stale-read rate, SQL fallback, Reserve correctness | Короткий fallback deadline, TTL jitter, singleflight; сравнить отказ и fallback | Read latency, SQL нагрузка и успешные заказы; ноль Reserve по кешу | Fallback может вызвать stampede в PostgreSQL; отображение остатков может быть устаревшим |
| F5: поздняя оплата/отмена | Задержать PaymentSucceeded до начала отмены; конкурирующие Commit/Release | State transitions, compensations, unmatched debit, terminal state | Durable cancel intent, cancellation barriers, idempotent refund | В итоге корректный COMPLETED или CANCELLED, единственный debit, отсутствие невозвращённых сумм у отменённых | Нужно больше состояний; нельзя объявлять неизвестный результат отменённым |
| F6: poison event / replay | Внести неверную схему или payload в конкретную партицию | Quarantine depth, errors, blocked partition, pending order age | Durable quarantine и контролируемый replay после исправления | Проблема видна, запись сохранена, процесс после replay корректен | Пропуск записи ради доступности ломает бизнес-порядок, если его не учитывает state machine |

Для F6 сначала сохранять проблемное событие и контекст в durable quarantine, затем разрешать продвижение offset согласно явной политике. Если DLQ находится в Kafka, outage Kafka не должен приводить к подтверждению потерянной DLQ-записи. Replay сохраняет исходный event_id/operation semantics и не обходит обычную идемпотентность. После изменения кода отдельно решить, что делать с ранее обработанными результатами: просто новый event_id способен создать второй эффект.

В F5 гонку приёма CancelOrder с PaymentSucceeded проверять через допустимые состояния Orders. Конкурентные Commit/Release — дополнительная проверка примитивов Inventory и внесённых нарушений протокола. В штатном MVP Orders не создаёт Release после начала FINALIZING; иначе для компенсации уже проданного товара понадобился бы отдельный протокол возврата склада.

#### 9. Методика сравнения

##### Протокол запуска

1. **Гипотеза до измерения.** Например: уменьшение времени удержания транзакции увеличит completed orders/s при том же CPU, без роста p99 и нарушений. Задать минимально полезный эффект и допустимые регрессии.
2. **Одинаковое исходное состояние.** Загрузить один seed dataset, очистить тестовые базы/кеш/топики либо создать изолированный namespace запуска. Сервисные operation keys и Kafka offsets не должны переходить из предыдущего прогона. Сброс допустим только для явно выделенного стенда.
3. **Manifest.** Git commit, dirty diff, версии и image digests, конфигурация, ресурсы, replicas, worker/pool caps, topics/partitions/groups, seed, stock, fault timeline, OS/Docker и ограничения машины.
4. **Прогрев.** Начать с двух минут, проверить выход в steady state. Прогревочные заказы обработать до конца и исключить из измеряемой когорты; состояние склада затем привести к одинаковому оговорённому старту. Cold-cache и warm-cache — разные сценарии.
5. **Основное окно.** Первое подтверждающее сравнение: десять минут одинаковой нагрузки. Диагностические ступени могут быть короче. Всё окно сценария, включая burst/fault, включается в анализ; исключаются только заранее оговорённые warmup и cooldown.
6. **Drain.** Прекратить новые заказы, оставить workers и зависимости работать. Ждать исчерпания jobs/outbox и завершения когорты до заданного лимита, например пяти минут. По истечении лимита перечислить unresolved и unknown; не выбрасывать их из данных.
7. **Verifier.** Проверить ledger, состояния, остатки, уникальность операций и полноту по списку request keys генератора. Неизвестные ответы CreateOrder сверить через GetOrderByRequestKey: timeout не означает, что заказ не был создан.
8. **Повторы.** Для первого A/B — минимум три пары, чередуя порядок A/B; для подтверждения предпочтительно пять. Использовать одинаковый seed внутри пары и заранее заданный набор разных seeds между парами. Повторить при нестабильной среде или небольшом эффекте.
9. **Сравнение.** Показать значения каждого прогона, median эффекта, разброс, объём выборки и ресурсную цену. Несколько прогонов не являются доказательством статистической значимости. Для близких результатов увеличить число пар и длительность.

##### Метрики результата

- `goodput = успешно завершённые заказы / длительность измеряемого окна`, с явным указанием, какие завершения относятся к окну. Отдельно показать completed во время drain.
- `business_success_ratio = завершённые ожидаемо успешные заказы когорты / все отправленные ожидаемо успешные заказы когорты` на контрольный момент.
- `accepted_completion_ratio` — вторичная метрика завершений среди принятых. Показывается вместе с accepted/offered, чтобы не скрывать load shedding.
- `recovery_time` — от снятия отказа до возврата backlog age и latency в заранее выбранный baseline диапазон и удержания в нём контрольное окно, например минуту.
- `amplification = физические попытки зависимого вызова / логические операции`.
- `resource_efficiency` — completed orders на CPU-second; также RSS и disk/network объём. Показать отдельно нагрузку сервисов, инфраструктуры и генератора.
- `invariant_violations` и `unresolved_after_drain` — отдельные обязательные результаты, а не примечание к производительности.

E2E latency в отчёте считается для когорты и имеет рядом число успешных, отказанных и незавершённых заказов. Percentile только успешно завершённых заказов без этих чисел может скрыть деградацию. Незавершённым не присваивать выдуманное время завершения; показывать их наблюдаемый возраст и долю цензурированных наблюдений.

Не усреднять p95/p99 между репликами. Для общего percentile объединять подходящие histogram buckets либо сырые наблюдения. Между прогонами можно показывать median их p95, прямо называя это статистикой прогонов, а не percentile всей выборки.

До/после сравнивать при одинаковом суммарном бюджете. Если добавлены CPU/память, выносить отдельный вариант «расширение ресурсов». Для реплик считать суммарные DB pools: четыре пула по 20 соединений дают другое давление на БД, чем один на 20. Проверять, что gRPC resolver/load-balancing policy действительно распределяет вызовы: один долгоживущий канал к одному backend не является тестом равномерного масштабирования.

При росте capacity искать максимальный поток, который выдерживает SLO и не показывает устойчивого роста backlog на всём контрольном окне. Быстрый приём с бесконечно растущей очередью — не устойчивый throughput.

Процент изменения считать только при ненулевом baseline и показывать абсолютную разницу. Если одна из конфигураций нарушает инварианты, её ускорение не принимается как улучшение. Если генератор или сбор данных не обеспечили заданную нагрузку/полноту, результат обозначается `inconclusive` с причиной.

#### 10. Минимальная автоматизация

Предлагаемая структура будущего репозитория:

```text
cmd/orders/
cmd/inventory/
cmd/payments/
cmd/loadgen/
cmd/labctl/
internal/orders/  internal/inventory/  internal/payments/
internal/lab/                  # Общие контракты и механизм экспериментов
internal/adapters/shop/        # Операции, данные и проверки магазина
api/proto/
migrations/orders/  migrations/inventory/  migrations/payments/
deploy/compose/
observability/dashboards/  observability/alerts/
experiments/scenarios/
experiments/targets/           # Профили подключаемых систем
reports/<experiment>/<run>/
docs/adr/
```

Это предлагаемая структура реализации; текущие файлы содержат только план и шаблон отчёта.

Владелец проекта подтвердил разделение общей лаборатории, адаптера и тестируемой системы; магазин проходит через ту же границу. Детализация контрактов и кандидат второго бизнеса описаны в [границе лаборатории](#appendix-lab-boundary). Обработка изображений пока предложение; подключение второго бизнеса не объявляется выполненным или принятым критерием G4.

P01 включает контейнерный toolchain, образы трёх сервисов и минимальный Compose-запуск вместе с Go entrypoints/config/health/shutdown/JSON-логами и Go CI. Основная инфраструктура подключается в P02, наблюдаемость — на соответствующем этапе. Инструменты запускаются по задаче; адаптер не требует отдельного контейнера.

[P01](#appendix-p01) фиксирует каркас сервисов; [P02](#appendix-p02) — инфраструктуру и данные. Зависимость задач P02: #12 → #16; одна задача соответствует одному PR. [Учёт задач](#appendix-task-tracking) определяет переходы; актуальное выполнение — в GitHub Issues и PR.

Этап 2: запуск Compose и loadgen вручную, экспорт JSON/CSV, Markdown отчёт по шаблону. Этап 4: hooks и crash scripts. Этап 7: `labctl` объединяет уже проверенные действия.

Проектируемые команды `labctl`:

```text
labctl prepare --scenario baseline --variant A
labctl run --scenario baseline --variant A
labctl verify --run <run_id>
labctl export --run <run_id>
labctl compare --baseline <A_runs> --candidate <B_runs>
```

Контракт scenario должен содержать: seed, dataset, arrival schedule, read/order mix, fault timeline, replicas/resources, worker/pool limits, warmup/measure/drain и критерии оценки. Это собственная схема runner, а не готовый формат какого-либо существующего инструмента.

Артефакты одного запуска:

```text
manifest.json       # Полная конфигурация, версии, среда и dirty diff reference
requests.jsonl      # Расписание, отправка, результат и request keys
business-results.csv
metrics.json        # Экспорт временных рядов и счётчиков
faults.jsonl        # Фактическое, а не только запланированное время воздействий
invariants.json     # Проверки и список unresolved/unknown
logs/               # Выбранные логи вокруг инцидента
charts/             # PNG/SVG throughput, latency, backlog, resources
report.md
```

Report формируется из численных данных, а не только screenshots Grafana. Ссылки на дашборды с окном времени полезны для диагностики; графики и исходные данные сохраняются локально, чтобы отчёт не зависел от retention Grafana/Loki/Prometheus.

Regression suite: короткий успешный процесс, конкуренция за остаток, повтор CreateOrder, повтор события, crash после SQL commit, неизвестный gRPC результат, поздняя оплата при отмене. Проверки направлены на бизнес-инварианты и реальные crash windows; не требуется писать тест на каждый элемент конфигурации.

#### 11. Шаблон отчёта

Отдельный файл `docs/experiments/REPORT_TEMPLATE.md` содержит форму для заполнения. Обязательная последовательность:

1. Проблема, гипотеза и критерии принятия изменения.
2. Среда, ресурсы, версии, данные, профиль нагрузки и fault timeline.
3. Единственное проверяемое изменение и фактический diff.
4. Показатели A/B с выборкой и разбросом; отдельно accepted, completed, rejected, unresolved.
5. Графики по одной шкале времени, включая drain и восстановление.
6. Проверка инвариантов и связь ошибок с order_id/operation_id.
7. Цена ресурсов, регрессии, накладные расходы и ограничения.
8. Вывод `accepted`, `rejected` или `inconclusive` с объяснением.
9. Следующий эксперимент.

#### 12. Первые десять задач

1. **Зафиксировать бизнес-правила.** Написать короткий ADR с состояниями заказа/резерва/платежа, моментом COMPLETED и условиями CANCELLED. Готово: неизвестный платёжный результат имеет явное поведение.
2. **Создать минимальный monorepo.** Три Go entrypoints, отдельные конфигурации, graceful shutdown, gRPC health и JSON logger. Готово: каждый процесс запускается и сообщает версию.
3. **Поднять инфраструктуру Compose.** PostgreSQL, Kafka и Redis в core-профиле; Grafana/Prometheus/Loki/Alloy в obs-профиле. Готово: healthchecks проходят, image versions закреплены, логи видны в Grafana.
4. **Создать базы, роли и миграции.** Отдельные credentials, constraints и ledger/operation keys. Готово: Inventory не может читать orders_db, первичный stock загружается повторяемо.
5. **Описать protobuf и envelope событий.** Ключи операций, ошибки, версии, trace/correlation. Готово: можно сгенерировать клиентов и однозначно классифицировать business reject, temporary error и unknown outcome.
6. **Реализовать CreateOrder/GetOrder/GetOrderByRequestKey.** Durable reserve job и уникальный client_request_id. Готово: повтор возвращает тот же order_id; другой payload с тем же ключом отклоняется.
7. **Реализовать атомарный Inventory.** Reserve, Commit, Release, cancellation barrier, GetReservation. Готово: конкурентные заказы на ограниченный остаток сохраняют баланс; повтор мутации не меняет его второй раз.
8. **Сделать детерминированный Payments.** Durable jobs, approve/decline, задержка вне транзакции, cancel barrier и refund. Готово: одинаковый seed воспроизводит исходы; на order_id максимум один debit.
9. **Соединить сервисы через Kafka.** Outbox/inbox, устойчивые publishers, сохранённые шаги Orders и минимальная сверка неизвестных результатов. Готово: сквозной успех/отказ и контрольный перезапуск дают правильные terminal states.
10. **Запустить первый loadgen + verifier.** JSON/CSV результатов, базовый Grafana dashboard, чтение через Redis и первый отчёт. Готово: один ступенчатый тест воспроизводится, каждый запрос учтён, узкое место выбрано по данным.

Приоритизация: сначала корректный сквозной путь и измеримость, затем дополнительные способы диагностики, затем эксперименты и их автоматизация. Любое добавление инфраструктуры должно помогать ответить на конкретный вопрос или закрывать выявленный пробел в корректности.

<a id="appendix-design-review"></a>

## Приложение: Принятый архитектурный пакет G0

Полный исходный документ: `docs/architecture/DESIGN.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### G0: принятый архитектурный пакет

Статус: **G0 полностью принят владельцем проекта 2026-10-13, Europe/Moscow**. Прямое решение пользователя: «G0 полностью принят». Этот документ фиксирует принятую архитектуру из docs/ROADMAP.md и docs/architecture/PROCESSES.md. Выполнение — в [GitHub Issues](https://github.com/Includeoyasi/faultline-lab/issues) и PR.

#### Визуальное описание G0

Подробные схемы находятся в [docs/architecture/PROCESSES.md](#appendix-g0-processes): последовательность успешного заказа и локальные транзакции, состояния трёх сущностей, недостаток товара, отказ оплаты, гонка отмены, поздняя оплата, outbox/inbox, неизвестный результат RPC, Redis и цикл эксперимента.

##### Общий путь и владелец каждого шага

```mermaid
flowchart TB
    C["Клиент: создать заказ"] --> O["Orders: сохранить заказ и reserve_job"]
    O --> A["Ответ клиенту: ACCEPTED"]
    O --> I["gRPC → Inventory: зарезервировать"]
    I --> K1["Orders → Kafka: OrderAwaitingPayment"]
    K1 --> P["Payments: симуляция оплаты"]
    P --> K2["Kafka → Orders: PaymentSucceeded"]
    K2 --> F{"Отмена уже принята?"}
    F -->|Нет| FIN["Orders: FINALIZING; gRPC CommitReservation"]
    FIN --> IS["Inventory: COMMITTED; событие InventoryCommitted"]
    IS --> DONE["Orders: COMPLETED"]
    F -->|Да| CANCEL["Сохранить CANCELLING; выполнить компенсации"]
    CANCEL --> SAFE["Резерв освобождён; списание отсутствует или возвращено"]
    SAFE --> END["Orders: CANCELLED"]
```

Схема показывает успешное резервирование и оплату. Недостаток товара и PaymentDeclined раскрыты отдельно в разделах 4–5 визуального описания. Ответ ACCEPTED выдаётся после commit в Orders, а дальнейшие шаги выполняются по сохранённым задачам.

##### Состояния заказа и допустимая отмена

```mermaid
stateDiagram-v2
    [*] --> ACCEPTED
    ACCEPTED --> RESERVING
    RESERVING --> PAYMENT_PENDING: Резерв подтверждён
    RESERVING --> REJECTED: Товара недостаточно
    PAYMENT_PENDING --> FINALIZING: Оплата успешна; отмена не принята
    FINALIZING --> COMPLETED: Продажа и оплата подтверждены
    ACCEPTED --> CANCELLING: Отмена принята
    RESERVING --> CANCELLING: Отмена принята
    PAYMENT_PENDING --> CANCELLING: Отмена или отказ оплаты
    CANCELLING --> CANCELLED: Компенсации подтверждены
    REJECTED --> [*]
    COMPLETED --> [*]
    CANCELLED --> [*]
    note right of FINALIZING
        CancelOrder отклоняется.
        Сбой означает ожидание и восстановление.
    end note
```

Успешная оплата и отмена конкурируют за один переход версии Orders. Победа отмены означает компенсацию, победа финализации — доведение покупки до COMPLETED. Timeout сам по себе не создаёт терминальное состояние.

##### Долговечность и очереди

```mermaid
flowchart LR
    TX["TX производителя: бизнес-данные + outbox"] --> PUB["Publisher внутри владельца"]
    PUB --> K["Kafka"]
    K --> IN["TX потребителя: inbox + durable job"]
    IN --> OFFSET["Подтвердить offset"]
    IN --> W["Воркер внутри сервиса"]
    W --> EFFECT["Короткая TX: идемпотентный эффект + outbox"]
```

Низкий Kafka lag означает продвижение доставки до сохранённых jobs. Чтобы увидеть реальную обработку, дополнительно измеряются очередь jobs и время завершения заказа. Отдельные crash windows показаны в разделе 8 визуального описания.

#### Принятые решения

| Решение | Принятый вариант | Почему | Как проверяем |
|---|---|---|---|
| Бизнес-модель | Мини-магазин, один склад, одна валюта; заказ → резерв → симулированная оплата → завершение/отмена | Даёт конкуренцию за ресурсы, асинхронный процесс и измеримые бизнес-инварианты | Успех, нехватка товара, отказ и поздний результат оплаты |
| Три сервиса | Orders управляет процессом; Inventory владеет складом; Payments владеет платёжным ledger | У каждого изменения данных есть один владелец; управляющий процесс не требует четвёртого сервиса | Проверка прав БД и запрет чужих SQL-запросов |
| gRPC | Приём/чтение заказа; Reserve/Commit/Release; запрос статуса, отмена и возврат платежа | Непосредственный результат и явные ключи повторяемых операций | Повтор RPC, deadline после commit и сверка неизвестного результата |
| Kafka | Orders запускает оплату фактом OrderAwaitingPayment; Payments и Inventory публикуют результаты; Orders продолжает процесс | Можно отдельно менять производительность приёма и обработки, воспроизводить backlog и повторы | Consumer pause, повторная доставка, накопление durable jobs |
| Согласованность | Локальные SQL-транзакции плюс долговечный процесс в Orders; outbox/inbox в MVP | Сохранённый шаг восстанавливается после перезапуска, повтор события не создаёт второй эффект | Падения до/после publish и consumer commit |
| PostgreSQL | Один экземпляр, три базы и отдельные роли | Малый стартовый объём инфраструктуры при сохранении владения данными | Сервис не может читать чужую базу; измеряется общий ресурсный bottleneck |
| Redis | Кеш чтений Inventory, без авторитетных складских или платёжных записей | Даёт предмет эксперимента с hit rate, stale reads и fallback | Недоступность Redis и устаревшая доступность не позволяют продать лишний товар |
| Наблюдаемость | В MVP Prometheus, Grafana, Loki, Alloy; расширение трассировки на этапе 3 | Нужны как скорость приёма, так и завершения бизнеса, диагностика backlog и ошибок | Искусственная задержка видна в метриках/логах и вызывает ожидаемый сигнал |
| Среда | Локальный Docker Compose, закреплённые версии и manifest ресурсов | Воспроизводимый старт без обязательного отдельного кластера | Повтор прогонов из одной конфигурации |

#### Принятые бизнес-правила

- CreateOrder подтверждает принятие только после сохранения заказа и задачи в PostgreSQL. Уход клиента не отменяет уже принятую работу.
- Ключ клиентского запроса неизменно определяет один заказ. Повтор с другим payload возвращает конфликт.
- Подтверждённый бизнес-отказ тоже сохраняется как результат операции. Повтор отказанного Reserve с тем же operation_id не создаёт новый резерв после изменения остатков.
- COMPLETED возможен только при успешной оплате ожидаемой суммы и подтверждённом COMMITTED резерве.
- Отмена клиентом разрешена до FINALIZING. Начало финализации и принятие отмены конкурируют за один переход версии Orders. После FINALIZING отмена отклоняется; заказ доводится до завершения.
- Если отмена победила, поздняя успешная оплата компенсируется возвратом и не запускает Commit. CANCELLED означает отсутствие активного резерва и невозвращённого списания.
- Release и Cancel сохраняют запрет поздней операции даже если исходный резерв или платёж ещё не был создан. Симулятор платежей перед фиксацией эффекта проверяет актуальное состояние.
- Независимого истечения резерва по TTL в MVP нет. Срок обработки и намерение отмены сохраняет Orders.
- Измерение незавершённых процессов обязательно. NEEDS_REVIEW, неизвестный исход и timeout не считаются успешным завершением.

Полные состояния, контракт событий, операции, crash windows и инварианты описаны в разделах 3–4 docs/ROADMAP.md. Эти правила особенно важны для ревью: они определяют поведение при гонках и границы будущих API.

#### Граница MVP

Включены три Go-сервиса, отдельные базы, весь обязательный стек, сквозной успех/отказ, минимальные компенсации, идемпотентность, outbox/inbox, восстановление сохранённого шага, кеш чтений, базовые дашборды/логи, loadgen, verifier и первый ручной отчёт.

Для следующих этапов оставлены: полная трассировка и набор алертов, quarantine/replay, расширенный каталог отказов, несколько реплик, автоматическое сравнение и генерация отчётов. Kubernetes — отдельный предмет согласования, если появится эксперимент, который требует его возможностей.

Авторизация, доставка, скидки, полноценный магазинный UI и возврат доставленного товара находятся вне согласуемого MVP.

#### Первый пакет реализации после G0

Цель — получить проверяемую основу проекта и затем первый работающий срез «создать заказ → атомарно зарезервировать товар».

Порядок небольших пакетов для код-ревью:

1. **P01 — первый PR после G0, до G1.** Три Go entrypoints, конфигурация, health, shutdown и JSON-логи; контейнерный toolchain, образы сервисов, минимальный Compose и Go CI. Дополнение среды согласовано после G0 для Windows/macOS; детали и статус — в [описании среды](#appendix-environment).
2. **P02:** Compose для PostgreSQL/Kafka/Redis, отдельные БД/роли, миграции и seed. [Согласованная декомпозиция](#appendix-p02) сохраняет этот объём: две задачи под два будущих PR: #12 — PostgreSQL/миграции/seed, #16 — Kafka/Redis и проверка полного стенда. Реализация не начата.
3. **P03:** Protobuf-контракты с operation keys, состояниями и классификацией ошибок.
4. **P04:** Create/Get/GetByRequestKey и долговечная задача Orders.
5. **P05:** Reserve/Commit/Release/GetReservation с ledger, блокировками и cancellation barrier.

Каждый пакет проходит своё код-ревью. После подготовки и принятия этой основы проводится отдельная приёмка контрольной точки G1; первый PR не называется G1 и не требует предварительного прохождения G1.

К G1 предъявляются код, точные команды запуска, результаты повторного запроса, конкурентной покупки ограниченного остатка и перезапуска на сохранённом шаге. Остальная часть сквозного процесса развивается после G1.

#### Решение по G0 и последующие изменения

G0 принят полностью, без дополнительных условий в сообщении пользователя от 2026-10-13. Решение сохранено также в docs/adr/0001-approved-g0.md и PROJECT_CONTEXT.md. Следующая контрольная точка — G1 после подготовки основы заказа и склада.

Существенные изменения принятой архитектуры по-прежнему согласуются с владельцем проекта отдельно. Обычные детали реализации внутри G0 выбираются самостоятельно и предъявляются в небольших пакетах для код-ревью. Принятие архитектуры не означает, что реализация, испытания или следующие контрольные точки уже завершены.

После G0 подтверждено направление развития общей лаборатории с адаптерами: [детализация границы](#appendix-lab-boundary). Три сервиса и их бизнес-правила сохраняются; второй бизнес и расширение приёмки пока находятся на стадии предложения.

<a id="appendix-g0-processes"></a>

## Приложение: Подробные визуальные схемы G0

Полный исходный документ: `docs/architecture/PROCESSES.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### G0: визуальное описание процессов

Статус: **G0 полностью принят владельцем проекта 2026-10-13, Europe/Moscow**. Схемы раскрывают принятые решения из docs/architecture/DESIGN.md и docs/ROADMAP.md. Они описывают согласованное планируемое поведение; реализация и измерения ещё не выполнены.

Порядок чтения: общий процесс → успешный заказ → состояния → отказы и отмена → восстановление → кеш → цикл эксперимента.

#### Легенда

- В sequence-диаграммах сплошная стрелка — gRPC-запрос; обратная пунктирная — ответ. Для Kafka стрелки подписаны как публикация, доставка или подтверждение offset.
- `TX Orders`, `TX Inventory`, `TX Payments` — отдельные короткие транзакции в базе владельца. Общей транзакции трёх сервисов нет.
- `job` — сохранённая задача внутреннего воркера сервиса; это не четвёртый бизнес-сервис.
- `outbox` — событие, сохранённое вместе с бизнес-изменением; `inbox` — запись о принятом событии, защищающая от повторной обработки.
- `operation_id` идентифицирует бизнес-операцию. Повтор имеет тот же ключ и смысл; новая сетевая попытка не создаёт новую операцию.
- У Kafka-стрелок есть задержка и возможная повторная доставка. Положение стрелки на схеме не обещает мгновенную обработку.

#### 1. Общий бизнес-процесс

Вопрос для ревью: кто владеет каждым шагом и в какой момент заказ можно считать завершённым?

```mermaid
flowchart TB
    A["Клиент создаёт заказ"] --> B["Orders: сохранить заказ и задачу"]
    B --> ACK["Клиент получает ACCEPTED"]
    B --> R["Orders вызывает Inventory.Reserve по gRPC"]
    R --> CHECK{"Inventory: хватает товара?"}
    CHECK -->|Нет| REJ["Orders: REJECTED"]
    CHECK -->|Да| HELD["Inventory: RESERVED; Orders: PAYMENT_PENDING"]
    HELD --> EVT["Orders → Kafka: OrderAwaitingPayment"]
    EVT --> PAY["Payments: симуляция оплаты"]
    PAY --> RESULT{"Результат оплаты"}
    RESULT -->|Успех| PAID["Kafka → Orders: PaymentSucceeded"]
    RESULT -->|Отказ| DECL["Kafka → Orders: PaymentDeclined"]
    PAID --> BRANCH{"Orders: какая ветка уже выбрана?"}
    BRANCH -->|Отмена принята| REF["Компенсировать оплату и освободить резерв"]
    BRANCH -->|Отмена не принята| FIN["Orders: FINALIZING"]
    FIN --> COM["gRPC CommitReservation"]
    COM --> INV["Inventory: COMMITTED + InventoryCommitted"]
    INV --> DONE["Orders: COMPLETED"]
    DECL --> RELEASE["Освободить резерв; подтвердить платёжный исход"]
    REF --> SAFE{"Подтверждены оба условия отмены?"}
    RELEASE --> SAFE
    SAFE -->|Да| CANCELLED["Orders: CANCELLED"]
    SAFE -->|Пока нет| WAIT["CANCELLING: повтор и сверка"]
    WAIT --> SAFE
```

`ACCEPTED` — подтверждение сохранения запроса. `COMPLETED` — подтверждённая покупка: деньги списаны на ожидаемую сумму, резерв превращён в продажу. Ветка клиентской отмены подробно показана в схемах 6–7.

#### 2. Успешный заказ: взаимодействия и транзакции

Вопрос для ревью: какие данные должны пережить падение после каждой стрелки?

```mermaid
sequenceDiagram
    autonumber
    participant C as Клиент
    participant O as Orders
    participant I as Inventory
    participant K as Kafka
    participant P as Payments

    C->>O: gRPC CreateOrder(request_key, items)
    Note over O: TX Orders: ACCEPTED + request_key + reserve_job
    O-->>C: order_id, ACCEPTED

    O->>I: gRPC Reserve(order_id, operation_id)
    Note over I: TX Inventory: резерв + ledger + InventoryReserved в outbox
    I-->>O: reservation_id, снимок суммы
    Note over O: TX Orders: PAYMENT_PENDING + OrderAwaitingPayment в outbox
    O-->>K: publish OrderAwaitingPayment
    K-->>P: delivery OrderAwaitingPayment
    Note over P: TX Payments: inbox + payment_job
    P-->>K: commit offset после TX

    Note over P: Воркер симулирует задержку вне SQL-транзакции
    Note over P: TX Payments: единственный debit + SUCCEEDED + outbox
    P-->>K: publish PaymentSucceeded
    K-->>O: delivery PaymentSucceeded
    Note over O: TX Orders: inbox + FINALIZING + commit_job
    O-->>K: commit offset после TX

    O->>I: gRPC CommitReservation(operation_id)
    Note over I: TX Inventory: COMMITTED + ledger + outbox
    I-->>O: COMMITTED
    I-->>K: publish InventoryCommitted
    K-->>O: delivery InventoryCommitted
    Note over O: TX Orders: inbox + COMPLETED + OrderCompleted в outbox
    O-->>K: publish OrderCompleted

    C->>O: gRPC GetOrder(order_id)
    O-->>C: COMPLETED
```

Publisher каждого сервиса работает из его outbox. Воркеры работают из сохранённых jobs. Они изображены внутри своего владельца, чтобы на схеме оставались ровно три бизнес-сервиса.

Inventory также публикует `InventoryReserved` и `InventoryReleased`; Orders использует их для наблюдения и сверки. Они опущены на успешной sequence-схеме, чтобы выделить события, непосредственно продолжающие этот путь.

#### 3. Состояния заказа

Вопрос для ревью: какие переходы являются допустимыми, а какие нужно запретить?

```mermaid
stateDiagram-v2
    [*] --> ACCEPTED: Заказ и задача сохранены
    ACCEPTED --> RESERVING: Воркер начал шаг
    RESERVING --> PAYMENT_PENDING: Резерв и сумма подтверждены
    RESERVING --> REJECTED: Подтверждена нехватка товара
    PAYMENT_PENDING --> FINALIZING: Успешная оплата; отмена не принята
    FINALIZING --> COMPLETED: InventoryCommitted; оплата подтверждена

    ACCEPTED --> CANCELLING: Принята отмена
    RESERVING --> CANCELLING: Принята отмена
    PAYMENT_PENDING --> CANCELLING: Отмена или отказ оплаты
    CANCELLING --> CANCELLED: Резерв безопасен; списания нет или оно возвращено

    RESERVING --> RESERVING: Неизвестный RPC; повтор того же ключа
    PAYMENT_PENDING --> PAYMENT_PENDING: Ожидание и сверка платежа
    FINALIZING --> FINALIZING: Повтор Commit и ожидание события
    CANCELLING --> CANCELLING: Повтор компенсации и сверка

    REJECTED --> [*]
    COMPLETED --> [*]
    CANCELLED --> [*]

    note right of FINALIZING
        Точка прекращения клиентской отмены.
        CancelOrder отклоняется.
        Временный сбой сохраняет FINALIZING.
    end note
```

Timeout, недоступность зависимости и рост Kafka lag не превращают заказ автоматически в REJECTED или CANCELLED. Пока результат не разрешён, заказ остаётся в соответствующем нетерминальном состоянии. Если позже добавляется NEEDS_REVIEW, он тоже учитывается как нерешённый процесс.

##### 3.1. Состояния резерва

```mermaid
stateDiagram-v2
    state "Записи ещё нет" as ABSENT
    [*] --> ABSENT
    ABSENT --> RESERVED: Reserve успешно сохранён
    ABSENT --> RELEASED: Release сохраняет запрет позднего Reserve
    RESERVED --> RESERVED: Повтор того же Reserve
    RESERVED --> COMMITTED: Commit фиксирует продажу
    RESERVED --> RELEASED: Release освобождает товар
    COMMITTED --> COMMITTED: Повтор Commit без второго эффекта
    RELEASED --> RELEASED: Повтор Release без второго эффекта

    note right of RELEASED
        Поздний Reserve запрещён.
        Commit из RELEASED отклоняется.
    end note

    note right of COMMITTED
        Release из COMMITTED отклоняется.
        Возврат проданного товара вне MVP.
    end note
```

##### 3.2. Состояния платежа

```mermaid
stateDiagram-v2
    state "Платежа ещё нет" as ABSENT
    [*] --> ABSENT
    ABSENT --> PROCESSING: Событие сохранено в inbox и job
    ABSENT --> CANCELLED: Cancel сохраняет запрет позднего списания
    PROCESSING --> SUCCEEDED: Единственный debit сохранён
    PROCESSING --> DECLINED: Подтверждён отказ
    PROCESSING --> CANCELLED: Cancel победил до фиксации debit
    SUCCEEDED --> REFUND_PENDING: Принят возврат
    REFUND_PENDING --> REFUNDED: Полный refund подтверждён
    REFUND_PENDING --> REFUND_PENDING: Повтор возврата с тем же ключом
    SUCCEEDED --> SUCCEEDED: Повтор исходной операции
    CANCELLED --> CANCELLED: Позднее событие не создаёт debit
    REFUNDED --> REFUNDED: Повтор платежа или возврата без нового эффекта

    note right of SUCCEEDED
        Cancel не стирает списание.
        Для отмены нужен Refund.
    end note
```

Задержка симулятора происходит вне SQL-транзакции. Перед фиксацией debit воркер проверяет актуальное состояние: сохранённая отмена имеет силу и для уже начатой симуляции.

#### 4. Недостаточно товара

Вопрос для ревью: может ли оплата начаться до гарантированного резерва?

```mermaid
sequenceDiagram
    participant O as Orders
    participant I as Inventory
    O->>I: gRPC Reserve(order_id, operation_id, items)
    Note over I: Проверить весь заказ одной транзакцией
    Note over I: Товара недостаточно: складские изменения откатываются
    Note over I: Результат отказа operation_id сохраняется долговечно
    I-->>O: Подтверждённый business reject
    Note over O: TX Orders: REJECTED + результат шага
    Note over O: OrderAwaitingPayment не создаётся
```

Для многопозиционного заказа недостаток одной позиции отменяет весь резерв. При этом подтверждённый результат операции сохраняется: повтор того же operation_id возвращает тот же отказ, даже если остатки позже изменились. Это достигается проверкой всех заблокированных позиций до мутаций либо откатом складской части к savepoint с последующей фиксацией результата операции. Технический timeout при Reserve не равен подтверждённой нехватке товара.

#### 5. Платёж отклонён

Вопрос для ревью: когда можно завершить отмену, если освобождение склада временно не работает?

```mermaid
sequenceDiagram
    participant P as Payments
    participant K as Kafka
    participant O as Orders
    participant I as Inventory

    Note over P: TX Payments: DECLINED + PaymentDeclined в outbox
    P-->>K: publish PaymentDeclined
    K-->>O: delivery PaymentDeclined
    Note over O: TX Orders: inbox + CANCELLING + release_job
    O->>I: gRPC ReleaseReservation(order_id, operation_id)

    alt Release подтверждён
        Note over I: TX Inventory: RELEASED + ledger + outbox
        I-->>O: RELEASED
        Note over O: Списания нет; резерв освобождён
        Note over O: TX Orders: CANCELLED + outbox
    else Timeout или временная недоступность
        I-->>O: Результат неизвестен или временная ошибка
        Note over O: Оставить CANCELLING; сохранить следующий повтор
        O->>I: Повтор Release с тем же operation_id
        I-->>O: Прежний или вновь подтверждённый RELEASED
        Note over O: TX Orders: CANCELLED после проверки условий
    end
```

Ветка повторов на схеме сокращена. Если зависимость остаётся недоступной, новые попытки не объявляются успешными: процесс остаётся CANCELLING и попадает в отчёт о незавершённых заказах.

#### 6. Гонка отмены и успешной оплаты

Вопрос для ревью: где принимается решение, исключающее одновременные Commit и компенсацию?

```mermaid
flowchart TB
    CANCEL["CancelOrder"] --> CAS["Orders DB: условный переход текущей версии"]
    PAID["PaymentSucceeded"] --> CAS
    CAS --> WHO{"Какой переход победил?"}

    WHO -->|Финализация| F["PAYMENT_PENDING → FINALIZING"]
    F --> DENY["CancelOrder: отмена отклонена"]
    F --> CJ["Сохранить commit_job"]
    CJ --> IC["Inventory: COMMITTED"]
    IC --> OK["Orders: COMPLETED"]

    WHO -->|Отмена| C["Допустимое состояние → CANCELLING"]
    C --> CANCELJOB["Сохранить задачи отмены платежа и освобождения резерва"]
    CANCELJOB --> STATE{"Платёж уже списан?"}
    STATE -->|Нет| BLOCK["CancelPayment: сохранить запрет списания"]
    STATE -->|Да| RF["RefundPayment: подтвердить возврат"]
    BLOCK --> PROOF["Проверить платёжный исход и RELEASED резерв"]
    RF --> PROOF
    PROOF --> CC["Orders: CANCELLED"]
```

Условный переход — короткая транзакция в Orders по актуальному состоянию/версии. Проигравший обработчик перечитывает результат и следует выбранной ветке. Принятая отмена не создаёт commit_job; выбранная FINALIZING не создаёт release_job.

На практике CANCELLING может начинаться и до PAYMENT_PENDING. Запрет позднего Reserve в Inventory и позднего списания в Payments защищает от запросов и событий, отправленных старым воркером до отмены.

#### 7. Поздняя успешная оплата после принятой отмены

Вопрос для ревью: достаточно ли игнорировать запоздавшее событие, если деньги уже были списаны?

```mermaid
sequenceDiagram
    participant C as Клиент
    participant O as Orders
    participant P as Payments
    participant I as Inventory
    participant K as Kafka

    Note over P: Платёж SUCCEEDED; debit уже сохранён
    P-->>K: publish PaymentSucceeded
    Note over K,O: Доставка в Orders задержана

    C->>O: gRPC CancelOrder(order_id)
    Note over O: TX Orders: CANCELLING + cancellation jobs
    O-->>C: Отмена принята; компенсации ещё выполняются
    O->>P: gRPC CancelPayment(order_id, operation_id)
    P-->>O: SUCCEEDED: нужен возврат
    O->>I: gRPC ReleaseReservation(order_id, operation_id)
    I-->>O: RELEASED подтверждён
    O->>P: gRPC RefundPayment(order_id, operation_id)
    Note over P: Возврат выполняется через сохранённый шаг
    Note over P: TX Payments: refund ledger + REFUNDED + outbox
    P-->>O: REFUNDED подтверждён
    Note over O: TX Orders: оба условия подтверждены; CANCELLED + outbox

    K-->>O: delivery старого PaymentSucceeded
    Note over O: Записать inbox; сохранить ветку отмены
    Note over O: Не создавать Commit; возврат уже подтверждён
    P-->>K: publish PaymentRefunded
    K-->>O: delivery PaymentRefunded
    Note over O: Идемпотентно подтвердить уже известный исход
```

На этой схеме Refund успевает подтвердиться через gRPC. При неизвестном результате возврата Orders остаётся CANCELLING и использует повтор с тем же ключом, GetPayment или событие PaymentRefunded. Одного намерения возврата недостаточно для CANCELLED.

#### 8. События и падения: outbox, inbox и jobs

Вопрос для ревью: что уже сохранено в каждом crash window и что будет повторено?

##### 8.1. Запись и публикация

```mermaid
flowchart TB
    A["Бизнес-операция владельца"] --> TX
    subgraph TX["Одна SQL-транзакция производителя"]
        B["Сохранить бизнес-состояние"] --> E["Сохранить событие в outbox"]
    end
    TX --> COMMIT["SQL commit"]
    COMMIT --> PUB["Внутренний publisher читает outbox"]
    PUB --> K["Kafka подтверждает публикацию"]
    K --> SENT["Пометить outbox как отправленный"]

    COMMIT -.-> CR1["Падение до publish: событие осталось в БД"]
    CR1 --> PUB
    K -.-> CR2["Падение до sent marker: публикация повторится"]
    CR2 --> PUB
```

Crash до SQL commit не сохраняет ни бизнес-изменение, ни событие. Crash после commit оставляет событие для publisher. Повторная публикация после успешной Kafka-записи ожидаема; event_id сохраняется.

##### 8.2. Приём события и долгая работа

```mermaid
flowchart TB
    K["Kafka доставляет event_id"] --> CHECK{"event_id уже в inbox?"}
    CHECK -->|Да| OLD["Нового job нет; прежний результат сохранён"]
    CHECK -->|Нет| TX
    subgraph TX["Одна SQL-транзакция потребителя"]
        IN["Вставить inbox с unique event_id"] --> JOB["Сохранить переход и durable job"]
    end
    TX --> COMMIT["SQL commit"]
    COMMIT --> OFF["Подтвердить Kafka offset"]
    OLD --> OFF
    COMMIT --> WORK["Внутренний воркер забирает job по lease"]
    WORK --> RUN["Долгая работа вне SQL-транзакции"]
    RUN --> EFFECT["Короткая TX: проверить состояние и operation_id; сохранить эффект и outbox"]
    EFFECT --> DONE["Job завершён; результат долговечен"]

    COMMIT -.-> CR["Падение до offset: Kafka доставит дубликат"]
    CR --> K
    WORK -.-> WCR["Падение воркера: lease истечёт"]
    WCR --> WORK
```

Ромб inbox — логическая схема, а не отдельный небезопасный SELECT перед INSERT. Проверка дубликата и создание job защищены unique constraint и одной транзакцией, включая конкурентную доставку.

Offset подтверждает, что продолжение работы сохранено, а не что платёж уже завершён. Поэтому нужны два вида наблюдаемости: Kafka lag и возраст/размер durable jobs. Повтор job использует прежний operation_id и повторно проверяет актуальное состояние перед эффектом.

#### 9. Неизвестный результат gRPC

Вопрос для ревью: означает ли timeout, что Inventory не успел изменить остатки?

```mermaid
sequenceDiagram
    participant O as Orders
    participant I as Inventory
    participant D as Inventory DB

    O->>I: Reserve(order_id, operation_id)
    I->>D: TX: резерв + операция + ledger + outbox
    D-->>I: COMMIT выполнен
    I--xO: Ответ потерян или deadline истёк
    Note over O: Результат неизвестен; не считать Reserve неуспешным

    O->>I: Повтор Reserve с тем же operation_id
    I->>D: Найти сохранённый результат операции
    D-->>I: Прежний reservation_id и сумма
    I-->>O: Тот же подтверждённый результат
    Note over O: Продолжить сохранённый шаг
```

Вместо повтора можно начать со сверки GetReservation. Но NOT_FOUND во время ещё выполняющегося запроса не доказывает отсутствие будущего эффекта. Безопасность обеспечивают прежний ключ операции, атомарность и cancellation barrier, если уже выбран путь отмены.

Точно так же timeout CreateOrder не означает, что заказа нет: клиент повторяет request_key или вызывает GetOrderByRequestKey. Timeout Commit и Refund требует сверки прежнего эффекта, а не новой операции с новым ключом.

#### 10. Чтение через Redis и запись в PostgreSQL

Вопрос для ревью: где допустима устаревшая информация и где нужен авторитетный результат?

```mermaid
flowchart TB
    READ["gRPC GetAvailability"] --> CACHE{"Redis: кеш доступен и ключ найден?"}
    CACHE -->|Да| ADVICE["Вернуть ориентировочную доступность"]
    CACHE -->|Нет или ошибка| DBREAD["Прочитать Inventory DB"]
    DBREAD --> FILL["Best-effort записать кеш с TTL"]
    FILL --> ADVICE

    WRITE["gRPC Reserve"] --> DBTX["Inventory DB: атомарно проверить и зарезервировать"]
    DBTX --> ANSWER["Авторитетный результат Reserve"]
    DBTX -.-> INVALID["Best-effort invalidation кеша после commit"]
```

Устаревший кеш может показать товар доступным, а Reserve затем корректно откажет. Это измеряемый stale read, не продажа лишнего товара. Ошибка кеша не должна отменять уже выполненную складскую транзакцию; fallback и его нагрузка на SQL становятся отдельным экспериментом.

#### 11. Цикл инженерного эксперимента

Вопрос для ревью: какие данные позволяют принять решение об эффективности изменения?

```mermaid
flowchart TB
    H["Проблема → гипотеза → критерии до запуска"] --> FIX["Зафиксировать версии, данные, seed и ресурсы"]
    FIX --> A["Прогон A: baseline"]
    A --> CHANGE["Изменить один проверяемый фактор"]
    CHANGE --> B["Прогон B: candidate"]
    B --> REPEAT["Повторить пары A/B; собрать разброс"]
    REPEAT --> DRAIN["Остановить вход; дождаться обработки или лимита drain"]
    DRAIN --> VER["Сверить все request keys, ledger и состояния"]
    VER --> DATA{"Данные и нагрузка достаточны?"}
    DATA -->|Нет| INC["INCONCLUSIVE: объяснить пробел и повторить"]
    DATA -->|Да| INV{"Инварианты соблюдены?"}
    INV -->|Нет| BAD["REJECTED: исправить нарушение"]
    INV -->|Да| MET["Сравнить goodput, latency, backlog, recovery и ресурсы"]
    MET --> REPORT["Отчёт: принять, отклонить или продолжить измерение"]
    REPORT --> NEXT["Следующий эксперимент"]
```

Каждый отдельный прогон A или B включает свой прогрев, измеряемое окно, drain и verifier. На схеме общие шаги после серии обозначают сбор уже проверенных результатов. Незавершённые заказы не исчезают из знаменателя; они остаются отдельным результатом.

#### 12. Процесс ревью G0

```mermaid
flowchart LR
    P["Предложение G0: архитектура и схемы"] --> R["Ревью владельца проекта"]
    R --> D{"Решение"}
    D -->|Изменить| C["Уточнить решения и обновить схемы"]
    C --> R
    D -->|Согласовано| G["Зафиксировать принятую версию G0"]
    G --> CODE["Первые небольшие пакеты кода для ревью"]
    CODE --> G1["G1: основа заказа и склада"]
```

G0 прошёл ветку согласования: пользователь полностью принял его 2026-10-13. Схема остаётся описанием процесса ревью и для последующих изменений. После замечаний схемы, контракты и бизнес-правила обновляются вместе; существенные изменения согласуются отдельно.

#### Контрольные вопросы владельцу проекта

| Область | Что должно быть однозначно после ревью |
|---|---|
| Успех | ACCEPTED отличается от COMPLETED; завершение требует оплаты и продажи |
| Отмена | Понятна точка FINALIZING и ветка позднего успеха оплаты |
| Корректность | Повторы не создают второй эффект; отменённый процесс не воскресает |
| Долговечность | Видно, что сохраняется перед ответом клиента и Kafka offset |
| Восстановление | Timeout остаётся неизвестным исходом до сверки; jobs/outbox переживают рестарт |
| Масштабирование | Приём, Kafka delivery и processing jobs имеют отдельные очереди и показатели |
| Кеш | Redis ускоряет чтения; складская запись проверяет PostgreSQL |
| Эксперимент | Вывод учитывает незавершённые процессы, инварианты и ресурсный бюджет |

<a id="appendix-development-process"></a>

## Приложение: Роли, ревью, контрольные точки и обновление контекста

Полный исходный документ: `docs/development/WORKFLOW.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Порядок разработки и ревью

#### Роли

Пользователь — владелец проекта: определяет приоритеты, утверждает ключевые решения, проводит код-ревью и согласует переходы между основными этапами.

Единое название роли человека во всех материалах — «владелец проекта»; при описании проверки кода используется «ревьюер». Эти обозначения сохраняют прежние полномочия и порядок согласования.

Codex (AI-Your-Bro) — исполнитель: готовит проектные решения, реализует согласованный объём, выполняет необходимые проверки, собирает доказательства корректности и эффективности, документирует результат и исправляет замечания.

Это распределение ролей задано пользователем. Архитектура и процесс G0 полностью приняты владельцем проекта 13 октября 2026 года; контрольные точки ниже входят в текущий согласованный порядок работы и могут быть изменены владельцем проекта.

#### Пакет изменений

Каждый пакет имеет одну понятную цель: например, атомарный Reserve, durable payment job или приём конкретного Kafka-события. Работа внутри согласованного этапа выполняется самостоятельно. Размер пакета выбирается так, чтобы владелец проекта мог проверить бизнес-поведение и существенный diff за одно ревью.

Перед началом каждой задачи исполнитель кратко обозначает ожидаемый результат и границы, основные части реализации и проверок, возможные трудности и зависимости. Если есть основания, даёт ориентировочный диапазон времени с допущениями, отдельно отмечая ожидание сборок и CI; при недостатке данных объясняет неопределённость. При существенном изменении объёма или появлении нового препятствия обновляет оценку и объясняет причину. Правило согласовано владельцем в чате 16 октября 2026 года; для уже порученной работы дополнительное согласование старта не требуется.

Для каждого пакета Codex представляет:

1. Задачу и соответствующий этап roadmap.
2. Изменённое поведение и существенные проектные решения.
3. Локальный diff либо PR, когда для проекта доступен настроенный remote.
4. Команды воспроизведения и фактические результаты проверок.
5. Известные ограничения и открытые вопросы, влияющие на принятие изменения.
6. Критерии готовности и предложение следующего пакета.

Обычные исправления по ревью выполняются в рамках того же пакета. После исправлений повторяются проверки, затронутые изменением. Невыполненные проверки явно обозначаются как невыполненные; объясняется причина и что необходимо для их запуска.

[Спецификации пакетов](#appendix-p02) описывают границы и критерии. Актуальные задачи, блокеры и ревью ведутся только в [GitHub](https://github.com/Includeoyasi/faultline-lab/issues); документы не дублируют текущий прогресс.

#### Стандарт код-ревью

[Правила код-ревью](#appendix-code-review) определяют Go-стиль, обработку ошибок, context, конкурентность, работу с данными, доказательства тестирования и порядок разбора замечаний. [Шаблон PR](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/.github/pull_request_template.md) содержит краткий применимый чек-лист.

После bootstrap изменения готовятся в отдельных ветках и предъявляются через PR. Исполнитель делает self-review и предоставляет фактические результаты; владелец проекта рассматривает предъявленный HEAD. ИИ не одобряет и не сливает собственный PR. Успешный CI подтверждает выполненные проверки; решение человека фиксируется отдельно. Обычный путь после приёмки — squash merge по явному поручению владельца проекта или действие самого владельца проекта.

По поручению владельца новые PR оформляются от исполнителя AI-Your-Bro, а владелец проекта проводит ревью и принимает решение. Защита main уже включена; авторизация исполнителя и права write без admin проверены через GitHub API. Новые PR от владельца не создаются; Git author metadata не заменяет проверку API actor. PR #4, созданный до перехода, уже слит (merge 8c392b8); его переоформление не требуется. CODEOWNERS действует в main после owner merge PR #5 и назначает владельца ревьюером всех файлов. Точный статус — в [безопасности репозитория](#appendix-repository-security).

Принятие в чате фиксируется в контексте, но не заменяет GitHub approval защищённой ветки. Обязательны актуальные Documentation и Go toolchain (AMD64/ARM64) checks, независимое owner review и остальные требования защиты main.

#### Контрольные точки

G0–G4 — контрольные точки приёмки. P01–P05 — небольшие пакеты реализации с отдельными PR и код-ревью. **Первый PR P01 находится между принятым G0 и будущей приёмкой G1.** Последовательность: G0 → P01 и его ревью → следующие PR P02–P05 и их ревью → G1. Приёмка отдельного PR не означает прохождения G1.

| Точка | Когда | Материалы для согласования | Критерий перехода |
|---|---|---|---|
| G0: архитектура | Перед реализацией | docs/architecture/DESIGN.md, схемы docs/architecture/PROCESSES.md, архитектура и состояния из docs/ROADMAP.md, границы MVP, порядок реализации | Владелец проекта согласовал предложенные решения или сформулировал необходимые изменения |
| G1: основа данных и операций | После этапа 1 | Код и миграции; API заказа и склада; повторный запрос; конкуренция за остатки; сохранённый шаг после рестарта | Ревью принято, обязательные проверки проходят, нет нерешённых нарушений складских инвариантов |
| G2: сквозной MVP | После этапа 2 | Сквозной процесс трёх сервисов, Kafka/outbox/inbox, платежи и компенсации, Redis и наблюдаемость, первый эксперимент и отчёт | Подтверждены успешный процесс, отказ оплаты, идемпотентность и контрольный перезапуск; каждый тестовый заказ учтён |
| G3: устойчивость | После этапа 4 | Сценарии неизвестного результата, поздних событий, отмены, replay; объяснение политики ошибок и recovery | Ревью и сценарии приняты; нерешённые процессы явно учитываются; отсутствуют ложные terminal states |
| G4: готовая лаборатория | После этапа 7 | Каталог экспериментов, результаты сравнения решений, масштабирование, автоматический запуск/сверка/отчёты, инструкции | Владелец проекта воспроизвёл или проверил заявленный цикл; отчёты отражают корректность, ресурсы и ограничения |

G0 полностью принят владельцем 13 октября 2026 года. Основания приёмки сохраняются как датированные факты; актуальные решения по следующим этапам смотрите в GitHub.

На контрольной точке Codex сначала доводит текущий пакет до проверяемого результата, затем предоставляет его владельцу проекта. Следующий этап, зависящий от согласования, начинается после явного решения. Замечания владельца проекта вносятся в пакет до повторного рассмотрения.

#### Изменения, требующие отдельного согласования

- Изменение ответственности сервисов, бизнес-инвариантов или момента окончательного завершения/отмены заказа.
- Несовместимое изменение API, событий, данных или политики идемпотентности.
- Добавление бизнес-сервиса либо инфраструктуры, существенно меняющей объём или среду запуска.
- Изменение утверждённых критериев эксперимента, которое влияет на вывод об эффективности.

Настройки, рефакторинг и необходимые проверки в рамках согласованного решения выполняются самостоятельно. Если найденный дефект требует смены архитектуры, Codex готовит конкретное предложение и объясняет, какой инвариант или сценарий заставляет пересмотреть решение.

#### Порядок принятия решений

План или гипотеза помечаются как предложенные. Согласование владельцем проекта фиксируется отдельно от фактической реализации и результатов проверки. Успешные тесты подтверждают проверенные свойства, а не заменяют решение ревьюера.

Решения, меняющие проект, сохраняются в ADR с причиной и последствиями. Итог пакета содержит текущий статус контрольной точки, оставшиеся замечания и ссылку на воспроизводимый результат.

По согласованному [ADR 0002](#appendix-adr-0002) рабочее обсуждение сохраняется в PR, существенное решение — в ADR, а отобранные примеры — на странице [«Инженерные решения и ревью»](#appendix-engineering-decisions). Запись связывает вопрос, варианты и причину выбора, изменение, доказательства и статус. Для рутинного исправления достаточно PR; отдельный ADR не требуется.

Замечания из чата пересказывает исполнитель от AI-Your-Bro, указывая владельца как автора требования и чат как источник. Он не создаёт комментарии или approvals от имени владельца. Согласование направления, готовая реализация, результаты тестов и приёмка конкретного PR фиксируются отдельно. Обзор и переносимый контекст обновляются при изменении существенного решения или его статуса.

[Спецификация Go smoke](#appendix-go-smoke-backlog) и [ADR 0003](#appendix-adr-0003) описывают выбранный подход; выполнение и ревью — в [Issue #9](https://github.com/Includeoyasi/faultline-lab/issues/9).

Исходные оценки часов в docs/ROADMAP.md описывают объём ручной разработки. Для работы Codex календарный прогноз уточняется по фактическим ограничениям среды, времени проверок и ревью. Контрольная точка достигается по критериям готовности.

#### Актуальность переносимого контекста

PROJECT_CONTEXT.md — самостоятельный пакет для другого ИИ-агента или чата без доступа к рабочему компьютеру. В нём есть краткое введение, датированная история решений и полные тексты проектных документов.

Каждый содержательный пакет работы завершать обновлением контекста: что изменено, какие решения приняты и что реально проверено, с датой, средой и commit. Актуальные задачи, стадии и блокеры обновляются только в GitHub. Обновить CONTEXT_STATE.json, соответствующие исходные документы и при необходимости tools/context_intro.md; увеличить revision и поставить дату из актуального пользовательского контекста Europe/Moscow. Затем выполнить tools/update_context.ps1 и проверку tools/update_context.ps1 -CheckOnly.

Проверка выявляет несоответствие переносимого файла его локальным источникам. Достоверность датированных утверждений и результатов проверок обеспечивается фактическими данными исполнителя; сама генерация не запускает тесты сервисов.

Перед передачей другому агенту использовать последнюю проверенную ревизию PROJECT_CONTEXT.md. Переданная копия является снимком: изменения, выполненные в другом чате, нужно вернуть вместе с обновлением контекста и согласовать с более новой локальной версией. Агент без доступа к файлам не помечает предложения как применённый код и проверки как выполненные.

<a id="appendix-report-template"></a>

## Приложение: Полный шаблон отчёта

Полный исходный документ: `docs/experiments/REPORT_TEMPLATE.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Отчёт об инженерном эксперименте

> Шаблон. Значения ниже заполняются после реальных прогонов. Пустые поля не являются измеренными нулевыми значениями.

#### 1. Вывод

- ID эксперимента: `<id>`
- Решение: `<accepted / rejected / inconclusive>`
- Проверяемое изменение: `<одно изменение>`
- Наблюдаемый эффект: `<абсолютная разница, процент при ненулевой базе, разброс>`
- Корректность: `<инварианты; unresolved/unknown; полнота данных>`
- Цена и регрессии: `<ресурсы, отказы, сложность, задержка>`

#### 2. Проблема и гипотеза

**Проблема:** `<бизнес-симптом и подтверждающий технический сигнал>`

**Гипотеза:** `<изменение X должно улучшить показатель Y за счёт механизма Z>`

**Критерии, заданные до запуска:**

- Минимально полезный эффект: `<порог>`
- Допустимые регрессии: `<пороги>`
- Инварианты: `<нулевая терпимость или конкретный проверяемый контракт>`
- Допустимое время восстановления: `<порог и определение восстановления>`
- Условия непригодного измерения: `<насыщение генератора, пропуски данных, различие среды>`

#### 3. Условия

| Параметр | A: baseline | B: candidate |
|---|---|---|
| Commit / dirty diff | `<значение>` | `<значение>` |
| Образы и digests | `<значение>` | `<значение>` |
| Суммарные CPU/RAM и Docker limits | `<значение>` | `<значение>` |
| Реплики, worker limits, SQL pools | `<значение>` | `<значение>` |
| Kafka topics/partitions/groups | `<значение>` | `<значение>` |
| Dataset, stock, seed | `<значение>` | `<значение>` |
| Offered arrival rate и read/order mix | `<значение>` | `<значение>` |
| Fault timeline: план и факт | `<значение>` | `<значение>` |
| Warmup / measure / drain | `<значение>` | `<значение>` |
| Повторы и порядок A/B | `<значение>` | `<значение>` |
| Telemetry / sampling / verifier load | `<значение>` | `<значение>` |

**Точные run IDs и ссылки на manifest:** `<список>`

**Фактический diff:** `<что изменено; подтверждение отсутствия других изменений>`

#### 4. Результаты

Показывать отдельно результаты каждой пары прогонов. Агрегат должен включать объём выборки и разброс. Median p95 отдельных прогонов так и называется; это не p95 всех запросов.

| Показатель | A | B | Разница | Разброс / выборка |
|---|---|---|---|---|
| Offered / sent / skipped unique orders | `<…>` | `<…>` | `<…>` | `<…>` |
| Accepted / unknown ingress result | `<…>` | `<…>` | `<…>` | `<…>` |
| Completed в окне / во время drain | `<…>` | `<…>` | `<…>` | `<…>` |
| Business rejects / technical failures | `<…>` | `<…>` | `<…>` | `<…>` |
| Unresolved / unknown после drain | `<…>` | `<…>` | `<…>` | `<…>` |
| Completed goodput | `<…>` | `<…>` | `<…>` | `<…>` |
| E2E p50 / p95 / p99 для завершённых | `<…>` | `<…>` | `<…>` | `<…>` |
| Business success ratio среди sent eligible | `<…>` | `<…>` | `<…>` | `<…>` |
| gRPC errors / deadlines / retries | `<…>` | `<…>` | `<…>` | `<…>` |
| SQL latency / pool acquire / lock wait | `<…>` | `<…>` | `<…>` | `<…>` |
| Kafka lag и oldest event age | `<…>` | `<…>` | `<…>` | `<…>` |
| Job/outbox backlog и oldest age | `<…>` | `<…>` | `<…>` | `<…>` |
| Recovery time / drain time | `<…>` | `<…>` | `<…>` | `<…>` |
| CPU-seconds / RSS / disk / network | `<…>` | `<…>` | `<…>` | `<…>` |
| Completed orders / CPU-second | `<…>` | `<…>` | `<…>` | `<…>` |

**Определения знаменателей и временных окон:** `<eligible задано заранее; completion cohort; treatment of retries; time origin>`

**Цензурированные наблюдения:** `<число незавершённых и наблюдаемый возраст; они не исключены из success ratio>`

#### 5. Графики

Добавить PNG/SVG с единой шкалой времени и отметками начала/конца fault:

1. Offered/sent/accepted/completed throughput.
2. gRPC и business latency; рядом доля незавершённых.
3. Kafka, durable jobs и outbox backlog/age.
4. CPU, память, SQL pool и lock waits.
5. Восстановление после снятия отказа; firing/resolved алертов.

**Ссылки на сохранённые графики, сырые данные и Grafana time window:** `<ссылки>`

#### 6. Корректность

| Проверка | A | B | Доказательство |
|---|---|---|---|
| Неотрицательные остатки и баланс available/reserved/sold | `<…>` | `<…>` | `<invariants.json / ledger>` |
| Не более одного debit на order_id | `<…>` | `<…>` | `<ledger>` |
| Возвраты не превышают списания | `<…>` | `<…>` | `<ledger>` |
| COMPLETED имеет оплату и committed reserve | `<…>` | `<…>` | `<сверка>` |
| CANCELLED не имеет резерва или невозвращённого debit | `<…>` | `<…>` | `<сверка>` |
| Все request keys учтены, включая unknown ingress | `<…>` | `<…>` | `<requests.jsonl + orders>` |
| Поздние/повторные события не создали повторных эффектов | `<…>` | `<…>` | `<operation keys / history>` |

**Нерешённые процессы:** `<order_id, состояние, возраст, причина; либо подтверждённое отсутствие>`

**Модель отказов и границы проверки:** `<сохранены ли диски, ограничен ли retention, работала ли сверка>`

#### 7. Интерпретация

- Подтверждает ли сигнал предполагаемый механизм улучшения: `<объяснение>`.
- Изменился ли ресурсный бюджет: `<объяснение>`.
- Скрывает ли хороший percentile больше отказов или незавершённых заказов: `<проверка>`.
- Насколько эффект больше межпрогонного разброса: `<оценка без необоснованной статистической уверенности>`.
- Что не удалось доказать: `<ограничения локальной среды, данных и длительности>`.
- Цена поддержки изменения: `<код, эксплуатация, настройка>`.

#### 8. Решение и следующий эксперимент

**Решение:** `<принять / отклонить / увеличить объём измерений>`.

**Причина:** `<сопоставление с заранее выбранными критериями>`.

**Следующий эксперимент:** `<гипотеза, изменяемый фактор, ожидаемый сигнал>`.

<a id="appendix-agents"></a>

## Приложение: Локальные инструкции ИИ-исполнителю

Полный исходный документ: `AGENTS.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Instructions for this project

- The human is the project owner and reviewer. Codex, authenticated as AI-Your-Bro, is the executor. Other AI agents are positioned as potential users, subscribers or testers; do not assign them implementation work.
- Read PROJECT_CONTEXT.md for specifications, decisions and dated evidence, docs/development/WORKFLOW.md for review rules, and docs/ROADMAP.md for scope. Live task, blocker and review status belongs only in GitHub Issues/PRs/milestones; do not maintain duplicate progress boards in Markdown or JSON.
- Follow the latest explicit owner instruction. A planning discussion alone does not authorize publication or implementation. On 2026-10-16, after merging PR #21, the owner explicitly assigned issue #12 (PostgreSQL foundation) and requested detailed explanations and interaction diagrams. The owner also requested removal of obsolete test images and prevention of their accumulation. Issue #16 and the remaining audit findings retain their agreed checkpoints.
- Verify the authenticated actor before every GitHub mutation/publication: it must be AI-Your-Bro. Git author metadata is insufficient. Use the project owner's credentials only for explicitly authorized repository administration, never to publish or approve executor changes.
- main requires independent owner approval, approval after the last push, dismissal of stale reviews, current successful Documentation and Go toolchain AMD64/ARM64 checks, resolved conversations and administrator enforcement. Do not weaken protection, force-push main or approve/merge your own changes. CODEOWNERS assigns Includeoyasi. See docs/development/REPOSITORY_SECURITY.md; the owner declined 2FA and must not be asked again without new cause.
- G0 was fully approved on 2026-10-13, Europe/Moscow: “G0 полностью принят.” Approved service boundaries and invariants are in docs/architecture/DESIGN.md and PROCESSES.md. Do not ask for G0 approval again. G0–G4 are acceptance gates; P01–P05 are implementation packets. P01 is between G0 and G1; acceptance of a PR does not automatically accept a gate.
- Work autonomously within the assigned scope. Prepare a concrete reviewable result, applicable checks and limitations before asking for gate approval. Significant changes to service boundaries, invariants, incompatible contracts or infrastructure scope require a concrete proposal and an owner decision.
- Before starting each implementation task, give the owner a short scope estimate: intended result and boundaries, main implementation/checks, likely difficulties and dependencies, and a qualified elapsed-time range when supportable, distinguishing build/CI waits. State uncertainty instead of inventing precision; update the estimate when a material surprise changes it. This owner instruction was confirmed in chat on 2026-10-16 and does not add an approval step for already authorized work.
- Keep exactly three business services: Orders, Inventory and Payments. Workflow orchestration lives inside Orders; do not add a business service for it. Lab tools and adapters are separate from the tested business.
- Follow docs/development/CODE_REVIEW.md. Address lifecycle, errors/context, concurrency, data ownership and meaningful tests. Never report a test, measurement or approval that did not occur. Compile-only results do not establish runtime compatibility.
- Addresses belong to deployment configuration. FL_GRPC_ADDR is required, with no hidden listen address or port fallback. Keep listener, peer DNS and host publication separate; host publication must bind loopback explicitly. Test loopback fixtures and observed addresses are not production defaults.
- Use clear PR headings such as “Что было сделано”, “Проверки” and “На что обратить внимание”. Explain observable behavior and evidence. Do not use the rejected question-style heading.
- Working review discussion belongs in PRs, substantial decisions in ADRs and selected examples in docs/development/ENGINEERING_DECISIONS.md. Attribute chat requirements to the owner and identify chat as their source; do not impersonate owner GitHub comments or approval.
- After a substantive decision, implementation result or environment check, update the canonical specifications and dated evidence in CONTEXT_STATE.json, increment its revision and regenerate PROJECT_CONTEXT.md with tools/update_context.ps1, then run -CheckOnly. Do not edit only the generated file. Do not copy live task stages, counters or assignments into context.
- Snapshot resource links use repository.snapshot_ref, an existing full source commit SHA. Keep live navigation separate. Commit source materials first, then record their SHA and regenerate the context; do not invent a future self-referencing commit hash.
- Context supplied to an offline user/tester is a dated snapshot. They must not infer live task status from it or claim changes/checks on the source machine. New evidence is reconciled by revision and provenance.
- On the owner's personal Windows computer, stop task-owned containers, quit Docker Desktop and stop its WSL VM after checks. Keep Docker autostart disabled and verify shutdown. Remove disposable run-specific images, including after build failure/cancellation; verify their tags are gone. Preserve unrelated workloads, installed tools, current reusable images and useful caches. Never use global prune as test cleanup. This cleanup refinement follows the owner's explicit instruction on 2026-10-16.

<a id="appendix-adr-0001"></a>

## Приложение: Запись принятия архитектуры G0

Полный исходный документ: `docs/adr/0001-approved-g0.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### ADR 0001: принятие архитектуры G0

- Дата решения: 2026-10-13, Europe/Moscow.
- Статус: accepted.
- Кто принял: пользователь, владелец проекта.
- Прямое основание: сообщение пользователя «G0 полностью принят».

#### Контекст

Пользователь поручает ИИ разработку pet-проекта и сам проводит код-ревью и согласует ключевые этапы. Подготовлены roadmap, архитектурный пакет и визуальные схемы процессов. Реализация к моменту решения ещё не создана.

#### Принятое решение

Принят G0, описанный в docs/architecture/DESIGN.md, docs/architecture/PROCESSES.md, docs/ROADMAP.md и docs/development/WORKFLOW.md:

- Ровно три бизнес-сервиса Go: Orders, Inventory, Payments.
- Orders хранит долговечный процесс; четвёртый сервис-оркестратор не добавляется.
- gRPC — непосредственные команды/чтения и компенсации; Kafka — бизнес-события и асинхронная платёжная обработка.
- Сначала один PostgreSQL с тремя отдельными базами/ролями и владением данными внутри сервиса.
- Redis кеширует чтения Inventory; резервирование доверяет PostgreSQL.
- Outbox, inbox, operation keys и durable jobs входят в MVP.
- Бизнес-инварианты, граница клиентской отмены FINALIZING и компенсации поздней успешной оплаты приняты в описанном виде.
- Начальная среда — Docker Compose; Prometheus/Grafana и Loki/Alloy входят в MVP, расширенная трассировка — в следующий этап.
- Нагрузка, сбои, verifier, воспроизводимые A/B-эксперименты и отчётность — основная цель проекта.
- Работа предъявляется небольшими пакетами для ревью; следующие точки — G1–G4.

#### Последствия

Повторно запрашивать согласование G0 не нужно. Можно готовить согласованную основу для G1. Принятие G0 не заменяет код-ревью, проверки или согласование последующих точек.

Точные версии библиотек и образов, protobuf-файлы, миграции и конкретные параметры исполнения пока не выбраны или не реализованы. Их выбор выполняется внутри принятой архитектуры и предъявляется на ревью; существенное изменение архитектуры требует нового согласования.

В том же сообщении пользователь запросил подробный, переносимый и регулярно обновляемый файл контекста для других ИИ-агентов. Это отдельная задача документирования, выполненная до начала реализации сервисов.

<a id="appendix-contributing"></a>

## Приложение: Участие, ревью и проверки

Полный исходный документ: `CONTRIBUTING.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Contributing to Faultline Lab

This repository documents a project built with AI implementation and human technical leadership. Changes should be understandable to a reviewer and supported by reproducible evidence.

#### Before changing anything

Read the [decisions and evidence snapshot](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/PROJECT_CONTEXT.md), [approved design](#appendix-design-review) and [development workflow](#appendix-development-process). The executor, AI-Your-Bro, follows [AGENTS.md](#appendix-agents).

The [code review standard](#appendix-code-review) covers idiomatic Go, errors, context, concurrency, data integrity, tests and evidence. It also defines blocking comments, optional suggestions and the baseline Go checks to activate with the first module. The detailed standard is currently written in Russian.

G0 is approved. G1–G4 remain future review gates. Significant changes to approved service boundaries, invariants or incompatible contracts require a concrete proposal and the project owner's decision.

#### A reviewable change

1. Name the problem, intended behavior and roadmap packet.
2. Keep the diff focused; explain decisions that affect correctness or operability.
3. Run checks appropriate to the change. Record exact commands, observed results and checks that could not run.
4. Update affected specifications and dated decisions/evidence in `CONTEXT_STATE.json`; regenerate `PROJECT_CONTEXT.md`. Live task progress is recorded only in GitHub Issues and PRs.
5. Use the pull request template to present the result. Address review comments before moving through the next gate.

AI assistance should be disclosed in the review packet. The human reviewer remains responsible for accepting decisions. Never claim tests, approvals or measurements that did not occur.

#### Documentation checks

Requires Git with the referenced source commit available locally, Node.js 20+ and PowerShell 7+ for the commands below. A shallow checkout may need an explicit fetch of that SHA or full history. The generator also supports Windows PowerShell 5.1 (`powershell` instead of `pwsh`).

```powershell
node --test tools/doc_safety.test.mjs
node tools/check_docs.mjs
pwsh -NoProfile -File tools/update_context.ps1
pwsh -NoProfile -File tools/update_context.ps1 -CheckOnly
```

The checker validates local links, Markdown fence balance, text encoding and accidental machine-specific paths. The context check compares the generated snapshot with all registered sources. These checks do not validate service behavior or fully parse Mermaid syntax.

#### Branches and commits

Use `main` as the baseline, short descriptive branches such as `feat/orders-bootstrap`, and focused commits with clear intent. After bootstrap, submit changes through a pull request; the AI does not approve or merge its own work. The human reviews the current change and authorizes merging; passing CI alone is insufficient. Prefer squash merge after acceptance. Suggested prefixes: `feat:`, `fix:`, `docs:`, `test:`, `chore:`. Include both canonical document changes and their generated context in the same review packet.

See the [runtime guide](#appendix-container-runtime) and [Go smoke](#go-smoke-guide) for executable checks.

<a id="appendix-license"></a>

## Приложение: MIT License

Полный исходный документ: `LICENSE`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

MIT License

Copyright (c) 2026 Dmitrii Morozov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

<a id="appendix-code-review"></a>

## Приложение: Правила Go code review и критерии принятия изменений

Полный исходный документ: `docs/development/CODE_REVIEW.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Правила код-ревью Faultline Lab

Приняты через слияние [PR #1](https://github.com/Includeoyasi/faultline-lab/pull/1) владельцем проекта 13 октября 2026 года, Europe/Moscow. Принятый HEAD: `214eb34`; merge commit: `c1e55fd`. Документ задаёт рабочий стандарт для Go-кода и распределённых процессов проекта. Принятие правил не является приёмкой G1.

**Цель ревью:** понятное, корректное и сопровождаемое изменение с проверяемыми доказательствами. Сначала проверяем бизнес-поведение и отказные сценарии, затем устройство кода. Форматирование выполняют инструменты.

#### 1. Порядок ревью и слияния

1. Одна ветка и один PR решают одну понятную задачу. Связанный рефакторинг допустим; посторонние изменения выносятся отдельно. Для крупной работы сначала показать границы и порядок небольших пакетов, сохраняя работоспособность каждого.
2. После первоначального bootstrap изменения готовятся в ветках от `main`. Обычный путь — PR, проверки, ревью владельца проекта, затем squash merge. ИИ самостоятельно не одобряет и не сливает собственные изменения.
3. До ревью исполнитель проверяет diff, удаляет отладочные фрагменты и предоставляет проблему, результат, решения, команды проверок и ограничения. Исходники, миграции, контракты и соответствующие тесты рассматриваются вместе.
4. Блокирующие замечания устраняются либо явно разрешаются владельцем проекта с причиной. Новые изменения после ревью предъявляются повторно; существенная правка делает прежнее согласование недостаточным. Одобрение и успешные проверки должны относиться к предъявленному HEAD.
5. Слияние требует явного решения человека, успешных применимых проверок и актуального контекста. Не применять auto-merge до такого решения. Если человек поручает ИИ выполнить merge, это отдельное явное поручение.
6. G0 уже принят. Обычный PR не заменяет приёмку G1–G4; переходы через эти точки регулирует [workflow](#appendix-development-process). Существенная смена принятой архитектуры требует решения владельца проекта и ADR.

13 октября 2026 года по прямому поручению владельца включена и повторно проверена защита main: независимое approval, сброс устаревших ревью, одобрение после последнего push, актуальный успешный CI, закрытые обсуждения и применение к администратору; force push и удаление запрещены. Новые PR создаются от исполнителя AI-Your-Bro с подтверждённой авторизацией и правами write без admin; GitHub approval даёт владелец проекта. Подробности и подготовленный CODEOWNERS — в [безопасности репозитория](#appendix-repository-security). Сообщение в чате сохраняет решение человека, но не заменяет обязательное GitHub approval. ИИ не выдаёт approval от владельца или от второго аккаунта своим изменениям.

#### 2. Как оформлять замечания

| Метка | Смысл | Влияние на слияние |
| :--- | :--- | :--- |
| `blocking` | Дефект, риск потери данных, нарушение контракта/инварианта, существенный пробел проверок | Исправить либо получить конкретное решение владельца проекта |
| `question` | Нужно понять поведение или основание решения | Ответить; если обнаружен риск, явно сделать замечание блокирующим |
| `suggestion` | Улучшение с обоснованным эффектом | Обсудить; само по себе не блокирует |
| `nit` | Необязательная косметика или личное предпочтение | Не блокирует |

Комментарий указывает место, конкретный сценарий, последствие и ожидаемое свойство. Например: `blocking: offset подтверждается до сохранения job; crash в этом окне теряет оплату. Сохранить inbox и job одной транзакцией до ACK и проверить рестарт.` Обсуждаем код и наблюдаемое поведение; спор разрешаем правилами, фактами и воспроизводимым примером. Косметика вне стандарта не становится блокирующей только из-за предпочтения ревьюера.

#### 3. Стиль и простота Go

- Код проходит `gofmt`; импорты обрабатываются `goimports`. Названия пакетов описывают ответственность; избегаем универсальных `utils`, `common` и лишнего повторения имени пакета в типах.
- Приняты обычные Go-имена: `orderID`, `ServeHTTP`, `ErrNotFound`; локальные имена соразмерны области видимости. Идентификаторы и комментарии к Go API — на английском; подробная проектная документация может быть на русском.
- Ранние возвраты делают успешный путь прямым. Не вводим жёсткие лимиты строк, длины функций или числа параметров: блокируем сложность, мешающую понять поведение, и показываем её причину.
- Интерфейс появляется у потребителя, когда ему действительно нужна граница. Предпочитаем небольшой интерфейс и конкретный возвращаемый тип; не добавляем интерфейс для каждой структуры или общий framework ради моков.
- Зависимости передаются явно через конструкторы. Избегаем изменяемых глобальных клиентов, скрытых побочных эффектов `init` и service locator. Допустимые исключения, например регистрация драйвера, должны быть понятны из кода.
- Generics, reflection, дополнительные слои и библиотеки требуют конкретной пользы. Вначале используем простое решение и возможности стандартной библиотеки; не проектируем абстракции для несуществующих задач.
- Комментарии объясняют контракт, причину и ограничение. Экспортируемый API документирован; generated code не редактируется вручную и проверяется воспроизводимой генерацией, когда она появится.

Основа стиля: [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) и [Google Go Style Guide](https://google.github.io/styleguide/go/guide). Проектные требования здесь являются нашим выбранным стандартом, а не утверждением, что все Go-команды используют одинаковые правила.

#### 4. Ошибки и context

- Значимая ошибка обрабатывается. Осознанно игнорируемый результат объясняется; не скрываем ошибки БД, публикации, сериализации или завершения ресурса. Ошибка очистки не должна незаметно заменить исходную причину сбоя.
- Добавляем полезный контекст через `fmt.Errorf("reserve inventory: %w", err)`, если вызывающему коду нужен исходный тип/сентинел. Проверяем цепочку через `errors.Is`/`errors.As`; не классифицируем ошибки по тексту сообщения.
- Доменные ошибки отделены от ошибок транспорта и временных отказов. gRPC-код, retryability и значение неизвестного результата заданы явно. Нехватка товара не становится `Internal`; timeout не становится доказанным бизнес-отказом.
- `panic` не используется для ожидаемых ошибок клиента, провайдера или БД. Решение завершить процесс находится на границе приложения; библиотечный код не вызывает `os.Exit`/`log.Fatal` и не обходит cleanup.
- `context.Context` передаётся первым аргументом операции и не бывает `nil`. Контекст запроса не сохраняется как состояние доменной сущности. Не заменяем его на `Background` глубоко внутри RPC/SQL, чтобы скрыть отмену.
- Deadline проходит в зависимости; каждый созданный cancel вызывается. Фоновый воркер использует управляемый контекст своего жизненного цикла и ограниченный бюджет попытки. После сохранения принятой работы уход клиента не удаляет job: это правило процесса, а не игнорирование context.
- Ошибка логируется у границы, где принимается решение о её обработке; не печатаем одну цепочку на каждом слое. Осмысленные логи попытки и окончательного исхода различимы.

Семантика отмены и времени жизни: [документация context](https://pkg.go.dev/context).

#### 5. Конкурентность и ресурсы

- У каждой goroutine есть владелец, условие завершения и способ дождаться её остановки. Конкурентная работа, очереди и количество попыток ограничены; overload не создаёт бесконечные goroutines.
- Владение состоянием и синхронизация явны. Не копируем используемые mutex/atomic-значения; указатели сами по себе не делают доступ безопасным. Канал закрывает сторона, владеющая отправкой.
- Блокирующие ожидания, retry/backoff и симулированные задержки учитывают отмену. Нет busy loop и сна в открытой SQL-транзакции. Синхронная функция предпочтительна, если фоновая работа не нужна по контракту.
- Соединения, SQL rows, файлы и таймеры имеют понятный срок жизни. Проверяем `rows.Err()` после итерации; транзакция заканчивается commit или rollback. Connection pool ограничен и учитывается суммарно по репликам.
- Shutdown прекращает приём работы, даёт ограниченное время текущим операциям и оставляет незавершённую принятую работу восстанавливаемой. Readiness и liveness имеют разные смыслы, когда зависимости и фоновые jobs появятся.
- `go test -race` обязателен для применимого Go-пакета изменений на поддерживаемом runner. Успех подтверждает только исполнившиеся пути; он не доказывает отсутствие всех гонок.

Ограничения детектора: [Go race detector](https://go.dev/doc/articles/race_detector). Правила SQL API: [database/sql](https://pkg.go.dev/database/sql); при выборе другого драйвера также проверяем его собственный контракт.

#### 6. Данные и распределённые эффекты

Здесь проверяется соответствие [принятому дизайну](#appendix-design-review) и [процессам G0](#appendix-g0-processes):

| Область | Обязательное свойство ревью |
| :--- | :--- |
| PostgreSQL | Сервис владеет своей БД; SQL параметризован; транзакции короткие; конкурентная мутация защищена условием/версией и ledger |
| Миграция | Понятны совместимость со старым кодом, блокировки, порядок выкладки и восстановление. Обратная миграция не объявляется безопасной без проверки потери данных |
| Деньги и количество | Денежные суммы хранятся в целых минимальных единицах с валютой; проверяются допустимый диапазон, переполнение и согласованная сумма. Количество и остатки не уходят в недопустимое состояние |
| gRPC | Валидируются входы, сроки и размеры; сохранённый отказ и повторный запрос имеют стабильный смысл; timeout означает неизвестный исход |
| Kafka | Схема/версия, event ID и ключ партиции понятны; дубликаты и поздние события ожидаемы; поведение poison event и replay определено |
| Outbox/inbox/jobs | Состояние и outbox записаны атомарно; sent marker следует после Kafka ACK; inbox и продолжение сохранены до offset commit; crash windows проверены |
| Идемпотентность | Один ключ с тем же смыслом даёт прежний результат, включая подтверждённый отказ; иной payload под прежним ключом не создаёт новый эффект |
| Отмена | Version fencing и cancellation barrier блокируют поздние эффекты; `COMPLETED` и `CANCELLED` подтверждены данными владельцев |
| Redis | Кэш не разрешает продажу; miss/stale/unavailable имеют определённое поведение; fallback не создаёт неограниченное давление на БД |

Внешний RPC, Kafka publish и задержка симулятора не выполняются внутри бизнес-транзакции БД. Мок, который всегда отвечает успешно, не подтверждает атомарность PostgreSQL или семантику Kafka. Если изменение затрагивает эти свойства, нужны соответствующие интеграционные сценарии. Низкий Kafka lag не заменяет проверку backlog долговечных задач и бизнес-инвариантов.

#### 7. Тесты и доказательства

- Тест проверяет наблюдаемый контракт и существенный риск. Для нового поведения нужны соответствующие случаи; для исправления дефекта — регрессия, когда её можно воспроизвести. Не добавляем тесты, повторяющие устройство функции, для простой правки текста или форматирования.
- Используем стандартный `testing`, subtests и table-driven cases, когда случаи действительно параллельны по структуре. Тесты изолированы: `t.Parallel` не делит изменяемый fixture/БД; cleanup выполняется через `t.Cleanup` или явного владельца.
- Синхронизация тестов задаётся событиями, каналами и ограниченным ожиданием. Случайный `Sleep` не служит доказательством завершения. Время, seed и симулированные исходы контролируются там, где это влияет на повторяемость.
- SQL-конкуренция, уникальные ограничения, outbox и Kafka ACK проверяются на реальных зависимостях в интеграционных тестах. Unit test и интеграционный прогон имеют отдельные команды и результаты; отсутствие среды обозначается `NOT_RUN`.
- Применимые сценарии: повторный ключ; другой payload под тем же ключом; конфликт за остаток; отмена/deadline; timeout после commit; duplicate/late event; рестарт в изменённом crash window; ошибка или отказ зависимости. Не требуем весь каталог для каждого PR — выбираем затронутые свойства.
- Coverage помогает найти непроверенные ветки, но фиксированный процент не заменяет проверку инвариантов. Fuzzing уместен для парсеров, валидации и кодеков; benchmark и profile — при заявлении об улучшении производительности.
- Результаты нагрузки сравниваются по принятой [методике roadmap](#appendix-roadmap) и [шаблону отчёта](#appendix-report-template): условия, ресурсы, drain, verifier и неизвестные исходы сохраняются. Числа не придумываются и не переносятся из чужих измерений.

Дополнительный инструмент проверки входов: [Go fuzzing](https://go.dev/doc/tutorial/fuzz).

#### 8. Наблюдаемость и безопасность

Конфигурация валидируется при запуске. Секреты, connection strings с паролями и чувствительные payload не попадают в репозиторий или логи. Входные размеры, очереди и сроки ограничены; SQL не собирается конкатенацией пользовательских значений. Fault hooks включаются явно в профиле эксперимента, отключены в обычном профиле и отражены в manifest.

Логи структурированы и позволяют связать заказ, событие, шаг и попытку. `order_id`, `trace_id`, `event_id` и `run_id` остаются полями логов, а не labels с неограниченной кардинальностью. Метрики различают технический ответ и бизнес-исход, Kafka lag и durable job backlog. При появлении алерта ревью проверяет условие, ожидаемое действие и ссылку на диагностику.

Новая зависимость обоснована задачей; `go.mod` и `go.sum` согласованы. Версии Go, инструментов, генераторов и CI actions закрепляются воспроизводимо. Для найденной вызываемой уязвимости показываем исправление либо решение владельца проекта о конкретном риске; техническая ошибка сканера не считается чистым результатом.

Основа проверки зависимостей: [Go security best practices](https://go.dev/doc/security/best-practices) и [govulncheck](https://go.dev/doc/tutorial/govulncheck).

#### 9. Автоматические проверки и момент подключения

Go-модуль, закреплённые инструменты и нативный Docker CI AMD64/ARM64 задают автоматическую базу ревью. Для main обязательны Documentation и оба Go toolchain job. Фактические результаты конкретного HEAD находятся в PR и датированных evidence.

| Проверка | Требование после появления Go-модуля |
| :--- | :--- |
| Формат и импорты | `gofmt` и `goimports`: CI проверяет отсутствие diff, а не исправляет файлы молча |
| Сборка | `go build ./...` на закреплённой версии Go |
| Статический анализ | `golangci-lint run ./...` с явным набором `govet`, `staticcheck`, `errcheck`, `ineffassign`, `unused` |
| Unit tests | `go test ./...`; чистый обязательный CI прогон использует `-count=1` |
| Гонки | `go test -race -count=1 ./...` на поддерживаемой среде |
| Зависимости | `go mod tidy` без diff в `go.mod`/`go.sum`; `go mod verify` |
| Уязвимости | `govulncheck ./...`; версия инструмента закреплена, результат и ошибки запуска различимы |
| Документация | `node tools/check_docs.mjs` и `tools/update_context.ps1 -CheckOnly` |

`govet` внутри golangci-lint покрывает выбранный запуск vet; отдельный `go vet ./...` можно использовать для диагностики без требования дублировать одинаковый анализ в CI. Не включаем все существующие линтеры. `nolint` имеет узкую область, название линтера и причину; массовое подавление для зелёного CI недопустимо. Набор и версии инструментов фиксируются в первом Go-пакете и проверяются его реальным запуском. Список анализаторов сверяется с [документацией golangci-lint](https://golangci-lint.run/docs/linters/).

Интеграционные тесты и проверка генерации подключаются вместе с соответствующими SQL/протоколами/Kafka; команды и профили будут добавлены при реализации. Результат `SKIPPED`, недоступная зависимость или отсутствие нужного теста не выдаются за успешную проверку. Предлагаемые настройки CI не называются включёнными GitHub required checks, пока настройка репозитория не выполнена и не проверена.

#### 10. Когда изменение готово к принятию

- Цель и поведение ясны, изменение входит в согласованный объём и сохраняет принятую архитектуру.
- Блокирующие дефекты исправлены; существенные отклонения и риски имеют явное решение владельца проекта.
- Применимые проверки прошли на предъявленном изменении; непроведённые проверки и их причины видны ревьюеру.
- Контракты, миграции и операционные последствия понятны; для затронутого риска есть достаточное доказательство.
- Документация и переносимый контекст актуальны. Человек рассмотрел этот результат и явно согласовал слияние.

Стандарт помогает принимать улучшения без стремления к недостижимому совершенству. Основание подхода к обязательным замечаниям и косметике — [Google: The Standard of Code Review](https://google.github.io/eng-practices/review/reviewer/standard.html). Для короткой проверки каждого PR используется [шаблон](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/.github/pull_request_template.md); он ссылается на этот документ и не дублирует всю спецификацию.

<a id="appendix-lab-boundary"></a>

## Приложение: Общая лаборатория, адаптеры и кандидат второго бизнеса

Полный исходный документ: `docs/architecture/LAB_BOUNDARY.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Лаборатория и подключаемые бизнес-системы

Статус на 13 октября 2026 года: **направление подтверждено владельцем проекта; ниже — детализация для ревью**. Реализации интерфейсов и второго подключения пока нет. Документ сохраняет обсуждение, для обсуждения и проверки пользователями и тестировщиками.

#### Что уже согласовано

Лаборатория должна подключаться к разным бизнес-системам. Эталонный мини-магазин использует ту же границу через встроенный, заранее настроенный адаптер. Владелец проекта подтвердил этот подход в голосовом обсуждении и поручил его зафиксировать.

Принятые правила [G0](#appendix-design-review) сохраняются: у магазина ровно три бизнес-сервиса — Orders, Inventory и Payments. Другой бизнес будет отдельной тестируемой системой. Он не добавляет четвёртый сервис к магазину.

[Общая среда Docker](#appendix-environment) для Windows и macOS согласована отдельно. Обработка изображений пока кандидат для второго подключения. Точные Go-интерфейсы, формат профиля и расширение критериев G4 требуют ревью; подробность этого документа не означает их принятия или реализации.

#### Три части системы

```mermaid
flowchart TB
    Profile["Профиль эксперимента<br/>нагрузка · данные · ресурсы · критерии"] --> Core
    Core["Общая лаборатория<br/>расписание · цикл запуска · сбор данных · сравнение"] --> ShopAdapter
    Core --> FileAdapter
    ShopAdapter["Адаптер магазина<br/>операции · данные · проверки"] --> Shop
    FileAdapter["Адаптер обработки изображений<br/>кандидат для второго подключения"] --> Files
    Shop["Мини-магазин<br/>Orders · Inventory · Payments"] --> Evidence
    Files["Другой бизнес<br/>свой API и свои правила"] --> Evidence
    Evidence["Наблюдения<br/>исходы операций · метрики · логи"] --> Core
    Core --> Report["Общий отчёт<br/>корректность · скорость · ресурсы · ограничения"]
```

Адаптер здесь — модуль клиентского кода внутри инструментов лаборатории, а не обязательный сетевой сервис или постоянно работающий контейнер. Магазин не импортирует лабораторию и не знает, какой генератор вызывает его API.

| Часть | Ответственность | Пример для магазина | Пример для обработки изображений |
| :--- | :--- | :--- | :--- |
| Общая лаборатория | Расписание нагрузки, лимиты генератора, seed, окна измерения и drain | Подать заданный поток бизнес-операций | Подать заданный поток бизнес-операций |
| Общая лаборатория | Запись фактических запусков, воздействий, исходов и полноты телеметрии | Сохранить наблюдения и manifest | Сохранить наблюдения и manifest |
| Общая лаборатория | Сравнение вариантов и оформление доказательств | Отчёт A/B по общему формату | Отчёт A/B по тому же формату |
| Адаптер | Клиенты API и последовательность бизнес-действий | Создать заказ, получить итог, при необходимости отменить | Отправить изображение, дождаться обработки, скачать результат |
| Адаптер | Подготовка данных и сверка неизвестных исходов | SKU, остатки, request keys и авторитетный итог заказа | Синтетические файлы, ключи задания и авторитетный итог обработки |
| Адаптер | Семантические проверки и бизнес-метрики | Складской и платёжный ledger, завершённые заказы | Доступность результата, формат, размеры, завершённые задания |
| Тестируемая система | Собственная бизнес-логика и долговечные данные | Сервисы и базы магазина | Сервисы и хранилища другого проекта |

В общем механизме не должно быть условий по названиям Orders, SKU, оплате или таблицам магазина. Сборщик знает источник телеметрии и окно времени; значение бизнес-показателя описывает профиль адаптера. Неподдерживаемые метрики не заполняются придуманными нулями.

#### Один цикл эксперимента

```mermaid
flowchart LR
    Prepare["Подготовка<br/>профиль · данные · доступные возможности"] --> Warmup["Прогрев"]
    Warmup --> Run["Измерение<br/>нагрузка + запланированные воздействия"]
    Run --> Drain["Drain<br/>остановить новые операции · дождаться исходов"]
    Drain --> Verify["Проверка<br/>сверить все операции и бизнес-правила"]
    Verify --> Export["Экспорт<br/>manifest · исходные данные · отчёт"]
    Export --> Compare["Сравнение A/B<br/>одни условия · одно изменение"]
```

Runner управляет этапами. Адаптер выполняет предметные действия на этих этапах. Драйвер среды выполняет допустимые воздействия на указанные компоненты: например, остановку контейнера или изменение числа реплик. Сопоставление логического компонента с Compose-сервисом находится в профиле среды; в ядре нет жёсткого списка сервисов магазина.

Перед стартом проверяются необходимые возможности. Сценарий с перезапуском воркера не запускается, если профиль даёт только публичный HTTP API. Дополнительные воздействия на сеть и ресурсы проектируются отдельно, когда появится соответствующий эксперимент.

Прогрев и измеряемая когорта различаются. После остановки новых операций система продолжает обработку до установленного deadline. Оставшиеся операции отражаются как незавершённые или с неизвестным исходом. Cleanup выполняется после сохранения доказательств и затрагивает только данные выделенного тестового запуска.

#### Предлагаемый контракт адаптера

Это описание обязанностей, а не готовые имена Go-методов. Конкретные типы выбираются при реализации первого вертикального среза.

| Обязанность | Что получает лаборатория |
| :--- | :--- |
| Описать возможности | Доступные операции, виды наблюдений, необходимые права, проверки и воздействия |
| Подготовить данные | Изолированный набор данных по seed и параметры, необходимые для операций |
| Выполнить операцию | Ключ логической операции, наблюдения API, подтверждённый отказ или идентификатор принятой работы |
| Наблюдать и сверить итог | Завершённый бизнес-исход либо явное отсутствие достаточных сведений |
| Проверить бизнес-правила | Результаты проверок, нарушения и доказательства по измеряемой когорте |
| Очистить тестовые данные | Результат очистки собственного namespace без уничтожения чужих данных |

Одна бизнес-операция может включать несколько запросов API, повторов и опросов статуса. Их счётчики ведутся отдельно. Быстрый ответ о принятии задачи не считается завершённым бизнесом. Timeout не превращается автоматически в отказ: адаптер сверяет авторитетный результат, если API позволяет это сделать.

Повтор допускается по известному контракту идемпотентности конкретной системы. Лаборатория не добавляет автоматически повторные изменения данных к неизвестному API. Плановая нагрузка, фактическая отправка, пропущенные старты и насыщение генератора сохраняются отдельно; ожидание ответов не должно незаметно менять заданный arrival rate.

Для проверки предлагаются статусы `PASS`, `FAIL`, `NOT_RUN`, `INCONCLUSIVE`. Отсутствующий обязательный способ сверки делает соответствующий вывод неполным. Успех проверки формата ответа не является доказательством корректности бизнеса.

#### Воспроизводимый профиль

Профиль тестируемой системы описывает адаптер и его версию, endpoints, доступные возможности, параметры данных, бизнес-проверки и источники телеметрии. Секреты передаются отдельно; в отчёте сохраняются только ссылки на настройки и обезличенные сведения.

Профиль эксперимента задаёт смесь операций, arrival schedule, seed, warmup/measure/drain, timeline воздействий, ресурсы, критерии достаточности данных и ожидаемый полезный эффект. Это проектируемый собственный формат, не существующий стандарт стороннего инструмента.

Manifest сохраняет фактические версии, image digests, конфигурацию, оборудование, ограничения Docker, workload и выполненные воздействия. Сравнения проводятся при одинаковых исходных условиях. Общий формат отчёта сохраняется, но бизнес-метрики разных доменов не становятся взаимозаменяемыми: заказы в секунду не равны обработанным изображениям в секунду.

Мастер автоматической настройки остаётся перспективой: импорт API, документации и тестов → предложение операций, данных и проверок → подтверждение смысловых правил владельцем системы → пробный запуск → сохранение проверенного профиля. Во время измерения выполняется фиксированный профиль; генерация новых гипотез ИИ не вмешивается в измеряемый цикл.

#### Кандидат второго бизнеса: обработка изображений

Предлагаемый процесс: клиент отправляет синтетическое изображение через HTTP, система принимает фоновое задание, воркер преобразует изображение, клиент получает статус и скачивает результат. Такой процесс отличается от заказа, резерва и оплаты; он проверяет смену протокола, данных и критериев корректности.

Источник второго бизнеса ещё не выбран: это может быть существующий демонстрационный проект или отдельный небольшой пример. Его API и правила нужно изучить перед написанием адаптера. Утверждённый стек магазина не навязывается чужой системе.

| Проверка | Требуемое основание |
| :--- | :--- |
| Каждое принятое задание учтено | Доступный авторитетный статус или экспорт; unresolved перечисляются явно |
| Успешный результат можно скачать и декодировать | Контракт статуса завершения и доступа к результату |
| Размеры и формат соответствуют запросу | Зафиксированные правила преобразования; эталонные синтетические входы |
| Повтор не создаёт лишний эффект | Только если целевая система обещает идемпотентность по ключу |
| После сбоя принятая работа восстанавливается | Только если долговечность и восстановление входят в требования системы |

Первый эксперимент-кандидат: одинаковый набор изображений и профиль поступления; baseline с малым числом воркеров → наблюдаемый backlog → изменение числа воркеров → повтор и проверка результатов. Если вместе с воркерами растёт ресурсный бюджет, это отдельный фактор, а не бесплатное улучшение. Перезапуск воркера и медленное хранилище — следующие сценарии, когда профиль предоставит нужные возможности.

#### Как проверить переносимость

Предлагаемое доказательство: подключить выбранный второй бизнес новым адаптером и профилем, провести полный цикл с нагрузкой, хотя бы одним доступным воздействием, сверкой исходов и сравнением вариантов. При этом общий scheduler, runner, экспорт и структура отчёта сохраняются. Допустимые исправления общего механизма объясняются общей потребностью, а не скрытым условием для конкретного бизнеса.

В review packet входят diff, оба профиля, manifest, исходные наблюдения, результаты предметных проверок и ограничения сравнения. Проверка на mock полезна для разработки runner, но полный запуск второго бизнеса остаётся отдельным свидетельством переносимости. Это предлагаемый будущий критерий, ещё не дополнение к принятой контрольной точке G4.

#### Предлагаемое размещение кода

Основная [структура roadmap](#appendix-roadmap) сохраняется. К ней по мере появления кода предлагаются:

```text
internal/lab/                  # Контракты, scheduler, runner, результаты
internal/adapters/shop/        # Клиенты, данные и проверки магазина
internal/adapters/images/      # Будущий адаптер выбранного второго бизнеса
experiments/targets/           # Профили систем и привязка компонентов к среде
experiments/scenarios/         # Воспроизводимые условия эксперимента
```

Адаптер зависит от контрактов лаборатории и внешних API системы. Ядро зависит от общих контрактов; оно не импортирует бизнес-реализацию или адаптер магазина. Точки запуска `cmd/loadgen` и `cmd/labctl` выбирают нужный адаптер. Пустые каталоги и универсальный framework заранее не создаются.

#### Следующие шаги

1. В P01 подготовить контейнерный каркас трёх сервисов; бизнес-правила и контрольные точки G0 не меняются.
2. При появлении loadgen/verifier реализовать первый срез через адаптер магазина: подготовка → операция → наблюдение → сверка.
3. Уточнить контракт по работающему срезу и сохранить общий цикл без зависимости от магазина.
4. Выбрать второй бизнес и его требования; подготовить конкретный review packet подключения.
5. Доказать переносимость полным экспериментом и затем оценить необходимость мастера настройки.

Открытые вопросы: выбор второго проекта, минимальный набор capabilities, точные типы исходов и доказательств, формат профилей и место проверки второго бизнеса в приёмке. Они не блокируют каркас трёх Go-сервисов.

<a id="appendix-environment"></a>

## Приложение: Переносимая Docker-среда и распределение контейнеров

Полный исходный документ: `docs/development/ENVIRONMENT.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Переносимая среда разработки и экспериментов

Этот документ описывает конфигурацию и команды. Датированные проверки сохранены в evidence; актуальное выполнение и ревью — в [GitHub](https://github.com/Includeoyasi/faultline-lab/issues).

#### Что согласовано

На Windows и macOS используется Docker Desktop с Linux-контейнерами. Docker Compose объединяет отдельные контейнеры сервисов и инфраструктуры. Сборка Go, тесты и инструменты проверок выполняются в контейнерной среде; установленный на хосте Go не обязателен для основного пути разработки.

Для x64 и Apple Silicon планируются нативные варианты образов `linux/amd64` и `linux/arm64`. Архитектура личного Mac пока не уточнялась. Эмуляция другой архитектуры не используется как основание для сопоставимых замеров производительности. Docker описывает варианты образов и ограничения эмуляции в [руководстве по multi-platform builds](https://docs.docker.com/build/building/multi-platform/).

#### Распределение контейнеров

```mermaid
flowchart TB
    Host["Windows / macOS<br/>Docker Desktop + Compose"] --> Shop
    Host --> Infra
    Host --> Obs
    Host --> Tools
    Shop["3 постоянных контейнера<br/>Orders · Inventory · Payments"]
    Infra["3 постоянных контейнера<br/>PostgreSQL · Kafka · Redis"]
    Obs["4 контейнера по необходимости<br/>Prometheus · Grafana · Loki · Alloy"]
    Tools["Временные контейнеры<br/>сборка / проверки / эксперимент"]
    Tools --> Adapter["Адаптер — модуль в инструменте<br/>отдельный контейнер ему не нужен"]
    Adapter --> Shop
```

Числа относятся к планируемому составу с одной репликой каждого сервиса. Экспортеры, дополнительные реплики и диагностические инструменты могут потребовать дополнительных контейнеров и будут отражены в manifest.

| Компонент | Размещение | Причина |
| :--- | :--- | :--- |
| Orders, Inventory, Payments | По одному контейнеру на сервис и на его реплику | Независимые рестарты, health, логи, ограничения ресурсов и масштабирование |
| Воркеры, outbox, управление заказом | Внутри своего бизнес-сервиса | Это части принятого сервиса G0; отдельный контейнер на каждый внутренний механизм не требуется |
| PostgreSQL | Один контейнер, три базы и отдельные роли | Уже принятый компромисс малого стенда; владение данными сохраняется |
| Kafka | Один брокер в отдельном контейнере на старте | Независимый backlog, рестарт и ресурсный предел; отказоустойчивость кластера этим составом не моделируется |
| Redis | Отдельный контейнер | Можно независимо проверить отказ кеша и нагрузку fallback на SQL |
| Prometheus, Grafana, Loki, Alloy | По одному контейнеру на компонент | Разные задачи, конфигурация, хранение и жизненный цикл; запускаются нужным набором |
| Toolchain | Общий образ с Go и закреплёнными инструментами, временные запуски | Не держать компилятор и линтер работающими всё время |
| Loadgen, verifier, экспорт и отчёт | Общий образ инструментов лаборатории, запуск по задаче | Не требуется отдельный постоянный контейнер для каждой CLI-команды |
| Адаптер бизнеса | Модуль внутри инструментов лаборатории | Обозначение на схеме не означает дополнительный сервис |

Docker рекомендует разделять ответственности по сервисам; несколько процессов внутри контейнера допустимы. Поэтому правило здесь — независимый жизненный цикл и измеряемая ответственность, а не формальное «один процесс». [Docker: несколько процессов в контейнере](https://docs.docker.com/engine/containers/multi-service_container/).

Объединить три бизнес-сервиса в один контейнер технически возможно, но для этого проекта невыгодно: обычные остановка и лимит ресурсов контейнера затронут все три, а масштабирование скопирует весь набор. Раздельные условия придётся воссоздавать дополнительным управлением процессами. Это затруднит эксперименты по поиску узкого места.

Объединение PostgreSQL, Kafka и Redis также не устраняет потребление памяти их процессами и кешами, но усложняет независимые сбои, health и обслуживание данных. Фактический расход ресурсов нужно измерить на готовом стенде; число контейнеров само по себе не определяет требуемую RAM.

Рекомендуемый состав: **6 постоянных контейнеров основного стенда**, **10 с полной базовой наблюдаемостью**, плюс временные инструменты во время выполнения задач. На этапе P01 достаточно каркаса трёх сервисов и toolchain; инфраструктура добавляется в следующих пакетах. Полная наблюдаемость обязательна для соответствующих экспериментов, а не для каждой сборки кода.

#### Профили и запуск по этапам

В roadmap уже предусмотрены `core` и `obs`. Предлагается сохранить их смысл: основной стенд и дополнительная наблюдаемость; конкретное устройство Compose и команды будут предъявлены вместе с реализацией. Toolchain и lab-инструменты запускаются явно по задаче. Не добавляются обязательные admin UI, reverse proxy или постоянный сервис адаптера без потребности эксперимента.

[Профили Compose](https://docs.docker.com/compose/how-tos/profiles/) позволяют включать нужные наборы компонентов. Компоненты без профиля запускаются по умолчанию; при реализации нужно проверить также явный запуск отдельного сервиса и его зависимости.

Дорогие линтеры, сборка и формирование отчёта не выполняются одновременно с измеряемой нагрузкой без специально оговорённого эксперимента. Сам генератор и наблюдаемость потребляют ресурсы: их лимиты, насыщение и влияние отражаются в отчёте. Разделение контейнеров на одном компьютере не создаёт физически независимые машины.

#### Требования к будущей реализации

- Multi-stage сборка: toolchain готовит бинарник, runtime-образ запускает сервис. Go, библиотеки, инструменты и образы закрепляются при реализации; выбранные версии записываются в manifest.
- Относительные пути и переносимые Docker/Compose-команды. Обязательный путь запуска не зависит от Bash, Homebrew или установленного на хосте Go. Внешняя проверка переносимого контекста пока остаётся PowerShell-инструментом.
- Для данных PostgreSQL и других компонентов используются явно описанные volumes; очистка тестового стенда — отдельное действие. В Compose нет абсолютных путей исходного компьютера.
- Каждый сервис корректно обрабатывает остановку, имеет health и пишет структурированные логи. Готовность зависимостей подтверждается проверкой; порядок старта сам по себе не доказывает готовности.
- Контейнеры обращаются друг к другу через имена сервисов Compose. На хост публикуются только нужные точки доступа. [Docker: сеть Compose](https://docs.docker.com/compose/how-tos/networking/).
- Профиль определяет, где выполняется генератор: внутри Compose или вне тестируемой среды. При подключении другого бизнеса это не должно требовать изменения его исходного кода.

В [P01](#appendix-design-review) к трём entrypoints, конфигурации, gRPC health, shutdown, JSON-логам и Go CI добавляются образы сервисов, контейнерный toolchain и минимальный Compose-запуск. PostgreSQL/Kafka/Redis, базы, роли, миграции и seed остаются P02. Сначала готовится рабочий набор файлов, затем выполняются доступные проверки и предъявляются ограничения.

#### Что известно о текущем компьютере

Установка проверена 15 октября: Windows x64, PowerShell 7.6.5; Go1.27.2, Docker Desktop4.95/Engine29.8.2/Compose5.6 и WSL2.7.13 установлены. Windows container smoke P01 PASS; после него Docker/WSL выключены. 16 октября выполнен [полный PostgreSQL smoke](#p02-docs-development-evidence-p02-postgres-md) в том же WSL-бюджете.

Характеристики Windows-хоста и подтверждённый WSL-бюджет 2 CPU / memory=2GB / swap=1GB сохранены в [ресурсных профилях](#appendix-resource-profiles). PostgreSQL foundation прошёл функциональные проверки в этом бюджете; полный состав Kafka/Redis и нагрузочные эксперименты ещё не измерены; macOS-окружение также не проверено. Возможный арендуемый сервер рассматривается позже при необходимости; он не выбран. Docker делает сборку и конфигурацию переносимыми, но одинаковой производительности разных компьютеров не обещает.

Для A/B по умолчанию используется одна машина с одинаковым бюджетом и нативной архитектурой. Сравнение между Windows и macOS возможно как отдельный эксперимент с явными ограничениями. В manifest сохраняются архитектура CPU, аппаратные ресурсы, OS, Docker, volumes, image digests и фактические лимиты.

Локальный runtime-запуск проверен на Windows через Docker Desktop/WSL2. [Доказательства smoke](#windows-smoke-proof) содержат версии, лимиты и исходы; личный Mac пока не проверен.

##### Личный компьютер после проверки

По поручению владельца среда запускается по необходимости. После smoke контейнеры/сеть удалены, Docker Desktop и WSL VM остановлены; автозапуск Docker выключен. Go, Docker, WSL, актуальные общие образы и полезные кеши сохраняются. По уточнению владельца 16 октября образы одноразовых smoke-прогонов удаляются, их отсутствие проверяется; общий prune не применяется. WSL ограничен 2 CPU, memory=2GB, swap=1GB. Нового измерения бизнес-производительности это не заменяет.

<a id="appendix-resource-profiles"></a>

## Приложение: Ресурсный снимок Windows и будущие профили macOS/сервера

Полный исходный документ: `docs/development/RESOURCE_PROFILES.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Ресурсные профили сред

Статус на 15 октября 2026 года: снимок Windows-хоста от 13 октября сохранён ниже как исторический. Go/WSL/Docker теперь установлены, Windows smoke PASS, текущий бюджет WSL и завершение среды зафиксированы отдельно. macOS/сервер не измерены.

#### Три уровня данных

1. **Хост:** CPU, доступная ОС память и свободное место. Это исходные условия; свободные ресурсы меняются с фоновой нагрузкой.
2. **Бюджет:** ресурсы, выделенные Docker и отдельным компонентам. Бюджет задаётся и фиксируется, а не выводится из общего объёма памяти компьютера.
3. **Потребление:** фактические CPU, память, диск, очереди и насыщение во время прогона. Измеряется на работающем стенде вместе с бизнес-результатами.

Профили не обещают одинаковую производительность разных машин. По умолчанию сравнение A/B выполняется на одной машине с одинаковыми лимитами. Изменение ресурсов или перенос на другой хост учитывается как отдельный фактор.

До установки 15 октября Docker/WSL отсутствовали. После установки выполнен [Windows smoke](#windows-smoke-proof): Go1.27.2, WSL2.7.13, Desktop4.95, Engine29.8.2, Compose5.6. Бюджет: WSL memory=2GB, processors=2, swap=1GB; Engine NCPU=2, MemTotal=1996603392 bytes. Перед запуском свободно около 2.06 GiB RAM; это отдельный изменчивый снимок, не постоянная ёмкость. После проверки Docker/WSL VM выключены, автозапуск Docker отключён.

#### Windows: исторический снимок 13 октября

| Параметр | Значение |
| :--- | :--- |
| ОС | Windows 11 Pro, build 26200, x64 |
| CPU | Intel Core i5-12400F, 6 физических ядер / 12 логических процессоров |
| Память, видимая ОС | 15.72 GiB |
| Свободная память в момент чтения | 3.38 GiB |
| Ёмкость тома рабочего каталога | 930.51 GiB |
| Свободное место в момент чтения | 134.32 GiB |
| Go / Docker в текущем PATH | Не обнаружены |
| Команда WSL | Обнаружена; работоспособность и конфигурация не проверены |
| Docker daemon и лимиты | Не проверены |
| Потребление стенда | На дату снимка не измерено; runtime ещё не был создан |

Основание: чтение Windows CIM классов процессора/ОС и сведений о томе через PowerShell. Первое чтение CIM было ограничено песочницей; повтор с разрешённым расширенным доступом выполнен успешно. Имена компьютера, серийные номера, идентификаторы оборудования и персональные пути не сохранялись.

Свободная память — снимок, а не постоянный лимит. Нельзя выделять стенду весь объём хоста, предполагая, что он свободен. После появления Docker нужно отдельно зафиксировать его backend, доступный бюджет, лимиты компонентов и потребление при простое/нагрузке. Для WSL2 ограничения виртуальной среды задаются отдельно от настроек других backend Docker Desktop. [Docker: настройки ресурсов](https://docs.docker.com/desktop/settings-and-maintenance/settings/).

#### macOS: основной второй компьютер

Профиль пока `NOT_MEASURED`. Нужны архитектура CPU, число доступных ядер, память, свободный диск, версия macOS, Docker/Compose, backend и лимиты Docker. Нативный вариант образа выбирается по фактической архитектуре ноутбука. Проверка на Windows не подтверждает запуск на macOS.

#### Сервер: возможное последующее расширение

Статус `FUTURE_OPTION`: сервер не выбран и не арендован; характеристики, стоимость и deployment не определены. Владелец проекта рассматривает его при необходимости масштабирования. Покупка, аренда и подключение не поручались.

Предполагается использовать те же образы и профили с конфигурацией среды. Когда появится конкретный сервер, нужно зафиксировать архитектуру, vCPU, RAM, диск, сеть, версии Docker/Compose и лимиты. Перенос на сервер не равен автоматически распределённому кластеру; необходимость нескольких узлов рассматривается отдельным предложением.

#### Следующее измерение

После подготовки каркаса и доступности Docker: проверить daemon/Compose → снять фактический бюджет → запустить минимальный профиль → измерить простой → добавить инфраструктуру → проверить ресурсы под контролируемой нагрузкой. Каждый прогон сохраняет свой снимок условий в manifest; текущая таблица не заменяет данные будущих запусков.

15 октября установка завершена; VirtualMachinePlatform Enabled, гипервизор и Linux backend работают. Перезагрузка не выполнялась. Пиковые CPU/RAM и бизнес-throughput не измерялись.

<a id="appendix-banner-options"></a>

## Приложение: Концепции баннера и выбранное оформление

Полный исходный документ: `docs/design/BANNER_OPTIONS.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Баннер Faultline Lab

Статус на 13 октября 2026 года: **владелец проекта выбрал вариант 2 — «Поток эксперимента»**. Основание: прямое сообщение «Мне нравится вариант 2». Выбранный SVG применяется в обоих README через общий файл [banner.svg](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/docs/assets/banner.svg).

#### Задание

Подготовить три альтернативных варианта с акцентом на технологичность, современность, доступность, мягкие пастельные тона, исследование и будущее. Сделаны векторные композиции с читаемым названием, тёмным текстом на светлом фоне и встроенными title/description. Внешних шрифтов, скриптов и удалённых графических зависимостей нет.

![Три концепции баннера](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/docs/assets/banner-concepts/comparison.png)

| Вариант | Образ | Палитра | Статус |
| :--- | :--- | :--- | :--- |
| 1. Карта системы | Сервисы, связи и инженерный стенд | Голубой, мятный, лавандовый | Сохранён как альтернатива |
| 2. Поток эксперимента | Гипотеза → эксперимент → доказательства | Слоновая кость, лаванда, персик, мята | Выбран владельцем проекта |
| 3. Исследовательский горизонт | Любопытство, исследование и будущее | Голубой, лавандовый, мягкий персик | Сохранён как альтернатива |

#### Выбранное решение

![Поток эксперимента](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/docs/assets/banner-concepts/02-experiment-flow.png)

Абстрактные линии проходят через область наблюдения. Это метафора исследования поведения системы и сравнения изменений; линии не являются графиком измеренных результатов. Композиция не зависит от домена магазина и подходит будущим подключаемым бизнесам.

Размер SVG — 1280 × 420; README задаёт ширину 100%. Исходник масштабируется, PNG служит предварительным просмотром. Основная подпись: **Explore how systems behave.** Нижняя последовательность: **Hypothesis → Experiment → Evidence**.

Исходники: [карта системы](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/docs/assets/banner-concepts/01-system-map.svg), [поток эксперимента](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/docs/assets/banner-concepts/02-experiment-flow.svg), [исследовательский горизонт](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/docs/assets/banner-concepts/03-research-horizon.svg). Из выбранного исходника копируется общий banner.svg; если его изменяют, исходник концепции и preview обновляются вместе.

Баннеры отрендерены локально и визуально просмотрены на общем листе. Проверены отсутствие обрезанного текста и соответствие композиций заданию. Отдельный preview выбранного варианта используется для проверки основного баннера. Выбранный вариант слит в main через [PR #4](https://github.com/Includeoyasi/faultline-lab/pull/4), merge commit 8c392b8. Проверка отрисовки нового README на GitHub пока не выполнялась.

<a id="appendix-p01"></a>

## Приложение: Первый PR с кодом: задачи, зависимости и критерии приёмки

Полный исходный документ: `docs/development/packets/P01.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### P01 — Go-каркас между G0 и G1

Историческая приёмка: 16 октября 2026 года владелец проекта слил [PR #8](https://github.com/Includeoyasi/faultline-lab/pull/8), принятый HEAD `56f3672`, merge `3ca290d`. Это свидетельство P01, а не текущая доска задач.

Основа: [первый пакет реализации](#appendix-design-review), [roadmap](#appendix-roadmap), [порядок ревью](#appendix-development-process), [стандарт Go](#appendix-code-review) и согласованный [Docker-подход](#appendix-environment).

#### Цель и границы PR

Результат всего P01: три Go-сервиса Orders, Inventory и Payments запускаются в отдельных контейнерах, валидируют конфигурацию, публикуют gRPC health, пишут JSON-логи и корректно завершаются. Сборка и проверки доступны через закреплённый контейнерный toolchain и CI.

Заключительный PR #8 объединил P01-05–08: bounded graceful shutdown, runtime-образы/Compose, итоговые Go/CI проверки и smoke. Все восемь задач приняты. [Устройство сервисов](#appendix-service-runtime), [контейнерная инструкция](#appendix-container-runtime) и [toolchain](#appendix-toolchain) сохраняют команды воспроизведения.

PostgreSQL/Kafka/Redis и миграции относятся к P02. Доменные protobuf, заказы, складские операции, платежи, адаптеры и runner здесь не реализуются. Использование стандартного gRPC health не означает готовность бизнес-API.

Терминология, план среды, выбранный баннер и этот backlog слиты в main через [документационный PR #4](https://github.com/Includeoyasi/faultline-lab/pull/4) по отдельному поручению владельца проекта. Он не является P01: первым PR с кодом остаётся каркас, описанный здесь. По последующему поручению новые PR публикуются от исполнителя AI-Your-Bro после проверки его авторизации.

#### Историческая приёмка P01, 16 октября 2026 года

Все статусы ниже — снимок после наблюдаемого merge PR #8 владельцем. P01 завершён; [Issue #3](https://github.com/Includeoyasi/faultline-lab/issues/3) закрыт. Это не приёмка G1.

| ID | Задача | Зависит от | Статус | Доказательство приёмки |
| :--- | :--- | :--- | :--- | :--- |
| P01-01 | Go-модуль и закреплённый toolchain | — | DONE | Сборка/запуск трёх бинарников, format/lint/tidy/verify/govulncheck в нативном CI AMD64/ARM64 PASS; PR #6 |
| P01-02 | Конфигурация с проверкой | 01 | DONE | Успешный старт с допустимой конфигурацией; отказ на некорректной |
| P01-03 | Структурированные логи | 01 | DONE | Разбор стартовых/финальных записей как JSON с ожидаемыми полями |
| P01-04 | Три entrypoints и gRPC health | 02, 03 | DONE | Проверки всех трёх endpoints; ошибка занятого адреса |
| P01-05 | Graceful shutdown с пределом ожидания | 04 | DONE | Завершение активной работы и принудительное прекращение по deadline |
| P01-06 | Образы и минимальный Compose | 01–05 | DONE | Сборка образов, health и остановка трёх контейнеров |
| P01-07 | Go-проверки и CI | 01–05 | DONE | Успешные обязательные проверки на предъявленном HEAD |
| P01-08 | Smoke-проверка, инструкция и review packet | 06, 07 | DONE | Воспроизводимые команды, результаты, ограничения и открытый PR |

#### Проверяемые критерии по задачам

##### P01-01: модуль и инструменты

- Есть Go-модуль `github.com/Includeoyasi/faultline-lab`, `go.sum` и три точки запуска `cmd/orders`, `cmd/inventory`, `cmd/payments`.
- Создано описание контейнерного toolchain и команды сборки/проверок. Сборка всех трёх бинарников проходит; нет ручной установки Go на хосте в основном документированном пути.
- Версии Go, библиотек, линтера, goimports, govulncheck и образов закреплены; команды сообщают фактически используемые версии. Выбор конкретных версий выполняется при реализации.
- `go mod tidy` не оставляет diff, `go mod verify` проходит. Общая основа запуска не импортирует бизнес-логику одного сервиса в другой и не создаёт пустой framework будущей лаборатории.

##### P01-02: конфигурация

- Слушающий адрес задаётся явно через обязательный `FL_GRPC_ADDR`, без IP/портов и fallback в Go-коде. По замечанию владельца 15 октября эта настройка заменяет прежний default. Уровень логирования и shutdown timeout имеют описанные defaults; настройки поступают из окружения. Есть несекретные образцы локального и контейнерного окружения.
- Некорректный адрес/порт, неизвестный уровень логирования и неположительный timeout отклоняются до объявления готовности; процесс выходит с ненулевым кодом и понятной диагностикой.
- Проверки охватывают отсутствие/пустое значение обязательного адреса, defaults остальных полей, допустимые настройки и каждый перечисленный ошибочный вход. Чувствительные значения конфигурации не попадают в диагностические записи.

##### P01-03: JSON-логи

- При штатном запуске и остановке прикладные записи разбираются как JSON; есть время, уровень, сообщение и имя сервиса.
- Стартовая запись содержит доступную build/version информацию; неизвестное значение обозначается явно и не подменяется придуманным commit.
- Ошибки конфигурации, bind и завершения позволяют определить сервис и причину. Секреты не выводятся. Реальные записи приводятся в review packet.

##### P01-04: запуск и здоровье

- Каждый из трёх сервисов слушает свой endpoint и отвечает через стандартный `grpc.health.v1.Health/Check`.
- Общая проверка здоровья возвращает `SERVING` после успешной инициализации; не публикуется готовность до открытия listener.
- Занятый адрес приводит к понятной ошибке и ненулевому коду завершения, а не к зависшему процессу.
- Сервисы запускаются независимо; отключение одного не останавливает два других. На этом этапе health подтверждает работоспособность каркаса, без заявлений о готовности БД или бизнеса.

##### P01-05: завершение

- SIGTERM и SIGINT инициируют остановку, убирают готовность и прекращают приём новой работы. Health переход проверяется до закрытия server, когда endpoint ещё доступен.
- Уже выполняемый RPC получает возможность завершиться в пределах shutdown timeout. Проверка использует управляемый тестовый handler, а не выдуманную бизнес-операцию.
- Если RPC не заканчивается, по deadline выполняется принудительная остановка; нет бесконечного ожидания. Проверка ограничена timeout с явно указанным запасом тестового harness.
- Штатная остановка заканчивается без принудительного убийства контейнера; итог и ошибки завершения отражаются в логах. Синхронизация тестов не зависит от случайного Sleep.

##### P01-06: контейнеры и Compose

- Есть воспроизводимая multi-stage сборка образов трёх сервисов с использованием toolchain из задачи 01. Runtime запускает бинарник сервиса, а не держит компилятор работающим.
- Минимальный Compose поднимает только необходимый для каркаса состав; сервисы имеют отдельные healthchecks, настройки и жизненный цикл.
- Сборка образов для `linux/amd64` и `linux/arm64` проходит либо соответствующий обязательный критерий обозначен невыполненным. Один успешный build не объявляется проверкой второй архитектуры.
- На доступном Windows-хосте проверены старт, health и остановка. Сборка ARM64 не выдаётся за запуск на macOS; native smoke на личном Mac отдельно остаётся `NOT_RUN` до доступа к этой машине.
- Конфигурация не содержит абсолютных путей компьютера и обязательных Bash-скриптов на хосте. Точные переносимые команды фиксируются после появления файлов.

##### P01-07: проверки качества

- Проверки соответствуют принятому [стандарту](#appendix-code-review): gofmt/goimports, build, golangci-lint с согласованным набором, unit tests, race detector, tidy/verify и govulncheck.
- Тесты проверяют риски конфигурации, bind/health и остановки; не требуется бизнес-регрессия ещё не реализованных заказов или Kafka.
- CI запускает Go-проверки вместе с существующей проверкой документов/контекста. Результаты относятся к предъявленному HEAD; старый успешный документационный CI не считается доказательством нового кода.
- Версии инструментов и фактические команды воспроизводимы. Невыполненное не отмечается успешным.

##### P01-08: доказательства и ревью

- Инструкция позволяет собрать, проверить, запустить, запросить health и остановить каркас с описанными prerequisites. Проверенный путь не требует установленного на хосте Go.
- Smoke фиксирует исходы каждого сервиса; при ошибке завершает проверку ненулевым кодом и сохраняет диагностические сведения.
- Review packet содержит цель P01, diff, команды и фактические результаты, версии, ограничения Windows/macOS/архитектур, состояние задач и обновлённый переносимый контекст.
- Открыт PR в `main`, выполняется self-review; AI не одобряет и не сливает собственную работу. Принятие пакета фиксирует владелец проекта после ревью конкретного HEAD.

#### Приёмка пакета

**Готово к ревью:** перечисленный код и образы существуют; обязательные проверки выполнены и проходят; PR и инструкции доступны; существенные ограничения перечислены. При отсутствии Docker код можно подготовить и показать draft PR, но запуск/health/stop не получают статус PASS, а пакет не объявляется полностью проверенным.

**Принято:** владелец проекта принял предъявленный HEAD, блокирующие замечания устранены, применимые проверки успешны, PR слит по решению человека. Только тогда пакет получает DONE. Готовность PR к ревью и его приёмка — отдельные состояния.

Тест на личном Mac отражается отдельным результатом совместимости. Его отсутствие не скрывается и не превращает сборку multi-platform в утверждение «запуск на обеих ОС подтверждён». Если владелец требует native smoke на Mac до merge, это фиксируется как дополнительное условие именно этого пакета.

P01 не достигает всей контрольной точки G1: базы, миграции и проверяемые операции заказа/склада остаются следующими пакетами P02–P05.

#### Известные ограничения выполнения

- Go 1.27.2 и Docker Desktop 4.95/Compose 5.6 установлены; Windows Linux-container smoke PASS. macOS остаётся отдельной непроведённой проверкой.
- Характеристики Windows сохранены в [ресурсном профиле](#appendix-resource-profiles). Текущий WSL-бюджет 2 CPU / memory=2GB / swap=1GB подтверждён отдельно; потребление бизнеса под нагрузкой не измерено.
- Среда macOS и возможный будущий сервер недоступны для проверки; их результаты не придумываются.
- Команды toolchain приведены в [инструкции](#appendix-toolchain). Запуск, health и Docker smoke описаны в [контейнерном пакете](#appendix-container-runtime); Linux CI и локальная Windows проверка фиксируются отдельно.

#### Где следить за задачами

[Учёт задач](#appendix-task-tracking) описывает правила. Границы инфраструктурных работ сохранены в [P02](#appendix-p02), актуальное выполнение — в GitHub. Исходные Windows-измерения находятся в [отчёте](#windows-smoke-proof).

<a id="appendix-task-tracking"></a>

## Приложение: GitHub Issues, статусы и переносимые снимки

Полный исходный документ: `docs/development/TASK_TRACKING.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Учёт задач и статусов

**Единственный источник актуального выполнения — [GitHub Issues](https://github.com/Includeoyasi/faultline-lab/issues), [PR](https://github.com/Includeoyasi/faultline-lab/pulls) и [milestones](https://github.com/Includeoyasi/faultline-lab/milestones).** Владелец выбрал этот порядок в ревью F02 16 октября 2026 года. Ручные копии текущего прогресса в репозитории не ведутся.

#### Назначение материалов

| Место | Назначение |
| :--- | :--- |
| [Roadmap](#appendix-roadmap) | Этапы и контрольные точки, не текущая очередь |
| Спецификация пакета | Границы, зависимости и критерии, без дублирования стадии выполнения |
| GitHub Issue | Текущая стадия, исполнитель, критерии выполнения, блокеры и связь с PR |
| GitHub PR | Diff, результаты проверок, замечания и решение владельца по HEAD |
| ADR | Существенное решение, варианты, последствия и основание принятия |
| [Контекст](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/PROJECT_CONTEXT.md) | Полные спецификации, датированные решения и доказательства для автономного чтения |

#### Стадии в Issue

| Стадия | Основание перехода |
| :--- | :--- |
| BACKLOG | Работа описана, место в очереди ещё не выбрано |
| TODO | Владелец выбрал работу для предстоящей реализации |
| IN_PROGRESS | Реально выполняются реализация или проверки |
| IN_REVIEW | PR готов, доказательства представлены, ожидается ревью человека |
| DONE | Критерии выполнены, изменение принято владельцем и слито |

Это рабочие обозначения, а не новые встроенные состояния GitHub. Пока отдельная доска Projects не требуется, стадия указывается в теле Issue. BLOCKED фиксируется отдельно с причиной и необходимым действием. Draft PR или успешная сборка сами по себе не означают IN_REVIEW/DONE.

Чекбокс означает фактически выполненный критерий с проверкой; у агрегирующего Issue отложенные пункты остаются открытыми. Закрытие заменённой карточки как not_planned не является выполненной разработкой. Нельзя удалить проваленный критерий, чтобы объявить завершение.

#### Навигация и связи

- [P01 #3](https://github.com/Includeoyasi/faultline-lab/issues/3) и [спецификация P01](#appendix-p01) связывают первоначальный каркас с историей его приёмки.
- [Go smoke #9](https://github.com/Includeoyasi/faultline-lab/issues/9) связан с F01 из [аудита #20](https://github.com/Includeoyasi/faultline-lab/issues/20). Его границы находятся в [спецификации](#appendix-go-smoke-backlog), решение — в [ADR 0003](#appendix-adr-0003).
- [P02 #12](https://github.com/Includeoyasi/faultline-lab/issues/12) задаёт PostgreSQL/миграции/seed, [#16](https://github.com/Includeoyasi/faultline-lab/issues/16) — Kafka/Redis и интеграционные проверки. Зависимость: #12 → #16, по одному PR на задачу; [спецификация P02](#appendix-p02).
- [Аудит #20](https://github.com/Includeoyasi/faultline-lab/issues/20) сохраняет F01–F12 и контрольные точки возврата к будущим вопросам; он не закрывается одним частичным PR.

#### Обновление и доказательства

Исполнитель обновляет Issue при начале работы, появлении/снятии блокера и перед ревью. Значимое решение или новый результат проверки сохраняется в каноническом документе и датированном контексте; простая смена стадии не требует править несколько файлов.

Результат проверки содержит PASS/FAIL/NOT_RUN, команду, среду, дату и commit. Сборка под архитектуру не заменяет запуск на ней. Исторические протоколы и принятия не переписываются под текущее состояние. Контекст не объявляется актуальным трекером, даже если дата снимка совпадает с сегодняшней.

Для спорного объёма, несовместимого контракта или изменения утверждённых инвариантов нужно решение владельца. ИИ не одобряет и не сливает собственный PR; актуальные требования приёмки — в [workflow](#appendix-development-process).

<a id="appendix-repository-security"></a>

## Приложение: Защита main, отдельный исполнитель и состояние подключения

Полный исходный документ: `docs/development/REPOSITORY_SECURITY.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Защита main и аккаунт исполнителя

Снимок на 15 октября 2026 года; исходная настройка выполнена 13 октября. Владелец проекта прямо поручил перейти на отдельный аккаунт исполнителя и настроить защиту main, затем предоставил логин AI-Your-Bro. Защита применена через GitHub API и проверена отдельным GET. Авторизация исполнителя проверена через GET /user: AI-Your-Bro, GitHub ID 342030453. Его права на репозиторий — write, push=true, admin=false; приглашение принято. Git Credential Manager сохранил OAuth-авторизацию в защищённом хранилище Windows. Владелец прямо отказался включать 2FA; этот выбор записан, а статус 2FA не выдан за подтверждённый API.

#### Действующие настройки main

| Правило | Проверенное значение |
| :--- | :--- |
| Изменения через PR | Включено через обязательное независимое ревью |
| Число обязательных approval | 1 |
| Сброс approval после новых изменений | Включён |
| Approval после последнего push другим участником | Обязателен |
| Ревью CODEOWNERS | Включено; CODEOWNERS действует в main после PR #5 |
| Required checks | Links and portable context; Go toolchain (amd64); Go toolchain (arm64) |
| Источник required check | GitHub Actions, app ID 15368 |
| Актуальность ветки относительно main | Обязательна перед merge |
| Закрытие обсуждений | Обязательно |
| Применение правил к администратору | Включено |
| Force push / удаление | Запрещены |

Это проверка фактических настроек, а не попытка выполнить запрещённый push или удалить ветку. Подписанные коммиты и обязательная линейная история не вводились. Оба Go checks добавлены в required checks 15 октября после успешного нативного CI P01-01. PATCH изменил только список проверок; отдельный GET подтвердил сохранение reviewer/CODEOWNERS/admin enforcement и запретов force push/deletion.

[.github/CODEOWNERS](#appendix-codeowners) находится в main после слияния PR #5, подтверждённого 15 октября 2026 года (merge d616a22). Все файлы, включая workflow и сам CODEOWNERS, назначены @Includeoyasi. GitHub использует CODEOWNERS базовой ветки и требует ревью указанного владельца. [Правила CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners).

#### Разделение ролей и авторизации

- Аккаунт владельца проекта: `Includeoyasi`; управление репозиторием, человеческое ревью и решение о merge.
- Исполнитель: [AI-Your-Bro](https://github.com/AI-Your-Bro), отдельный machine account под контролем владельца, с подтверждённой авторизацией и правами write без admin. Он публикует ветки и PR; ИИ не использует его для одобрения собственных изменений.
- Административная авторизация владельца использована только для прямо порученной настройки защиты. Новые PR от аккаунта владельца после этого поручения не создаются.

Git author name/email не меняют автора PR: GitHub определяет его по аутентифицированному аккаунту API/UI. Перед первой публикацией надо проверить реальную identity API и scoped Git credentials; изменение подписи коммита само по себе не завершает переход.

GitHub разрешает machine account, зарегистрированный человеком, но запрещает автоматическую регистрацию. Владелец принимает условия и контролирует email; пароль, recovery codes и токены не помещаются в Git, документы или чат. API текущей OAuth-авторизации не вернул поле two_factor_authentication. На предложение включить 2FA владелец ответил «я не хочу ее включать»; ИИ этот выбор не отменяет и не сообщает о включённой 2FA. [Условия аккаунтов GitHub](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service#3-account-requirements).

#### Завершение подключения

1. Владелец предоставляет логин зарегистрированного вручную machine account. Выполнено: AI-Your-Bro, существование публичного профиля проверено.
2. Приглашение AI-Your-Bro в faultline-lab с правами write отправлено и принято (invitation ID 337527166). От имени исполнителя проверен GET collaborators/AI-Your-Bro/permission: write, push=true, admin=false.
3. OAuth через Git Credential Manager device flow выполнен; GET /user подтвердил AI-Your-Bro. Credentials находятся в Windows Credential Manager и не записываются в репозиторий. Администрация владельца и публикация исполнителя используют отдельные credentials; перед мутацией API проверяется фактический actor. Для работы с workflow нужны соответствующие права токена; администрация ради обычных PR не выдаётся.
4. Локальная Git identity этого репозитория: AI-Your-Bro, 342030453+AI-Your-Bro@users.noreply.github.com. Локальная credential.https://github.com.username выбирает AI-Your-Bro. Глобальная Git identity и авторизация других репозиториев не менялись. На macOS нужно отдельно подключить исполнителя и проверить identity; локальные credentials автоматически туда не переносятся.
5. [PR #4](https://github.com/Includeoyasi/faultline-lab/pull/4) уже слит: merge commit 8c392b8, принятый HEAD 83e0076. Переоформлять его не нужно; переход применяется к новым PR. Новая локальная ветка chore/executor-security обновлена до этого main.
6. [PR #5](https://github.com/Includeoyasi/faultline-lab/pull/5) слит владельцем: принятый HEAD 704a7ec, merge d616a22. ИИ не выполнял approve/merge. GitHub ранее подтвердил автора PR и коммиттера как AI-Your-Bro.
7. CODEOWNERS подтверждён в main, переносимый контекст обновлён. На 15 октября P01-01 в ревью; G1–G4 не приняты. Required checks перечислены выше и подтверждены серверным GET; одно наличие workflow не считается доказательством настройки.

#### Текущее ограничение

Защита main и CODEOWNERS действуют. Авторизация новой Windows-копии или macOS требует проверки исполнителя; credentials не переносятся через Git. Включение 2FA не заявляется: действует зафиксированное решение владельца. Текущий PR реализации и фактические Go-проверки находятся в CONTEXT_STATE.json.

<a id="appendix-codeowners"></a>

## Приложение: Действующий CODEOWNERS для человеческого ревью

Полный исходный документ: `.github/CODEOWNERS`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### The project owner reviews all changes, including CI and this policy.
* @Includeoyasi

<a id="appendix-toolchain"></a>

## Приложение: Контейнерный toolchain и команды

Полный исходный документ: `docs/development/TOOLCHAIN.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Go-модуль и контейнерный toolchain

Закреплённый toolchain используется CI и сборкой runtime-образов. Команды и границы проверок описаны ниже; фактические результаты конкретного изменения находятся в PR и датированных evidence.

#### Состав и версии

| Компонент | Источник закрепления | Версия |
| :--- | :--- | :--- |
| Go | [go.mod](#appendix-go-mod), [Dockerfile](#appendix-toolchain-dockerfile) | 1.27.2 |
| goimports | tool directive в go.mod, go.sum | golang.org/x/tools v0.52.0 |
| govulncheck | tool directive в go.mod, go.sum | golang.org/x/vuln v1.9.0 |
| golangci-lint | Dockerfile, официальный образ с digest | v2.14.0 |
| База toolchain | Dockerfile, индекс с AMD64/ARM64 | golang:1.27.2-trixie, sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e |
| Образ линтера | Dockerfile, индекс с AMD64/ARM64 | golangci/golangci-lint:v2.14.0, sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f |

Версии проверены по официальным источникам 15 октября 2026 года: [Go releases](https://go.dev/dl/), [golangci-lint releases](https://github.com/golangci/golangci-lint/releases/tag/v2.14.0), Go module proxy и Docker Registry. Оба закреплённых digest относятся к multi-platform индексам, а не к одному варианту CPU.

gRPC v1.84.0 добавлен для реального health API в P01-04; версия проверена через официальный Go proxy. JSON-логи использует стандартный log/slog. Инструменты закреплены через [Go tool directives](https://go.dev/doc/modules/managing-dependencies#tools). PostgreSQL использует pgx v5.11.0; Kafka/Redis библиотеки добавляются при первом использовании.

#### Основной путь Windows и macOS

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

`data-build` собирает host CLI для [подготовки PostgreSQL](#p02-docs-development-postgres-md); `data-cross-check` проверяет компиляцию для Windows AMD64, Linux AMD64/ARM64 и macOS AMD64/ARM64. Host OS/архитектура задаются через GOOS/GOARCH.

Toolchain — временный сервис профиля tools. Явный `compose run toolchain` активирует его; команды не поднимают PostgreSQL/Kafka/Redis или runtime сервисов. Shell-скрипт выполняется внутри Linux-контейнера. Исходники подключаются из относительного пути, кеши модулей и сборки — в отдельных named volumes. Используются стандартные capabilities Docker для записи в bind mount; privileged mode и Docker socket не включены, no-new-privileges включён. Ограничения runtime-контейнеров сервисов описаны отдельно в [контейнерном каркасе](#appendix-container-runtime).

Внутри образа Git доверяет только рабочему каталогу /workspace: bind mount может иметь другого владельца, а Go читает Git для настоящих VCS build metadata. Настройка действует в образе инструментов и не меняет Git-конфигурацию хоста.

#### CI и границы доказательств

[Go workflow](#appendix-go-workflow) использует нативные ubuntu-24.04 (AMD64) и ubuntu-24.04-arm (ARM64). На каждой архитектуре он собирает один и тот же Dockerfile, выполняет check, service-smoke и проверяет отказ на неизвестную команду. Затем проверяет отсутствие изменений отслеживаемых исходников. Успешная сборка ARM64 в CI не доказывает запуск Docker Desktop на личном Mac.

На Windows установлены host Go1.27.2, WSL2.7.13 и Docker Desktop4.95. [Runtime smoke и shutdown](#windows-smoke-proof) PASS; основной путь сборки остаётся контейнерным. Личная macOS не проверена, performance capacity не измерена. После работы Docker Desktop/WSL VM выключаются; установленные инструменты и полезные кеши сохраняются.

Версии, директивы tools и образы обновляются совместно отдельным проверяемым изменением. GOTOOLCHAIN=local запрещает неявную загрузку другого Go, контейнер использует GOFLAGS=-mod=readonly. После изменения go.mod/go.sum контейнер пересобирается. Текущая защита main и required checks отражаются в [безопасности репозитория](#appendix-repository-security), а не выводятся из наличия workflow.

<a id="appendix-go-mod"></a>

## Приложение: Go-модуль и версии инструментов

Полный исходный документ: `go.mod`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
module github.com/Includeoyasi/faultline-lab

go 1.27.2

tool (
	golang.org/x/tools/cmd/goimports
	golang.org/x/vuln/cmd/govulncheck
)

require (
	github.com/jackc/pgx/v5 v5.11.0
	golang.org/x/sys v0.49.0
	google.golang.org/grpc v1.84.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/telemetry v0.0.0-20261010193011-f71dc92da45e // indirect
	golang.org/x/text v0.43.0 // indirect
	golang.org/x/tools v0.52.0 // indirect
	golang.org/x/vuln v1.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
```

<a id="appendix-go-sum"></a>

## Приложение: Контрольные суммы модулей

Полный исходный документ: `go.sum`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```text
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
github.com/google/go-cmdtest v0.4.1-0.20220921163831-55ab3332a786 h1:rcv+Ippz6RAtvaGgKxc+8FQIpxHgsF+HBzPyYL2cyVU=
github.com/google/go-cmdtest v0.4.1-0.20220921163831-55ab3332a786/go.mod h1:apVn/GCasLZUVpAJ6oWAuyP7Ne7CEsQbTnc0plM3m+o=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/renameio v0.1.0 h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=
github.com/google/renameio v0.1.0/go.mod h1:KWCgfxg9yswjAJkECMjeO8J8rahYeXnNhOm40UhjYkI=
github.com/jackc/pgpassfile v1.0.0 h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=
github.com/jackc/pgpassfile v1.0.0/go.mod h1:CEx0iS5ambNFdcRtxPj5JhEz+xB6uRky5eyVu/W2HEg=
github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 h1:iCEnooe7UlwOQYpKFhBabPMi4aNAfoODPEFNiAnClxo=
github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761/go.mod h1:5TJZWKEWniPve33vlWYSoGYefn3gLQRzjfDlhSJ9ZKM=
github.com/jackc/pgx/v5 v5.11.0 h1:IzBBtyK9AHqf98cctWFifYSci2hgQR/cd56wB4p+ogg=
github.com/jackc/pgx/v5 v5.11.0/go.mod h1:mal1tBGAFfLHvZzaYh77YS/eC6IX9OWbRV1QIIM0Jn4=
github.com/jackc/puddle/v2 v2.2.2 h1:PR8nw+E/1w0GLuRFSmiioY6UooMp6KJv0/61nB7icHo=
github.com/jackc/puddle/v2 v2.2.2/go.mod h1:vriiEXHvEE654aYKXXjOvZM39qJ0q+azkZFrfEOc3H4=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
golang.org/x/mod v0.42.0 h1:ICeL3t2L30ljSPHHxDnR4fqVcBt/zR2zJAb9juavmQg=
golang.org/x/mod v0.42.0/go.mod h1:u5fzIrRInkL8MaO1cbRkUHT64v303YIT5KEqpNr7jRs=
golang.org/x/net v0.61.0 h1:RR+6j/BTBCrggOLugxa+66R8FvoN6MHlzDlZIbN9YFw=
golang.org/x/net v0.61.0/go.mod h1:WgDCOTH8iwtB66T4/GfI94AHo/hhdb7Hah9NDseQfiI=
golang.org/x/sync v0.24.0 h1:iqiDrGoDgyuvAKrhLQthUEzU/R378N78gXbBjxed7oA=
golang.org/x/sync v0.24.0/go.mod h1:sUUOizhqBxiL6pEWpqNLUiaJn1ShEbZ6BBqskPbjZm0=
golang.org/x/sys v0.49.0 h1:XbzkgYJHdqh/8m2Uu0W/dQv8nktxx4BFHp1M0gROTXA=
golang.org/x/sys v0.49.0/go.mod h1:hNLxWAXmnKAxqDtdwIYC4bM9oQPEecfsnNMuSxOs3og=
golang.org/x/telemetry v0.0.0-20261010193011-f71dc92da45e h1:P6NGsugec4oASJ+yKMNgTndeo0h2CbixZLvQIe9h5R4=
golang.org/x/telemetry v0.0.0-20261010193011-f71dc92da45e/go.mod h1:GNu+nb+E5axfGmP+Uu3ZVQu5lE/Xk42MjZOHWUD94/o=
golang.org/x/text v0.43.0 h1:1QivrlhwAsnMOcqgOdl9my7s9OyZl8O5/UGUTJVXoWc=
golang.org/x/text v0.43.0/go.mod h1:hbSIYA/amXcRXTFSZqh8tWHEoqRvvbKPJ2uz/KU+a+s=
golang.org/x/tools v0.52.0 h1:ryxg4U63oaUYk7uOLTF37n43f4qgtGvqV4Ioit1BaEs=
golang.org/x/tools v0.52.0/go.mod h1:ItcP5UNaVGG+Yc4eIYmaFLNJ4soEopkx5Y2d4lcEurs=
golang.org/x/tools/go/expect v0.1.1-deprecated h1:jpBZDwmgPhXsKZC6WhL20P4b/wmnpsEAGHaNy0n/rJM=
golang.org/x/tools/go/expect v0.1.1-deprecated/go.mod h1:eihoPOH+FgIqa3FpoTwguz/bVUSGBlGQU67vpBeOrBY=
golang.org/x/tools/go/packages/packagestest v0.1.1-deprecated h1:1h2MnaIAIXISqTFKdENegdpAgUXz6NrPEsbIeWaBRvM=
golang.org/x/tools/go/packages/packagestest v0.1.1-deprecated/go.mod h1:RVAQXBGNv1ib0J382/DPCRS/BPnsGebyM1Gj5VSDpG8=
golang.org/x/vuln v1.9.0 h1:12VaqnI9kJfqwRoPptxoqsUhQtdfrSbdksT7X/KY/sc=
golang.org/x/vuln v1.9.0/go.mod h1:E/Frwm0xp49fTyBW+6CQHdkvobknbwUM49yEpMGUauw=
gonum.org/v1/gonum v0.17.0 h1:VbpOemQlsSMrYmn7T2OUvQ4dqxQXU+ouZFQsZOx50z4=
gonum.org/v1/gonum v0.17.0/go.mod h1:El3tOrEuMpv2UdMrbNlKEh9vd86bmQ6vqIcDwxEOc1E=
google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 h1:qEHAMpSaUhtD0p3NbEEI83HwNGFxEwaSJ1G9PLnCBZE=
google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800/go.mod h1:4Hqkh8ycfw05ld/3BWL7rJOSfebL2Q+DVDeRgYgxUU8=
google.golang.org/grpc v1.84.0 h1:soMyaPJ8pAak5PIQ0DGBUir0XRo2fRoMqhNWMLlLxO0=
google.golang.org/grpc v1.84.0/go.mod h1:ljCht0DrxQrXBDRTZp52Qxh3Ffk8CdYm2sj4O2QN2C0=
google.golang.org/protobuf v1.36.11 h1:fV6ZwhNocDyBLK0dj+fg8ektcVegBBuEolpbTQyBNVE=
google.golang.org/protobuf v1.36.11/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
```

<a id="appendix-orders-main"></a>

## Приложение: Orders gRPC health entrypoint

Полный исходный документ: `cmd/orders/main.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Includeoyasi/faultline-lab/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	exitCode := service.Execute(ctx, "orders", os.Args[1:], os.LookupEnv, os.Stderr)
	stop()
	os.Exit(exitCode)
}
```

<a id="appendix-inventory-main"></a>

## Приложение: Inventory gRPC health entrypoint

Полный исходный документ: `cmd/inventory/main.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Includeoyasi/faultline-lab/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	exitCode := service.Execute(ctx, "inventory", os.Args[1:], os.LookupEnv, os.Stderr)
	stop()
	os.Exit(exitCode)
}
```

<a id="appendix-payments-main"></a>

## Приложение: Payments gRPC health entrypoint

Полный исходный документ: `cmd/payments/main.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Includeoyasi/faultline-lab/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	exitCode := service.Execute(ctx, "payments", os.Args[1:], os.LookupEnv, os.Stderr)
	stop()
	os.Exit(exitCode)
}
```

<a id="appendix-toolchain-dockerfile"></a>

## Приложение: Единый multi-stage Dockerfile: toolchain и runtime

Полный исходный документ: `build/Dockerfile`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```dockerfile
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
COPY cmd ./cmd
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
COPY internal ./internal
COPY migrations ./migrations
RUN CGO_ENABLED=0 go build -trimpath -o /out/labdata ./cmd/labdata \
    && CGO_ENABLED=0 go test -c -o /out/database.test ./internal/database

FROM scratch AS labdata
COPY --from=labdata-build /out/labdata /app/labdata
COPY --from=labdata-build /out/database.test /app/database.test
USER 65532:65532
ENTRYPOINT ["/app/labdata"]
```

<a id="appendix-compose-toolchain"></a>

## Приложение: Compose для временных инструментов

Полный исходный документ: `compose.yaml`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```yaml
x-runtime: &runtime
  env_file:
    - ${FL_RUNTIME_ENV_FILE:-configs/env/container.env.example}
  read_only: true
  cap_drop: [ALL]
  security_opt:
    - no-new-privileges:true
  stop_grace_period: ${FL_STOP_GRACE_PERIOD:-15s}
  networks: [runtime]
  healthcheck:
    test: [CMD, /app/service, healthcheck]
    interval: 5s
    timeout: 3s
    retries: 3
    start_period: 5s

x-service-build: &service-build
  context: .
  dockerfile: build/Dockerfile
  target: runtime

services:
  orders:
    <<: *runtime
    build:
      <<: *service-build
      args:
        SERVICE: orders
  inventory:
    <<: *runtime
    build:
      <<: *service-build
      args:
        SERVICE: inventory
  payments:
    <<: *runtime
    build:
      <<: *service-build
      args:
        SERVICE: payments
  toolchain:
    profiles: [tools]
    build:
      context: .
      dockerfile: build/Dockerfile
      target: toolchain
    working_dir: /workspace
    volumes:
      - .:/workspace
      - go-modules:/go/pkg/mod
      - go-build:/root/.cache/go-build
    security_opt:
      - no-new-privileges:true

volumes:
  go-modules:
  go-build:

networks:
  runtime:
    internal: true
```

<a id="appendix-go-commands"></a>

## Приложение: Команды внутри Linux-контейнера

Полный исходный документ: `tools/go.sh`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```sh
#!/bin/sh
# This script runs inside the Linux toolchain container, not on the host.
set -eu

if [ "$#" -gt 1 ]; then
  echo "usage: lab-go {versions|build|check|format|service-smoke|smoke-build|smoke-cross-check|smoke-cli-tests|data-build|data-cross-check}" >&2
  exit 2
fi

versions() {
  go version
  golangci-lint --version
  go version -m "$(command -v goimports)"
  govulncheck -version
}

build() {
  mkdir -p bin
  for service in orders inventory payments; do
    CGO_ENABLED=0 go build -trimpath -o "bin/$service" "./cmd/$service"
    go version -m "bin/$service"
  done
}

host_build() {
  command_name=$1
  target_os=${GOOS:-$(go env GOOS)}
  target_arch=${GOARCH:-$(go env GOARCH)}
  case "$target_os/$target_arch" in
    windows/amd64|linux/amd64|linux/arm64|darwin/amd64|darwin/arm64) ;;
    *) echo "unsupported smoke target" >&2; exit 2 ;;
  esac
  extension=
  if [ "$target_os" = windows ]; then extension=.exe; fi
  mkdir -p bin
  CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" go build -trimpath -o "bin/$command_name-$target_os-$target_arch$extension" "./cmd/$command_name"
}

source_files() {
  file_list=$(mktemp)
  trap 'rm -f "$file_list"' EXIT
  find . -type d \( -name .git -o -name .local -o -name vendor -o -name bin -o -name tmp -o -name artifacts \) -prune \
    -o -type f -name '*.go' -print0 > "$file_list"
}

case "${1:-versions}" in
  versions)
    versions
    ;;
  build)
    build
    ;;
  check)
    versions
    source_files
    formatting=$(xargs -0 -r gofmt -l < "$file_list")
    imports=$(xargs -0 -r goimports -l < "$file_list")
    if [ -n "$formatting$imports" ]; then
      printf 'Formatting/import changes required:\n%s\n%s\n' "$formatting" "$imports" >&2
      exit 1
    fi
    go mod tidy -diff
    go mod verify
    go build ./...
    golangci-lint run --timeout=5m
    go test -count=1 -timeout=3m ./...
    go test -race -count=1 -timeout=5m ./...
    govulncheck ./...
    ;;
  format)
    source_files
    xargs -0 -r goimports -w < "$file_list"
    ;;
  service-smoke)
    build
    LAB_BIN_DIR=/workspace/bin go test -count=1 -v -timeout=2m ./tests/smoke
    ;;
  smoke-build)
    host_build smoke
    ;;
  smoke-cross-check)
    for target in windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
      GOOS=${target%/*} GOARCH=${target#*/} host_build smoke
    done
    ;;
  data-build)
    host_build labdata
    ;;
  data-cross-check)
    for target in windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
      GOOS=${target%/*} GOARCH=${target#*/} host_build labdata
    done
    ;;
  smoke-cli-tests)
    mkdir -p bin
    CGO_ENABLED=0 go test -c -o bin/smoke-cli-tests ./tests/smokecli
    ;;
  *)
    echo "usage: lab-go {versions|build|check|format|service-smoke|smoke-build|smoke-cross-check|smoke-cli-tests|data-build|data-cross-check}" >&2
    exit 2
    ;;
esac
```

<a id="appendix-golangci-config"></a>

## Приложение: Согласованные Go-линтеры

Полный исходный документ: `.golangci.yml`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```yaml
version: "2"

linters:
  default: none
  enable:
    - govet
    - staticcheck
    - errcheck
    - ineffassign
    - unused
```

<a id="appendix-go-workflow"></a>

## Приложение: Проверки Go в нативном Linux CI

Полный исходный документ: `.github/workflows/go.yml`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```yaml
name: Go toolchain

on:
  push:
  pull_request:

permissions:
  contents: read

concurrency:
  group: go-${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  toolchain:
    name: Go toolchain (${{ matrix.arch }})
    runs-on: ${{ matrix.runner }}
    timeout-minutes: 20
    strategy:
      fail-fast: false
      matrix:
        include:
          - arch: amd64
            runner: ubuntu-24.04
          - arch: arm64
            runner: ubuntu-24.04-arm
    steps:
      - uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7
        with:
          persist-credentials: false
      - name: Show Docker prerequisites
        run: |
          docker version
          docker compose version
      - name: Build pinned toolchain
        run: docker compose build toolchain
      - name: Run Go quality checks
        run: docker compose run --rm toolchain check
      - name: Verify service processes, health and isolation
        run: docker compose run --rm toolchain service-smoke
      - name: Reject invalid toolchain command
        run: |
          if docker compose run --rm toolchain invalid-command; then
            status=0
          else
            status=$?
          fi
          test "$status" -eq 2
      - name: Build, start, inspect and stop runtime containers
        run: |
          docker compose run --rm toolchain smoke-build
          bin/smoke-linux-${{ matrix.arch }}
      - name: Verify SIGINT report and scoped cleanup on the host
        run: |
          docker compose run --rm toolchain smoke-cli-tests
          LAB_SMOKE_BINARY="$PWD/bin/smoke-linux-${{ matrix.arch }}" LAB_SMOKE_ROOT="$PWD" bin/smoke-cli-tests -test.v -test.timeout=4m
      - name: Cross-compile smoke for supported host platforms
        if: matrix.arch == 'amd64'
        run: docker compose run --rm toolchain smoke-cross-check
      - name: Verify PostgreSQL roles, migrations, seed and persistence
        run: |
          docker compose run --rm toolchain data-build
          bin/labdata-linux-${{ matrix.arch }} smoke
      - name: Cross-compile data tool for supported hosts
        if: matrix.arch == 'amd64'
        run: docker compose run --rm toolchain data-cross-check
      - name: Retain smoke reports and diagnostics
        if: always()
        uses: actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9 # v7.0.2
        with:
          name: runtime-smoke-${{ matrix.arch }}-${{ github.run_id }}-${{ github.run_attempt }}
          path: artifacts/smoke/
          if-no-files-found: warn
          retention-days: 14
      - name: Retain PostgreSQL evidence
        if: always()
        uses: actions/upload-artifact@cf430e030ddbb5b0abf93d22962f4752f3646cd9 # v7.0.2
        with:
          name: postgres-${{ matrix.arch }}-${{ github.run_id }}-${{ github.run_attempt }}
          path: artifacts/postgres/**/report.json
          if-no-files-found: warn
          retention-days: 14
      - name: Ensure checks did not change source files
        run: git diff --exit-code
```

<a id="appendix-dockerignore"></a>

## Приложение: Исключения Docker build context

Полный исходный документ: `.dockerignore`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```text
.git
.local
.env
.env.*
**/.env
**/.env.*
**/*.env
*.pem
*.key
credentials.json
bin
dist
coverage
data
artifacts
tmp
node_modules
```

<a id="appendix-service-runtime"></a>

## Приложение: Конфигурация, JSON-логи и gRPC health: поведение и проверки

Полный исходный документ: `docs/development/SERVICE_RUNTIME.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Сервисы: конфигурация, логи и жизненный цикл

Три независимых процесса Orders, Inventory и Payments проверяют настройки, пишут JSON в stderr и обслуживают стандартный `grpc.health.v1.Health/Check`. Health пока подтверждает только работу gRPC-каркаса. Бизнес-методов, PostgreSQL/Kafka/Redis и доменных protobuf здесь нет.

Конфигурация и базовый процесс введены в PR #6/#7; контейнерный жизненный цикл — в [PR #8](https://github.com/Includeoyasi/faultline-lab/pull/8). Это исторические источники реализации. Актуальные задачи и ревью — в [GitHub](https://github.com/Includeoyasi/faultline-lab/issues).

#### Как устроен запуск

```mermaid
flowchart LR
    Env["Окружение + defaults"] --> Validate["Проверка конфигурации"]
    Validate -->|ошибка| Error["JSON ERROR → exit 1"]
    Validate --> Bind["Открыть TCP listener"]
    Bind -->|ошибка| Error
    Bind --> Ready["Health SERVING<br/>JSON: service started"]
    Ready --> Serve["Обслуживать gRPC health"]
    Serve -->|SIGINT / SIGTERM| Unready["Health NOT_SERVING"]
    Unready --> Drain["GracefulStop<br/>ожидание активных RPC"]
    Drain -->|завершены| Stop["JSON: stop_mode=graceful<br/>exit 0"]
    Drain -->|истёк timeout| Force["Stop + завершение процесса<br/>stop_mode=forced / exit 1"]
```

[cmd/orders](#appendix-orders-main), [cmd/inventory](#appendix-inventory-main) и [cmd/payments](#appendix-payments-main) задают имя сервиса и владеют обработкой сигналов/кодом выхода. Адресов и портов в прикладном Go-коде нет: они поступают из окружения. [internal/service](#appendix-service-run) содержит только общие механизмы запуска; бизнес-логика одного сервиса не импортируется другим. Лаборатория и адаптеры ещё не реализованы.

#### Настройки

Одинаковые имена переменных действуют в окружении каждого отдельного процесса. Файлы `.env` приложение самостоятельно не читает.

| Переменная | Default | Допустимые значения |
| :--- | :--- | :--- |
| `FL_GRPC_ADDR` | Нет; обязательная настройка окружения | `host:port`, порт 1–65535; IPv6 в квадратных скобках; пустой host разрешён только как явно выбранный wildcard |
| `FL_LOG_LEVEL` | `INFO` | `DEBUG`, `INFO`, `WARN`, `ERROR`, регистр не важен |
| `FL_SHUTDOWN_TIMEOUT` | `10s` | Положительная Go duration, например `500ms` или `15s` |
| `FL_HEALTH_TARGET` | Нет; обязателен только для команды `healthcheck` | Конкретный host:port, не wildcard; задаётся отдельно от listener |

Отсутствующий или пустой `FL_GRPC_ADDR` приводит к JSON-ошибке и exit 1 до открытия listener. Приложение не выбирает сеть за оператора и не имеет скрытого адреса по умолчанию. Для уровня логов и timeout действуют defaults; явно пустое значение отклоняется. Адрес проверяется синтаксически без DNS-запроса. Имя host должно быть IP или DNS-именем; IPv6 zone ID пока не поддерживается. Возможность разрешить имя/открыть адрес проверяется фактическим bind. URI, имена портов и произвольный текст отклоняются. Ошибка называет переменную и правило, не повторяя введённое значение; весь конфиг в логи не выводится.

`FL_SHUTDOWN_TIMEOUT` ограничивает ожидание активных RPC после удаления готовности. Graceful stop даёт завершиться начатой работе; по deadline запрашивается принудительный Stop и процесс завершается с кодом 1. Оставшаяся goroutine некооперативного handler прекращается вместе с процессом. Реализация и управляемые тесты описаны в [контейнерном пакете](#appendix-container-runtime).

Адрес прослушивания, адрес подключения клиента и публикация порта на хосте — разные настройки:

| Назначение | Где задаётся | Пример / правило |
| :--- | :--- | :--- |
| Listener локального процесса | [local.env.example](#appendix-local-env-example) | `localhost:50051`; для нескольких процессов на одном хосте задаются разные порты |
| Listener в контейнере | [container.env.example](#appendix-container-env-example) | `:50051` — явный wildcard внутри сетевого пространства контейнера; разные контейнеры могут использовать один внутренний порт |
| Endpoint другого сервиса | Будущая deployment-конфигурация, когда появятся клиенты | Имена Compose DNS, например `inventory:50051`; IP контейнеров не закрепляются |
| Доступ с хоста | Явная публикация при запуске | Для локальной разработки привязка только к loopback; без `-p` порт не публикуется на интерфейсах хоста |

Образцы содержат только несекретные настройки. Личные копии `*.env` игнорируются Git и исключены из Docker build context; окружение передаёт launcher, приложение не загружает dotenv самостоятельно. Loopback в тестах — изолированный fixture с динамическим портом, а в наблюдённых логах — фактический адрес запущенного тестового процесса.

Использование DNS и настройка адреса не заменяют разграничение сетевого доступа и аутентификацию. Сейчас это локальный gRPC-стенд без TLS/аутентификации; production-размещение потребует отдельной политики TLS/mTLS и доступа. Привязка wildcard в контейнере не является такой политикой. Правила публикации и разрешения имён основаны на [Docker port publishing](https://docs.docker.com/engine/network/port-publishing/) и [Compose networking](https://docs.docker.com/compose/how-tos/networking/).

#### Логи

Используется стандартный [log/slog](https://pkg.go.dev/log/slog), отдельный logger каждого сервиса. Каждая прикладная запись — одна JSON-строка с `time`, `level`, `msg`, `service`.

- `service started` (INFO): фактический `grpc_address`, `version`, `revision`, `vcs_modified`, `go_version` из Go build metadata. Неизвестное значение — `unknown`; локальный модуль может иметь версию `(devel)`. Commit не задаётся придуманной константой.
- `service stopping` (INFO): health=NOT_SERVING и shutdown_timeout перед началом drain.
- `service stopped` (INFO): stop_mode=graceful; при deadline — WARN с stop_mode=forced и последующая ERROR-диагностика, exit 1.
- `configuration rejected` (ERROR): безопасная причина ошибки настроек; до создания рабочего logger используется INFO logger, поэтому ошибочный уровень не скрывает диагностику.
- `service failed` (ERROR): ошибка bind или обслуживания/остановки.

Уровень фильтрует записи: при WARN/ERROR штатные старт/стоп INFO не выводятся. По умолчанию INFO они доступны. Секреты и payload не добавляются. Идентификаторы запросов/заказов появятся вместе с соответствующими операциями.

#### Наблюдаемые записи

Каждый runtime smoke сохраняет реальные JSONL-записи старта, удаления готовности и остановки. В review packet приводятся ссылка на соответствующий CI и его артефакты; [контейнерная инструкция](#appendix-container-runtime) описывает состав отчёта. Process smoke также печатает наблюдённые JSON в лог job. Временный merge ref GitHub может отличаться от HEAD ветки; metadata не подменяются придуманным SHA.

#### Сборка и проверки

Основной путь требует Docker с Linux-контейнерами и Compose v2; Go/Bash на хосте не нужны. Из корня:

```text
docker compose build toolchain
docker compose run --rm toolchain check
docker compose run --rm toolchain service-smoke
```

`check` проверяет формат, модули, сборку, линтеры, unit/process tests, race detector и govulncheck. `service-smoke` собирает все три бинарника и проверяет их как отдельные процессы внутри одного временного контейнера инструментов. Это проверка изоляции процессов, а не готовых runtime-контейнеров. Прежняя команда `bootstrap-smoke` удалена: теперь программы продолжают обслуживать запросы.

[Process smoke](#appendix-service-process-tests) проверяет:

- запуск всех трёх процессов и `SERVING` для общего health (`service=""`);
- `NotFound` для неизвестного имени health-сервиса;
- обязательные поля JSON и build metadata;
- для каждого бинарника — отсутствующий/пустой адрес, ошибки адреса, уровня, timeout и занятого порта: одна JSON-диагностика, exit 1, отсутствие готовности;
- после остановки Orders — доступность Inventory и Payments;
- на Linux — нормальные SIGINT/SIGTERM exit и forced exit 1 по deadline при удерживаемом health Watch.

Тесты синхронизируются по стартовой записи и реальному RPC с deadline, без случайных Sleep. Порты тестов выделяются динамически; unit fixture использует порт 0, публичная конфигурация его отклоняет. Под Windows smoke завершает дочерние процессы через Kill: это проверка изоляции, а не доставки Unix-сигнала. Завершение по отмене контекста и JSON stop отдельно проверяются unit-тестом на доступных ОС; Linux CI дополнительно проверяет сигнал реальному процессу. Race detector охватывает также сам сервер в unit-тесте; дочерние бинарники smoke собираются обычной сборкой.

Для отдельных минимальных runtime-контейнеров используйте [контейнерную инструкцию](#appendix-container-runtime): сборка, запуск, health, остановка и автоматический smoke без host Go. Она заменяет прежний запуск серверов внутри образа компилятора.

#### Границы доказательств

Нативный Windows Go smoke может подтвердить TCP/health процессов, но не Docker Desktop. Linux AMD64/ARM64 CI проверяет toolchain-контейнер и процессы на каждой архитектуре; он не доказывает запуск на личном Mac. Graceful drain проверяется управляемыми unit/process tests; runtime-контейнеры — отдельным Docker smoke. Работа БД, бизнес-инварианты, нагрузка и потребление ресурсов ещё не проверяются. Итоговые PASS/NOT_RUN фиксируются после реального прогона, отдельно от этого описания поведения.

<a id="appendix-service-config"></a>

## Приложение: Валидация конфигурации

Полный исходный документ: `internal/service/config.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
// Package service contains the shared startup mechanics of the three services.
package service

import (
	"errors"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"
)

// Config holds process settings, including the active RPC drain budget.
type Config struct {
	GRPCAddress     string
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

// LoadConfig requires an explicit listen address and applies defaults for the
// remaining settings. Supplied values are not included in validation errors.
func LoadConfig(lookup func(string) (string, bool)) (Config, error) {
	cfg := Config{LogLevel: slog.LevelInfo, ShutdownTimeout: 10 * time.Second}
	cfg.GRPCAddress, _ = lookup("FL_GRPC_ADDR")
	if !validAddress(cfg.GRPCAddress) {
		return Config{}, errors.New("FL_GRPC_ADDR is required and must be host:port with a port from 1 to 65535")
	}
	if value, ok := lookup("FL_LOG_LEVEL"); ok {
		switch strings.ToUpper(value) {
		case "DEBUG":
			cfg.LogLevel = slog.LevelDebug
		case "INFO":
			cfg.LogLevel = slog.LevelInfo
		case "WARN":
			cfg.LogLevel = slog.LevelWarn
		case "ERROR":
			cfg.LogLevel = slog.LevelError
		default:
			return Config{}, errors.New("FL_LOG_LEVEL must be DEBUG, INFO, WARN or ERROR")
		}
	}
	if value, ok := lookup("FL_SHUTDOWN_TIMEOUT"); ok {
		timeout, err := time.ParseDuration(value)
		if err != nil || timeout <= 0 {
			return Config{}, errors.New("FL_SHUTDOWN_TIMEOUT must be a positive Go duration")
		}
		cfg.ShutdownTimeout = timeout
	}
	return cfg, nil
}

func validAddress(address string) bool {
	host, port, err := net.SplitHostPort(address)
	if err != nil || port == "" {
		return false
	}
	for _, char := range port {
		if char < '0' || char > '9' {
			return false
		}
	}
	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return false
	}
	if host == "" || net.ParseIP(host) != nil {
		return true
	}
	// Validate DNS syntax locally; startup must not require a DNS lookup.
	host = strings.TrimSuffix(host, ".")
	if len(host) == 0 || len(host) > 253 {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, char := range label {
			switch {
			case char >= 'a' && char <= 'z', char >= 'A' && char <= 'Z', char >= '0' && char <= '9', char == '-':
			default:
				return false
			}
		}
	}
	return true
}
```

<a id="appendix-service-logging"></a>

## Приложение: JSON-логи и build metadata

Полный исходный документ: `internal/service/log.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
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
```

<a id="appendix-service-run"></a>

## Приложение: Запуск gRPC health и немедленная остановка

Полный исходный документ: `internal/service/run.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// ErrShutdownTimeout indicates that active RPCs required a forced shutdown.
var ErrShutdownTimeout = errors.New("shutdown timeout exceeded; active RPCs were interrupted")

// Run validates configuration, serves gRPC health and drains active RPCs on
// cancellation. It returns ErrShutdownTimeout if the drain budget expires.
func Run(ctx context.Context, name string, lookup func(string) (string, bool), output io.Writer) error {
	logger := newLogger(output, name, slog.LevelInfo)
	cfg, err := LoadConfig(lookup)
	if err != nil {
		logger.ErrorContext(ctx, "configuration rejected", "error", err.Error())
		return err
	}
	logger = newLogger(output, name, cfg.LogLevel)
	if err := serve(ctx, cfg, logger); err != nil {
		logger.ErrorContext(ctx, "service failed", "error", err.Error())
		return err
	}
	return nil
}

func serve(ctx context.Context, cfg Config, logger *slog.Logger) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	listener, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		return fmt.Errorf("listen for gRPC: %w", err)
	}
	// grpc.Server.Serve owns and closes the listener, including on Stop.
	server := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	healthpb.RegisterHealthServer(server, healthServer)
	return serveGRPC(ctx, listener, server, healthServer, cfg.ShutdownTimeout, logger)
}

// serveGRPC owns the server, listener and all shutdown goroutines. The explicit
// server parameter also lets lifecycle tests register a controlled RPC handler.
func serveGRPC(ctx context.Context, listener net.Listener, server *grpc.Server, healthServer *health.Server, timeout time.Duration, logger *slog.Logger) error {
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	info, ok := debug.ReadBuildInfo()
	attrs := append(buildAttributes(info, ok), slog.String("grpc_address", listener.Addr().String()))
	logger.LogAttrs(ctx, slog.LevelInfo, "service started", attrs...)
	serveResult := make(chan error, 1)
	go func() { serveResult <- server.Serve(listener) }()
	select {
	case err := <-serveResult:
		server.Stop()
		if err == nil {
			err = errors.New("server stopped unexpectedly")
		}
		return fmt.Errorf("serve gRPC: %w", err)
	case <-ctx.Done():
		healthServer.Shutdown()
		logger.InfoContext(ctx, "service stopping", "health", "NOT_SERVING", "shutdown_timeout", timeout.String())
		drained := make(chan struct{})
		go func() { server.GracefulStop(); close(drained) }()
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		select {
		case <-drained:
		case <-timer.C:
			// Give a completed drain precedence if completion and deadline race.
			select {
			case <-drained:
			default:
				// GracefulStop can wait for a handler that ignores cancellation.
				// Do not join it on this fatal path: Run returns an error and the
				// entrypoint exits the process, terminating remaining RPC work.
				go server.Stop()
				logger.WarnContext(ctx, "service stopped", "stop_mode", "forced")
				return ErrShutdownTimeout
			}
		}
		if err := <-serveResult; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("stop gRPC: %w", err)
		}
		logger.InfoContext(ctx, "service stopped", "stop_mode", "graceful")
		return nil
	}
}
```

<a id="appendix-service-config-tests"></a>

## Приложение: Проверки конфигурации

Полный исходный документ: `internal/service/config_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package service

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func environment(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) { value, ok := values[key]; return value, ok }
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{"non-network defaults", map[string]string{"FL_GRPC_ADDR": "localhost:50051"}, Config{"localhost:50051", slog.LevelInfo, 10 * time.Second}},
		{"overrides", map[string]string{"FL_GRPC_ADDR": "0.0.0.0:8080", "FL_LOG_LEVEL": "debug", "FL_SHUTDOWN_TIMEOUT": "250ms"}, Config{"0.0.0.0:8080", slog.LevelDebug, 250 * time.Millisecond}},
		{"ipv6", map[string]string{"FL_GRPC_ADDR": "[::1]:8080", "FL_LOG_LEVEL": "WARN"}, Config{"[::1]:8080", slog.LevelWarn, 10 * time.Second}},
		{"dns", map[string]string{"FL_GRPC_ADDR": "orders.local:50051", "FL_LOG_LEVEL": "ERROR"}, Config{"orders.local:50051", slog.LevelError, 10 * time.Second}},
		{"wildcard", map[string]string{"FL_GRPC_ADDR": ":50051"}, Config{":50051", slog.LevelInfo, 10 * time.Second}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := LoadConfig(environment(tt.env))
			if err != nil || got != tt.want {
				t.Fatalf("LoadConfig() = %+v, %v; want %+v", got, err, tt.want)
			}
		})
	}
}

func TestLoadConfigRejectsInvalidValues(t *testing.T) {
	t.Parallel()
	inputs := map[string][]string{
		"FL_GRPC_ADDR":        {"", "localhost", "::1:8080", "localhost:http", ":0", ":65536", ":-1", ":+80", "bad host:80", "-host:80", "host..local:80", "host_:80", "secret://user:password@host:80"},
		"FL_LOG_LEVEL":        {"", "TRACE", "INFO+1", "super-secret-value"},
		"FL_SHUTDOWN_TIMEOUT": {"", "0s", "-1s", "10", "999999999999999999999h", "super-secret-value"},
	}
	for key, values := range inputs {
		for _, value := range values {
			t.Run(key+"/"+value, func(t *testing.T) {
				t.Parallel()
				values := map[string]string{"FL_GRPC_ADDR": "localhost:50051"}
				values[key] = value
				_, err := LoadConfig(environment(values))
				if err == nil || !strings.Contains(err.Error(), key) {
					t.Fatalf("expected diagnostic naming %s, got %v", key, err)
				}
				if value != "" && strings.Contains(err.Error(), value) {
					t.Fatalf("diagnostic echoed the supplied value: %v", err)
				}
			})
		}
	}
	if _, err := LoadConfig(environment(nil)); err == nil || !strings.Contains(err.Error(), "FL_GRPC_ADDR is required") {
		t.Fatalf("missing listen address must fail closed, got %v", err)
	}
}
```

<a id="appendix-service-log-tests"></a>

## Приложение: Проверки JSON и build metadata

Полный исходный документ: `internal/service/log_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"runtime/debug"
	"testing"
	"time"
)

func TestJSONLoggingAndLevel(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer
	logger := newLogger(&output, "orders", slog.LevelWarn)
	logger.Info("hidden")
	if output.Len() != 0 {
		t.Fatal("INFO was not filtered at WARN")
	}
	logger.Error("bind failed", "error", "address already in use")
	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if record["service"] != "orders" || record["level"] != "ERROR" || record["msg"] != "bind failed" || record["error"] == "" {
		t.Fatalf("unexpected record: %v", record)
	}
	if _, err := time.Parse(time.RFC3339Nano, record["time"].(string)); err != nil {
		t.Fatalf("invalid log timestamp: %v", err)
	}
}

func TestBuildAttributes(t *testing.T) {
	t.Parallel()
	for _, available := range []bool{false, true} {
		t.Run(map[bool]string{false: "unknown", true: "available"}[available], func(t *testing.T) {
			info := &debug.BuildInfo{GoVersion: "go1.test", Main: debug.Module{Version: "v0.test"}, Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "actual-revision"}, {Key: "vcs.modified", Value: "true"},
			}}
			var output bytes.Buffer
			newLogger(&output, "inventory", slog.LevelInfo).LogAttrs(context.Background(), slog.LevelInfo, "service started", buildAttributes(info, available)...)
			var record map[string]any
			if err := json.Unmarshal(output.Bytes(), &record); err != nil {
				t.Fatal(err)
			}
			for key, value := range map[string]string{"version": "v0.test", "revision": "actual-revision", "vcs_modified": "true", "go_version": "go1.test"} {
				if !available {
					value = "unknown"
				}
				if record[key] != value {
					t.Errorf("%s = %v, want %s", key, record[key], value)
				}
			}
		})
	}
}
```

<a id="appendix-service-run-tests"></a>

## Приложение: Проверки bind, TCP health и отмены контекста

Полный исходный документ: `internal/service/run_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestRunRejectsConfigurationBeforeListening(t *testing.T) {
	t.Parallel()
	var output bytes.Buffer
	err := Run(context.Background(), "payments", environment(map[string]string{"FL_GRPC_ADDR": "localhost:50053", "FL_LOG_LEVEL": "super-secret-value"}), &output)
	if err == nil || strings.Contains(output.String(), "super-secret-value") {
		t.Fatalf("expected sanitized configuration failure, got %v / %s", err, &output)
	}
	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record["service"] != "payments" || record["level"] != "ERROR" || record["msg"] != "configuration rejected" {
		t.Fatalf("unexpected diagnostic: %v", record)
	}
}

func TestRunRejectsOccupiedAddress(t *testing.T) {
	t.Parallel()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := listener.Close(); err != nil {
			t.Error(err)
		}
	})
	var output bytes.Buffer
	err = Run(context.Background(), "orders", environment(map[string]string{"FL_GRPC_ADDR": listener.Addr().String()}), &output)
	var operationError *net.OpError
	if !errors.As(err, &operationError) || operationError.Op != "listen" {
		t.Fatalf("expected wrapped listen error, got %v", err)
	}
	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	if record["service"] != "orders" || record["msg"] != "service failed" || record["level"] != "ERROR" {
		t.Fatalf("unexpected bind diagnostic: %v", record)
	}
}

func TestServeWithCancelledContextDoesNotAnnounceReadiness(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	err := serve(ctx, Config{GRPCAddress: "127.0.0.1:0"}, newLogger(&output, "orders", slog.LevelInfo))
	if !errors.Is(err, context.Canceled) || output.Len() != 0 {
		t.Fatalf("expected cancellation without startup, got %v / %s", err, &output)
	}
}

type recordWriter chan []byte

func (records recordWriter) Write(data []byte) (int, error) {
	records <- append([]byte(nil), data...)
	return len(data), nil
}

type writerFunc func([]byte) (int, error)

func (write writerFunc) Write(data []byte) (int, error) { return write(data) }

func TestServeCancellationDuringStartup(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var output bytes.Buffer
	writer := writerFunc(func(data []byte) (int, error) {
		// Cancel before the Serve goroutine is started. Stop may win that race;
		// grpc.ErrServerStopped is then an expected shutdown result.
		cancel()
		return output.Write(data)
	})
	if err := serve(ctx, Config{GRPCAddress: "127.0.0.1:0", ShutdownTimeout: time.Second}, newLogger(writer, "orders", slog.LevelInfo)); err != nil {
		t.Fatalf("cancellation during startup failed: %v", err)
	}
	if !strings.Contains(output.String(), `"msg":"service stopped"`) {
		t.Fatalf("missing final log: %s", &output)
	}
}

func TestServeHealthAndCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	records := make(recordWriter, 4)
	result := make(chan error, 1)
	go func() {
		// Port zero is only a test fixture; public configuration rejects it.
		result <- serve(ctx, Config{GRPCAddress: "127.0.0.1:0", ShutdownTimeout: time.Second}, newLogger(records, "orders", slog.LevelInfo))
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case err := <-result:
			if err != nil {
				t.Errorf("serve: %v", err)
			}
		case <-time.After(10 * time.Second):
			t.Error("server did not exit after cancellation")
		}
	})
	var started map[string]any
	select {
	case record := <-records:
		if err := json.Unmarshal(record, &started); err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("missing startup record")
	}
	address, ok := started["grpc_address"].(string)
	if !ok || started["msg"] != "service started" {
		t.Fatalf("invalid startup record: %v", started)
	}
	conn, err := grpc.NewClient("passthrough:///"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
	if err != nil || response.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("health = %v, %v", response, err)
	}
	cancel()
	select {
	case record := <-records:
		if !strings.Contains(string(record), `"health":"NOT_SERVING"`) {
			t.Fatalf("missing readiness transition: %s", record)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("missing stopping record")
	}
	select {
	case record := <-records:
		var stopped map[string]any
		if err := json.Unmarshal(record, &stopped); err != nil {
			t.Fatal(err)
		}
		if stopped["msg"] != "service stopped" || stopped["service"] != "orders" || stopped["stop_mode"] != "graceful" {
			t.Fatalf("invalid stop record: %v", stopped)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("missing stop record")
	}
	// Cleanup owns the single receive from result, so it always joins the server.
}
```

<a id="appendix-service-process-tests"></a>

## Приложение: Проверки трёх реальных процессов

Полный исходный документ: `tests/smoke/services_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
// Package smoke tests the actual service binaries over TCP, without databases.
package smoke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type logCapture struct {
	mu      sync.Mutex
	pending string
	records []map[string]any
	err     error
	started chan string
}

func (capture *logCapture) Write(data []byte) (int, error) {
	capture.mu.Lock()
	defer capture.mu.Unlock()
	capture.pending += string(data)
	for {
		line, rest, found := strings.Cut(capture.pending, "\n")
		if !found {
			break
		}
		capture.pending = rest
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			capture.err = fmt.Errorf("non-JSON process output: %w", err)
			continue
		}
		capture.records = append(capture.records, record)
		if record["msg"] == "service started" {
			address, _ := record["grpc_address"].(string)
			select {
			case capture.started <- address:
			default:
			}
		}
	}
	return len(data), nil
}

func (capture *logCapture) check(t *testing.T, service string) []map[string]any {
	t.Helper()
	capture.mu.Lock()
	defer capture.mu.Unlock()
	if capture.err != nil || capture.pending != "" {
		t.Fatalf("invalid process logs: %v; unfinished line %q", capture.err, capture.pending)
	}
	for _, record := range capture.records {
		level, hasLevel := record["level"].(string)
		message, hasMessage := record["msg"].(string)
		if record["service"] != service || !hasLevel || level == "" || !hasMessage || message == "" {
			t.Fatalf("missing common log fields: %v", record)
		}
		timestamp, ok := record["time"].(string)
		if !ok {
			t.Fatalf("missing timestamp: %v", record)
		}
		if _, err := time.Parse(time.RFC3339Nano, timestamp); err != nil {
			t.Fatalf("invalid timestamp: %v", err)
		}
	}
	return append([]map[string]any(nil), capture.records...)
}

type process struct {
	cmd     *exec.Cmd
	logs    *logCapture
	done    chan struct{}
	waitErr error // Read only after done closes.
	address string
}

func childEnvironment(overrides map[string]string) []string {
	var result []string
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if !strings.HasPrefix(strings.ToUpper(key), "FL_") {
			result = append(result, entry)
		}
	}
	for key, value := range overrides {
		result = append(result, key+"="+value)
	}
	return result
}

func freeAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	return address
}

func start(t *testing.T, binary string, overrides map[string]string) *process {
	t.Helper()
	p := &process{cmd: exec.Command(binary), logs: &logCapture{started: make(chan string, 1)}, done: make(chan struct{})}
	values := map[string]string{"FL_GRPC_ADDR": freeAddress(t)}
	for key, value := range overrides {
		values[key] = value
	}
	p.cmd.Env = childEnvironment(values)
	p.cmd.Stdout, p.cmd.Stderr = p.logs, p.logs
	if err := p.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	go func() { p.waitErr = p.cmd.Wait(); close(p.done) }()
	t.Cleanup(func() {
		select {
		case <-p.done:
			return
		default:
		}
		if err := p.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			t.Error(err)
		}
		select {
		case <-p.done:
		case <-time.After(10 * time.Second):
			t.Error("child did not exit after cleanup")
		}
	})
	select {
	case p.address = <-p.logs.started:
	case <-p.done:
		t.Fatalf("service exited before startup: %v; logs: %v", p.waitErr, p.logs.check(t, filepath.Base(binary)))
	case <-time.After(10 * time.Second):
		t.Fatal("service did not announce startup")
	}
	return p
}

func checkHealth(t *testing.T, address, service string) {
	t.Helper()
	conn, err := grpc.NewClient("passthrough:///"+address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{Service: service}, grpc.WaitForReady(true))
	if service == "unknown-business-api" {
		if status.Code(err) != codes.NotFound {
			t.Fatalf("unknown health service: got %v, want NotFound", err)
		}
		return
	}
	if err != nil || response.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("health %s = %v, %v", address, response, err)
	}
}

func stop(t *testing.T, p *process, service string, signal os.Signal) {
	t.Helper()
	// Windows cannot send os.Interrupt to an arbitrary child. Its forced cleanup
	// verifies isolation; normal signal exit/logs are verified on both Linux CIs.
	var err error
	if runtime.GOOS == "windows" {
		err = p.cmd.Process.Kill()
	} else {
		err = p.cmd.Process.Signal(signal)
	}
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-p.done:
	case <-time.After(10 * time.Second):
		t.Fatal("service did not stop")
	}
	records := p.logs.check(t, service)
	if runtime.GOOS != "windows" {
		if p.waitErr != nil || len(records) != 3 || records[1]["msg"] != "service stopping" || records[1]["health"] != "NOT_SERVING" || records[2]["msg"] != "service stopped" || records[2]["stop_mode"] != "graceful" {
			t.Fatalf("unexpected normal exit: %v; logs: %v", p.waitErr, records)
		}
		logJSON(t, records[1])
		logJSON(t, records[2])
	}
}

func TestServiceProcesses(t *testing.T) {
	binaries := os.Getenv("LAB_BIN_DIR")
	if binaries == "" {
		binaries = t.TempDir()
		build := exec.CommandContext(t.Context(), "go", "build", "-o", binaries+string(os.PathSeparator), "./cmd/...")
		build.Dir = "../.."
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build service binaries: %v\n%s", err, output)
		}
	}
	services := []string{"orders", "inventory", "payments"}
	processes := make([]*process, len(services))
	for i, name := range services {
		binary := filepath.Join(binaries, name)
		if runtime.GOOS == "windows" {
			binary += ".exe"
		}
		processes[i] = start(t, binary, nil)
		checkHealth(t, processes[i].address, "")
		checkHealth(t, processes[i].address, "unknown-business-api")
		records := processes[i].logs.check(t, name)
		for _, key := range []string{"version", "revision", "vcs_modified", "go_version"} {
			if records[0][key] == nil || records[0][key] == "" {
				t.Errorf("%s startup lacks %s", name, key)
			}
		}
		t.Logf("%s health SERVING", name)
		logJSON(t, records[0])
		t.Run(name+"/healthcheck-command", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			probe := exec.CommandContext(ctx, binary, "healthcheck")
			probe.Env = childEnvironment(map[string]string{"FL_HEALTH_TARGET": processes[i].address})
			if output, err := probe.CombinedOutput(); err != nil || len(output) != 0 {
				t.Fatalf("probe = %v; output %s", err, output)
			}
		})
		for key, value := range map[string]string{"FL_GRPC_ADDR": "secret://user:password@host:80", "FL_LOG_LEVEL": "super-secret-value", "FL_SHUTDOWN_TIMEOUT": "0s"} {
			t.Run(name+"/reject/"+key, func(t *testing.T) {
				values := map[string]string{"FL_GRPC_ADDR": processes[i].address}
				values[key] = value
				assertFailure(t, binary, name, values, key, value)
			})
		}
		t.Run(name+"/missing-address", func(t *testing.T) {
			assertFailure(t, binary, name, nil, "FL_GRPC_ADDR is required", "")
		})
		t.Run(name+"/empty-address", func(t *testing.T) {
			assertFailure(t, binary, name, map[string]string{"FL_GRPC_ADDR": ""}, "FL_GRPC_ADDR is required", "")
		})
		t.Run(name+"/occupied", func(t *testing.T) {
			assertFailure(t, binary, name, map[string]string{"FL_GRPC_ADDR": processes[i].address}, "listen for gRPC", "")
		})
		if runtime.GOOS != "windows" {
			t.Run(name+"/shutdown-deadline", func(t *testing.T) {
				assertProcessDeadline(t, binary, name)
			})
		}
	}
	stop(t, processes[0], "orders", syscall.SIGTERM)
	for i := 1; i < len(services); i++ {
		checkHealth(t, processes[i].address, "")
		t.Logf("%s remains SERVING after orders stopped", services[i])
	}
	for i := 1; i < len(services); i++ {
		stop(t, processes[i], services[i], os.Interrupt)
	}
}

func assertProcessDeadline(t *testing.T, binary, service string) {
	t.Helper()
	p := start(t, binary, map[string]string{"FL_SHUTDOWN_TIMEOUT": "500ms"})
	conn, err := grpc.NewClient("passthrough:///"+p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	watch, err := healthpb.NewHealthClient(conn).Watch(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
	if err != nil {
		t.Fatal(err)
	}
	response, err := watch.Recv()
	if err != nil || response.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("initial Watch = %v, %v", response, err)
	}
	if err := p.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	// Keep the stream open after NOT_SERVING to exercise the real process's
	// deadline, error logging and exit code without inventing a business RPC.
	response, err = watch.Recv()
	if err != nil || response.GetStatus() != healthpb.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("shutdown Watch = %v, %v", response, err)
	}
	select {
	case <-p.done:
	case <-ctx.Done():
		t.Fatal("process did not enforce its shutdown deadline")
	}
	var exitError *exec.ExitError
	if !errors.As(p.waitErr, &exitError) || exitError.ExitCode() != 1 {
		t.Fatalf("forced process exit = %v, want 1", p.waitErr)
	}
	records := p.logs.check(t, service)
	if len(records) != 4 || records[2]["stop_mode"] != "forced" || records[3]["msg"] != "service failed" {
		t.Fatalf("missing forced shutdown diagnostics: %v", records)
	}
	t.Logf("%s: open Watch forced shutdown and exit 1 verified", service)
}

func assertFailure(t *testing.T, binary, service string, overrides map[string]string, reason, secret string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Env = childEnvironment(overrides)
	output, err := cmd.CombinedOutput()
	var exitError *exec.ExitError
	if ctx.Err() != nil || !errors.As(err, &exitError) || exitError.ExitCode() != 1 {
		t.Fatalf("expected prompt exit 1, got %v; output %s", err, output)
	}
	var record map[string]any
	if err := json.Unmarshal(output, &record); err != nil {
		t.Fatalf("failure was not one JSON record: %v; %s", err, output)
	}
	if record["service"] != service || record["level"] != "ERROR" || !strings.Contains(fmt.Sprint(record["error"]), reason) || (secret != "" && strings.Contains(string(output), secret)) {
		t.Fatalf("unexpected failure record: %v", record)
	}
	logJSON(t, record)
}

func logJSON(t *testing.T, record map[string]any) {
	t.Helper()
	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("observed JSON: %s", data)
}
```

<a id="appendix-local-env-example"></a>

## Приложение: Явное локальное окружение: образец

Полный исходный документ: `configs/env/local.env.example`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```dotenv
# Example for one process on the development host; inject via the launcher.
# Use distinct ports for simultaneous processes (e.g. 50051, 50052, 50053).
# Copy to an ignored *.env file when customizing; never add secrets to examples.
FL_GRPC_ADDR=localhost:50051
FL_HEALTH_TARGET=localhost:50051
FL_LOG_LEVEL=INFO
FL_SHUTDOWN_TIMEOUT=10s
```

<a id="appendix-container-env-example"></a>

## Приложение: Явное контейнерное окружение: образец

Полный исходный документ: `configs/env/container.env.example`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```dotenv
# Explicit wildcard listener inside the container's network namespace.
# Each service can use the same internal port in its own container.
# Host port publication and peer DNS endpoints belong to deployment config.
FL_GRPC_ADDR=:50051
FL_HEALTH_TARGET=localhost:50051
FL_LOG_LEVEL=INFO
FL_SHUTDOWN_TIMEOUT=10s
```

<a id="appendix-gitignore"></a>

## Приложение: Исключение личных настроек и артефактов из Git

Полный исходный документ: `.gitignore`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```text
# Local configuration and credentials
.env
.env.*
*.env
!.env.example
*.pem
*.key
credentials.json

# Go build and test output
/bin/
/dist/
/coverage/
*.test
coverage.out
coverage.html

# Local lab data; intentional reports belong in docs/experiments/
/.local/
/data/
/artifacts/
/tmp/

# Editor and operating system files
.idea/
.vscode/
.DS_Store
Thumbs.db
Desktop.ini
*.swp
*.swo

# Dependency caches
node_modules/
__pycache__/
.venv/
```

<a id="appendix-service-command"></a>

## Приложение: Entrypoint command and bounded health probe

Полный исходный документ: `internal/service/command.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Execute runs the service or a bounded local health probe. It returns the exit
// code; only cmd entrypoints call os.Exit, including after a forced shutdown.
func Execute(ctx context.Context, name string, args []string, lookup func(string) (string, bool), output io.Writer) int {
	if len(args) == 0 {
		if err := Run(ctx, name, lookup, output); err != nil {
			return 1
		}
		return 0
	}
	logger := newLogger(output, name, slog.LevelInfo)
	if len(args) != 1 || args[0] != "healthcheck" {
		logger.ErrorContext(ctx, "invalid command", "usage", "service [healthcheck]")
		return 2
	}
	if err := probe(ctx, lookup); err != nil {
		logger.ErrorContext(ctx, "health check failed", "error", err.Error())
		return 1
	}
	return 0
}

func probe(ctx context.Context, lookup func(string) (string, bool)) error {
	target, _ := lookup("FL_HEALTH_TARGET")
	host, _, err := net.SplitHostPort(target)
	if err != nil || !validAddress(target) || host == "" || net.ParseIP(host).IsUnspecified() {
		return errors.New("FL_HEALTH_TARGET is required and must contain a concrete host and port")
	}
	// This command only probes the local plaintext lab. TLS for a deployed
	// business system will be an explicit, separately reviewed configuration.
	conn, err := grpc.NewClient("passthrough:///"+target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New("create health check client")
	}
	defer func() { _ = conn.Close() }() // Sole owner; a second Close is the only error case.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New("gRPC health request failed or exceeded its 2s deadline")
	}
	if response.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		return errors.New("gRPC health is not SERVING")
	}
	return nil
}
```

<a id="appendix-service-command-tests"></a>

## Приложение: Health probe and command tests

Полный исходный документ: `internal/service/command_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package service

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthcheckCommand(t *testing.T) {
	t.Parallel()
	fixture := startShutdownFixture(t, time.Second)
	lookup := environment(map[string]string{"FL_HEALTH_TARGET": fixture.listener.Addr().String()})
	var output bytes.Buffer
	if code := Execute(t.Context(), "fixture", []string{"healthcheck"}, lookup, &output); code != 0 || output.Len() != 0 {
		t.Fatalf("healthy probe = %d; output %s", code, &output)
	}
	fixture.health.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	if code := Execute(t.Context(), "fixture", []string{"healthcheck"}, lookup, &output); code != 1 || !strings.Contains(output.String(), "not SERVING") {
		t.Fatalf("unhealthy probe = %d; output %s", code, &output)
	}
}

func TestHealthcheckRejectsInvalidConfiguration(t *testing.T) {
	t.Parallel()
	for _, target := range []string{"", ":50051", "0.0.0.0:50051", "[::]:50051", "secret://user:password@host:80"} {
		t.Run(target, func(t *testing.T) {
			var output bytes.Buffer
			code := Execute(t.Context(), "orders", []string{"healthcheck"}, environment(map[string]string{"FL_HEALTH_TARGET": target}), &output)
			if code != 1 || !strings.Contains(output.String(), "FL_HEALTH_TARGET is required") || (target != "" && strings.Contains(output.String(), target)) {
				t.Fatalf("probe = %d; output %s", code, &output)
			}
		})
	}
	var output bytes.Buffer
	if code := Execute(t.Context(), "orders", []string{"secret-argument"}, environment(nil), &output); code != 2 || strings.Contains(output.String(), "secret-argument") {
		t.Fatalf("unknown command = %d; output %s", code, &output)
	}
}

func TestHealthcheckHonorsCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	var output bytes.Buffer
	lookup := environment(map[string]string{"FL_HEALTH_TARGET": "localhost:1"})
	if code := Execute(ctx, "orders", []string{"healthcheck"}, lookup, &output); code != 1 || !strings.Contains(output.String(), "health check failed") {
		t.Fatalf("cancelled probe = %d; output %s", code, &output)
	}
}
```

<a id="appendix-shutdown-tests"></a>

## Приложение: Controlled active RPC and deadline tests

Полный исходный документ: `internal/service/shutdown_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package service

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	testpb "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/status"
)

type closingListener struct {
	net.Listener
	closed chan struct{}
	once   sync.Once
}

func (listener *closingListener) Close() error {
	err := listener.Listener.Close()
	listener.once.Do(func() { close(listener.closed) })
	return err
}

type controlledRPC struct {
	testpb.UnimplementedTestServiceServer
	entered chan struct{}
	release chan struct{}
	exited  chan struct{}
}

func (rpc *controlledRPC) EmptyCall(context.Context, *testpb.Empty) (*testpb.Empty, error) {
	close(rpc.entered)
	defer close(rpc.exited)
	// Deliberately ignores cancellation to prove that shutdown returns on its
	// deadline. The test owns release and always lets this handler exit afterward.
	<-rpc.release
	return &testpb.Empty{}, nil
}

type shutdownFixture struct {
	listener      *closingListener
	server        *grpc.Server
	health        *health.Server
	work          *controlledRPC
	conn          *grpc.ClientConn
	cancel        context.CancelFunc
	stopping      chan struct{}
	continueDrain func()
	releaseWork   func()
	done          chan struct{}
	runErr        error // Read after done closes.
	logs          bytes.Buffer
}

func startShutdownFixture(t *testing.T, timeout time.Duration) *shutdownFixture {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	gate := make(chan struct{})
	fixture := &shutdownFixture{
		listener: &closingListener{Listener: listener, closed: make(chan struct{})},
		server:   grpc.NewServer(), health: health.NewServer(), cancel: cancel,
		stopping: make(chan struct{}), done: make(chan struct{}),
		work: &controlledRPC{entered: make(chan struct{}), release: make(chan struct{}), exited: make(chan struct{})},
	}
	fixture.continueDrain = sync.OnceFunc(func() { close(gate) })
	fixture.releaseWork = sync.OnceFunc(func() { close(fixture.work.release) })
	healthpb.RegisterHealthServer(fixture.server, fixture.health)
	testpb.RegisterTestServiceServer(fixture.server, fixture.work)
	writer := writerFunc(func(data []byte) (int, error) {
		length, err := fixture.logs.Write(data)
		if strings.Contains(string(data), `"msg":"service stopping"`) {
			close(fixture.stopping)
			// Observe real health over TCP before GracefulStop closes the listener.
			<-gate
		}
		return length, err
	})
	go func() {
		fixture.runErr = serveGRPC(ctx, fixture.listener, fixture.server, fixture.health, timeout, newLogger(writer, "fixture", slog.LevelInfo))
		close(fixture.done)
	}()
	t.Cleanup(func() {
		cancel()
		fixture.continueDrain()
		fixture.releaseWork()
		await(t, fixture.done, "server cleanup")
		fixture.server.Stop()
		if fixture.conn != nil {
			if err := fixture.conn.Close(); err != nil {
				t.Error(err)
			}
		}
	})
	fixture.conn, err = grpc.NewClient("passthrough:///"+listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	assertHealth(t, fixture.conn, healthpb.HealthCheckResponse_SERVING)
	return fixture
}

func await(t *testing.T, signal <-chan struct{}, operation string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", operation)
	}
}

func assertHealth(t *testing.T, conn *grpc.ClientConn, want healthpb.HealthCheckResponse_ServingStatus) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
	if err != nil || response.GetStatus() != want {
		t.Fatalf("health = %v, %v; want %v", response, err, want)
	}
}

func beginWork(t *testing.T, fixture *shutdownFixture) <-chan error {
	t.Helper()
	result := make(chan error, 1)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)
	go func() {
		_, err := testpb.NewTestServiceClient(fixture.conn).EmptyCall(ctx, &testpb.Empty{})
		result <- err
	}()
	await(t, fixture.work.entered, "active RPC")
	return result
}

func beginDrain(t *testing.T, fixture *shutdownFixture) {
	t.Helper()
	fixture.cancel()
	await(t, fixture.stopping, "readiness removal")
	assertHealth(t, fixture.conn, healthpb.HealthCheckResponse_NOT_SERVING)
	fixture.continueDrain()
	await(t, fixture.listener.closed, "listener closure")
	connection, err := net.DialTimeout("tcp", fixture.listener.Addr().String(), time.Second)
	if err == nil {
		_ = connection.Close()
		t.Fatal("listener accepted a new connection during drain")
	}
}

func TestGracefulShutdownFinishesActiveRPC(t *testing.T) {
	t.Parallel()
	fixture := startShutdownFixture(t, 3*time.Second)
	result := beginWork(t, fixture)
	beginDrain(t, fixture)
	select {
	case err := <-result:
		t.Fatalf("active RPC ended before work was released: %v", err)
	default:
	}
	fixture.releaseWork()
	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("active RPC was interrupted: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("released RPC did not complete")
	}
	await(t, fixture.done, "graceful shutdown")
	if fixture.runErr != nil || !strings.Contains(fixture.logs.String(), `"stop_mode":"graceful"`) {
		t.Fatalf("shutdown = %v; logs: %s", fixture.runErr, &fixture.logs)
	}
}

func TestShutdownDeadlineInterruptsBlockedRPC(t *testing.T) {
	t.Parallel()
	const budget = 200 * time.Millisecond
	fixture := startShutdownFixture(t, budget)
	result := beginWork(t, fixture)
	started := time.Now()
	beginDrain(t, fixture)
	await(t, fixture.done, "bounded forced shutdown")
	if !errors.Is(fixture.runErr, ErrShutdownTimeout) || time.Since(started) < budget || !strings.Contains(fixture.logs.String(), `"stop_mode":"forced"`) {
		t.Fatalf("shutdown = %v after %v; logs: %s", fixture.runErr, time.Since(started), &fixture.logs)
	}
	select {
	case err := <-result:
		if status.Code(err) != codes.Unavailable {
			t.Fatalf("forced RPC = %v, want Unavailable", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("client did not observe forced interruption")
	}
	// Release the intentionally noncooperative handler so the test also joins
	// gRPC internals. Real cmd entrypoints exit with code 1 on this fatal path.
	fixture.releaseWork()
	await(t, fixture.work.exited, "blocked handler cleanup")
}
```

<a id="appendix-container-runtime"></a>

## Приложение: Runtime build, health, shutdown, smoke and review limitations

Полный исходный документ: `docs/development/CONTAINER_RUNTIME.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Контейнерный каркас P01

Этот пакет завершает реализацию P01-05–08: ограниченное graceful shutdown, отдельные runtime-образы трёх сервисов, Compose, CI и smoke с диагностикой. G0 принят; P01 выполняется до G1. Бизнес-операции и инфраструктура PostgreSQL/Kafka/Redis относятся к следующим пакетам.

Исторический источник каркаса — [PR #8](https://github.com/Includeoyasi/faultline-lab/pull/8), принятый владельцем 16 октября 2026 года. Исходные Windows-проверки сохранены в [evidence](#windows-smoke-proof); актуальное ревью — в GitHub.

#### Состав и границы

```mermaid
flowchart TB
    Host["Windows / macOS / Linux<br/>Docker CLI + Compose"] --> Compose["Compose: три runtime-контейнера"]
    Base["Закреплённый Go toolchain<br/>только сборка и проверки"] --> Build["Multi-stage build<br/>статические Go-бинарники"]
    Build --> O["Orders<br/>scratch / UID 65532"]
    Build --> I["Inventory<br/>scratch / UID 65532"]
    Build --> P["Payments<br/>scratch / UID 65532"]
    Compose --> O
    Compose --> I
    Compose --> P
    O --- Net["Внутренняя сеть runtime<br/>без публикации host ports"]
    I --- Net
    P --- Net
```

[Dockerfile](#appendix-toolchain-dockerfile) использует закреплённый [toolchain](#appendix-toolchain) как stage общего multi-stage Dockerfile. Сервис выбирается build argument, разрешены только orders/inventory/payments. Финальный `scratch` содержит один статический бинарник и запускает его напрямую: без компилятора, shell, root и промежуточного процесса вместо PID 1. CA bundle и timezone database сейчас не требуются; перед появлением TLS-клиентов состав образа пересматривается явно.

[Compose](#appendix-compose-toolchain) задаёт read-only root filesystem, `cap_drop: ALL`, `no-new-privileges`, внутреннюю сеть и отсутствие публикации портов. Runtime не получает host mounts или Docker socket. Toolchain остаётся временным сервисом профиля tools. Сборка runtime использует stage toolchain напрямую; активация сервиса профиля tools для неё не требуется. Обычный up запускает только три runtime-сервиса.

#### Требования и запуск

Нужны Git, Docker с Linux-контейнерами и Compose v2 с поддержкой up --wait. Для локальной публикации на loopback используйте Docker Engine 28+. Go и Bash на хосте не нужны. Готовый Go smoke не требует Node.js или Go SDK на хосте; [сборка и запуск](#go-smoke-guide).

Из корня репозитория:

```text
docker compose build orders inventory payments
.\bin\smoke-windows-amd64.exe -check-config
docker compose up -d --wait orders inventory payments
docker compose ps
docker compose exec -T orders /app/service healthcheck
docker compose exec -T inventory /app/service healthcheck
docker compose exec -T payments /app/service healthcheck
docker compose stop orders inventory payments
docker compose down
```

`up --wait` ожидает healthy. Команда `healthcheck` внутри того же бинарника вызывает стандартный gRPC `Health/Check`, молчит при успехе и возвращает 0. Ошибка/NOT_SERVING — безопасная JSON-диагностика и exit 1; неизвестная команда — exit 2. Проба ограничена двумя секундами, Compose даёт ей три секунды. Проба не подтверждает готовность базы или бизнес-процесса.

Окружение задаётся через [container.env.example](#appendix-container-env-example). Все три контейнера могут использовать один внутренний порт; адресов/портов по умолчанию в Go-коде нет. `FL_GRPC_ADDR` обязателен для сервера; отдельный `FL_HEALTH_TARGET` обязателен только для команды healthcheck. Wildcard допустим как явно выбранный listener, но отклоняется как цель подключения.

При изменении порта listener обновляйте и target собственной health-пробы. Для нескольких локальных процессов target каждого должен указывать на его listener; в Compose каждый сервис имеет отдельное сетевое пространство.

Для личной конфигурации скопируйте образец в игнорируемый `configs/env/container.env` и задайте в корневом игнорируемом `.env`:

```dotenv
FL_RUNTIME_ENV_FILE=configs/env/container.env
```

Compose выбирает файл и передаёт окружение процессам; приложение dotenv не читает. Личные env исключены из Git и Docker build context. Не печатайте полный конфиг при работе с секретами; для валидации используйте `docker compose config --quiet`.

Порты на хост не публикуются. Если позже понадобится локальный доступ, добавляется отдельная deployment-настройка с явной привязкой к loopback. Endpoint клиента, listener и host port — разные параметры; детали в [описании сервисов](#appendix-service-runtime). Текущий приватный стенд не заявляет готовность production TLS/mTLS.

#### Остановка и её проверка

`SIGINT`/`SIGTERM` переводит health в NOT_SERVING; затем `GracefulStop` закрывает listener и ждёт активные RPC. При успехе процесс пишет `stop_mode=graceful` и выходит с кодом 0. Если `FL_SHUTDOWN_TIMEOUT` истёк, запрашивается принудительный Stop, пишется `stop_mode=forced` и ошибка, entrypoint завершает процесс с кодом 1.

На принудительном пути приложение не ждёт бесконечно goroutine GracefulStop: gRPC может удерживать её на handler, который игнорирует cancellation. Владельцем оставшейся работы становится завершающийся процесс. Все будущие бизнес-handlers обязаны соблюдать context. Unit fixture намеренно нарушает это правило, проверяет ограниченный возврат и затем сам освобождает handler, чтобы не оставлять goroutine в тестах.

Управляемые unit-тесты проверяют завершение активного RPC, deadline, NOT_SERVING через реальный TCP endpoint до закрытия listener и отказ новых соединений. Синхронизация — каналы и ограниченные ожидания. Linux process tests проверяют SIGTERM/SIGINT для настоящих бинарников и forced exit 1 при открытом health Watch. Постоянный Watch может занять весь drain budget; Compose использует короткий unary Check.

Default приложения — 10s; независимый `FL_STOP_GRACE_PERIOD` Compose — 15s. [Preflight Go CLI](#go-smoke-guide) до запуска проверяет внешний бюджет с минимум 2s запаса; несовместимый профиль отклоняется. Таймауты проверяют свойства завершения; это не результат исследования производительности.

#### Автоматический smoke и отчёт

[Go CLI](#go-smoke-guide) выполняет build → start → health/security → isolation → stop/logs → cleanup на собственном уникальном Compose project. Руководство описывает контейнерную сборку для пяти платформ, отмену дерева процессов, независимый cleanup, формат отчёта и адресное восстановление после аварии.

CI выполняет Go quality checks и этот host-side smoke на нативных Linux AMD64/ARM64, сохраняя JSON/JSONL и диагностику на 14 дней. Сборка macOS не заменяет фактическую проверку на Mac.

#### Локальная проверка Windows

15 октября установлены host Go 1.27.2 (официальный MSI), WSL 2.7.13.0 и Docker Desktop 4.95.0 в профиль пользователя. Linux engine 29.8.2 / Compose 5.6.0 успешно выполнили smoke на чистом source d81b51c. [Отчёт, JSONL и подтверждение остановки](#windows-smoke-proof) включены в репозиторий и переносимый контекст. Перезагрузка не понадобилась.

#### Личный компьютер: запуск по необходимости

WSL ограничен двумя CPU, memory=2GB и swap=1GB; Engine видит 2 CPU и 1996603392 bytes RAM. Это бюджет smoke, не измерение capacity бизнеса. AutoStart Docker выключен и запись из автозагрузки Windows удалена.

После установки откройте новый терминал, чтобы он получил обновлённый PATH. Запуск и остановка:

```text
docker desktop start --timeout 120
.\bin\smoke-windows-amd64.exe
docker desktop stop --timeout 60
wsl --shutdown
```

Последняя команда относится к Windows и завершает виртуальную среду WSL. В этом профиле WSL использует только лаборатория; если появятся другие дистрибутивы с работающими задачами, учитывайте их перед общим shutdown. Smoke удаляет только собственные контейнеры и сеть. Docker Desktop и WSL VM после работы нужно остановить также при ошибке проверки, сохранив диагностику. Программы, образы, библиотеки и кеши оставляем для следующих запусков.

После этого прогона подтверждены: контейнеров нет, Docker Desktop остановлен, wsl --shutdown выполнен, процессов Docker/vmmemWSL нет, автозапуск выключен. Доказательства относятся к этой дате и исходному commit; новые проверки и приёмка изменений фиксируются отдельно.

<a id="windows-smoke-proof"></a>

## Приложение: Windows smoke: итог и ограничения

Полный исходный документ: `docs/development/evidence/p01/WINDOWS_SMOKE.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### P01: проверка контейнеров на Windows

15 октября 2026 года выполнен обязательный Windows container smoke P01-06. Это запуск Linux-контейнеров через Docker Desktop/WSL2 на личном Windows-компьютере, а не проверка Windows container images или macOS.

Исходный commit: `d81b51ce2eb46f56d698fd8bf5f9966b91a16036`, рабочее дерево чистое. Команда: `node tools/smoke.mjs`. Run ID: `8119b50d-f910-412a-9bd2-ab50fccd1b71`. Результат: **PASS**.

| Компонент | Проверенная версия / значение |
| :--- | :--- |
| Хост | Windows 11 Pro, build 26200, x64 |
| Go на хосте | 1.27.2 windows/amd64; официальный MSI, SHA256 и подпись Google LLC проверены |
| WSL | 2.7.13.0 |
| Docker Desktop | 4.95.0, установка в профиль пользователя |
| Docker Engine | 29.8.2, linux/x86_64 |
| Compose | 5.6.0 |
| Node smoke | v24.19.0 |
| Бюджет | 2 CPU, Engine MemTotal 1996603392 bytes |
| WSL config | memory=2GB, processors=2, swap=1GB |

| Сервис | Health | Ограничения | Остановка |
| :--- | :--- | :--- | :--- |
| orders | healthy; unary Check PASS | UID 65532:65532, read-only, capabilities dropped, no-new-privileges, internal network, no host ports/mounts | exit 0, graceful, OOM=false |
| inventory | healthy; unary Check PASS | UID 65532:65532, read-only, capabilities dropped, no-new-privileges, internal network, no host ports/mounts | exit 0, graceful, OOM=false |
| payments | healthy; unary Check PASS | UID 65532:65532, read-only, capabilities dropped, no-new-privileges, internal network, no host ports/mounts | exit 0, graceful, OOM=false |

После остановки Orders оставшиеся сервисы продолжили отвечать: **PASS**. Удаление собственного Compose-проекта: **PASS**.

#### Исходные данные

- [report.json](#windows-smoke-report): команды, время, исходы, source revision и реальные image IDs.
- Логи: [Orders](#windows-smoke-orders), [Inventory](#windows-smoke-inventory), [Payments](#windows-smoke-payments).
- [environment.json](#windows-smoke-environment): наблюдённый backend и бюджет.
- [cleanup.json](#windows-smoke-cleanup): подтверждение остановки среды.
- SHA256 report.json: `b08d9deabadc65c765d9173482e9171782aeff28affb5e7fd9bf0b201d5add65`.

#### Состояние личного компьютера после проверки

Тестовые контейнеры и сеть удалены. Docker Desktop штатно остановлен, выполнен wsl --shutdown; процессов Docker и vmmemWSL при финальной проверке нет. AutoStart=false, запись Docker Desktop из автозагрузки Windows удалена. Go, Docker, WSL, образы и build cache сохранены для дальнейшей работы. Перезагрузка не потребовалась и не выполнялась.

При следующих проверках запускайте среду по необходимости и освобождайте ресурсы после окончания. Инструкция: [контейнерный каркас](#appendix-container-runtime).

Этот smoke подтверждает lifecycle каркаса. Бизнес-операции, базы, пропускная способность и пиковое потребление под нагрузкой не измерялись. Личный Mac остаётся NOT_RUN. Результат не означает принятия P01 или G1 владельцем.

<a id="windows-smoke-report"></a>

## Приложение: Полный Windows smoke report

Полный исходный документ: `docs/development/evidence/p01/windows-2026-10-15/report.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "run_id": "8119b50d-f910-412a-9bd2-ab50fccd1b71",
  "started_at": "2026-10-15T19:45:57.496Z",
  "result": "PASS",
  "host_platform": "win32",
  "node": "v24.19.0",
  "project": "faultline-smoke-8119b50d",
  "services": {
    "orders": {
      "id": "0e35127ad00aa13be4f551edf6b032c00ffc4894081efb4e4c33cc723fff9b04",
      "image_id": "sha256:c15eb12d4ab1369a2c49fd8da5f6ca178c530d39cf0f19c37d4ad513eac4574f",
      "architecture": "linux/x86_64",
      "healthy": true,
      "security": {
        "user": "65532:65532",
        "read_only": true,
        "capabilities_dropped": true,
        "no_new_privileges": true,
        "host_ports": [],
        "mounts": [],
        "internal_network": true
      },
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "stop_mode": "graceful"
    },
    "inventory": {
      "id": "926b6baf437c4b0d0a01302d158872112689ddf89472c58a56f0d66c8f621d84",
      "image_id": "sha256:1fb289e7407f8559f671e9ec276faac6715f7c6e364220f70b30839624c7f778",
      "architecture": "linux/x86_64",
      "healthy": true,
      "security": {
        "user": "65532:65532",
        "read_only": true,
        "capabilities_dropped": true,
        "no_new_privileges": true,
        "host_ports": [],
        "mounts": [],
        "internal_network": true
      },
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "stop_mode": "graceful"
    },
    "payments": {
      "id": "8a8c51c0112100bbdf7b981f9ebc1f9c4e099baaa24654a85eb85e956f55b868",
      "image_id": "sha256:3481c4f6aea54924e0f16b1da34fb2082bef2e937f5477fe077f5bce8b4eec35",
      "architecture": "linux/x86_64",
      "healthy": true,
      "security": {
        "user": "65532:65532",
        "read_only": true,
        "capabilities_dropped": true,
        "no_new_privileges": true,
        "host_ports": [],
        "mounts": [],
        "internal_network": true
      },
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "stop_mode": "graceful"
    }
  },
  "commands": [
    {
      "program": "docker",
      "args": [
        "version",
        "--format",
        "{{.Server.Version}}"
      ],
      "started_at": "2026-10-15T19:45:57.501Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:45:57.858Z"
    },
    {
      "program": "docker",
      "args": [
        "info",
        "--format",
        "{{.OSType}}/{{.Architecture}}"
      ],
      "started_at": "2026-10-15T19:45:57.858Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:45:59.547Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "version",
        "--short"
      ],
      "started_at": "2026-10-15T19:45:59.547Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:46:00.081Z"
    },
    {
      "program": "git",
      "args": [
        "rev-parse",
        "HEAD"
      ],
      "started_at": "2026-10-15T19:46:00.081Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:46:00.163Z"
    },
    {
      "program": "git",
      "args": [
        "status",
        "--porcelain"
      ],
      "started_at": "2026-10-15T19:46:00.163Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:46:00.244Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "config",
        "--quiet"
      ],
      "started_at": "2026-10-15T19:46:00.245Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:46:00.877Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "build",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-15T19:46:00.877Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:16.700Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "up",
        "-d",
        "--wait",
        "--wait-timeout",
        "60",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-15T19:48:16.701Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:23.023Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "exec",
        "-T",
        "orders",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-15T19:48:23.024Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:23.550Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-15T19:48:23.550Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:24.049Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "0e35127ad00aa13be4f551edf6b032c00ffc4894081efb4e4c33cc723fff9b04"
      ],
      "started_at": "2026-10-15T19:48:24.049Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:24.274Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b526971f08818f9bea79d032f95822db2be12065b6e61d0bb72556e38d5fc2a6"
      ],
      "started_at": "2026-10-15T19:48:24.278Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:24.501Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-15T19:48:24.501Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:25.030Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-15T19:48:25.030Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:25.523Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "926b6baf437c4b0d0a01302d158872112689ddf89472c58a56f0d66c8f621d84"
      ],
      "started_at": "2026-10-15T19:48:25.523Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:25.745Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b526971f08818f9bea79d032f95822db2be12065b6e61d0bb72556e38d5fc2a6"
      ],
      "started_at": "2026-10-15T19:48:25.745Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:25.964Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-15T19:48:25.964Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:26.622Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-15T19:48:26.622Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:27.459Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "8a8c51c0112100bbdf7b981f9ebc1f9c4e099baaa24654a85eb85e956f55b868"
      ],
      "started_at": "2026-10-15T19:48:27.459Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:27.893Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b526971f08818f9bea79d032f95822db2be12065b6e61d0bb72556e38d5fc2a6"
      ],
      "started_at": "2026-10-15T19:48:27.893Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:28.184Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "stop",
        "orders"
      ],
      "started_at": "2026-10-15T19:48:28.185Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:29.105Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-15T19:48:29.106Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:29.676Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "0e35127ad00aa13be4f551edf6b032c00ffc4894081efb4e4c33cc723fff9b04"
      ],
      "started_at": "2026-10-15T19:48:29.676Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:29.836Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-15T19:48:29.836Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:30.294Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "0e35127ad00aa13be4f551edf6b032c00ffc4894081efb4e4c33cc723fff9b04"
      ],
      "started_at": "2026-10-15T19:48:30.294Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:30.517Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-15T19:48:30.523Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:31.082Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-15T19:48:31.082Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:31.590Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "926b6baf437c4b0d0a01302d158872112689ddf89472c58a56f0d66c8f621d84"
      ],
      "started_at": "2026-10-15T19:48:31.590Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:31.821Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b526971f08818f9bea79d032f95822db2be12065b6e61d0bb72556e38d5fc2a6"
      ],
      "started_at": "2026-10-15T19:48:31.822Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:32.080Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-15T19:48:32.080Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:32.626Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-15T19:48:32.626Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:33.148Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "8a8c51c0112100bbdf7b981f9ebc1f9c4e099baaa24654a85eb85e956f55b868"
      ],
      "started_at": "2026-10-15T19:48:33.148Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:33.400Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b526971f08818f9bea79d032f95822db2be12065b6e61d0bb72556e38d5fc2a6"
      ],
      "started_at": "2026-10-15T19:48:33.400Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:33.648Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "stop",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-15T19:48:33.648Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:34.375Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-15T19:48:34.375Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:34.879Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "926b6baf437c4b0d0a01302d158872112689ddf89472c58a56f0d66c8f621d84"
      ],
      "started_at": "2026-10-15T19:48:34.879Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:35.106Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-15T19:48:35.106Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:35.618Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "926b6baf437c4b0d0a01302d158872112689ddf89472c58a56f0d66c8f621d84"
      ],
      "started_at": "2026-10-15T19:48:35.618Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:35.845Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-15T19:48:35.847Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:36.370Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "8a8c51c0112100bbdf7b981f9ebc1f9c4e099baaa24654a85eb85e956f55b868"
      ],
      "started_at": "2026-10-15T19:48:36.370Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:36.590Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-15T19:48:36.590Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:37.106Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "8a8c51c0112100bbdf7b981f9ebc1f9c4e099baaa24654a85eb85e956f55b868"
      ],
      "started_at": "2026-10-15T19:48:37.106Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:37.335Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-15T19:48:37.336Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:37.841Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "0e35127ad00aa13be4f551edf6b032c00ffc4894081efb4e4c33cc723fff9b04"
      ],
      "started_at": "2026-10-15T19:48:37.841Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:38.064Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-15T19:48:38.066Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:38.651Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "926b6baf437c4b0d0a01302d158872112689ddf89472c58a56f0d66c8f621d84"
      ],
      "started_at": "2026-10-15T19:48:38.651Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:38.814Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-15T19:48:38.815Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:39.321Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "8a8c51c0112100bbdf7b981f9ebc1f9c4e099baaa24654a85eb85e956f55b868"
      ],
      "started_at": "2026-10-15T19:48:39.321Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:39.633Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8119b50d",
        "down",
        "--remove-orphans"
      ],
      "started_at": "2026-10-15T19:48:39.634Z",
      "exit_code": 0,
      "finished_at": "2026-10-15T19:48:40.359Z"
    }
  ],
  "limitations": [
    "No business operations or capacity measurement"
  ],
  "docker_version": "29.8.2",
  "docker_architecture": "linux/x86_64",
  "compose_version": "5.6.0",
  "source_revision": "d81b51ce2eb46f56d698fd8bf5f9966b91a16036",
  "source_dirty": false,
  "isolation_after_orders_stop": "PASS",
  "cleanup": "PASS",
  "finished_at": "2026-10-15T19:48:40.359Z"
}
```

<a id="windows-smoke-orders"></a>

## Приложение: Windows: Orders JSONL

Полный исходный документ: `docs/development/evidence/p01/windows-2026-10-15/orders.jsonl`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{"time":"2026-10-15T19:48:17.515896056Z","level":"INFO","msg":"service started","service":"orders","version":"(devel)","revision":"unknown","vcs_modified":"unknown","go_version":"go1.27.2","grpc_address":"[::]:50051"}
{"time":"2026-10-15T19:48:28.731075342Z","level":"INFO","msg":"service stopping","service":"orders","health":"NOT_SERVING","shutdown_timeout":"10s"}
{"time":"2026-10-15T19:48:28.731311009Z","level":"INFO","msg":"service stopped","service":"orders","stop_mode":"graceful"}
```

<a id="windows-smoke-inventory"></a>

## Приложение: Windows: Inventory JSONL

Полный исходный документ: `docs/development/evidence/p01/windows-2026-10-15/inventory.jsonl`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{"time":"2026-10-15T19:48:17.642915372Z","level":"INFO","msg":"service started","service":"inventory","version":"(devel)","revision":"unknown","vcs_modified":"unknown","go_version":"go1.27.2","grpc_address":"[::]:50051"}
{"time":"2026-10-15T19:48:33.994730098Z","level":"INFO","msg":"service stopping","service":"inventory","health":"NOT_SERVING","shutdown_timeout":"10s"}
{"time":"2026-10-15T19:48:33.995400269Z","level":"INFO","msg":"service stopped","service":"inventory","stop_mode":"graceful"}
```

<a id="windows-smoke-payments"></a>

## Приложение: Windows: Payments JSONL

Полный исходный документ: `docs/development/evidence/p01/windows-2026-10-15/payments.jsonl`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{"time":"2026-10-15T19:48:17.801803375Z","level":"INFO","msg":"service started","service":"payments","version":"(devel)","revision":"unknown","vcs_modified":"unknown","go_version":"go1.27.2","grpc_address":"[::]:50051"}
{"time":"2026-10-15T19:48:33.994900713Z","level":"INFO","msg":"service stopping","service":"payments","health":"NOT_SERVING","shutdown_timeout":"10s"}
{"time":"2026-10-15T19:48:33.995212132Z","level":"INFO","msg":"service stopped","service":"payments","stop_mode":"graceful"}
```

<a id="windows-smoke-environment"></a>

## Приложение: Windows: фактический бюджет

Полный исходный документ: `docs/development/evidence/p01/windows-2026-10-15/environment.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "desktop": "4.95.0",
  "docker": "29.8.2",
  "remaining_containers_before_stop": [],
  "cpus": 2,
  "memory_bytes": 1996603392,
  "wsl_config": {
    "memory": "2GB",
    "processors": 2,
    "swap": "1GB"
  },
  "platform": "linux/x86_64",
  "observed_at": "2026-10-15T22:49:35.7064530+03:00"
}
```

<a id="windows-smoke-cleanup"></a>

## Приложение: Windows: подтверждение остановки

Полный исходный документ: `docs/development/evidence/p01/windows-2026-10-15/cleanup.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "containers_remaining": 0,
  "wsl_shutdown_exit": 0,
  "images_and_build_cache_preserved": true,
  "desktop_stop_exit": 0,
  "remaining_development_processes": [],
  "installed_tools_preserved": true,
  "docker_autostart": false,
  "docker_run_entry_present": false,
  "observed_at": "2026-10-15T22:49:53.4008625+03:00",
  "reboot_performed": false
}
```

<a id="appendix-engineering-decisions"></a>

## Приложение: Инженерные решения, вклад владельца и доказательства ревью

Полный исходный документ: `docs/development/ENGINEERING_DECISIONS.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Инженерные решения и ревью

Короткая подборка обсуждений: какой вопрос поднял владелец, какое решение подготовил исполнитель и чем подтверждён результат. Подробности остаются в PR, ADR и отчётах. Состояние на 16 октября 2026 года; GitHub хранит актуальные обсуждения и приёмку конкретного HEAD.

#### Примеры из проекта

| Обсуждение | Решение и результат | Статус | Материалы |
| :--- | :--- | :--- | :--- |
| Владелец заметил сетевые IP/порты в прикладном коде | FL_GRPC_ADDR стал обязательным; окружение, listener и публикация портов разделены. Ошибочная конфигурация проверяется до bind | Принято владельцем в PR #7 | [Сводка исправления](https://github.com/Includeoyasi/faultline-lab/pull/7#issuecomment-6192239609), [конфигурация](#appendix-service-config), [описание](#appendix-service-runtime) |
| Владелец поставил вопрос о Node.js в Go-проекте и поручил вынести перенос в backlog | Предложен Go CLI с сохранением smoke-проверок и отчёта; проверка документации остаётся отдельным Node.js-инструментом | Выбор Go согласован владельцем 16 октября; выполнение — в #9 | [Задача #9](https://github.com/Includeoyasi/faultline-lab/issues/9), [ADR 0003](#appendix-adr-0003), [критерии](#appendix-go-smoke-backlog) |
| Владелец потребовал освобождать ресурсы личного ПК после работы | Исполнитель ограничил WSL до 2 CPU / memory=2GB / swap=1GB; после smoke выключил Docker Desktop/WSL VM и автозапуск, сохранив программы и кеши | Принято владельцем в PR #8 | [Windows evidence](#windows-smoke-proof), [cleanup.json](#windows-smoke-cleanup), [PR #8](https://github.com/Includeoyasi/faultline-lab/pull/8) |

Замечания в этих примерах получены в чате с владельцем. Их пересказ и ссылки подготовлены исполнителем AI-Your-Bro. Авторство GitHub-комментариев видно в самих PR; пересказ не выдаётся за комментарий, написанный владельцем на GitHub. Принятие PR #7 подтверждено слиянием владельцем проекта; PR #8 также принят через слияние владельцем проекта 16 октября, Europe/Moscow.

#### Очистка образов после экспериментов

16 октября владелец в чате заметил рост списка Docker images. Причина: уникальный Compose project изолировал каждый smoke-прогон, но `down` удалял только контейнеры и сеть. Исправление сохраняет изоляцию и добавляет `--rmi local`, проверку отсутствия тегов своего прогона и cleanup после частичной сборки. Исторические 21 smoke-образ удалены адресно после проверки отсутствия контейнеров, использующих их; общие PostgreSQL/toolchain и полезные кеши сохранены. Общий `prune` не используется. Методика и результаты — в [отчёте PostgreSQL](#p02-docs-development-evidence-p02-postgres-md); это дополнение поручено владельцем вместе с #12.

#### Как сохраняем следующие обсуждения

```mermaid
flowchart LR
    Q["Вопрос или замечание<br/>PR / чат"] --> V["Варианты и обоснование"]
    V --> C["Изменение и commit"]
    C --> T["Проверка и отчёт"]
    T --> R["Ревью конкретного HEAD"]
    R -->|принято| M["Merge и обновление статуса"]
    R -->|замечания| C
    V -. значимое решение .-> A["ADR: причины и последствия"]
```

- **PR** хранит рабочее обсуждение, изменения и результаты проверок. Для замечания из чата исполнитель указывает источник и автора требования. Комментарий сам по себе не означает approval.
- **ADR** сохраняет существенное решение: контекст, варианты, выбор, последствия, статус и основание согласования. Рутинные исправления остаются в PR. Изменение согласованных границ сначала выносится владельцу; в согласованном объёме исполнитель работает самостоятельно.
- **Эта страница** содержит отобранные примеры и ссылки. Добавляем запись при появлении содержательного решения или изменении его статуса; полный чат сюда не переносится.

Для каждого примера достаточно пяти полей: вопрос → варианты/причина выбора → изменение → доказательства → статус и источник приёмки. Исполнитель сохраняет выбранные решения и датированные доказательства в переносимом контексте.

#### Как читать статусы

**Proposed** — вариант ещё обсуждается. Стадия задачи и её ревью берутся из GitHub. **Принято** — есть датированное решение владельца; для кода также применяются требования review/merge из [workflow](#appendix-development-process). Успешный тест подтверждает только проверенное поведение.

#### Реестр ADR

| Запись | Предмет | Статус решения |
| :--- | :--- | :--- |
| [0001](#appendix-adr-0001) | Архитектура G0 | accepted |
| [0002](#appendix-adr-0002) | PR, ADR и обзорная страница для фиксации решений | accepted: порядок согласован в чате, публикация принята в PR #8 |
| [0003](#appendix-adr-0003) | Go вместо Node.js для контейнерного smoke | accepted: выбран Go; выполнение — в #9 |

<a id="appendix-adr-0002"></a>

## Приложение: ADR 0002: согласованный порядок фиксации решений

Полный исходный документ: `docs/adr/0002-record-review-decisions.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### ADR 0002: фиксация инженерных решений и ревью

- Дата: 2026-10-15, Europe/Moscow.
- Статус решения: **accepted**.
- Решение принял владелец проекта; запись подготовил исполнитель AI-Your-Bro.
- Основание: после предложения использовать PR для обсуждений, ADR для существенных решений и обзорную страницу владелец ответил в чате: «Хорошо».
- Публикация принята: владелец проекта слил PR #8 16 октября 2026 года (Europe/Moscow), merge `3ca290d`. Согласие с порядком работы и фактическая приёмка PR зафиксированы отдельно.

#### Контекст

Проект развивается в сотрудничестве владельца с ИИ. Нужно сохранять причины изменений, замечания владельца, выбранные варианты и реальные результаты, чтобы читатель мог проследить инженерную работу по репозиторию. Часть обсуждений происходит в чате и не имеет собственного GitHub URL.

#### Варианты

1. Оставлять обсуждения только в чате: сведения недоступны читателю репозитория.
2. Складывать полные обсуждения в один журнал: растёт объём, трудно найти решение и проверить статус.
3. Хранить рабочее обсуждение в PR, существенные решения в ADR и краткую подборку со ссылками на обзорной странице.

#### Принятое решение

Выбран третий вариант. [Обзорная страница](#appendix-engineering-decisions) связывает вопрос владельца, обоснование, изменение, доказательства и статус. PR остаётся местом ревью кода; ADR описывает одно существенное решение с последствиями. Исходное назначение ADR — сохранять решение и его обоснование: [описание формата](https://adr.github.io/).

Владелец оставляет замечания в PR или чате. Исполнитель поддерживает записи и ссылки. При переносе из чата он публикует пересказ от AI-Your-Bro с явным указанием источника и автора требования. GitHub-комментарии, approvals и решения от аккаунта владельца не создаются исполнителем.

#### Последствия

- Для следующего агента сохраняется и технический результат, и причина выбора.
- Статусы proposed, BACKLOG, на ревью и принято различаются; приёмка ссылается на фактическое основание.
- Рутинные исправления не требуют отдельного ADR. Обзорная страница содержит значимые примеры, а не полный чат.
- После изменения решения обновляются PR/ADR, обзорная страница и переносимый контекст.
- Границы сервисов G0, контрольные точки и право владельца на review/merge сохраняются.

<a id="appendix-adr-0003"></a>

## Приложение: ADR 0003: proposed Go container smoke

Полный исходный документ: `docs/adr/0003-go-container-smoke.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### ADR 0003: Go для контейнерного smoke

- Обсуждение: 15 октября 2026 года; решение владельца: 16 октября 2026 года, Europe/Moscow.
- Статус решения: **accepted** — выбран перенос на Go. Выполнение и приёмка кода отслеживаются в [Issue #9](https://github.com/Includeoyasi/faultline-lab/issues/9) и его PR.
- Источник: прямое поручение владельца в чате взять #9 с F02/F04/F07/F08; ранее в ревью F01 из [аудита #20](https://github.com/Includeoyasi/faultline-lab/issues/20) выбран перенос вместо отдельного ремонта Node.js runner.

#### Причина

Исходный Node.js smoke из PR #8 проверял контейнеры и собирал JSON/JSONL. Он использовал runtime документационных проверок, но оставлял основной инструмент лаборатории на другом языке. Аудит также воспроизвёл потерю отчёта и оставшиеся контейнеры при прерывании процесса.

#### Выбор и последствия

Сохраняем сценарий, коды выхода и изоляцию, переносим orchestration в Go CLI. Он собирается закреплённым Docker toolchain и запускается готовым host-бинарником без Node.js и обязательного host Go SDK. Общий Go CI проверяет инструмент вместе с сервисами. Документационная проверка остаётся на Node.js.

Отмена — часть переноса: отдельный бюджет cleanup, ранний checkpoint отчёта и завершение дерева команды. Docker CLI может запускать Compose/buildx как плагины; остановки только родительского процесса недостаточно. Unix использует process group; Windows — Job Object, назначенный при создании процесса. Для Windows применяется уже закреплённый x/sys/windows, а не новый модуль. Платформенный код требует отдельных проверок и оправдан необходимостью ограничить дочерние процессы без гонки Start→Assign.

Альтернатива — небольшой ремонт старого runner — позволяла быстрее закрыть сигнал, но сохраняла отдельный runtime и последующий повторный перенос. Владелец выбрал совмещённую работу. Сервисные границы и бизнес-контракты не меняются.

Формат отчёта получает явную версию 2; изменения перечислены в [руководстве](#go-smoke-guide). macOS cross-build не выдаётся за runtime-проверку личного Mac. Политика полноценной платформенной provenance остаётся отдельным F03.

#### Проверяемость

[Спецификация](#appendix-go-smoke-backlog) задаёт критерии эквивалентности, ошибок, отмены, cleanup и пяти target-сборок. Новое решение считается принятым направлением; готовность конкретного кода определяется результатами его PR и решением владельца. Официальные основания: [os/exec](https://pkg.go.dev/os/exec), [Windows Job Objects](https://learn.microsoft.com/en-us/windows/win32/procthread/job-objects).

<a id="appendix-go-smoke-backlog"></a>

## Приложение: Backlog #9: перенос smoke на Go, объём и критерии

Полный исходный документ: `docs/development/backlog/GO_SMOKE.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Backlog: перенос контейнерного smoke на Go

Выполнение и ревью — в [GitHub Issue #9](https://github.com/Includeoyasi/faultline-lab/issues/9). Выбранный подход: [ADR 0003](#appendix-adr-0003).

#### Цель

Перенести управление контейнерным smoke из `tools/smoke.mjs` в Go-инструмент: единый язык основных инструментов лаборатории, общие Go-проверки и запуск smoke без Node.js. Проверки самих сервисов в `*_test.go` уже написаны на Go.

#### Место в плане

Перенос относится к инструментам лаборатории и связан с F01/F04 из [аудита #20](https://github.com/Includeoyasi/faultline-lab/issues/20). Порядок и приёмка конкретного PR ведутся в GitHub. Эта спецификация не расширяет первоначальный P01 и не означает приёмку G1.

#### Объём

- Перенести Docker CLI orchestration, сбор результатов, JSON/JSONL-отчёт и cleanup в небольшой Go CLI.
- Сохранить текущую проверку трёх сервисов, явный тестовый Compose-профиль и уникальный проект запуска.
- Предусмотреть документированный путь сборки через закреплённый контейнерный toolchain и запуск готового бинарника; host Go SDK остаётся необязательным для пользователя smoke.
- Обновить CI, команды запуска и переносимый контекст; удалить `tools/smoke.mjs` после проверки замены.
- Проверка документации `tools/check_docs.mjs` остаётся отдельным Node.js-инструментом. Задача не удаляет весь JavaScript из репозитория и не добавляет бизнес-функциональность.

#### Критерии приёмки

- Go CLI запускает Docker напрямую с массивом аргументов, без обязательного shell, Docker socket в контейнере и дополнительных Go-библиотек без обоснования.
- Сохранены build → start → health/security → orders isolation → stop/logs → cleanup; исходы и реальные image IDs присутствуют в отчёте.
- Сохранены коды 0=PASS, 1=FAIL, 2=NOT_RUN и совместимый формат JSON/JSONL. Изменение схемы, если понадобится, явно описано и проверено.
- У команд есть сроки ожидания; ошибка или поддерживаемое ОС прерывание запуска приводит к ограниченному cleanup собственного Compose-проекта. Чужие ресурсы не удаляются.
- Значимые ошибки проверены: отсутствующий Docker/backend, ошибка build/health/stop и ошибка cleanup; диагностика сохраняется, секреты/персональные пути не публикуются.
- Сборка CLI проверена для Windows AMD64, Linux AMD64/ARM64 и macOS AMD64/ARM64; реальные smoke выполнены на доступных Windows/Linux. Mac runtime отдельно остаётся NOT_RUN до доступа к машине.
- CLI проходит применимые format/build/lint/unit/race проверки общего Go CI. Перед заменой показана эквивалентность результатов прежнего и нового runner на одной среде.
- Основной путь smoke больше не требует Node.js или host Go SDK; сборка/получение бинарника, prerequisites и запуск описаны явно.
- Обновлены инструкции, workflow, документы решений и PROJECT_CONTEXT.md; после локальных проверок Docker Desktop/WSL выключены по правилу личного ПК, установленные инструменты и полезные кеши сохранены.
- Отдельный PR принят владельцем; задача закрывается после review/merge, а не после одной успешной сборки.

#### Источники и проверяемая исходная точка

- [Исходный PR #8](https://github.com/Includeoyasi/faultline-lab/pull/8).
- [Node.js smoke](https://github.com/Includeoyasi/faultline-lab/blob/445d324ea48419133fa37b15d44242bd6bbf5b08/tools/smoke.mjs).
- [Windows baseline и исходные данные](https://github.com/Includeoyasi/faultline-lab/blob/445d324ea48419133fa37b15d44242bd6bbf5b08/docs/development/evidence/p01/WINDOWS_SMOKE.md).

Основание — обсуждение владельцем и исполнителем выбора Node.js/Go в чате. Постановку документирует AI-Your-Bro; GitHub approval от имени владельца не создаётся. Требования бизнес-сервисов G0 и их ровно три границы не меняются.

<a id="appendix-p02"></a>

## Приложение: Предложенный P02: задачи, зависимости, критерии и порядок PR

Полный исходный документ: `docs/development/packets/P02.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### P02 — инфраструктура и основа данных

Границы двух задач [milestone P02](https://github.com/Includeoyasi/faultline-lab/milestone/1). Текущие стадии, исполнение и ревью ведутся только в GitHub.

Основа: [архитектура](#appendix-design-review), [roadmap](#appendix-roadmap) и [порядок разработки](#appendix-development-process).

#### [P02: Подготовить PostgreSQL, миграции и исходные данные](https://github.com/Includeoyasi/faultline-lab/issues/12)

##### Описание

Подготовить воспроизводимую основу данных для Orders, Inventory и Payments: PostgreSQL в Compose, отдельные базы и роли, версионируемые миграции и начальный каталог Inventory. Новый стенд должен подготавливаться документированными командами, а повторная подготовка и перезапуски — сохранять существующие данные.

##### Критерии приёмки

- PostgreSQL запускается из чистого checkout через Compose с закреплённым образом и постоянным volume. Настройки окружения вынесены из Go-кода; секреты отсутствуют в Git и логах.
- Созданы orders_db, inventory_db и payments_db с отдельными прикладными ролями и правами миграции. Разрешённые подключения и операции проверены на реальной БД.
- Прикладные роли не имеют superuser/DDL-прав и доступа к чужим базам; отрицательные проверки подтверждают запреты.
- Миграции применяются по сервисам с нуля; повторный запуск не меняет схему. Одновременный запуск не дублирует версии; ошибка возвращает ненулевой код и не помечает миграцию успешно применённой.
- Минимальная схема Inventory и seed создают воспроизводимые SKU, цены в минимальных денежных единицах и исходные остатки. Повтор seed не дублирует товары и не перезаписывает изменённые остатки.
- Штатный перезапуск сохраняет схему и данные. Сброс выполняется отдельной явной командой только для выделенного стенда.
- Проверки прав, миграций, seed и сохранности данных включены в CI и проходят вместе с существующими Go-проверками.
- Инструкция позволяет воспроизвести подготовку, проверки и остановку через контейнерные инструменты; указаны фактические результаты и ограничения среды.

**Порядок:** первый PR P02; после #3.

#### [P02: Подключить Kafka и Redis с интеграционными проверками и CI](https://github.com/Includeoyasi/faultline-lab/issues/16)

##### Описание

Дополнить стенд из #12 одним Kafka broker и Redis. Проверить доставку и сохранность сообщений, доступ к кешу и TTL, затем автоматизировать проверку полного состава инфраструктуры. Сохранить инструкции и результаты, по которым можно воспроизвести запуск, оценить потребление ресурсов и безопасно остановить стенд.

##### Критерии приёмки

- Kafka и Redis запускаются в Compose с закреплёнными образами, явной конфигурацией и проверками готовности. Адреса задаются через DNS стенда; секреты не публикуются.
- Реальный producer/consumer подтверждает отправку и чтение сообщений в изолированном топике. После штатного перезапуска брокера записи сохраняются на volume и доступны для чтения.
- Для Redis проверены разрешённый и запрещённый доступ, запись, чтение и истечение TTL. Ошибки подключения завершаются в пределах заданного таймаута; тест использует собственные ключи.
- Интеграционный прогон проверяет совместный запуск PostgreSQL, Kafka и Redis, перезапуски и сохранность данных. У него отдельный Compose project, ограниченные сроки команд и диагностика при сбое.
- В Linux CI AMD64/ARM64 проходят существующие Go-проверки и новые инфраструктурные проверки. Выполнен smoke на доступном Windows/WSL2; macOS отражается отдельным результатом и остаётся NOT_RUN до реальной проверки.
- Cleanup после успеха и ошибки затрагивает только ресурсы своего тестового запуска. Обычная остановка сохраняет данные; после локальных проверок Docker Desktop/WSL остановлены по правилу личного ПК.
- Инструкция и отчёт содержат команды запуска/проверки/остановки, версии и digests, лимиты, наблюдаемое потребление ресурсов и результаты PASS/FAIL/NOT_RUN. Переносимый контекст обновлён.

- По F09 из [#20](https://github.com/Includeoyasi/faultline-lab/issues/20) зафиксированы гарантия порядка событий, требования к Kafka/Go-клиенту, способ защиты от старого publisher после смены lease и сценарий проверки takeover. Реализация и полная проверка outbox выполняются на его этапе; один SQL lease или producer idempotence не считаются достаточной защитой.

**Порядок:** второй PR P02; после #12.

Порядок: #12 → #16. Проверки и документация входят в обе задачи. Перенос Go smoke отслеживается в #9. Текущий WSL-бюджет — 2 CPU, memory=2GB, swap=1GB; его достаточность проверяется при реализации, автоматически ресурсы не увеличиваются.

<a id="go-smoke-guide"></a>

## Приложение: Go container smoke guide

Полный исходный документ: `docs/development/SMOKE.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Контейнерный smoke на Go

Smoke собирает и запускает три сервиса, проверяет health, изоляцию и остановку, сохраняет отчёт и удаляет ресурсы тестового запуска.

#### Сборка без Go SDK на хосте

Нужны Git, Docker с Linux backend и Compose. Соберите бинарник для своей ОС через [контейнерный toolchain](#appendix-toolchain) и запустите его на хосте.

Windows AMD64, из корня checkout:

```powershell
docker compose build toolchain
docker compose run --rm -e GOOS=windows -e GOARCH=amd64 toolchain smoke-build
.\bin\smoke-windows-amd64.exe
```

Для macOS укажите `GOOS=darwin` и архитектуру хоста `arm64` либо `amd64`; для Linux — `GOOS=linux` и соответствующую архитектуру. Имена файлов: `bin/smoke-darwin-arm64`, `bin/smoke-darwin-amd64`, `bin/smoke-linux-arm64`, `bin/smoke-linux-amd64`. При необходимости на Unix выполните `chmod +x` для собранного бинарника. Команда `docker compose run --rm toolchain smoke-cross-check` собирает все пять вариантов; сборка не доказывает runtime-совместимость машины.

Если подходящий Go уже установлен, допустим `go run ./cmd/smoke`. В CI используется готовый бинарник, собранный контейнерным toolchain.

#### Профили и таймауты

Обычный запуск сервисов использует выбранный `FL_RUNTIME_ENV_FILE` и независимый `FL_STOP_GRACE_PERIOD` Compose. Перед `docker compose up` проверьте эффективный профиль:

```powershell
.\bin\smoke-windows-amd64.exe -check-config
docker compose up -d --wait orders inventory payments
```

Проверка не запускает сервисы и не требует работающего backend. Она читает результат `docker compose config --format json` только в памяти: полный конфиг может содержать секреты и не выводится. У каждого сервиса `stop_grace_period` должен покрывать `FL_SHUTDOWN_TIMEOUT` и минимум две дополнительные секунды на завершение вне RPC-drain. Это минимальный запас каркаса, а не универсальная гарантия для будущих зависимостей. Для них бюджет пересматривается по наблюдениям. Значения 10s/15s допустимы; 20s/15s отклоняются.

Для намеренно неправильного профиля эксперимента есть `-check-config -allow-unsafe-shutdown`: несогласованные бюджеты дают явные предупреждения. Некорректные duration и отсутствующие сервисы всё равно отклоняются. Обычный `docker compose up` сам Go-проверку не вызывает: оператор должен выполнить preflight; автоматический smoke делает это обязательно.

Сам smoke всегда выбирает `configs/env/container.env.example` и grace=15s, независимо от личных env-файлов. `-root` задаёт checkout, остальные флаги смотрите через `-help`.

#### Проверяемый сценарий

1. Проверить Linux Docker/Compose и Git; сохранить версии и исходный commit/dirty.
2. Проверить бюджеты явного smoke-профиля, собрать и поднять три сервиса.
3. Проверить health, UID, read-only root, capabilities, no-new-privileges, внутреннюю сеть и отсутствие host ports/mounts.
4. Остановить Orders, проверить exit/logs; Inventory и Payments должны продолжать отвечать.
5. Остановить остальные, проверить JSON start → NOT_SERVING → graceful stop и exit 0.
6. Собрать диагностику, удалить только собственный Compose project и проверить отсутствие его контейнеров и сетей.

#### Отмена и восстановление

Ctrl+C, а на Unix также SIGTERM, отменяет запуск. Новые этапы не начинаются. На Unix Docker CLI и плагины принадлежат отдельной группе процессов; на Windows — ненаследуемому Job Object, назначенному при создании процесса. Отмена завершает дерево команды до очистки. После выхода родителя Unix executor также завершает оставшихся участников своей группы, включая потомка с перенаправленным stdio. Windows использует уже закреплённый `golang.org/x/sys/windows` для Job API; дополнительный модуль не добавляется.

У обычной команды бюджет 60s, у build — 10 минут. Очистка получает новый независимый контекст на 60s; завершение OS-процессов и чтение их потоков дополнительно ограничены. Отменённый основной контекст не блокирует `down`. Повторная отмена не прерывает этот ограниченный cleanup. Образы с уникальными именами своего прогона удаляются через `down --rmi local`, также после ошибки или отмены сборки. Cleanup проверяет отсутствие контейнеров, сетей и тегов образов прогона. Чужие проекты и полезные общие кеши сохраняются.

`artifacts/smoke/<run-id>/report.json` появляется до первого запуска Docker, затем обновляется по этапам и командам. После SIGKILL, завершения ОС или недоступности диска обработчик может не исполниться. Незаконченный отчёт и `recovery_command` позволяют убрать конкретный project, например `docker compose -f compose.yaml -p <project-from-report> down --remove-orphans --rmi local`. Никогда не заменяйте адресную очистку глобальным prune.

#### Формат результата

Коды сохраняются: **0=PASS, 1=FAIL, 2=NOT_RUN**. Недоступные prerequisites дают NOT_RUN; прерывание — FAIL с `cancelled=true`. Ошибка cleanup делает весь прогон FAIL.

JSON schema version 2 сохраняет `run_id`, `project`, `result`, `services`, `commands`, `cleanup`, source revision/dirty и реальные image IDs. Вместо `node` записывается `go_version`; добавлены `phase`, `cancelled` и `recovery_command`. Поле архитектуры из сервиса удалено: `docker_architecture` обозначает только backend; выбор полного контракта платформ остаётся в F03. Итоговые JSONL и ошибки команд находятся рядом. Неизвестный OOM не записывается как false: наблюдённые exit/OOM/ID сохраняются также при неуспешной проверке.

Отчёт не содержит полного inspect Env или разрешённого Compose-конфига. Пути checkout и домашнего каталога редактируются в диагностике. Smoke использует несекретный профиль; сторонние произвольные логи перед публикацией нужно проверять отдельно. Это проверка каркаса, а не измерение бизнес-корректности или capacity.

После локальной проверки остановите Docker Desktop и его WSL VM по [правилу личного компьютера](https://github.com/Includeoyasi/faultline-lab/blob/1c54d9c6aee2be4da7c38b7775b6d862402caac1/docs/development/CONTAINER_RUNTIME.md#личный-компьютер-запуск-по-необходимости), сохранив установленные инструменты и кеши.

<a id="go-smoke-main"></a>

## Приложение: Go smoke entrypoint

Полный исходный документ: `cmd/smoke/main.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
// Command smoke checks the containerized service skeleton from the host.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/Includeoyasi/faultline-lab/internal/smoke"
)

func main() { os.Exit(mainCode()) }

func mainCode() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return run(ctx, os.Args[1:], os.Stdout, os.Stderr)
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("smoke", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("root", ".", "repository directory")
	checkConfig := flags.Bool("check-config", false, "validate the effective local Compose shutdown budgets without starting services")
	allowUnsafe := flags.Bool("allow-unsafe-shutdown", false, "explicitly allow incompatible shutdown budgets in check-config experiments")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 || (*allowUnsafe && !*checkConfig) {
		_, _ = fmt.Fprintln(stderr, "unexpected arguments; use -help")
		return 2
	}
	if *checkConfig {
		if err := smoke.CheckConfig(ctx, *root, *allowUnsafe, stdout); err != nil {
			_, _ = fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	}
	report, err := smoke.Run(ctx, smoke.Options{Root: *root, Output: stdout})
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
	}
	return report.ExitCode()
}
```

<a id="go-smoke-main-tests"></a>

## Приложение: Smoke CLI tests

Полный исходный документ: `cmd/smoke/main_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package main

import (
	"bytes"
	"context"
	"testing"
)

func TestCLIUsage(t *testing.T) {
	for _, tc := range []struct {
		args []string
		code int
	}{
		{[]string{"-help"}, 0},
		{[]string{"unexpected"}, 2},
		{[]string{"-unknown"}, 2},
		{[]string{"-allow-unsafe-shutdown"}, 2},
	} {
		var output bytes.Buffer
		if code := run(context.Background(), tc.args, &output, &output); code != tc.code {
			t.Fatalf("%v: code=%d output=%s", tc.args, code, output.String())
		}
	}
}
```

<a id="go-smoke-run"></a>

## Приложение: Smoke lifecycle

Полный исходный документ: `internal/smoke/run.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
// Package smoke owns a single isolated Docker Compose smoke run and its cleanup.
package smoke

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var services = []string{"orders", "inventory", "payments"}

// Options controls execution budgets, never service deployment defaults.
type Options struct {
	Root           string
	Artifacts      string
	Output         io.Writer
	CommandTimeout time.Duration
	BuildTimeout   time.Duration
	CleanupTimeout time.Duration
}

type runner struct {
	options           Options
	directory         string
	report            Report
	execute           executeFunc
	resourcesMayExist bool
	redact            *strings.Replacer
}

// Run invokes Docker directly on the host, without a shell or Docker socket mount.
func Run(ctx context.Context, options Options) (Report, error) {
	root, err := filepath.Abs(options.Root)
	if err != nil {
		return Report{}, errors.New("cannot resolve repository directory")
	}
	options.Root = root
	executor := processExecutor{directory: root, env: smokeEnvironment()}
	return run(ctx, options, executor.execute)
}

func smokeEnvironment() []string {
	env := make([]string, 0, len(os.Environ())+2)
	for _, value := range os.Environ() {
		key, _, _ := strings.Cut(value, "=")
		if strings.EqualFold(key, "FL_RUNTIME_ENV_FILE") || strings.EqualFold(key, "FL_STOP_GRACE_PERIOD") {
			continue
		}
		env = append(env, value)
	}
	return append(env, "FL_RUNTIME_ENV_FILE=configs/env/container.env.example", "FL_STOP_GRACE_PERIOD=15s")
}

func run(ctx context.Context, options Options, execute executeFunc) (Report, error) {
	r, err := newRunner(options, execute)
	if err != nil {
		return Report{}, err
	}
	if err = r.checkpoint(); err != nil {
		return r.report, errors.New("cannot save initial smoke report")
	}
	err = r.exercise(ctx)
	if err != nil {
		r.report.Result = "FAIL"
		if r.report.Phase == "prerequisites" && ctx.Err() == nil {
			r.report.Result = "NOT_RUN"
		}
		r.report.FailedPhase, r.report.Error = r.report.Phase, r.redact.Replace(err.Error())
	} else {
		r.report.Result = "PASS"
	}
	// The exercise has returned and joined its current child. Cleanup gets a new
	// budget: cancellation of the run must not cancel resource removal as well.
	cleanupCtx, cancel := context.WithTimeout(context.Background(), r.options.CleanupTimeout)
	cleanupErr := r.cleanup(cleanupCtx)
	cancel()
	if cleanupErr != nil {
		r.report.Result, r.report.Cleanup = "FAIL", "FAIL"
		r.report.CleanupError = r.redact.Replace(cleanupErr.Error())
	}
	if ctx.Err() != nil {
		r.report.Cancelled, r.report.Result = true, "FAIL"
		if r.report.Error == "" {
			r.report.Error = ctx.Err().Error()
		}
	}
	finished := time.Now().UTC()
	r.report.FinishedAt = &finished
	saveErr := r.checkpoint()
	if saveErr != nil {
		r.report.Result = "FAIL"
	}
	_, _ = fmt.Fprintf(r.options.Output, "%s: artifacts/smoke/%s/report.json\n", r.report.Result, r.report.RunID)
	return r.report, errors.Join(err, cleanupErr, saveErr)
}

func newRunner(options Options, execute executeFunc) (*runner, error) {
	if options.Output == nil {
		options.Output = io.Discard
	}
	if options.CommandTimeout <= 0 {
		options.CommandTimeout = time.Minute
	}
	if options.BuildTimeout <= 0 {
		options.BuildTimeout = 10 * time.Minute
	}
	if options.CleanupTimeout <= 0 {
		options.CleanupTimeout = time.Minute
	}
	if options.Artifacts == "" {
		options.Artifacts = filepath.Join(options.Root, "artifacts", "smoke")
	}
	var token [16]byte
	if _, err := rand.Read(token[:]); err != nil {
		return nil, errors.New("cannot create smoke run ID")
	}
	id := hex.EncodeToString(token[:])
	id = id[:8] + "-" + id[8:12] + "-" + id[12:16] + "-" + id[16:20] + "-" + id[20:]
	project := "faultline-smoke-" + id[:8]
	directory := filepath.Join(options.Artifacts, id)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return nil, errors.New("cannot create smoke artifact directory")
	}
	pairs := []string{options.Root, "[workspace]", filepath.ToSlash(options.Root), "[workspace]"}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		pairs = append(pairs, home, "[home]", filepath.ToSlash(home), "[home]")
	}
	return &runner{options: options, directory: directory, execute: execute, redact: strings.NewReplacer(pairs...), report: Report{
		SchemaVersion: 2, RunID: id, Project: project, StartedAt: time.Now().UTC(), Result: "FAIL", Phase: "prerequisites", Cleanup: "PENDING",
		HostPlatform: runtime.GOOS, GoVersion: runtime.Version(), Services: map[string]serviceReport{}, Commands: []commandReport{},
		Recovery:    []string{"docker", "compose", "-f", "compose.yaml", "-p", project, "down", "--remove-orphans", "--rmi", "local"},
		Limitations: []string{"No business operations or capacity measurement", "docker_architecture describes the backend, not individual image platforms (F03)", "SIGKILL or host failure requires recovery_command; signal cleanup cannot be guaranteed in those cases"},
	}}, nil
}

func (r *runner) command(ctx context.Context, program string, args []string, timeout time.Duration, sensitive bool) (commandOutput, error) {
	if err := ctx.Err(); err != nil {
		return commandOutput{}, err
	}
	commandCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	observation := commandReport{Program: program, Args: args, StartedAt: time.Now().UTC()}
	result, err := r.execute(commandCtx, program, args)
	observation.FinishedAt, observation.ExitCode = time.Now().UTC(), result.exitCode
	if err != nil {
		observation.Error = r.redact.Replace(err.Error())
		if commandCtx.Err() != nil {
			observation.Error = r.redact.Replace(errors.Join(commandCtx.Err(), err).Error())
		}
		if !sensitive {
			name := fmt.Sprintf("command-%d.log", len(r.report.Commands)+1)
			diagnostic := r.redact.Replace(result.stdout + result.stderr)
			if writeErr := os.WriteFile(filepath.Join(r.directory, name), []byte(diagnostic), 0o600); writeErr != nil {
				observation.Error += "; cannot save command diagnostics"
			}
		}
		err = fmt.Errorf("%s %s failed (exit %d; %s)", program, strings.Join(args, " "), result.exitCode, observation.Error)
	}
	r.report.Commands = append(r.report.Commands, observation)
	return result, errors.Join(err, r.checkpoint())
}

func (r *runner) docker(ctx context.Context, args ...string) (commandOutput, error) {
	return r.command(ctx, "docker", args, r.options.CommandTimeout, false)
}

func (r *runner) composeArgs(args ...string) []string {
	return append([]string{"compose", "-f", "compose.yaml", "-p", r.report.Project}, args...)
}

func (r *runner) compose(ctx context.Context, args ...string) (commandOutput, error) {
	return r.docker(ctx, r.composeArgs(args...)...)
}

func (r *runner) exercise(ctx context.Context) error {
	steps := []struct {
		name string
		work func(context.Context) error
	}{
		{"prerequisites", r.prerequisites},
		{"configuration", r.configuration},
		{"build", func(ctx context.Context) error {
			r.resourcesMayExist = true // A partial build can leave run-specific images.
			_, err := r.command(ctx, "docker", r.composeArgs(append([]string{"build"}, services...)...), r.options.BuildTimeout, false)
			return err
		}},
		{"start", func(ctx context.Context) error {
			r.resourcesMayExist = true
			_, err := r.compose(ctx, append([]string{"up", "-d", "--wait", "--wait-timeout", "60"}, services...)...)
			return err
		}},
		{"check", r.checkAll},
		{"isolation", r.isolate},
		{"stop", r.stopAll},
	}
	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		r.report.Phase = step.name
		if err := r.checkpoint(); err != nil {
			return errors.New("cannot checkpoint smoke phase")
		}
		_, _ = fmt.Fprintf(r.options.Output, "Runtime smoke %s: %s\n", r.report.Project, step.name)
		if err := step.work(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) prerequisites(ctx context.Context) error {
	version, err := r.docker(ctx, "version", "--format", "{{.Server.Version}}")
	if err != nil {
		return err
	}
	r.report.DockerVersion = strings.TrimSpace(version.stdout)
	platform, err := r.docker(ctx, "info", "--format", "{{.OSType}}/{{.Architecture}}")
	if err != nil {
		return err
	}
	r.report.DockerArch = strings.TrimSpace(platform.stdout)
	if !strings.HasPrefix(r.report.DockerArch, "linux/") {
		return errors.New("docker must run Linux containers")
	}
	compose, err := r.docker(ctx, "compose", "version", "--short")
	if err != nil {
		return err
	}
	r.report.Compose = strings.TrimSpace(compose.stdout)
	revision, err := r.command(ctx, "git", []string{"rev-parse", "HEAD"}, r.options.CommandTimeout, false)
	if err != nil {
		return err
	}
	r.report.Revision = strings.TrimSpace(revision.stdout)
	dirty, err := r.command(ctx, "git", []string{"status", "--porcelain"}, r.options.CommandTimeout, true)
	r.report.Dirty = strings.TrimSpace(dirty.stdout) != ""
	return err
}

func (r *runner) configuration(ctx context.Context) error {
	resolved, err := r.command(ctx, "docker", r.composeArgs("config", "--format", "json"), r.options.CommandTimeout, true)
	if err != nil {
		return err
	}
	_, err = validateConfig([]byte(resolved.stdout), false)
	return err
}

func (r *runner) cleanup(ctx context.Context) error {
	if !r.resourcesMayExist {
		r.report.Cleanup = "NOT_NEEDED"
		return nil
	}
	for _, service := range services {
		logCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, err := r.collectLogs(logCtx, service)
		cancel()
		if err != nil {
			observation := r.report.Services[service]
			observation.DiagnosticError = r.redact.Replace(err.Error())
			r.report.Services[service] = observation
		}
	}
	// Recheck engine state after the joined up command and scoped down. A second
	// observation catches resources still being finalized by the Docker engine.
	for {
		if _, err := r.compose(ctx, "down", "--remove-orphans", "--rmi", "local"); err != nil {
			return err
		}
		for attempt := 0; attempt < 2; attempt++ {
			for _, args := range [][]string{
				{"ps", "--all", "--quiet", "--filter", "label=com.docker.compose.project=" + r.report.Project},
				{"network", "ls", "--quiet", "--filter", "label=com.docker.compose.project=" + r.report.Project},
				{"image", "ls", "--quiet", "--filter", "reference=" + r.report.Project + "-*:latest"},
			} {
				remaining, err := r.docker(ctx, args...)
				if err != nil {
					return err
				}
				if strings.TrimSpace(remaining.stdout) != "" {
					goto retry
				}
			}
			if attempt == 0 {
				timer := time.NewTimer(200 * time.Millisecond)
				select {
				case <-ctx.Done():
					timer.Stop()
					return ctx.Err()
				case <-timer.C:
				}
			}
		}
		r.report.Cleanup = "PASS"
		return nil
	retry:
		if err := ctx.Err(); err != nil {
			return err
		}
		timer := time.NewTimer(200 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}
```

<a id="go-smoke-report"></a>

## Приложение: Smoke report

Полный исходный документ: `internal/smoke/report.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package smoke

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Report version 2 preserves result/cleanup/services/commands and exit meanings.
// It replaces Node metadata with Go metadata and adds cancellation/recovery data.
type Report struct {
	SchemaVersion int                      `json:"schema_version"`
	RunID         string                   `json:"run_id"`
	Project       string                   `json:"project"`
	StartedAt     time.Time                `json:"started_at"`
	FinishedAt    *time.Time               `json:"finished_at,omitempty"`
	Result        string                   `json:"result"`
	Phase         string                   `json:"phase"`
	FailedPhase   string                   `json:"failed_phase,omitempty"`
	Error         string                   `json:"error,omitempty"`
	Cancelled     bool                     `json:"cancelled"`
	Cleanup       string                   `json:"cleanup"`
	CleanupError  string                   `json:"cleanup_error,omitempty"`
	HostPlatform  string                   `json:"host_platform"`
	GoVersion     string                   `json:"go_version"`
	DockerVersion string                   `json:"docker_version,omitempty"`
	DockerArch    string                   `json:"docker_architecture,omitempty"`
	Compose       string                   `json:"compose_version,omitempty"`
	Revision      string                   `json:"source_revision,omitempty"`
	Dirty         bool                     `json:"source_dirty"`
	Isolation     string                   `json:"isolation_after_orders_stop,omitempty"`
	Services      map[string]serviceReport `json:"services"`
	Commands      []commandReport          `json:"commands"`
	Recovery      []string                 `json:"recovery_command"`
	Limitations   []string                 `json:"limitations"`
}

type serviceReport struct {
	ID              string         `json:"id,omitempty"`
	ImageID         string         `json:"image_id,omitempty"`
	Healthy         bool           `json:"healthy"`
	Stopped         bool           `json:"stopped"`
	ExitCode        *int           `json:"exit_code,omitempty"`
	OOMKilled       *bool          `json:"oom_killed,omitempty"`
	Status          string         `json:"status,omitempty"`
	StopMode        string         `json:"stop_mode,omitempty"`
	Security        map[string]any `json:"security,omitempty"`
	DiagnosticError string         `json:"diagnostic_error,omitempty"`
}

type commandReport struct {
	Program    string    `json:"program"`
	Args       []string  `json:"args"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	ExitCode   int       `json:"exit_code"`
	Error      string    `json:"error,omitempty"`
}

// ExitCode keeps the original smoke's public exit contract.
func (r Report) ExitCode() int {
	switch r.Result {
	case "PASS":
		return 0
	case "NOT_RUN":
		return 2
	default:
		return 1
	}
}

func (r *runner) checkpoint() error {
	data, err := json.MarshalIndent(r.report, "", "  ")
	if err != nil {
		return err
	}
	temporary := filepath.Join(r.directory, "report.json.tmp")
	if err = os.WriteFile(temporary, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(temporary, filepath.Join(r.directory, "report.json"))
}
```

<a id="go-smoke-checks"></a>

## Приложение: Smoke container checks

Полный исходный документ: `internal/smoke/checks.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package smoke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var containerIDPattern = regexp.MustCompile(`^[a-f0-9]{12,64}$`)

type container struct {
	Image  string
	Config struct{ User string }
	State  struct {
		Status    string
		ExitCode  int
		OOMKilled bool
		Health    struct{ Status string }
	}
	HostConfig struct {
		ReadonlyRootfs bool
		Privileged     bool
		CapDrop        []string
		SecurityOpt    []string
		PortBindings   map[string]json.RawMessage
	}
	Mounts          []json.RawMessage
	NetworkSettings struct {
		Networks map[string]struct{ NetworkID string }
	}
}

func (r *runner) inspect(ctx context.Context, service string) (string, container, error) {
	ids, err := r.compose(ctx, "ps", "--all", "-q", service)
	if err != nil {
		return "", container{}, err
	}
	id := strings.TrimSpace(ids.stdout)
	if !containerIDPattern.MatchString(id) {
		return "", container{}, fmt.Errorf("missing unique container for %s", service)
	}
	result, err := r.docker(ctx, "inspect", id)
	if err != nil {
		return "", container{}, err
	}
	var data []container
	if err = json.Unmarshal([]byte(result.stdout), &data); err != nil || len(data) != 1 {
		return "", container{}, fmt.Errorf("invalid container inspect for %s", service)
	}
	return id, data[0], nil
}

func (r *runner) checkAll(ctx context.Context) error {
	for _, service := range services {
		if err := r.checkRunning(ctx, service); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) checkRunning(ctx context.Context, service string) error {
	if _, err := r.compose(ctx, "exec", "-T", service, "/app/service", "healthcheck"); err != nil {
		return err
	}
	id, data, err := r.inspect(ctx, service)
	if err != nil {
		return err
	}
	host := data.HostConfig
	observation := r.report.Services[service]
	observation.ID, observation.ImageID, observation.Status = id, data.Image, data.State.Status
	observation.OOMKilled = &data.State.OOMKilled
	r.report.Services[service] = observation
	valid := data.State.Status == "running" && data.State.Health.Status == "healthy" &&
		data.Config.User == "65532:65532" && host.ReadonlyRootfs && !host.Privileged &&
		slices.Contains(host.CapDrop, "ALL") && slices.ContainsFunc(host.SecurityOpt, func(v string) bool {
		return v == "no-new-privileges" || v == "no-new-privileges:true" || v == "no-new-privileges=true"
	}) &&
		len(host.PortBindings) == 0 && len(data.Mounts) == 0 && len(data.NetworkSettings.Networks) > 0
	if !valid {
		return fmt.Errorf("%s violates health or runtime security requirements", service)
	}
	for _, network := range data.NetworkSettings.Networks {
		result, err := r.docker(ctx, "network", "inspect", network.NetworkID)
		if err != nil {
			return err
		}
		var networks []struct{ Internal bool }
		if err = json.Unmarshal([]byte(result.stdout), &networks); err != nil || len(networks) != 1 || !networks[0].Internal {
			return fmt.Errorf("%s must use an internal network", service)
		}
	}
	observation.ID, observation.ImageID, observation.Healthy = id, data.Image, true
	observation.Security = map[string]any{"user": data.Config.User, "read_only": true, "capabilities_dropped": true, "no_new_privileges": true, "host_ports": []string{}, "mounts": []string{}, "internal_network": true}
	r.report.Services[service] = observation
	return nil
}

func (r *runner) collectLogs(ctx context.Context, service string) (string, error) {
	result, err := r.compose(ctx, "logs", "--no-color", "--no-log-prefix", service)
	if err != nil {
		return "", err
	}
	logs := r.redact.Replace(result.stdout + result.stderr)
	if err = os.WriteFile(filepath.Join(r.directory, service+".jsonl"), []byte(logs), 0o600); err != nil {
		return "", errors.New("cannot save service logs")
	}
	return logs, nil
}

func (r *runner) checkStopped(ctx context.Context, service string) error {
	id, data, err := r.inspect(ctx, service)
	if err != nil {
		return err
	}
	observation := r.report.Services[service]
	observation.ID, observation.ImageID, observation.Status = id, data.Image, data.State.Status
	observation.ExitCode, observation.OOMKilled = &data.State.ExitCode, &data.State.OOMKilled
	r.report.Services[service] = observation
	if data.State.Status != "exited" || data.State.ExitCode != 0 || data.State.OOMKilled {
		return fmt.Errorf("%s failed normal SIGTERM shutdown", service)
	}
	logs, err := r.collectLogs(ctx, service)
	if err != nil {
		return err
	}
	step := 0
	for _, line := range strings.Split(logs, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var record struct {
			Service, Time, Level, Msg, Health string
			StopMode                          string `json:"stop_mode"`
		}
		if err = json.Unmarshal([]byte(line), &record); err != nil || record.Service != service || record.Time == "" || record.Level == "" || record.Msg == "" {
			return fmt.Errorf("%s emitted an invalid JSON record", service)
		}
		switch {
		case step == 0 && record.Msg == "service started":
			step++
		case step == 1 && record.Msg == "service stopping" && record.Health == "NOT_SERVING":
			step++
		case step == 2 && record.Msg == "service stopped" && record.StopMode == "graceful":
			step++
		}
	}
	if step != 3 {
		return fmt.Errorf("%s lacks ordered start/readiness/stop records", service)
	}
	observation.Stopped, observation.ExitCode, observation.StopMode = true, &data.State.ExitCode, "graceful"
	r.report.Services[service] = observation
	return nil
}

func (r *runner) isolate(ctx context.Context) error {
	if _, err := r.compose(ctx, "stop", "orders"); err != nil {
		return err
	}
	if err := r.checkStopped(ctx, "orders"); err != nil {
		return err
	}
	for _, service := range []string{"inventory", "payments"} {
		if err := r.checkRunning(ctx, service); err != nil {
			return err
		}
	}
	r.report.Isolation = "PASS"
	return nil
}

func (r *runner) stopAll(ctx context.Context) error {
	if _, err := r.compose(ctx, "stop", "inventory", "payments"); err != nil {
		return err
	}
	for _, service := range []string{"inventory", "payments"} {
		if err := r.checkStopped(ctx, service); err != nil {
			return err
		}
	}
	return nil
}
```

<a id="go-smoke-config"></a>

## Приложение: Compose shutdown preflight

Полный исходный документ: `internal/smoke/config.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package smoke

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// ShutdownMargin reserves time outside the application's RPC drain deadline.
const ShutdownMargin = 2 * time.Second

// CheckConfig checks the user's effective deployment configuration, including
// the selected env file. Unlike Run it does not force the smoke fixture profile.
// Resolved config can contain secrets: it stays in memory and is never logged.
func CheckConfig(ctx context.Context, root string, allowUnsafe bool, output io.Writer) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	executor := processExecutor{directory: root, env: os.Environ()}
	resolved, err := executor.execute(ctx, "docker", []string{"compose", "-f", "compose.yaml", "config", "--format", "json"})
	if err != nil {
		return fmt.Errorf("cannot resolve Compose config (exit %d); check Docker CLI and compose.yaml", resolved.exitCode)
	}
	warnings, err := validateConfig([]byte(resolved.stdout), allowUnsafe)
	if err != nil {
		return err
	}
	for _, warning := range warnings {
		_, _ = fmt.Fprintln(output, "WARNING:", warning)
	}
	_, err = fmt.Fprintln(output, "Compose shutdown budgets checked; no services started")
	return err
}

func validateConfig(data []byte, allowUnsafe bool) ([]string, error) {
	var config struct {
		Services map[string]struct {
			Environment map[string]string `json:"environment"`
			Grace       string            `json:"stop_grace_period"`
		} `json:"services"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid resolved Compose JSON: %w", err)
	}
	var warnings []string
	for _, name := range services {
		service, ok := config.Services[name]
		if !ok {
			return nil, fmt.Errorf("compose config lacks service %s", name)
		}
		appValue, present := service.Environment["FL_SHUTDOWN_TIMEOUT"]
		if !present {
			appValue = "10s" // The application's documented default.
		}
		app, err := time.ParseDuration(appValue)
		if err != nil || app <= 0 {
			return nil, fmt.Errorf("%s: FL_SHUTDOWN_TIMEOUT must be a positive duration", name)
		}
		graceValue := service.Grace
		if graceValue == "" {
			graceValue = "10s" // Compose's default when the property is omitted.
		}
		grace, err := time.ParseDuration(graceValue)
		if err != nil || grace <= 0 {
			return nil, fmt.Errorf("%s: stop_grace_period must be a positive duration", name)
		}
		// Subtraction avoids overflow for an otherwise valid very long duration.
		if grace < app || grace-app < ShutdownMargin {
			message := fmt.Sprintf("%s: container grace %s must cover app timeout %s plus %s margin", name, grace, app, ShutdownMargin)
			if !allowUnsafe {
				return nil, fmt.Errorf("%s", message)
			}
			warnings = append(warnings, message)
		}
	}
	return warnings, nil
}
```

<a id="go-smoke-exec"></a>

## Приложение: Bounded command output

Полный исходный документ: `internal/command/exec.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package command

import (
	"bytes"
)

const outputLimit = 8 << 20

type commandOutput struct {
	stdout, stderr string
	exitCode       int
}

// processExecutor owns the command tree, including Docker CLI plugins.
type processExecutor struct {
	directory string
	env       []string
}

type boundedBuffer struct {
	buffer    bytes.Buffer
	truncated bool
}

func (b *boundedBuffer) Len() int       { return b.buffer.Len() }
func (b *boundedBuffer) String() string { return b.buffer.String() }

func (b *boundedBuffer) Write(p []byte) (int, error) {
	length := len(p)
	remaining := outputLimit - b.Len()
	if length > remaining {
		p = p[:remaining]
		b.truncated = true
	}
	_, _ = b.buffer.Write(p)
	return length, nil
}
```

<a id="go-smoke-unix"></a>

## Приложение: Unix command process groups

Полный исходный документ: `internal/command/exec_unix.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
//go:build !windows

package command

import (
	"context"
	"errors"
	"os/exec"
	"syscall"
	"time"
)

func (p processExecutor) execute(ctx context.Context, program string, args []string) (commandOutput, error) {
	cmd := exec.CommandContext(ctx, program, args...)
	cmd.Dir, cmd.Env, cmd.WaitDelay = p.directory, p.env, 2*time.Second
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return err
	}
	stdout, stderr := &boundedBuffer{}, &boundedBuffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	err := cmd.Run()
	if cmd.ProcessState != nil {
		// A plugin may close its output handles and outlive a successful parent.
		// The command owns its process group even after the parent exits.
		if killErr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); killErr != nil && !errors.Is(killErr, syscall.ESRCH) {
			err = errors.Join(err, killErr)
		}
	}
	code := -1
	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	}
	if stdout.truncated || stderr.truncated {
		err = errors.Join(err, errors.New("command output exceeded limit"))
	}
	return commandOutput{stdout.String(), stderr.String(), code}, errors.Join(err, ctx.Err())
}
```

<a id="go-smoke-windows"></a>

## Приложение: Windows command jobs

Полный исходный документ: `internal/command/exec_windows.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// JOB_LIST is available since Windows 10. Assigning at CreateProcess avoids the
// race where Docker starts its plugin before AssignProcessToJobObject runs.
const procThreadAttributeJobList = 0x0002000D

type jobAccounting struct {
	TotalUserTime, TotalKernelTime, ThisUserTime, ThisKernelTime int64
	PageFaults, TotalProcesses, ActiveProcesses, Terminated      uint32
}

func (p processExecutor) execute(ctx context.Context, program string, args []string) (commandOutput, error) {
	failed := commandOutput{exitCode: -1}
	if err := ctx.Err(); err != nil {
		return failed, err
	}
	executable, err := exec.LookPath(program)
	if err != nil {
		return failed, err
	}
	application, err := windows.UTF16PtrFromString(executable)
	if err != nil {
		return failed, err
	}
	commandLine, err := windows.UTF16PtrFromString(windows.ComposeCommandLine(append([]string{executable}, args...)))
	if err != nil {
		return failed, err
	}
	directory, err := windows.UTF16PtrFromString(p.directory)
	if err != nil {
		return failed, err
	}
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return failed, err
	}
	defer func() { _ = windows.CloseHandle(job) }()
	var limits windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	limits.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err = windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits))); err != nil {
		return failed, err
	}
	stdin, err := os.Open(os.DevNull)
	if err != nil {
		return failed, err
	}
	defer func() { _ = stdin.Close() }()
	outReader, outWriter, err := os.Pipe()
	if err != nil {
		return failed, err
	}
	defer func() { _ = outReader.Close(); _ = outWriter.Close() }()
	errReader, errWriter, err := os.Pipe()
	if err != nil {
		return failed, err
	}
	defer func() { _ = errReader.Close(); _ = errWriter.Close() }()
	files := []*os.File{stdin, outWriter, errWriter}
	handles := make([]windows.Handle, 0, len(files))
	defer func() {
		for _, handle := range handles {
			_ = windows.CloseHandle(handle)
		}
	}()
	for _, file := range files {
		var duplicate windows.Handle
		if err = windows.DuplicateHandle(windows.CurrentProcess(), windows.Handle(file.Fd()), windows.CurrentProcess(), &duplicate, 0, true, windows.DUPLICATE_SAME_ACCESS); err != nil {
			return failed, err
		}
		handles = append(handles, duplicate)
	}
	attributes, err := windows.NewProcThreadAttributeList(2)
	if err != nil {
		return failed, err
	}
	defer attributes.Delete()
	if err = attributes.Update(windows.PROC_THREAD_ATTRIBUTE_HANDLE_LIST, unsafe.Pointer(&handles[0]), uintptr(len(handles))*unsafe.Sizeof(handles[0])); err != nil {
		return failed, err
	}
	if err = attributes.Update(procThreadAttributeJobList, unsafe.Pointer(&job), unsafe.Sizeof(job)); err != nil {
		return failed, err
	}
	startup := windows.StartupInfoEx{ProcThreadAttributeList: attributes.List()}
	startup.Cb, startup.Flags = uint32(unsafe.Sizeof(startup)), windows.STARTF_USESTDHANDLES
	startup.StdInput, startup.StdOutput, startup.StdErr = handles[0], handles[1], handles[2]
	env := append([]string(nil), p.env...)
	sort.Slice(env, func(i, j int) bool { return strings.ToUpper(env[i]) < strings.ToUpper(env[j]) })
	environment := utf16.Encode([]rune(strings.Join(env, "\x00") + "\x00\x00"))
	var process windows.ProcessInformation
	err = windows.CreateProcess(application, commandLine, nil, nil, true,
		windows.EXTENDED_STARTUPINFO_PRESENT|windows.CREATE_UNICODE_ENVIRONMENT|windows.CREATE_NO_WINDOW,
		&environment[0], directory, &startup.StartupInfo, &process)
	runtime.KeepAlive(handles)
	runtime.KeepAlive(job)
	if err != nil {
		return failed, err
	}
	defer func() { _ = windows.CloseHandle(process.Process) }()
	_ = windows.CloseHandle(process.Thread)
	for _, handle := range handles {
		_ = windows.CloseHandle(handle)
	}
	handles = nil
	_ = outWriter.Close()
	_ = errWriter.Close()
	stdout, stderr := &boundedBuffer{}, &boundedBuffer{}
	drained := make(chan struct{}, 2)
	go func() { _, _ = io.Copy(stdout, outReader); drained <- struct{}{} }()
	go func() { _, _ = io.Copy(stderr, errReader); drained <- struct{}{} }()
	err = waitJob(ctx, job, process.Process)
	// No child inherits the job handle. Closing it after wait also kills any
	// remaining process on a failed termination, before cleanup can be attempted.
	_ = windows.TerminateJobObject(job, 1)
	for range 2 {
		timer := time.NewTimer(2 * time.Second)
		select {
		case <-drained:
			timer.Stop()
		case <-timer.C:
			_ = outReader.Close()
			_ = errReader.Close()
			<-drained
			err = errors.Join(err, errors.New("command output drain timed out"))
		}
	}
	var code uint32
	codeErr := windows.GetExitCodeProcess(process.Process, &code)
	if codeErr == nil && code != 0 {
		err = errors.Join(err, fmt.Errorf("process exited with code %d", code))
	}
	if stdout.truncated || stderr.truncated {
		err = errors.Join(err, errors.New("command output exceeded limit"))
	}
	return commandOutput{stdout.String(), stderr.String(), int(code)}, errors.Join(err, codeErr, ctx.Err())
}

func waitJob(ctx context.Context, job, process windows.Handle) error {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	waitCtx := ctx
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	for {
		var accounting jobAccounting
		if err := windows.QueryInformationJobObject(job, windows.JobObjectBasicAccountingInformation, uintptr(unsafe.Pointer(&accounting)), uint32(unsafe.Sizeof(accounting)), nil); err != nil {
			return err
		}
		state, err := windows.WaitForSingleObject(process, 0)
		if err != nil {
			return err
		}
		if state == windows.WAIT_OBJECT_0 && accounting.ActiveProcesses == 0 {
			return ctx.Err()
		}
		select {
		case <-waitCtx.Done():
			if cancel != nil {
				return errors.New("command tree did not terminate within 5s")
			}
			if err = windows.TerminateJobObject(job, 1); err != nil {
				return err
			}
			waitCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		case <-ticker.C:
		}
	}
}
```

<a id="go-smoke-run-tests"></a>

## Приложение: Smoke failure and cancellation tests

Полный исходный документ: `internal/smoke/run_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package smoke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// dockerFixture models observable engine state and one independent control
// project. It tests ownership and failure paths without pretending to be Docker.
type dockerFixture struct {
	t                 *testing.T
	fail, cancelAt    string
	cancel            context.CancelFunc
	project           string
	owned, foreign    bool
	images            bool
	stickyImages      bool
	stopped           map[string]bool
	operations        []string
	unsafe, oom       bool
	cleanupBlocks     bool
	lateResourceReads int
}

func (f *dockerFixture) execute(ctx context.Context, program string, args []string) (commandOutput, error) {
	f.t.Helper()
	if program == "git" {
		if args[0] == "rev-parse" {
			return commandOutput{stdout: strings.Repeat("a", 40)}, nil
		}
		return commandOutput{}, nil
	}
	if program != "docker" {
		f.t.Fatalf("unexpected executable %s", program)
	}
	op := args[0]
	if op == "compose" && len(args) > 5 {
		project := args[4]
		if !strings.HasPrefix(project, "faultline-smoke-") || (f.project != "" && f.project != project) {
			f.t.Fatalf("unexpected project %s", project)
		}
		f.project = project
		args = args[5:]
		op = args[0]
	}
	f.operations = append(f.operations, op)
	if op == "up" {
		f.owned = true
	}
	if op == "build" {
		f.images = true
	}
	if op == f.cancelAt {
		f.cancel()
		<-ctx.Done()
		return commandOutput{exitCode: -1}, ctx.Err()
	}
	if op == f.fail {
		return commandOutput{exitCode: 42, stderr: "controlled command failure"}, errors.New("fixture failure")
	}
	switch op {
	case "version":
		return commandOutput{stdout: "29.8.2"}, nil
	case "info":
		return commandOutput{stdout: "linux/arm64"}, nil
	case "compose": // compose version
		return commandOutput{stdout: "5.6.0"}, nil
	case "config":
		return commandOutput{stdout: string(resolvedConfig("10s", "15s"))}, nil
	case "build", "up", "exec":
		return commandOutput{}, nil
	case "stop":
		for _, name := range args[1:] {
			f.stopped[name] = true
		}
		return commandOutput{}, nil
	case "ps":
		if slices.Contains(args, "-q") {
			index := slices.Index(services, args[len(args)-1])
			return commandOutput{stdout: strings.Repeat(string(rune('a'+index)), 64)}, nil
		}
		if f.owned {
			return commandOutput{stdout: strings.Repeat("a", 64)}, nil
		}
		if f.lateResourceReads > 0 {
			f.lateResourceReads--
			f.owned = true
			return commandOutput{stdout: strings.Repeat("a", 64)}, nil
		}
		return commandOutput{}, nil
	case "inspect":
		index := int(args[1][0] - 'a')
		name := services[index]
		state, code, oom := "running", 0, false
		if f.stopped[name] {
			state = "exited"
			if f.oom {
				code, oom = 137, true
			}
		}
		user := "65532:65532"
		if f.unsafe {
			user = "root"
		}
		data := fmt.Sprintf(`[{"Image":"sha256:%s","Config":{"User":%q},"State":{"Status":%q,"ExitCode":%d,"OOMKilled":%t,"Health":{"Status":"healthy"}},"HostConfig":{"ReadonlyRootfs":true,"CapDrop":["ALL"],"SecurityOpt":["no-new-privileges:true"]},"Mounts":[],"NetworkSettings":{"Networks":{"runtime":{"NetworkID":"fixture-network"}}}}]`, strings.Repeat("d", 64), user, state, code, oom)
		return commandOutput{stdout: data}, nil
	case "network":
		if args[1] == "inspect" {
			return commandOutput{stdout: `[{"Internal":true}]`}, nil
		}
		return commandOutput{}, nil
	case "logs":
		name := args[len(args)-1]
		return commandOutput{stdout: fmt.Sprintf("{\"service\":%q,\"time\":\"2026-10-16T00:00:00Z\",\"level\":\"INFO\",\"msg\":\"service started\"}\n{\"service\":%q,\"time\":\"2026-10-16T00:00:01Z\",\"level\":\"INFO\",\"msg\":\"service stopping\",\"health\":\"NOT_SERVING\"}\n{\"service\":%q,\"time\":\"2026-10-16T00:00:02Z\",\"level\":\"INFO\",\"msg\":\"service stopped\",\"stop_mode\":\"graceful\"}\n", name, name, name)}, nil
	case "down":
		if !slices.Contains(args, "--rmi") || args[len(args)-1] != "local" {
			f.t.Fatal("cleanup must remove only run-local images")
		}
		if f.cleanupBlocks {
			<-ctx.Done()
			return commandOutput{exitCode: -1}, ctx.Err()
		}
		f.owned = false
		if !f.stickyImages {
			f.images = false
		}
		return commandOutput{}, nil
	case "image":
		if f.images {
			return commandOutput{stdout: "owned-image"}, nil
		}
		return commandOutput{}, nil
	default:
		f.t.Fatalf("unexpected Docker operation %s %v", op, args)
		return commandOutput{}, nil
	}
}

func fixture(t *testing.T) (*dockerFixture, Options) {
	t.Helper()
	return &dockerFixture{t: t, foreign: true, stopped: map[string]bool{}}, Options{Root: t.TempDir(), CleanupTimeout: 2 * time.Second}
}

func verifySavedReport(t *testing.T, options Options, report Report) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(options.Root, "artifacts", "smoke", report.RunID, "report.json"))
	if err != nil {
		t.Fatal(err)
	}
	var saved Report
	if err = json.Unmarshal(data, &saved); err != nil || saved.Result != report.Result || saved.Cleanup != report.Cleanup || saved.Cancelled != report.Cancelled || saved.FinishedAt == nil {
		t.Fatalf("saved report mismatch: %s; %v", data, err)
	}
	if strings.Contains(string(data), options.Root) || len(saved.Recovery) == 0 {
		t.Fatal("report exposes workspace or lacks recovery command")
	}
}

func TestSuccessfulRunRetainsServiceChecks(t *testing.T) {
	f, options := fixture(t)
	f.lateResourceReads = 1
	report, err := run(context.Background(), options, f.execute)
	if err != nil || report.ExitCode() != 0 || report.Cleanup != "PASS" || f.owned || !f.foreign || report.Isolation != "PASS" {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	for _, service := range services {
		observed := report.Services[service]
		if !observed.Healthy || !observed.Stopped || observed.ID == "" || observed.ImageID == "" || observed.ExitCode == nil || *observed.ExitCode != 0 {
			t.Fatalf("incomplete %s observation: %+v", service, observed)
		}
	}
	if count := strings.Count(strings.Join(f.operations, ","), "down"); count != 2 {
		t.Fatalf("late engine resource was not cleaned again: %v", f.operations)
	}
	verifySavedReport(t, options, report)
}

func TestCancellationStopsStagesAndCleansOwnProject(t *testing.T) {
	for _, op := range []string{"build", "up", "exec", "stop"} {
		t.Run(op, func(t *testing.T) {
			f, options := fixture(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			f.cancelAt, f.cancel = op, cancel
			report, err := run(ctx, options, f.execute)
			if err == nil || !report.Cancelled || report.ExitCode() != 1 || f.owned || !f.foreign {
				t.Fatalf("cancel report=%+v err=%v", report, err)
			}
			if report.Cleanup != "PASS" || f.images {
				t.Fatalf("cleanup=%s", report.Cleanup)
			}
			if op == "up" && slices.Contains(f.operations, "exec") {
				t.Fatal("new stage started after cancellation")
			}
			verifySavedReport(t, options, report)
		})
	}
}

func TestImageRemovalConflictFailsCleanup(t *testing.T) {
	f, options := fixture(t)
	f.stickyImages = true
	options.CleanupTimeout = 80 * time.Millisecond
	report, err := run(context.Background(), options, f.execute)
	if err == nil || report.Cleanup != "FAIL" || !f.images || !f.foreign || !slices.Contains(report.Recovery, "--rmi") {
		t.Fatalf("unremoved images hidden: %+v %v", report, err)
	}
	verifySavedReport(t, options, report)
}

func TestCommandFailuresAndCleanupDeadline(t *testing.T) {
	for _, op := range []string{"version", "build", "exec", "stop", "down", "cleanup-timeout"} {
		t.Run(op, func(t *testing.T) {
			f, options := fixture(t)
			f.fail = op
			if op == "cleanup-timeout" {
				f.cleanupBlocks, options.CleanupTimeout = true, 30*time.Millisecond
			}
			report, err := run(context.Background(), options, f.execute)
			wantCode := 1
			if op == "version" {
				wantCode = 2
			}
			if err == nil || report.ExitCode() != wantCode || !f.foreign {
				t.Fatalf("report=%+v err=%v", report, err)
			}
			if op == "down" || op == "cleanup-timeout" {
				if report.Cleanup != "FAIL" || !f.owned {
					t.Fatal("cleanup failure was hidden")
				}
			} else if f.owned {
				t.Fatal("owned resources left behind")
			}
			verifySavedReport(t, options, report)
		})
	}
}

func TestFailedChecksRetainObservedFacts(t *testing.T) {
	for _, scenario := range []string{"unsafe", "oom"} {
		t.Run(scenario, func(t *testing.T) {
			f, options := fixture(t)
			f.unsafe, f.oom = scenario == "unsafe", scenario == "oom"
			report, err := run(context.Background(), options, f.execute)
			observed := report.Services["orders"]
			if err == nil || report.ExitCode() != 1 || observed.ID == "" || observed.ImageID == "" {
				t.Fatalf("missing failed inspection: %+v %v", observed, err)
			}
			if f.oom && (observed.ExitCode == nil || *observed.ExitCode != 137 || observed.OOMKilled == nil || !*observed.OOMKilled) {
				t.Fatalf("OOM observations lost: %+v", observed)
			}
			verifySavedReport(t, options, report)
		})
	}
}

func TestCommandDiagnosticsRetainRedactedCause(t *testing.T) {
	for _, sensitive := range []bool{false, true} {
		t.Run(fmt.Sprint(sensitive), func(t *testing.T) {
			_, options := fixture(t)
			r, err := newRunner(options, func(context.Context, string, []string) (commandOutput, error) {
				return commandOutput{exitCode: -1, stdout: "resolved-config-must-not-leak"}, fmt.Errorf("cannot create process in %s: executable missing", options.Root)
			})
			if err != nil {
				t.Fatal(err)
			}
			_, err = r.command(context.Background(), "docker", []string{"compose", "config"}, time.Second, sensitive)
			if err == nil || !strings.Contains(err.Error(), "executable missing") || strings.Contains(err.Error(), options.Root) {
				t.Fatalf("cause lost or not redacted: %v", err)
			}
			if !strings.Contains(r.report.Commands[0].Error, "[workspace]") {
				t.Fatal("report lost redacted process cause")
			}
			if sensitive {
				if _, err = os.Stat(filepath.Join(r.directory, "command-1.log")); !errors.Is(err, os.ErrNotExist) {
					t.Fatal("sensitive output persisted")
				}
			}
		})
	}
}
```

<a id="go-smoke-config-tests"></a>

## Приложение: Shutdown budget tests

Полный исходный документ: `internal/smoke/config_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package smoke

import (
	"encoding/json"
	"testing"
)

func resolvedConfig(app, grace string) []byte {
	config := map[string]any{"services": map[string]any{}}
	for _, name := range services {
		config["services"].(map[string]any)[name] = map[string]any{"environment": map[string]string{"FL_SHUTDOWN_TIMEOUT": app}, "stop_grace_period": grace}
	}
	data, _ := json.Marshal(config)
	return data
}

func TestShutdownBudgets(t *testing.T) {
	for _, tc := range []struct {
		name, app, grace string
		allow, valid     bool
		warnings         int
	}{
		{"defaults", "10s", "15s", false, true, 0},
		{"minimum-margin", "10s", "12s", false, true, 0},
		{"external-too-short", "20s", "15s", false, false, 0},
		{"no-margin", "15s", "15s", false, false, 0},
		{"explicit-experiment", "20s", "15s", true, true, 3},
		{"invalid-app", "secret-value", "15s", true, false, 0},
		{"invalid-grace", "10s", "invalid", true, false, 0},
		{"zero-app", "0s", "15s", false, false, 0},
		{"negative-grace", "10s", "-1s", false, false, 0},
		{"overflow-margin", "2562047h47m16.854775807s", "2562047h47m16.854775807s", false, false, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			warnings, err := validateConfig(resolvedConfig(tc.app, tc.grace), tc.allow)
			if (err == nil) != tc.valid || len(warnings) != tc.warnings {
				t.Fatalf("warnings=%v error=%v", warnings, err)
			}
		})
	}
	for _, data := range []string{`not-json`, `{"services":{}}`} {
		if _, err := validateConfig([]byte(data), false); err == nil {
			t.Fatal("invalid/missing services accepted")
		}
	}
}
```

<a id="go-smoke-exec-tests"></a>

## Приложение: Command tree cancellation tests

Полный исходный документ: `internal/command/exec_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package command

import (
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCommandHelper(t *testing.T) {
	if os.Getenv("FAULTLINE_EXEC_HELPER") != "1" {
		return
	}
	if os.Getenv("FAULTLINE_EXEC_CHILD") == "1" {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			os.Exit(3)
		}
		if err := os.WriteFile(os.Getenv("FAULTLINE_EXEC_READY"), []byte(listener.Addr().String()), 0o600); err != nil {
			os.Exit(3)
		}
		connection, err := listener.Accept()
		if err != nil {
			os.Exit(4)
		}
		_, _ = io.Copy(io.Discard, connection)
		_ = connection.Close()
		_ = listener.Close()
		os.Exit(0)
	}
	child := exec.Command(os.Args[0], "-test.run=^TestCommandHelper$")
	child.Env = append(os.Environ(), "FAULTLINE_EXEC_CHILD=1")
	child.Stdout, child.Stderr = os.Stdout, os.Stderr
	if os.Getenv("FAULTLINE_EXEC_PARENT_EXITS") == "1" {
		child.Stdout, child.Stderr = nil, nil
	}
	if err := child.Start(); err != nil {
		os.Exit(2)
	}
	if os.Getenv("FAULTLINE_EXEC_PARENT_EXITS") == "1" {
		for {
			if _, err := os.Stat(os.Getenv("FAULTLINE_EXEC_RELEASE")); err == nil {
				os.Exit(0)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	if err := child.Wait(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func TestExecutorCancelsPluginTree(t *testing.T) {
	ready := filepath.Join(t.TempDir(), "ready")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	executor := processExecutor{directory: t.TempDir(), env: append(os.Environ(), "FAULTLINE_EXEC_HELPER=1", "FAULTLINE_EXEC_READY="+ready)}
	done := make(chan error, 1)
	go func() {
		_, err := executor.execute(ctx, os.Args[0], []string{"-test.run=^TestCommandHelper$"})
		done <- err
	}()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	var address []byte
	for {
		if data, err := os.ReadFile(ready); err == nil && len(data) > 0 {
			address = data
			break
		}
		select {
		case err := <-done:
			t.Fatalf("helper exited before child readiness: %v", err)
		case <-ctx.Done():
			t.Fatal("child readiness timed out")
		case <-ticker.C:
		}
	}
	connection, err := net.DialTimeout("tcp", string(address), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = connection.Close() }()
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("cancelled command succeeded")
		}
	case <-time.After(7 * time.Second):
		t.Fatal("command tree did not settle")
	}
	// A living plugin retains this socket even after its parent was killed.
	if err = connection.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	var value [1]byte
	_, err = connection.Read(value[:])
	if networkError, ok := err.(net.Error); err == nil || (ok && networkError.Timeout()) {
		t.Fatalf("plugin survived command cancellation: %v", err)
	}
}

func TestExecutorMissingProgram(t *testing.T) {
	executor := processExecutor{directory: t.TempDir(), env: os.Environ()}
	output, err := executor.execute(context.Background(), "faultline-nonexistent-test-command", nil)
	if err == nil || output.exitCode != -1 {
		t.Fatalf("output=%+v err=%v", output, err)
	}
}

func TestOutputCaptureIsBounded(t *testing.T) {
	buffer := &boundedBuffer{}
	data := make([]byte, outputLimit+123)
	if count, err := buffer.Write(data); err != nil || count != len(data) || buffer.Len() != outputLimit || !buffer.truncated {
		t.Fatalf("unbounded output: n=%d err=%v len=%d", count, err, buffer.Len())
	}
	buffer = &boundedBuffer{}
	reader := struct{ io.Reader }{strings.NewReader(strings.Repeat("x", outputLimit+123))}
	if count, err := io.Copy(buffer, reader); err != nil || count != outputLimit+123 || buffer.Len() != outputLimit || !buffer.truncated {
		t.Fatalf("io.Copy bypassed limit: n=%d err=%v len=%d", count, err, buffer.Len())
	}
}
```

<a id="doc-path-hygiene"></a>

## Приложение: Documentation path hygiene

Полный исходный документ: `tools/doc_safety.mjs`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```javascript
// Deliberately narrow hygiene check, not a general secret or filesystem scanner.
export function containsMachinePath(text) {
  const withoutWebURLs = text.replace(/\bhttps?:\/\/[^\s"'<>`]+/gi, '');
  return /\b[A-Za-z]:[\\/]+Users[\\/]+/i.test(withoutWebURLs)
    || /(?:^|[^A-Za-z0-9_.-])\/(?:Users|home)\/[^/\s"'<>()[\]{}]+/m.test(withoutWebURLs);
}
```

<a id="doc-path-hygiene-tests"></a>

## Приложение: Path hygiene examples

Полный исходный документ: `tools/doc_safety.test.mjs`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```javascript
import assert from 'node:assert/strict';
import test from 'node:test';
import { containsMachinePath } from './doc_safety.mjs';

test('detects synthetic Windows, macOS and Linux homes', () => {
  const windows = ['C:', 'Users', 'example', 'private.txt'].join('\\');
  const mac = ['', 'Users', 'example', 'private.txt'].join('/');
  const linux = ['', 'home', 'example', 'private.txt'].join('/');
  for (const value of [windows, JSON.stringify({ path: windows }), windows.replaceAll('\\', '/'), mac, linux, `path="${mac}"`, `file://${linux}`, linux.slice(0, linux.lastIndexOf('/'))]) {
    assert.equal(containsMachinePath(value), true, value);
  }
});

test('allows container paths, relative paths, regular URLs and ordinary prose', () => {
  for (const value of ['/app/service', '/workspace', '/root/.cache/go-build', 'docs/Users/guide.md', 'src/home/example', 'https://example.test/home/example/report', 'https://example.test/Users/example', 'Users and home directories']) {
    assert.equal(containsMachinePath(value), false, value);
  }
});
```

<a id="go-smoke-cli-interrupt"></a>

## Приложение: Host CLI SIGINT regression

Полный исходный документ: `tests/smokecli/interrupt_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
//go:build !windows

// Package smokecli runs opt-in host integration tests of the compiled CLI.
package smokecli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/smoke"
)

func TestSIGINTCleansRunningProject(t *testing.T) {
	binary, root := os.Getenv("LAB_SMOKE_BINARY"), os.Getenv("LAB_SMOKE_ROOT")
	if binary == "" || root == "" {
		t.Skip("host integration requires LAB_SMOKE_BINARY and LAB_SMOKE_ROOT; CI runs the compiled test on the host")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, "-root", root)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	projects, signals, scanned := make(chan string, 1), make(chan error, 1), make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasSuffix(line, ": prerequisites") {
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					projects <- strings.TrimSuffix(fields[2], ":")
				}
			}
			if strings.HasSuffix(line, ": check") {
				signals <- cmd.Process.Signal(os.Interrupt)
			}
		}
		scanned <- scanner.Err()
	}()
	scanErr := <-scanned
	waitErr := cmd.Wait()
	if scanErr != nil {
		t.Logf("CLI stdout ended with: %v", scanErr)
	}
	var project string
	select {
	case project = <-projects:
	default:
		t.Fatalf("CLI never identified its project: %v; %s", waitErr, stderr.String())
	}
	// Even a failed regression must leave no resources owned by this test.
	t.Cleanup(func() {
		cleanupCtx, stop := context.WithTimeout(context.Background(), time.Minute)
		defer stop()
		cleanup := exec.CommandContext(cleanupCtx, "docker", "compose", "-f", filepath.Join(root, "compose.yaml"), "-p", project, "down", "--remove-orphans", "--rmi", "local")
		if output, err := cleanup.CombinedOutput(); err != nil {
			t.Errorf("test recovery failed: %v; %s", err, output)
		}
	})
	select {
	case err := <-signals:
		if err != nil {
			t.Fatal(err)
		}
	default:
		t.Fatalf("CLI did not reach interruption phase: %v; %s", waitErr, stderr.String())
	}
	var exit *exec.ExitError
	if !errors.As(waitErr, &exit) || exit.ExitCode() != 1 {
		t.Fatalf("unexpected interrupt exit: %v; %s", waitErr, stderr.String())
	}
	paths, err := filepath.Glob(filepath.Join(root, "artifacts", "smoke", "*", "report.json"))
	if err != nil {
		t.Fatal(err)
	}
	var observed *smoke.Report
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var report smoke.Report
		if err = json.Unmarshal(data, &report); err != nil {
			t.Fatal(err)
		}
		if report.Project == project {
			observed = &report
		}
	}
	if observed == nil || !observed.Cancelled || observed.Cleanup != "PASS" || observed.Result != "FAIL" || observed.FinishedAt == nil {
		t.Fatalf("missing or invalid cancellation report: %+v", observed)
	}
	for _, args := range [][]string{
		{"ps", "--all", "--quiet", "--filter", "label=com.docker.compose.project=" + project},
		{"network", "ls", "--quiet", "--filter", "label=com.docker.compose.project=" + project},
		{"image", "ls", "--quiet", "--filter", "reference=" + project + "-*:latest"},
	} {
		output, err := exec.CommandContext(ctx, "docker", args...).Output()
		if err != nil || strings.TrimSpace(string(output)) != "" {
			t.Fatalf("owned resources remain: %s; %v", output, err)
		}
	}
}
```

<a id="source-smoke-exec-unix-test"></a>

## Приложение: Unix orphan plugin lifecycle regression

Полный исходный документ: `internal/command/exec_unix_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
//go:build !windows

package command

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExecutorStopsPluginAfterSuccessfulParent(t *testing.T) {
	directory := t.TempDir()
	ready, release := filepath.Join(directory, "ready"), filepath.Join(directory, "release")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	executor := processExecutor{directory: directory, env: append(os.Environ(),
		"FAULTLINE_EXEC_HELPER=1", "FAULTLINE_EXEC_PARENT_EXITS=1",
		"FAULTLINE_EXEC_READY="+ready, "FAULTLINE_EXEC_RELEASE="+release)}
	done := make(chan error, 1)
	go func() {
		_, err := executor.execute(ctx, os.Args[0], []string{"-test.run=^TestCommandHelper$"})
		done <- err
	}()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	var address []byte
	for {
		if data, err := os.ReadFile(ready); err == nil && len(data) > 0 {
			address = data
			break
		}
		select {
		case err := <-done:
			t.Fatalf("helper exited before child readiness: %v", err)
		case <-ctx.Done():
			t.Fatal("child readiness timed out")
		case <-ticker.C:
		}
	}
	connection, err := net.DialTimeout("tcp", string(address), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = connection.Close() }()
	if err = os.WriteFile(release, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	select {
	case err = <-done:
		if err != nil {
			t.Fatalf("successful parent: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("executor retained an orphan plugin")
	}
	if err = connection.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	var value [1]byte
	_, err = connection.Read(value[:])
	if networkError, ok := err.(net.Error); err == nil || (ok && networkError.Timeout()) {
		t.Fatalf("plugin survived successful parent: %v", err)
	}
}
```

<a id="go-smoke-evidence-go-equivalence-report"></a>

## Приложение: Go smoke evidence: go-equivalence-report.json

Полный исходный документ: `docs/development/evidence/go-smoke/go-equivalence-report.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "schema_version": 2,
  "run_id": "2c542d76-0c2a-d94f-8058-90277740c358",
  "project": "faultline-smoke-2c542d76",
  "started_at": "2026-10-16T10:25:24.2558023Z",
  "finished_at": "2026-10-16T10:26:53.8025839Z",
  "result": "PASS",
  "phase": "stop",
  "cancelled": false,
  "cleanup": "PASS",
  "host_platform": "windows",
  "go_version": "go1.27.2",
  "docker_version": "29.8.2",
  "docker_architecture": "linux/x86_64",
  "compose_version": "5.6.0",
  "source_revision": "511b9a85c7006cc7042263e0b87e7eb7fb8b116c",
  "source_dirty": true,
  "isolation_after_orders_stop": "PASS",
  "services": {
    "inventory": {
      "id": "1693539fb05a68bb4fe8367af0b5c7ecf909c426a69c72f6eb357eb2109170e5",
      "image_id": "sha256:1fe8e5bf09a0afa078cd254be22e782717995183c5ce4570ca50db37e48223d4",
      "healthy": true,
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "status": "exited",
      "stop_mode": "graceful",
      "security": {
        "capabilities_dropped": true,
        "host_ports": [],
        "internal_network": true,
        "mounts": [],
        "no_new_privileges": true,
        "read_only": true,
        "user": "65532:65532"
      }
    },
    "orders": {
      "id": "9a58597ce51b32ff8a88470b41730c0abadf895ed842ff4854c91ad32f31529f",
      "image_id": "sha256:9a1ba3867da05da489dc28fd48e8f7c7e3fc1507085a720275e2d0d9bc423ec6",
      "healthy": true,
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "status": "exited",
      "stop_mode": "graceful",
      "security": {
        "capabilities_dropped": true,
        "host_ports": [],
        "internal_network": true,
        "mounts": [],
        "no_new_privileges": true,
        "read_only": true,
        "user": "65532:65532"
      }
    },
    "payments": {
      "id": "0b439bfe5f9413d6a426068d79ba0ed1edb560221b100c2086d80dbcda219467",
      "image_id": "sha256:31f1150261d44cf46ffed9779eae20bc6b190c58593d50711a7c8450f3bdb916",
      "healthy": true,
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "status": "exited",
      "stop_mode": "graceful",
      "security": {
        "capabilities_dropped": true,
        "host_ports": [],
        "internal_network": true,
        "mounts": [],
        "no_new_privileges": true,
        "read_only": true,
        "user": "65532:65532"
      }
    }
  },
  "commands": [
    {
      "program": "docker",
      "args": [
        "version",
        "--format",
        "{{.Server.Version}}"
      ],
      "started_at": "2026-10-16T10:25:24.2586247Z",
      "finished_at": "2026-10-16T10:25:24.5590754Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "info",
        "--format",
        "{{.OSType}}/{{.Architecture}}"
      ],
      "started_at": "2026-10-16T10:25:24.5608582Z",
      "finished_at": "2026-10-16T10:25:25.4423204Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "version",
        "--short"
      ],
      "started_at": "2026-10-16T10:25:25.4475039Z",
      "finished_at": "2026-10-16T10:25:25.8784425Z",
      "exit_code": 0
    },
    {
      "program": "git",
      "args": [
        "rev-parse",
        "HEAD"
      ],
      "started_at": "2026-10-16T10:25:25.8796159Z",
      "finished_at": "2026-10-16T10:25:25.9525844Z",
      "exit_code": 0
    },
    {
      "program": "git",
      "args": [
        "status",
        "--porcelain"
      ],
      "started_at": "2026-10-16T10:25:25.9536769Z",
      "finished_at": "2026-10-16T10:25:26.0467329Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "config",
        "--format",
        "json"
      ],
      "started_at": "2026-10-16T10:25:26.0486798Z",
      "finished_at": "2026-10-16T10:25:26.5904105Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "build",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T10:25:26.5940295Z",
      "finished_at": "2026-10-16T10:26:30.2916155Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "up",
        "-d",
        "--wait",
        "--wait-timeout",
        "60",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T10:26:30.2944158Z",
      "finished_at": "2026-10-16T10:26:36.4737484Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "exec",
        "-T",
        "orders",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:26:36.4758897Z",
      "finished_at": "2026-10-16T10:26:37.0629287Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T10:26:37.0645474Z",
      "finished_at": "2026-10-16T10:26:37.5982237Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "9a58597ce51b32ff8a88470b41730c0abadf895ed842ff4854c91ad32f31529f"
      ],
      "started_at": "2026-10-16T10:26:37.5997991Z",
      "finished_at": "2026-10-16T10:26:37.8801755Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "86a6ffea1c9d19df885e8352b93f30451b2f502ba574f193be06487d2e96274d"
      ],
      "started_at": "2026-10-16T10:26:37.8812956Z",
      "finished_at": "2026-10-16T10:26:38.1606769Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:26:38.162372Z",
      "finished_at": "2026-10-16T10:26:38.7176138Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T10:26:38.7193256Z",
      "finished_at": "2026-10-16T10:26:39.3576591Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "1693539fb05a68bb4fe8367af0b5c7ecf909c426a69c72f6eb357eb2109170e5"
      ],
      "started_at": "2026-10-16T10:26:39.3592638Z",
      "finished_at": "2026-10-16T10:26:39.6367909Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "86a6ffea1c9d19df885e8352b93f30451b2f502ba574f193be06487d2e96274d"
      ],
      "started_at": "2026-10-16T10:26:39.6378901Z",
      "finished_at": "2026-10-16T10:26:39.890402Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:26:39.8925602Z",
      "finished_at": "2026-10-16T10:26:40.4496294Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T10:26:40.450692Z",
      "finished_at": "2026-10-16T10:26:41.0533871Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "0b439bfe5f9413d6a426068d79ba0ed1edb560221b100c2086d80dbcda219467"
      ],
      "started_at": "2026-10-16T10:26:41.0554864Z",
      "finished_at": "2026-10-16T10:26:41.3146838Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "86a6ffea1c9d19df885e8352b93f30451b2f502ba574f193be06487d2e96274d"
      ],
      "started_at": "2026-10-16T10:26:41.3162477Z",
      "finished_at": "2026-10-16T10:26:41.5515734Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "stop",
        "orders"
      ],
      "started_at": "2026-10-16T10:26:41.5543346Z",
      "finished_at": "2026-10-16T10:26:42.2502596Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T10:26:42.2514313Z",
      "finished_at": "2026-10-16T10:26:42.8240457Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "9a58597ce51b32ff8a88470b41730c0abadf895ed842ff4854c91ad32f31529f"
      ],
      "started_at": "2026-10-16T10:26:42.8258869Z",
      "finished_at": "2026-10-16T10:26:43.0796172Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "orders"
      ],
      "started_at": "2026-10-16T10:26:43.0807302Z",
      "finished_at": "2026-10-16T10:26:43.6377388Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:26:43.6397745Z",
      "finished_at": "2026-10-16T10:26:44.2714292Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T10:26:44.2731569Z",
      "finished_at": "2026-10-16T10:26:44.8315476Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "1693539fb05a68bb4fe8367af0b5c7ecf909c426a69c72f6eb357eb2109170e5"
      ],
      "started_at": "2026-10-16T10:26:44.8326319Z",
      "finished_at": "2026-10-16T10:26:45.0676528Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "86a6ffea1c9d19df885e8352b93f30451b2f502ba574f193be06487d2e96274d"
      ],
      "started_at": "2026-10-16T10:26:45.0696814Z",
      "finished_at": "2026-10-16T10:26:45.3028022Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:26:45.3038982Z",
      "finished_at": "2026-10-16T10:26:45.8592771Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T10:26:45.8604517Z",
      "finished_at": "2026-10-16T10:26:46.3761683Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "0b439bfe5f9413d6a426068d79ba0ed1edb560221b100c2086d80dbcda219467"
      ],
      "started_at": "2026-10-16T10:26:46.3772369Z",
      "finished_at": "2026-10-16T10:26:46.6072718Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "86a6ffea1c9d19df885e8352b93f30451b2f502ba574f193be06487d2e96274d"
      ],
      "started_at": "2026-10-16T10:26:46.6087618Z",
      "finished_at": "2026-10-16T10:26:46.8397519Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "stop",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T10:26:46.843144Z",
      "finished_at": "2026-10-16T10:26:47.5621048Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T10:26:47.5638759Z",
      "finished_at": "2026-10-16T10:26:48.1632927Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "1693539fb05a68bb4fe8367af0b5c7ecf909c426a69c72f6eb357eb2109170e5"
      ],
      "started_at": "2026-10-16T10:26:48.1650802Z",
      "finished_at": "2026-10-16T10:26:48.3984032Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "inventory"
      ],
      "started_at": "2026-10-16T10:26:48.4006521Z",
      "finished_at": "2026-10-16T10:26:48.9538034Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T10:26:48.9549741Z",
      "finished_at": "2026-10-16T10:26:49.4914866Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "0b439bfe5f9413d6a426068d79ba0ed1edb560221b100c2086d80dbcda219467"
      ],
      "started_at": "2026-10-16T10:26:49.4928709Z",
      "finished_at": "2026-10-16T10:26:49.723309Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "payments"
      ],
      "started_at": "2026-10-16T10:26:49.7243684Z",
      "finished_at": "2026-10-16T10:26:50.2768365Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "orders"
      ],
      "started_at": "2026-10-16T10:26:50.2803204Z",
      "finished_at": "2026-10-16T10:26:50.8539136Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "inventory"
      ],
      "started_at": "2026-10-16T10:26:50.8590591Z",
      "finished_at": "2026-10-16T10:26:51.4117227Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "payments"
      ],
      "started_at": "2026-10-16T10:26:51.4133432Z",
      "finished_at": "2026-10-16T10:26:51.9819447Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-2c542d76",
        "down",
        "--remove-orphans"
      ],
      "started_at": "2026-10-16T10:26:51.984209Z",
      "finished_at": "2026-10-16T10:26:52.6648361Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "ps",
        "--all",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-2c542d76"
      ],
      "started_at": "2026-10-16T10:26:52.6658416Z",
      "finished_at": "2026-10-16T10:26:52.9020802Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "ls",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-2c542d76"
      ],
      "started_at": "2026-10-16T10:26:52.9049175Z",
      "finished_at": "2026-10-16T10:26:53.1388709Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "ps",
        "--all",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-2c542d76"
      ],
      "started_at": "2026-10-16T10:26:53.3417179Z",
      "finished_at": "2026-10-16T10:26:53.571252Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "ls",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-2c542d76"
      ],
      "started_at": "2026-10-16T10:26:53.5723396Z",
      "finished_at": "2026-10-16T10:26:53.8016341Z",
      "exit_code": 0
    }
  ],
  "recovery_command": [
    "docker",
    "compose",
    "-f",
    "compose.yaml",
    "-p",
    "faultline-smoke-2c542d76",
    "down",
    "--remove-orphans"
  ],
  "limitations": [
    "No business operations or capacity measurement",
    "docker_architecture describes the backend, not individual image platforms (F03)",
    "SIGKILL or host failure requires recovery_command; signal cleanup cannot be guaranteed in those cases"
  ]
}
```

<a id="go-smoke-evidence-go-smoke-equivalence"></a>

## Приложение: Go smoke evidence: go-smoke-equivalence.json

Полный исходный документ: `docs/development/evidence/go-smoke/go-smoke-equivalence.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "date": "2026-10-16",
  "platform": "windows/amd64",
  "legacy_run_id": "8e8a1be6-fb98-45c2-906f-ac267292eb74",
  "go_run_id": "2c542d76-0c2a-d94f-8058-90277740c358",
  "legacy_result": "PASS",
  "go_result": "PASS",
  "legacy_cleanup": "PASS",
  "go_cleanup": "PASS",
  "health_stop_exit_oom_mode_equal": true,
  "same_image_ids": false,
  "source_revision": "511b9a85c7006cc7042263e0b87e7eb7fb8b116c",
  "source_dirty": true,
  "scope": "Functional equivalence, not performance comparison"
}
```

<a id="go-smoke-evidence-go-smoke-preflight"></a>

## Приложение: Go smoke evidence: go-smoke-preflight.json

Полный исходный документ: `docs/development/evidence/go-smoke/go-smoke-preflight.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "date": "2026-10-16",
  "app_timeout": "20s",
  "invalid_container_grace": "15s",
  "invalid_exit": 1,
  "explicit_unsafe_exit": 0,
  "valid_container_grace": "25s",
  "valid_exit": 0,
  "services_started": false
}
```

<a id="go-smoke-evidence-go-smoke-windows-interrupt"></a>

## Приложение: Go smoke evidence: go-smoke-windows-interrupt.json

Полный исходный документ: `docs/development/evidence/go-smoke/go-smoke-windows-interrupt.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "date": "2026-10-16",
  "platform": "windows/amd64",
  "signal": "CTRL_BREAK_EVENT",
  "interrupted_phase": "check",
  "report_run_id": "0d41daa4-2151-8b9e-8268-0a597b5164c0",
  "result": "FAIL",
  "cancelled": true,
  "cleanup": "PASS",
  "owned_containers_remaining": 0,
  "owned_networks_remaining": 0,
  "independent_control_preserved": true,
  "source_revision": "511b9a85c7006cc7042263e0b87e7eb7fb8b116c",
  "source_dirty": true
}
```

<a id="go-smoke-evidence-no-backend-report"></a>

## Приложение: Go smoke evidence: no-backend-report.json

Полный исходный документ: `docs/development/evidence/go-smoke/no-backend-report.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "schema_version": 2,
  "run_id": "9dca928d-11b5-79c0-e9bc-81ea0e3e0475",
  "project": "faultline-smoke-9dca928d",
  "started_at": "2026-10-16T11:05:21.8245087Z",
  "finished_at": "2026-10-16T11:05:22.0645715Z",
  "result": "NOT_RUN",
  "phase": "prerequisites",
  "failed_phase": "prerequisites",
  "error": "docker version --format {{.Server.Version}} failed (exit 1; process exited with code 1)",
  "cancelled": false,
  "cleanup": "NOT_NEEDED",
  "host_platform": "windows",
  "go_version": "go1.27.2",
  "source_dirty": false,
  "services": {},
  "commands": [
    {
      "program": "docker",
      "args": [
        "version",
        "--format",
        "{{.Server.Version}}"
      ],
      "started_at": "2026-10-16T11:05:21.8263028Z",
      "finished_at": "2026-10-16T11:05:22.0620351Z",
      "exit_code": 1,
      "error": "process exited with code 1"
    }
  ],
  "recovery_command": [
    "docker",
    "compose",
    "-f",
    "compose.yaml",
    "-p",
    "faultline-smoke-9dca928d",
    "down",
    "--remove-orphans"
  ],
  "limitations": [
    "No business operations or capacity measurement",
    "docker_architecture describes the backend, not individual image platforms (F03)",
    "SIGKILL or host failure requires recovery_command; signal cleanup cannot be guaranteed in those cases"
  ]
}
```

<a id="go-smoke-evidence-node-baseline-report"></a>

## Приложение: Go smoke evidence: node-baseline-report.json

Полный исходный документ: `docs/development/evidence/go-smoke/node-baseline-report.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "run_id": "8e8a1be6-fb98-45c2-906f-ac267292eb74",
  "started_at": "2026-10-16T10:21:13.626Z",
  "result": "PASS",
  "host_platform": "win32",
  "node": "v24.19.0",
  "project": "faultline-smoke-8e8a1be6",
  "services": {
    "orders": {
      "id": "eab66dbca7d7258fe172b1f9dab449675f8f2f4cdb9a5895359f2a9717475fe4",
      "image_id": "sha256:ded9ced431e0e1d36d42176fd3db15a6ab41888d2dd25f7135e773d323042fd1",
      "architecture": "linux/x86_64",
      "healthy": true,
      "security": {
        "user": "65532:65532",
        "read_only": true,
        "capabilities_dropped": true,
        "no_new_privileges": true,
        "host_ports": [],
        "mounts": [],
        "internal_network": true
      },
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "stop_mode": "graceful"
    },
    "inventory": {
      "id": "2c963ec0e72eafd83485056026435505a8cac4017f6bfd0b3d2f8a0fa71b40c2",
      "image_id": "sha256:c54585b2eeb99651e465809010e8431b0cf7c397edf8c66f01405fa4a8a9dabd",
      "architecture": "linux/x86_64",
      "healthy": true,
      "security": {
        "user": "65532:65532",
        "read_only": true,
        "capabilities_dropped": true,
        "no_new_privileges": true,
        "host_ports": [],
        "mounts": [],
        "internal_network": true
      },
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "stop_mode": "graceful"
    },
    "payments": {
      "id": "3857972f28e8210fb75c052c292b9ece1db471fa87d7ab48a9ed999130dc408d",
      "image_id": "sha256:55af9ccea8d6ea6a94463b651a0dbe070b55841b048ec00a9f960ed5af0e84ef",
      "architecture": "linux/x86_64",
      "healthy": true,
      "security": {
        "user": "65532:65532",
        "read_only": true,
        "capabilities_dropped": true,
        "no_new_privileges": true,
        "host_ports": [],
        "mounts": [],
        "internal_network": true
      },
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "stop_mode": "graceful"
    }
  },
  "commands": [
    {
      "program": "docker",
      "args": [
        "version",
        "--format",
        "{{.Server.Version}}"
      ],
      "started_at": "2026-10-16T10:21:13.629Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:21:13.934Z"
    },
    {
      "program": "docker",
      "args": [
        "info",
        "--format",
        "{{.OSType}}/{{.Architecture}}"
      ],
      "started_at": "2026-10-16T10:21:13.934Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:21:16.251Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "version",
        "--short"
      ],
      "started_at": "2026-10-16T10:21:16.251Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:21:16.666Z"
    },
    {
      "program": "git",
      "args": [
        "rev-parse",
        "HEAD"
      ],
      "started_at": "2026-10-16T10:21:16.667Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:21:16.808Z"
    },
    {
      "program": "git",
      "args": [
        "status",
        "--porcelain"
      ],
      "started_at": "2026-10-16T10:21:16.808Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:21:16.910Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "config",
        "--quiet"
      ],
      "started_at": "2026-10-16T10:21:16.910Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:21:17.629Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "build",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T10:21:17.629Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:44.002Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "up",
        "-d",
        "--wait",
        "--wait-timeout",
        "60",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T10:22:44.003Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:50.231Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "exec",
        "-T",
        "orders",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:22:50.232Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:50.802Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T10:22:50.802Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:51.310Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "eab66dbca7d7258fe172b1f9dab449675f8f2f4cdb9a5895359f2a9717475fe4"
      ],
      "started_at": "2026-10-16T10:22:51.311Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:51.534Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "677de20f67e81ad9506ca619e0b0290d65c7dae7f68264683b2e36debca5e8c3"
      ],
      "started_at": "2026-10-16T10:22:51.534Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:51.758Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:22:51.758Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:52.295Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T10:22:52.295Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:52.810Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "2c963ec0e72eafd83485056026435505a8cac4017f6bfd0b3d2f8a0fa71b40c2"
      ],
      "started_at": "2026-10-16T10:22:52.810Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:53.035Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "677de20f67e81ad9506ca619e0b0290d65c7dae7f68264683b2e36debca5e8c3"
      ],
      "started_at": "2026-10-16T10:22:53.036Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:53.254Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:22:53.254Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:53.808Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T10:22:53.808Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:54.318Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "3857972f28e8210fb75c052c292b9ece1db471fa87d7ab48a9ed999130dc408d"
      ],
      "started_at": "2026-10-16T10:22:54.318Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:54.544Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "677de20f67e81ad9506ca619e0b0290d65c7dae7f68264683b2e36debca5e8c3"
      ],
      "started_at": "2026-10-16T10:22:54.544Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:54.765Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "stop",
        "orders"
      ],
      "started_at": "2026-10-16T10:22:54.765Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:55.399Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T10:22:55.399Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:55.926Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "eab66dbca7d7258fe172b1f9dab449675f8f2f4cdb9a5895359f2a9717475fe4"
      ],
      "started_at": "2026-10-16T10:22:55.926Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:56.144Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T10:22:56.144Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:56.653Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "eab66dbca7d7258fe172b1f9dab449675f8f2f4cdb9a5895359f2a9717475fe4"
      ],
      "started_at": "2026-10-16T10:22:56.653Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:56.882Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:22:56.885Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:57.422Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T10:22:57.422Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:57.948Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "2c963ec0e72eafd83485056026435505a8cac4017f6bfd0b3d2f8a0fa71b40c2"
      ],
      "started_at": "2026-10-16T10:22:57.948Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:58.174Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "677de20f67e81ad9506ca619e0b0290d65c7dae7f68264683b2e36debca5e8c3"
      ],
      "started_at": "2026-10-16T10:22:58.174Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:58.393Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:22:58.393Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:58.947Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T10:22:58.947Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:59.449Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "3857972f28e8210fb75c052c292b9ece1db471fa87d7ab48a9ed999130dc408d"
      ],
      "started_at": "2026-10-16T10:22:59.450Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:59.674Z"
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "677de20f67e81ad9506ca619e0b0290d65c7dae7f68264683b2e36debca5e8c3"
      ],
      "started_at": "2026-10-16T10:22:59.674Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:22:59.900Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "stop",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T10:22:59.900Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:00.551Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T10:23:00.551Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:01.122Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "2c963ec0e72eafd83485056026435505a8cac4017f6bfd0b3d2f8a0fa71b40c2"
      ],
      "started_at": "2026-10-16T10:23:01.122Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:01.341Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T10:23:01.342Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:01.848Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "2c963ec0e72eafd83485056026435505a8cac4017f6bfd0b3d2f8a0fa71b40c2"
      ],
      "started_at": "2026-10-16T10:23:01.848Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:02.075Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T10:23:02.077Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:02.587Z"
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "3857972f28e8210fb75c052c292b9ece1db471fa87d7ab48a9ed999130dc408d"
      ],
      "started_at": "2026-10-16T10:23:02.587Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:02.811Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T10:23:02.811Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:03.338Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "3857972f28e8210fb75c052c292b9ece1db471fa87d7ab48a9ed999130dc408d"
      ],
      "started_at": "2026-10-16T10:23:03.338Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:03.565Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T10:23:03.566Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:04.101Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "eab66dbca7d7258fe172b1f9dab449675f8f2f4cdb9a5895359f2a9717475fe4"
      ],
      "started_at": "2026-10-16T10:23:04.101Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:04.331Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T10:23:04.331Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:04.849Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "2c963ec0e72eafd83485056026435505a8cac4017f6bfd0b3d2f8a0fa71b40c2"
      ],
      "started_at": "2026-10-16T10:23:04.849Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:05.074Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T10:23:05.075Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:05.574Z"
    },
    {
      "program": "docker",
      "args": [
        "logs",
        "3857972f28e8210fb75c052c292b9ece1db471fa87d7ab48a9ed999130dc408d"
      ],
      "started_at": "2026-10-16T10:23:05.574Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:05.797Z"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-8e8a1be6",
        "down",
        "--remove-orphans"
      ],
      "started_at": "2026-10-16T10:23:05.797Z",
      "exit_code": 0,
      "finished_at": "2026-10-16T10:23:06.444Z"
    }
  ],
  "limitations": [
    "No business operations or capacity measurement"
  ],
  "docker_version": "29.8.2",
  "docker_architecture": "linux/x86_64",
  "compose_version": "5.6.0",
  "source_revision": "511b9a85c7006cc7042263e0b87e7eb7fb8b116c",
  "source_dirty": true,
  "isolation_after_orders_stop": "PASS",
  "cleanup": "PASS",
  "finished_at": "2026-10-16T10:23:06.444Z"
}
```

<a id="go-smoke-evidence-VALIDATION"></a>

## Приложение: Go smoke evidence: VALIDATION.md

Полный исходный документ: `docs/development/evidence/go-smoke/VALIDATION.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### Проверка Go smoke — 16 октября 2026 года

Датированные результаты для [#9](https://github.com/Includeoyasi/faultline-lab/issues/9) и F01/F02/F04/F07/F08 из [#20](https://github.com/Includeoyasi/faultline-lab/issues/20). Это доказательства конкретного изменения; актуальное выполнение и приёмка находятся в GitHub.

#### Среда и происхождение

Windows AMD64, Go 1.27.2; Docker Desktop 4.95, Linux backend 29.8.2/x86_64, Compose 5.6.0. Toolchain работает в контейнере. Ресурсы WSL ограничены двумя CPU и 2 GiB памяти; эти проверки не измеряют производительность бизнеса.

Исходники заключительного запуска: `8e49367c868af33023a6c1f7a8f94db447a828eb`, чистая рабочая копия. [Итоговый отчёт](#go-smoke-evidence-windows-final-report) содержит `source_dirty=false`, реальные container/image IDs, health/security, изоляцию Orders, graceful stop с exit 0, OOM=false и cleanup=PASS. Рядом сохранены JSONL [Orders](#go-smoke-evidence-windows-final-orders), [Inventory](#go-smoke-evidence-windows-final-inventory), [Payments](#go-smoke-evidence-windows-final-payments).

#### Фактические проверки

| Проверка | Результат и границы |
|---|---|
| `docker compose run --rm toolchain check` | PASS: format/imports, tidy/verify, build, lint, unit, race, govulncheck. Linux AMD64 внутри Docker Desktop. |
| `docker compose run --rm toolchain service-smoke` | PASS: процессы трёх сервисов, health, изоляция и остановка. Существующие Go-тесты подтверждают завершение активного RPC и ограниченный forced stop зависшего RPC/Watch. |
| `docker compose run --rm toolchain smoke-cross-check` | PASS: Windows AMD64, Linux AMD64/ARM64, macOS AMD64/ARM64. Это проверка сборки, а не runtime всех платформ. |
| Нативный Windows `go test -count=1 -timeout=90s ./internal/smoke ./cmd/smoke` | PASS, включая завершение дерева процессов через Job Object. |
| `bin/smoke-windows-amd64.exe` | PASS на чистом исходном коммите; все собственные контейнеры/сеть удалены. |
| Windows `CTRL_BREAK_EVENT` на этапе check | Ожидаемый FAIL, `cancelled=true`, cleanup=PASS, собственных ресурсов не осталось; независимый контрольный контейнер продолжил работу. [Наблюдения](#go-smoke-evidence-go-smoke-windows-interrupt), [отчёт](#go-smoke-evidence-windows-interrupt-report). |
| Отмена build/start/check/stop и ошибки команд/cleanup | PASS в Go fixtures с явной моделью Docker. Они не выдаются за четыре реальных OS-interrupt прогона. Проверены независимый cleanup context, deadline, позднее появление собственного ресурса и сохранение чужого. |
| Эффективный Compose-профиль | 20s/15s отклонён до запуска; 20s/25s принят; намеренный 20s/15s допустим только через явный unsafe preflight. [Наблюдения](#go-smoke-evidence-go-smoke-preflight). |
| Backend выключен | Ожидаемый NOT_RUN / exit 2 / cleanup=NOT_NEEDED, [отчёт](#go-smoke-evidence-no-backend-report). Новые сервисы не запускались. |
| Гигиена путей | `node --test tools/doc_safety.test.mjs`: PASS. Windows, JSON-escaped Windows, macOS/Linux homes отклоняются; контейнерные/относительные пути и HTTP URL допускаются. |
| Переносимый контекст | Генерация и `-CheckOnly` в PowerShell 7 и Windows PowerShell 5.1; изменяемый `snapshot_ref=main` отклонён. SHA исходных материалов проверяется через Git. |

Тест Unix process group отдельно проверяет потомка, удерживающего TCP-соединение после успешного выхода родителя с перенаправленным stdio. Executor завершает остаток своей группы. Ограничение буфера проверяется через обычный Write и `io.Copy`; причина ошибки executor сохраняется с редактированием домашних путей, sensitive config stdout/stderr не публикуются.

#### Сравнение с прежним runner

На одной Windows/Docker среде последовательно выполнены прежний Node runner и новый Go runner: оба PASS, cleanup PASS, одинаковые health/stop/exit/OOM/isolation результаты. [Сводка](#go-smoke-evidence-go-smoke-equivalence), [Node report](#go-smoke-evidence-node-baseline-report), [Go report](#go-smoke-evidence-go-equivalence-report).

Эти два запуска и прерывание выполнены при HEAD `511b9a85c7006cc7042263e0b87e7eb7fb8b116c` с незакоммиченными изменениями (`source_dirty=true`), до заключительного чистого запуска. Image IDs различаются после пересборки; идентичность образов или сравнительная производительность не заявляются. Старый runner удалён после проверки замены; исходная версия доступна в истории Git.

#### Ограничения и завершение локальной работы

Linux host smoke и настоящий SIGINT запускаются отдельными CI jobs для нативных AMD64/ARM64; их результаты и артефакты находятся в PR. Проверки выше не подменяют эти прогоны. macOS runtime — **NOT_RUN**, доступа к Mac в этом запуске не было. Полный контракт платформ F03 остаётся отдельным решением; поле backend не обозначает архитектуру каждого image.

После локальных проверок собственных контейнеров осталось 0, Docker Desktop выключен, Docker WSL VM и процессы разработки завершены. Автозапуск отключён; программы, образы и полезные кеши сохранены. [Подтверждение](#go-smoke-evidence-windows-cleanup), время наблюдения указано внутри. Проверка NOT_RUN после выключения backend не поднимает его снова.

Полный разрешённый Compose config и inspect Env не сохранялись. Публикуемые JSON/JSONL не содержат персональных путей; произвольные сторонние логи требуют отдельной проверки перед публикацией. Для SIGKILL/аварии хоста восстановление выполняется адресно по `recovery_command`, без обещания запуска signal handler.

<a id="go-smoke-evidence-windows-cleanup"></a>

## Приложение: Go smoke evidence: windows-cleanup.json

Полный исходный документ: `docs/development/evidence/go-smoke/windows-cleanup.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "observed_at": "2026-10-16T14:05:21.6991190+03:00",
  "containers_remaining": 0,
  "desktop_stop_exit": 0,
  "docker_vm_running": false,
  "remaining_development_processes": [],
  "docker_autostart": false,
  "docker_run_entry_present": false,
  "unrelated_workloads_preserved": true,
  "installed_tools_preserved": true,
  "images_and_build_cache_preserved": true,
  "reboot_performed": false
}
```

<a id="go-smoke-evidence-windows-final-inventory"></a>

## Приложение: Go smoke evidence: windows-final-inventory.jsonl

Полный исходный документ: `docs/development/evidence/go-smoke/windows-final-inventory.jsonl`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```jsonl
{"time":"2026-10-16T11:02:26.441897119Z","level":"INFO","msg":"service started","service":"inventory","version":"(devel)","revision":"unknown","vcs_modified":"unknown","go_version":"go1.27.2","grpc_address":"[::]:50051"}
{"time":"2026-10-16T11:02:42.489483535Z","level":"INFO","msg":"service stopping","service":"inventory","health":"NOT_SERVING","shutdown_timeout":"10s"}
{"time":"2026-10-16T11:02:42.490028932Z","level":"INFO","msg":"service stopped","service":"inventory","stop_mode":"graceful"}
```

<a id="go-smoke-evidence-windows-final-orders"></a>

## Приложение: Go smoke evidence: windows-final-orders.jsonl

Полный исходный документ: `docs/development/evidence/go-smoke/windows-final-orders.jsonl`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```jsonl
{"time":"2026-10-16T11:02:26.566500793Z","level":"INFO","msg":"service started","service":"orders","version":"(devel)","revision":"unknown","vcs_modified":"unknown","go_version":"go1.27.2","grpc_address":"[::]:50051"}
{"time":"2026-10-16T11:02:37.303113807Z","level":"INFO","msg":"service stopping","service":"orders","health":"NOT_SERVING","shutdown_timeout":"10s"}
{"time":"2026-10-16T11:02:37.303206973Z","level":"INFO","msg":"service stopped","service":"orders","stop_mode":"graceful"}
```

<a id="go-smoke-evidence-windows-final-payments"></a>

## Приложение: Go smoke evidence: windows-final-payments.jsonl

Полный исходный документ: `docs/development/evidence/go-smoke/windows-final-payments.jsonl`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```jsonl
{"time":"2026-10-16T11:02:26.705166656Z","level":"INFO","msg":"service started","service":"payments","version":"(devel)","revision":"unknown","vcs_modified":"unknown","go_version":"go1.27.2","grpc_address":"[::]:50051"}
{"time":"2026-10-16T11:02:42.48859265Z","level":"INFO","msg":"service stopping","service":"payments","health":"NOT_SERVING","shutdown_timeout":"10s"}
{"time":"2026-10-16T11:02:42.488961345Z","level":"INFO","msg":"service stopped","service":"payments","stop_mode":"graceful"}
```

<a id="go-smoke-evidence-windows-final-report"></a>

## Приложение: Go smoke evidence: windows-final-report.json

Полный исходный документ: `docs/development/evidence/go-smoke/windows-final-report.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "schema_version": 2,
  "run_id": "609bf33e-2a05-f1d0-0f0c-2326b75f5f9f",
  "project": "faultline-smoke-609bf33e",
  "started_at": "2026-10-16T11:01:12.803945Z",
  "finished_at": "2026-10-16T11:02:49.1915435Z",
  "result": "PASS",
  "phase": "stop",
  "cancelled": false,
  "cleanup": "PASS",
  "host_platform": "windows",
  "go_version": "go1.27.2",
  "docker_version": "29.8.2",
  "docker_architecture": "linux/x86_64",
  "compose_version": "5.6.0",
  "source_revision": "8e49367c868af33023a6c1f7a8f94db447a828eb",
  "source_dirty": false,
  "isolation_after_orders_stop": "PASS",
  "services": {
    "inventory": {
      "id": "f4db4ff30e18a1dd35da4324c85eb4a2d061f1b679ab4fc541d926fbf6a5c22f",
      "image_id": "sha256:329580dc579c18c33f8c9a7d403157683fad9ee043860c05a3470958de6d6a22",
      "healthy": true,
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "status": "exited",
      "stop_mode": "graceful",
      "security": {
        "capabilities_dropped": true,
        "host_ports": [],
        "internal_network": true,
        "mounts": [],
        "no_new_privileges": true,
        "read_only": true,
        "user": "65532:65532"
      }
    },
    "orders": {
      "id": "9185e572e18b34c15df15fcc6a7b5982a02167ee39f08eec3c98175ad87c15c9",
      "image_id": "sha256:58c2aa983f0225f41d369dc3a37671a00ea0c7fd5c2abf4d9d3d85e761e4b8f5",
      "healthy": true,
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "status": "exited",
      "stop_mode": "graceful",
      "security": {
        "capabilities_dropped": true,
        "host_ports": [],
        "internal_network": true,
        "mounts": [],
        "no_new_privileges": true,
        "read_only": true,
        "user": "65532:65532"
      }
    },
    "payments": {
      "id": "0708c6493b8cdef80faa97d48266c7865bb5263ae5a826bcc0c2411b4f39af83",
      "image_id": "sha256:730e62b427afc0af0e8dcbcfe8fe161822e8f6ee8ee05a057a5e32b51b3928fc",
      "healthy": true,
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "status": "exited",
      "stop_mode": "graceful",
      "security": {
        "capabilities_dropped": true,
        "host_ports": [],
        "internal_network": true,
        "mounts": [],
        "no_new_privileges": true,
        "read_only": true,
        "user": "65532:65532"
      }
    }
  },
  "commands": [
    {
      "program": "docker",
      "args": [
        "version",
        "--format",
        "{{.Server.Version}}"
      ],
      "started_at": "2026-10-16T11:01:12.8061309Z",
      "finished_at": "2026-10-16T11:01:13.1215169Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "info",
        "--format",
        "{{.OSType}}/{{.Architecture}}"
      ],
      "started_at": "2026-10-16T11:01:13.1226089Z",
      "finished_at": "2026-10-16T11:01:14.0397636Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "version",
        "--short"
      ],
      "started_at": "2026-10-16T11:01:14.0407758Z",
      "finished_at": "2026-10-16T11:01:14.4041291Z",
      "exit_code": 0
    },
    {
      "program": "git",
      "args": [
        "rev-parse",
        "HEAD"
      ],
      "started_at": "2026-10-16T11:01:14.405134Z",
      "finished_at": "2026-10-16T11:01:14.477643Z",
      "exit_code": 0
    },
    {
      "program": "git",
      "args": [
        "status",
        "--porcelain"
      ],
      "started_at": "2026-10-16T11:01:14.4787349Z",
      "finished_at": "2026-10-16T11:01:14.5504047Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "config",
        "--format",
        "json"
      ],
      "started_at": "2026-10-16T11:01:14.5524348Z",
      "finished_at": "2026-10-16T11:01:15.1123095Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "build",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T11:01:15.1149278Z",
      "finished_at": "2026-10-16T11:02:25.7186139Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "up",
        "-d",
        "--wait",
        "--wait-timeout",
        "60",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T11:02:25.7217247Z",
      "finished_at": "2026-10-16T11:02:31.9319969Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "exec",
        "-T",
        "orders",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T11:02:31.9348411Z",
      "finished_at": "2026-10-16T11:02:32.5146103Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T11:02:32.5156425Z",
      "finished_at": "2026-10-16T11:02:33.0495287Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "9185e572e18b34c15df15fcc6a7b5982a02167ee39f08eec3c98175ad87c15c9"
      ],
      "started_at": "2026-10-16T11:02:33.0506295Z",
      "finished_at": "2026-10-16T11:02:33.3058091Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "97225697ac40cbfa4a75774d931a6b349ede98fe9dbc7ba45ac720616de1fd67"
      ],
      "started_at": "2026-10-16T11:02:33.3086318Z",
      "finished_at": "2026-10-16T11:02:33.545552Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T11:02:33.5466078Z",
      "finished_at": "2026-10-16T11:02:34.2305222Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T11:02:34.2315766Z",
      "finished_at": "2026-10-16T11:02:34.8477697Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "f4db4ff30e18a1dd35da4324c85eb4a2d061f1b679ab4fc541d926fbf6a5c22f"
      ],
      "started_at": "2026-10-16T11:02:34.8494415Z",
      "finished_at": "2026-10-16T11:02:35.1119187Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "97225697ac40cbfa4a75774d931a6b349ede98fe9dbc7ba45ac720616de1fd67"
      ],
      "started_at": "2026-10-16T11:02:35.1130388Z",
      "finished_at": "2026-10-16T11:02:35.3674833Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T11:02:35.3690686Z",
      "finished_at": "2026-10-16T11:02:35.9240251Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T11:02:35.9256394Z",
      "finished_at": "2026-10-16T11:02:36.4592815Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "0708c6493b8cdef80faa97d48266c7865bb5263ae5a826bcc0c2411b4f39af83"
      ],
      "started_at": "2026-10-16T11:02:36.46026Z",
      "finished_at": "2026-10-16T11:02:36.7210212Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "97225697ac40cbfa4a75774d931a6b349ede98fe9dbc7ba45ac720616de1fd67"
      ],
      "started_at": "2026-10-16T11:02:36.7227503Z",
      "finished_at": "2026-10-16T11:02:36.9736557Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "stop",
        "orders"
      ],
      "started_at": "2026-10-16T11:02:36.9764457Z",
      "finished_at": "2026-10-16T11:02:37.6159581Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T11:02:37.6196494Z",
      "finished_at": "2026-10-16T11:02:38.1988892Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "9185e572e18b34c15df15fcc6a7b5982a02167ee39f08eec3c98175ad87c15c9"
      ],
      "started_at": "2026-10-16T11:02:38.201707Z",
      "finished_at": "2026-10-16T11:02:38.4347498Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "orders"
      ],
      "started_at": "2026-10-16T11:02:38.4381476Z",
      "finished_at": "2026-10-16T11:02:38.9984874Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T11:02:38.9997754Z",
      "finished_at": "2026-10-16T11:02:39.5732388Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T11:02:39.5760253Z",
      "finished_at": "2026-10-16T11:02:40.0882447Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "f4db4ff30e18a1dd35da4324c85eb4a2d061f1b679ab4fc541d926fbf6a5c22f"
      ],
      "started_at": "2026-10-16T11:02:40.090613Z",
      "finished_at": "2026-10-16T11:02:40.3246392Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "97225697ac40cbfa4a75774d931a6b349ede98fe9dbc7ba45ac720616de1fd67"
      ],
      "started_at": "2026-10-16T11:02:40.3262248Z",
      "finished_at": "2026-10-16T11:02:40.557235Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T11:02:40.5587891Z",
      "finished_at": "2026-10-16T11:02:41.1101859Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T11:02:41.1113663Z",
      "finished_at": "2026-10-16T11:02:41.6427474Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "0708c6493b8cdef80faa97d48266c7865bb5263ae5a826bcc0c2411b4f39af83"
      ],
      "started_at": "2026-10-16T11:02:41.6451706Z",
      "finished_at": "2026-10-16T11:02:41.8798591Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "97225697ac40cbfa4a75774d931a6b349ede98fe9dbc7ba45ac720616de1fd67"
      ],
      "started_at": "2026-10-16T11:02:41.8815022Z",
      "finished_at": "2026-10-16T11:02:42.1404997Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "stop",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T11:02:42.1422942Z",
      "finished_at": "2026-10-16T11:02:42.8612791Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T11:02:42.8637033Z",
      "finished_at": "2026-10-16T11:02:43.4660331Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "f4db4ff30e18a1dd35da4324c85eb4a2d061f1b679ab4fc541d926fbf6a5c22f"
      ],
      "started_at": "2026-10-16T11:02:43.4675653Z",
      "finished_at": "2026-10-16T11:02:43.7328487Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "inventory"
      ],
      "started_at": "2026-10-16T11:02:43.7361231Z",
      "finished_at": "2026-10-16T11:02:44.3214316Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T11:02:44.324505Z",
      "finished_at": "2026-10-16T11:02:44.8678404Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "0708c6493b8cdef80faa97d48266c7865bb5263ae5a826bcc0c2411b4f39af83"
      ],
      "started_at": "2026-10-16T11:02:44.8688639Z",
      "finished_at": "2026-10-16T11:02:45.1035758Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "payments"
      ],
      "started_at": "2026-10-16T11:02:45.1050878Z",
      "finished_at": "2026-10-16T11:02:45.6578265Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "orders"
      ],
      "started_at": "2026-10-16T11:02:45.6608044Z",
      "finished_at": "2026-10-16T11:02:46.2150267Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "inventory"
      ],
      "started_at": "2026-10-16T11:02:46.2166575Z",
      "finished_at": "2026-10-16T11:02:46.7915107Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "payments"
      ],
      "started_at": "2026-10-16T11:02:46.7931247Z",
      "finished_at": "2026-10-16T11:02:47.3457435Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-609bf33e",
        "down",
        "--remove-orphans"
      ],
      "started_at": "2026-10-16T11:02:47.3474183Z",
      "finished_at": "2026-10-16T11:02:48.0249597Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "ps",
        "--all",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-609bf33e"
      ],
      "started_at": "2026-10-16T11:02:48.0275166Z",
      "finished_at": "2026-10-16T11:02:48.2684199Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "ls",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-609bf33e"
      ],
      "started_at": "2026-10-16T11:02:48.27066Z",
      "finished_at": "2026-10-16T11:02:48.5065366Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "ps",
        "--all",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-609bf33e"
      ],
      "started_at": "2026-10-16T11:02:48.7080415Z",
      "finished_at": "2026-10-16T11:02:48.939512Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "ls",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-609bf33e"
      ],
      "started_at": "2026-10-16T11:02:48.9405844Z",
      "finished_at": "2026-10-16T11:02:49.1904664Z",
      "exit_code": 0
    }
  ],
  "recovery_command": [
    "docker",
    "compose",
    "-f",
    "compose.yaml",
    "-p",
    "faultline-smoke-609bf33e",
    "down",
    "--remove-orphans"
  ],
  "limitations": [
    "No business operations or capacity measurement",
    "docker_architecture describes the backend, not individual image platforms (F03)",
    "SIGKILL or host failure requires recovery_command; signal cleanup cannot be guaranteed in those cases"
  ]
}
```

<a id="go-smoke-evidence-windows-interrupt-report"></a>

## Приложение: Go smoke evidence: windows-interrupt-report.json

Полный исходный документ: `docs/development/evidence/go-smoke/windows-interrupt-report.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "schema_version": 2,
  "run_id": "0d41daa4-2151-8b9e-8268-0a597b5164c0",
  "project": "faultline-smoke-0d41daa4",
  "started_at": "2026-10-16T10:49:59.829639Z",
  "finished_at": "2026-10-16T10:50:15.8692974Z",
  "result": "FAIL",
  "phase": "check",
  "failed_phase": "check",
  "error": "docker compose -f compose.yaml -p faultline-smoke-0d41daa4 ps --all -q orders failed (exit 1; context canceled)",
  "cancelled": true,
  "cleanup": "PASS",
  "host_platform": "windows",
  "go_version": "go1.27.2",
  "docker_version": "29.8.2",
  "docker_architecture": "linux/x86_64",
  "compose_version": "5.6.0",
  "source_revision": "511b9a85c7006cc7042263e0b87e7eb7fb8b116c",
  "source_dirty": true,
  "services": {},
  "commands": [
    {
      "program": "docker",
      "args": [
        "version",
        "--format",
        "{{.Server.Version}}"
      ],
      "started_at": "2026-10-16T10:49:59.8314426Z",
      "finished_at": "2026-10-16T10:50:00.1633603Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "info",
        "--format",
        "{{.OSType}}/{{.Architecture}}"
      ],
      "started_at": "2026-10-16T10:50:00.1644835Z",
      "finished_at": "2026-10-16T10:50:01.024469Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "version",
        "--short"
      ],
      "started_at": "2026-10-16T10:50:01.0255365Z",
      "finished_at": "2026-10-16T10:50:01.3847432Z",
      "exit_code": 0
    },
    {
      "program": "git",
      "args": [
        "rev-parse",
        "HEAD"
      ],
      "started_at": "2026-10-16T10:50:01.3858188Z",
      "finished_at": "2026-10-16T10:50:01.4596761Z",
      "exit_code": 0
    },
    {
      "program": "git",
      "args": [
        "status",
        "--porcelain"
      ],
      "started_at": "2026-10-16T10:50:01.4608542Z",
      "finished_at": "2026-10-16T10:50:01.5528102Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-0d41daa4",
        "config",
        "--format",
        "json"
      ],
      "started_at": "2026-10-16T10:50:01.5549608Z",
      "finished_at": "2026-10-16T10:50:02.1193906Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-0d41daa4",
        "build",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T10:50:02.1210997Z",
      "finished_at": "2026-10-16T10:50:04.8631976Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-0d41daa4",
        "up",
        "-d",
        "--wait",
        "--wait-timeout",
        "60",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T10:50:04.8648684Z",
      "finished_at": "2026-10-16T10:50:11.0350184Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-0d41daa4",
        "exec",
        "-T",
        "orders",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T10:50:11.0378213Z",
      "finished_at": "2026-10-16T10:50:11.7853655Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-0d41daa4",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T10:50:11.7874186Z",
      "finished_at": "2026-10-16T10:50:12.0221344Z",
      "exit_code": 1,
      "error": "context canceled"
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-0d41daa4",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "orders"
      ],
      "started_at": "2026-10-16T10:50:12.0238006Z",
      "finished_at": "2026-10-16T10:50:12.6291675Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-0d41daa4",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "inventory"
      ],
      "started_at": "2026-10-16T10:50:12.6323995Z",
      "finished_at": "2026-10-16T10:50:13.2083822Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-0d41daa4",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "payments"
      ],
      "started_at": "2026-10-16T10:50:13.2117141Z",
      "finished_at": "2026-10-16T10:50:13.7669729Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-0d41daa4",
        "down",
        "--remove-orphans"
      ],
      "started_at": "2026-10-16T10:50:13.7687867Z",
      "finished_at": "2026-10-16T10:50:14.663702Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "ps",
        "--all",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-0d41daa4"
      ],
      "started_at": "2026-10-16T10:50:14.6652611Z",
      "finished_at": "2026-10-16T10:50:14.9379534Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "ls",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-0d41daa4"
      ],
      "started_at": "2026-10-16T10:50:14.9388931Z",
      "finished_at": "2026-10-16T10:50:15.1964348Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "ps",
        "--all",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-0d41daa4"
      ],
      "started_at": "2026-10-16T10:50:15.3983351Z",
      "finished_at": "2026-10-16T10:50:15.6304179Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "ls",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-0d41daa4"
      ],
      "started_at": "2026-10-16T10:50:15.6340613Z",
      "finished_at": "2026-10-16T10:50:15.8681773Z",
      "exit_code": 0
    }
  ],
  "recovery_command": [
    "docker",
    "compose",
    "-f",
    "compose.yaml",
    "-p",
    "faultline-smoke-0d41daa4",
    "down",
    "--remove-orphans"
  ],
  "limitations": [
    "No business operations or capacity measurement",
    "docker_architecture describes the backend, not individual image platforms (F03)",
    "SIGKILL or host failure requires recovery_command; signal cleanup cannot be guaranteed in those cases"
  ]
}
```

<a id="p02-docs-development-evidence-p02-windows-runtime-report-json"></a>

## Приложение: docs/development/evidence/p02/windows-runtime-report.json

Полный исходный документ: `docs/development/evidence/p02/windows-runtime-report.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "schema_version": 2,
  "run_id": "5562c957-f2db-dcae-2c62-0ec73c2e8a91",
  "project": "faultline-smoke-5562c957",
  "started_at": "2026-10-16T16:13:18.2802558Z",
  "finished_at": "2026-10-16T16:15:07.6836875Z",
  "result": "PASS",
  "phase": "stop",
  "cancelled": false,
  "cleanup": "PASS",
  "host_platform": "windows",
  "go_version": "go1.27.2",
  "docker_version": "29.8.2",
  "docker_architecture": "linux/x86_64",
  "compose_version": "5.6.0",
  "source_revision": "843e95ed2b65a44729fbf297b986363c588d51ca",
  "source_dirty": true,
  "isolation_after_orders_stop": "PASS",
  "services": {
    "inventory": {
      "id": "644de19c878f17aca005976d38e2d97e7eafa3f5e8a6755ef16ce8a082d16eb4",
      "image_id": "sha256:ff5317081044d8eb1a830157220cb1ddfed80d641e84a66a00ad80b1f98b407a",
      "healthy": true,
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "status": "exited",
      "stop_mode": "graceful",
      "security": {
        "capabilities_dropped": true,
        "host_ports": [],
        "internal_network": true,
        "mounts": [],
        "no_new_privileges": true,
        "read_only": true,
        "user": "65532:65532"
      }
    },
    "orders": {
      "id": "09392ec030f308b80f50dbf6219b73e0ca479ba27ec514c0d768b6d24668b24c",
      "image_id": "sha256:6a5ee48e44cf9e56088f49e9d51c79060c104f7049862e9118b15d99ad88ee01",
      "healthy": true,
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "status": "exited",
      "stop_mode": "graceful",
      "security": {
        "capabilities_dropped": true,
        "host_ports": [],
        "internal_network": true,
        "mounts": [],
        "no_new_privileges": true,
        "read_only": true,
        "user": "65532:65532"
      }
    },
    "payments": {
      "id": "b689b51321ea12943d358eb138c9af42df46ecc045f01f455dd66f425b3e67f8",
      "image_id": "sha256:524f361a5e6cfd7478cff626cb95d703a288da2321fbed121500204344235c3a",
      "healthy": true,
      "stopped": true,
      "exit_code": 0,
      "oom_killed": false,
      "status": "exited",
      "stop_mode": "graceful",
      "security": {
        "capabilities_dropped": true,
        "host_ports": [],
        "internal_network": true,
        "mounts": [],
        "no_new_privileges": true,
        "read_only": true,
        "user": "65532:65532"
      }
    }
  },
  "commands": [
    {
      "program": "docker",
      "args": [
        "version",
        "--format",
        "{{.Server.Version}}"
      ],
      "started_at": "2026-10-16T16:13:18.2829842Z",
      "finished_at": "2026-10-16T16:13:18.5890779Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "info",
        "--format",
        "{{.OSType}}/{{.Architecture}}"
      ],
      "started_at": "2026-10-16T16:13:18.5901195Z",
      "finished_at": "2026-10-16T16:13:20.3082358Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "version",
        "--short"
      ],
      "started_at": "2026-10-16T16:13:20.3098121Z",
      "finished_at": "2026-10-16T16:13:20.7366132Z",
      "exit_code": 0
    },
    {
      "program": "git",
      "args": [
        "rev-parse",
        "HEAD"
      ],
      "started_at": "2026-10-16T16:13:20.7373355Z",
      "finished_at": "2026-10-16T16:13:20.8154235Z",
      "exit_code": 0
    },
    {
      "program": "git",
      "args": [
        "status",
        "--porcelain"
      ],
      "started_at": "2026-10-16T16:13:20.8171969Z",
      "finished_at": "2026-10-16T16:13:20.9311169Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "config",
        "--format",
        "json"
      ],
      "started_at": "2026-10-16T16:13:20.9331637Z",
      "finished_at": "2026-10-16T16:13:21.5216634Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "build",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T16:13:21.5267703Z",
      "finished_at": "2026-10-16T16:14:41.3349157Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "up",
        "-d",
        "--wait",
        "--wait-timeout",
        "60",
        "orders",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T16:14:41.3400836Z",
      "finished_at": "2026-10-16T16:14:48.0541367Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "exec",
        "-T",
        "orders",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T16:14:48.0562843Z",
      "finished_at": "2026-10-16T16:14:48.7619294Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T16:14:48.7641022Z",
      "finished_at": "2026-10-16T16:14:49.4464479Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "09392ec030f308b80f50dbf6219b73e0ca479ba27ec514c0d768b6d24668b24c"
      ],
      "started_at": "2026-10-16T16:14:49.4503083Z",
      "finished_at": "2026-10-16T16:14:49.7276956Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b875a819cf10b612795177065d8ac30ca4ec40468fc9a9111815c7e3c87377ce"
      ],
      "started_at": "2026-10-16T16:14:49.7287347Z",
      "finished_at": "2026-10-16T16:14:50.0018181Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T16:14:50.005557Z",
      "finished_at": "2026-10-16T16:14:50.6381462Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T16:14:50.6397081Z",
      "finished_at": "2026-10-16T16:14:51.2315804Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "644de19c878f17aca005976d38e2d97e7eafa3f5e8a6755ef16ce8a082d16eb4"
      ],
      "started_at": "2026-10-16T16:14:51.2337506Z",
      "finished_at": "2026-10-16T16:14:51.5016289Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b875a819cf10b612795177065d8ac30ca4ec40468fc9a9111815c7e3c87377ce"
      ],
      "started_at": "2026-10-16T16:14:51.5045374Z",
      "finished_at": "2026-10-16T16:14:51.7751067Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T16:14:51.776127Z",
      "finished_at": "2026-10-16T16:14:52.4067784Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T16:14:52.4083153Z",
      "finished_at": "2026-10-16T16:14:53.0207251Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "b689b51321ea12943d358eb138c9af42df46ecc045f01f455dd66f425b3e67f8"
      ],
      "started_at": "2026-10-16T16:14:53.0218701Z",
      "finished_at": "2026-10-16T16:14:53.2907252Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b875a819cf10b612795177065d8ac30ca4ec40468fc9a9111815c7e3c87377ce"
      ],
      "started_at": "2026-10-16T16:14:53.2918639Z",
      "finished_at": "2026-10-16T16:14:53.5679535Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "stop",
        "orders"
      ],
      "started_at": "2026-10-16T16:14:53.5701404Z",
      "finished_at": "2026-10-16T16:14:54.2851639Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "ps",
        "--all",
        "-q",
        "orders"
      ],
      "started_at": "2026-10-16T16:14:54.2872939Z",
      "finished_at": "2026-10-16T16:14:54.8733926Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "09392ec030f308b80f50dbf6219b73e0ca479ba27ec514c0d768b6d24668b24c"
      ],
      "started_at": "2026-10-16T16:14:54.8745041Z",
      "finished_at": "2026-10-16T16:14:55.1413054Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "orders"
      ],
      "started_at": "2026-10-16T16:14:55.1453533Z",
      "finished_at": "2026-10-16T16:14:55.7769553Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "exec",
        "-T",
        "inventory",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T16:14:55.7790284Z",
      "finished_at": "2026-10-16T16:14:56.4059657Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T16:14:56.4080812Z",
      "finished_at": "2026-10-16T16:14:57.000445Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "644de19c878f17aca005976d38e2d97e7eafa3f5e8a6755ef16ce8a082d16eb4"
      ],
      "started_at": "2026-10-16T16:14:57.0016884Z",
      "finished_at": "2026-10-16T16:14:57.272168Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b875a819cf10b612795177065d8ac30ca4ec40468fc9a9111815c7e3c87377ce"
      ],
      "started_at": "2026-10-16T16:14:57.2746073Z",
      "finished_at": "2026-10-16T16:14:57.5486864Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "exec",
        "-T",
        "payments",
        "/app/service",
        "healthcheck"
      ],
      "started_at": "2026-10-16T16:14:57.5498305Z",
      "finished_at": "2026-10-16T16:14:58.1625596Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T16:14:58.1657216Z",
      "finished_at": "2026-10-16T16:14:58.7317792Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "b689b51321ea12943d358eb138c9af42df46ecc045f01f455dd66f425b3e67f8"
      ],
      "started_at": "2026-10-16T16:14:58.7358853Z",
      "finished_at": "2026-10-16T16:14:59.0352535Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "inspect",
        "b875a819cf10b612795177065d8ac30ca4ec40468fc9a9111815c7e3c87377ce"
      ],
      "started_at": "2026-10-16T16:14:59.0368331Z",
      "finished_at": "2026-10-16T16:14:59.3222731Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "stop",
        "inventory",
        "payments"
      ],
      "started_at": "2026-10-16T16:14:59.3298208Z",
      "finished_at": "2026-10-16T16:15:00.1251358Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "ps",
        "--all",
        "-q",
        "inventory"
      ],
      "started_at": "2026-10-16T16:15:00.1264533Z",
      "finished_at": "2026-10-16T16:15:00.7662995Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "644de19c878f17aca005976d38e2d97e7eafa3f5e8a6755ef16ce8a082d16eb4"
      ],
      "started_at": "2026-10-16T16:15:00.7683632Z",
      "finished_at": "2026-10-16T16:15:01.0351817Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "inventory"
      ],
      "started_at": "2026-10-16T16:15:01.0368127Z",
      "finished_at": "2026-10-16T16:15:01.6481027Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "ps",
        "--all",
        "-q",
        "payments"
      ],
      "started_at": "2026-10-16T16:15:01.6520038Z",
      "finished_at": "2026-10-16T16:15:02.2409544Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "inspect",
        "b689b51321ea12943d358eb138c9af42df46ecc045f01f455dd66f425b3e67f8"
      ],
      "started_at": "2026-10-16T16:15:02.2420895Z",
      "finished_at": "2026-10-16T16:15:02.5110147Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "payments"
      ],
      "started_at": "2026-10-16T16:15:02.5115283Z",
      "finished_at": "2026-10-16T16:15:03.1179026Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "orders"
      ],
      "started_at": "2026-10-16T16:15:03.1195078Z",
      "finished_at": "2026-10-16T16:15:03.7469717Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "inventory"
      ],
      "started_at": "2026-10-16T16:15:03.7497576Z",
      "finished_at": "2026-10-16T16:15:04.3971364Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "logs",
        "--no-color",
        "--no-log-prefix",
        "payments"
      ],
      "started_at": "2026-10-16T16:15:04.4025292Z",
      "finished_at": "2026-10-16T16:15:05.0339073Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "compose",
        "-f",
        "compose.yaml",
        "-p",
        "faultline-smoke-5562c957",
        "down",
        "--remove-orphans",
        "--rmi",
        "local"
      ],
      "started_at": "2026-10-16T16:15:05.0355312Z",
      "finished_at": "2026-10-16T16:15:05.8296656Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "ps",
        "--all",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-5562c957"
      ],
      "started_at": "2026-10-16T16:15:05.8327177Z",
      "finished_at": "2026-10-16T16:15:06.1021426Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "ls",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-5562c957"
      ],
      "started_at": "2026-10-16T16:15:06.1031938Z",
      "finished_at": "2026-10-16T16:15:06.3765194Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "image",
        "ls",
        "--quiet",
        "--filter",
        "reference=faultline-smoke-5562c957-*:latest"
      ],
      "started_at": "2026-10-16T16:15:06.3790039Z",
      "finished_at": "2026-10-16T16:15:06.6576401Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "ps",
        "--all",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-5562c957"
      ],
      "started_at": "2026-10-16T16:15:06.8591973Z",
      "finished_at": "2026-10-16T16:15:07.1309866Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "network",
        "ls",
        "--quiet",
        "--filter",
        "label=com.docker.compose.project=faultline-smoke-5562c957"
      ],
      "started_at": "2026-10-16T16:15:07.1320467Z",
      "finished_at": "2026-10-16T16:15:07.4000639Z",
      "exit_code": 0
    },
    {
      "program": "docker",
      "args": [
        "image",
        "ls",
        "--quiet",
        "--filter",
        "reference=faultline-smoke-5562c957-*:latest"
      ],
      "started_at": "2026-10-16T16:15:07.4024797Z",
      "finished_at": "2026-10-16T16:15:07.6815516Z",
      "exit_code": 0
    }
  ],
  "recovery_command": [
    "docker",
    "compose",
    "-f",
    "compose.yaml",
    "-p",
    "faultline-smoke-5562c957",
    "down",
    "--remove-orphans",
    "--rmi",
    "local"
  ],
  "limitations": [
    "No business operations or capacity measurement",
    "docker_architecture describes the backend, not individual image platforms (F03)",
    "SIGKILL or host failure requires recovery_command; signal cleanup cannot be guaranteed in those cases"
  ]
}
```

<a id="p02-docs-development-evidence-p02-windows-cleanup-json"></a>

## Приложение: docs/development/evidence/p02/windows-cleanup.json

Полный исходный документ: `docs/development/evidence/p02/windows-cleanup.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "observed_at": "2026-10-16T19:22:00.0959015+03:00",
  "containers_remaining": 0,
  "run_images_remaining": 0,
  "remaining_images": [
    {
      "Repository": "faultline-lab-toolchain",
      "Tag": "latest",
      "ID": "c7cd1c13080d"
    },
    {
      "Repository": "postgres",
      "Tag": "<none>",
      "ID": "afc7e2d44132"
    }
  ],
  "remaining_volumes": [
    "faultline-lab_go-build",
    "faultline-lab_go-modules"
  ],
  "docker_autostart": false,
  "historical_smoke_images_removed": 21,
  "global_prune_used": false,
  "desktop_stop_exit": 0,
  "wsl_terminate_exit": 0,
  "running_wsl": "",
  "remaining_docker_processes": []
}
```

<a id="p02-docs-development-evidence-p02-windows-postgres-report-json"></a>

## Приложение: docs/development/evidence/p02/windows-postgres-report.json

Полный исходный документ: `docs/development/evidence/p02/windows-postgres-report.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "project": "faultline-data-test-90eda236ac29dee2",
  "started_at": "2026-10-16T16:12:11.1591171Z",
  "finished_at": "2026-10-16T16:13:14.6017804Z",
  "result": "PASS",
  "cleanup": "PASS",
  "events": [
    {
      "stage": "prerequisites",
      "output": "29.8.2\n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "build dbtool",
      "output": "#1 [internal] load local bake definitions\n#1 reading from stdin 679B 0.0s done\n#1 DONE 0.0s\n\n#2 [internal] load build definition from Dockerfile\n#2 transferring dockerfile: 1.49kB done\n#2 DONE 0.0s\n\n#3 [internal] load metadata for docker.io/library/golang:1.27.2-trixie@sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e\n#3 DONE 0.0s\n\n#4 [internal] load metadata for docker.io/golangci/golangci-lint:v2.14.0@sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f\n#4 DONE 0.1s\n\n#5 [internal] load .dockerignore\n#5 transferring context: 172B done\n#5 DONE 0.0s\n\n#6 [toolchain 1/6] FROM docker.io/library/golang:1.27.2-trixie@sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e\n#6 resolve docker.io/library/golang:1.27.2-trixie@sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e 0.0s done\n#6 DONE 0.0s\n\n#7 [lint 1/1] FROM docker.io/golangci/golangci-lint:v2.14.0@sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f\n#7 resolve docker.io/golangci/golangci-lint:v2.14.0@sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f 0.0s done\n#7 DONE 0.0s\n\n#8 [internal] load build context\n#8 transferring context: 19.98kB done\n#8 DONE 0.0s\n\n#9 [toolchain 3/6] COPY --from=lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint\n#9 CACHED\n\n#10 [toolchain 5/6] RUN git config --system --add safe.directory /workspace     \u0026\u0026 go mod download \u0026\u0026 go install tool\n#10 CACHED\n\n#11 [toolchain 2/6] WORKDIR /workspace\n#11 CACHED\n\n#12 [toolchain 4/6] COPY go.mod go.sum ./\n#12 CACHED\n\n#13 [toolchain 6/6] COPY --chmod=0755 tools/go.sh /usr/local/bin/lab-go\n#13 CACHED\n\n#14 [labdata-build 1/4] COPY cmd ./cmd\n#14 DONE 0.0s\n\n#15 [labdata-build 2/4] COPY internal ./internal\n#15 DONE 0.0s\n\n#16 [labdata-build 3/4] COPY migrations ./migrations\n#16 DONE 0.0s\n\n#17 [labdata-build 4/4] RUN CGO_ENABLED=0 go build -trimpath -o /out/labdata ./cmd/labdata     \u0026\u0026 CGO_ENABLED=0 go test -c -o /out/database.test ./internal/database\n#17 DONE 32.6s\n\n#18 [labdata 1/2] COPY --from=labdata-build /out/labdata /app/labdata\n#18 DONE 0.0s\n\n#19 [labdata 2/2] COPY --from=labdata-build /out/database.test /app/database.test\n#19 DONE 0.1s\n\n#20 exporting to image\n#20 exporting layers\n#20 exporting layers 0.6s done\n#20 exporting manifest sha256:4f6a2734b813e35e9ee271f0457db11898fd5ed94633b9e73dc06eb3fa0d1982\n#20 exporting manifest sha256:4f6a2734b813e35e9ee271f0457db11898fd5ed94633b9e73dc06eb3fa0d1982 0.0s done\n#20 exporting config sha256:38329ed58f971fe607196f92eec4155b639aeb8160b488c7f75aa2bc336b6c41 0.0s done\n#20 exporting attestation manifest sha256:f3f289782cd254f18b91727f7e2118c8d1b013517e6363a243f4878103b54e9d 0.0s done\n#20 exporting manifest list sha256:250ab652b2723be95dc55e3ee4e400147d220ef8ed9d6b473d1dc19bbee3b66e 0.0s done\n#20 naming to docker.io/library/faultline-data-test-90eda236ac29dee2-dbtool:latest done\n#20 unpacking to docker.io/library/faultline-data-test-90eda236ac29dee2-dbtool:latest\n#20 unpacking to docker.io/library/faultline-data-test-90eda236ac29dee2-dbtool:latest 0.2s done\n#20 DONE 0.9s\n\n#21 resolving provenance for metadata file\n#21 DONE 0.0s\n Image faultline-data-test-90eda236ac29dee2-dbtool Building \n Image faultline-data-test-90eda236ac29dee2-dbtool Built \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "start PostgreSQL",
      "output": " Network faultline-data-test-90eda236ac29dee2_data Creating \n Network faultline-data-test-90eda236ac29dee2_data Creating \n Volume faultline-data-test-90eda236ac29dee2_postgres-data Creating \n Volume faultline-data-test-90eda236ac29dee2_postgres-data Creating \n Volume faultline-data-test-90eda236ac29dee2_postgres-data Created \n Volume faultline-data-test-90eda236ac29dee2_postgres-data Created \n Network faultline-data-test-90eda236ac29dee2_data Created \n Network faultline-data-test-90eda236ac29dee2_data Created \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Creating \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Created \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Starting \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Started \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Waiting \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Healthy \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "bootstrap",
      "output": "bootstrap: PASS\n Container faultline-data-test-90eda236ac29dee2-dbtool-run-ee0c5b992210 Creating \n Container faultline-data-test-90eda236ac29dee2-dbtool-run-ee0c5b992210 Created \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "bootstrap",
      "output": "bootstrap: PASS\n Container faultline-data-test-90eda236ac29dee2-dbtool-run-75f6c25d8e83 Creating \n Container faultline-data-test-90eda236ac29dee2-dbtool-run-75f6c25d8e83 Created \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "rights, migrations and seed",
      "output": "=== RUN   TestPostgresFoundation\n=== RUN   TestPostgresFoundation/concurrent_migrations_and_immutable_history\n=== RUN   TestPostgresFoundation/role_isolation_and_denied_DDL\n=== RUN   TestPostgresFoundation/failed_migration_rollback_and_checksum_drift\n=== RUN   TestPostgresFoundation/repeat_seed_preserves_catalog_and_modified_stock\n--- PASS: TestPostgresFoundation (0.89s)\n    --- PASS: TestPostgresFoundation/concurrent_migrations_and_immutable_history (0.30s)\n    --- PASS: TestPostgresFoundation/role_isolation_and_denied_DDL (0.47s)\n    --- PASS: TestPostgresFoundation/failed_migration_rollback_and_checksum_drift (0.03s)\n    --- PASS: TestPostgresFoundation/repeat_seed_preserves_catalog_and_modified_stock (0.08s)\nPASS\n Container faultline-data-test-90eda236ac29dee2-dbtool-run-b590a77d5301 Creating \n Container faultline-data-test-90eda236ac29dee2-dbtool-run-b590a77d5301 Created \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "fingerprint before restart",
      "output": "6bd627a0b8eff5f313bc0f372ce2714f586efe3fed6d7bdbf7f78f826ac2d408\n Container faultline-data-test-90eda236ac29dee2-dbtool-run-dcb0b6dc0438 Creating \n Container faultline-data-test-90eda236ac29dee2-dbtool-run-dcb0b6dc0438 Created \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "restart PostgreSQL",
      "output": " Container faultline-data-test-90eda236ac29dee2-postgres-1 Restarting \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Started \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "wait after restart",
      "output": " Container faultline-data-test-90eda236ac29dee2-postgres-1 Running \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Waiting \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Healthy \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "repeat bootstrap",
      "output": "bootstrap: PASS\n Container faultline-data-test-90eda236ac29dee2-dbtool-run-eb397cba85da Creating \n Container faultline-data-test-90eda236ac29dee2-dbtool-run-eb397cba85da Created \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "repeat migrate",
      "output": "migrate: PASS\n Container faultline-data-test-90eda236ac29dee2-dbtool-run-3546396728d9 Creating \n Container faultline-data-test-90eda236ac29dee2-dbtool-run-3546396728d9 Created \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "repeat seed",
      "output": "seed: PASS\n Container faultline-data-test-90eda236ac29dee2-dbtool-run-66da6ca1c68f Creating \n Container faultline-data-test-90eda236ac29dee2-dbtool-run-66da6ca1c68f Created \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "fingerprint after restart",
      "output": "6bd627a0b8eff5f313bc0f372ce2714f586efe3fed6d7bdbf7f78f826ac2d408\n Container faultline-data-test-90eda236ac29dee2-dbtool-run-4b47419947d6 Creating \n Container faultline-data-test-90eda236ac29dee2-dbtool-run-4b47419947d6 Created \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "resource snapshot",
      "output": "{\"BlockIO\":\"991kB / 766kB\",\"CPUPerc\":\"5.12%\",\"Container\":\"faultline-data-test-90eda236ac29dee2-postgres-1\",\"ID\":\"1aaba13b0ff1\",\"MemPerc\":\"4.58%\",\"MemUsage\":\"23.43MiB / 512MiB\",\"Name\":\"faultline-data-test-90eda236ac29dee2-postgres-1\",\"NetIO\":\"25.6kB / 23kB\",\"PIDs\":\"9\"}\n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "PostgreSQL logs",
      "output": "postgres-1  | The files belonging to this database system will be owned by user \"postgres\".\npostgres-1  | This user must also own the server process.\npostgres-1  | \npostgres-1  | The database cluster will be initialized with locale \"en_US.utf8\".\npostgres-1  | The default database encoding has accordingly been set to \"UTF8\".\npostgres-1  | The default text search configuration will be set to \"english\".\npostgres-1  | \npostgres-1  | Data page checksums are enabled.\npostgres-1  | \npostgres-1  | fixing permissions on existing directory /var/lib/postgresql/18/docker ... ok\npostgres-1  | creating subdirectories ... ok\npostgres-1  | selecting dynamic shared memory implementation ... posix\npostgres-1  | selecting default \"max_connections\" ... 100\npostgres-1  | selecting default \"shared_buffers\" ... 128MB\npostgres-1  | selecting default time zone ... Etc/UTC\npostgres-1  | creating configuration files ... ok\npostgres-1  | running bootstrap script ... ok\npostgres-1  | performing post-bootstrap initialization ... ok\npostgres-1  | syncing data to disk ... ok\npostgres-1  | \npostgres-1  | Success. You can now start the database server using:\npostgres-1  | \npostgres-1  |     pg_ctl -D /var/lib/postgresql/18/docker -l logfile start\npostgres-1  | \npostgres-1  | waiting for server to start....2026-10-16 16:12:49.302 UTC [37] LOG:  starting PostgreSQL 18.6 (Debian 18.6-1.pgdg12+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 12.2.0-14+deb12u1) 12.2.0, 64-bit\npostgres-1  | 2026-10-16 16:12:49.302 UTC [37] LOG:  listening on Unix socket \"/var/run/postgresql/.s.PGSQL.5432\"\npostgres-1  | 2026-10-16 16:12:49.311 UTC [43] LOG:  database system was shut down at 2026-10-16 16:12:49 UTC\npostgres-1  | 2026-10-16 16:12:49.322 UTC [37] LOG:  database system is ready to accept connections\npostgres-1  |  done\npostgres-1  | server started\npostgres-1  | \npostgres-1  | /usr/local/bin/docker-entrypoint.sh: ignoring /docker-entrypoint-initdb.d/*\npostgres-1  | \npostgres-1  | waiting for server to shut down....2026-10-16 16:12:49.457 UTC [37] LOG:  received fast shutdown request\npostgres-1  | 2026-10-16 16:12:49.460 UTC [37] LOG:  aborting any active transactions\npostgres-1  | 2026-10-16 16:12:49.463 UTC [37] LOG:  background worker \"logical replication launcher\" (PID 46) exited with exit code 1\npostgres-1  | 2026-10-16 16:12:49.464 UTC [41] LOG:  shutting down\npostgres-1  | 2026-10-16 16:12:49.467 UTC [41] LOG:  checkpoint starting: shutdown immediate\npostgres-1  | 2026-10-16 16:12:49.481 UTC [41] LOG:  checkpoint complete: wrote 0 buffers (0.0%), wrote 3 SLRU buffers; 0 WAL file(s) added, 0 removed, 0 recycled; write=0.005 s, sync=0.002 s, total=0.018 s; sync files=2, longest=0.002 s, average=0.001 s; distance=0 kB, estimate=0 kB; lsn=0/1759450, redo lsn=0/1759450\npostgres-1  | 2026-10-16 16:12:49.494 UTC [37] LOG:  database system is shut down\npostgres-1  |  done\npostgres-1  | server stopped\npostgres-1  | \npostgres-1  | PostgreSQL init process complete; ready for start up.\npostgres-1  | \npostgres-1  | 2026-10-16 16:12:49.583 UTC [1] LOG:  starting PostgreSQL 18.6 (Debian 18.6-1.pgdg12+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 12.2.0-14+deb12u1) 12.2.0, 64-bit\npostgres-1  | 2026-10-16 16:12:49.583 UTC [1] LOG:  listening on IPv4 address \"0.0.0.0\", port 5432\npostgres-1  | 2026-10-16 16:12:49.583 UTC [1] LOG:  listening on IPv6 address \"::\", port 5432\npostgres-1  | 2026-10-16 16:12:49.585 UTC [1] LOG:  listening on Unix socket \"/var/run/postgresql/.s.PGSQL.5432\"\npostgres-1  | 2026-10-16 16:12:49.593 UTC [57] LOG:  database system was shut down at 2026-10-16 16:12:49 UTC\npostgres-1  | 2026-10-16 16:12:49.597 UTC [1] LOG:  database system is ready to accept connections\npostgres-1  | 2026-10-16 16:12:57.001 UTC [105] FATAL:  permission denied for database \"inventory_db\"\npostgres-1  | 2026-10-16 16:12:57.001 UTC [105] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.007 UTC [106] FATAL:  permission denied for database \"payments_db\"\npostgres-1  | 2026-10-16 16:12:57.007 UTC [106] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.063 UTC [108] FATAL:  permission denied for database \"inventory_db\"\npostgres-1  | 2026-10-16 16:12:57.063 UTC [108] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.072 UTC [109] FATAL:  permission denied for database \"payments_db\"\npostgres-1  | 2026-10-16 16:12:57.072 UTC [109] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.080 UTC [116] ERROR:  permission denied for database orders_db\npostgres-1  | 2026-10-16 16:12:57.081 UTC [116] ERROR:  permission denied for schema public at character 14\npostgres-1  | 2026-10-16 16:12:57.082 UTC [116] ERROR:  permission denied to create temporary tables in database \"orders_db\" at character 19\npostgres-1  | 2026-10-16 16:12:57.082 UTC [116] ERROR:  permission denied to set role \"orders_migrator\"\npostgres-1  | 2026-10-16 16:12:57.082 UTC [116] ERROR:  permission denied for schema lab_migrations at character 13\npostgres-1  | 2026-10-16 16:12:57.101 UTC [116] ERROR:  must be owner of table privilege_probe\npostgres-1  | 2026-10-16 16:12:57.102 UTC [116] ERROR:  permission denied for table privilege_probe\npostgres-1  | 2026-10-16 16:12:57.169 UTC [118] FATAL:  permission denied for database \"postgres\"\npostgres-1  | 2026-10-16 16:12:57.169 UTC [118] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.186 UTC [119] FATAL:  permission denied for database \"template1\"\npostgres-1  | 2026-10-16 16:12:57.186 UTC [119] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.192 UTC [121] FATAL:  permission denied for database \"orders_db\"\npostgres-1  | 2026-10-16 16:12:57.192 UTC [121] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.273 UTC [123] FATAL:  permission denied for database \"payments_db\"\npostgres-1  | 2026-10-16 16:12:57.273 UTC [123] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.283 UTC [124] FATAL:  permission denied for database \"orders_db\"\npostgres-1  | 2026-10-16 16:12:57.283 UTC [124] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.297 UTC [126] FATAL:  permission denied for database \"payments_db\"\npostgres-1  | 2026-10-16 16:12:57.297 UTC [126] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.303 UTC [127] ERROR:  permission denied for database inventory_db\npostgres-1  | 2026-10-16 16:12:57.303 UTC [127] ERROR:  permission denied for schema public at character 14\npostgres-1  | 2026-10-16 16:12:57.304 UTC [127] ERROR:  permission denied to create temporary tables in database \"inventory_db\" at character 19\npostgres-1  | 2026-10-16 16:12:57.304 UTC [127] ERROR:  permission denied to set role \"inventory_migrator\"\npostgres-1  | 2026-10-16 16:12:57.304 UTC [127] ERROR:  permission denied for schema lab_migrations at character 13\npostgres-1  | 2026-10-16 16:12:57.317 UTC [127] ERROR:  must be owner of table privilege_probe\npostgres-1  | 2026-10-16 16:12:57.318 UTC [127] ERROR:  permission denied for table privilege_probe\npostgres-1  | 2026-10-16 16:12:57.325 UTC [129] FATAL:  permission denied for database \"postgres\"\npostgres-1  | 2026-10-16 16:12:57.325 UTC [129] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.360 UTC [130] FATAL:  permission denied for database \"template1\"\npostgres-1  | 2026-10-16 16:12:57.360 UTC [130] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.368 UTC [131] FATAL:  permission denied for database \"orders_db\"\npostgres-1  | 2026-10-16 16:12:57.368 UTC [131] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.373 UTC [132] FATAL:  permission denied for database \"inventory_db\"\npostgres-1  | 2026-10-16 16:12:57.373 UTC [132] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.384 UTC [134] FATAL:  permission denied for database \"orders_db\"\npostgres-1  | 2026-10-16 16:12:57.384 UTC [134] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.391 UTC [135] FATAL:  permission denied for database \"inventory_db\"\npostgres-1  | 2026-10-16 16:12:57.391 UTC [135] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.403 UTC [137] ERROR:  permission denied for database payments_db\npostgres-1  | 2026-10-16 16:12:57.403 UTC [137] ERROR:  permission denied for schema public at character 14\npostgres-1  | 2026-10-16 16:12:57.403 UTC [137] ERROR:  permission denied to create temporary tables in database \"payments_db\" at character 19\npostgres-1  | 2026-10-16 16:12:57.403 UTC [137] ERROR:  permission denied to set role \"payments_migrator\"\npostgres-1  | 2026-10-16 16:12:57.404 UTC [137] ERROR:  permission denied for schema lab_migrations at character 13\npostgres-1  | 2026-10-16 16:12:57.417 UTC [137] ERROR:  must be owner of table privilege_probe\npostgres-1  | 2026-10-16 16:12:57.417 UTC [137] ERROR:  permission denied for table privilege_probe\npostgres-1  | 2026-10-16 16:12:57.424 UTC [139] FATAL:  permission denied for database \"postgres\"\npostgres-1  | 2026-10-16 16:12:57.424 UTC [139] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.458 UTC [140] FATAL:  permission denied for database \"template1\"\npostgres-1  | 2026-10-16 16:12:57.458 UTC [140] DETAIL:  User does not have CONNECT privilege.\npostgres-1  | 2026-10-16 16:12:57.482 UTC [142] ERROR:  division by zero\npostgres-1  | 2026-10-16 16:12:57.575 UTC [145] ERROR:  new row for relation \"stock\" violates check constraint \"stock_available_check\"\npostgres-1  | 2026-10-16 16:12:57.575 UTC [145] DETAIL:  Failing row contains (SKU-001, -1).\npostgres-1  | 2026-10-16 16:12:57.576 UTC [145] ERROR:  new row for relation \"products\" violates check constraint \"products_price_minor_check\"\npostgres-1  | 2026-10-16 16:12:57.576 UTC [145] DETAIL:  Failing row contains (SKU-001, Notebook, -1, RUB).\npostgres-1  | 2026-10-16 16:12:57.576 UTC [145] ERROR:  new row for relation \"products\" violates check constraint \"products_currency_check\"\npostgres-1  | 2026-10-16 16:12:57.576 UTC [145] DETAIL:  Failing row contains (SKU-001, Notebook, 99900, rub).\npostgres-1  | 2026-10-16 16:12:57.576 UTC [145] ERROR:  insert or update on table \"stock\" violates foreign key constraint \"stock_sku_fkey\"\npostgres-1  | 2026-10-16 16:12:57.576 UTC [145] DETAIL:  Key (sku)=(missing) is not present in table \"products\".\npostgres-1  | 2026-10-16 16:12:59.429 UTC [1] LOG:  received fast shutdown request\npostgres-1  | 2026-10-16 16:12:59.431 UTC [1] LOG:  aborting any active transactions\npostgres-1  | 2026-10-16 16:12:59.435 UTC [1] LOG:  background worker \"logical replication launcher\" (PID 60) exited with exit code 1\npostgres-1  | 2026-10-16 16:12:59.435 UTC [55] LOG:  shutting down\npostgres-1  | 2026-10-16 16:12:59.437 UTC [55] LOG:  checkpoint starting: shutdown immediate\npostgres-1  | 2026-10-16 16:12:59.576 UTC [55] LOG:  checkpoint complete: wrote 2868 buffers (35.0%), wrote 3 SLRU buffers; 0 WAL file(s) added, 0 removed, 1 recycled; write=0.046 s, sync=0.084 s, total=0.141 s; sync files=933, longest=0.004 s, average=0.001 s; distance=13179 kB, estimate=13179 kB; lsn=0/24380A0, redo lsn=0/24380A0\npostgres-1  | 2026-10-16 16:12:59.590 UTC [1] LOG:  database system is shut down\npostgres-1  | \npostgres-1  | PostgreSQL Database directory appears to contain a database; Skipping initialization\npostgres-1  | \npostgres-1  | 2026-10-16 16:12:59.890 UTC [1] LOG:  starting PostgreSQL 18.6 (Debian 18.6-1.pgdg12+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 12.2.0-14+deb12u1) 12.2.0, 64-bit\npostgres-1  | 2026-10-16 16:12:59.890 UTC [1] LOG:  listening on IPv4 address \"0.0.0.0\", port 5432\npostgres-1  | 2026-10-16 16:12:59.890 UTC [1] LOG:  listening on IPv6 address \"::\", port 5432\npostgres-1  | 2026-10-16 16:12:59.894 UTC [1] LOG:  listening on Unix socket \"/var/run/postgresql/.s.PGSQL.5432\"\npostgres-1  | 2026-10-16 16:12:59.908 UTC [19] LOG:  database system was shut down at 2026-10-16 16:12:59 UTC\npostgres-1  | 2026-10-16 16:12:59.917 UTC [1] LOG:  database system is ready to accept connections\n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "cleanup",
      "output": " Container faultline-data-test-90eda236ac29dee2-postgres-1 Stopping \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Stopped \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Removing \n Container faultline-data-test-90eda236ac29dee2-postgres-1 Removed \n Volume faultline-data-test-90eda236ac29dee2_postgres-data Removing \n Network faultline-data-test-90eda236ac29dee2_data Removing \n Volume faultline-data-test-90eda236ac29dee2_postgres-data Removed \n Image faultline-data-test-90eda236ac29dee2-dbtool:latest Removing \n Image faultline-data-test-90eda236ac29dee2-dbtool:latest Removed \n Network faultline-data-test-90eda236ac29dee2_data Removed \n",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "verify cleanup",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "verify cleanup",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "verify cleanup",
      "passed": true,
      "exit_code": 0
    },
    {
      "stage": "verify cleanup",
      "passed": true,
      "exit_code": 0
    }
  ]
}
```

<a id="p02-docs-development-evidence-p02-postgres-image-json"></a>

## Приложение: docs/development/evidence/p02/postgres-image.json

Полный исходный документ: `docs/development/evidence/p02/postgres-image.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "image": "postgres:18.6-bookworm",
  "digest": "sha256:afc7e2d441324c0388fa80c3d24f733b4194a4eb7f47dd8ee2b08eb1a24a647c",
  "platforms": [
    "linux/amd64",
    "linux/arm",
    "linux/arm64",
    "linux/386",
    "linux/ppc64le"
  ]
}
```

<a id="p02-docs-development-evidence-p02-postgres-md"></a>

## Приложение: docs/development/evidence/p02/POSTGRES.md

Полный исходный документ: `docs/development/evidence/p02/POSTGRES.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### PostgreSQL — проверка основы данных

Задача: [#12](https://github.com/Includeoyasi/faultline-lab/issues/12). Дата: 16 октября 2026 года, Europe/Moscow. Методика и команды — в [инструкции](#p02-docs-development-postgres-md).

Среда: Windows AMD64, Go 1.27.2, Docker Desktop 4.95 / Engine 29.8.2 / Compose 5.6, WSL2 с лимитами 2 CPU / 2 GiB RAM / 1 GiB swap. PostgreSQL 18.6-bookworm закреплён по [digest и списку платформ](#p02-docs-development-evidence-p02-postgres-image-json). Исходная база ветки — `843e95ed2b65a44729fbf297b986363c588d51ca`; локальные проверки выполнялись с изменениями рабочего дерева. Они не выдаются за прогон на чистом финальном коммите: его проверяет CI PR.

#### Результаты

| Проверка | Результат |
| :--- | :--- |
| Контейнерный `toolchain check` | PASS: формат/imports, tidy/verify, build, lint без замечаний, unit/race и govulncheck без найденных уязвимостей. Доработанные Windows CLI/cleanup дополнительно проверены native Go tests; финальный HEAD проверяет CI. |
| Реальный PostgreSQL | PASS: конкурентные миграции, неизменность истории, положительная/отрицательная матрица ролей, DDL/TEMP/SET ROLE/journal запреты, default grants, откат ошибки, checksum drift, повтор seed и ограничения таблиц. |
| Полный Windows `labdata smoke` | PASS, [исходный очищенный отчёт](#p02-docs-development-evidence-p02-windows-postgres-report-json), project `faultline-data-test-90eda236ac29dee2`; 19:12:11–19:13:14 MSK. |
| Перезапуск и повторная подготовка | PASS: `bootstrap/migrate/seed` после restart сохранили полную историю и изменённый каталог/остатки. SHA-256 до/после: `6bd627a0b8eff5f313bc0f372ce2714f586efe3fed6d7bdbf7f78f826ac2d408`. |
| Явный reset | PASS: `faultline-data-local` удалён адресно; отдельно созданный `faultline-p02-cleanup-control` volume сохранился и затем удалён собственной командой. |
| Native Windows unit/process tests | PASS: `internal/command`, `internal/smoke`, `internal/datalab`, `cmd/labdata`; в том числе Job Object, отмена/ошибка partial build, неудачная очистка, read-only каталог и exit 1 при ошибке команды миграции. |
| Регрессия runtime smoke трёх сервисов | PASS, [отчёт](#p02-docs-development-evidence-p02-windows-runtime-report-json); health, изоляция, нормальное завершение, cleanup. Проверено отсутствие тегов образов своего прогона. |
| `data-cross-check` | PASS: Windows AMD64, Linux AMD64/ARM64, macOS AMD64/ARM64. Это сборки, не пять реальных runtime-сред. |
| Linux AMD64/ARM64 runtime и SIGINT | Выполняются workflow PR; результаты и артефакты в [GitHub Actions](https://github.com/Includeoyasi/faultline-lab/actions). Локальные Windows результаты их не заменяют. |
| Личный Mac | NOT_RUN: доступа к машине не было. |

Первый полный прогон выявил read-only атрибут пустого временного каталога credentials в OneDrive. Он завершился FAIL вместо ложного PASS. Исправление восстанавливает права только каталога своего прогона перед удалением семи известных файлов и каталога; регрессия покрывает этот случай. Успешный окончательный отчёт приведён выше. SQL/init ошибки также не скрывались: неверная YAML-запись tmpfs была исправлена до полного прогона.

#### Ресурсы и завершение

В момент снимка PostgreSQL после функциональных тестов: **23.43 MiB / 512 MiB**, **5.12% CPU**. Это одиночный снимок из отчёта, не peak, baseline нагрузки или предел производительности; генератора бизнес-нагрузки ещё нет. Функциональная проверка укладывается в текущий WSL-бюджет; весь состав Kafka/Redis/наблюдаемости пока не измерен.

По замечанию владельца удалены 21 исторический smoke-образ, не использовавшийся контейнерами. Общие слои объясняют небольшой фактически освобождённый объём около 76 MB, несмотря на множество строк в UI. Исправленный runtime cleanup использует scoped `down --rmi local` и проверяет исчезновение тегов, включая отменённую/неудачную сборку. Fixture проверяет случай, когда Docker оставляет image, несмотря на успешный `down`.

[Подтверждение завершения](#p02-docs-development-evidence-p02-windows-cleanup-json): контейнеров 0, образов отдельных прогонов 0; остались актуальный toolchain и закреплённый PostgreSQL, volumes только для кешей Go. Docker Desktop остановлен, Docker WSL VM завершена, запущенных WSL-дистрибутивов и Docker-процессов нет; автозапуск выключен. Общий prune не использовался.

<a id="p02-internal-command-command-go"></a>

## Приложение: internal/command/command.go

Полный исходный документ: `internal/command/command.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
// Package command executes bounded commands with ownership of their process tree.
package command

import "context"

// Executor owns child processes and their bounded output for one explicit cwd/env.
type Executor struct {
	Directory string
	Env       []string
}

// Output captures both streams and the observed process exit code.
type Output struct {
	Stdout, Stderr string
	ExitCode       int
}

// Execute starts a program directly and joins its process tree before returning.
func (e Executor) Execute(ctx context.Context, program string, args []string) (Output, error) {
	inner := processExecutor{directory: e.Directory, env: e.Env}
	output, err := inner.execute(ctx, program, args)
	return Output{output.stdout, output.stderr, output.exitCode}, err
}
```

<a id="p02-internal-database-migrate-go"></a>

## Приложение: internal/database/migrate.go

Полный исходный документ: `internal/database/migrate.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strconv"
	"time"

	"github.com/Includeoyasi/faultline-lab/migrations"
	"github.com/jackc/pgx/v5"
)

type migration struct {
	version             int
	name, checksum, sql string
}

var migrationName = regexp.MustCompile(`^([0-9]{4})_[a-z0-9_]+\.sql$`)

func loadMigrations(files fs.FS, service string) ([]migration, error) {
	entries, err := fs.ReadDir(files, service)
	if err != nil {
		return nil, errors.New("cannot read migrations")
	}
	var result []migration
	previous := 0
	for _, entry := range entries {
		match := migrationName.FindStringSubmatch(entry.Name())
		if entry.IsDir() || match == nil {
			return nil, errors.New("invalid migration filename")
		}
		version, _ := strconv.Atoi(match[1])
		if version <= previous {
			return nil, errors.New("migration versions must be positive and strictly increasing")
		}
		data, err := fs.ReadFile(files, service+"/"+entry.Name())
		if err != nil || len(data) == 0 {
			return nil, errors.New("cannot read nonempty migration")
		}
		hash := sha256.Sum256(data)
		result = append(result, migration{version, entry.Name(), hex.EncodeToString(hash[:]), string(data)})
		previous = version
	}
	if len(result) == 0 {
		return nil, errors.New("service has no migrations")
	}
	return result, nil
}

// Migrate atomically applies pending SQL and verifies immutable service history.
func Migrate(ctx context.Context, cfg Config, service Service) error {
	return migrate(ctx, cfg, service, migrations.Files)
}

func rollback(tx pgx.Tx) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = tx.Rollback(ctx)
}

func migrate(ctx context.Context, cfg Config, service Service, files fs.FS) error {
	items, err := loadMigrations(files, service.Name)
	if err != nil {
		return err
	}
	conn, err := Connect(ctx, cfg, service.Database, service.MigrationUser)
	if err != nil {
		return err
	}
	defer Close(conn)
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return SafeError("migration transaction", err)
	}
	defer rollback(tx)
	// Lock before creating or reading the journal. It is local to this database.
	if _, err = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(706020121)`); err != nil {
		return SafeError("migration lock", err)
	}
	if _, err = tx.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS lab_migrations;
CREATE TABLE IF NOT EXISTS lab_migrations.versions (
 version integer PRIMARY KEY, name text NOT NULL, checksum text NOT NULL,
 applied_at timestamptz NOT NULL DEFAULT now()
)`); err != nil {
		return SafeError("migration journal", err)
	}
	rows, err := tx.Query(ctx, "SELECT version,name,checksum FROM lab_migrations.versions ORDER BY version")
	if err != nil {
		return SafeError("migration history", err)
	}
	var applied []migration
	for rows.Next() {
		var item migration
		if err = rows.Scan(&item.version, &item.name, &item.checksum); err != nil {
			rows.Close()
			return SafeError("migration history", err)
		}
		applied = append(applied, item)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return SafeError("migration history", err)
	}
	if len(applied) > len(items) {
		return errors.New("database migration history is ahead of this tool")
	}
	for index, item := range items {
		if index < len(applied) {
			if old := applied[index]; old.version != item.version || old.name != item.name || old.checksum != item.checksum {
				return fmt.Errorf("migration history/checksum mismatch at version %d", item.version)
			}
			continue
		}
		if _, err = tx.Exec(ctx, item.sql); err != nil {
			return SafeError(fmt.Sprintf("migration %d", item.version), err)
		}
		if _, err = tx.Exec(ctx, "INSERT INTO lab_migrations.versions(version,name,checksum) VALUES($1,$2,$3)", item.version, item.name, item.checksum); err != nil {
			return SafeError("migration recording", err)
		}
	}
	schema, app := identifierSQL(service.Name), identifierSQL(service.AppUser)
	for _, statement := range []string{
		"GRANT USAGE ON SCHEMA " + schema + " TO " + app,
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA " + schema + " TO " + app,
		"ALTER DEFAULT PRIVILEGES IN SCHEMA " + schema + " GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO " + app,
	} {
		if _, err = tx.Exec(ctx, statement); err != nil {
			return SafeError("application grants", err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return SafeError("migration commit", err)
	}
	return nil
}

// Seed inserts missing Inventory fixtures without overwriting existing values.
func Seed(ctx context.Context, cfg Config) error {
	service, err := cfg.Service("inventory")
	if err != nil {
		return err
	}
	conn, err := Connect(ctx, cfg, service.Database, service.AppUser)
	if err != nil {
		return err
	}
	defer Close(conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return SafeError("seed transaction", err)
	}
	defer rollback(tx)
	data, err := migrations.Files.ReadFile("seed/inventory.sql")
	if err != nil {
		return errors.New("cannot read inventory seed")
	}
	if _, err = tx.Exec(ctx, string(data)); err != nil {
		return SafeError("inventory seed", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return SafeError("seed commit", err)
	}
	return nil
}
```

<a id="p02-internal-database-integration-test-go"></a>

## Приложение: internal/database/integration_test.go

Полный исходный документ: `internal/database/integration_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package database

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/Includeoyasi/faultline-lab/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// These tests change their dedicated dataset. Only the isolated host smoke opts in.
func TestPostgresFoundation(t *testing.T) {
	if os.Getenv("FL_DATA_INTEGRATION") != "1" {
		t.Skip("requires isolated PostgreSQL smoke")
	}
	cfg, err := Load(os.Getenv("FL_DATA_CONFIG"), os.Getenv("FL_DATA_SECRETS"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	connect := func(t *testing.T, db, user string) *pgx.Conn {
		t.Helper()
		conn, e := Connect(ctx, cfg, db, user)
		if e != nil {
			t.Fatal(e)
		}
		t.Cleanup(func() { Close(conn) })
		return conn
	}
	scalar := func(t *testing.T, conn *pgx.Conn, query string) int {
		t.Helper()
		var value int
		if e := conn.QueryRow(ctx, query).Scan(&value); e != nil {
			t.Fatal(SafeError("test query", e))
		}
		return value
	}
	denied := func(t *testing.T, conn *pgx.Conn, sql string) {
		t.Helper()
		_, e := conn.Exec(ctx, sql)
		var pgerr *pgconn.PgError
		if !errors.As(e, &pgerr) || pgerr.Code != "42501" {
			t.Fatalf("expected permission denial, got %v", SafeError("negative permission test", e))
		}
	}
	t.Run("concurrent migrations and immutable history", func(t *testing.T) {
		for _, service := range cfg.Services {
			var wg sync.WaitGroup
			results := make(chan error, 2)
			for attempt := 0; attempt < 2; attempt++ {
				wg.Go(func() { results <- Migrate(ctx, cfg, service) })
			}
			wg.Wait()
			close(results)
			for e := range results {
				if e != nil {
					t.Fatal(e)
				}
			}
			conn := connect(t, service.Database, service.MigrationUser)
			want := 1
			if service.Name == "inventory" {
				want = 2
			}
			if got := scalar(t, conn, "SELECT count(*) FROM lab_migrations.versions"); got != want {
				t.Fatalf("versions=%d want=%d", got, want)
			}
		}
		before, e := Fingerprint(ctx, cfg)
		if e != nil {
			t.Fatal(e)
		}
		for _, service := range cfg.Services {
			if e = Migrate(ctx, cfg, service); e != nil {
				t.Fatal(e)
			}
		}
		after, e := Fingerprint(ctx, cfg)
		if e != nil || before != after {
			t.Fatal("repeat modified history", e)
		}
	})
	t.Run("role isolation and denied DDL", func(t *testing.T) {
		for _, service := range cfg.Services {
			for _, user := range []string{service.AppUser, service.MigrationUser} {
				for _, target := range cfg.Services {
					conn, e := Connect(ctx, cfg, target.Database, user)
					if target.Name == service.Name {
						if e != nil {
							t.Fatal(e)
						}
						Close(conn)
					} else {
						if conn != nil {
							Close(conn)
						}
						if e == nil || !strings.Contains(e.Error(), "42501") {
							t.Fatalf("cross-database access allowed or wrong error: %v", e)
						}
					}
				}
			}
			app := connect(t, service.Database, service.AppUser)
			for _, sql := range []string{"CREATE SCHEMA forbidden", "CREATE TABLE public.forbidden(id int)", "CREATE TEMP TABLE forbidden(id int)", "SET ROLE " + identifierSQL(service.MigrationUser), "DELETE FROM lab_migrations.versions"} {
				denied(t, app, sql)
			}
			if got := scalar(t, app, `SELECT count(*) FROM pg_roles WHERE rolname=current_user AND (rolsuper OR rolcreatedb OR rolcreaterole OR rolreplication OR rolbypassrls)`); got != 0 {
				t.Fatal("elevated app role")
			}
			migration := connect(t, service.Database, service.MigrationUser)
			table := identifierSQL(service.Name) + ".privilege_probe"
			if _, e := migration.Exec(ctx, "CREATE TABLE "+table+"(id int PRIMARY KEY)"); e != nil {
				t.Fatal(SafeError("probe create", e))
			}
			for _, sql := range []string{"INSERT INTO " + table + " VALUES(1)", "UPDATE " + table + " SET id=2", "DELETE FROM " + table + " WHERE id=2"} {
				if _, e := app.Exec(ctx, sql); e != nil {
					t.Fatal(SafeError("app DML", e))
				}
			}
			denied(t, app, "ALTER TABLE "+table+" ADD COLUMN forbidden int")
			denied(t, app, "TRUNCATE "+table)
			if _, e := migration.Exec(ctx, "DROP TABLE "+table); e != nil {
				t.Fatal(SafeError("probe drop", e))
			}
			for _, maintenance := range []string{cfg.AdminDatabase, "template1"} {
				conn, e := Connect(ctx, cfg, maintenance, service.AppUser)
				if conn != nil {
					Close(conn)
				}
				if e == nil || !strings.Contains(e.Error(), "42501") {
					t.Fatalf("maintenance access: %v", e)
				}
			}
		}
	})
	t.Run("failed migration rollback and checksum drift", func(t *testing.T) {
		service, _ := cfg.Service("inventory")
		conn := connect(t, service.Database, service.MigrationUser)
		files := fstest.MapFS{}
		entries, e := fs.ReadDir(migrations.Files, "inventory")
		if e != nil {
			t.Fatal(e)
		}
		for _, entry := range entries {
			data, e := fs.ReadFile(migrations.Files, "inventory/"+entry.Name())
			if e != nil {
				t.Fatal(e)
			}
			files["inventory/"+entry.Name()] = &fstest.MapFile{Data: data}
		}
		files["inventory/0003_failure.sql"] = &fstest.MapFile{Data: []byte("CREATE TABLE inventory.must_rollback(id int); SELECT 1/0;")}
		if e = migrate(ctx, cfg, service, files); e == nil || !strings.Contains(e.Error(), "22012") {
			t.Fatalf("expected failed migration: %v", e)
		}
		if scalar(t, conn, "SELECT count(*) FROM lab_migrations.versions") != 2 || scalar(t, conn, "SELECT count(*) FROM information_schema.tables WHERE table_schema='inventory' AND table_name='must_rollback'") != 0 {
			t.Fatal("failed migration persisted")
		}
		delete(files, "inventory/0003_failure.sql")
		files["inventory/0001_schema.sql"].Data = []byte("CREATE SCHEMA inventory; -- changed")
		if e = migrate(ctx, cfg, service, files); e == nil || !strings.Contains(e.Error(), "checksum") {
			t.Fatal("accepted changed migration", e)
		}
	})
	t.Run("repeat seed preserves catalog and modified stock", func(t *testing.T) {
		if e := Seed(ctx, cfg); e != nil {
			t.Fatal(e)
		}
		service, _ := cfg.Service("inventory")
		conn := connect(t, service.Database, service.AppUser)
		if scalar(t, conn, "SELECT count(*) FROM inventory.products") != 3 || scalar(t, conn, "SELECT available FROM inventory.stock WHERE sku='SKU-001'") != 100 {
			t.Fatal("unexpected initial catalog")
		}
		if _, e := conn.Exec(ctx, "UPDATE inventory.stock SET available=7 WHERE sku='SKU-001'; UPDATE inventory.products SET price_minor=99900 WHERE sku='SKU-001'"); e != nil {
			t.Fatal(SafeError("test update", e))
		}
		if e := Seed(ctx, cfg); e != nil {
			t.Fatal(e)
		}
		if scalar(t, conn, "SELECT count(*) FROM inventory.products") != 3 || scalar(t, conn, "SELECT available FROM inventory.stock WHERE sku='SKU-001'") != 7 || scalar(t, conn, "SELECT price_minor FROM inventory.products WHERE sku='SKU-001'") != 99900 {
			t.Fatal("repeat seed overwrote data")
		}
		for _, sql := range []string{"UPDATE inventory.stock SET available=-1 WHERE sku='SKU-001'", "UPDATE inventory.products SET price_minor=-1 WHERE sku='SKU-001'", "UPDATE inventory.products SET currency='rub' WHERE sku='SKU-001'", "INSERT INTO inventory.stock VALUES('missing',1)"} {
			if _, e := conn.Exec(ctx, sql); e == nil {
				t.Fatal("invalid inventory data accepted")
			}
		}
	})
}
```

<a id="p02-internal-database-fingerprint-go"></a>

## Приложение: internal/database/fingerprint.go

Полный исходный документ: `internal/database/fingerprint.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// Fingerprint covers ordered migration history and the complete small catalog.
// It is evidence of persistence for this foundation, not a general backup tool.
func Fingerprint(ctx context.Context, cfg Config) (string, error) {
	var values []string
	for _, service := range cfg.Services {
		conn, err := Connect(ctx, cfg, service.Database, service.MigrationUser)
		if err != nil {
			return "", err
		}
		var history string
		err = conn.QueryRow(ctx, `SELECT coalesce(json_agg(v ORDER BY version)::text,'[]') FROM lab_migrations.versions v`).Scan(&history)
		Close(conn)
		if err != nil {
			return "", SafeError("history fingerprint", err)
		}
		values = append(values, service.Name, history)
	}
	service, err := cfg.Service("inventory")
	if err != nil {
		return "", err
	}
	conn, err := Connect(ctx, cfg, service.Database, service.AppUser)
	if err != nil {
		return "", err
	}
	defer Close(conn)
	var catalog string
	err = conn.QueryRow(ctx, `SELECT coalesce(json_agg(v ORDER BY sku)::text,'[]') FROM (SELECT p.*,s.available FROM inventory.products p JOIN inventory.stock s USING(sku)) v`).Scan(&catalog)
	if err != nil {
		return "", SafeError("catalog fingerprint", err)
	}
	values = append(values, catalog)
	data, err := json.Marshal(values)
	if err != nil {
		return "", errors.New("cannot encode fingerprint")
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
```

<a id="p02-internal-database-config-test-go"></a>

## Приложение: internal/database/config_test.go

Полный исходный документ: `internal/database/config_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package database

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/jackc/pgx/v5/pgconn"
)

func testConfig(t *testing.T) Config {
	t.Helper()
	cfg, err := Load("../../configs/postgres/local.json", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}
func TestConfigValidation(t *testing.T) {
	for _, timeout := range []string{"0", "-1s", "500us", "61s", "invalid"} {
		t.Run(timeout, func(t *testing.T) {
			cfg := testConfig(t)
			cfg.Timeout = timeout
			if cfg.Validate() == nil {
				t.Fatal("accepted invalid budget")
			}
		})
	}
	cfg := testConfig(t)
	cfg.Services[0].AppUser = cfg.Services[1].AppUser
	if cfg.Validate() == nil {
		t.Fatal("accepted shared role")
	}
	cfg = testConfig(t)
	cfg.Services[0].Database = "db; DROP"
	if cfg.Validate() == nil {
		t.Fatal("accepted SQL identifier")
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"unknown":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path, t.TempDir()); err == nil {
		t.Fatal("accepted unknown settings")
	}
}
func TestSecretsPreservedAndMalformedRejected(t *testing.T) {
	cfg := testConfig(t)
	if err := GenerateSecrets(cfg); err != nil {
		t.Fatal(err)
	}
	first, err := cfg.Secret(cfg.AdminUser)
	if err != nil {
		t.Fatal(err)
	}
	if err = GenerateSecrets(cfg); err != nil {
		t.Fatal(err)
	}
	second, _ := cfg.Secret(cfg.AdminUser)
	if first != second {
		t.Fatal("credentials rotated unexpectedly")
	}
	if err = os.WriteFile(filepath.Join(cfg.Secrets, cfg.AdminUser+".password"), []byte("malformed-sensitive-value"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err = GenerateSecrets(cfg); err == nil || strings.Contains(err.Error(), "malformed-sensitive-value") {
		t.Fatal("malformed secret accepted or exposed")
	}
}
func TestSafeError(t *testing.T) {
	for _, err := range []error{&pgconn.PgError{Code: "42501", Message: "secret DSN", Detail: "password"}, errors.New("password"), errors.Join(context.Canceled, errors.New("password"))} {
		safe := SafeError("operation", err)
		if strings.Contains(safe.Error(), "password") || strings.Contains(safe.Error(), "DSN") {
			t.Fatal("sensitive diagnostics leaked")
		}
	}
	if !errors.Is(SafeError("operation", context.DeadlineExceeded), context.DeadlineExceeded) {
		t.Fatal("deadline lost")
	}
}
func TestMigrationFiles(t *testing.T) {
	for _, files := range []fstest.MapFS{
		{"orders/0000_invalid.sql": {Data: []byte("SELECT 1")}},
		{"orders/0001_a.sql": {Data: []byte("SELECT 1")}, "orders/0001_b.sql": {Data: []byte("SELECT 2")}},
		{"orders/unversioned.sql": {Data: []byte("SELECT 1")}},
		{"orders/0001_empty.sql": {Data: []byte{}}},
	} {
		if _, err := loadMigrations(files, "orders"); err == nil {
			t.Fatal("invalid migration set accepted")
		}
	}
	files := fstest.MapFS{"orders/0001_a.sql": {Data: []byte("SELECT 1")}}
	before, err := loadMigrations(files, "orders")
	if err != nil {
		t.Fatal(err)
	}
	files["orders/0001_a.sql"].Data = []byte("SELECT 2")
	after, err := loadMigrations(files, "orders")
	if err != nil {
		t.Fatal(err)
	}
	if before[0].checksum == after[0].checksum {
		t.Fatal("SQL edit did not alter checksum")
	}
}
```

<a id="p02-internal-database-config-go"></a>

## Приложение: internal/database/config.go

Полный исходный документ: `internal/database/config.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
// Package database prepares the lab's PostgreSQL databases without business RPCs.
package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Service assigns database and role ownership to one business service.
type Service struct {
	Name          string `json:"name"`
	Database      string `json:"database"`
	AppUser       string `json:"app_user"`
	MigrationUser string `json:"migration_user"`
}

// Config contains the explicit local connection profile; credentials stay in files.
type Config struct {
	Host          string    `json:"host"`
	Port          uint16    `json:"port"`
	SSLMode       string    `json:"sslmode"`
	Timeout       string    `json:"timeout"`
	AdminDatabase string    `json:"admin_database"`
	AdminUser     string    `json:"admin_user"`
	Services      []Service `json:"services"`
	Secrets       string    `json:"-"`
}

var identifier = regexp.MustCompile(`^[a-z][a-z0-9_]{0,62}$`)
var secretValue = regexp.MustCompile(`^[a-f0-9]{64}$`)

// Load reads a strict JSON profile and validates it before any connection.
func Load(path, secrets string) (Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return cfg, errors.New("cannot open database configuration")
	}
	defer func() { _ = file.Close() }()
	decoder := json.NewDecoder(io.LimitReader(file, 64<<10))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&cfg); err != nil {
		return cfg, errors.New("invalid database configuration")
	}
	var extra any
	if err = decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return cfg, errors.New("database configuration contains trailing data")
	}
	cfg.Secrets = secrets
	return cfg, cfg.Validate()
}

// Validate rejects incomplete profiles, shared identities and unbounded timeouts.
func (c Config) Validate() error {
	budget, err := time.ParseDuration(c.Timeout)
	if err != nil || budget < time.Millisecond || budget > time.Minute || c.Host == "" || c.Port == 0 || c.SSLMode != "disable" || c.Secrets == "" {
		return errors.New("database configuration requires host, port, local sslmode=disable, secrets directory and timeout in [1ms,1m]")
	}
	if !identifier.MatchString(c.AdminDatabase) || !identifier.MatchString(c.AdminUser) || len(c.Services) != 3 {
		return errors.New("invalid administrator or service configuration")
	}
	names, databases, users := map[string]bool{}, map[string]bool{c.AdminDatabase: true}, map[string]bool{c.AdminUser: true}
	for _, service := range c.Services {
		if (service.Name != "orders" && service.Name != "inventory" && service.Name != "payments") || names[service.Name] || !identifier.MatchString(service.Database) || databases[service.Database] {
			return errors.New("configuration must contain three distinct service databases")
		}
		names[service.Name], databases[service.Database] = true, true
		for _, user := range []string{service.AppUser, service.MigrationUser} {
			if !identifier.MatchString(user) || users[user] {
				return errors.New("database roles must be valid and distinct")
			}
			users[user] = true
		}
	}
	return nil
}

// Service finds one of the three configured service identities.
func (c Config) Service(name string) (Service, error) {
	for _, service := range c.Services {
		if service.Name == name {
			return service, nil
		}
	}
	return Service{}, errors.New("unknown service")
}

// Secret reads a generated hex password without including its value in errors.
func (c Config) Secret(user string) (string, error) {
	if !identifier.MatchString(user) {
		return "", errors.New("invalid secret identifier")
	}
	data, err := os.ReadFile(filepath.Join(c.Secrets, user+".password"))
	if err != nil || !secretValue.MatchString(strings.TrimSpace(string(data))) {
		return "", errors.New("database secret is unavailable or malformed")
	}
	return strings.TrimSpace(string(data)), nil
}

// GenerateSecrets preserves existing credentials. The private directory protects
// host access; readable files allow non-root containers to consume bind secrets.
func GenerateSecrets(cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.Secrets, 0o700); err != nil {
		return errors.New("cannot create private secret directory")
	}
	info, err := os.Lstat(cfg.Secrets)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("secret directory must be a regular directory")
	}
	if err = os.Chmod(cfg.Secrets, 0o700); err != nil {
		return errors.New("cannot protect secret directory")
	}
	users := []string{cfg.AdminUser}
	for _, service := range cfg.Services {
		users = append(users, service.AppUser, service.MigrationUser)
	}
	for _, user := range users {
		path := filepath.Join(cfg.Secrets, user+".password")
		if info, err := os.Lstat(path); err == nil {
			if !info.Mode().IsRegular() {
				return errors.New("secret must be a regular file")
			}
			if _, err = cfg.Secret(user); err != nil {
				return err
			}
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			return errors.New("cannot inspect secret file")
		}
		var value [32]byte
		if _, err := rand.Read(value[:]); err != nil {
			return errors.New("cannot generate database secret")
		}
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return errors.New("cannot create secret file")
		}
		_, writeErr := file.WriteString(hex.EncodeToString(value[:]) + "\n")
		if closeErr := file.Close(); errors.Join(writeErr, closeErr) != nil {
			return errors.New("cannot save secret file")
		}
	}
	return nil
}

// Connect uses explicit credentials, SCRAM authentication and bounded SQL waits.
func Connect(ctx context.Context, cfg Config, database, user string) (*pgx.Conn, error) {
	password, err := cfg.Secret(user)
	if err != nil {
		return nil, err
	}
	address := url.URL{Scheme: "postgres", Host: net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port))), Path: database, User: url.UserPassword(user, password)}
	query := url.Values{"sslmode": {cfg.SSLMode}, "require_auth": {"scram-sha-256"}}
	address.RawQuery = query.Encode()
	connection, err := pgx.ParseConfig(address.String())
	if err != nil {
		return nil, errors.New("cannot configure database connection")
	}
	connection.ConnectTimeout, _ = time.ParseDuration(cfg.Timeout)
	connection.RuntimeParams["statement_timeout"] = strconv.FormatInt(connection.ConnectTimeout.Milliseconds(), 10)
	connection.RuntimeParams["lock_timeout"] = strconv.FormatInt(connection.ConnectTimeout.Milliseconds(), 10)
	connection.RuntimeParams["timezone"] = "UTC"
	conn, err := pgx.ConnectConfig(ctx, connection)
	if err != nil {
		return nil, SafeError("database connection", err)
	}
	return conn, nil
}

// Close bounds connection cleanup independently of the operation's context.
func Close(conn *pgx.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = conn.Close(ctx)
}

// SafeError excludes SQL, parameters, connection strings and server details.
func SafeError(stage string, err error) error {
	var databaseError *pgconn.PgError
	if errors.As(err, &databaseError) {
		return fmt.Errorf("%s failed (SQLSTATE %s)", stage, databaseError.Code)
	}
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("%s: %w", stage, context.Canceled)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: %w", stage, context.DeadlineExceeded)
	}
	return fmt.Errorf("%s failed", stage)
}
```

<a id="p02-internal-database-bootstrap-go"></a>

## Приложение: internal/database/bootstrap.go

Полный исходный документ: `internal/database/bootstrap.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package database

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

func identifierSQL(value string) string { return pgx.Identifier{value}.Sanitize() }

// Bootstrap is repeatable after partial failure. CREATE DATABASE cannot be part
// of one transaction, so a session lock serializes this dedicated lab bootstrap.
func Bootstrap(ctx context.Context, cfg Config) error {
	admin, err := Connect(ctx, cfg, cfg.AdminDatabase, cfg.AdminUser)
	if err != nil {
		return err
	}
	defer Close(admin)
	if _, err = admin.Exec(ctx, "SELECT pg_advisory_lock(706020120)"); err != nil {
		return SafeError("bootstrap lock", err)
	}
	for _, service := range cfg.Services {
		for _, role := range []string{service.MigrationUser, service.AppUser} {
			var exists bool
			if err = admin.QueryRow(ctx, "SELECT EXISTS(SELECT FROM pg_roles WHERE rolname=$1)", role).Scan(&exists); err != nil {
				return SafeError("role lookup", err)
			}
			if !exists {
				if _, err = admin.Exec(ctx, "CREATE ROLE "+identifierSQL(role)); err != nil {
					return SafeError("role creation", err)
				}
			}
			password, err := cfg.Secret(role)
			if err != nil {
				return err
			}
			// Passwords are generated hex, and still quoted defensively. Server
			// statement/parameter logging is disabled by the lab deployment.
			statement := "ALTER ROLE " + identifierSQL(role) + " WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS NOINHERIT PASSWORD '" + strings.ReplaceAll(password, "'", "''") + "'"
			if _, err = admin.Exec(ctx, statement); err != nil {
				return SafeError("role configuration", err)
			}
			var memberships int
			if err = admin.QueryRow(ctx, "SELECT count(*) FROM pg_auth_members m JOIN pg_roles r ON r.oid=m.member WHERE r.rolname=$1", role).Scan(&memberships); err != nil {
				return SafeError("role membership check", err)
			}
			if memberships != 0 {
				return errors.New("managed roles must not inherit other roles")
			}
		}
		var owner string
		err = admin.QueryRow(ctx, "SELECT pg_get_userbyid(datdba) FROM pg_database WHERE datname=$1", service.Database).Scan(&owner)
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = admin.Exec(ctx, "CREATE DATABASE "+identifierSQL(service.Database)+" OWNER "+identifierSQL(service.MigrationUser))
			owner = service.MigrationUser
		}
		if err != nil {
			return SafeError("database preparation", err)
		}
		if owner != service.MigrationUser {
			return errors.New("existing database has an unexpected owner; no ownership takeover performed")
		}
		for _, statement := range []string{
			"REVOKE ALL ON DATABASE " + identifierSQL(service.Database) + " FROM PUBLIC",
			"GRANT CONNECT ON DATABASE " + identifierSQL(service.Database) + " TO " + identifierSQL(service.MigrationUser) + ", " + identifierSQL(service.AppUser),
		} {
			if _, err = admin.Exec(ctx, statement); err != nil {
				return SafeError("database privileges", err)
			}
		}
		conn, err := Connect(ctx, cfg, service.Database, cfg.AdminUser)
		if err != nil {
			return err
		}
		_, err = conn.Exec(ctx, "REVOKE ALL ON SCHEMA public FROM PUBLIC")
		Close(conn)
		if err != nil {
			return SafeError("public schema privileges", err)
		}
	}
	for _, database := range []string{cfg.AdminDatabase, "template1"} {
		if _, err = admin.Exec(ctx, "REVOKE ALL ON DATABASE "+identifierSQL(database)+" FROM PUBLIC"); err != nil {
			return SafeError("maintenance database privileges", err)
		}
	}
	return nil
}
```

<a id="p02-configs-postgres-local-json"></a>

## Приложение: configs/postgres/local.json

Полный исходный документ: `configs/postgres/local.json`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```json
{
  "host": "postgres",
  "port": 5432,
  "sslmode": "disable",
  "timeout": "10s",
  "admin_database": "postgres",
  "admin_user": "postgres",
  "services": [
    { "name": "orders", "database": "orders_db", "app_user": "orders_app", "migration_user": "orders_migrator" },
    { "name": "inventory", "database": "inventory_db", "app_user": "inventory_app", "migration_user": "inventory_migrator" },
    { "name": "payments", "database": "payments_db", "app_user": "payments_app", "migration_user": "payments_migrator" }
  ]
}
```

<a id="p02-cmd-labdata-main-test-go"></a>

## Приложение: cmd/labdata/main_test.go

Полный исходный документ: `cmd/labdata/main_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestOperationErrorsExitNonzero(t *testing.T) {
	var out, errs bytes.Buffer
	code := run(context.Background(), []string{"migrate", "-config", "../../configs/postgres/local.json", "-secrets", t.TempDir(), "-service", "unknown"}, &out, &errs)
	if code != 1 || !strings.Contains(errs.String(), "unknown service") || strings.Contains(out.String(), "PASS") {
		t.Fatalf("code=%d out=%q err=%q", code, out.String(), errs.String())
	}
}

func TestUsageAndSecrets(t *testing.T) {
	for _, test := range []struct {
		args []string
		code int
	}{
		{nil, 2}, {[]string{"unknown"}, 2}, {[]string{"secrets", "-help"}, 0},
		{[]string{"secrets", "-config", "../../configs/postgres/local.json", "-secrets", t.TempDir()}, 0},
	} {
		var out, errs bytes.Buffer
		if code := run(context.Background(), test.args, &out, &errs); code != test.code {
			t.Fatalf("args=%v code=%d err=%s", test.args, code, errs.String())
		}
	}
}
```

<a id="p02-cmd-labdata-main-go"></a>

## Приложение: cmd/labdata/main.go

Полный исходный документ: `cmd/labdata/main.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
// Command labdata prepares and verifies the isolated PostgreSQL foundation.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/database"
	"github.com/Includeoyasi/faultline-lab/internal/datalab"
)

func main() {
	os.Exit(mainCode())
}

func mainCode() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return run(ctx, os.Args[1:], os.Stdout, os.Stderr)
}
func setting(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func run(ctx context.Context, args []string, out, errs io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(errs, "usage: labdata {secrets|bootstrap|migrate|seed|fingerprint|smoke|reset} [flags]")
		return 2
	}
	flags := flag.NewFlagSet("labdata "+args[0], flag.ContinueOnError)
	flags.SetOutput(errs)
	config := flags.String("config", setting("FL_DATA_CONFIG", "configs/postgres/local.json"), "public database configuration")
	secrets := flags.String("secrets", setting("FL_DATA_SECRETS", ".local/postgres"), "private credentials directory")
	service := flags.String("service", "", "single migration service; empty applies all")
	project := flags.String("project", "", "isolated Compose project for explicit reset")
	confirm := flags.String("confirm", "", "repeat the project name to confirm data deletion")
	if err := flags.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		return 2
	}
	var err error
	switch args[0] {
	case "smoke":
		err = datalab.Run(ctx, out)
	case "reset":
		err = datalab.Reset(ctx, *project, *confirm, out)
	case "secrets", "bootstrap", "migrate", "seed", "fingerprint":
		var cfg database.Config
		cfg, err = database.Load(*config, *secrets)
		if err != nil {
			break
		}
		operationCtx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		switch args[0] {
		case "secrets":
			err = database.GenerateSecrets(cfg)
		case "bootstrap":
			err = database.Bootstrap(operationCtx, cfg)
		case "migrate":
			selected := cfg.Services
			if *service != "" {
				var item database.Service
				item, err = cfg.Service(*service)
				selected = []database.Service{item}
			}
			if err == nil {
				for _, item := range selected {
					if err = database.Migrate(operationCtx, cfg, item); err != nil {
						break
					}
				}
			}
		case "seed":
			err = database.Seed(operationCtx, cfg)
		case "fingerprint":
			var value string
			value, err = database.Fingerprint(operationCtx, cfg)
			if err == nil {
				_, err = fmt.Fprintln(out, value)
			}
		}
	default:
		_, _ = fmt.Fprintln(errs, "unknown command; use labdata without arguments for usage")
		return 2
	}
	if err != nil {
		_, _ = fmt.Fprintln(errs, err)
		return 1
	}
	if args[0] != "fingerprint" {
		_, _ = fmt.Fprintln(out, args[0]+": PASS")
	}
	return 0
}
```

<a id="p02-migrations-inventory-0002-catalog-sql"></a>

## Приложение: migrations/inventory/0002_catalog.sql

Полный исходный документ: `migrations/inventory/0002_catalog.sql`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```sql
CREATE TABLE inventory.products (
    sku text PRIMARY KEY CHECK (length(sku) BETWEEN 1 AND 64),
    title text NOT NULL CHECK (length(title) BETWEEN 1 AND 200),
    price_minor bigint NOT NULL CHECK (price_minor >= 0),
    currency text NOT NULL CHECK (currency ~ '^[A-Z]{3}$')
);

CREATE TABLE inventory.stock (
    sku text PRIMARY KEY REFERENCES inventory.products(sku),
    available bigint NOT NULL CHECK (available >= 0)
);
```

<a id="p02-migrations-inventory-0001-schema-sql"></a>

## Приложение: migrations/inventory/0001_schema.sql

Полный исходный документ: `migrations/inventory/0001_schema.sql`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```sql
CREATE SCHEMA inventory;
```

<a id="p02-migrations-embed-go"></a>

## Приложение: migrations/embed.go

Полный исходный документ: `migrations/embed.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
// Package migrations contains the reviewed, transactional database migrations.
package migrations

import "embed"

// Files is the immutable migration input shipped with the database tool.
//
//go:embed orders/*.sql inventory/*.sql payments/*.sql seed/*.sql
var Files embed.FS
```

<a id="p02-migrations-seed-inventory-sql"></a>

## Приложение: migrations/seed/inventory.sql

Полный исходный документ: `migrations/seed/inventory.sql`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```sql
INSERT INTO inventory.products (sku, title, price_minor, currency) VALUES
    ('SKU-001', 'Notebook', 19900, 'RUB'),
    ('SKU-002', 'Coffee', 45000, 'RUB'),
    ('SKU-003', 'Headphones', 399000, 'RUB')
ON CONFLICT (sku) DO NOTHING;

INSERT INTO inventory.stock (sku, available) VALUES
    ('SKU-001', 100), ('SKU-002', 50), ('SKU-003', 25)
ON CONFLICT (sku) DO NOTHING;
```

<a id="p02-migrations-orders-0001-schema-sql"></a>

## Приложение: migrations/orders/0001_schema.sql

Полный исходный документ: `migrations/orders/0001_schema.sql`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```sql
CREATE SCHEMA orders;
```

<a id="p02-migrations-payments-0001-schema-sql"></a>

## Приложение: migrations/payments/0001_schema.sql

Полный исходный документ: `migrations/payments/0001_schema.sql`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```sql
CREATE SCHEMA payments;
```

<a id="p02-internal-datalab-run-test-go"></a>

## Приложение: internal/datalab/run_test.go

Полный исходный документ: `internal/datalab/run_test.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package datalab

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/Includeoyasi/faultline-lab/internal/command"
)

func TestResetRejectsAmbiguousTargets(t *testing.T) {
	for _, test := range []struct{ project, confirm string }{
		{"", ""}, {"faultline-lab", "faultline-lab"}, {"foreign", "foreign"},
		{"faultline-data-local", ""}, {"faultline-data-local", "faultline-data-other"},
		{"faultline-data-../other", "faultline-data-../other"},
	} {
		if err := Reset(context.Background(), test.project, test.confirm, io.Discard); err == nil {
			t.Fatalf("accepted unsafe reset target %q", test.project)
		}
	}
}

func TestRunCleanupAndReports(t *testing.T) {
	for _, scenario := range []string{"success", "read-only-directory", "partial-build", "cancel-build", "cleanup-failure"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			config, e := os.ReadFile("../../configs/postgres/local.json")
			if e != nil {
				t.Fatal(e)
			}
			if e = os.MkdirAll(filepath.Join(root, "configs/postgres"), 0o700); e != nil {
				t.Fatal(e)
			}
			if e = os.WriteFile(filepath.Join(root, "configs/postgres/local.json"), config, 0o600); e != nil {
				t.Fatal(e)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			cleaned := false
			execute := func(commandCtx context.Context, program string, args []string) (command.Output, error) {
				if program != "docker" {
					t.Fatal("unexpected executable")
				}
				if slices.Contains(args, "build") && (scenario == "partial-build" || scenario == "cancel-build") {
					if scenario == "cancel-build" {
						cancel()
					}
					return command.Output{ExitCode: 1}, errors.New("controlled failure in " + root)
				}
				if slices.Contains(args, "down") {
					if scenario == "read-only-directory" {
						paths, e := filepath.Glob(filepath.Join(root, "artifacts/postgres/*/secrets"))
						if e != nil || len(paths) != 1 {
							t.Fatal("secret directory missing", e)
						}
						if e = os.Chmod(paths[0], 0o500); e != nil {
							t.Fatal(e)
						}
					}
					if commandCtx.Err() != nil {
						t.Fatal("cleanup reused canceled context")
					}
					if !slices.Contains(args, "--volumes") || !slices.Contains(args, "--rmi") || args[len(args)-1] != "local" {
						t.Fatal("wrong cleanup scope")
					}
					cleaned = true
					if scenario == "cleanup-failure" {
						return command.Output{ExitCode: 1}, errors.New("cleanup unavailable")
					}
				}
				return command.Output{}, nil
			}
			err := run(ctx, io.Discard, root, execute)
			wantSuccess := scenario == "success" || scenario == "read-only-directory"
			if !cleaned || wantSuccess != (err == nil) {
				t.Fatalf("cleanup=%t err=%v", cleaned, err)
			}
			paths, e := filepath.Glob(filepath.Join(root, "artifacts/postgres/*/report.json"))
			if e != nil || len(paths) != 1 {
				t.Fatal("report missing", e)
			}
			data, e := os.ReadFile(paths[0])
			if e != nil {
				t.Fatal(e)
			}
			var saved report
			if e = json.Unmarshal(data, &saved); e != nil {
				t.Fatal(e)
			}
			if strings.Contains(string(data), root) {
				t.Fatal("report exposes host path")
			}
			if _, e = os.Stat(filepath.Join(filepath.Dir(paths[0]), "secrets")); !errors.Is(e, os.ErrNotExist) {
				t.Fatal("temporary credentials remain")
			}
			if wantSuccess && saved.Result != "PASS" || !wantSuccess && saved.Result != "FAIL" {
				t.Fatal("wrong verdict")
			}
			if (scenario == "cleanup-failure") != (saved.Cleanup == "FAIL") {
				t.Fatal("wrong cleanup verdict")
			}
		})
	}
}
```

<a id="p02-internal-datalab-run-go"></a>

## Приложение: internal/datalab/run.go

Полный исходный документ: `internal/datalab/run.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
// Package datalab runs checks on its own disposable Compose project.
package datalab

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/command"
	"github.com/Includeoyasi/faultline-lab/internal/database"
)

var projectName = regexp.MustCompile(`^faultline-data-[a-z0-9][a-z0-9-]{0,47}$`)

func composeArgs(project string, args ...string) []string {
	return append([]string{"compose", "-f", "compose.yaml", "-f", "compose.postgres.yaml", "-p", project}, args...)
}

// Reset requires an exact explicit project confirmation; ordinary down preserves data.
func Reset(ctx context.Context, project, confirm string, out io.Writer) error {
	if !projectName.MatchString(project) || project != confirm {
		return errors.New("reset requires a faultline-data-* project and identical -confirm value")
	}
	root, err := os.Getwd()
	if err != nil {
		return errors.New("cannot locate checkout")
	}
	executor := command.Executor{Directory: root, Env: os.Environ()}
	bounded, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	result, err := executor.Execute(bounded, "docker", composeArgs(project, "down", "--volumes", "--remove-orphans", "--rmi", "local"))
	if err != nil {
		return errors.New("scoped reset failed; inspect Docker for the named project")
	}
	_, err = io.WriteString(out, result.Stdout)
	return err
}

type event struct {
	Stage    string `json:"stage"`
	Output   string `json:"output,omitempty"`
	Passed   bool   `json:"passed"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exit_code"`
}
type report struct {
	Project  string    `json:"project"`
	Started  time.Time `json:"started_at"`
	Finished time.Time `json:"finished_at"`
	Result   string    `json:"result"`
	Cleanup  string    `json:"cleanup"`
	Events   []event   `json:"events"`
}

func Run(ctx context.Context, out io.Writer) (finalErr error) {
	root, err := os.Getwd()
	if err != nil {
		return errors.New("cannot locate checkout")
	}
	return run(ctx, out, root, nil)
}

func run(ctx context.Context, out io.Writer, root string, execute func(context.Context, string, []string) (command.Output, error)) (finalErr error) {
	var err error
	var token [8]byte
	if _, err = rand.Read(token[:]); err != nil {
		return errors.New("cannot generate project ID")
	}
	project := "faultline-data-test-" + hex.EncodeToString(token[:])
	dir := filepath.Join(root, "artifacts", "postgres", project)
	if err = os.MkdirAll(dir, 0o700); err != nil {
		return errors.New("cannot create report directory")
	}
	cfg, err := database.Load(filepath.Join(root, "configs/postgres/local.json"), filepath.Join(dir, "secrets"))
	if err != nil {
		return err
	}
	removeCredentials := func() error {
		var result error
		// Windows/OneDrive can mark a bind-source directory read-only.
		// Restore write access only to this run's owned credential directory.
		if e := os.Chmod(cfg.Secrets, 0o700); e != nil && !errors.Is(e, os.ErrNotExist) {
			result = errors.Join(result, errors.New("cannot restore temporary directory permissions"))
		}
		for _, role := range append([]string{cfg.AdminUser}, roleNames(cfg)...) {
			if e := os.Remove(filepath.Join(cfg.Secrets, role+".password")); e != nil && !errors.Is(e, os.ErrNotExist) {
				result = errors.Join(result, errors.New("cannot remove temporary credential"))
			}
		}
		if e := os.Remove(cfg.Secrets); e != nil && !errors.Is(e, os.ErrNotExist) {
			result = errors.Join(result, errors.New("cannot remove temporary credential directory"))
		}
		return result
	}
	defer func() { finalErr = errors.Join(finalErr, removeCredentials()) }()
	if err = database.GenerateSecrets(cfg); err != nil {
		return err
	}
	pairs := []string{root, "[workspace]", filepath.ToSlash(root), "[workspace]"}
	if home, e := os.UserHomeDir(); e == nil && home != "" {
		pairs = append(pairs, home, "[home]", filepath.ToSlash(home), "[home]")
	}
	for _, role := range append([]string{cfg.AdminUser}, roleNames(cfg)...) {
		value, e := cfg.Secret(role)
		if e != nil {
			return e
		}
		pairs = append(pairs, value, "[secret]")
	}
	redact := strings.NewReplacer(pairs...)
	env := []string{}
	for _, item := range os.Environ() {
		key, _, _ := strings.Cut(item, "=")
		if !strings.EqualFold(key, "FL_PG_SECRETS_DIR") {
			env = append(env, item)
		}
	}
	env = append(env, "FL_PG_SECRETS_DIR="+cfg.Secrets)
	executor := command.Executor{Directory: root, Env: env}
	if execute == nil {
		execute = executor.Execute
	}
	evidence := report{Project: project, Started: time.Now().UTC(), Result: "FAIL", Cleanup: "PENDING", Events: []event{}}
	save := func() error {
		data, e := json.MarshalIndent(evidence, "", "  ")
		if e != nil {
			return e
		}
		return os.WriteFile(filepath.Join(dir, "report.json"), append(data, '\n'), 0o600)
	}
	if err = save(); err != nil {
		return errors.New("cannot save initial report")
	}
	call := func(parent context.Context, stage string, budget time.Duration, args ...string) (string, error) {
		_, _ = fmt.Fprintln(out, "PostgreSQL smoke:", stage)
		commandCtx, cancel := context.WithTimeout(parent, budget)
		defer cancel()
		result, e := execute(commandCtx, "docker", args)
		observation := event{Stage: stage, Output: redact.Replace(result.Stdout + result.Stderr), Passed: e == nil, ExitCode: result.ExitCode}
		if e != nil {
			observation.Error = redact.Replace(e.Error())
			if commandCtx.Err() != nil {
				observation.Error = redact.Replace(errors.Join(commandCtx.Err(), e).Error())
			}
			e = fmt.Errorf("%s failed: %s; see redacted report", stage, observation.Error)
		}
		evidence.Events = append(evidence.Events, observation)
		return result.Stdout, errors.Join(e, save())
	}
	mayExist := false
	defer func() {
		if mayExist {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			_, e := call(cleanupCtx, "cleanup", 60*time.Second, composeArgs(project, "down", "--volumes", "--remove-orphans", "--rmi", "local")...)
			if e == nil {
				for _, args := range [][]string{
					{"ps", "-aq", "--filter", "label=com.docker.compose.project=" + project},
					{"network", "ls", "-q", "--filter", "label=com.docker.compose.project=" + project},
					{"volume", "ls", "-q", "--filter", "label=com.docker.compose.project=" + project},
					{"image", "ls", "-q", "--filter", "reference=" + project + "-*:latest"},
				} {
					var remaining string
					remaining, e = call(cleanupCtx, "verify cleanup", 5*time.Second, args...)
					if e != nil || strings.TrimSpace(remaining) != "" {
						e = errors.New("owned resources remain after cleanup")
						break
					}
				}
			}
			evidence.Cleanup = "PASS"
			if e != nil {
				evidence.Cleanup = "FAIL"
			}
			finalErr = errors.Join(finalErr, e)
		} else {
			evidence.Cleanup = "NOT_NEEDED"
		}
		credentialErr := removeCredentials()
		if credentialErr != nil {
			evidence.Cleanup = "FAIL"
		}
		finalErr = errors.Join(finalErr, credentialErr)
		evidence.Finished = time.Now().UTC()
		if finalErr == nil {
			evidence.Result = "PASS"
		}
		finalErr = errors.Join(finalErr, save())
		_, _ = fmt.Fprintf(out, "%s: artifacts/postgres/%s/report.json\n", evidence.Result, project)
	}()
	if _, err = call(ctx, "prerequisites", 15*time.Second, "version", "--format", "{{.Server.Version}}"); err != nil {
		return err
	}
	mayExist = true
	if _, err = call(ctx, "build dbtool", 10*time.Minute, composeArgs(project, "build", "dbtool")...); err != nil {
		return err
	}
	if _, err = call(ctx, "start PostgreSQL", 90*time.Second, composeArgs(project, "up", "-d", "--wait", "--wait-timeout", "75", "postgres")...); err != nil {
		return err
	}
	defer func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = call(logCtx, "PostgreSQL logs", 5*time.Second, composeArgs(project, "logs", "--no-color", "postgres")...)
	}()
	tool := func(stage string, args ...string) (string, error) {
		return call(ctx, stage, 90*time.Second, composeArgs(project, append([]string{"run", "--rm", "--no-deps", "dbtool"}, args...)...)...)
	}
	for attempt := 0; attempt < 2; attempt++ {
		if _, err = tool("bootstrap", "bootstrap"); err != nil {
			return err
		}
	}
	if _, err = call(ctx, "rights, migrations and seed", 3*time.Minute, composeArgs(project, "run", "--rm", "--no-deps", "-e", "FL_DATA_INTEGRATION=1", "--entrypoint", "/app/database.test", "dbtool", "-test.run=^TestPostgresFoundation$", "-test.v", "-test.timeout=150s")...); err != nil {
		return err
	}
	before, err := tool("fingerprint before restart", "fingerprint")
	if err != nil {
		return err
	}
	if _, err = call(ctx, "restart PostgreSQL", 60*time.Second, composeArgs(project, "restart", "postgres")...); err != nil {
		return err
	}
	if _, err = call(ctx, "wait after restart", 60*time.Second, composeArgs(project, "up", "-d", "--wait", "--wait-timeout", "45", "postgres")...); err != nil {
		return err
	}
	for _, operation := range []string{"bootstrap", "migrate", "seed"} {
		if _, err = tool("repeat "+operation, operation); err != nil {
			return err
		}
	}
	after, err := tool("fingerprint after restart", "fingerprint")
	if err != nil {
		return err
	}
	if strings.TrimSpace(before) != strings.TrimSpace(after) {
		return errors.New("data or migration history changed across restart/repeat")
	}
	_, err = call(ctx, "resource snapshot", 15*time.Second, "stats", "--no-stream", "--format", "{{json .}}", project+"-postgres-1")
	return err
}
func roleNames(cfg database.Config) []string {
	var result []string
	for _, service := range cfg.Services {
		result = append(result, service.AppUser, service.MigrationUser)
	}
	return result
}
```

<a id="p02-compose-postgres-yaml"></a>

## Приложение: compose.postgres.yaml

Полный исходный документ: `compose.postgres.yaml`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```yaml
services:
  postgres:
    image: postgres:18.6-bookworm@sha256:afc7e2d441324c0388fa80c3d24f733b4194a4eb7f47dd8ee2b08eb1a24a647c
    user: "999:999"
    read_only: true
    cap_drop: [ALL]
    security_opt: [no-new-privileges:true]
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/postgres.password
      POSTGRES_INITDB_ARGS: --auth-host=scram-sha-256 --auth-local=scram-sha-256
    secrets:
      - source: postgres-password
        target: postgres.password
    command:
      - postgres
      - -c
      - shared_buffers=64MB
      - -c
      - max_connections=30
      - -c
      - log_statement=none
      - -c
      - log_min_error_statement=panic
      - -c
      - log_parameter_max_length=0
      - -c
      - log_parameter_max_length_on_error=0
    volumes:
      - postgres-data:/var/lib/postgresql
    tmpfs:
      - /var/run/postgresql:uid=999,gid=999,mode=0770
      - /tmp:mode=1777,size=32m
    mem_limit: 512m
    cpus: 0.5
    stop_grace_period: 30s
    networks: [data]
    healthcheck:
      test: [CMD, pg_isready, -U, postgres, -d, postgres]
      interval: 2s
      timeout: 2s
      retries: 30
      start_period: 5s
  dbtool:
    profiles: [tools]
    build:
      context: .
      dockerfile: build/Dockerfile
      target: labdata
    read_only: true
    cap_drop: [ALL]
    security_opt: [no-new-privileges:true]
    environment:
      FL_DATA_CONFIG: /etc/faultline/postgres.json
      FL_DATA_SECRETS: /run/secrets
    volumes:
      - ./configs/postgres/local.json:/etc/faultline/postgres.json:ro
    secrets:
      - {source: postgres-password, target: postgres.password}
      - {source: orders-app-password, target: orders_app.password}
      - {source: orders-migrator-password, target: orders_migrator.password}
      - {source: inventory-app-password, target: inventory_app.password}
      - {source: inventory-migrator-password, target: inventory_migrator.password}
      - {source: payments-app-password, target: payments_app.password}
      - {source: payments-migrator-password, target: payments_migrator.password}
    tmpfs: ["/tmp:mode=1777,size=32m"]
    networks: [data]
    mem_limit: 256m
    cpus: 1
volumes:
  postgres-data:
networks:
  data:
    internal: true
secrets:
  postgres-password:
    file: ${FL_PG_SECRETS_DIR:-.local/postgres}/postgres.password
  orders-app-password:
    file: ${FL_PG_SECRETS_DIR:-.local/postgres}/orders_app.password
  orders-migrator-password:
    file: ${FL_PG_SECRETS_DIR:-.local/postgres}/orders_migrator.password
  inventory-app-password:
    file: ${FL_PG_SECRETS_DIR:-.local/postgres}/inventory_app.password
  inventory-migrator-password:
    file: ${FL_PG_SECRETS_DIR:-.local/postgres}/inventory_migrator.password
  payments-app-password:
    file: ${FL_PG_SECRETS_DIR:-.local/postgres}/payments_app.password
  payments-migrator-password:
    file: ${FL_PG_SECRETS_DIR:-.local/postgres}/payments_migrator.password
```

<a id="p02-internal-smoke-exec-adapter-go"></a>

## Приложение: internal/smoke/exec_adapter.go

Полный исходный документ: `internal/smoke/exec_adapter.go`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

```go
package smoke

import (
	"context"

	"github.com/Includeoyasi/faultline-lab/internal/command"
)

type commandOutput struct {
	stdout, stderr string
	exitCode       int
}

type executeFunc func(context.Context, string, []string) (commandOutput, error)

type processExecutor struct {
	directory string
	env       []string
}

func (p processExecutor) execute(ctx context.Context, program string, args []string) (commandOutput, error) {
	output, err := (command.Executor{Directory: p.directory, Env: p.env}).Execute(ctx, program, args)
	return commandOutput{output.Stdout, output.Stderr, output.ExitCode}, err
}
```

<a id="p02-docs-development-postgres-md"></a>

## Приложение: docs/development/POSTGRES.md

Полный исходный документ: `docs/development/POSTGRES.md`. Текст включён непосредственно ниже; чтение исходного пути не требуется.

### PostgreSQL: подготовка данных и проверка изоляции

Основа данных для [задачи #12](https://github.com/Includeoyasi/faultline-lab/issues/12): один PostgreSQL, три отдельные базы, прикладная и миграционная роль для каждой базы, миграции и небольшой каталог Inventory. Команды выполняются из корня checkout. Фактические проверки и ограничения среды сохраняются в [отчёте](#p02-docs-development-evidence-p02-postgres-md).

#### Состав и взаимодействия

```mermaid
flowchart LR
    Operator[Оператор / CI] --> CLI[labdata на хосте]
    CLI -->|Docker CLI| Compose[Docker Compose]
    Compose --> Tool[Временный dbtool]
    Compose --> PG[PostgreSQL 18.6]
    Secrets[Локальные файлы паролей] -->|read-only secrets| Tool
    Secrets -->|только пароль postgres| PG
    Config[Публичный JSON конфигурации] --> Tool
    Tool -->|bootstrap: администратор| PG
    Tool -->|migrate: роль своего сервиса| PG
    Tool -->|seed: inventory_app| PG
    PG --> Volume[(Постоянный volume)]
    CLI --> Report[JSON-отчёт проверок]
```

`dbtool` и PostgreSQL соединены внутренней сетью `data`: порт базы на хост не публикуется. PostgreSQL работает от UID 999, dbtool — от UID 65532; файловая система контейнеров read-only, запись разрешена только в volume/tmpfs. У PostgreSQL лимит 512 MiB и 0.5 CPU; у временного dbtool — 256 MiB и 1 CPU. Это ограничения локального стенда, не измеренная производственная ёмкость.

Три бизнес-сервиса пока сохраняют свой gRPC-каркас. Подключение их бизнес-операций к базам будет следующим изменением. Kafka/Redis относятся к отдельной [задаче #16](https://github.com/Includeoyasi/faultline-lab/issues/16).

PostgreSQL добавляется через `compose.postgres.yaml`. Базовый `compose.yaml` и существующий smoke трёх сервисов запускаются независимо от подготовки паролей базы.

#### Кто что может делать

| Роль | Доступ | Назначение |
| :--- | :--- | :--- |
| `postgres` | Административный доступ к выделенному экземпляру | Создание баз и ролей, настройка прав; используется только инструментом подготовки |
| `orders_migrator` / `inventory_migrator` / `payments_migrator` | Владелец только своей базы | Создание и изменение схемы, журнал миграций |
| `orders_app` / `inventory_app` / `payments_app` | CONNECT к своей базе, USAGE своей схемы, SELECT/INSERT/UPDATE/DELETE бизнес-таблиц | Будущая работа приложений; Inventory seed использует эти же ограниченные права |

У прикладных ролей нет SUPERUSER, CREATEDB, CREATEROLE, REPLICATION, BYPASSRLS, членства в миграционных ролях и DDL/TEMP/TRUNCATE. Доступ к чужим базам, `postgres`, `template1` запрещён. Стандартные права PUBLIC отозваны явно. Журнал `lab_migrations.versions` находится вне бизнес-схемы и прикладной роли недоступен. Права на будущие таблицы выдаются через `ALTER DEFAULT PRIVILEGES` от миграционной роли.

Один экземпляр PostgreSQL выбран для экономии ресурсов личного компьютера. Базы логически изолированы, но делят CPU, память и отказ экземпляра. Это позволяет проверять права и миграции; для эксперимента с независимыми отказами самих БД потребуется отдельная конфигурация.

#### Подготовка стенда

Сначала собирается host-бинарник средствами контейнерного toolchain. Для Windows AMD64:

```powershell
docker compose build toolchain
docker compose run --rm -e GOOS=windows -e GOARCH=amd64 toolchain data-build
.\bin\labdata-windows-amd64.exe secrets
docker compose -f compose.yaml -f compose.postgres.yaml -p faultline-data-local build dbtool
docker compose -f compose.yaml -f compose.postgres.yaml -p faultline-data-local up -d --wait postgres
docker compose -f compose.yaml -f compose.postgres.yaml -p faultline-data-local run --rm --no-deps dbtool bootstrap
docker compose -f compose.yaml -f compose.postgres.yaml -p faultline-data-local run --rm --no-deps dbtool migrate
docker compose -f compose.yaml -f compose.postgres.yaml -p faultline-data-local run --rm --no-deps dbtool seed
```

Для macOS в команде сборки задаются `GOOS=darwin` и `GOARCH=arm64` либо `amd64`, для Linux — `GOOS=linux` и архитектура хоста. Имя результата: `bin/labdata-<os>-<arch>`; на Unix перед запуском при необходимости выполнить `chmod +x`. `data-cross-check` собирает все пять вариантов; компиляция не подтверждает работу Docker на конкретном компьютере. При установленном Go допустим `go run ./cmd/labdata <команда>`.

`secrets` создаёт семь криптографически случайных паролей в `.local/postgres`, без вывода значений. Повторный вызов сохраняет их; повреждённый существующий файл отклоняется. Конфигурация адреса, порта, имён баз и ролей — в [local.json](#p02-configs-postgres-local-json). Пароли не входят в конфигурацию. `FL_PG_SECRETS_DIR` позволяет передать Compose другой каталог; `labdata secrets -secrets <каталог>` должен использовать тот же путь. В контейнере настройки читаются через `FL_DATA_CONFIG` и `FL_DATA_SECRETS`.

На POSIX каталог паролей имеет режим 0700, файлы — 0644 для чтения непривилегированным контейнером; доступ с хоста ограничен каталогом. На Windows доступ зависит от ACL профиля пользователя: `chmod` не заменяет настройку Windows ACL. Эти файлы и артефакты исключены из Git, CI публикует только очищенный `report.json`.

Подготовка использует SCRAM-SHA-256. TLS отключён в конфигурации внутренней Docker-сети локального стенда. Для удалённого развёртывания потребуется отдельный профиль TLS, секретов и сетевого доступа; этот профиль для публикации PostgreSQL в интернет не предназначен. Серверное логирование SQL и параметров отключено, включая SQL при ошибках: это предотвращает попадание пароля из `ALTER ROLE` в журнал.

#### Порядок подготовки

```mermaid
sequenceDiagram
    actor Operator as Оператор
    participant Host as labdata на хосте
    participant DB as PostgreSQL
    participant Admin as dbtool bootstrap
    participant Migrator as dbtool migrate
    participant Seed as dbtool seed
    Operator->>Host: secrets
    Host-->>Operator: Сохранены локальные пароли
    Operator->>DB: compose up postgres
    Operator->>Admin: bootstrap
    Admin->>DB: Session lock → роли → базы → права
    Operator->>Migrator: migrate [по всем сервисам]
    Migrator->>DB: BEGIN → advisory lock → проверка истории
    Migrator->>DB: SQL + версии + SHA-256 → COMMIT
    Operator->>Seed: seed
    Seed->>DB: BEGIN → INSERT ON CONFLICT DO NOTHING → COMMIT
```

Bootstrap повторяем после частичного сбоя: создание баз не может быть одной транзакцией PostgreSQL. Session lock сериализует подготовку, существующая база с неожиданным владельцем отклоняется, данные не удаляются. Пароли не меняются при обычном повторе подготовки.

Миграции лежат в `migrations/<service>/NNNN_name.sql` и встроены в dbtool. Их версии, имена, SHA-256 и время применения записываются в каждой базе. Изменение уже применённого SQL останавливает команду; исправление оформляется следующей миграцией. Параллельные команды сериализуются advisory lock до чтения журнала. Все новые миграции одного сервиса и записи их истории применяются в одной транзакции; ошибка откатывает весь этот набор. Разные базы применяются независимо: ошибка Payments не откатывает уже подготовленную Inventory. Повторная команда завершает оставшуюся работу.

Это небольшой транзакционный механизм для SQL данного стенда. Миграции не должны содержать `BEGIN/COMMIT/ROLLBACK` или команды, требующие выполнения вне транзакции, например `CREATE INDEX CONCURRENTLY`. Автоматического down/удаления схемы нет. Для запуска одного сервиса: `dbtool migrate -service inventory`.

#### Что такое seed и зачем повторяемость

Здесь **seed — начальные тестовые данные**, а не случайное число:

| SKU | Товар | Цена в копейках | Валюта | Начальный остаток |
| :--- | :--- | ---: | :--- | ---: |
| SKU-001 | Notebook | 19900 | RUB | 100 |
| SKU-002 | Coffee | 45000 | RUB | 50 |
| SKU-003 | Headphones | 399000 | RUB | 25 |

Фиксированные SKU и значения дают известную исходную точку: тест может проверить конкретный товар, цену и остаток. Цены хранятся целыми числами в минимальных денежных единицах. Ограничения БД запрещают отрицательные цены/остатки, некорректную валюту и остаток без товара.

Повторяемость имеет два свойства:

1. На пустой базе seed создаёт одинаковый каталог.
2. На существующей базе повтор seed не создаёт дубликаты и не возвращает изменённые значения к исходным. Например, остаток 100 → 7 остаётся 7; изменённая цена также сохраняется. Это обеспечивается первичным ключом SKU и `ON CONFLICT DO NOTHING`.

**Повтор seed не равен сбросу эксперимента.** Для новой чистой исходной точки нужен отдельный набор данных/стенд или явный reset с повторной подготовкой. Нельзя запускать seed и считать уже использованную базу начальной: это исказит сравнение решений. Добавленный вручную товар не удаляется.

Числовой seed **будущего генератора нагрузки** решает другую задачу: вместе со стабильным логическим номером операции и версией алгоритма он задаёт воспроизводимый план операций. ID запуска хранится отдельно, план получает контрольную сумму — это согласованное решение F12 из [аудита #20](https://github.com/Includeoyasi/faultline-lab/issues/20), ещё не реализованное в #12. Воспроизводимый план сам по себе не гарантирует одинакового порядка конкурентного выполнения. Случайные пароли и ID изолированных прогонов, напротив, намеренно создаются заново.

#### Остановка и явное удаление данных

Штатная остановка сохраняет volume:

```text
docker compose -f compose.yaml -f compose.postgres.yaml -p faultline-data-local down --remove-orphans
```

Reset удаляет данные **указанного** выделенного проекта. Имя должно начинаться с `faultline-data-` и быть повторено в подтверждении:

```powershell
.\bin\labdata-windows-amd64.exe reset -project faultline-data-local -confirm faultline-data-local
```

Reset удаляет контейнеры, сети, volume и локальный образ dbtool данного проекта; общий PostgreSQL, toolchain и кеши сборки сохраняются. Команда ограничена минутой. Обычная остановка никогда не вызывает reset. Каталог `.local/postgres` не удаляется reset: существующие учётные данные пригодны для следующей подготовки. Для полной смены паролей нужен отдельный согласованный порядок ротации.

#### Автоматическая проверка

```powershell
.\bin\labdata-windows-amd64.exe smoke
```

Smoke создаёт свой `faultline-data-test-<случайный ID>`, временные пароли и отчёт до первого Docker-вызова. Он собирает dbtool, поднимает PostgreSQL, дважды запускает bootstrap и проверяет:

- конкурентные миграции с нуля, отсутствие повторных версий и неизменность истории при повторе;
- матрицу подключения прикладных/миграционных ролей к трём базам, запреты DDL/TEMP/SET ROLE/журнала, разрешённые DML и права будущих таблиц;
- откат ошибочного SQL вместе с записью версии и обнаружение изменённого checksum;
- исходный каталог, отсутствие дубликатов, сохранение изменённых цен и остатков при повторе seed, ограничения таблиц;
- совпадение контрольной суммы каталога, остатков и полной истории миграций до и после перезапуска и повторной подготовки.

Контрольная сумма здесь проверяет сохранность небольшого набора данных. Она сравнивается внутри одного набора данных: новый стенд имеет другие времена применения миграций и может получить другой отпечаток. Она не заменяет backup и не является будущей контрольной суммой плана нагрузки.

После успеха, ошибки или перехваченной отмены cleanup получает новый контекст на 60s и удаляет только ресурсы своего одноразового проекта, включая volume и образ dbtool. Затем проверяет отсутствие этих ресурсов и удаляет временные пароли. SIGKILL, авария ОС и недоступный Docker могут помешать очистке; для такого одноразового набора восстановление означает адресное удаление оставшихся ресурсов, а не восстановление данных. Для повторного Compose cleanup можно заново создать файлы паролей в каталоге прогона командой `secrets -secrets <каталог>` и передать его через `FL_PG_SECRETS_DIR`, затем вызвать reset с точным project из отчёта.

`artifacts/postgres/<project>/report.json` содержит этапы, exit codes, очищенные stdout/stderr и причины ошибок, логи PostgreSQL, отпечатки данных, Docker resource snapshot и итог cleanup. Пароли и домашние пути скрываются. Результаты проверок прав/seed записаны в выводе реальных Go integration tests. CI выполняет тот же host smoke на Linux AMD64 и ARM64; один успешный Windows-прогон не доказывает совместимость macOS.

На личном компьютере после проверок останавливаются Docker Desktop и его WSL VM. Образы отдельных smoke-прогонов не сохраняются; закреплённый PostgreSQL, актуальный toolchain и полезные общие кеши используются повторно.
