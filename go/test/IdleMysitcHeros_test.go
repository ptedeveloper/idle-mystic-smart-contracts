package test

import (
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	nft_contracts "github.com/onflow/flow-nft/lib/go/contracts"
	"log"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

//
const (
	idleMysticIdleMysticPath                   = "../../contracts/IdleMysitcHeros.cdc"
	idleMysticSetupAccountPath                 = "../../transactions/Heros/setup_account.cdc"
	idleMysticMintPath                         = "../../transactions/Heros/mint_hero_extend.cdc"
	idleMysticTransferHeroPath                 = "../../transactions/Heros/transfer_hero_item.cdc"
	idleMysticHeroInheritPath                 = "../../transactions/Heros/inherit.cdc"
	idleMysticInspectIdleMysitcHerosSupplyPath = "../../scripts/Heros/read_hero_items_supply.cdc"
	idleMysticInspectCollectionLenPath         = "../../scripts/Heros/read_collection_length.cdc"
	idleMysticInspectCollectionIdsPath         = "../../scripts/Heros/read_collection_ids.cdc"
	idleMysticHerosInfoPath                    = "../../scripts/Heros/read_heroinfo_dna.cdc"
)

//


func IdleMysticDeployContracts(b *emulator.Blockchain, t *testing.T) (flow.Address, flow.Address, crypto.Signer) {
	accountKeys := test.AccountKeyGenerator()

	// should be able to deploy a contract as a new account with no keys
	nftCode := loadNonFungibleToken()
	nftAddr, err := b.CreateAccount(
		nil,
		[]sdktemplates.Contract{
			{
				Name:   "NonFungibleToken",
				Source: string(nftCode),
			},
		},
	)
	require.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// should be able to deploy a contract as a new account with one key
	// idleMysitcAccountKey, idleMysitcSigner := accountKeys.NewWithSigner()
	idleMysitcAccountKey, idleMysitcSigner := accountKeys.NewWithSigner()
	idleMysitcCode := loadIdleMysitcHeros(nftAddr.String())
	// idleMysitcCode := loadidleMysitc(nftAddr.String())
	// idleMysitcAddr, err := b.CreateAccount(
	idleMysitcAddr, err := b.CreateAccount(
		[]*flow.AccountKey{idleMysitcAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "IdleMysitcHeros",
				Source: string(idleMysitcCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// simplify the workflow by having the contract address also be our initial test collection
	idleMysitcSetupAccount(t, b, idleMysitcAddr, idleMysitcSigner, nftAddr, idleMysitcAddr)

	return nftAddr, idleMysitcAddr, idleMysitcSigner
}

//init account colloction
func IdleMysticSetupHerosAccount(t *testing.T, b *emulator.Blockchain, SaleuserAddress sdk.Address, SaleUserSigner crypto.Signer, nftAddr sdk.Address, idleMysitcAddr sdk.Address) {
	tx := flow.NewTransaction().
		SetScript(GenerateSetupAccountScript(nftAddr.String(), idleMysitcAddr.String())).
		SetGasLimit(1000).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(SaleuserAddress)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, SaleuserAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), SaleUserSigner},
		false,
	)
}
func CreateAccount(t *testing.T, b *emulator.Blockchain, nftAddr sdk.Address, idleMysitcAddr sdk.Address) (sdk.Address, crypto.Signer) {
	SaleuserAddress, SaleUserSigner, _ := createAccount(t, b)
	IdleMysticSetupHerosAccount(t, b, SaleuserAddress, SaleUserSigner, nftAddr, idleMysitcAddr)
	return SaleuserAddress, SaleUserSigner
}
func IdleMysticMint(t *testing.T, b *emulator.Blockchain, nftAddr, idleMysitcAddr,acceptAddr flow.Address, idleMysticSigner crypto.Signer, ) {
	tx := flow.NewTransaction().
		SetScript(MintIdleMysticHeroScript(nftAddr.String(), idleMysitcAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(idleMysitcAddr)
	//recipientAddr := cadence.NewAddress(adminAddr)
	_birth := time.Now().Local().String()

	tx.AddArgument(cadence.NewAddress(acceptAddr))
	tx.AddArgument(cadence.NewString(_birth))
	tx.AddArgument(cadence.NewInt(12))
	tx.AddArgument(cadence.NewBool(true))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), idleMysticSigner},
		false,
	)
}
func IdleMysticTransfer(b *emulator.Blockchain, t *testing.T, nftAddr, idleMysitcAddr flow.Address, idleMysticSigner crypto.Signer, NFTID uint64, recipientAddr flow.Address, shouldFail bool) {
	tx := flow.NewTransaction().
		SetScript(idleMysticGenerateTransferidleMysitcScript(nftAddr.String(), idleMysitcAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(idleMysitcAddr)
	tx.AddArgument(cadence.NewAddress(recipientAddr))
	tx.AddArgument(cadence.NewUInt64(NFTID))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), idleMysticSigner},
		shouldFail,
	)
}


