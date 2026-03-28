from aiogram import Router
from aiogram.types import Message, KeyboardButton, ReplyKeyboardMarkup, WebAppInfo
from aiogram.filters import Command

router = Router()

@router.message(Command("webapp"))
async def cmd_webapp(message: Message):
    web_app = WebAppInfo(url="http://localhost:8080") 
    button = KeyboardButton(text="Открыть мини-приложение", web_app=web_app)
    keyboard = ReplyKeyboardMarkup(keyboard=[[button]], resize_keyboard=True)
    await message.answer("Удобное бронирование через веб-приложение:", reply_markup=keyboard)