import asyncio
from aiogram import Bot, Dispatcher
from handlers import user, admin
from utils.logger import setup_logger

TOKEN = "8687856104:AAHaweuoCAMSBrrrgDBYAGy6bEwiSK07H-M"

bot = Bot(token=TOKEN)
dp = Dispatcher()

async def main():
    setup_logger()

    # регистрация роутеров
    dp.include_router(user.router)
    dp.include_router(admin.router)

    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())