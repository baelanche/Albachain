# Albachain

## How to run

1. run ./generate.sh
2. run ./start.sh
3. run ./albaPublish.sh
4. run node enrollAdmin.js
5. run node server.js

## Clean up after run

1. run ./teardown.sh
2. clean the mongoDB
    * `$ sudo mongo`
    * `> use test`
    * `> db.users.remove({})`