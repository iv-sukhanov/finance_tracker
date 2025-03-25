# Finance Tracker Bot

## Description

Finance Tracker Bot is a Telegram bot designed to help users manage their financial transactions. It enables users to add, categorize, and analyze their expenses efficiently.

## Table of Contents

1. [Installation & Setup](#installation--setup)
2. [Running the Project](#running-the-project)
3. [Deployment](#deployment)
4. [Logging & Debugging](#logging--debugging)
5. [Database](#database)
6. [Overview](#overview)
7. [Additional Notes](#additional-notes)

## Installation & Setup

### Prerequisites

Ensure the following tools are installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.msys2.org) (optional)

### Clone the Repository

```sh
git clone https://github.com/iv-sukhanov/finance_tracker.git
cd finance_tracker
```

### Configure Environment Variables

Create an `.env/.dev` file and define the required variables:

```env
TELEGRAM_BOT_TOKEN=your_token
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_DB=finance_tracker
DATABASE_URL=postgres://your_user:your_password@db:5432/finance_tracker?sslmode=disable
```

Optional variables:

```env
LOG_LEVEL=(DEBUG|INFO|WARN|ERROR|FATAL|PANIC)
TELEGRAM_DEBUG_MODE=(true|false)
TELEGRAM_USERNAME=your_tg_username
APP_NAME=finance_tracker_bot
```

- `LOG_LEVEL`: Sets the application's log level (default: INFO).
- `TELEGRAM_DEBUG_MODE`: Enables Telegram API debug logs (default: false).
- `TELEGRAM_USERNAME`: Adds your username to internal error messages.
- `APP_NAME`: Appends the app name to logs.

## Running the Project

### Using Docker Compose

```sh
docker compose up -d --build
```

### Using Make

```sh
make compose-up
```

or

```sh
make compose-restart
```

This will initialize the PostgreSQL database, apply migrations, and start the bot.

## Deployment

The project uses GitHub Actions for automated deployment.

### Manual Deployment

1. Add `SSH_PRIVATE_KEY` and `VPS_HOST` to GitHub secrets.
2. Trigger the workflow:

```sh
gh workflow run main.yml
```

Alternatively, set `WORKFLOW_ID=workflow_id` and deploy using:

```sh
make gh-action
```

## Logging & Debugging

- View real-time logs:
    ```sh
    docker compose logs -f
    ```
- View database logs:
    ```sh
    docker compose logs -f db
    ```

## Database

The project uses PostgreSQL for data storage. The database container is managed by Docker Compose, and migrations are applied automatically.

### Schema Overview

![Database Schema](/doc/schema.png)

- **Tables**: `users`, `spending_categories`, `spending_records`
- **Relationships**:
  - `users` → `spending_categories`: One-to-Many
  - `spending_categories` → `spending_records`: One-to-Many

## Overview

### Motivation

As a poor student, I have always wanted to keep track of how much money I spend on different things so that I can later reconsider my expenses and manage my finances more wisely. However, it was too difficult for me to remember all my expenses and update the Excel table on my laptop every day. Bank applications do not provide enough details about my spending (especially BoC 🤭🤭), and other budgeting apps felt like too much overhead, so they didn’t really suit me.

So, I decided to build a little bot to help me track my spending using Telegram — an app I already use constantly. This way, I don’t have to get used to something new, and using Telegram reminds me to add spending records, making it a perfect fit for my needs. The bot is currently hosted, and I personally use it to track my expenses.

### Features

The bot allows users to:

- Add and view spending categories.
- Record and analyze expenses.
- Generate Excel reports for detailed analysis.

The bot is hosted on a DigitalOcean droplet and is available for testing [here](https://t.me/tgSukhanov_bot). But please please don't steal the data, otherwise you will know how much money I spend on beer and delivery food ;)

### How It Works

The bot processes user input by identifying commands and delegating tasks to appropriate goroutines. Each goroutine manages its session state using atomic operations to prevent data races.

The project structure includes:

- `cmd`: Contains the main entry point.
- `internal`:
  - `bot`: Handles Telegram bot functionality.
  - `repository`: Manages database interactions.
  - `service`: Provides business logic and utility functions.
  - `utils`: Contains general-purpose helper functions.

## Additional Notes

Feedback on the code, architecture, or any other aspect is highly appreciated. Your suggestions will help me improve my skills.
