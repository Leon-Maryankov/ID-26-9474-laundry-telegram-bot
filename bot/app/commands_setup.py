from aiogram.types import BotCommand
from app.command_registry import commands_registry

async def setup_bot_commands(bot):
    """
    Collects all registered commands from the command registry and sets them in BotFather.
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