(cd ./backend/docker && ./createLedger.sh && docker-compose up -d)
(cd ./backend/scripts && ./createChannel.sh)
(cd ./backend/chaincycle && ./deployChaincode.sh)
