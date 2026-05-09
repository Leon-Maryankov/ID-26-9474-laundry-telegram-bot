from aiogram import Router
from aiogram.types import Message
from aiogram.filters import Command
from aiogram.fsm.context import FSMContext
from app.services.api_client import api_client

from app.command_registry import register_command

router = Router()

from aiogram.fsm.state import State

@register_command(router, "mybookings", "🧺 Посмотреть брони")
async def cmd_mybookings(message: Message, state: FSMContext):
    current_state = await state.get_state()

    if current_state is not None:
        await message.answer(
            "⏳ Вы сейчас находитесь в процессе регистрации.\n"
            "Пожалуйста, сначала введите почту, чтобы завершить её."
        )
        return
    
    tg_id = message.from_user.id
    # может произойти так, что бот почему то выключился в процессе регистрации
    # пользователь еще не ввел почту, и при этом находится в state none
    try:
        resp = await api_client.check_user_exists(tg_id)
        if resp.get("answer") is False:
            await message.answer(f"⚠️ Пользователь не найден! Похоже, вы не завершили регистрацию. Нажмите /start")
            return
    except Exception as e:
        await message.answer(f"⚠️ Не удалось связаться с сервером. Попробуй позже.")
        return

    await message.answer("Ваши брони: (заглушка)")