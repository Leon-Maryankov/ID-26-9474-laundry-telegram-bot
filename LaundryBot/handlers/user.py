from aiogram import Router
from aiogram.types import Message
from aiogram.filters import Command

router = Router()

HELP_TEXT = """
📌 Бот для записи в прачечную

Команды:
/start — запуск  
/help — помощь  
/book — записаться  
/cancel — отменить запись  
/my — мои записи  
"""

@router.message(Command("start"))
async def start(message: Message):
    await message.answer("👋 Добро пожаловать в бот прачечной!\nНапишите /help")

@router.message(Command("help"))
async def help_cmd(message: Message):
    await message.answer(HELP_TEXT)

@router.message(Command("book"))
async def book(message: Message):
    await message.answer("Выберите время (пока заглушка)")

@router.message(Command("cancel"))
async def cancel(message: Message):
    await message.answer("Запись отменена (пока заглушка)")