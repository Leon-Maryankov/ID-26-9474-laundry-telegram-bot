from aiogram import Router
from aiogram.types import Message
from aiogram.filters import Command
from aiogram import Bot

router = Router()


@router.message(Command("help"))
async def cmd_help(message: Message, bot: Bot):
    bot_commands = await bot.get_my_commands()
    
    if bot_commands:
        commands_list = []
        for cmd in bot_commands:
            commands_list.append(f"/{cmd.command} - {cmd.description}")
        
        help_text = f"""
📖 **Доступные команды:**

{chr(10).join(commands_list)}
        """
    else:
        help_text = """
⚠️ Telegram API Error while parsing commands
        """
    
    await message.answer(help_text, parse_mode="Markdown")