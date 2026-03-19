from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton

def get_time_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="18:00", callback_data="time_18")],
        [InlineKeyboardButton(text="19:00", callback_data="time_19")],
        [InlineKeyboardButton(text="20:00", callback_data="time_20")],
    ])
    return keyboard