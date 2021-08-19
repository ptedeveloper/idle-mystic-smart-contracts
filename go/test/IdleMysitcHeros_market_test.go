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
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
	"time"
)

const (
	kittyItemsMarketIdleMysitcHerossMarketPath = "/contracts/IdleMysitcHerossMarket.cdc"
	kittyItemsMarketSetupAccountPath           = "/transactions/setup_account.cdc"
	kittyItemsMarketSellItemPath               = "/transactions/sell_market_item.cdc"
	kittyItemsMarketBuyItemPath                = "/transactions/buy_market_item.cdc"
	kittyItemsMarketRemoveItemPath             = "/transactions/remove_market_item.cdc"
	MarketSetupAccountPath                     = "../../transactions/HerosMarket/setup_account.cdc"
	MarketSalePath                             = "../../transactions/HerosMarket/sell_market_item.cdc"
	MarketBuyPath                             = "../../transactions/HerosMarket/buy_market_item.cdc"
	MarketRemovePath                             = "../../transactions/HerosMarket/remove_market_item.cdc"
	MarketListPath                             = "../../scripts/HerosMarket/read_collection_ids.cdc"
	MarkettPath                                = "../../contracts/IdleMysitcHerosMarket.cdc"
)

func MarketDeployContracts(b *emulator.Blockchain, t *testing.T) (flow.Address,flow.Address,flow.Address,flow.Address,flow.Address,
	crypto.Signer,crypto.Signer,crypto.Signer) {
	accountKeys := test.AccountKeyGenerator()

	// should be able to deploy a contract as a new account NonFungibleToken
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

	// should be able to deploy a contract as a new account IdleMysitcHeros

	idleMysitcAccountKey, idleMysitcHerosSign := accountKeys.NewWithSigner()
	idleMysitcCode := loadIdleMysitcHeros(nftAddr.String())
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

	// should be able to deploy a contract as a new account FungibleToken
	ftAccountKey, _ := accountKeys.NewWithSigner()
	ftCode := loadFungibleToken()
	ftAddr, err := b.CreateAccount(
		[]*flow.AccountKey{ftAccountKey},
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

	// should be able to deploy a contract as a new account IdleMysitcToken
	tokenAccountKey, tokenSign := accountKeys.NewWithSigner()
	idleMysitctokenCode := loadIdleMysitcHerosToken(ftAddr.String())
	idleMysitctokenAddr, err := b.CreateAccount(
		[]*flow.AccountKey{tokenAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "IdleMysitcToken",
				Source: string(idleMysitctokenCode),
			},
		},
	)
	assert.NoError(t, err)
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// should be able to deploy a contract as a new account IdleMysitcHerosMarket
	MarketCode := replaceMarketPlace(string(readFile(MarkettPath)), nftAddr.String(), ftAddr.String(), idleMysitcAddr.String(), idleMysitctokenAddr.String(), "")
	MarketAccountKey, MarketSigner := accountKeys.NewWithSigner()
	MarketAddr, err := b.CreateAccount(
		[]*flow.AccountKey{MarketAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "IdleMysitcHerosMarket",
				Source: string(MarketCode),
			},
		},
	)
	assert.NoError(t, err)
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// simplify the workflow by having the contract address also be our initial test collection
	MarketSetupAccount(t, b, MarketAddr, MarketSigner, MarketAddr)

	return nftAddr,ftAddr, idleMysitcAddr,idleMysitctokenAddr, MarketAddr,MarketSigner,idleMysitcHerosSign,tokenSign
}

func TestMarketDeployContracts(t *testing.T) {
	b := newEmulator()
	MarketDeployContracts(b, t)
}

