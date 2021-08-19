import IdleMysitcHerosMarket from "../../contracts/IdleMysitcHerosMarket.cdc"

transaction(itemID: UInt64) {
    let marketCollection: &IdleMysitcHerosMarket.Collection

    prepare(signer: AuthAccount) {
        self.marketCollection = signer.borrow<&IdleMysitcHerosMarket.Collection>(from: IdleMysitcHerosMarket.CollectionStoragePath)
            ?? panic("Missing or mis-typed IdleMysitcHerosMarket Collection")
    }

    execute {
        let offer <-self.marketCollection.remove(itemID: itemID)
        destroy offer
    }
}
