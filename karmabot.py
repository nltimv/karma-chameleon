import os
from slack_bolt import App
from slack_bolt.adapter.socket_mode import SocketModeHandler
import psycopg2
import re

app = App(token=os.environ.get("SLACK_BOT_TOKEN"))

db_params = {
    "dbname": os.environ.get("DB_NAME"),
    "user": os.environ.get("DB_USER"),
    "password": os.environ.get("DB_PASSWORD"),
    "host": os.environ.get("DB_HOST"),
    "port": os.environ.get("DB_PORT"),
}

def update_karma(user_id, increment):
    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("SELECT karma FROM users WHERE user_id = %s", (user_id,))
    row = cursor.fetchone()

    if row is None:
        cursor.execute("INSERT INTO users (user_id, karma) VALUES (%s, %s)", (user_id, increment))
        conn.commit()
        karma = increment
    else:
        karma = row[0] + increment
        cursor.execute("UPDATE users SET karma = %s WHERE user_id = %s", (karma, user_id))
        conn.commit()

    cursor.close()
    conn.close()

    return karma

def get_karma(user_id):
    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("SELECT karma FROM users WHERE user_id = %s", (user_id,))
    row = cursor.fetchone()

    cursor.close()
    conn.close()

    return row[0] if row else 0

@app.message(re.compile(".*<@([a-zA-Z0-9_]+)>\\s?(\\+\\+\\+|\\-\\-\\-).*"))
def process_double_karma_message(say, context):
    user_id = context.matches[0]
    action = context.matches[1]
    current_user = context.user_id

    if user_id == current_user and action == "+++":
        say("Nice try! You can't boost your own ego. 😜")
        return

    increment = 2 if action == "+++" else -2
    karma = update_karma(user_id, increment)

    if increment == 2:
        say(f"<@{user_id}>'s karma got a double boost! 🚀 New karma count: {karma}")
    else:
        say(f"<@{user_id}>'s karma took a double hit! 💔 New karma count: {karma}")

@app.message(re.compile(".*<@([a-zA-Z0-9_]+)>\\s?(\\+\\+|\\-\\-).*"))
def process_karma_message(say, context):
    user_id = context.matches[0]
    action = context.matches[1]
    current_user = context.user_id

    if user_id == current_user and action == "++":
        say("Nice try! You can't boost your own ego. 😜")
        return

    increment = 1 if action == "++" else -1
    karma = update_karma(user_id, increment)

    if increment == 1:
        say(f"<@{user_id}>'s karma is on the rise! 🚀 New karma count: {karma}")
    else:
        say(f"<@{user_id}>'s karma took a hit! 💔 New karma count: {karma}")

@app.message(re.compile("<@([a-zA-Z0-9_]+)> karma"))
def get_user_karma(say, context):
    user_id = context.matches[0]
    karma = get_karma(user_id)
    say(f"<@{user_id}>'s current karma: {karma}")

@app.message(".*")
def ping(message, say):
    return

def create_table():
    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY, user_id VARCHAR(255) NOT NULL UNIQUE, karma INT NOT NULL )")
    conn.commit()

    cursor.close()
    conn.close()

if __name__ == "__main__":
    print("Initializing...")
    create_table()
    print("Starting...")
    SocketModeHandler(app, os.environ["SLACK_APP_TOKEN"]).start()
    