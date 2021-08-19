import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import IdleMysitcHeros from  "../../contracts/IdleMysitcHeros.cdc"

// This transction uses the NFTMinter resource to mint a new NFT.
//
// It must be run with the account that has the minter resource
// stored at path /storage/NFTMinter.

transaction(recipient: Address, birth: String, batch:Int, first:Bool) {
    
    // local variable for storing the minter reference
    let minter: &IdleMysitcHeros.NFTMinter
    let adminRef: &IdleMysitcHeros.Admin

    prepare(signer: AuthAccount) {
            
        self.adminRef = signer.borrow<&IdleMysitcHeros.Admin>(from: IdleMysitcHeros.AdminStoragePath)
            ?? panic("No admin resource in storage")

        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.borrow<&IdleMysitcHeros.NFTMinter>(from: IdleMysitcHeros.MinterStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")

    }

    execute {

        let heroInfo = self.adminRef.createHeroInfoData(birth: birth, batch:batch, first:first)

        // // get the public account object for the recipient
        let recipient = getAccount(recipient)

        // // borrow the recipient's public NFT collection reference
        let receiver = recipient
            .getCapability(IdleMysitcHeros.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")
        
        // // mint the NFT and deposit it to the recipient's collection
        self.minter.mintNFT(recipient: receiver)
    }
}
