SYS_CHANNEL="sys-channel"
CHANNEL_NAME="newChannel"

chmod -R 0755 ../*

# Delete existing artifacts
rm -rf ../hyperledger/crypto-config/*
rm genesis.block $CHANNEL_NAME.tx
rm -rf ../hyperledger/transactions/*

#Generate Crypto artifactes for organizations
cryptogen generate --config=./crypto-config.yml --output=../hyperledger/crypto-config/

# Generate System Genesis block
echo
echo "----------  Generating genesis block  ----------"
echo
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL  -outputBlock ./genesis.block

# Generate channel configuration block
echo
echo "----------  Generating channel transaction  ----------"
echo
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./$CHANNEL_NAME.tx -channelID $CHANNEL_NAME

echo
echo "----------  Generating anchor peer update for Org1MSP  ----------"
echo
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ../hyperledger/transactions/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

echo
echo "----------  Generating anchor peer update for Org2MSP  ----------"
echo
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ../hyperledger/transactions/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP