# Трамплин (Trumplin) — карьерная платформа
## Стек
- **Frontend:** `Next.js (React)`
- **Backend:** `Go`
- **Database:** `PostgreSQL`
- **Карта:** **Yandex Maps JS** + **Yandex Geocoder HTTP** (для офлайн-адресов) + кеш в PostgreSQL
- **Docker:** `docker-compose`

## Важно про роли (по ТЗ)
- **Соискатель (`APPLICANT`)**
  - отклики на возможности
  - приватность (скрыть отклики/резюме, открыть профиль для нетворкинга)
  - нетворкинг (контакты)
- **Работодатель (`EMPLOYER`)**
  - создание возможностей
  - просмотр откликов на свои возможности и изменение статусов откликов
- **Куратор / Админ (`CURATOR` / `ADMIN`)**
  - верификация компаний
  - модерация возможностей (approve/reject и статусы)

## Логотип
Логотип лежит в `frontend/public/logo.png` и используется в шапке фронтенда.

## Переменные окружения
Шаблон: `.env.example` в корне проекта.

Рекомендация: скопируй `.env.example` в `.env` и заполни значения.

### Что обязательно
- `TRUMPLIN_DATABASE_DSN` — DSN подключения к PostgreSQL
- `YANDEX_GEOCODER_KEY` — ключ для HTTP Geocoder (геокодирование офлайн адресов)
- `YANDEX_JAVASCRIPT_API_KEY` — ключ для карты (JS API)

### Что не обязательно (имеют dev-default)
- `TRUMPLIN_JWT_SECRET` — если не задан, backend стартует с dev-заглушкой
  - В проде/демо обязательно замени.

### Админ-куратор
Чтобы залогиниться как куратор (модерация), задайте:
- `TRUMPLIN_ADMIN_EMAIL`
- `TRUMPLIN_ADMIN_PASSWORD`

При старте backend создаст пользователя-админа, если его еще нет.

## Запуск

### Вариант A: запуск через Docker (рекомендуется)
1. Проверь, что Docker Desktop работает.
2. В корне проекта:
   - `docker compose up --build`
3. Открой:
   - Frontend: `http://localhost:3000`
   - Backend health: `http://localhost:8080/api/health`

Заметки по cookies/CORS:
- Backend настроен на CORS origin `http://localhost:3000` (переменная `TRUMPLIN_CORS_ORIGIN`).
- Аутентификация сделана через `httpOnly` cookie.

### Вариант B: запуск вручную (локальная разработка)

#### 1) PostgreSQL
Подними PostgreSQL любым способом. Docker тоже подходит, но можно и локально.

#### 2) Backend
В терминале (PowerShell) выставь как минимум:
- `TRUMPLIN_DATABASE_DSN`
- `YANDEX_GEOCODER_KEY`
- `TRUMPLIN_ADMIN_EMAIL`
- `TRUMPLIN_ADMIN_PASSWORD`
- `YANDEX_GEOCODER_KEY` (для геокодинга)

Далее:
```powershell
cd backend/cmd/server
go run main.go
```

Backend автоматически:
- применяет миграции
- создаёт админа (если заданы admin env vars)
- добавляет демо-данные для карты (пару `APPROVED` возможностей)

После старта проверь:
- `http://localhost:8080/api/health`

#### 3) Frontend
В отдельном терминале:
```powershell
cd frontend
npm run dev
```

В `frontend` в `.env.local` (или переменными окружения) установи:
- `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080`
- `NEXT_PUBLIC_YANDEX_API_KEY=...`

Открой:
- `http://localhost:3000`

## Как пользоваться (быстрый сценарий)

### 1) Публичная витрина (главная)
На главной:
- карта + лента возможностей
- маркер на карте подсвечивает карточку **на hover**
- по клику карточка становится “закрепленной”
- “избранное” хранится в `localStorage` (требование ТЗ)

### 2) Регистрация
Доступно:
- `http://localhost:3000/login`
- `http://localhost:3000/register`

Роль при регистрации:
- `APPLICANT` — соискатель
- `EMPLOYER` — работодатель

### 3) Работодатель → создание возможностей
На `GET/POST /api/employer/opportunities` работает RBAC.
Создание возможно только если:
- `companies.verification_status = APPROVED`

Поэтому в MVP два пути:
- сначала войти куратором и верифицировать компанию
- либо использовать демо-работодателя (он уже `APPROVED` и имеет демо-вакансии)

В интерфейсе (кабинет `EMPLOYER`) есть форма создания возможности (MVP: `CITY` + skills + зарплата опционально).

### 4) Куратор/Админ → модерация
В кабинете (роль `ADMIN/CURATOR`) доступны:
- список компаний `PENDING` + кнопки `Одобрить/Отклонить`
- список возможностей `PENDING` + кнопки `Одобрить/Отклонить`

После модерации обновляй список/страницу — карточки на главной появятся.

## Основные эндпоинты API (как устроено)
### Публично
- `GET /api/public/opportunities?city=...`

### Auth
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `POST /api/auth/refresh`
- `GET /api/me`

### Employer
- `GET /api/employer/opportunities`
- `POST /api/employer/opportunities`

### Applicant
- `GET /api/applicant/applications`
- `POST /api/applicant/applications`
- `PATCH /api/applicant/privacy`
- `GET /api/applicant/contacts`
- `POST /api/applicant/contacts`

### Curator/Admin
- `GET /api/curator/companies/pending`
- `PATCH /api/curator/companies/{companyId}/verification`
- `GET /api/curator/opportunities/pending`
- `PATCH /api/curator/opportunities/{opportunityId}/status`

## Маппинг данных на карту
- В БД координаты (`lat/lng`) сохраняются при создании возможности через геокодинг.
- На витрине показываются только возможности со статусом `APPROVED` и известными координатами.
- Для избранного используется другой цвет маркера (в `YandexMap`).

## Примечания безопасности
- Токены выдаются и хранятся через `httpOnly` cookie.
- RBAC проверяется в backend (никогда не “доверяем” фронтенду).
- Password хранится как `argon2id` hash.
- Добавлены refresh-token rotation и таблица `refresh_tokens`.

## Контакты/презентация
По требованиям задания вы должны приложить:
- презентацию (до 5 МБ)
- видеоролик (до 5 минут)

