# Finance Tracker Bot

## Description

Finance Tracker Bot is a Telegram bot that helps users track their financial transactions. It allows users to add, categorize, and analyze their expenses.

## Table of contents

1. [Installation & Setup](#installation--setup)
2. [Running the Rroject](#running-the-project)
3. [Deployment](#deployment)
4. [Logging & Debugging](#logging--debugging)
5. [Database](#database)
6. [Overview](#overview)

## Installation & Setup

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.msys2.org) (optional)

### Clone the Repository

```sh
git clone https://github.com/iv-sukhanov/finance_tracker.git
cd finance_tracker
```

### Environment Variables

Create an `.env/.dev` file and configure the required environment variables:

```
TELEGRAM_BOT_TOKEN=your_token
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_DB=finance_tracker
DATABASE_URL=postgres://your_user:your_password@db:5432/finance_tracker?sslmode=disable
```

Additionally there are optional variables:

```
LOG_LEVEL=(DEBUG|INFO|WARN|ERROR|FATAL|PANIC)
TELEGRAM_DEBUG_MODE=(true|false)
TELEGRAM_USERNAME=your_tg_username
APP_NAME=finance_tracker_bot
```

`LOG_LEVEL` defines the log level of the application *default is INFO

`TELEGRAM_DEBUG_MODE` defines if the telegram API debug logs would be shown *default is false

`TELEGRAM_USERNAME` would add your username to the internal error messages to the users

`APP_NAME` would add the app name to some logs 

## Running the Project

### Using Docker Compose

```sh
docker compose up -d --build
```

### Using make

```sh
make compose-up
```

or

```sh
make compose-restart
```

This will start the PostgreSQL database, run migrations, and start the bot.

## Deployment

This project uses GitHub Actions for automated deployment. 

To trigger a deployment manually you first need to specify `SSH_PRIVATE_KEY` and `VPS_HOST` variables in the github secrets.

Then the workflow dispatch could be triggered running:

```sh
gh workflow run main.yml
```

or it is possible to add an environmental variable `WORKFLOW_ID=workflow_id` and deploy by running:

```sh
make gh-action
```

## Logging & Debugging

- View real-time logs:
  ```sh
  docker compose logs -f
  ```
- Check database logs:
  ```sh
  docker compose logs -f db
  ```

## Database

The project uses postgres to store the data. The posrgres container is run automatically
with docker compose, then a migrate container applies the init migration. After that the go application connects to the db.

### Schema

Here is the database schema:

![database schema](/doc/schema.png)

There are tree tables: users, spending_categories, and spending_records.

The relation between the users table and the spending_catigories is one to many, and the relation between the spending_catigories and the spending_records table is one to many as well.

## Overview

### Motivation
As a poor student I have always wanted to keep track of the amount of money I spend on different things to then reconsider the spendings and manage the finances more wisely. But it was too hard for me to keep in mind all the spendings and then update the EXEL table on the laptop every day. Bank applications do not provide enough details about the spendings (*BoC especially*🤭🤭), and different applications felt like too much overhead so they did not really suit me.

So I decided to build a little bot that would help me to keep track of the spendings using telegram - the application that I already constantly use. So I do not need to get used to something new, and usage of telegram reminds me to add spending records, so it suits my demands. It is hosted right now and I personally use it, tracking the spendings.

### Functionality

As it was written above, the project is a telegram bot that allows to add and observe:

    * spendings categories
    * spending records

Then it allows to display the collected data for the chosen time period for the specific category. Additionally, you can request an EXEL file with the data for the further examination.

The bot is hosted right now on a digital ocean droplet (at least it is hosted on the moment of writing this README), so you are welcome to test it [here](https://t.me/tgSukhanov_bot), but please please don't steal the data, otherwise you will know how much money I spend on beer and delivery food💀💀.

### How it works

In two words, the bot gets text, then it checks if it is a new command or an input for the previously initiated command. Then it transmits the data to the proper goroutine, or creates a new one (if it is a new command). Once the command is done, or a timout reached, the go-routine is shut.

The whole processing and database manipulating happens in separate go-routines, for each go-routine process there is a special session struct that stores a state as an atomic int32, so there should not be any data races. 

I tried to comment all essential parts of the code, so you are welcome to look at the code, I am sure it would be pretty easy to figure out how everything works.

Anyway, it is a full go project, so all the code is inside the go foulder. Then there is a `cmd` foulder with just a main file, nothing interesing there, the db connection and the bot are inited there, and then started.

In the `internal` foulder there are 4 packages:

    * bot
    * repository
    * service
    * utils

`utils` package contains some util functions that are pretty general and not that related to the project.

`repository` is responsible for all code related to the database, there are tests using [test-containers](https://golang.testcontainers.org). The code is just put and retrieve the data from the db.

`service` is the application service. It provides functions wrappers for the repository and functions to create exel files.

`bot` is the most outer arcitecture layer. It represents the whole bot functionality. The `Run()` function there gets the updates from the telegram API and then feeds the recieved text to the proper go-routine, or creates a new one, if needed.