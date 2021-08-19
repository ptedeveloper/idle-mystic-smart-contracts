import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import IdleMysitcHeros from  "../../contracts/IdleMysitcHeros.cdc"

// This script returns the size of an account's IdleMysitcHeros collection.

pub fun main(address: Address): Int {
    let account = getAccount(address)

    let collectionRef = account.getCapability(IdleMysitcHeros.CollectionPublicPath)!
        .borrow<&{NonFungibleToken.CollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.getIDs().length
}
