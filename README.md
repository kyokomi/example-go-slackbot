example-go-slackbot
===============================

[![wercker status](https://app.wercker.com/status/4956e5bb2ee4285dd24b02197e9d0975/s/master "wercker status")](https://app.wercker.com/project/byKey/4956e5bb2ee4285dd24b02197e9d0975)

## Application Environment
- `SLACK_BOT_TOKEN`: required
- `DOCOMO_APIKEY`: `github.com/kyokomi/go-docomo/docomo` plugin
- `REDISTOGO_URL`: `redis://<id>:<password>@<host>:<port>` https://redistogo.com/

## Wercker Environment
- DockerHub push
    - `DOCKER_HUB_USERNAME`
    - `DOCKER_HUB_PASSWORD`
    - `DOCKER_HUB_REPOSITORY`: (example) `kyokomi/go-slackbot`
- Arukas restart
    - `ARUKAS_JSON_API_TOKEN`
    - `ARUKAS_JSON_API_SECRET`
    - `ARUKAS_CONTAINER_ID`

## License
[MIT](https://github.com/kyokomi/example-go-slackbot/blob/master/LICENSE)
