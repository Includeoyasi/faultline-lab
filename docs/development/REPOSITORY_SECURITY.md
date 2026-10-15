# Защита main и аккаунт исполнителя

Снимок на 15 октября 2026 года; исходная настройка выполнена 13 октября. Владелец проекта прямо поручил перейти на отдельный аккаунт исполнителя и настроить защиту main, затем предоставил логин AI-Your-Bro. Защита применена через GitHub API и проверена отдельным GET. Авторизация исполнителя проверена через GET /user: AI-Your-Bro, GitHub ID 342030453. Его права на репозиторий — write, push=true, admin=false; приглашение принято. Git Credential Manager сохранил OAuth-авторизацию в защищённом хранилище Windows. Владелец прямо отказался включать 2FA; этот выбор записан, а статус 2FA не выдан за подтверждённый API.

## Действующие настройки main

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

[.github/CODEOWNERS](../../.github/CODEOWNERS) находится в main после слияния PR #5, подтверждённого 15 октября 2026 года (merge d616a22). Все файлы, включая workflow и сам CODEOWNERS, назначены @Includeoyasi. GitHub использует CODEOWNERS базовой ветки и требует ревью указанного владельца. [Правила CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners).

## Разделение ролей и авторизации

- Владелец: Includeoyasi, управление репозиторием, человеческое ревью и решение о merge.
- Исполнитель: [AI-Your-Bro](https://github.com/AI-Your-Bro), отдельный machine account под контролем владельца, с подтверждённой авторизацией и правами write без admin. Он публикует ветки и PR; ИИ не использует его для одобрения собственных изменений.
- Административная авторизация владельца использована только для прямо порученной настройки защиты. Новые PR от аккаунта владельца после этого поручения не создаются.

Git author name/email не меняют автора PR: GitHub определяет его по аутентифицированному аккаунту API/UI. Перед первой публикацией надо проверить реальную identity API и scoped Git credentials; изменение подписи коммита само по себе не завершает переход.

GitHub разрешает machine account, зарегистрированный человеком, но запрещает автоматическую регистрацию. Владелец принимает условия и контролирует email; пароль, recovery codes и токены не помещаются в Git, документы или чат. API текущей OAuth-авторизации не вернул поле two_factor_authentication. На предложение включить 2FA владелец ответил «я не хочу ее включать»; ИИ этот выбор не отменяет и не сообщает о включённой 2FA. [Условия аккаунтов GitHub](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service#3-account-requirements).

## Завершение подключения

1. Владелец предоставляет логин зарегистрированного вручную machine account. Выполнено: AI-Your-Bro, существование публичного профиля проверено.
2. Приглашение AI-Your-Bro в faultline-lab с правами write отправлено и принято (invitation ID 337527166). От имени исполнителя проверен GET collaborators/AI-Your-Bro/permission: write, push=true, admin=false.
3. OAuth через Git Credential Manager device flow выполнен; GET /user подтвердил AI-Your-Bro. Credentials находятся в Windows Credential Manager и не записываются в репозиторий. Администрация владельца и публикация исполнителя используют отдельные credentials; перед мутацией API проверяется фактический actor. Для работы с workflow нужны соответствующие права токена; администрация ради обычных PR не выдаётся.
4. Локальная Git identity этого репозитория: AI-Your-Bro, 342030453+AI-Your-Bro@users.noreply.github.com. Локальная credential.https://github.com.username выбирает AI-Your-Bro. Глобальная Git identity и авторизация других репозиториев не менялись. На macOS нужно отдельно подключить исполнителя и проверить identity; локальные credentials автоматически туда не переносятся.
5. [PR #4](https://github.com/Includeoyasi/faultline-lab/pull/4) уже слит: merge commit 8c392b8, принятый HEAD 83e0076. Переоформлять его не нужно; переход применяется к новым PR. Новая локальная ветка chore/executor-security обновлена до этого main.
6. [PR #5](https://github.com/Includeoyasi/faultline-lab/pull/5) слит владельцем: принятый HEAD 704a7ec, merge d616a22. ИИ не выполнял approve/merge. GitHub ранее подтвердил автора PR и коммиттера как AI-Your-Bro.
7. CODEOWNERS подтверждён в main, переносимый контекст обновлён. На 15 октября P01-01 в ревью; G1–G4 не приняты. Required checks перечислены выше и подтверждены серверным GET; одно наличие workflow не считается доказательством настройки.

## Текущее ограничение

Защита main и CODEOWNERS действуют. Авторизация новой Windows-копии или macOS требует проверки исполнителя; credentials не переносятся через Git. Включение 2FA не заявляется: действует зафиксированное решение владельца. Текущий PR реализации и фактические Go-проверки находятся в CONTEXT_STATE.json.
