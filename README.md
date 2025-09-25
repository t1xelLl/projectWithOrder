# Project With Order 🚀

Микросервис для обработки заказов с использованием Kafka, PostgreSQL и Redis. Система принимает заказы через Kafka, сохраняет в базу данных и предоставляет REST API для доступа к данным.

## 📋 Описание проекта

Проект представляет собой микросервисную архитектуру для обработки заказов в реальном времени:

- **Consumer**: Слушает сообщения из Kafka и сохраняет заказы в БД
- **REST API**: Предоставляет HTTP endpoints для доступа к данным заказов
- **Кэширование**: Автоматическое кэширование в Redis для ускорения доступа
- **Восстановление кэша**: Автоматическая загрузка данных в кэш при запуске

## 🏗️ Архитектура
Kafka Producer → Kafka Topic → Go Consumer → PostgreSQL → Redis → REST API

## 🛠️ Технологии

- **Backend**: Go 1.24
- **База данных**: PostgreSQL 15
- **Кэш**: Redis 7.2
- **Брокер сообщений**: Apache Kafka
- **HTTP сервер**: Gin Framework
- **Миграции**: Migrate
- **Контейнеризация**: Docker + Docker Compose
- **Логирование**: Logrus

## ⚙️ Установка и запуск

### Предварительные требования

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.24+ (для локальной разработки)

1. **Клонируйте репозиторий**:
```bash
git clone https://github.com/t1xelLl/projectWithOrder.git
cd projectWithOrder
```
2. **Установите необходимые пакеты Go**
```bash
go mod tidy
```
3. **Запуск приложения**
```bash
docker compose up -d 
```
4. **Запуск Kafka Producer**
``` bash
docker-compose exec kafka kafka-console-producer --broker-list localhost:9092 --topic order
```
5. **Доступ к API**: API будет доступен по адресу http://localhost:8080 .

## Демо видео
