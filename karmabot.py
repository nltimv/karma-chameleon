import os
import requests
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

def is_valid_user(user_id, token):
    api_url = f"https://slack.com/api/users.info?user={user_id}"

    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/x-www-form-urlencoded"
    }

    response = requests.get(api_url, headers=headers)
    user_info = response.json()

    return user_info.get("ok", False)

def update_karma(user_id, team_id, increment):
    if not is_valid_user(user_id, os.environ.get("SLACK_BOT_TOKEN")):
        return 0

    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("SELECT karma FROM users WHERE user_id = %s AND team_id = %s", (user_id, team_id))
    row = cursor.fetchone()

    if row is None:
        cursor.execute("INSERT INTO users (user_id, team_id, karma) VALUES (%s, %s, %s)", (user_id, team_id, increment))
    else:
        karma = row[0] + increment
        cursor.execute("UPDATE users SET karma = %s WHERE user_id = %s AND team_id = %s", (karma, user_id, team_id))

    conn.commit()
    cursor.close()
    conn.close()

    return karma

def get_karma(user_id, team_id):
    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("SELECT karma FROM users WHERE user_id = %s AND team_id = %s", (user_id, team_id))
    row = cursor.fetchone()

    cursor.close()
    conn.close()

    return row[0] if row else 0

@app.message(re.compile(".*<@([a-zA-Z0-9_]+)>\\s?(\\+\\+\\+|\\-\\-\\-|\\+\\+|\\-\\-).*"))
def process_karma_message(say, context):
    user_id = context.matches[0]
    current_user = context.user_id
    team_id = context.team_id
    bot_token = os.environ.get("SLACK_BOT_TOKEN")

    increment = context.matches[1]

    if not is_valid_user(user_id, bot_token):
        return  # Do nothing if the user ID is invalid

    increment_value = \
              2 if increment == "+++"  \
        else  1 if increment == "++" \
        else -1 if increment == "--" \
        else -2 if increment == "---" \
        else  0
    
    if user_id == current_user and increment_value > 0:
        say(f"Nice try! You can't boost your own ego. 😜")
        return
    
    karma = update_karma(user_id, team_id, increment_value)

    if increment_value == 2:
        say(f"<@{user_id}>'s karma got a double boost! 🚀 New karma count: {karma}")
    elif increment_value == 1:
        say(f"<@{user_id}>'s karma is on the rise! 🚀 New karma count: {karma}")
    elif increment_value == -1:
        say(f"<@{user_id}>'s karma took a hit! 💔 New karma count: {karma}")
    elif increment_value == -2:
        say(f"<@{user_id}>'s karma took a double hit! 💔 New karma count: {karma}")

@app.message(re.compile("<@([a-zA-Z0-9_]+)>\\s?karma"))
def get_user_karma(say, context):
    user_id = context.matches[0]
    team_id = context.team_id
    bot_token = os.environ.get("SLACK_BOT_TOKEN")

    if not is_valid_user(user_id, bot_token):
        return  # Do nothing if the user ID is invalid

    karma = get_karma(user_id, team_id)
    say(f"<@{user_id}>'s current karma: {karma}")

@app.message(".*")
def default():
    return

def create_table():
    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY, user_id VARCHAR(255) NOT NULL, team_id VARCHAR(255) NOT NULL, karma INT NOT NULL )")
    conn.commit()

    cursor.close()
    conn.close()

if __name__ == "__main__":
    print("Initializing...")
    create_table()
    print("Starting...")
    SocketModeHandler(app, os.environ["SLACK_APP_TOKEN"]).start()