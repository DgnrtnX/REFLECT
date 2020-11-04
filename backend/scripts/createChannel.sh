export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=$REFLECT/backend/hyperledger/crypto-config/ordererOrganizations/example.com/orderers/orderer1.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=$REFLECT/backend/hyperledger/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export PEER0_ORG2_CA=$REFLECT/backend/hyperledger/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export FABRIC_CFG_PATH=$REFLECT/backend/docker/

export CHANNEL_NAME="newChannel"

setGlobalsForPeer0Org1(){
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=$REFLECT/backend/hyperledger/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
}

setGlobalsForPeer1Org1(){
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=$REFLECT/backend/hyperledger/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:8051
    
}

setGlobalsForPeer0Org2(){
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=$REFLECT/backend/hyperledger/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
    
}

setGlobalsForPeer1Org2(){
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=$REFLECT/backend/hyperledger/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:10051
    
}

createChannel(){
    # rm -rf ../hyperledger/blocks/*
    setGlobalsForPeer0Org1
    
    peer channel create -o localhost:7050 -c $CHANNEL_NAME \
    --ordererTLSHostnameOverride orderer1.example.com \
    -f ../docker/${CHANNEL_NAME}.tx --outputBlock ../hyperledger/blocks/${CHANNEL_NAME}.block \
    --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA
}


joinChannel(){
    setGlobalsForPeer0Org1
    peer channel join -b ../hyperledger/blocks/$CHANNEL_NAME.block
    
    setGlobalsForPeer1Org1
    peer channel join -b ../hyperledger/blocks/$CHANNEL_NAME.block
    
    setGlobalsForPeer0Org2
    peer channel join -b ../hyperledger/blocks/$CHANNEL_NAME.block
    
    setGlobalsForPeer1Org2
    peer channel join -b ../hyperledger/blocks/$CHANNEL_NAME.block
    
}

updateAnchorPeers(){
    setGlobalsForPeer0Org1
    peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer1.example.com \
    -c $CHANNEL_NAME -f ../hyperledger/transactions/${CORE_PEER_LOCALMSPID}anchors.tx \
    --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA
    
    setGlobalsForPeer0Org2
    peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer1.example.com \
    -c $CHANNEL_NAME -f ../hyperledger/transactions/${CORE_PEER_LOCALMSPID}anchors.tx \
    --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA
    
}

createChannel
joinChannel
updateAnchorPeers