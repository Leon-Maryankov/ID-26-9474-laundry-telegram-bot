from aiogram import Router, F, types
from aiogram.fsm.state import State, StatesGroup
from aiogram.fsm.context import FSMContext
from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, WebAppInfo
from aiogram.filters import CommandStart
from app.command_registry import register_command
from app.config import settings
from app.services.api_client import api_client

router = Router()

class RegistrationStates(StatesGroup):
    waiting_for_email = State()

# 🖼 Хелпер: показать кнопку WebApp
async def show_webapp_button(message: types.Message):
    kb = InlineKeyboardMarkup(inline_keyboard=[[
        InlineKeyboardButton(
            text="🧺 Забронировать слот",
            web_app=WebAppInfo(url=settings.WEBAPP_URL)
        )
    ]])
    await message.answer("Нажми кнопку ниже, чтобы выбрать время стирки 👇", reply_markup=kb)

@register_command(router, "start", " Запустить бота")
async def cmd_start(message: types.Message, state: FSMContext):
    tg_id = message.from_user.id

    try:
        resp = await api_client.check_user_exists(tg_id)
        if resp.get("answer") is True:
            await show_webapp_button(message)
            await state.clear()
            return
    except Exception as e:
        await message.answer(f"⚠️ Не удалось связаться с сервером. Попробуй позже.")
        return

    await message.answer(
        "👋 Привет! Я пока не вижу тебя в базе.\n"
        "Отправь свою учебную почту (например, `ivanov@edu.spbstu.ru`):"
    )
    await state.set_state(RegistrationStates.waiting_for_email)

@router.message(RegistrationStates.waiting_for_email, F.text)
async def process_email(message: types.Message, state: FSMContext):
    email = message.text.strip()

    if "@" not in email or "." not in email.split("@")[-1]:
        await message.answer("❌ Некорректный формат почты. Попробуй ещё раз:")
        return

    if not email.endswith("@edu.spbstu.ru"):
        await message.answer("❌ Используй только почту @edu.spbstu.ru")
        return

    try:
        resp = await api_client.register_user(
            telegramID=message.from_user.id,
            name=message.from_user.first_name or "",
            surname=message.from_user.last_name or "",
            email=email
        )

        if resp.get("status") == "good":
            await message.answer("✅ Регистрация успешна! Добро пожаловать.")
            await show_webapp_button(message)
            await state.clear()  # ✅ КРИТИЧНО: выходим из состояния ожидания
        else:
            await message.answer(" Ошибка на сервере. Попробуй /start заново.")
            await state.clear()  # Сбрасываем, чтобы не зависнуть в ожидании
            
    except Exception as e:
        await message.answer(f"⚠️ Ошибка при регистрации: {e}")
        await state.clear()