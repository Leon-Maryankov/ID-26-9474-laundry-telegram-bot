import aiohttp
from app.config import settings

class APIClient:
    def __init__(self):
        self.base_url = settings.GO_API_URL
        self.api_token = "token" 

    async def check_user_exists(self, telegramID: int):
        async with aiohttp.ClientSession() as session:
            url = f"{self.base_url}/{self.api_token}/doesUserExist"
            async with session.post(url, json={"telegramId": telegramID}) as resp:
                return await resp.json()

    async def register_user(self, telegramID, name, surname, email):
        async with aiohttp.ClientSession() as session:
            data = {
                "telegramId": telegramID,
                "name": name,
                "surname": surname,
                "email": email
            }
            url = f"{self.base_url}/{self.api_token}/registerUser"
            async with session.post(url, json=data) as resp:
                return await resp.json()


    # Другие методы: get_slots, book, cancel, get_bookings

    async def get_bookings(self, user_id: int):
        """Get bookings (slots) attached to a user by internal user id.

        Calls POST {base}/{token}/getSlotsAttachedToUser with JSON {"userId": user_id}.
        Returns parsed JSON response from the server.
        """
        async with aiohttp.ClientSession() as session:
            url = f"{self.base_url}/{self.api_token}/getSlotsAttachedToUser"
            async with session.post(url, json={"userId": user_id}) as resp:
                return await resp.json()

    async def get_available_slots(self):
        """Get available free slots from the server."""
        async with aiohttp.ClientSession() as session:
            url = f"{self.base_url}/{self.api_token}/getAvailableSlots"
            async with session.post(url, json={}) as resp:
                return await resp.json()

    async def occupy_slot(self, slot_id: int, user_id: int):
        """Occupy a slot for a user.

        Calls POST {base}/{token}/occupySlot with JSON {"slotId": slot_id, "userId": user_id}.
        Returns parsed JSON response from the server.
        """
        async with aiohttp.ClientSession() as session:
            url = f"{self.base_url}/{self.api_token}/occupySlot"
            async with session.post(url, json={"slotId": slot_id, "userId": user_id}) as resp:
                return await resp.json()

    async def cancel_slot(self, slot_id: int, user_id: int):
        """Cancel a booked slot for a user.

        Calls POST {base}/{token}/cancelSlot with JSON {"slotId": slot_id, "userId": user_id}.
        Returns parsed JSON response from the server.
        """
        async with aiohttp.ClientSession() as session:
            url = f"{self.base_url}/{self.api_token}/cancelSlot"
            async with session.post(url, json={"slotId": slot_id, "userId": user_id}) as resp:
                return await resp.json()

api_client = APIClient()