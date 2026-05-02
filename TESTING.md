# Интеграционное тестирование Go-сервера

Данные тесты проверяют работу API сервера через реальные HTTP-запросы (регистрация, бронирование, отмена, лимиты).

## Предварительные требования

- Установленные **Docker** и **Docker Compose**
- **Python 3.10** или новее (с `pip`)
- Клонированный репозиторий проекта

## Запуск тестов

### 1. Запустить сервер в контейнере
Из корневой папки проекта выполните:
```bash
docker-compose up --build
### 2. Установить SQLite внутри контейнера (делается один раз)
```bash
docker exec -it laundry-server apk add sqlite
### 3. Перейти в папку с интеграционными тестами
```bash
cd tests/integration
### 4. Установить Python-зависимости
```bash
pip install -r requirements-test.txt
### 5. Запустить тесты
```bash
pytest -v

###Ожидаемый результат
Все 5 тестов должны завершиться успешно:
```bash
test_api.py::test_register_and_check_user PASSED
test_api.py::test_get_available_slots PASSED
test_api.py::test_occupy_slot PASSED
test_api.py::test_cancel_slot PASSED
test_api.py::test_washing_limit PASSED
