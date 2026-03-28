import asyncio
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from aiogram import Bot
from app.services.api_client import api_client

scheduler = AsyncIOScheduler()

async def check_and_notify(bot: Bot):
    # Запросить у Go список броней, требующих напоминания
    # Отправить сообщение каждому пользователю
    pass

def start_scheduler(bot: Bot):
    scheduler.add_job(check_and_notify, "interval", minutes=1, args=[bot])
    scheduler.start()