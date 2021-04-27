cmds
```

ssh -tt -i ~/privateKey user@host docker exec -it $(docker ps | grep  unique_text | cut -c1-10) /bin/bash deploy.sh

sed -i.bak 's/127.0.0.1 goreman.default.*/127.0.0.1 goreman.default2/g'  1.host

```

#### requests
```
# login
curl --location --request POST 'http://127.0.0.1:6701/login' \
--header 'Content-Type: application/json' \
--data-raw '{"email": "admin@nutshell.io", "password": "123456"}'

# friends

# connect websocket

# send message
```