func MarketSetupAccount(
	t *testing.T, b *emulator.Blockchain,
	SaleuserAddress sdk.Address, SaleUserSigner crypto.Signer, marketAddr sdk.Address, ) {

	tx := flow.NewTransaction().
		SetScript(GenerateMarketSetupAccountScript(marketAddr.String())).
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
func GenerateMarketSetupAccountScript(idleMysitcAddr string) []byte {
	return replaceMarketPlaceholders(
		string(readFile(MarketSetupAccountPath)),
		idleMysitcAddr,
	)
}
func replaceMarketPlaceholders(code, idleMysitcAddress string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			idleMysitcAddress: idleMysticMarketAddressPlaceHolder,
		},
	))
}
func ReadFileAndReplaceAddre(filepath, Adminaddress string, ) (script []byte, err error) {
	Adminaddress = "0x" + Adminaddress
	script, readerr := ioutil.ReadFile(filepath)
	str_script := string(script)
	res := strings.Replace(str_script, "0xADDRESS", Adminaddress, -1)
	res = strings.Replace(str_script, "0xADDRESS", Adminaddress, -1)
	res = strings.Replace(res, "0xNFTADDRESS", Adminaddress, -1)
	res = strings.Replace(res, "0xSHARDEDADDRESS", Adminaddress, -1)
	res = strings.Replace(res, "0xFUNGIBLETOKENADDRESS", Adminaddress, -1)
	res = strings.Replace(res, "0xMARKETADDRESS", Adminaddress, -1)
	res = strings.Replace(res, "0xNonFungibleToken", Adminaddress, -1)
	res = strings.Replace(res, "\"../../contracts/NonFungibleToken.cdc\"", Adminaddress, -1)
	res = strings.Replace(res, "\"../../contracts/IdleMysitcHeros.cdc\"", Adminaddress, -1)
	res = strings.Replace(res, "\"../../contracts/IdleMysitcHerosMarket.cdc\"", Adminaddress, -1)
	script = []byte(res)
	return script, readerr
}

func replaceMarketPlace(code, nftAddress, ftAddress, idleMysitcAddress, tokenAddress, marketAddress string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			nftAddress:        nftAddressPlaceholder,
			ftAddress:         ftAddressPlaceholder,
			idleMysitcAddress: idleMysitcAddressPlaceHolder,
			tokenAddress:      idleMysticTokenAddressPlaceHolder,
			marketAddress:     idleMysticMarketAddressPlaceHolder,
		},
	))

}

//// Create a new account with the Kibble and IdleMysitcHeross resources set up BUT no IdleMysitcHerossMarket resource.
func MarketCreatePurchaserAccount(b *emulator.Blockchain, t *testing.T, nftAddr, ftAddr, idleMysitcAddr, idleMysitctokenAddr sdk.Address, idleMysitcSigner crypto.Signer, ) (sdk.Address, crypto.Signer) {
	SaleuserAddress, SaleUserSigner, _ := createAccount(t, b)
	idleMysitcSetupAccount(t, b, idleMysitcAddr, idleMysitcSigner, nftAddr, idleMysitcAddr)
	idleMysitcTokenSetupAccount(t, b, idleMysitctokenAddr, idleMysitcSigner, ftAddr, idleMysitctokenAddr)
	MarketSetupAccount(t, b, idleMysitcAddr, idleMysitcSigner, nftAddr)

	return SaleuserAddress, SaleUserSigner
}
//
//func IdleMysitcHerossMarketListItem(b *emulator.Blockchain, t *testing.T, SaleuserAddress sdk.Address, SaleUserSigner crypto.Signer, shouldFail bool) {
//	tx := flow.NewTransaction().
//		SetScript(MarketGenerateSellItemScript(SaleuserAddress.String())).
//		SetGasLimit(9999).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(SaleuserAddress)
//	tx.AddArgument(cadence.NewAddress(SaleuserAddress))
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, SaleuserAddress},
//		[]crypto.Signer{b.ServiceKey().Signer(), SaleUserSigner},
//		shouldFail,
//	)
//}

