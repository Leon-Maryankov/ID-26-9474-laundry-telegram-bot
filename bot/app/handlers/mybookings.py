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

    user = resp.get("user") or {}
    user_id = user.get("id")
    if not user_id:
        await message.answer("⚠️ Не удалось получить идентификатор пользователя.")
        return

    try:
        bookings_resp = await api_client.get_bookings(user_id)
    except Exception:
        await message.answer("⚠️ Не удалось получить брони с сервера. Попробуй позже.")
        return

    if bookings_resp.get("answer") is False:
        await message.answer("⚠️ Сервер вернул ошибку при получении броней.")
        return

    slots = bookings_resp.get("slots", []) or []
    if len(slots) == 0:
        await message.answer("У вас нет активных броней.")
        return

    lines = []
    for s in slots:
        # handle different JSON key casings
        date = s.get("Date") or s.get("date") or ""
        if "T" in date:
            date = date.replace("T", " ").split("+")[0].split("Z")[0]
        mid = s.get("MachineId") or s.get("machineId") or "?"
        mtype = s.get("MachineType") or s.get("machineType") or "?"
        slot_id = s.get("Id") or s.get("id") or "?"
        lines.append(f"- {date} | Машина {mid} ({mtype}) | id:{slot_id}")

    await message.answer("Ваши брони:\n" + "\n".join(lines))