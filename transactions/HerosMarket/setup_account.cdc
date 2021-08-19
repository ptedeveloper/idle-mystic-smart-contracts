import IdleMysitcHerosMarket from "../../contracts/IdleMysitcHerosMarket.cdc"

// This transaction configures an account to hold SaleOffer items.

transaction {
    prepare(signer: AuthAccount) {

        // if the account doesn't already have a collection
        if signer.borrow<&IdleMysitcHerosMarket.Collection>(from: IdleMysitcHerosMarket.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- IdleMysitcHerosMarket.createEmptyCollection() as! @IdleMysitcHerosMarket.Collection
            
            // save it to the account
            signer.save(<-collection, to: IdleMysitcHerosMarket.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&IdleMysitcHerosMarket.Collection{IdleMysitcHerosMarket.CollectionPublic}>(IdleMysitcHerosMarket.CollectionPublicPath, target: IdleMysitcHerosMarket.CollectionStoragePath)
        }
    }
}
