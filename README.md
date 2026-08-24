# Go Backend Practice

Репозиторий, в котором я собираю решения практических задач и прикладных паттернов по ходу изучения языка Go.

Репозиторий разделен на два основных направления:
1. `services/` — курсовые домашние задания и законченные учебные мини-сервисы.
2. `patterns/` — отработка изолированных бэкенд-паттернов, сетевых ручек (`net/http`) и Concurrency-пайплайнов для секций лайвкодинга.

Здесь нет текстов заданий (в целях соблюдения авторских прав), только **мой код** и описание реализации.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Status](https://img.shields.io/badge/Status-Active_Learning-green)

---

## 📂 Структура репозитория

```text
go-backend-practice/
├── README.md
├── services/                  # Курсовые проекты и мини-сервисы
│   └── task-manager-bot/      # Telegram-бот на Webhooks с in-memory состоянием
└── patterns/                  # Бэкенд-паттерны и livecoding задачи
```

---

## 🚀 Сервисы (`services/`)

### 1. [Task Manager Bot](./services/task-manager-bot)
**Telegram-бот для управления задачами.**

Полноценная реализация бэкенда для бота. Сервис работает через **Webhooks**, хранит состояние в памяти и обрабатывает команды конкурентно.

* **Цель:** Настроить взаимодействие с внешним API через **Webhooks** (вместо Long Polling), разобраться с пакетом `net/http` и обеспечить безопасный доступ к данным в памяти (`sync.RWMutex`).
* **Функционал:** Создание, назначение и выполнение задач.

[🔗 Перейти к коду](./services/task-manager-bot)

---

## ⚡ Практика паттернов (`patterns/`)

Отработка прикладных сценариев для highload-нагрузок и собеседований:
* **HTTP & Web:** Чистый `net/http` REST API, безопасная пагинация, Middleware-цепочки, Graceful Shutdown.
* **Concurrency & Pipelines:** Worker Pool, Fan-In/Fan-Out, семафоры, отмена контекстов (`context.WithCancel`).
* **Highload In-Memory:** Фоновые воркеры под 10k RPS, TTL-кэши с автоочисткой, Rate Limiting (Token Bucket), Singleflight.

---

## 🛠 Ключевой фокус изучения

* **Concurrency:** Безопасная работа с данными в многопоточной среде (Goroutines, Channels, Mutexes, Atomics, Data Race avoidance).
* **Architecture & Clean Code:** Разделение логики на слои, обработка ошибок, отсутствие утечек горутин и памяти.
* **Web & Networking:** Эффективная работа со стандартной библиотекой Go, контексты и таймауты.
