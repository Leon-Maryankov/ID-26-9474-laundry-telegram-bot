from dotenv import load_dotenv
import os

load_dotenv('common.env')   # общие настройки
load_dotenv()               # личные (переопределяют)

class Settings:
    BOT_TOKEN = os.getenv("BOT_TOKEN")
    PROXY = os.getenv("PROXY")
    GO_API_URL = os.getenv("GO_API_URL", "http://localhost:8080")
    WEBAPP_URL = os.getenv("WEBAPP_URL")

settings = Settings()