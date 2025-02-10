# Post-comment-system

## Выбор хранилища

Для выбора inmemory хранилища требуется передать флаг `-storage=inmemory`, для PostgreSQL следует передать `-storage=postgres`. По умолчанию в docker-compose стоит флаг `-storage=postgres`

## Запуск
1. Создаем .env, пример можно взять из .env.example
2. В docker-compose проверяем, что выбрано нужно нам хранилище
3. Запускаем контейнеры с помощью `docker-compose up -d`
4. Радуемся =)

## Структура
```
+---graph                                        # Сгенерированные файлы и модели graphql
+---internal                                     # Файлы проекта
|   +---repository                               # Репозиторий для управления сущностями
|   |   |   comment_repository.go                # интерфейс для взаимодействия с комментариями
|   |   |   post_repository.go                   # интерфейс для взаимодействия с постами
|   |   |
|   |   +---inmemory                             # имплементация интерфейсов репозитория для inmemory хранилища
|   |   |       comment_repo.go
|   |   |       post_repo.go
|   |   |
|   |   \---postgres                             # имплементация интерфейса репозитория для postgresql хранилища
|   |           comment_repo.go
|   |           post_repo.go
|   |
|   +---service                                  # Сервисный слой с бизнес логикой
|   |   +---comment
|   |   |       comment_service.go
|   |   |
|   |   +---post
|   |   |       post_service.go
|   |   |
|   |   \---subscriber_manager                   # Сервис для отправления уведомления о новых сообщениях всем подписчикам
|   |           manager.go
|   |
|   \---storage
|       +---inmemory                             # Реализация inmemory хранилища
|       |       storage.go
|       |
|       \---postgres                             # Реализация подключения к postgresql хранилищу
|           |   storage.go
|           |
|           \---migrations
|                   V0001__init.sql
|                   V0002__add_users.sql
|
\---tests
    +---inmemory                                 # тесты для inmemory хранилища
    |       inmemory_comment_test.go
    |       inmemory_post_test.go
    |
    \---postgres                                 # тесты для postgresql хранилища
            postgres_comment_test.go
            postgres_post_test.go
```
