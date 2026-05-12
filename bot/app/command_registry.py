commands_registry = {}

from aiogram.filters import Command

def register_command(router, command_name: str, description: str = ""):
    """Decorator registers both command in register and message handler"""
    def decorator(func):
        router.message(Command(command_name))(func)
        
        commands_registry[command_name] = {
            "description": description,
            "handler": func
        }
        
        return func
    return decorator