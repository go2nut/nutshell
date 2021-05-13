## A Universe in a Nutshell for Development
Nutshell provider a lightweight dev/test environment base on the docker. Support request auto route between multi environments.
- Env discovery based on etcd (single env is also support with inner etcd)
- Simple configure to define the request between app
- Integrate goreman for process manage

### how to use

###### edit .nutshell config for your projects

```
➜  nutshell git:(master) ✗ ls _example/.nutshell
apps.Procfile   apps.yaml
```

###### start docker ans env named test1
```
sudo docker run -p 6700-6702:6700-6702/tcp -v /data/code/go_home/src/nutshell:/go/src/nutshell -e "nutshell_env=env1" -e "nutshell_ip=172.31.27.138" -e "nutshell_http_port=6701" -e "nutshell_grpc_port=6702" -e "nutshell_etcd=127.0.0.1" -e "nutshell_ws=/go/src/nutshell/_example" --name test1.nutshell --hostname test1.nutshell  --dns=127.0.0.1  -w /go/src/project1 --rm -it nutshell:latest
```

###### login and add local dns into /etc/resolve.conf 
```
sudo docker exec -it test1.nutshell bash
namespace 127.0.0.1 // add this line as first in /etc/resolve.conf
```

###### start nutctl
```
cd _example
nutctl start
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

```
docker build --tag=nutshell:latest .
```

##### ports usage:
    - 6700 for customize usage like debug
    - 6701 for http proxy port
    - 6702 for grpc proxy port

