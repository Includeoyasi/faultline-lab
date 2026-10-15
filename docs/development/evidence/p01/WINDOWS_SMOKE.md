# P01: проверка контейнеров на Windows

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

## Исходные данные

- [report.json](windows-2026-10-15/report.json): команды, время, исходы, source revision и реальные image IDs.
- Логи: [Orders](windows-2026-10-15/orders.jsonl), [Inventory](windows-2026-10-15/inventory.jsonl), [Payments](windows-2026-10-15/payments.jsonl).
- [environment.json](windows-2026-10-15/environment.json): наблюдённый backend и бюджет.
- [cleanup.json](windows-2026-10-15/cleanup.json): подтверждение остановки среды.
- SHA256 report.json: `b08d9deabadc65c765d9173482e9171782aeff28affb5e7fd9bf0b201d5add65`.

## Состояние личного компьютера после проверки

Тестовые контейнеры и сеть удалены. Docker Desktop штатно остановлен, выполнен wsl --shutdown; процессов Docker и vmmemWSL при финальной проверке нет. AutoStart=false, запись Docker Desktop из автозагрузки Windows удалена. Go, Docker, WSL, образы и build cache сохранены для дальнейшей работы. Перезагрузка не потребовалась и не выполнялась.

При следующих проверках запускайте среду по необходимости и освобождайте ресурсы после окончания. Инструкция: [контейнерный каркас](../../CONTAINER_RUNTIME.md).

Этот smoke подтверждает lifecycle каркаса. Бизнес-операции, базы, пропускная способность и пиковое потребление под нагрузкой не измерялись. Личный Mac остаётся NOT_RUN. Результат не означает принятия P01 или G1 владельцем.