////
func replaceIdleMysticAddressPlaceholders(code, nftAddress, idleMysitcAddress string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			nftAddress:        nftAddressPlaceholder,
			idleMysitcAddress: idleMysitcAddressPlaceHolder,
		},
	))
}

func loadNonFungibleToken() []byte {
	return nft_contracts.NonFungibleToken()
}

func idleMysitcSetupAccount(
	t *testing.T, b *emulator.Blockchain,
	SaleuserAddress sdk.Address, SaleUserSigner crypto.Signer, nftAddr sdk.Address, idleMysitcAddr sdk.Address, ) {
	tx := flow.NewTransaction().
		SetScript(GenerateSetupAccountScript(nftAddr.String(), idleMysitcAddr.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(SaleuserAddress)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, SaleuserAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), SaleUserSigner},
		false,
	)
}
func loadIdleMysitcHeros(nftAddr string) []byte {

	script:= []byte(replaceImports(
		string(readFile(idleMysticIdleMysticPath)),
		map[string]*regexp.Regexp{
			nftAddr: nftAddressPlaceholder,
		},
	))
	return	script
}
func GenerateSetupAccountScript(nftAddr, idleMysitcAddr string) []byte {
	return replaceIdleMysticAddressPlaceholders(
		string(readFile(idleMysticSetupAccountPath)),
		nftAddr,
		idleMysitcAddr,
	)
}

func MintIdleMysticHeroScript(nftAddr, idleMysitcAddr string) []byte {
	return replaceIdleMysticAddressPlaceholders(
		string(readFile(idleMysticMintPath)),
		nftAddr,
		idleMysitcAddr,
	)
}

func idleMysticGenerateTransferidleMysitcScript(nftAddr, idleMysitcAddr string) []byte {
	return replaceIdleMysticAddressPlaceholders(
		string(readFile(idleMysticTransferHeroPath)),
		nftAddr,
		idleMysitcAddr,
	)
}
func GenerateInheritHerosScript(nftAddr, idleMysitcAddr string) []byte {
	return replaceIdleMysticAddressPlaceholders(
		string(readFile(idleMysticHeroInheritPath)),
		nftAddr,
		idleMysitcAddr,
	)
}

func idleMysticGenerateInspectIdleMysitcHerosSupplyScript(nftAddr, idleMysitcAddr string) []byte {
	// This scripts returns the number of IdleMysitcHeros currently in existence.

	return replaceIdleMysticAddressPlaceholders(
		string(readFile(idleMysticInspectIdleMysitcHerosSupplyPath)),
		nftAddr,
		idleMysitcAddr,
	)
}

func idleMysticGenerateInspectCollectionLenScript(nftAddr, idleMysitcAddr string) []byte {
	// This script returns the size of an account's IdleMysitcHeros collection.

	return replaceIdleMysticAddressPlaceholders(
		string(readFile(idleMysticInspectCollectionLenPath)),
		nftAddr,
		idleMysitcAddr,
	)
}

func idleMysticGenerateInspectCollectionIdsScript(nftAddr, idleMysitcAddr string) []byte {
	//collection contain nft  ids
	return replaceIdleMysticAddressPlaceholders(
		string(readFile(idleMysticInspectCollectionIdsPath)),
		nftAddr,
		idleMysitcAddr,
	)
}
func HeroinfoScript(nftAddr, idleMysitcAddr string) []byte {
	//collection contain nft  ids
	return replaceIdleMysticAddressPlaceholders(
		string(readFile(idleMysticHerosInfoPath)),
		nftAddr,
		idleMysitcAddr,
	)
}
func InheritHeros(b *emulator.Blockchain, t *testing.T, nftAddr, idleMysitcAddr flow.Address, idleMysticSigner crypto.Signer, parentAID,parentBID uint64, recipientAddr flow.Address, shouldFail bool) {
	tx := flow.NewTransaction().
		SetScript(GenerateInheritHerosScript(nftAddr.String(), idleMysitcAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(idleMysitcAddr)
	tx.AddArgument(cadence.NewAddress(recipientAddr))
	tx.AddArgument(cadence.NewString(time.Now().Local().String()))
	tx.AddArgument(cadence.NewUInt64(parentAID))
	tx.AddArgument(cadence.NewUInt64(parentBID))
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), idleMysticSigner},
		shouldFail,
	)
}

//test func

func TestIdleMysticDeployContracts(t *testing.T) {
	b := newEmulator()
	IdleMysticDeployContracts(b, t)
}

