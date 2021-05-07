## A Universe in a Nutshell for Development
Nutshell provider a lightweight dev/test environment base on the docker. Support request auto route between multi environments.
- Env discovery based on etcd (single env is also support with inner etcd)
- Simple configure to define the request between app
- Integrate goreman for process manage

### how to use

##### full example

#### quick start

###### start docker
```
docker run -d -p 10.0.1.38:9001:80 -v /Users/sunwei/go/src/nutshell/workspace -e "nutshell_env=env1" -e "nutshell_ip=172.168.0.12" -e "nutshell_http_port=8001" -e "nutshell_grpc_port=8002" -e "nutshell_etcd=1.1.1.1" --name test1.nutshell --hostname test1.nutshell -w "/workspace" --rm -it nutshell:latest bash


```

###### config
```
etcd path: /nutshell/envs/dev1
conent: {address:"192.168.2.10", port:8082}
```

###### login and dev
```
docker exec -it test1.nutshell zsh
```

## how it works

#### config app and env

#### request route

app proxy:
 - match the app with url prefix, then route the request to target address
 - the app defined in local is 

env proxy:
- all *.nutshell pattern host will be parsed to 127.0.0.1 via local dns server on port 53
- then the request will be route into a grpc proxy on 80
- the proxy detect the url prefix (/shard.User) and host, 
  - use url prefix to match app 
  - use host to match env

#### setup apps

## how to build
docker build --tag=nutshell:latest .

sudo docker run -p 172.31.27.138:9001:80 -v /data/code/go_home/src/nutshell/workspace -e "nutshell_env=env1" -e "nutshell_ip=172.168.0.12" -e "nutshell_http_port=8001" -e "nutshell_grpc_port=8002" -e "nutshell_etcd=1.1.1.1" --name test1.nutshell --hostname test1.nutshell -w "/nutshell" --rm -it nutshell:latest bash

sudo docker 
