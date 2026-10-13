# Защита main и аккаунт исполнителя

Снимок на 13 октября 2026 года. Владелец проекта прямо поручил перейти на отдельный аккаунт исполнителя и настроить защиту main, затем предоставил логин AI-Your-Bro. Защита применена через GitHub API и проверена отдельным GET. Авторизация исполнителя проверена через GET /user: AI-Your-Bro, GitHub ID 342030453. Его права на репозиторий — write, push=true, admin=false; приглашение принято. Git Credential Manager сохранил OAuth-авторизацию в защищённом хранилище Windows. Владелец прямо отказался включать 2FA; этот выбор записан, а статус 2FA не выдан за подтверждённый API.

## Действующие настройки main

| Правило | Проверенное значение |
| :--- | :--- |
| Изменения через PR | Включено через обязательное независимое ревью |
| Число обязательных approval | 1 |
| Сброс approval после новых изменений | Включён |
| Approval после последнего push другим участником | Обязателен |
| Ревью CODEOWNERS | Включено; файл владельцев ещё не находится в main |
| Required check | Links and portable context |
| Источник required check | GitHub Actions, app ID 15368 |
| Актуальность ветки относительно main | Обязательна перед merge |
| Закрытие обсуждений | Обязательно |
| Применение правил к администратору | Включено |
| Force push / удаление | Запрещены |

Это проверка фактических настроек, а не попытка выполнить запрещённый push или удалить ветку. Подписанные коммиты и обязательная линейная история не вводились. Go checks добавляются после появления проверенного workflow в P01; отсутствующие проверки не объявляются действующими.

Подготовлен [.github/CODEOWNERS](../../.github/CODEOWNERS): все файлы, включая workflow и сам CODEOWNERS, назначены @Includeoyasi. Пока он не слит в main, действует требование одного независимого approval без привязки к этому файлу. После merge GitHub использует CODEOWNERS базовой ветки и требует ревью указанного владельца. [Правила CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners).

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
6. [PR #5](https://github.com/Includeoyasi/faultline-lab/pull/5) с CODEOWNERS и документами безопасности открыт от проверенного AI-Your-Bro; у Includeoyasi запрошено ревью. GitHub также подтвердил авторство и committer первого коммита fbbaef0 как AI-Your-Bro. Владелец ревьюирует конкретный HEAD и выполняет approval. Merge — только по его решению; ИИ сам не принимает свои изменения.
7. После merge проверить CODEOWNERS в main, реальные required checks, автора нового PR и обновить контекст. Основная реализация P01 пока TODO, 0/8; G1–G4 не приняты.

## Текущее ограничение

Защита main действует уже сейчас. CODEOWNERS и документы опубликованы в chore/executor-security и [PR #5](https://github.com/Includeoyasi/faultline-lab/pull/5) от проверенного AI-Your-Bro, на основе main 8c392b8. Новый PR и push от владельца не выполнялись. CODEOWNERS станет действующим в main после человеческого ревью и merge; текущие commit/PR и проверенные результаты находятся в CONTEXT_STATE.json. Основная реализация P01 ещё не начата.
