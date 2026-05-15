# LaundryMachineBot

Telegram-бот для автоматизации записи в прачечную общежития
# Прачечная – Telegram-бот для бронирования стиральных машин

## Общая архитектура

- **Go-сервер** – REST API, работающий с SQLite. Управляет пользователями, слотами, бронированием и уведомлениями.
- **Telegram-бот** (Python + aiogram) – взаимодействует с пользователями, отправляет кнопку для открытия WebApp и рассылает напоминания.
- **WebApp** (HTML/CSS/JS) – размещён на Netlify и открывается внутри Telegram. Предоставляет интерфейс для просмотра слотов и бронирования.

## Локальный запуск (для разработки и тестирования)

### 1. Клонирование репозитория
```bash
git clone https://github.com/Leon-Maryankov/ID-26-9474-laundry-telegram-bot.git
cd ID-26-9474-laundry-telegram-bot
```
### 2. Запуск сервера (Docker)
Убедитесь, что установлены Docker и Docker Compose.
```bash
docker-compose up --build
```
Сервер будет доступен по адресу http://localhost:8080.

### 3. Настройка окружения для Telegram-бота
#### 3.1. Создайте виртуальное окружение и активируйте его
```bash
cd bot
python -m venv venv
# Windows:
venv\Scripts\activate
# macOS/Linux:
source venv/bin/activate
```
#### 3.2. Установите зависимости
```bash
pip install -r requirements.txt
```
Если возникают ошибки с SOCKS-зависимостями, выполните:
```bash
pip install pysocks
pip install 'aiohttp[socks]'
```
#### 3.3. Настройте переменные окружения
Скопируйте файл .env.example в .env и отредактируйте его:
```bash
cp .env.example .env   # Windows: copy .env.example .env
```
#### Пример содержимого .env:
```text
BOT_TOKEN=ваш_токен_от_BotFather
PROXY=socks5://127.0.0.1:10808   # если нужен прокси, иначе удалите строку
GO_API_URL=http://localhost:8080
WEBAPP_URL=https://laundrytgbot.netlify.app   # или ваш публичный адрес на Netlify
```
#### Примечание:
BOT_TOKEN – получите у @BotFather.
PROXY – укажите только если Telegram API заблокирован (например, через ваш VPN). Если не нужен – удалите строку или оставьте пустой.
WEBAPP_URL – должен использовать HTTPS. В разработке можно временно использовать локальный сервер, но Telegram требует публичный HTTPS. Рекомендуется задеплоить WebApp на Netlify (см. раздел 5).
#### 3.4. Запустите бота
```bash
python -m app.main
```

### 4. Запуск WebApp локально (опционально, для отладки)
Если вы не хотите каждый раз деплоить на Netlify, можете запустить локальный сервер:
```bash
cd netlifyWebApp
python -m http.server 8000
```
Тогда в .env укажите WEBAPP_URL=http://localhost:8000 (но учтите, что Telegram не примет HTTP-ссылку). Для тестирования через браузер или внутри Telegram лучше использовать публичный HTTPS (Netlify/ngrok).

### 5. Деплой WebApp на Netlify
- Замените содержимое папки netlifyWebApp вашими файлами (index.html, style.css, app.js).
- Запушьте изменения на GitHub.
- Авторизуйтесь на Netlify, выберите репозиторий и настройте публикацию:
- Base directory – оставьте пустым
- Build command – пусто
- Publish directory – netlifyWebApp

После деплоя вы получите адрес вида https://ваш-проект.netlify.app. Укажите его в WEBAPP_URL в .env бота.
