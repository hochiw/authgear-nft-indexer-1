package task

import (
	"fmt"
	"log"
	"math/big"

	"github.com/authgear/authgear-nft-indexer/pkg/config"
	"github.com/authgear/authgear-nft-indexer/pkg/model"
	ethmodel "github.com/authgear/authgear-nft-indexer/pkg/model/eth"
	"github.com/authgear/authgear-nft-indexer/pkg/util/hexstring"
	"github.com/jrallison/go-workers"
)

type SycnNFTCollectionTaskCollectionQuery interface {
	QueryNFTCollections() ([]ethmodel.NFTCollection, error)
}

type SyncETHNFTCollectionTaskHandler struct {
	Config             config.Config
	NftCollectionQuery SycnNFTCollectionTaskCollectionQuery
}

func (h *SyncETHNFTCollectionTaskHandler) Handler(message *workers.Msg) {
	collections, err := h.NftCollectionQuery.QueryNFTCollections()
	if err != nil {
		panic(fmt.Sprintf("SyncNFTCollections: failed to query NFT collections: %s", err))
	}

	if len(collections) == 0 {
		panic("SyncNFTCollections: no NFT collections found")
	}

	nftContractAddressesByNetwork := make(map[model.BlockchainNetwork][]string, 0)
	smallestBlockByNetwork := make(map[model.BlockchainNetwork]*big.Int, 0)
	for _, collection := range collections {
		blockchainNetwork := model.BlockchainNetwork{
			Blockchain: collection.Blockchain,
			Network:    collection.Network,
		}
		nftContractAddressesByNetwork[blockchainNetwork] = append(nftContractAddressesByNetwork[blockchainNetwork], collection.ContractAddress)

		if smallestBlockByNetwork[blockchainNetwork] == nil || smallestBlockByNetwork[blockchainNetwork].Cmp(collection.SyncedBlockHeight.ToMathBig()) > 0 {
			smallestBlockByNetwork[blockchainNetwork] = collection.SyncedBlockHeight.ToMathBig()
		}
	}

	for network, contractAddresses := range nftContractAddressesByNetwork {
		smallestBlock := smallestBlockByNetwork[network]

		syncedBlockHex, err := hexstring.NewFromBigInt(smallestBlock)
		if err != nil {
			log.Fatalf("SyncNFTCollections: failed to convert smallest block to hex string: %s", err)
			continue
		}

		_, err = workers.EnqueueWithOptions(h.Config.Worker.TransferQueueName, "Sync", SyncETHNFTTransfersMessageArgs{
			Blockchain:        network.Blockchain,
			Network:           network.Network,
			ContractAddresses: contractAddresses,
			SyncedBlock:       syncedBlockHex.String(),
			PageKey:           "",
		}, workers.EnqueueOptions{Retry: true})

		if err != nil {
			panic(fmt.Sprintf("SyncNFTCollections: failed to enqueue NFT transfers: %s", err))
		}
	}
}
