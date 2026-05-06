# LoL Tracker

## Общая идея

Локальный клиент, который:

- читает данные из League of Legends Live Client API
- преобразует их
- отправляет на локальный HTTP сервер для теста

---

# АРХИТЕКТУРА

```text
api LeagueClient
        ↓
analytics/ GameTracker
        ↓
models/ PlayerAnalytics, Events
        ↓
transport/ HTTP sender
        ↓
server/ mock backend :8080
```

---

# ЧТО УЖЕ РЕАЛИЗОВАНО

## 1. Сбор данных

Файл: `internal/api`

- подключение к:

  ```
  https://127.0.0.1:2999
  ```

- Получение метрик игры

- Проверяет доступность игры

---

## 2. Предобработка

Файл: `internal/analytics`

Функция:

```go
GetPlayerAnalytics()
```

Что делает:

- находит текущего игрока
- считает:
  - kills / deaths / assists
  - CS (фарм)
  - gold
  - ward score
- собирает:
  - предметы (items)
  - способности (Q/W/E/R)
  - события игры

**Превращает raw data → structured analytics**

---

## 3. Модель данных

Файл: `internal/models`

Есть 3 ключевых сущности:

### PlayerAnalytics

- champion
- kills/deaths/assists
- cs
- gold
- items
- abilities
- ward score

### EventAnalytics

- тип события (kill, death и т.д.)
![alt text](image-1.png)
- время
- killer/victim

### ServerPayload

- timestamp
- player
- events

✔ это основной контракт между клиентом и сервером

---

## 4. Transport (отправка данных)

Файл: `internal/transport`

- HTTP POST отправка на сервер
- JSON payload
- configurable URL

✔ сейчас отправляет в:

```
http://localhost:8080/v1/update
```

---

## 5. Display (debug UI)

Файл: `internal/display`

- выводит live статус в консоль:
  - champion
  - KDA
  - CS
  - gold

✔ это только локальный debug слой

---

## 6. Mock Server

Файл: `cmd/server`

- принимает POST `/v1/update`
- хранит последний payload
- GET `/v1/update` возвращает данные
- логирует JSON красиво

✔ используется как тестовый backend

---

## 7. Main loop (tracker)

Файл: `cmd/tracker`

- polling каждые N секунд
- проверка game availability
- сбор данных
- анализ
- вывод
- отправка на сервер

✔ это runtime engine системы

---

# ПРОБЛЕМА СЕЙЧАС

## Нет формального контракта

Сейчас:

```text
Go structs = контракт
```

Проблемы:

- нельзя гарантировать совместимость
- frontend/backend зависят от Go кода
- нет версионности API
- сложно масштабировать

---

# ЧТО ДОЛЖЕН ДАТЬ OPENAPI

OpenAPI добавляет:

## 1. Единый контракт системы

```text
Client ↔ OpenAPI spec ↔ Server
```

---

## 2. Формальное описание API

- `/v1/update` (POST)
- `/v1/update` (GET debug)
- `/health`

---

## 3. Описание данных

- PlayerAnalytics
- EventAnalytics
- ServerPayload

---

## 4. Возможности будущего

- генерация Go моделей
- генерация JS/TS клиента
- validation запросов
- versioning API (/v1 → v2)

---

# ИТОГ (очень коротко)

## Сейчас у тебя:

✔ рабочий LoL telemetry client
✔ сбор live данных
✔ аналитика игрока
✔ события игры
✔ отправка на локальный сервер
✔ debug mock backend

---

## Чего не хватает:

- формального API контракта
- стандарта обмена данными
- генерации типов
- масштабируемости под backend/ML/frontend

---

## Что решает OpenAPI:

✔ фиксирует структуру данных
✔ делает систему расширяемой
✔ убирает зависимость от Go structs
✔ позволяет подключать любые языки

![alt text](image.png)