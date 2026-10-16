# PostgreSQL — проверка основы данных

Задача: [#12](https://github.com/Includeoyasi/faultline-lab/issues/12). Дата: 16 октября 2026 года, Europe/Moscow. Методика и команды — в [инструкции](../../POSTGRES.md).

Среда: Windows AMD64, Go 1.27.2, Docker Desktop 4.95 / Engine 29.8.2 / Compose 5.6, WSL2 с лимитами 2 CPU / 2 GiB RAM / 1 GiB swap. PostgreSQL 18.6-bookworm закреплён по [digest и списку платформ](postgres-image.json). Исходная база ветки — `843e95ed2b65a44729fbf297b986363c588d51ca`; локальные проверки выполнялись с изменениями рабочего дерева. Они не выдаются за прогон на чистом финальном коммите: его проверяет CI PR.

## Результаты

| Проверка | Результат |
| :--- | :--- |
| Контейнерный `toolchain check` | PASS: формат/imports, tidy/verify, build, lint без замечаний, unit/race и govulncheck без найденных уязвимостей. Доработанные Windows CLI/cleanup дополнительно проверены native Go tests; финальный HEAD проверяет CI. |
| Реальный PostgreSQL | PASS: конкурентные миграции, неизменность истории, положительная/отрицательная матрица ролей, DDL/TEMP/SET ROLE/journal запреты, default grants, откат ошибки, checksum drift, повтор seed и ограничения таблиц. |
| Полный Windows `labdata smoke` | PASS, [исходный очищенный отчёт](windows-postgres-report.json), project `faultline-data-test-90eda236ac29dee2`; 19:12:11–19:13:14 MSK. |
| Перезапуск и повторная подготовка | PASS: `bootstrap/migrate/seed` после restart сохранили полную историю и изменённый каталог/остатки. SHA-256 до/после: `6bd627a0b8eff5f313bc0f372ce2714f586efe3fed6d7bdbf7f78f826ac2d408`. |
| Явный reset | PASS: `faultline-data-local` удалён адресно; отдельно созданный `faultline-p02-cleanup-control` volume сохранился и затем удалён собственной командой. |
| Native Windows unit/process tests | PASS: `internal/command`, `internal/smoke`, `internal/datalab`, `cmd/labdata`; в том числе Job Object, отмена/ошибка partial build, неудачная очистка, read-only каталог и exit 1 при ошибке команды миграции. |
| Регрессия runtime smoke трёх сервисов | PASS, [отчёт](windows-runtime-report.json); health, изоляция, нормальное завершение, cleanup. Проверено отсутствие тегов образов своего прогона. |
| `data-cross-check` | PASS: Windows AMD64, Linux AMD64/ARM64, macOS AMD64/ARM64. Это сборки, не пять реальных runtime-сред. |
| Linux AMD64/ARM64 runtime и SIGINT | Выполняются workflow PR; результаты и артефакты в [GitHub Actions](https://github.com/Includeoyasi/faultline-lab/actions). Локальные Windows результаты их не заменяют. |
| Личный Mac | NOT_RUN: доступа к машине не было. |

Первый полный прогон выявил read-only атрибут пустого временного каталога credentials в OneDrive. Он завершился FAIL вместо ложного PASS. Исправление восстанавливает права только каталога своего прогона перед удалением семи известных файлов и каталога; регрессия покрывает этот случай. Успешный окончательный отчёт приведён выше. SQL/init ошибки также не скрывались: неверная YAML-запись tmpfs была исправлена до полного прогона.

## Ресурсы и завершение

В момент снимка PostgreSQL после функциональных тестов: **23.43 MiB / 512 MiB**, **5.12% CPU**. Это одиночный снимок из отчёта, не peak, baseline нагрузки или предел производительности; генератора бизнес-нагрузки ещё нет. Функциональная проверка укладывается в текущий WSL-бюджет; весь состав Kafka/Redis/наблюдаемости пока не измерен.

По замечанию владельца удалены 21 исторический smoke-образ, не использовавшийся контейнерами. Общие слои объясняют небольшой фактически освобождённый объём около 76 MB, несмотря на множество строк в UI. Исправленный runtime cleanup использует scoped `down --rmi local` и проверяет исчезновение тегов, включая отменённую/неудачную сборку. Fixture проверяет случай, когда Docker оставляет image, несмотря на успешный `down`.

[Подтверждение завершения](windows-cleanup.json): контейнеров 0, образов отдельных прогонов 0; остались актуальный toolchain и закреплённый PostgreSQL, volumes только для кешей Go. Docker Desktop остановлен, Docker WSL VM завершена, запущенных WSL-дистрибутивов и Docker-процессов нет; автозапуск выключен. Общий prune не использовался.
