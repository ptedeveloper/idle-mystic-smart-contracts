import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import IdleMysitcHeros from  "../../contracts/IdleMysitcHeros.cdc"

// This script returns the metadata for an NFT in an account's collection.

pub fun main(address: Address, infoID: UInt64): [UInt64] {

    // get the public account object for the token owner
    let owner = getAccount(address)

    let collectionBorrow = owner.getCapability(IdleMysitcHeros.CollectionPublicPath)!
        .borrow<&{IdleMysitcHeros.HerosCollectionPublic}>()
        ?? panic("Could not borrow HerosCollectionPublic")

    // borrow a reference to a specific NFT in the collection
    let heroInfo = collectionBorrow.borrowHeroInfoData(infoID: infoID)
        ?? panic("No such itemID in that collection")
    return heroInfo.childs
}
