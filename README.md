# Мануал по деплою бэкенда

## 1. Общая информация

### Назначение

Данный документ описывает процесс сборки и деплоя бэкенд-приложения в целевое окружение.

### Проект

Бэкенд-приложение: dokkee-backend

### Технологический стек

Go

PostgreSQL

Docker

Docker Compose

migrate (для миграций)

Linux Ubuntu 22.04 (для серверов) 

## 2. Описание окружения

### 2.1 Локальное окружение разработчика

Используется для разработки и отладки приложения.

Параметры:

ОС: Windows

Go: (версия из go.mod)

PostgreSQL: через Docker

Команда запуска:

make env-up
make migrate-up
make backend-build
make backend-up

Фронтенд должен обращаться на localhost:8000(либо можно сменить в докеркомпозе).

### 2.2 Docker-окружение

Бэкенд теперь может запускаться в Docker-контейнере в одной сети с PostgreSQL.

Команды запуска:

make env-up
make migrate-up
make backend-build
make backend-up

Бэкенд будет доступен на порту из переменной PORT (по умолчанию 8000).

## 3. Тестирование

Проект включает unit-тесты для основных компонентов.

### Запуск тестов

go test ./...

Тесты покрывают:
- Сервис авторизации (создание пользователя, генерация/парсинг токенов)
- Репозиторий (работа с БД через mocks)
- HTTP-хендлеры (sign-up, sign-in)

Тесты используют testify для mocks и sqlmock для БД.

## 4. Требования для деплоя

Перед началом деплоя необходимо убедиться, что на сервере установлены:

Docker

Docker Compose

PostgreSQL (или через Docker)

Go (если сборка на сервере)

Git

Проверочные команды:

docker --version
docker compose version
psql --version (или docker exec)
git --version

## 5. Структура проекта

<img width="235" height="343" alt="image" src="https://github.com/user-attachments/assets/31d206f4-1093-49fd-96d2-106bbcdaa7d5" />


## 6. Docker-конфигурация

Проект включает Dockerfile для бэкенда и docker-compose.yaml для управления сервисами.

### Dockerfile

Multistage-сборка на основе golang:1.25-alpine и alpine:latest.

Аргументы:
- PORT: порт для EXPOSE (по умолчанию 8000)

### docker-compose.yaml

Сервисы:
- dokkee-postgres: PostgreSQL
- dokkee-backend: бэкенд-приложение
- dokkee-postgres-migrate: для миграций
- port-forwarder: для локального доступа к БД

Переменные окружения:
- PORT: порт бэкенда (дефолт 8000)
- POSTGRES_*: для подключения к БД
- salt, sign_key: для JWT

### Команды Makefile для Docker

- make backend-build: пересобрать образ бэкенда
- make backend-up: запустить бэкенд
- make backend-down: остановить бэкенд
- make backend-logs: посмотреть логи

## 7. Процесс сборки приложения

### 7.1 Установка зависимостей

go mod download

### 7.2 Сборка проекта

go build -o bin/server cmd/main.go

## 8. Конфигурация

Пример config.yml:

port: "8000"

Обязательно создать .env файл по примеру .env.example с необходимыми переменными окружения:
- POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB: для подключения к БД
- salt, sign_key: для генерации паролей и токенов


Для Docker: POSTGRES_HOST=dokkee-postgres (внутри сети).


## 9. Возможные проблемы

Проблема: БД не подключается

Причина: неверный config.yml или порт не открыт.

Проблема: Тесты не проходят

Причина: отсутствуют зависимости (go mod tidy), или mocks не настроены.

Проблема: Docker-контейнер не стартует

Причина: переменные окружения не заданы в .env, или порт занят.

Решение: проверить config.yml, make env-port-forward.

Проблема: Приложение падает

Причина: ошибка в коде или конфиге.

Решение: docker logs dokkee-backend

## 10. Ответственные

Разработчик бэкенда: Айрат

РИДМИ будет дополняться по мере разработки.

