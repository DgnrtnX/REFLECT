(cd ./backend/docker && ./createLedger.sh && docker-compose up -d)
sleep 5
(cd ./backend/scripts && ./createChannel.sh)
# sleep 5
# (cd ./backend/chaincycle && ./deployChaincode.sh)
