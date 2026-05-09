from aiogram import Router
from aiogram.types import Message
from aiogram.filters import Command

from app.command_registry import register_command

router = Router()

@register_command(router, "mybookings", "🧺 Посмотреть брони")
async def cmd_mybookings(message: Message):
    # Запросить у Go брони пользователя и показать
    await message.answer("Ваши брони: (заглушка)")