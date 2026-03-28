import asyncio
import logging
from aiogram import Bot, Dispatcher
from aiogram.enums import ParseMode
from aiogram.fsm.storage.memory import MemoryStorage
from app.config import settings
from app.handlers import start, booking, mybookings, webapp
from app.scheduler.reminders import start_scheduler

logging.basicConfig(level=logging.INFO)

async def main():
    bot = Bot(token=settings.BOT_TOKEN, parse_mode=ParseMode.HTML)
    dp = Dispatcher(storage=MemoryStorage())

    # Подключаем роутеры
    dp.include_router(start.router)
    dp.include_router(booking.router)
    dp.include_router(mybookings.router)
    dp.include_router(webapp.router)

    # Запускаем планировщик напоминаний
    start_scheduler(bot)

    # Запускаем бота
    await bot.delete_webhook(drop_pending_updates=True)
    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())