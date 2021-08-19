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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

const (
	FungibleTokenPath                   = "../../contracts/FungibleToken.cdc"
	idleMysticTokenMysticPath           = "../../contracts/IdleMysitcToken.cdc"
	idleMysticTokenMintPath             = "../../transactions/Stone/mint_tokens.cdc"
	idleMysticTokenSetupAccountMintPath = "../../transactions/Stone/setup_account.cdc"
	idleMysticTransferTokenPath         = "../../transactions/Stone/transfer_tokens.cdc"

	idleMysticTokenGetBalancePath       = "../../scripts/Stone/get_balance.cdc"
)

func IdleMysticTokenDeployContracts(b *emulator.Blockchain, t *testing.T) (flow.Address, flow.Address, crypto.Signer) {
	accountKeys := test.AccountKeyGenerator()

	// should be able to deploy a contract as a new account with no keys
	ftCode := loadFungibleToken()
	ftAddr, err := b.CreateAccount(
		nil,
		[]sdktemplates.Contract{
			{
				Name:   "FungibleToken",
				Source: string(ftCode),
			},
		},
	)
	require.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// should be able to deploy a contract as a new account with one key

	idleMysitcAccountKey, idleMysitcSigner := accountKeys.NewWithSigner()
	idleMysitcCode := loadIdleMysitcHerosToken(ftAddr.String())

	idleMysitctokenAddr, err := b.CreateAccount(
		[]*flow.AccountKey{idleMysitcAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "IdleMysitcToken",
				Source: string(idleMysitcCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// simplify the workflow by having the contract address also be our initial test collection
	idleMysitcTokenSetupAccount(t, b, idleMysitctokenAddr, idleMysitcSigner, ftAddr, idleMysitctokenAddr)

	return ftAddr, idleMysitctokenAddr, idleMysitcSigner
}
func GenerateTokenSetupAccountScript(ftAddr, idleMysitcAddr string) []byte {
	return replaceIdleMysticTokenAddressPlaceholders(
		string(readFile(idleMysticTokenSetupAccountMintPath)),
		ftAddr,
		idleMysitcAddr,
	)
}

func idleMysitcTokenSetupAccount(
	t *testing.T, b *emulator.Blockchain,
	SaleuserAddress sdk.Address, SaleUserSigner crypto.Signer,
	nftAddr sdk.Address, idleMysitcAddr sdk.Address, ) {
	tx := flow.NewTransaction().
		SetScript(GenerateTokenSetupAccountScript(nftAddr.String(), idleMysitcAddr.String())).
		SetGasLimit(9999).
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


func IdleMysticTokenMint(t *testing.T, b *emulator.Blockchain, nftAddr, idleMysitcAddr,receptAddr flow.Address, idleMysticSigner crypto.Signer, ) {
	tx := flow.NewTransaction().
		SetScript(MintIdleMysticHerosTokenScript(nftAddr.String(), idleMysitcAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(idleMysitcAddr)

	tx.AddArgument(cadence.NewAddress(receptAddr))
	amount, _ := cadence.NewUFix64("10.0")
	tx.AddArgument(amount)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), idleMysticSigner},
		false,
	)
}

func IdleMysticTokenTransfer(b *emulator.Blockchain, t *testing.T, nftAddr, idleMysitcAddr flow.Address, idleMysticSigner crypto.Signer, recipientAddr flow.Address, shouldFail bool,amounts string) {
	tx := flow.NewTransaction().
		SetScript(idleMysticGenerateTransfertokenScript(nftAddr.String(), idleMysitcAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(idleMysitcAddr)
	tx.AddArgument(cadence.NewAddress(recipientAddr))
	tx.AddArgument(CadenceUFix64("10.0"))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), idleMysticSigner},
		shouldFail,
	)
}
func loadFungibleToken() []byte {
	return readFile(FungibleTokenPath)
}

func loadIdleMysitcHerosToken(ftAddr string) []byte {

	return []byte(replaceImports(
		string(readFile(idleMysticTokenMysticPath)),
		map[string]*regexp.Regexp{
			ftAddr: ftAddressPlaceholder,
		},
	))
}

func replaceIdleMysticTokenAddressPlaceholders(code, ftAddress, idleMysitcAddress string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			ftAddress:         ftAddressPlaceholder,
			idleMysitcAddress: idleMysticTokenAddressPlaceHolder,
		},
	))
}

