import FungibleToken from "../../contracts/FungibleToken.cdc"
import IdleMysitcToken from "../../contracts/IdleMysitcToken.cdc"

// This transaction is a template for a transaction
// to add a Vault resource to their account
// so that they can use the IdleMysitcToken

transaction {

    prepare(signer: AuthAccount) {

        if signer.borrow<&IdleMysitcToken.Vault>(from: IdleMysitcToken.VaultStoragePath) == nil {
            // Create a new IdleMysitcToken Vault and put it in storage
            signer.save(<-IdleMysitcToken.createEmptyVault(), to: IdleMysitcToken.VaultStoragePath)

            // Create a public capability to the Vault that only exposes
            // the deposit function through the Receiver interface
            signer.link<&IdleMysitcToken.Vault{FungibleToken.Receiver}>(
                IdleMysitcToken.ReceiverPublicPath,
                target: IdleMysitcToken.VaultStoragePath
            )

            // Create a public capability to the Vault that only exposes
            // the balance field through the Balance interface
            signer.link<&IdleMysitcToken.Vault{FungibleToken.Balance}>(
                IdleMysitcToken.BalancePublicPath,
                target: IdleMysitcToken.VaultStoragePath
            )
        }
    }
}
