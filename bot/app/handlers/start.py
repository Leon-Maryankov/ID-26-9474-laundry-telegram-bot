# bot/app/handlers/start.py
from aiogram import Router, types
from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, WebAppInfo
from app.command_registry import register_command

router = Router()

@register_command(router, "start", "🚀 Запустить бота")
async def cmd_start(message: types.Message):
    kb = InlineKeyboardMarkup(inline_keyboard=[[
        InlineKeyboardButton(
            text="🧺 Забронировать слот",
            web_app=WebAppInfo(url="https://polylaundry.netlify.app") 
        )
    ]])
    await message.answer("Нажми кнопку ниже, чтобы выбрать время стирки 👇", reply_markup=kb)