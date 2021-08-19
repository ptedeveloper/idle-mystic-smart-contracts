import IdleMysitcHeros from  "../../contracts/IdleMysitcHeros.cdc"

// This scripts returns the number of IdleMysitcHeros currently in existence.

pub fun main(): UInt64 {    
    return IdleMysitcHeros.totalSupply
}