func idleMystictokenGenerateGetBalanceScript(nftAddr, idleMysitcAddr string) []byte {
	// This script returns an account's IdleMysitcToken balance.

	return replaceIdleMysticTokenAddressPlaceholders(
		string(readFile(idleMysticTokenGetBalancePath)),
		nftAddr,
		idleMysitcAddr,
	)
}
func idleMysticGenerateTransfertokenScript(ftAddr, idleMysitcAddr string) []byte {
	return replaceIdleMysticTokenAddressPlaceholders(
		string(readFile(idleMysticTransferTokenPath)),
		ftAddr,
		idleMysitcAddr,
	)
}

func MintIdleMysticHerosTokenScript(nftAddr, idleMysitcAddr string) []byte {
	return replaceIdleMysticTokenAddressPlaceholders(
		string(readFile(idleMysticTokenMintPath)),
		nftAddr,
		idleMysitcAddr,
	)
}

func TestIdleMysticTokenDeployContracts(t *testing.T) {
	b := newEmulator()
	IdleMysticTokenDeployContracts(b, t)
}

func TestCreateIdleMysitcHerosToken(t *testing.T) {
	b := newEmulator()

	ftAddr, tokenAddr, tokenSigner := IdleMysticTokenDeployContracts(b, t)
	// This script returns the size of an account's IdleMysitcHeros collection.
	amount := executeScriptAndCheck(
		t,
		b,
		idleMystictokenGenerateGetBalanceScript(ftAddr.String(), tokenAddr.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(tokenAddr))},
	)
	//Initial amount
	assert.Equal(t, amount.String(), "0.00000000")

	//assert.Equal(t, cadence.NewInt(0), length.(cadence.Int))

	t.Run("Should be able to mint a idleMystictoken", func(t *testing.T) {
		//mint token
		IdleMysticTokenMint(t, b, ftAddr, tokenAddr, tokenAddr,tokenSigner)

		// Assert that the account's amount is correct

		amount := executeScriptAndCheck(
			t,
			b,
			idleMystictokenGenerateGetBalanceScript(ftAddr.String(), tokenAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(tokenAddr))},
		)
		//now amount
		assert.Equal(t, amount.String(), "10.00000000")
	})

}

func TestTransferToken(t *testing.T) {
	b := newEmulator()

	ftAddr, idleMysitcAddr, idleMysticSigner := IdleMysticTokenDeployContracts(b, t)
	// This script returns the size of an account's IdleMysitcHeros collection.
	amount := executeScriptAndCheck(
		t,
		b,
		idleMystictokenGenerateGetBalanceScript(ftAddr.String(), idleMysitcAddr.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))},
	)
	//Initial amount
	assert.Equal(t, amount.String(), "0.00000000")

	SaleuserAddress, SaleUserSigner, _ := createAccount(t, b)
	acceptAddress, acceptSigner, _ := createAccount(t, b)
	//create a new Collection
	t.Run("Should be able to create a new empty  Collection", func(t *testing.T) {
		idleMysitcTokenSetupAccount(t, b, SaleuserAddress, SaleUserSigner, ftAddr, idleMysitcAddr)
		idleMysitcTokenSetupAccount(t, b, acceptAddress, acceptSigner, ftAddr, idleMysitcAddr)


		IdleMysticTokenMint(t, b, ftAddr, idleMysitcAddr,idleMysitcAddr, idleMysticSigner)

		amount := executeScriptAndCheck(
			t,
			b,
			idleMystictokenGenerateGetBalanceScript(ftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))},
		)
		//now amount

		assert.Equal(t, amount.String(),"10.00000000")
		//
		//
	})

	//transfer an Token Stone
	t.Run("Should be able  transfer token to another accounts amount ", func(t *testing.T) {

		amount := executeScriptAndCheck(
			t,
			b,
			idleMystictokenGenerateGetBalanceScript(ftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))},
		)
		assert.Equal(t,"10.00000000" ,amount.String(),)
		//transfer an Token (Stone)
		IdleMysticTokenTransfer(b, t, ftAddr, idleMysitcAddr, idleMysticSigner,  acceptAddress, false,"10.0")

		amount = executeScriptAndCheck(
			t,
			b,
			idleMystictokenGenerateGetBalanceScript(ftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))},
		)
		//now amount
		assert.Equal(t,"0.00000000" ,amount.String(),)
	})

}

