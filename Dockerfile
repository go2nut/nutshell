FROM golang:1.13

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# install zsh
RUN apt-get -q -q update && \
    apt-get -y install git jq curl zsh vim telnet etcd
RUN chsh -s /bin/zsh
RUN sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"

#RUN git clone git@github.com:stolk/imcat.git  /go/src/imcat
#RUN cd /go/src/imcat && make && mv imcat /go/bin/imcat
#RUN apt-get install graphviz

# install nginx
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# install git
RUN mkdir /srv/git

# This can speed up future builds because of cache, only rebuild when vendors are
# added.
ENV GOPROXY=http://goproxy.cn
ENV GO111MODULE=on
VOLUME ["/srv/git"]

RUN go install github.com/mattn/goreman
RUN go install github.com/loov/watchrun
#RUN go install github.com/etcd-io/etcd
#RUN go install github.com/etcd-io/etcd/etcdctl

RUN mkdir /nutshell /nutshell/bin /nutshell/_example
ADD . /go/src/nutshell/
WORKDIR /go/src/nutshell/

COPY _example/config.yaml /nutshell/_example/config.yaml
COPY _example/Procfile /nutshell/_example/example.Procfile
RUN go build -o /nutshell/bin/example _example/apps/main.go
RUN chmod 777 /nutshell/bin/example

COPY nutlet/nutlet.sh /nutshell/bin/nutlet.sh
COPY nutlet/nutlet.Procfile /nutshell/nutlet/nutlet.Procfile
RUN go build -o /nutshell/bin/nutlet nutlet/main.go
RUN chmod 777 /nutshell/bin/nutlet

WORKDIR /nutshell/

EXPOSE 80
#ENTRYPOINT cd /nutshell && sh nutlet.sh
ENTRYPOINT watch -n 5 "pwd"
