# Discord Proxy Flows

This is an application that interacts with Discord and FLows, acting as a proxy between them to allow receiving messages and media files sent between a Discord bot and FLows.

## Configuration

The application uses environment variables to configure its behavior. The following environment variables must be set:

- `DISCORD_PROXY_PORT`: The port on which the application runs. Default value is 8000.

- `DISCORD_PROXY_BOT_TOKEN`: The authentication token of the Discord bot. Obtain the bot token from the Discord developer control panel.

- `DISCORD_PROXY_FLOWS_URL`: The URL of FLows that the application needs to access. Make sure to provide a valid URL.

- `DISCORD_PROXY_CHANNEL_UUID`: The UUID of the Discord channel in FLows.

Make sure to set these environment variables before running the application.

## Running the Application

Make sure you have Go installed on your machine. You can run the application using the following command:

```
go run cmd/main.go
```

This will start the application, and it will begin to receive messages and media files sent by the bot and the Discord channel in FLows.
