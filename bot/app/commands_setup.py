from aiogram.types import BotCommand
from app.command_registry import commands_registry

BOT_DESCRIPTION = """🫧 Привет! Это бот прачечной общежития.

🔑 Регистрация:
1. Нажми /start
2. Введи свою почту @edu.spbstu.ru
3. Готово! Теперь можешь бронировать слоты.

📌 Как работает бронирование:
1. Ты выбираешь дату и время.
2. Система резервирует за тобой машинку.
3. Проверяй свои брони в /mybookings и отменяй их в /cancel.

🗓 Основные команды:
🧺 /book — Забронировать слот
👕 /mybookings — Мои записи
📋 /list_slots — Доступные слоты
❌ /cancel — Отменить бронь
❓ /help — Справка"""

async def setup_bot_commands(bot):
    """
    Collects all registered commands from the command registry and sets them in BotFather.
    Also sets the bot description.
    """
    commands = [
        BotCommand(command=cmd, description=desc["description"]) 
        for cmd, desc in commands_registry.items()
    ]
    
    if commands:
        await bot.set_my_commands(commands)
        print(f"Registered {len(commands)} commands in BotFather")
        for cmd in commands:
            print(f"   /{cmd.command} - {cmd.description}")
    else:
        print("⚠️ No commands to register in BotFather")
    
    # Устанавливаем описание бота
    try:
        await bot.set_my_description(description=BOT_DESCRIPTION)
        print("✅ Описание бота установлено успешно")
    except Exception as e:
        print(f"⚠️ Ошибка при установке описания бота: {e}")