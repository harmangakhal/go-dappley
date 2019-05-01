package sdk

import (
	"context"
	"github.com/dappley/go-dappley/core"
	"github.com/dappley/go-dappley/core/pb"
	"github.com/dappley/go-dappley/rpc/pb"
	"github.com/dappley/go-dappley/storage"
)

type DappSdkUtxoIndex struct {
	conn      *DappSdkConn
	sdkWallet *DappSdkWallet
	utxoIndex *core.UTXOIndex
}

func NewDappleySdkUtxoIndex(conn *DappSdkConn, sdkWallet *DappSdkWallet) *DappSdkUtxoIndex {
	return &DappSdkUtxoIndex{
		conn:      conn,
		utxoIndex: core.NewUTXOIndex(core.NewUTXOCache(storage.NewRamStorage())),
		sdkWallet: sdkWallet,
	}
}

func (sdkui *DappSdkUtxoIndex) update() error {

	for _, addr := range sdkui.sdkWallet.addrs {

		kp := sdkui.sdkWallet.wm.GetKeyPairByAddress(addr)
		_, err := core.NewUserPubKeyHash(kp.PublicKey)
		if err != nil {
			return err
		}

		utxos, err := sdkui.getUtxoByAddr(addr)
		if err != nil {
			return err
		}

		for _, utxoPb := range utxos {
			utxo := core.UTXO{}
			utxo.FromProto(utxoPb)
			sdkui.utxoIndex.AddUTXO(utxo.TXOutput, utxo.Txid, utxo.TxIndex)
		}
	}

	return nil
}

func (sdkui *DappSdkUtxoIndex) getUtxoByAddr(addr core.Address) ([]*corepb.Utxo, error) {
	resp, err := sdkui.conn.rpcClient.RpcGetUTXO(context.Background(), &rpcpb.GetUTXORequest{
		Address: addr.String(),
	})

	if err != nil || resp == nil {
		return nil, err
	}

	return resp.Utxos, nil
}
