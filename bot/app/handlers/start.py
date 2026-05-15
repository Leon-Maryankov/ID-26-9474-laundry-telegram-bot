import logging
from aiogram import Router, F, types
from aiogram.fsm.state import State, StatesGroup
from aiogram.fsm.context import FSMContext
from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, WebAppInfo
from aiogram.filters import CommandStart
from app.config import settings
from app.services.api_client import api_client

router = Router()
logger = logging.getLogger(__name__)

class RegistrationStates(StatesGroup):
    waiting_for_email = State()

async def show_webapp_button(message: types.Message):
    """Отправляет кнопку для открытия WebApp"""
    logger.info(f"Отправляем кнопку с URL: {settings.WEBAPP_URL}")
    kb = InlineKeyboardMarkup(inline_keyboard=[[
        InlineKeyboardButton(
            text="🧺 Забронировать слот",
            web_app=WebAppInfo(url=settings.WEBAPP_URL)
        )
    ]])
    await message.answer("Нажми кнопку ниже, чтобы выбрать время стирки 👇", reply_markup=kb)
    logger.info("Кнопка отправлена")

@router.message(CommandStart())
async def cmd_start(message: types.Message, state: FSMContext):
    tg_id = message.from_user.id
    try:
        resp = await api_client.check_user_exists(tg_id)
        logger.info(f"Ответ сервера для tg_id {tg_id}: {resp}")
        if resp.get("answer") is True:
            await show_webapp_button(message)
            await state.clear()
            return
    except Exception as e:
        logger.error(f"Ошибка при check_user_exists: {type(e).__name__}: {e}", exc_info=True)
        await message.answer("⚠️ Не удалось связаться с сервером. Попробуй позже.")
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
            # ➡️ ВОТ ЭТО ВАЖНО: вызываем кнопку WebApp
            await show_webapp_button(message)
            await state.clear()
        else:
            await message.answer("❌ Ошибка на сервере. Попробуй /start заново.")
            await state.clear()
    except Exception as e:
        logger.error(f"Ошибка регистрации: {e}", exc_info=True)
        await message.answer(f"⚠️ Ошибка при регистрации: {e}")
        await state.clear()