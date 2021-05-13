#### how to start

```
nutctl start

open browser: http://127.0.0.1:8002/token=token_100
```

#### requests example

```

# login
curl --location --request POST 'http://127.0.0.1:8001/login' \
--header 'Content-Type: application/json' \
--data-raw '{"email": "admin@nutshell.io", "password": "123456"}'

# friends

# connect websocket

# send message
```