func MarketNFTMint(t *testing.T, b *emulator.Blockchain, idleMysitcAddr flow.Address, idleMysticSigner crypto.Signer, ) {
	tx := flow.NewTransaction().
		SetScript(MintMarketNFTScript(idleMysitcAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(idleMysitcAddr)
	//recipientAddr := cadence.NewAddress(adminAddr)
	_birth := time.Now().Local().String()

	tx.AddArgument(cadence.NewAddress(idleMysitcAddr))
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
func MintMarketNFTScript(nftAddr string) []byte {
	script, _ := ReadFileAndReplaceAddre(idleMysticMintPath, nftAddr)
	return script
}

func TestMarket(t *testing.T) {
	b := newEmulator()

	nftAddr,ftAddr,idleMysitcAddr,tokenAddr, marketaddr, _ ,idleMysticSigner,tokenSigner:= MarketDeployContracts(b, t)

	t.Run("Should be able to Whole market process done  ", func(t *testing.T) {

		//nftAddr, idleMysitcAddr, idleMysticSigner := IdleMysticDeployContracts(b, t)
		//ftAddr, tokenAddr, tokenSigner := IdleMysticTokenDeployContracts(b, t)

		SaleuserAddress, SaleuserSigner, _ := createAccount(t, b)
		BuyuserAddress, BuyUserSigner, _ := createAccount(t, b)
		MarketSetupAccount(t, b, SaleuserAddress, SaleuserSigner, marketaddr)
		MarketSetupAccount(t, b, BuyuserAddress, BuyUserSigner, marketaddr)

		// NFT collection set up

		IdleMysticSetupHerosAccount(t, b, SaleuserAddress, SaleuserSigner, nftAddr, idleMysitcAddr)
		IdleMysticSetupHerosAccount(t, b, BuyuserAddress, BuyUserSigner, nftAddr, idleMysitcAddr)
		// mint nft
		IdleMysticMint(t, b, nftAddr, idleMysitcAddr, SaleuserAddress, idleMysticSigner)
		res := executeScriptAndCheck(t, b, idleMysticGenerateInspectCollectionIdsScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))})

		assert.Equal(t, "[1]", res.(cadence.Array).String())

		// FT (stone) collection set up

		idleMysitcTokenSetupAccount(t, b, SaleuserAddress, SaleuserSigner, ftAddr, tokenAddr)
		idleMysitcTokenSetupAccount(t, b, BuyuserAddress, BuyUserSigner, ftAddr, tokenAddr)
		// mint ft (stone)
		IdleMysticTokenMint(t, b, ftAddr, tokenAddr, BuyuserAddress, tokenSigner)
		amount := executeScriptAndCheck(
			t,
			b,
			idleMystictokenGenerateGetBalanceScript(ftAddr.String(), tokenAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(BuyuserAddress))},
		)
		//now amount
		assert.Equal(t, amount.String(), "10.00000000")

		// sale NFT
		MarketSaleNFT(b, t, nftAddr, ftAddr, idleMysitcAddr, tokenAddr, marketaddr, SaleuserAddress, SaleuserSigner, 1, false)


		//check Market list
		//IdleMysitcHerossMarketListItem(b,t,SaleuserAddress,SaleuserSigner,false)
		res = executeScriptAndCheck(t, b, MarketGenerateSellItemScript(marketaddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))})
		assert.Equal(t, "[1]", res.(cadence.Array).String())


		//purchase NFT id 1
		PurchaseNFT(b, t, nftAddr, ftAddr, idleMysitcAddr, tokenAddr, marketaddr, SaleuserAddress,BuyuserAddress, BuyUserSigner, 1, false)


		//check Market list is empty
		res = executeScriptAndCheck(t, b, MarketGenerateSellItemScript(marketaddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))})
		assert.Equal(t, "[]", res.(cadence.Array).String())


		//check buyuser collection NFT list
		res = executeScriptAndCheck(t, b, idleMysticGenerateInspectCollectionIdsScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(BuyuserAddress))})
		assert.Equal(t, "[1]", res.(cadence.Array).String())

		//remove MarketRemoveNFT
		// mint NFT id 2 to saleuser
		IdleMysticMint(t, b, nftAddr, idleMysitcAddr, SaleuserAddress, idleMysticSigner)
		res = executeScriptAndCheck(t, b, idleMysticGenerateInspectCollectionIdsScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))})
		assert.Equal(t, "[2]", res.(cadence.Array).String())
		//sale NFT id 2
		MarketSaleNFT(b, t, nftAddr, ftAddr, idleMysitcAddr, tokenAddr, marketaddr, SaleuserAddress, SaleuserSigner, 2, false)
		res = executeScriptAndCheck(t, b, MarketGenerateSellItemScript(marketaddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))})
		assert.Equal(t, "[2]", res.(cadence.Array).String())

		// now remove NFT id 2
		MarketRemoveNFT(b, t, nftAddr, ftAddr, idleMysitcAddr, tokenAddr, marketaddr, SaleuserAddress, SaleuserSigner, 2, false)
		// check market is empty
		res = executeScriptAndCheck(t, b, MarketGenerateSellItemScript(marketaddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))})
		assert.Equal(t, "[]", res.(cadence.Array).String())
		//check saleUser collection now has get NFT id 2
		res = executeScriptAndCheck(t, b, idleMysticGenerateInspectCollectionIdsScript(nftAddr.String(), idleMysitcAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))})
		assert.Equal(t, "[2]", res.(cadence.Array).String())

	})

}



