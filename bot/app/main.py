import asyncio
import logging
from aiogram import Bot, Dispatcher
from aiogram.client.session.aiohttp import AiohttpSession
from aiogram.client.bot import DefaultBotProperties
from aiogram.enums import ParseMode
from app.config import settings
from app.handlers import start

# Настраиваем логирование, чтобы видеть ошибки
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

async def main():
    session = AiohttpSession(proxy=settings.PROXY)

    bot = Bot(
        token=settings.BOT_TOKEN,
        session=session,
        default=DefaultBotProperties(parse_mode=ParseMode.HTML)
    )
    
    dp = Dispatcher()

    dp.include_router(start.router)

    await bot.delete_webhook(drop_pending_updates=True)
    
    logger.info("🤖 Бот запущен и ждет команды...")
    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())