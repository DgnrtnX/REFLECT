channelName="newChannel"
# Delete existing artifacts
# (cd ../hyperledger/crypto-config/ && chmod u+x *)
rm -rf ../hyperledger/crypto-config/*
# rm genesis.block $channelName.tx

#Generate Crypto artifactes for organizations
cryptogen generate --config=crypto-config.yml --output=../hyperledger/crypto-config/

# channel name defaults to "$channelName"
echo "---------------------"
echo "connecting " $channelName
echo "---------------------"


##########################
#Errors for comands below
#Try to debug the issue
#and uncomment the single commented lines
###########################

# # Generate System Genesis block
# echo "Creating genesis block"

# configtxgen -profile OrdererGenesis -configPath . -channelID sys-channel  -outputBlock genesis.block

# # Generate channel configuration block
# echo "Creating transactions "
# configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ../hyperledger/transactions/newChannel.tx -channelID $channelName

# echo "#######    Generating anchor peer update for Org1    ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ../hyperledger/transactions/Org1MSPanchors.tx -channelID $channelName -asOrg Org1MSP

# echo "#######    Generating anchor peer update for Org2    ##########"
# configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ../hyperledger/transactions/Org2MSPanchors.tx -channelID $channelName -asOrg Org2MSP