func TestCreateIdleMysitcHero(t *testing.T) {
	b := newEmulator()

	nftAddr, idleMysitcAddr, idleMysticSigner := IdleMysticDeployContracts(b, t)
	// This scripts returns the number of IdleMysitcHeros currently in existence.
	supply := executeScriptAndCheck(t, b, idleMysticGenerateInspectIdleMysitcHerosSupplyScript(nftAddr.String(), idleMysitcAddr.String()), nil)
	assert.Equal(t, cadence.NewUInt64(0), supply.(cadence.UInt64))
	// This script returns the size of an account's IdleMysitcHeros collection.
	length := executeScriptAndCheck(
		t,
		b,
		idleMysticGenerateInspectCollectionLenScript(nftAddr.String(), idleMysitcAddr.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))},
	)
	assert.Equal(t, cadence.NewInt(0), length.(cadence.Int))

	t.Run("Should be able to mint a idleMystic", func(t *testing.T) {
		IdleMysticMint(t, b, nftAddr, idleMysitcAddr,idleMysitcAddr, idleMysticSigner)

		// Assert that the account's collection is correct
		nftid := executeScriptAndCheck(
			t,
			b,
			idleMysticGenerateInspectCollectionLenScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))},
		)
		assert.Equal(t, cadence.NewInt(1), nftid.(cadence.Int))

	})

}

func TestTransferNFT(t *testing.T) {
	b := newEmulator()

	nftAddr, idleMysitcAddr, idleMysticSigner := IdleMysticDeployContracts(b, t)

	SaleuserAddress, SaleUserSigner, _ := createAccount(t, b)
	acceptAddress, acceptSigner, _ := createAccount(t, b)
	// create a new Collection
	t.Run("Should be able to create a new empty NFT Collection", func(t *testing.T) {
		IdleMysticSetupHerosAccount(t, b, SaleuserAddress, SaleUserSigner, nftAddr, idleMysitcAddr)
		IdleMysticSetupHerosAccount(t, b, acceptAddress, acceptSigner, nftAddr, idleMysitcAddr)

		length := executeScriptAndCheck(
			t,
			b, idleMysticGenerateInspectCollectionLenScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))},
		)
		assert.Equal(t, cadence.NewInt(0), length.(cadence.Int))

	})

	//transfer an NFT
	t.Run("Should be able to withdraw an NFT and deposit to another accounts collection", func(t *testing.T) {
		//mint nft
		IdleMysticMint(t, b, nftAddr, idleMysitcAddr, idleMysitcAddr,idleMysticSigner)
		//transfer NFTID 1 to acceptaddress
		IdleMysticTransfer(b, t, nftAddr, idleMysitcAddr, idleMysticSigner, 1, acceptAddress, false)

		//Assert that the account's collection is [1]
		res := executeScriptAndCheck(t, b, idleMysticGenerateInspectCollectionIdsScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(acceptAddress))})

		assert.Equal(t, "[1]", res.(cadence.Array).String())

	})

}


func TestHerosInherit(t *testing.T) {
	b := newEmulator()

	nftAddr, idleMysitcAddr, idleMysticSigner := IdleMysticDeployContracts(b, t)
	// This scripts returns the number of IdleMysitcHeros currently in existence.
	supply := executeScriptAndCheck(t, b, idleMysticGenerateInspectIdleMysitcHerosSupplyScript(nftAddr.String(), idleMysitcAddr.String()), nil)
	assert.Equal(t, cadence.NewUInt64(0), supply.(cadence.UInt64))
	// This script returns the size of an account's IdleMysitcHeros collection.
	length := executeScriptAndCheck(
		t,
		b,
		idleMysticGenerateInspectCollectionLenScript(nftAddr.String(), idleMysitcAddr.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))},
	)
	assert.Equal(t, cadence.NewInt(0), length.(cadence.Int))

	t.Run("Should be able to mint idleMysticHeros and inherit ", func(t *testing.T) {
		IdleMysticMint(t, b, nftAddr, idleMysitcAddr,idleMysitcAddr, idleMysticSigner)
		IdleMysticMint(t, b, nftAddr, idleMysitcAddr,idleMysitcAddr, idleMysticSigner)

		// Assert that the account's collection is correct
		//check saleUser collection now has get NFT id 1, 2
		res := executeScriptAndCheck(t, b, idleMysticGenerateInspectCollectionIdsScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))})
		assert.Equal(t, "[1, 2]", res.(cadence.Array).String())
		//start inherit
		InheritHeros(b, t, nftAddr, idleMysitcAddr, idleMysticSigner, 1,2, idleMysitcAddr, false)
		//check collection has child NFT id 3
		res = executeScriptAndCheck(t, b, idleMysticGenerateInspectCollectionIdsScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))})
		assert.Equal(t, "[1, 2, 3]", res.(cadence.Array).String())

		//check child info code
		res = executeScriptAndCheck(t, b, HeroinfoScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{
				jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr)),
				jsoncdc.MustEncode(cadence.NewUInt64(1)),
			})
		log.Println("parentA info code",res)
		res = executeScriptAndCheck(t, b, HeroinfoScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{
				jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr)),
				jsoncdc.MustEncode(cadence.NewUInt64(2)),
			})
		log.Println("parentB info code",res)
		res = executeScriptAndCheck(t, b, HeroinfoScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{
				jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr)),
				jsoncdc.MustEncode(cadence.NewUInt64(3)),
			})
		log.Println("child info code  ",res)
		assert.NotEmpty(t,res)

	})

}