func MarketGenerateSellItemScript(address string) []byte {
	script, _ := ReadFileAndReplaceAddre(MarketListPath, address)
	return script

}
func  PurchaseNFT(b *emulator.Blockchain, t *testing.T,
	nftAddr, ftAddr, idleMysitcAddr, idleMysitctokenAddr, marketAddr, saleuserAddr,buyuserAddr flow.Address,
	buyuserSign crypto.Signer,
	NFTID uint64,
	shouldFail bool) {

	tx := flow.NewTransaction().
		SetScript(
			replaceMarketPlace(string(readFile(MarketBuyPath)),
				nftAddr.String(),
				ftAddr.String(),
				idleMysitcAddr.String(),
				idleMysitctokenAddr.String(),
				marketAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(buyuserAddr)
	tx.AddArgument(cadence.NewUInt64(NFTID))
	tx.AddArgument(cadence.NewAddress(saleuserAddr))
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, buyuserAddr,},
		[]crypto.Signer{b.ServiceKey().Signer(), buyuserSign,},
		shouldFail,
	)
}

func  MarketRemoveNFT(b *emulator.Blockchain, t *testing.T,
	nftAddr, ftAddr, idleMysitcAddr, idleMysitctokenAddr, marketAddr, saleuserAddr flow.Address,
	saleuserSign crypto.Signer,
	NFTID uint64,
	shouldFail bool) {

	tx := flow.NewTransaction().
		SetScript(
			replaceMarketPlace(string(readFile(MarketRemovePath)),
				nftAddr.String(),
				ftAddr.String(),
				idleMysitcAddr.String(),
				idleMysitctokenAddr.String(),
				marketAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(saleuserAddr)
	tx.AddArgument(cadence.NewUInt64(NFTID))
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, saleuserAddr,},
		[]crypto.Signer{b.ServiceKey().Signer(), saleuserSign,},
		shouldFail,
	)
}


func GenerateMarketNFTListScript(address string) []byte {
	script, _ := ReadFileAndReplaceAddre(MarketListPath, address)
	return script
}

func MarketSaleNFT(b *emulator.Blockchain, t *testing.T,
	nftAddr, ftAddr, idleMysitcAddr, idleMysitctokenAddr, marketAddr, saleuserAddr flow.Address,
	saleuserSigner crypto.Signer,
	//saleuserKey *sdk.AccountKey,
	NFTID uint64,
	shouldFail bool) {

	tx := flow.NewTransaction().
		SetScript(
			replaceMarketPlace(string(readFile(MarketSalePath)),
				nftAddr.String(),
				ftAddr.String(),
				idleMysitcAddr.String(),
				idleMysitctokenAddr.String(),
				marketAddr.String())).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(saleuserAddr)
	tx.AddArgument(cadence.NewUInt64(NFTID))
	tx.AddArgument(CadenceUFix64("10.0"))
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address,saleuserAddr, },
		[]crypto.Signer{b.ServiceKey().Signer(),saleuserSigner, },
		shouldFail,
	)

}
