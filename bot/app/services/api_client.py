import aiohttp
from app.config import settings

class APIClient:
    def __init__(self):
        self.base_url = settings.GO_API_URL

    async def register_user(self, tg_id, first_name, last_name, email):
        async with aiohttp.ClientSession() as session:
            data = {
                "tg_id": tg_id,
                "first_name": first_name,
                "last_name": last_name,
                "email": email
            }
            async with session.post(f"{self.base_url}/api/register", json=data) as resp:
                return await resp.json()

    # Другие методы: get_slots, book, cancel, get_bookings

api_client = APIClient()