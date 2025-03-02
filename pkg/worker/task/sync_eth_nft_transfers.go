package task

import (
	"encoding/json"
	"fmt"

	apimodel "github.com/authgear/authgear-nft-indexer/pkg/api/model"
	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/authgear/authgear-nft-indexer/pkg/util/hexstring"
	"github.com/jrallison/go-workers"
	"github.com/uptrace/bun/extra/bunbig"
)

const TransferPageSize = 1000

type SycnNFTTransferTransferMutator interface {
	InsertNFTTransfers(transfers []ethmodel.NFTTransfer) error
}

type SycnNFTTransferAlchemyAPI interface {
	GetNFTTransfers(blockchainNetwork model.BlockchainNetwork, contractAddresses []string, syncedBlock string, blockNum string, pageKey string, maxCount int64) (*apimodel.AssetTransferResponse, error)
}

type SyncETHNFTTransfersMessageArgs struct {
	Blockchain        string   `json:"blockchain"`
	Network           string   `json:"network"`
	ContractAddresses []string `json:"contract_address"`
	SyncedBlock       string   `json:"synced_block"`
	PageKey           string   `json:"page_key"`
}

type SyncETHNFTTransferTaskHandler struct {
	AlchemyAPI         SycnNFTTransferAlchemyAPI
	Config             config.Config
	NftTransferMutator SycnNFTTransferTransferMutator
}

func (h *SyncETHNFTTransferTaskHandler) Handler(message *workers.Msg) {
	args := message.Args()
	if args == nil {
		panic("SyncNFTTranfers: missing args")
	}

	argsJSON, err := args.Json.Encode()
	if err != nil {
		panic(fmt.Sprintf("SyncNFTTranfers: failed to serialize args: %s", err))
	}

	var castedArgs SyncETHNFTTransfersMessageArgs
	err = json.Unmarshal(argsJSON, &castedArgs)
	if err != nil {
		panic(fmt.Sprintf("SyncNFTTranfers: failed to unmarshal args: %s", err))
	}

	blockchainNetwork := model.BlockchainNetwork{
		Blockchain: castedArgs.Blockchain,
		Network:    castedArgs.Network,
	}

	res, err := h.AlchemyAPI.GetNFTTransfers(blockchainNetwork, castedArgs.ContractAddresses, castedArgs.SyncedBlock, "latest", castedArgs.PageKey, TransferPageSize)
	if err != nil {
		panic(fmt.Sprintf("SyncNFTTranfers: failed to get NFT transfers: %s", err))
	}

	nftTransfers := make([]ethmodel.NFTTransfer, 0, len(res.Result.Transfers))
	for _, transfer := range res.Result.Transfers {

		tokenID := hexstring.MustParse(transfer.TokenID)
		blockNum := hexstring.MustParse(transfer.BlockNum)

		nftTransfers = append(nftTransfers, ethmodel.NFTTransfer{
			Blockchain:      castedArgs.Blockchain,
			Network:         castedArgs.Network,
			ContractAddress: transfer.RawContract.Address,
			TokenID:         *bunbig.FromMathBig(tokenID.ToBigInt()),
			BlockNumber:     *bunbig.FromMathBig(blockNum.ToBigInt()),
			FromAddress:     transfer.From,
			ToAddress:       transfer.To,
			TransactionHash: transfer.Hash,
		})
	}

	// Rico: NFT Owners are automatically updated when we insert new records to the NFT transfers table.
	// While we know the sequential ordering for the NFT Transfers based on the block number,
	// we can't be sure that a single NFT token wouldn't be transferred multiple times in a single block.
	// Given the low chance of it happening and the fact that it will be automatically resolve if the token is transferred again,
	// we will skip it for now.
	err = h.NftTransferMutator.InsertNFTTransfers(nftTransfers)
	if err != nil {
		panic(fmt.Sprintf("SyncNFTTranfers: failed to insert NFT transfers: %s", err))
	}

	if res.Result.PageKey != "" {
		_, err := workers.Enqueue(h.Config.Worker.TransferQueueName, "Sync", SyncETHNFTTransfersMessageArgs{
			Blockchain:        castedArgs.Blockchain,
			Network:           castedArgs.Network,
			ContractAddresses: castedArgs.ContractAddresses,
			SyncedBlock:       castedArgs.SyncedBlock,
			PageKey:           res.Result.PageKey,
		})

		if err != nil {
			panic(fmt.Sprintf("SyncNFTTranfers: failed to enqueue NFT transfers: %s", err))
		}
	}
}
