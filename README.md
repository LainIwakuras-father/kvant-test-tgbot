# **ТЕСТОВОЕ ЗАДАНИЕ - ТЕЛЕГРАМ БОТ И HTTP API ДЛЯ ПЕРЕСЫЛКИ СООБЩЕНИЙ**

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white&style=for-the-badge" alt="Go"></a>
  <a href="https://gin-gonic.com/"><img src="https://img.shields.io/badge/Gin-00B386?logo=go&logoColor=white&style=for-the-badge" alt="Gin"></a>
  <a href="https://www.docker.com/"><img src="https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white&style=for-the-badge" alt="Docker"></a>
  <a><img src="https://img.shields.io/badge/REST%20API-FF6F00?logo=rest&logoColor=white&style=for-the-badge" alt="REST API"></a>
  <a href="https://swagger.io/"><img src="https://img.shields.io/badge/Swagger-85EA2D?logo=swagger&logoColor=black&style=for-the-badge" alt="Swagger"></a>
</p>

## Что делает приложение

Приложение состоит из двух частей, работающих одновременно:

- **Telegram-бот** (@ValentinkaBot или другое имя)
  - Приветствует нового пользователя командой `/start`
  - Сохраняет всех, кто написал боту (chat_id) в памяти
  - Подтверждает получение любого текстового сообщения от пользователя
- **HTTP API** (на Gin)
  - Защищённый секретным ключом (X-Secret-Key в заголовке)
  - Позволяет отправить сообщение **конкретному** пользователю по chat_id
  - Позволяет сделать **рассылку** (broadcast) всем пользователям, которые когда-либо взаимодействовали с ботом

- Эндпоинты:

## 🛠 Стек технологий
- **Язык программирования**: Go 1.25+
- **HTTP-фреймворк**: Gin
- **TelegramBot-фреймфорк**: go-telegram-bot-api
- **Контейнеризация**: Docker + Docker Compose
- **Документация**: Swagger (опционально)
- **Логгирование**: стандартная библиотека log/slog

---
 

## 📋 Системные требования

### Docker (Деплой)
- **Docker** — [Скачать](https://docs.docker.com/desktop/)

### Backend (Локальная разработка)
- **Go** `1.25.5+` — [Скачать](https://go.dev/dl/)

## Особенности разработки

- **Hot-reload** с помощью [air](https://github.com/air-verse/air)
- При изменении любого .go-файла сервер автоматически перезапускается
- Очень удобно при разработке бота и API одновременно

- Логи разделены по компонентам:
- `component=bot` — всё, что связано с Telegram
- `component=api` — запросы и обработка HTTP

- Используется **multi-stage Dockerfile**:
- dev-стадия → air + delve (debug)


## Как запустить


1. Склонируй репозиторий
2. Создай файл `.env` в корне проекта:

```env
BOT_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
SECRET_KEY=super-secret-token-2026
PORT=8080
DEBUG=true
```
3. Запустить через Docker 

```bash
docker-compose up -d 
```

---



## ⚠️ Важно

- Чтобы Docker работал, необходимо запустить **Docker Desktop**.
- Если y Docker возникают ошибки, попробуйте перезагрузить **Docker Desktop**.

---

### 🎉 Готово! Приложение запущено и готово к работе

## как оно выгладет
