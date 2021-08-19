package test
//
//import (
//	nft_contracts "github.com/onflow/flow-nft/lib/go/contracts"
//	"regexp"
//	"testing"
//
//	"github.com/onflow/cadence"
//	jsoncdc "github.com/onflow/cadence/encoding/json"
//	emulator "github.com/onflow/flow-emulator"
//	sdk "github.com/onflow/flow-go-sdk"
//	"github.com/onflow/flow-go-sdk/crypto"
//	sdktemplates "github.com/onflow/flow-go-sdk/templates"
//	"github.com/onflow/flow-go-sdk/test"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//
//	"github.com/onflow/flow-go-sdk"
//)
//
//const (
//	TransactionsRootPath  = "../../transactions"
//	idleMysitcScriptsRootPath = "../../transactions/scripts"
//
//	//	idleMysticRootPath                   = "../../"
//	//	idleMysticIdleMysticPath             = "../../contracts/IdleMysitcHeros.cdc"
//	//	idleMysticSetupAccountPath           = idleMysticRootPath + "transactions/Heros/setup_account.cdc"
//	//	idleMysticMintIdleMysitcHerosPath          = idleMysticRootPath + "/transactions/mint_kitty_item.cdc"
//	//	idleMysticTransferIdleMysitcHerosPath      = idleMysticRootPath + "/transactions/transfer_kittgo y_item.cdc"
//	//	idleMysticInspectIdleMysitcHerosSupplyPath = idleMysticRootPath + "/scripts/read_kitty_items_supply.cdc"
//	//	idleMysticInspectCollectionLenPath   = idleMysticRootPath + "/scripts/read_collection_length.cdc"
//	//	idleMysticInspectCollectionIdsPath   = idleMysticRootPath + "/scripts/read_collection_ids.cdc"
//	//	typeID = 1000
//	LeagueHerosContractPath = "../../contracts/IdleMysitcHeros.cdc"
//	SetupAccountPath        =  "../../transactions/Heros/setup_account.cdc"
//	CreatePlayPath          = TransactionsRootPath + "/admin/create_play.cdc"
//	CreateSetPath           = TransactionsRootPath + "/admin/create_set.cdc"
//	AddPlayPath             = TransactionsRootPath + "/admin/add_play_to_set.cdc"
//	MintPath                = TransactionsRootPath + "/admin/mint_moment.cdc"
//	LockSetPath             = TransactionsRootPath + "/admin/lock_set.cdc"
//	RetirePlayALLPath       = TransactionsRootPath + "/admin/retire_all.cdc"
//	StartNewSchedulePath    = TransactionsRootPath + "/admin/start_new_schedule.cdc"
//	TransferPath            = TransactionsRootPath + "/user/transfer_moment.cdc"
//
//	InspectSupplyPath        = idleMysitcScriptsRootPath + "/get_totalSupply.cdc"
//	InspectCollectionLenPath = idleMysitcScriptsRootPath + "/collections/get_collection_length.cdc"
//	GetAllPlaysLengthPath    = idleMysitcScriptsRootPath + "/plays/get_all_plays_length.cdc"
//	GetAllPlaysPath          = idleMysitcScriptsRootPath + "/plays/get_all_plays.cdc"
//	GetNextPlayIDPath        = idleMysitcScriptsRootPath + "/plays/get_nextPlayID.cdc"
//	GetPlayMetadataPath      = idleMysitcScriptsRootPath + "/plays/get_play_metadata.cdc"
//	GetSetNamePath           = idleMysitcScriptsRootPath + "/sets/get_setName.cdc"
//	GetPlaysInSetPath        = idleMysitcScriptsRootPath + "/sets/get_plays_in_set.cdc"
//	GetCurrentSchedulePath   = idleMysitcScriptsRootPath + "/get_currentSchedule.cdc"
//	GetEditionRetiredPath    = idleMysitcScriptsRootPath + "/sets/get_edition_retired.cdc"
//
//)
//
//func DeployContracts(b *emulator.Blockchain, t *testing.T) (flow.Address, flow.Address, crypto.Signer) {
//	accountKeys := test.AccountKeyGenerator()
//
//	// should be able to deploy a contract as a new account with no keys
//	nftCode := loadNonFungibleToken()
//	nftAddr, err := b.CreateAccount(
//		nil,
//		[]sdktemplates.Contract{
//			{
//				Name:   "NonFungibleToken",
//				Source: string(nftCode),
//			},
//		},
//	)
//	require.NoError(t, err)
//
//	_, err = b.CommitBlock()
//	assert.NoError(t, err)
//
//	// should be able to deploy a contract as a new account with one key
//	// idleMysitcAccountKey, idleMysitcSigner := accountKeys.NewWithSigner()
//	idleMysitcAccountKey, idleMysitcSigner := accountKeys.NewWithSigner()
//	idleMysitcCode := loadLeagueHeros(nftAddr.String())
//	// idleMysitcCode := loadidleMysitc(nftAddr.String())
//	// idleMysitcAddr, err := b.CreateAccount(
//	idleMysitcAddr, err := b.CreateAccount(
//		[]*flow.AccountKey{idleMysitcAccountKey},
//		[]sdktemplates.Contract{
//			{
//				Name:   "IdleMysitcHeros",
//				Source: string(idleMysitcCode),
//			},
//		},
//	)
//	assert.NoError(t, err)
//
//	_, err = b.CommitBlock()
//	assert.NoError(t, err)
//	//t.Log("address",idleMysitcAddr,nftAddr)
//	// simplify the workflow by having the contract address also be our initial test collection
//	idleMysitcSetupAccount(t, b, idleMysitcAddr, idleMysitcSigner, nftAddr, idleMysitcAddr)
//
//	return nftAddr, idleMysitcAddr, idleMysitcSigner
//}
//
//func idleMysitcSetupAccount(
//	t *testing.T, b *emulator.Blockchain,
//	SaleuserAddress sdk.Address, SaleUserSigner crypto.Signer, nftAddr sdk.Address, idleMysitcAddr sdk.Address,
//) {
//	tx := flow.NewTransaction().
//		SetScript(GenerateSetupAccountScript(nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(SaleuserAddress)
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, SaleuserAddress},
//		[]crypto.Signer{b.ServiceKey().Signer(), SaleUserSigner},
//		false,
//	)
//}
//
//func createPlay(
//	t *testing.T, b *emulator.Blockchain,
//	nftAddr, idleMysitcAddr flow.Address,
//	idleMysitcSigner crypto.Signer,
//) {
//	tx := flow.NewTransaction().
//		SetScript(GenerateScript(CreatePlayPath, nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(idleMysitcAddr)
//
//	nameKey := cadence.NewString("Name")
//	nameValue := cadence.NewString("Fank")
//	nameKey2 := cadence.NewString("Name")
//	nameValue2 := cadence.NewString("Lili")
//	nameKey3 := cadence.NewString("Name")
//	nameValue3 := cadence.NewString("virgil")
//	// FankPlayID := uint32(1)
//	metadata := []cadence.KeyValuePair{{Key: nameKey, Value: nameValue}}
//	play := cadence.NewDictionary(metadata)
//
//	_ = tx.AddArgument(play)
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
//		[]crypto.Signer{b.ServiceKey().Signer(), idleMysitcSigner},
//		false,
//	)
//	tx = flow.NewTransaction().
//		SetScript(GenerateScript(CreatePlayPath, nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(idleMysitcAddr)
//	metadata = []cadence.KeyValuePair{{Key: nameKey2, Value: nameValue2}}
//	play = cadence.NewDictionary(metadata)
//
//	_ = tx.AddArgument(play)
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
//		[]crypto.Signer{b.ServiceKey().Signer(), idleMysitcSigner},
//		false,
//	)
//
//	tx = flow.NewTransaction().
//		SetScript(GenerateScript(CreatePlayPath, nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(idleMysitcAddr)
//	metadata = []cadence.KeyValuePair{{Key: nameKey3, Value: nameValue3}}
//	play = cadence.NewDictionary(metadata)
//
//	_ = tx.AddArgument(play)
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
//		[]crypto.Signer{b.ServiceKey().Signer(), idleMysitcSigner},
//		false,
//	)
//
//	length := executeScriptAndCheck(
//		t,
//		b,
//		idleMysitcGenerateScript(GetAllPlaysLengthPath, nftAddr.String(), idleMysitcAddr.String()),
//		nil,
//	)
//	assert.EqualValues(t, cadence.NewInt(3), length)
//
//}
//
//func createMatch(
//	t *testing.T, b *emulator.Blockchain,
//	nftAddr, idleMysitcAddr flow.Address,
//	idleMysitcSigner crypto.Signer,
//) {
//	fiveKillSetID := uint32(1)
//	tx := flow.NewTransaction().
//		SetScript(GenerateScript(CreateSetPath, nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(idleMysitcAddr)
//
//	_ = tx.AddArgument(cadence.NewString("FiveKill"))
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
//		[]crypto.Signer{b.ServiceKey().Signer(), idleMysitcSigner},
//		false,
//	)
//
//	setName := executeScriptAndCheck(
//		t,
//		b,
//		idleMysitcGenerateScript(GetSetNamePath, nftAddr.String(), idleMysitcAddr.String()),
//		[][]byte{jsoncdc.MustEncode(cadence.UInt32(fiveKillSetID))},
//	)
//	assert.EqualValues(t, cadence.NewString("FiveKill"), setName)
//}
//
//func addPlaytoMatch(
//	t *testing.T, b *emulator.Blockchain,
//	nftAddr, idleMysitcAddr flow.Address,
//	idleMysitcSigner crypto.Signer,
//) {
//	tx := flow.NewTransaction().
//		SetScript(GenerateScript(AddPlayPath, nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(idleMysitcAddr)
//
//	_ = tx.AddArgument(cadence.NewUInt32(1))
//	_ = tx.AddArgument(cadence.NewUInt32(1))
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
//		[]crypto.Signer{b.ServiceKey().Signer(), idleMysitcSigner},
//		false,
//	)
//
//	plays := executeScriptAndCheck(
//		t,
//		b,
//		idleMysitcGenerateScript(GetPlaysInSetPath, nftAddr.String(), idleMysitcAddr.String()),
//		[][]byte{jsoncdc.MustEncode(cadence.UInt32(1))},
//	)
//	assert.EqualValues(t, cadence.NewArray([]cadence.Value{cadence.NewUInt32(1)}), plays)
//}
//
//func retirePlay(
//	t *testing.T, b *emulator.Blockchain,
//	nftAddr, idleMysitcAddr flow.Address,
//	idleMysitcSigner crypto.Signer,
//) {
//
//	// check play in match
//	play := executeScriptAndCheck(
//		t,
//		b,
//		idleMysitcGenerateScript(GetPlaysInSetPath, nftAddr.String(), idleMysitcAddr.String()),
//		[][]byte{jsoncdc.MustEncode(cadence.UInt32(1))},
//	)
//	assert.EqualValues(t, cadence.NewArray([]cadence.Value{cadence.NewUInt32(1)}), play)
//
//	tx := flow.NewTransaction().
//		SetScript(GenerateScript(RetirePlayALLPath, nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(idleMysitcAddr)
//
//	_ = tx.AddArgument(cadence.NewUInt32(1))
//	// _ = tx.AddArgument(cadence.NewUInt32(1))
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
//		[]crypto.Signer{b.ServiceKey().Signer(), idleMysitcSigner},
//		false,
//	)
//
//	isRetired := executeScriptAndCheck(
//		t,
//		b,
//		idleMysitcGenerateScript(GetEditionRetiredPath, nftAddr.String(), idleMysitcAddr.String()),
//		[][]byte{jsoncdc.MustEncode(cadence.UInt32(1)), jsoncdc.MustEncode(cadence.UInt32(1))},
//	)
//	assert.EqualValues(t, cadence.Bool(true), isRetired)
//}
//
//func idleMysitcMintItem(
//	t *testing.T, b *emulator.Blockchain,
//	nftAddr, idleMysitcAddr flow.Address,
//	idleMysitcSigner crypto.Signer,
//) {
//	tx := flow.NewTransaction().
//		SetScript(GenerateScript(MintPath, nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(idleMysitcAddr)
//
//	_ = tx.AddArgument(cadence.NewUInt32(1))
//	_ = tx.AddArgument(cadence.NewUInt32(1))
//	_ = tx.AddArgument(cadence.NewAddress(idleMysitcAddr))
//	_ = tx.AddArgument(cadence.NewString("https://this_is_a_pic.jpg"))
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
//		[]crypto.Signer{b.ServiceKey().Signer(), idleMysitcSigner},
//		false,
//	)
//
//	total := executeScriptAndCheck(
//		t,
//		b,
//		idleMysitcGenerateScript(InspectSupplyPath, nftAddr.String(), idleMysitcAddr.String()),
//		nil,
//	)
//	assert.EqualValues(t, cadence.NewUInt64(1), total)
//}
//
//func idleMysitcTransferItem(
//	t *testing.T, b *emulator.Blockchain,
//	nftAddr, idleMysitcAddr flow.Address, idleMysitcSigner crypto.Signer,
//	typeID uint64, recipientAddr flow.Address, shouldFail bool,
//) {
//
//	tx := flow.NewTransaction().
//		SetScript(idleMysitcGenerateTransferidleMysitccript(nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(idleMysitcAddr)
//
//	_ = tx.AddArgument(cadence.NewAddress(recipientAddr))
//	_ = tx.AddArgument(cadence.NewUInt64(typeID))
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
//		[]crypto.Signer{b.ServiceKey().Signer(), idleMysitcSigner},
//		shouldFail,
//	)
//}
//
//func newSchedule(
//	t *testing.T, b *emulator.Blockchain,
//	nftAddr, idleMysitcAddr flow.Address, idleMysitcSigner crypto.Signer,
//) {
//
//	tx := flow.NewTransaction().
//		SetScript(GenerateScript(StartNewSchedulePath, nftAddr.String(), idleMysitcAddr.String())).
//		SetGasLimit(100).
//		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
//		SetPayer(b.ServiceKey().Address).
//		AddAuthorizer(idleMysitcAddr)
//
//	signAndSubmit(
//		t, b, tx,
//		[]flow.Address{b.ServiceKey().Address, idleMysitcAddr},
//		[]crypto.Signer{b.ServiceKey().Signer(), idleMysitcSigner},
//		false,
//	)
//
//	schedule := executeScriptAndCheck(
//		t,
//		b,
//		idleMysitcGenerateScript(GetCurrentSchedulePath, nftAddr.String(), idleMysitcAddr.String()),
//		nil,
//	)
//	assert.EqualValues(t, cadence.NewUInt32(1), schedule)
//}
//
//
////func TestDeployContracts(t *testing.T) {
////	b := newEmulator()
////	DeployContracts(b, t)
////}
//
//func TestCreateLeague(t *testing.T) {
//	b := newEmulator()
//
//	nftAddr, idleMysitcAddr, idleMysitcSigner := DeployContracts(b, t)
//
//	supply := executeScriptAndCheck(
//		t, b,
//		idleMysitcGenerateInspectidleMysitcsupplyScript(nftAddr.String(), idleMysitcAddr.String()),
//		nil,
//	)
//	assert.EqualValues(t, cadence.NewUInt64(0), supply)
//
//	// assert that the account collection is empty
//	length := executeScriptAndCheck(
//		t,
//		b,
//		idleMysitcGenerateInspectCollectionLenScript(nftAddr.String(), idleMysitcAddr.String()),
//		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(idleMysitcAddr))},
//	)
//	assert.EqualValues(t, cadence.NewInt(0), length)
//
//	t.Run("Should be able to add a play", func(t *testing.T) {
//		createPlay(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
//	})
//
//	t.Run("Should be able to add a match", func(t *testing.T) {
//		createMatch(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
//	})
//
//	t.Run("Should be able to add play to match", func(t *testing.T) {
//		addPlaytoMatch(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
//	})
//
//	t.Run("Should be able to mint item", func(t *testing.T) {
//		idleMysitcMintItem(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
//	})
//}
//
////func TestTransferNFT(t *testing.T) {
////	b := newEmulator()
////
////	nftAddr, idleMysitcAddr, idleMysitcSigner := DeployContracts(b, t)
////
////	SaleuserAddress, SaleUserSigner, _ := createAccount(t, b)
////
////	// create a new Collection for new account
////	t.Run("Should be able to create a new empty NFT Collection", func(t *testing.T) {
////		idleMysitcSetupAccount(t, b, SaleuserAddress, SaleUserSigner, nftAddr, idleMysitcAddr)
////
////		length := executeScriptAndCheck(
////			t,
////			b, idleMysitcGenerateInspectCollectionLenScript(nftAddr.String(), idleMysitcAddr.String()),
////			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))},
////		)
////		assert.EqualValues(t, cadence.NewInt(0), length)
////	})
////
////	t.Run("Shouldn not be able to withdraw an NFT that does not exist in a collection", func(t *testing.T) {
////		nonExistentID := uint64(3333333)
////
////		idleMysitcTransferItem(
////			t, b,
////			nftAddr, idleMysitcAddr, idleMysitcSigner,
////			nonExistentID, SaleuserAddress, true,
////		)
////	})
////
////	// transfer an NFT
////	t.Run("Should be able to withdraw an NFT and deposit to another accounts collection", func(t *testing.T) {
////		createPlay(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
////		createMatch(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
////		addPlaytoMatch(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
////		idleMysitcMintItem(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
////		// Cheat: we have minted one item, its ID will be zero
////		idleMysitcTransferItem(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner, 1, SaleuserAddress, false)
////	})
////
////	length := executeScriptAndCheck(
////		t,
////		b, idleMysitcGenerateInspectCollectionLenScript(nftAddr.String(), idleMysitcAddr.String()),
////		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(SaleuserAddress))},
////	)
////	assert.EqualValues(t, cadence.NewInt(1), length)
////
////}
//
//func TestRetire(t *testing.T) {
//	b := newEmulator()
//
//	nftAddr, idleMysitcAddr, idleMysitcSigner := DeployContracts(b, t)
//
//	t.Run("Should be able to add a play", func(t *testing.T) {
//		createPlay(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
//		createMatch(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
//		addPlaytoMatch(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
//		retirePlay(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
//	})
//}
//
//func TestNewSchedule(t *testing.T) {
//	b := newEmulator()
//
//	nftAddr, idleMysitcAddr, idleMysitcSigner := DeployContracts(b, t)
//
//	t.Run("Should be able to add a play", func(t *testing.T) {
//		newSchedule(t, b, nftAddr, idleMysitcAddr, idleMysitcSigner)
//	})
//}
//
//func replaceidleMysitcAddressPlaceholders(code, nftAddress, idleMysitcAddress string) []byte {
//	return []byte(replaceImports(
//		code,
//		map[string]*regexp.Regexp{
//			nftAddress:    nftAddressPlaceholder,
//			idleMysitcAddress: idleMysitcAddressPlaceHolder,
//		},
//	))
//}
////
//func loadNonFungibleToken() []byte {
//	return nft_contracts.NonFungibleToken()
//}
//
//func loadLeagueHeros(nftAddr string) []byte {
//	return []byte(replaceImports(
//		string(readFile(LeagueHerosContractPath)),
//		map[string]*regexp.Regexp{
//			nftAddr: nftAddressPlaceholder,
//		},
//	))
//	// return replaceidleMysitcAddressPlaceholders(
//	// 	string(readFile(LeagueHerosContractPath)),
//	// 	nftAddr,
//	// 	"import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"",
//	// )
//}
//
//func GenerateSetupAccountScript(nftAddr, idleMysitcAddr string) []byte {
//	return replaceidleMysitcAddressPlaceholders(
//		string(readFile(SetupAccountPath)),
//		nftAddr,
//		idleMysitcAddr,
//	)
//}
//
//func GenerateScript(Path, nftAddr, idleMysitcAddr string) []byte {
//	return replaceidleMysitcAddressPlaceholders(
//		string(readFile(Path)),
//		nftAddr,
//		idleMysitcAddr,
//	)
//}
//
//func GenerateMintidleMysitcscript(nftAddr, idleMysitcAddr string) []byte {
//	return replaceidleMysitcAddressPlaceholders(
//		string(readFile(MintPath)),
//		nftAddr,
//		idleMysitcAddr,
//	)
//}
//
//func idleMysitcGenerateTransferidleMysitccript(nftAddr, idleMysitcAddr string) []byte {
//	return replaceidleMysitcAddressPlaceholders(
//		string(readFile(TransferPath)),
//		nftAddr,
//		idleMysitcAddr,
//	)
//}
//
//func idleMysitcGenerateInspectidleMysitcsupplyScript(nftAddr, idleMysitcAddr string) []byte {
//	return replaceidleMysitcAddressPlaceholders(
//		string(readFile(InspectSupplyPath)),
//		nftAddr,
//		idleMysitcAddr,
//	)
//}
//
//func idleMysitcGenerateInspectCollectionLenScript(nftAddr, idleMysitcAddr string) []byte {
//	return replaceidleMysitcAddressPlaceholders(
//		string(readFile(InspectCollectionLenPath)),
//		nftAddr,
//		idleMysitcAddr,
//	)
//}
//
//func idleMysitcGenerateScript(path, nftAddr, idleMysitcAddr string) []byte {
//	return replaceidleMysitcAddressPlaceholders(
//		string(readFile(path)),
//		nftAddr,
//		idleMysitcAddr,
//	)
//}
