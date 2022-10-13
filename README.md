# lost_ark_server_status
Gets changes in the statuses of the Lost Ark servers by scraping the [official site](https://www.playlostark.com/en-us/support/server-status). Posts changes to a given [Discord Webhook](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks).

## Config file
|Key|Default|Description|
|---|-------|-----------|
|`webhook_url`|N/A (required)|The Webhook URL to send status changes to|
|`wait_between_requests`|`30000000000` (30s)|The time in nanoseconds to wait between requests|
|`log_level`|`info`|Logging level (see [golog](https://github.com/kataras/golog))|
|`filter`|`[]`|Which servers to post the changes of to Discord. List of strings for exact server names (e.g. ["Mari", "Valtan"]).|
