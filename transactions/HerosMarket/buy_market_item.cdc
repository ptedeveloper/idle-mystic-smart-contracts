import FungibleToken from "../../contracts/FungibleToken.cdc"
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import IdleMysitcToken from "../../contracts/IdleMysitcToken.cdc"
import IdleMysitcHeros from "../../contracts/IdleMysitcHeros.cdc"
import IdleMysitcHerosMarket from "../../contracts/IdleMysitcHerosMarket.cdc"

transaction(itemID: UInt64, marketCollectionAddress: Address) {
    let paymentVault: @FungibleToken.Vault
    let herosCollection: &IdleMysitcHeros.Collection{NonFungibleToken.Receiver}
    let marketCollection: &IdleMysitcHerosMarket.Collection{IdleMysitcHerosMarket.CollectionPublic}

    prepare(signer: AuthAccount) {
        self.marketCollection = getAccount(marketCollectionAddress)
            .getCapability<&IdleMysitcHerosMarket.Collection{IdleMysitcHerosMarket.CollectionPublic}>(
                IdleMysitcHerosMarket.CollectionPublicPath
            )
            .borrow()
            ?? panic("Could not borrow market collection from market address")

        let saleItem = self.marketCollection.borrowSaleItem(itemID: itemID)
                    ?? panic("No item with that ID")
        let price = saleItem.price

        let mainStoneVault = signer.borrow<&IdleMysitcToken.Vault>(from: IdleMysitcToken.VaultStoragePath)
            ?? panic("Cannot borrow IdleMysitcToken vault from acct storage")
        self.paymentVault <- mainStoneVault.withdraw(amount: price)

        self.herosCollection = signer.borrow<&IdleMysitcHeros.Collection{NonFungibleToken.Receiver}>(
            from: IdleMysitcHeros.CollectionStoragePath
        ) ?? panic("Cannot borrow IdleMysitcHeros collection receiver from acct")
    }

    execute {
        self.marketCollection.purchase(
            itemID: itemID,
            buyerCollection: self.herosCollection,
            buyerPayment: <- self.paymentVault
        )
    }
}
