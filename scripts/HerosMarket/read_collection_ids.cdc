import IdleMysitcHerosMarket from 0xNonFungibleToken

// This script returns an array of all the NFT IDs for sale 
// in an account's SaleOffer collection.

pub fun main(address: Address): [UInt64] {
    let marketCollectionRef = getAccount(address)
        .getCapability<&IdleMysitcHerosMarket.Collection{IdleMysitcHerosMarket.CollectionPublic}>(
            IdleMysitcHerosMarket.CollectionPublicPath
        )
        .borrow()
        ?? panic("Could not borrow market collection from market address")
    
    return marketCollectionRef.getSaleOfferIDs()
}
