from aiogram import Router
from aiogram.types import Message
from aiogram.filters import Command

router = Router()

ADMIN_ID = 123456789  # сюда свой id

def is_admin(user_id: int):
    return user_id == ADMIN_ID

@router.message(Command("add_slot"))
async def add_slot(message: Message):
    if not is_admin(message.from_user.id):
        return await message.answer("Нет доступа")

    await message.answer("Слот добавлен (заглушка)")

@router.message(Command("all"))
async def all_bookings(message: Message):
    if not is_admin(message.from_user.id):
        return await message.answer("Нет доступа")

    await message.answer("Все записи (заглушка)")