SYS_CHANNEL="sys-channel"
CHANNEL_NAME="mychannel"

chmod -R u+x ../../*

# Delete existing artifacts
rm -rf ../hyperledger/crypto-config/
rm genesis.block $CHANNEL_NAME.tx
rm -rf ../hyperledger/transactions/*

#Generate Crypto artifactes for organizations
cryptogen generate --config=./crypto-config.yaml --output=./crypto-config/

echo "-----------------"
echo "connecting " $CHANNEL_NAME
echo "-----------------"

# Generate System Genesis block
echo
echo "generating genesis block"
echo
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL  -outputBlock ./genesis.block

# Generate channel configuration block
echo
echo "generating channel transaction"
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./$CHANNEL_NAME.tx -channelID $CHANNEL_NAME

echo
echo "----------  Generating anchor peer update for Org1MSP  ----------"
echo
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

echo
echo "----------  Generating anchor peer update for Org2MSP  ----------"
echo
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP