# Project Structure: agilityfc-bot

The `agilityfc-bot` project is structured as follows:

- **cmd/**: Contains the entry point for the application.
  - **main.go**: The main file where the application execution starts.

- **config/**: Configuration package for the project.
  - **config.go**: Configuration related functions and variables.
  - **vars.go**: Variables used throughout the configuration.

- **internal/bot/**: Core logic for the bot functionality.
  - **bot.go**: Main bot logic and initialization.
  - **commands.go**: Definitions and handling of bot commands.
  - **handlers.go**: Event handlers for the bot.
  - **reactions.go**: Logic for bot reactions.

- **config.json**: Configuration file for the project settings.
