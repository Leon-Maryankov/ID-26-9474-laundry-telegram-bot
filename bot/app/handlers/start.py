# bot/app/handlers/start.py
from aiogram import Router, types
from aiogram.filters import CommandStart
from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, WebAppInfo

router = Router()

@router.message(CommandStart())
async def cmd_start(message: types.Message):
    kb = InlineKeyboardMarkup(inline_keyboard=[[
        InlineKeyboardButton(
            text="🧺 Забронировать слот",
            web_app=WebAppInfo(url="https://example.com")  # пока можно заглушку
        )
    ]])
    await message.answer("Нажми кнопку ниже, чтобы выбрать время стирки 👇", reply_markup=kb)