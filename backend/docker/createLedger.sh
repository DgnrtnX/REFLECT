SYS_CHANNEL="sys-channel"
CHANNEL_NAME="supplychainchannel"

echo
echo "======================================= Generating crypto material =========================================="
echo
# generate crypto material
cryptogen generate --config=./crypto-config.yml --output=../hyperledger/crypto-config/

if [ "$?" -ne 0 ]; then
    echo "Failed to generate crypto material!"
    exit 1
fi
echo

echo "======================================= Creating genesis block ==============================================="
echo
# generate genesis block for orderer
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL  -outputBlock ../hyperledger/blocks/genesis.block

if [ "$?" -ne 0 ]; then
    echo "Failed to generate orderer genesis block!"
    exit 1
fi
echo

echo "============================= Generating channel configuration transaction ==================================="
echo
# generate channel configuration transaction
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ../hyperledger/transactions/$CHANNEL_NAME.tx -channelID $CHANNEL_NAME

if [ "$?" -ne 0 ]; then
    echo "Failed to generate channel configuration transaction!"
    exit 1
fi
echo

echo "============================= Generating anchor peer update for Org1MSP =============================="
echo
# generate anchor peer transaction
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ../hyperledger/transactions/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

if [ "$?" -ne 0 ]; then
    echo "Failed to generate anchor peer update for Org1MSP!"
    exit 1
fi
echo

echo "============================= Generating anchor peer update for Org2MSP ====================================="
echo
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ../hyperledger/transactions/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP

if [ "$?" -ne 0 ]; then
    echo "Failed to generate anchor peer update for Org2MSP!"
    exit 1
fi
echo