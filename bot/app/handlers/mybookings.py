from aiogram import Router
from aiogram.types import Message
from aiogram.filters import Command

router = Router()

@router.message(Command("mybookings"))
async def cmd_mybookings(message: Message):
    # Запросить у Go брони пользователя и показать
    await message.answer("Ваши брони: (заглушка)")