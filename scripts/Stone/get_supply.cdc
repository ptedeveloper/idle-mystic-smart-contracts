import IdleMysitcToken from "../../contracts/IdleMysitcToken.cdc"

// This script returns the total amount of IdleMysitcToken currently in existence.

pub fun main(): UFix64 {

    let supply = IdleMysitcToken.totalSupply

    log(supply)

    return supply
}
