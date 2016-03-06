# Foxx-installer

Foxx-installer is a standalone utility to install/upgrade/replace  Arangodb Foxx applications.
It is intended to be used within a Docker container such that Foxx apps can be deployed as docker container.

## Usage

```
foxx-installer install \
    --database=<dbname> \
    --server-url<arangodburl> \
    --app-path=<foxxapppath> \
    --mountpoint=/yourapp
    [--replace]
```

## Environment variables

The following environment variable can be used instead of command line arguments.

- `FI_DATABASE` - Replaces `--database`
- `FI_SERVER_URL` - Replaces `--server-url`
- `FI_APP_PATH` - Replaces `--app-path`
- `FI_MOUNTPOINT` - Replaces `--mountpoint`
- `FI_REPLACE` - Replaces `--replace` (if set to "1")

## Docker usage

This is a sample Dockerfile that uses `foxx-installer` to install a app location in `./app/`.

```
FROM pulcy/foxx-installer

COPY ./app/ /app/

ENV FI_APP_PATH="/app/"
ENV FI_MOUNTPOINT="/myapp"

CMD ["install", "--database=mydb", "--server-url=http://arangodb.local:8529"]
```

Build the docker image and run it with a simple `docker run yourimagename`.
