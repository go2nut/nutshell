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

RUN mkdir /nutshell /nutshell/_example /nutshell/bin
ADD . /go/src/nutshell/
WORKDIR /go/src/nutshell/

RUN go build -o /go/bin/nutlet nutlet/main.go
RUN chmod 777 /go/bin/nutlet

ENV PATH="/go/bin:${PATH}"

WORKDIR /go/src/nutshell/_example
EXPOSE 80
ENTRYPOINT nutlet
