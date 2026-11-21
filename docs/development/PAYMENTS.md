# Приём заказа на оплату

После подготовки секретов по [инструкции PostgreSQL](POSTGRES.md) запустите профиль с одним идентификатором набора данных. Пример для PowerShell:

```powershell
$env:FL_ORDERS_DATASET_ID = 'payment-demo'
$env:FL_INVENTORY_DATASET_ID = $env:FL_ORDERS_DATASET_ID
$stack = @('-p','payment-demo','-f','compose.yaml','-f','compose.postgres.yaml','-f','compose.infrastructure.yaml','-f','compose.inventory.yaml','-f','compose.orders.yaml','-f','compose.publishing.yaml','-f','compose.payments.yaml')
docker compose @stack up -d --wait postgres kafka
docker compose @stack run --rm dbtool bootstrap
docker compose @stack run --rm dbtool migrate
docker compose @stack run --rm dbtool seed
docker compose @stack up -d --build orders inventory payments
```

Для заказа с успешным резервом Payments читает подтверждённый Kafka-факт, проверяет принадлежность набору данных и сохраняет снимок суммы в минимальных денежных единицах. В таблице `payments.jobs` одна запись одновременно служит inbox и заданием. Повтор идентичного события безопасен. Конфликт ключа или испорченный контракт прерывает обработчик до подтверждения offset. Диагностику показывает `docker compose @stack logs payments`. Для несовместимого события потребуется исправление причины остановки.

Новая группа Payments начинает чтение с начала доступной истории топика. Это позволяет обработать накопленные заказы после первого запуска потребителя.

Симулятор вычисляет исход и задержку из seed и order_id. `FL_PAYMENTS_DECLINE_BPS` задаёт долю отказов в базисных пунктах (0 — все успешны, 10000 — все отклонены). `FL_PAYMENTS_MAX_DELAY` задаёт верхнюю границу задержки от 0 до часа. Результат выбора и срок исполнения сохраняются при приёме. Перезапуск или смена конфигурации касается лишь новых решений. Работник блокирует только готовое задание. Одна транзакция фиксирует состояние, при успехе — списание и при любом исходе — запись outbox.

Проверка в payments_db под ролью payments_app:

```sql
SELECT j.order_id, j.state, j.amount_minor, j.currency,
       d.amount_minor AS debited_minor, o.event_type, o.published_at
FROM payments.jobs j
LEFT JOIN payments.debits d USING (order_id)
LEFT JOIN payments.outbox o USING (order_id);
```

Для SUCCEEDED сумма списания совпадает с суммой резерва. Для DECLINED поле debited_minor пусто. Заполненное published_at означает, что relay получил подтверждение Kafka и отметил событие. Топик результата — `payments.events.v1`. Orders сохраняет исход и продолжает работу со складом, как описано в [руководстве Orders](ORDERS.md).
