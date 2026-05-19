from aiogram import Router
from aiogram.types import Message
from aiogram.fsm.context import FSMContext
from app.services.api_client import api_client
from app.command_registry import register_command
from app.handlers.start import show_webapp_button
from collections import defaultdict
import datetime

router = Router()

@register_command(router, "book", "🧺 Забронировать слот")
async def cmd_book(message: Message, state: FSMContext):
    current_state = await state.get_state()
    if current_state is not None:
        await message.answer(
            "⏳ Вы сейчас находитесь в процессе регистрации.\n"
            "Пожалуйста, сначала завершите её, затем используйте /book."
        )
        return

    tg_id = message.from_user.id
    try:
        resp = await api_client.check_user_exists(tg_id)
    except Exception:
        await message.answer("⚠️ Не удалось связаться с сервером. Попробуй позже.")
        return

    if resp.get("answer") is False:
        await message.answer("⚠️ Пользователь не найден! Пожалуйста, начните с /start и зарегистрируйтесь.")
        return

    # Try to parse slot id from message text: /book 123
    text = (message.text or "").strip()
    parts = text.split()
    slot_id = None
    if len(parts) >= 2:
        # find first token that parses as int
        for token in parts[1:]:
            try:
                slot_id = int(token)
                break
            except Exception:
                continue

    if slot_id is None:
        # No id provided: show webapp button as before and usage hint
        await show_webapp_button(message)
        await message.answer("💡 Для бронирования через команду: /book <номер_слота>\nПример: /book 123")
        return

    # Resolve internal user id
    user = resp.get("user") or {}
    user_id = user.get("id")
    if not user_id:
        await message.answer("⚠️ Не удалось определить внутренний идентификатор пользователя.")
        return

    # Attempt to occupy slot
    try:
        occ = await api_client.occupy_slot(slot_id, user_id)
    except Exception:
        await message.answer("⚠️ Ошибка при обращении к серверу. Попробуйте позже.")
        return

    if occ.get("status") == "good":
        await message.answer(f"✅ Слот {slot_id} успешно забронирован.")
        return

    # status bad
    reason = occ.get("reason") or "Неизвестная причина"
    await message.answer(f"⚠️ Не удалось забронировать слот {slot_id}: {reason}")


@register_command(router, "list_slots", "📋 Показать доступные слоты")
async def cmd_list_slots(message: Message, state: FSMContext):
    """Fetch available slots, sort and present them grouped by type and machine id."""
    try:
        resp = await api_client.get_available_slots()
    except Exception as e:
        print(f"Error: {e}")
        await message.answer("⚠️ Не удалось получить список слотов с сервера. Попробуйте позже.")
        return

    status = resp.get("status")
    if status != "good":
        await message.answer("⚠️ Сервер вернул ошибку при запросе слотов.")
        return

    slots_list = resp.get("Slots") or resp.get("slots") or []
    if not slots_list:
        await message.answer("Нет доступных слотов.")
        return

    # Normalize and prepare slots for sorting
    normalized = []
    for s in slots_list:
        try:
            sid = int(s.get("Id") or s.get("id") or 0)
        except Exception:
            continue

        mtype = (s.get("MachineType") or s.get("machineType") or "").lower()
        try:
            mid = int(s.get("MachineId") or s.get("machineId") or 0)
        except Exception:
            mid = 0

        date = s.get("Date") or s.get("date") or ""
        # format date into readable form
        date_str = ""
        if isinstance(date, (int, float)):
            date_str = datetime.datetime.fromtimestamp(int(date)).strftime('%Y-%m-%d %H:%M')
        elif isinstance(date, str):
            if "T" in date:
                date_str = date.replace("T", " ").split("+")[0].split("Z")[0]
            else:
                date_str = date
        else:
            date_str = str(date)

        normalized.append({
            "Id": sid,
            "MachineType": mtype,
            "MachineId": mid,
            "Date": date_str,
        })

    if not normalized:
        await message.answer("Нет доступных слотов.")
        return

    # sort by type (washing, dryer), then machine id, then slot id
    type_order = {"washing": 0, "dryer": 1}
    normalized.sort(key=lambda x: (type_order.get(x["MachineType"], 2), x["MachineId"], x["Id"]))

    # group and build message
    lines = []
    groups = defaultdict(list)
    for s in normalized:
        groups[(s["MachineType"], s["MachineId"])].append(s)

    # Build message without asking for input
    for (mtype, mid), entries in groups.items():
        # header
        human_type = "🧺 Стиральная" if mtype == "washing" else "🔥 Сушильная"
        lines.append(f"\n{human_type} машина #{mid + 1}:")
        lines.append("━" * 30)
        
        # sort entries by Id
        entries.sort(key=lambda e: e["Id"])
        for e in entries:
            lines.append(f"  🆔 {e['Id']}  |  📅 {e['Date']}")
    
    lines.append("\n" + "━" * 30)
    # ИСПРАВЛЕНО: убраны угловые скобки
    lines.append("💡 Для бронирования используйте команду:\n/book номер_слота")
    lines.append("Номер лота - это номер после значка 🆔")
    lines.append("Пример: /book 123")

    # ИСПРАВЛЕНО: явно указываем parse_mode=None
    await message.answer("\n".join(lines), parse_mode=None)


@register_command(router, "cancel", "❌ Отменить бронь")
async def cmd_cancel(message: Message, state: FSMContext):
    """Cancel a booked slot: /cancel <slot_id>."""
    tg_id = message.from_user.id
    try:
        resp = await api_client.check_user_exists(tg_id)
    except Exception:
        await message.answer("⚠️ Не удалось связаться с сервером. Попробуй позже.")
        return

    if resp.get("answer") is False:
        await message.answer("⚠️ Пользователь не найден! Пожалуйста, начните с /start и зарегистрируйтесь.")
        return

    # Parse slot id from message text: /cancel 123
    text = (message.text or "").strip()
    parts = text.split()
    slot_id = None
    if len(parts) >= 2:
        for token in parts[1:]:
            try:
                slot_id = int(token)
                break
            except Exception:
                continue

    if slot_id is None:
        await message.answer("💡 Для отмены бронирования используйте команду:\n/cancel <номер\\_слота>\nПример: /cancel 123", parse_mode=None)

    # Resolve internal user id
    user = resp.get("user") or {}
    user_id = user.get("id")
    if not user_id:
        await message.answer("⚠️ Не удалось определить внутренний идентификатор пользователя.")
        return

    # Attempt to cancel slot
    try:
        cancel = await api_client.cancel_slot(slot_id, user_id)
    except Exception:
        await message.answer("⚠️ Ошибка при обращении к серверу. Попробуйте позже.")
        return

    if cancel.get("status") == "good":
        await message.answer(f"✅ Бронирование слота {slot_id} успешно отменено.")
        return

    # status bad
    reason = cancel.get("reason") or "Неизвестная причина"
    await message.answer(f"⚠️ Не удалось отменить бронирование слота {slot_id}: {reason}")