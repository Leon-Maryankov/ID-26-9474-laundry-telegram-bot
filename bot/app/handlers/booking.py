from aiogram import Router, F
from aiogram.types import Message, CallbackQuery
from aiogram.fsm.context import FSMContext
from aiogram.filters import Command

router = Router()

@router.message(Command("book"))
async def cmd_book(message: Message, state: FSMContext):
    # Начать процесс бронирования (выбор даты, слота)
    await message.answer("Выберите дату (например, 2025-04-01):")