## Introduction

<p>Built a Distributed Chat system with Golang, Cassandra and Redis. Scalling Websocket servers with redis Pub/Sub and storing message in Cassandra. Made this project for educational purpose. </p>

<b>Detailed Explanation on System Architecture: https://dev.to/leoantony72/distributed-chat-application-22oo</b>

&nbsp;

## Running the application

```bash
docker-compose up -d
```

## CRUD Routes

|  Routes   | Method |  Description   |
| :-------: | :----: | :------------: |
|   /user   |  POST  | Create a User  |
|   /room   |  POST  | Create a Room  |
| /joinroom |  POST  | Join in a Room |

## Websocket Connection

```websocket URL
ws:/localhost/chat?id=username
```

## Examples

- <i>/user</i>

```JSON
{
    "username":"JOHN"
}
```

- <i>/room</i>

```JSON
{
    "name":"anime",
    "user":"JOHN"
}
```

- <i>/joinroom</i>

```JSON
{
    "name":"anime",
    "user":"JOHN"
}
```

- <i>Private Chat - Websocket</i>

```JSON
{
    "msg":"hello Boi",
    "receiver":"2GKtMkzDZDerO2a1gl5IHK6OTPY",
    "is_group":false
}
```

- <i>Group Chat - Websocket</i>

```JSON
{
    "msg":"hello Boi",
    "is_group":true,
    "group_name":"anime"
}
```

---

<i>Hope you will try this out, and leave me your feedback also feel free to improve this project by making a PR😁.</i>
