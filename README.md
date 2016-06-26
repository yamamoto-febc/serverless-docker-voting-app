# Serverless Docker Example Voting App

This is a serverless app built with Docker. Read more in the [Serverless Docker repository](https://github.com/bfirsh/serverless-docker).

## Architecture

It consists of a simple entrypoint server that listens for HTTP requests. All of the other functionality of the app is run on-demand as Docker containers for each HTTP request:

 - **vote**: The voting web app, as a CGI container that serves a single HTTP request.
 - **record-vote-task**: A container which processes a vote in the background, run by the vote app.
 - **result**: The result web app, as a CGI container.

## Running

このディレクトリにて以下を実行:

```bash
# swarmバックエンド用KVSマシンの作成(KVSにはconsulを使う)
$ docker-machine create -d virtualbox swarm-kvs
# consul起動
$ eval $(docker-machine env swarm-kvs)
$ docker run -d -p "8500:8500" \
             -h "consul" progrium/consul -server -bootstrap

# swarm マスターノード作成
$ docker-machine create -d virtualbox \
                 --swarm --swarm-master \
                 --swarm-discovery="consul://$(docker-machine ip swarm-kvs):8500" \
                 --engine-opt="cluster-store=consul://$(docker-machine ip swarm-kvs):8500" \
                 --engine-opt="cluster-advertise=eth1:2376" \
                 swarm-master

# swarm エージェントノード作成1
$ docker-machine create -d virtualbox \
                 --swarm \
                 --swarm-discovery="consul://$(docker-machine ip swarm-kvs):8500" \
                 --engine-opt="cluster-store=consul://$(docker-machine ip swarm-kvs):8500" \
                 --engine-opt="cluster-advertise=eth1:2376" \
                 swarm-agent01

# swarmエージェントノード作成2
$ docker-machine create -d virtualbox \
                 --swarm \
                 --swarm-discovery="consul://$(docker-machine ip swarm-kvs):8500" \
                 --engine-opt="cluster-store=consul://$(docker-machine ip swarm-kvs):8500" \
                 --engine-opt="cluster-advertise=eth1:2376" \
                 swarm-agent02

#  イメージのビルド(master)
$ eval $(docker-machine env swarm-master)
$ make build
#  イメージのビルド(agent01)
$ eval $(docker-machine env swarm-agent01)
$ make build
#  イメージのビルド(agent02)
$ eval $(docker-machine env swarm-agent02)
$ make build

# swarm masterへの接続
$ eval $(docker-machine env --swarm swarm-master)

# 実行(docker-compose up)
$ make

```

起動後、別のコンソールなどで`entrypoint`コンテナが起動しているホストのIPアドレスを調べて、以下へアクセス
  - 投票画面: `http://[調べたIP]/vote`
  - 結果画面: `http://[調べたIP]/result`

