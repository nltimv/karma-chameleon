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

    try:
        response = requests.get(api_url, headers=headers)
        response.raise_for_status()
        user_info = response.json()
        return user_info.get("ok", False)
    except requests.exceptions.RequestException as e:
        print(f"Error validating user ID: {e}")
        return False

def get_usergroup_members(usergroup_id, token):
    api_url = f"https://slack.com/api/usergroups.users.list?usergroup={usergroup_id}"

    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/x-www-form-urlencoded"
    }

    try:
        response = requests.get(api_url, headers=headers)
        response.raise_for_status()
        usergroup_members = response.json().get("users", [])
        return usergroup_members
    except requests.exceptions.RequestException as e:
        print(f"Error getting user group members: {e}")
        return []

def update_user_karma(user_id, team_id, increment):
    if not is_valid_user(user_id, os.environ.get("SLACK_BOT_TOKEN")):
        return 0

    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("SELECT karma FROM users WHERE user_id = %s AND team_id = %s", (user_id, team_id))
    row = cursor.fetchone()

    if row is None:
        karma = increment
        cursor.execute("INSERT INTO users (user_id, team_id, karma) VALUES (%s, %s, %s)", (user_id, team_id, increment))
    else:
        karma = row[0] + increment
        cursor.execute("UPDATE users SET karma = %s WHERE user_id = %s AND team_id = %s", (karma, user_id, team_id))

    conn.commit()
    cursor.close()
    conn.close()

    return karma

def update_group_karma(group_id, team_id, increment, giver_user_id):

    # Give karma to each member of the user group, except the giver
    usergroup_members = get_usergroup_members(group_id, os.environ.get("SLACK_BOT_TOKEN"))
    for member_id in usergroup_members:
        if member_id != giver_user_id or increment < 0:
            _ = update_user_karma(member_id, team_id, increment)

    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("SELECT karma FROM groups WHERE group_id = %s AND team_id = %s", (group_id, team_id))
    row = cursor.fetchone()

    if row is None:
        karma = increment
        cursor.execute("INSERT INTO groups (group_id, team_id, karma) VALUES (%s, %s, %s)", (group_id, team_id, increment))
    else:
        karma = row[0] + increment
        cursor.execute("UPDATE groups SET karma = %s WHERE group_id = %s AND team_id = %s", (karma, group_id, team_id))

    conn.commit()
    cursor.close()
    conn.close()

    return karma

def get_user_karma(user_id, team_id):
    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("SELECT karma FROM users WHERE user_id = %s AND team_id = %s", (user_id, team_id))
    row = cursor.fetchone()

    cursor.close()
    conn.close()

    return row[0] if row else 0

def get_group_karma(group_id, team_id):
    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("SELECT karma FROM groups WHERE group_id = %s AND team_id = %s", (group_id, team_id))
    row = cursor.fetchone()

    cursor.close()
    conn.close()

    return row[0] if row else 0

@app.message(re.compile(".*<@([a-zA-Z0-9_]+)>\\s?(\\+\\+\\+|\\-\\-\\-|\\+\\+|\\-\\-).*"))
def process_karma_user_message(say, context):
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
    
    karma = update_user_karma(user_id, team_id, increment_value)

    if increment_value == 2:
        say(f"<@{user_id}>'s karma got a double boost! 🚀 New karma count: {karma}")
    elif increment_value == 1:
        say(f"<@{user_id}>'s karma is on the rise! 🚀 New karma count: {karma}")
    elif increment_value == -1:
        say(f"<@{user_id}>'s karma took a hit! 💔 New karma count: {karma}")
    elif increment_value == -2:
        say(f"<@{user_id}>'s karma took a double hit! 💔 New karma count: {karma}")

@app.message(re.compile("<@([a-zA-Z0-9_]+)>\\s?karma"))
def process_get_karma_user_message(say, context):
    user_id = context.matches[0]
    team_id = context.team_id
    bot_token = os.environ.get("SLACK_BOT_TOKEN")

    if not is_valid_user(user_id, bot_token):
        return  # Do nothing if the user ID is invalid

    karma = get_user_karma(user_id, team_id)
    say(f"<@{user_id}>'s current karma: {karma}")

@app.message(re.compile(r"<!subteam\^([a-zA-Z0-9_]+)\|?[a-zA-Z0-9_\-\.]*>\s?(\+\+\+|\-\-\-|\+\+|\-\-).*"))
def process_karma_group_message(say, context):
    group_id = context.matches[0]
    team_id = context.team_id
    bot_token = os.environ.get("SLACK_BOT_TOKEN")
    print(group_id)

    increment = context.matches[1]

    increment_value = \
              2 if increment == "+++"  \
        else  1 if increment == "++" \
        else -1 if increment == "--" \
        else -2 if increment == "---" \
        else  0

    usergroup_members = get_usergroup_members(group_id, bot_token)

    # Do nothing if the user group doesn't exist
    if len(usergroup_members) == 0:
        return

    # If the user giving karma is the only member of the group, send a sassy message
    if len(usergroup_members) == 1 and usergroup_members[0] == context.user_id:
        say(f"Nice try! Creating a user group for youself so you can get group karma? You're smart, but not smart enough! 😜")
        return

    # Give karma to each member of the user group, except the giver
    for member_id in usergroup_members:
        if member_id != context.user_id or increment_value < 0:
            _ = update_user_karma(member_id, team_id, increment_value)

    karma = update_group_karma(group_id, team_id, increment_value, context.user_id)

    if increment_value == 2:
        say(f"The karma of <!subteam^{group_id}> and its members got a double boost! 🚀 New group karma count: {karma}")
    elif increment_value == 1:
        say(f"The karma of <!subteam^{group_id}> and its members is on the rise! 🚀 New group karma count: {karma}")
    elif increment_value == -1:
        say(f"The karma of <!subteam^{group_id}> and its members took a hit! 💔 New group karma count: {karma}")
    elif increment_value == -2:
        say(f"The karma of <!subteam^{group_id}> and its members took a double hit! 💔 New group karma count: {karma}")

@app.message(re.compile(r"<!subteam\^([a-zA-Z0-9_]+)\|?[a-zA-Z0-9_\-\.]*>\s?karma"))
def process_get_karma_group_message(say, context):
    group_id = context.matches[0]
    team_id = context.team_id
    bot_token = os.environ.get("SLACK_BOT_TOKEN")
    print(group_id)

    usergroup_members = get_usergroup_members(group_id, bot_token)

    # Do nothing if the user group doesn't exist
    if len(usergroup_members) == 0:
        return

    karma = get_group_karma(group_id, team_id)
    say(f"Current karma for group <!subteam^{group_id}>: {karma}")

@app.message(".*")
def default():
    return

def create_tables():
    conn = psycopg2.connect(**db_params)
    cursor = conn.cursor()

    cursor.execute("CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY, user_id VARCHAR(255) NOT NULL, team_id VARCHAR(255) NOT NULL, karma INT NOT NULL )")
    cursor.execute("CREATE TABLE IF NOT EXISTS groups ( id SERIAL PRIMARY KEY, group_id VARCHAR(255) NOT NULL, team_id VARCHAR(255) NOT NULL, karma INT NOT NULL )")

    conn.commit()
    cursor.close()
    conn.close()

if __name__ == "__main__":
    print("Initializing...")
    create_tables()
    print("Starting...")
    SocketModeHandler(app, os.environ["SLACK_APP_TOKEN"]).start()