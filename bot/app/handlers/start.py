from aiogram import Router, F
from aiogram.types import Message
from aiogram.filters import CommandStart

router = Router()

@router.message(CommandStart())
async def cmd_start(message: Message):
    # Отправить приветствие, запросить email (или использовать Telegram Login)
    await message.answer("Привет! Для использования бота нужно зарегистрироваться. Отправь свой email.")