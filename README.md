# Albachain

Blockchain work record management with Hyperledger fabric

## Precondition

* curl, docker, docker-compose, go, nodejs, python
* hyperledger fabric-docker images are installed
* GOPATH are configured
    * `$ echo $PATH`
* hyperledger bineries are installed (cryptogen, configtxgen ... etcs)
* mongoDB

## How to run

1. run `./generate.sh` (once)
    * create `./network/config`, `./network/crypto-config`
2. run `./start.sh`
3. run `./albaPublish.sh`
    * install, instantiate, (test) ...
4. run `node enrollAdmin.js`
5. run `node server.js`
    * You can register the user via 'join' on web server
    * then, the wallet is created at `./application/wallet`

## Clean up after run

1. run `./teardown.sh`
    * clean up docker container, docker images
    * down docker network
    * remove all wallets
2. clean the mongoDB
    * `$ sudo mongo`
    * `> use test`
    * `> db.users.drop()` or `> db.users.remove({})`

<hr/>

## Trouble Shooting

### MongoDB 설치 오류

~~해당 문제는 Ubuntu 의 버전에 따라 차이가 있다.~~  
~~기존 Ubuntu 에 설치되어 있는 mongodb 와 설치하려는 mongodb 의 충돌로 인하여 동작이 제대로 되지 않을 때 수행한다.~~  
처음에 받은 이미지에 MongoDB 가 이미 설치되어 있었던 것으로 보인다. 버전 충돌로 인해 삭제 후 재설치했다.

1. MongoDB 삭제
```
$ apt remove mongodb-org*
$ apt-get purge mongodb-org*
$ rm -r /var/log/mongodb
$ rm -r /var/lib/mongodb
$ rm -r /etc/mongodb.conf
```

2. 확인 (+삭제)
   * `$ apt list --installed | grep mongo` 를 통해 mongodb 가 남아있는지 확인한다.
   * `$ dpkg --list` 를 통해 mongodb 관련 패키지를 찾아 제거한다.
   * ex) `$ sudo apt-get --purge remove mongodb*`

3. 재설치

