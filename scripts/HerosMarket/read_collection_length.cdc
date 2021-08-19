import IdleMysitcHerosMarket from 0xNonFungibleToken

// This script returns the size of an account's SaleOffer collection.

pub fun main(marketCollectionAddress: Address): Int {
    // let acct = getAccount(account)
    let marketCollectionRef = getAccount(marketCollectionAddress)
        .getCapability<&IdleMysitcHerosMarket.Collection{IdleMysitcHerosMarket.CollectionPublic}>(
             IdleMysitcHerosMarket.CollectionPublicPath
        )
        .borrow()
        ?? panic("Could not borrow market collection from market address")
    
    return marketCollectionRef.getSaleOfferIDs().length
}
