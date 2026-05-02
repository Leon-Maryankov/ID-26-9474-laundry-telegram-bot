import requests
import pytest

def test_register_and_check_user(api_base, create_test_user):
    user_id = create_test_user
    assert user_id is not None
    resp = requests.post(f"{api_base}/doesUserExist", json={"telegramId": 999111})
    data = resp.json()
    assert data["status"] == "good"
    assert data["answer"] is True

def test_get_available_slots(api_base):
    resp = requests.post(f"{api_base}/getAvailableSlots", json={})
    data = resp.json()
    assert data["status"] == "good"
    assert "Slots" in data  

def test_occupy_slot(api_base, create_test_user, create_free_slot):
    user_id = create_test_user
    slot_id = create_free_slot
    resp = requests.post(f"{api_base}/occupySlot", json={"slotId": slot_id, "userId": user_id})
    assert resp.json()["status"] == "good"
    resp2 = requests.post(f"{api_base}/getSlotsAttachedToUser", json={"userId": user_id})
    slots = resp2.json()["slots"]
    assert any(s["Id"] == slot_id for s in slots)  

def test_cancel_slot(api_base, create_test_user, create_free_slot):
    user_id = create_test_user
    slot_id = create_free_slot
    requests.post(f"{api_base}/occupySlot", json={"slotId": slot_id, "userId": user_id})
    cancel = requests.post(f"{api_base}/cancelSlot", json={"slotId": slot_id, "userId": user_id})
    assert cancel.json()["status"] == "good"
    free = requests.post(f"{api_base}/getAvailableSlots", json={}).json()["Slots"]
    assert any(s["Id"] == slot_id for s in free)  

def test_washing_limit(api_base, create_test_user):
    import subprocess
    import time
    
    def docker_exec_sql(sql: str) -> None:
        cmd = f"docker exec laundry-server sqlite3 /root/auxiliary/database.db \"{sql}\""
        subprocess.run(cmd, shell=True, check=True)
    
    user_id = create_test_user
    now = int(time.time())
    slot_ids = []
    for i in range(3):
        t = now + (i+1)*3600
        docker_exec_sql(f"INSERT INTO slots (date, machineId, machineType, userId) VALUES ({t}, {i}, 'washing', NULL);")
        result = subprocess.run(f"docker exec laundry-server sqlite3 /root/auxiliary/database.db \"SELECT id FROM slots WHERE date = {t};\"", shell=True, capture_output=True, text=True)
        slot_ids.append(int(result.stdout.strip()))
    try:
        for sid in slot_ids[:2]:
            resp = requests.post(f"{api_base}/occupySlot", json={"slotId": sid, "userId": user_id})
            assert resp.json()["status"] == "good"
        resp3 = requests.post(f"{api_base}/occupySlot", json={"slotId": slot_ids[2], "userId": user_id})
        assert resp3.json()["status"] == "bad"
        assert resp3.json()["reason"] == "too many washing machines"
    finally:
        for sid in slot_ids:
            docker_exec_sql(f"DELETE FROM slots WHERE id = {sid};")