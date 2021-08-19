import IdleMysitcToken from "../../contracts/IdleMysitcToken.cdc"
import FungibleToken from "../../contracts/FungibleToken.cdc"
// This script returns an account's IdleMysitcToken balance.

pub fun main(address: Address): UFix64 {
    let account = getAccount(address)
    
    let vaultRef = account.getCapability(IdleMysitcToken.BalancePublicPath)!.borrow<&IdleMysitcToken.Vault{FungibleToken.Balance}>()
        ?? panic("Could not borrow Balance reference to the Vault")

    return vaultRef.balance
}
