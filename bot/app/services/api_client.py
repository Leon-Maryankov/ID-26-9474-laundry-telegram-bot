import aiohttp
from app.config import settings

class APIClient:
    def __init__(self):
        self.base_url = settings.GO_API_URL
        self.api_token = "token" 

    async def register_user(self, telegramID, first_name, last_name, email):
        async with aiohttp.ClientSession() as session:
            data = {
                "telegramId": telegramID,
                "name": first_name,
                "surname": last_name,
                "email": email
            }
            url = f"{self.base_url}/{self.api_token}/registerUser"
            async with session.post(url, json=data) as resp:
                return await resp.json()

    # Другие методы: get_slots, book, cancel, get_bookings

api_client = APIClient()