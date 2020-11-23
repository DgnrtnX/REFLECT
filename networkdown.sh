rm -rf ./backend/hyperledger/crypto-config/*
rm ./backend/hyperledger/blocks/* 
rm -rf ./backend/hyperledger/transactions/*
rm -rf ./backend/chaincycle/Reflect.tar.gz
rm -rf ./backend/chaincycle/log.txt
(cd ./backend/docker && docker-compose down)
