import pytest
import requests
import subprocess

BASE_URL = "http://localhost:8080/token"

def docker_exec_sql(sql: str) -> None:
    cmd = f"docker exec laundry-server sqlite3 /root/auxiliary/database.db \"{sql}\""
    subprocess.run(cmd, shell=True, check=True)

@pytest.fixture
def api_base():
    return BASE_URL

@pytest.fixture
def create_test_user(api_base):
    payload = {"telegramId": 999111, "name": "Api", "surname": "Test", "email": "api@test.com"}
    requests.post(f"{api_base}/registerUser", json=payload)
    resp = requests.post(f"{api_base}/doesUserExist", json={"telegramId": 999111})
    if resp.json()["answer"]:
        return resp.json()["user"]["id"]
    return None

@pytest.fixture
def create_free_slot():
    docker_exec_sql("INSERT INTO slots (date, machineId, machineType, userId) VALUES (strftime('%s', 'now', '+1 hour'), 0, 'washing', NULL);")
    result = subprocess.run("docker exec laundry-server sqlite3 /root/auxiliary/database.db \"SELECT id FROM slots WHERE userId IS NULL LIMIT 1;\"", shell=True, capture_output=True, text=True)
    slot_id = int(result.stdout.strip())
    yield slot_id
    docker_exec_sql(f"DELETE FROM slots WHERE id = {slot_id};")