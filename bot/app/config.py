import os
from dotenv import load_dotenv

load_dotenv()

class Settings:
    BOT_TOKEN = os.getenv("8687856104:AAHaweuoCAMSBrrrgDBYAGy6bEwiSK07H-M")
    GO_API_URL = os.getenv("GO_API_URL", "http://localhost:8080")

settings = Settings()