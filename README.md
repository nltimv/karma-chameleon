# Karma Chameleon

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/nltimv/karma-chameleon/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/nltimv/karma-chameleon/tree/main) [![license: MIT license](https://img.shields.io/badge/license-MIT%20license-blue.svg)](LICENSE) [![security: bandit](https://img.shields.io/badge/security-bandit-yellow.svg)](https://github.com/PyCQA/bandit)

![Logo](assets/img/logo.jpg){width=200}

> Spreading good vibes, one karma leap at a time! 🦎✨

Karma Chameleon is a bot for Slack teams that enables users to give and receive karma points. Users can boost the karma of individual team members or entire user groups, fostering positivity and recognition within the team. The bot stores karma information in a PostgreSQL database. Most of the code has been generated using ChatGPT. Give karma, spread good vibes, and let the Slack Karma Bot enhance your team's collaborative spirit!

## Installation
To use this bot, you'll need to add a Slack app to your workspace. The bot has not been submitted to the Slack App Directory, so you'll need to create your own app and host the bot yourself. Use the `slack-app-manifest.yml` in the root of the repo when creating the app to configure the app, including the necessary permissions.

### Deployment
#### Helm chart (recommended)
The Helm chart is the recommended way to install the bot on a production environment. The Helm chart requires an existing Kubernetes cluster with an existing PostgreSQL deployment. **A PostgreSQL instance is NOT included in the Helm chart!** Also note that PostgreSQL is the only supported database type at the moment.

The Helm chart can be deployed using the following command (filling in the placeholders with your own values):
```
 $ helm install karma-chameleon \
   oci://ghcr.io/nltimv/helm/karma-chameleon \
   --set karmaChameleon.db.databaseName=<database name> \
   --set karmaChameleon.db.host=<PostgreSQL server host/ip> \
   --set karmaChameleon.db.port=<PostgreSQL server port (if not 5432)> \
   --set karmaChameleon.db.username=<database username> \
   --set karmaChameleon.db.password=<database password> \
   --set karmaChameleon.slack.appToken<your Slack app token> \
   --set karmaChameleon.slack.botToken=<your Slack bot token> \
   -n karma-chameleon --create-namespace --wait
```

For more details about values, refer to the [values.yaml](charts/karma-chameleon/values.yaml) file.

#### Docker Compose 
> **Warning!**
  The Docker Compose method is primarily intended for development purposes. Therefore, it's not recommended to use this for production deployments!

You can use Docker Compose to quickly try out the bot or to use it for development purposes. The deployment includes an empty PostgreSQL server and database, which is NOT persistent. To deploy using Docker compose:
  1. Check out this repository.
      ``` 
      $ git clone https://github.com/nltimv/karma-chameleon.git
      ```
  2. Copy the `.env.overrides.template` file to `.env.overrides`
  3. Enter your Slack app and bot token in the `.env.overrides` file
  4. Run Docker Compose
     ```
     $ docker-compose up
     ```
  5. To stop the application, use
     ```
     $ docker-compose down
     ```

## Usage
To use the bot, you'll need to add it to a public or private channel in Slack. Then send one of the following messages to interact with the bot:

### User interactions
Replace `@user` with the handle of the user you want to do karma operations on.
| Message       | Description                                           |
| ------------- | ----------------------------------------------------- |
| `@user ++`    | Awards one karma to @user. Does not work on yourself. |
| `@user +++`   | Awards two karma to @user. Does not work on yourself. |
| `@user --`    | Takes away one karma from @user.                      |
| `@user ---`   | Takes away two karma from @user.                      |
| `@user karma` | Returns the current karma of @user.                   |

### User group interactions
Replace `@group` with the handle of the user group you want to do karma operations on. (User groups is a premium feature of Slack, so you will not be able to use this if you are on the Free plan)
| Message        | Description                                                                                                               |
| -------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `@group ++`    | Awards one karma to @group and all its members, excluding yourself. Does not work if you're the only member of the group. |
| `@group +++`   | Awards two karma to @group and all its members, excluding yourself. Does not work if you're the only member of the group. |
| `@group --`    | Takes away one karma from @group and all its members.                                                                     |
| `@group ---`   | Takes away two karma from @group and all its members.                                                                     |
| `@group karma` | Returns the current karma of @group.                                                                                      |

## Development

To start development of the bot, fork and check out the repo and install all dependencies. Make sure you have Python 3 installed on your computer.

```
$ git clone https://github.com/<username|org>/karma-chameleon.git
$ cd karma-chameleon
$ python -m pip install -r requirements.txt
```

### Pre-commit

This repository has pre-commit hooks configured to ensure code quality and security. Make sure the changes you submit pass all checks. To run the hooks on every commit, run this from the root of the repository:
```
$ pre-commit install
```
For more information about pre-commit, see [pre-commit.com](https://pre-commit.com/).

### Contributing

First of all, thank you for your interest in this project!

Before you start making changes, it is good practice to create an issue first describing the changes you want to make. There can be many reason I won't accept changes (for example, the change could go against the direction I want this project to go). By creating an issue, you can get an indication of whether your proposed changes/features are desirable, before you spend your time ans effort into developing something.

Before you submit a pull request, make sure you at least check the following:
 1. Check that pre-commit does not give any errors; (see above)
 2. Check that the Helm chart has been updated, if applicable;
 3. Make sure that both new installations and upgrades of existing installations work with your changes;
 4. Make sure to update this README, if applicable