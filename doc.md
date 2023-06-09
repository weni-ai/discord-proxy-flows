Your users will need a Discord account and a compatible device to send/receive messages. This type of channel is only available if your instance has been configured with the Discord proxy application.

Register your app by following the official instructions.

- Discord has official documentation on how to register your bot.
https://discord.com/developers/docs/getting-started

- Go to the developer dashboard and click New Application and build your app:
https://discord.com/developers/applications

1- Go to the "Bot" tab and click "Add Bot". You will have to confirm by clicking "Yes, do it!"

2- Keep the default settings for Public Bot (checked) and Require OAuth2 code grant (unchecked).

The next step is to copy the token.

Click "Reset Token" and get your token.

3- Add the channel to the Weni platform

In settings go to: Add channel -> Discord

Place your discord_bot_token provided by discord developer page and your bot url.

4- Run your discord proxy

Configure and Run application provided in this repository https://github.com/weni-ai/discord-proxy-flows
to stablish communication between discord flows channel and discord bot


