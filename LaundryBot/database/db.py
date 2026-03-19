import sqlite3

conn = sqlite3.connect("laundry.db")
cursor = conn.cursor()

def init_db():
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS bookings (
        user_id INTEGER,
        time TEXT
    )
    """)
    conn.commit()

def add_booking(user_id, time):
    cursor.execute(
        "INSERT INTO bookings (user_id, time) VALUES (?, ?)",
        (user_id, time)
    )
    conn.commit()

def get_bookings(user_id):
    cursor.execute(
        "SELECT time FROM bookings WHERE user_id = ?",
        (user_id,)
    )
    return cursor.fetchall()

def delete_booking(user_id):
    cursor.execute(
        "DELETE FROM bookings WHERE user_id = ?",
        (user_id,)
    )
    conn.commit()