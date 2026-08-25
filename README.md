<div align="center">

# 🎓 MentorHub

**Платформа сопровождения первокурсников университета KBTU**

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org)
[![Angular](https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular)](https://angular.io)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker)](https://docker.com)

[🚀 Быстрый старт](#-быстрый-старт) · [📖 API](#-api-документация) · [🏗 Архитектура](#-архитектура) · [👥 Роли](#-роли-системы)

</div>

---

## 📌 О проекте

**MentorHub** — информационная система для организации адаптации первокурсников в KBTU. Платформа автоматизирует назначение задач, отслеживание прогресса, проведение встреч и коммуникацию между четырьмя категориями участников.

### Ключевые возможности

- 📋 **Управление задачами** — Head создаёт шаблоны задач и массово назначает их студентам; ментор проверяет выполнение
- 📅 **Встречи и объявления** — Mentor планирует встречи с группой и публикует объявления
- 🔔 **Уведомления** — автоматические in-app уведомления при изменении статусов задач
- ❓ **Q&A** — первокурсники задают вопросы своему ментору прямо в системе
- 📊 **Аналитика** — дашборды для каждой роли с актуальной статистикой
- 🚨 **Жалобы** — эдвайзеры подают жалобы, Head рассматривает и выносит решение
- 📚 **FAQ** — база знаний для первокурсников

---

## 👥 Роли системы

| Роль | Описание |
|------|----------|
| 👑 **Head** | Руководитель. Управляет всеми пользователями, учебными годами, факультетами, шаблонами задач и жалобами |
| 📊 **Advisor** | Эдвайзер. Надзирает за менторами, отслеживает прогресс студентов, отправляет напоминания |
| 🎓 **Mentor** | Ментор. Ведёт группу первокурсников, проверяет задачи, проводит встречи |
| 🧑‍🎓 **Freshman** | Первокурсник. Выполняет задачи, читает объявления, общается с ментором |

---

## 🛠 Технологический стек

### Backend
- **Go 1.22** — основной язык
- **chi v5** — HTTP-роутер
- **pgx/v5** — драйвер PostgreSQL
- **golang-jwt/jwt v5** — JWT аутентификация
- **zerolog** — структурированное логирование
- **viper** — конфигурация через `.env`
- **swaggo/swag** — Swagger UI документация

### База данных
- **PostgreSQL 16** — основная СУБД
- **golang-migrate** — управление миграциями

### Frontend
- **Angular** — SPA-фронтенд
- **Nginx** — production-сервер

### Инфраструктура
- **Docker + Docker Compose** — контейнеризация
- **Multi-stage Dockerfile** — оптимизированный образ (~10 MB)

---

## 🚀 Быстрый старт

### Требования

- [Docker](https://docs.docker.com/get-docker/) и Docker Compose
- [Go 1.22+](https://golang.org/dl/) (для локальной разработки)
- [golang-migrate](https://github.com/golang-migrate/migrate) (для локальных миграций)

### 1. Клонировать репозиторий

```bash
git clone https://github.com/ayekaa06/mentor-hub-KBTU.git
cd mentor-hub-KBTU
```

### 2. Настроить переменные окружения

```bash
cp .env.example .env
# Отредактируйте .env при необходимости
```

Пример `.env`:
```env
SERVER_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_NAME=mentorhub_db
DB_USER=mentorhub
DB_PASSWORD=mentorhub_secret

JWT_SECRET=your-super-secret-key-change-in-production
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
```

### 3. Запустить через Docker Compose (рекомендуется)

```bash
# Поднять всё: БД + Backend + Frontend
docker compose up -d

# Применить миграции (первый запуск)
docker compose --profile migrate run --rm mentorhub-migrate
```

После запуска доступно:
- 🌐 **Фронтенд:** http://localhost:8082
- 📖 **Swagger UI:** http://localhost:8082/swagger/index.html
- ❤️ **Health check:** http://localhost:8080/health

### 4. Локальная разработка (без Docker)

```bash
# Поднять только PostgreSQL
make docker-up

# Применить миграции
make migrate-up

# Запустить сервер
make run
```

### 5. Заполнить тестовыми данными

```bash
go run ./cmd/seed
```

---

## 🔑 Тестовые аккаунты

После запуска seed-инструмента:

| Роль | Email | Пароль |
|------|-------|--------|
| Head | `head@mentorhub.com` | `password123` |
| Advisor | `advisor@mentorhub.com` | `password123` |
| Mentor | `mentor@mentorhub.com` | `password123` |
| Freshman 1–10 | `freshman1@mentorhub.com` ... | `password123` |

Из seed-миграции `000004` (пароль `Test1234!`):
- Mentor: `arman.seitkali@mentorhub.kz`
- Advisor: `gulnar.bekova@mentorhub.kz`
- Freshmen: `aibek.dzhaksybekov@mentorhub.kz` ... `aliya.suleimenova@mentorhub.kz`

---

## 📖 API Документация

Swagger UI: **http://localhost:8082/swagger/index.html**

### Пример — логин и запрос

```bash
# Получить токен
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"head@mentorhub.com","password":"password123"}'

# Запрос с токеном
curl http://localhost:8080/api/v1/head/dashboard \
  -H "Authorization: Bearer <access_token>"
```

### Группы эндпоинтов

| Префикс | Роль | Примеры |
|---------|------|---------|
| `/api/v1/auth/*` | Все | login, register, refresh, logout, me |
| `/api/v1/head/*` | Head | users, task-templates, analytics, complaints |
| `/api/v1/advisor/*` | Advisor | mentors, analytics, reminders |
| `/api/v1/mentor/*` | Mentor | group, tasks, meetings, questions |
| `/api/v1/freshman/*` | Freshman | tasks, meetings, faq, notifications |

---

## 🏗 Архитектура

```
MentorHub/
├── cmd/
│   ├── server/         # Точка входа HTTP-сервера
│   └── seed/           # CLI seed-инструмент
├── internal/
│   ├── config/         # Конфигурация (viper)
│   ├── domain/         # Бизнес-сущности
│   ├── dto/            # Request/Response структуры
│   ├── handler/        # HTTP handlers, по ролям
│   ├── middleware/     # Auth JWT, RBAC, CORS, Logger
│   ├── pkg/
│   │   ├── hasher/     # bcrypt обёртка
│   │   ├── jwt/        # JWT Manager
│   │   └── response/   # Унифицированные JSON-ответы
│   ├── repository/
│   │   ├── interfaces.go
│   │   └── postgres/   # pgx реализации
│   └── service/        # Бизнес-логика, по ролям
├── migrations/         # SQL миграции (up/down)
├── docs/               # Swagger (автогенерация)
├── Practice-master/    # Angular frontend
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

Бэкенд следует **Layered Architecture**:

```
Handler → Service → Repository → Domain
```

---

## 🗄 Миграции

```bash
make migrate-up                        # Применить все
make migrate-down                      # Откатить последнюю
make migrate-create name=<name>        # Создать новую
```

| Версия | Содержание |
|--------|------------|
| `000001` | users, tokens, academic_years, faculties, groups |
| `000002` | task_templates, tasks, meetings, announcements |
| `000003` | notifications, complaints, faq_items, questions |
| `000004` | Seed data: 10 первокурсников, 4 шаблона задач |

---

## ⚙️ Makefile

```bash
make run            # Dev-сервер
make build          # Бинарник в ./bin/
make test           # Тесты + покрытие
make lint           # golangci-lint

make docker-up      # PostgreSQL
make docker-all     # Всё
make docker-down    # Остановить
make docker-logs    # Логи backend
```

---

## 🔒 Безопасность

- Пароли: **bcrypt** cost 12
- Access token TTL: **15 минут**, Refresh token: **7 дней**
- Refresh token хранится в БД и инвалидируется при logout
- `.env` исключён из git (`.gitignore`)
- RBAC middleware блокирует доступ к чужим ресурсам

---

## 📜 Лицензия

MIT License — свободное использование в образовательных целях.

---

<div align="center">
Разработано в рамках учебной практики <strong>KBTU 2025–2026</strong>
</div>
