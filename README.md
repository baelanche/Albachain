# Albachain

Blockchain work record management with Hyperledger fabric

## Precondition

* curl, docker, docker-compose, go, nodejs, python
* hyperledger fabric-docker images are installed
* GOPATH are configured
    * `$ echo $PATH`
* hyperledger bineries are installed (cryptogen, configtxgen ... etcs)

## How to run

1. run ./generate.sh (once)
    * create `./network/config`, `./network/crypto-config`
2. run ./start.sh
3. run ./albaPublish.sh
    * install, instantiate, (test) ...
4. run node enrollAdmin.js
5. run node server.js
    * You can register the user via 'join' on web server
    * then, the wallet is created at `./application/wallet`

## Clean up after run

1. run ./teardown.sh
    * clean up docker container, docker images
    * down docker network
    * remove all wallet (line : 35)
2. clean the mongoDB
    * `$ sudo mongo`
    * `> use test`
    * `> db.users.remove({})`