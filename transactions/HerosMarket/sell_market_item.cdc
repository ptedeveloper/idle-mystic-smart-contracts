import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import IdleMysitcToken from "../../contracts/IdleMysitcToken.cdc"
import IdleMysitcHeros from "../../contracts/IdleMysitcHeros.cdc"
import IdleMysitcHerosMarket from "../../contracts/IdleMysitcHerosMarket.cdc"

transaction(itemID: UInt64, price: UFix64) {
    let stoneVault: Capability<&IdleMysitcToken.Vault{FungibleToken.Receiver}>
    let herosCollection: Capability<&IdleMysitcHeros.Collection{NonFungibleToken.Provider, IdleMysitcHeros.HerosCollectionPublic}>
    let marketCollection: &IdleMysitcHerosMarket.Collection

    prepare(signer: AuthAccount) {
        // we need a provider capability, but one is not provided by default so we create one.
        let HerosCollectionProviderPrivatePath = /private/herosCollectionProvider

        self.stoneVault = signer.getCapability<&IdleMysitcToken.Vault{FungibleToken.Receiver}>(IdleMysitcToken.ReceiverPublicPath)!
        assert(self.stoneVault.borrow() != nil, message: "Missing or mis-typed IdleMysitcToken receiver")

        if !signer.getCapability<&IdleMysitcHeros.Collection{NonFungibleToken.Provider, IdleMysitcHeros.HerosCollectionPublic}>(HerosCollectionProviderPrivatePath)!.check() {
            signer.link<&IdleMysitcHeros.Collection{NonFungibleToken.Provider, IdleMysitcHeros.HerosCollectionPublic}>(HerosCollectionProviderPrivatePath, target: IdleMysitcHeros.CollectionStoragePath)
        }

        self.herosCollection = signer.getCapability<&IdleMysitcHeros.Collection{NonFungibleToken.Provider, IdleMysitcHeros.HerosCollectionPublic}>(HerosCollectionProviderPrivatePath)!
        assert(self.herosCollection.borrow() != nil, message: "Missing or mis-typed HerosCollection provider")

        self.marketCollection = signer.borrow<&IdleMysitcHerosMarket.Collection>(from: IdleMysitcHerosMarket.CollectionStoragePath)
            ?? panic("Missing or mis-typed IdleMysitcHerosMarket Collection")
    }

    execute {
        let offer <- IdleMysitcHerosMarket.createSaleOffer (
            sellerItemProvider: self.herosCollection,
            itemID: itemID,
            infoID: self.herosCollection.borrow()!.borrowHero(id: itemID)!.infoID,
            sellerPaymentReceiver: self.stoneVault,
            price: price
        )
        self.marketCollection.insert(offer: <-offer)
    }
}
