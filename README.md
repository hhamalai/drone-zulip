# drone-zulip

Drone plugin to publish pipeline events to Zulip chat

Configure drone project secret with Zulip bot apikey ```zulip_bot_apikey``` from Zulip UI.

###Send notifications to Zulip stream
```
steps:
...
- name: notify
  image: ghcr.io/hhamalai/drone-zulip:latest
  environment:
    PLUGIN_URL: https://mychatapp.example/api/v1/messages
    PLUGIN_TYPE: stream
    PLUGIN_RECIPIENT: MyStream
    PLUGIN_TOPIC: builds
    PLUGIN_BOT_EMAIL: drone-bot@mychatapp.example
    PLUGIN_BOT_APIKEY:
      from_secret: zulip_bot_apikey
```

###Send notifications to Zulip user
```
steps:
...
- name: notify
  image: ghcr.io/hhamalai/drone-zulip:latest
  environment:
    PLUGIN_URL: https://mychatapp.example/api/v1/messages
    PLUGIN_TYPE: private
    PLUGIN_RECIPIENT: "[42]"
    PLUGIN_BOT_EMAIL: drone-bot@mychatapp.example
    PLUGIN_BOT_APIKEY:
      from_secret: zulip_bot_apikey
```