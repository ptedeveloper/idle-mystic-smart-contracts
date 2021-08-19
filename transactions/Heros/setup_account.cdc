import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import IdleMysitcHeros from  "../../contracts/IdleMysitcHeros.cdc"


// This transaction configures an account to hold Kitty Items.

transaction {
    prepare(signer: AuthAccount) {
        // if the account doesn't already have a collection
        if signer.borrow<&IdleMysitcHeros.Collection>(from: IdleMysitcHeros.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- IdleMysitcHeros.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: IdleMysitcHeros.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&IdleMysitcHeros.Collection{NonFungibleToken.CollectionPublic, IdleMysitcHeros.HerosCollectionPublic}>(IdleMysitcHeros.CollectionPublicPath, target: IdleMysitcHeros.CollectionStoragePath)
            
        }
    }
}
