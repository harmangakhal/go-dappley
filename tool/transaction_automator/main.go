package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"

	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/dappley/go-dappley/client"
	"github.com/dappley/go-dappley/common"
	"github.com/dappley/go-dappley/config"
	"github.com/dappley/go-dappley/config/pb"
	"github.com/dappley/go-dappley/core"
	"github.com/dappley/go-dappley/logic"
	"github.com/dappley/go-dappley/rpc/pb"
)

var (
	password             = "testpassword"
	maxWallet            = 10
	initialAmount        = uint64(10)
	maxDefaultSendAmount = uint64(5)
	sendInterval         = time.Duration(1000) //ms
	checkBalanceInterval = time.Duration(10)   //s
	fundTimeout          = time.Duration(time.Minute * 5)
	currBalance          = make(map[string]uint64)
)

func main() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})

	var filePath string
	flag.StringVar(&filePath, "f", "conf/default_cli.conf", "CLI config file path")
	flag.Parse()

	cliConfig := &configpb.CliConfig{}
	config.LoadConfig(filePath, cliConfig)
	conn := initRpcClient(int(cliConfig.GetPort()))

	adminClient := rpcpb.NewAdminServiceClient(conn)
	rpcClient := rpcpb.NewRpcServiceClient(conn)

	addresses := createWallet()

	fundFromMiner(adminClient, rpcClient, addresses)
	logger.WithFields(logger.Fields{
		"initial_total_amount": initialAmount,
		"send_interval":        fmt.Sprintf("%d ms", sendInterval),
	}).Info("Funding is completed. Script starts.")
	displayBalances(rpcClient, addresses)

	ticker := time.NewTicker(time.Millisecond * sendInterval).C
	currHeight := getBlockHeight(rpcClient)
	for {
		select {
		case <-ticker:
			height := getBlockHeight(rpcClient)
			if height > currHeight {
				displayBalances(rpcClient, addresses)
				currHeight = height
			} else {
				sendRandomTransactions(adminClient, addresses)
			}
		}
	}
}

func initRpcClient(port int) *grpc.ClientConn {
	//prepare grpc client
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprint(":", port), grpc.WithInsecure())
	if err != nil {
		logger.WithError(err).Panic("Connection to RPC server failed.")
	}
	return conn
}

func createWallet() []core.Address {
	wm, err := logic.GetWalletManager(client.GetWalletFilePath())
	if err != nil {
		logger.Panic("Cannot get wallet manager.")
	}
	addresses := wm.GetAddresses()
	numOfWallets := len(addresses)
	for i := numOfWallets; i < maxWallet; i++ {
		_, err := logic.CreateWalletWithpassphrase(password)
		if err != nil {
			logger.WithError(err).Panic("Cannot create new wallet.")
		}
	}
	wm, err = logic.GetWalletManager(client.GetWalletFilePath())
	addresses = wm.GetAddresses()
	logger.WithFields(logger.Fields{
		"addresses": addresses,
	}).Info("Wallets are created")
	return addresses
}

func fundFromMiner(adminClient rpcpb.AdminServiceClient, rpcClient rpcpb.RpcServiceClient, addresses []core.Address) {
	logger.Info("Requesting fund from miner...")

	if len(addresses) == 0 {
		logger.Panic("There is no wallet to receive fund.")
	}

	fundAddr := addresses[0].String()

	requestFundFromMiner(adminClient, fundAddr)
	bal, isSufficient := checkSufficientInitialAmount(rpcClient, fundAddr)
	if isSufficient {
		//continue if the initial amount is sufficient
		return
	}
	logger.WithFields(logger.Fields{
		"address":    fundAddr,
		"balance":    bal,
		"target_fund": initialAmount,
	}).Info("Current wallet balance is insufficient. Waiting for more funds...")
	waitTilInitialAmountIsSufficient(adminClient, rpcClient, fundAddr)
}

func checkSufficientInitialAmount(rpcClient rpcpb.RpcServiceClient, addr string) (uint64, bool) {
	balance, err := getBalance(rpcClient, addr)
	if err != nil {
		logger.WithError(err).WithFields(logger.Fields{
			"address": addr,
		}).Panic("Failed to get balance.")
	}
	return uint64(balance), uint64(balance) >= initialAmount
}

func waitTilInitialAmountIsSufficient(adminClient rpcpb.AdminServiceClient, rpcClient rpcpb.RpcServiceClient, addr string) {
	checkBalanceTicker := time.NewTicker(time.Second * 5).C
	timeout := time.NewTicker(fundTimeout).C
	for {
		select {
		case <-checkBalanceTicker:
			bal, isSufficient := checkSufficientInitialAmount(rpcClient, addr)
			if isSufficient {
				//continue if the initial amount is sufficient
				return
			}
			logger.WithFields(logger.Fields{
				"address":     addr,
				"balance":     bal,
				"target_fund": initialAmount,
			}).Info("Current wallet balance is insufficient. Waiting for more funds...")
			requestFundFromMiner(adminClient, addr)
		case <-timeout:
			logger.WithFields(logger.Fields{
				"target_fund": initialAmount,
			}).Panic("Timed out while waiting for sufficient fund from miner!")
		}
	}
}

func requestFundFromMiner(adminClient rpcpb.AdminServiceClient, fundAddr string) {

	sendFromMinerRequest := rpcpb.SendFromMinerRequest{}
	sendFromMinerRequest.To = fundAddr
	sendFromMinerRequest.Amount = common.NewAmount(initialAmount).Bytes()

	_, err := adminClient.RpcSendFromMiner(context.Background(), &sendFromMinerRequest)
	if err != nil {
		logger.WithError(err).WithFields(logger.Fields{
			"fund_address": fundAddr,
		}).Panic("Failed to get test fund from miner.")
	}
}

func sendRandomTransactions(adminClient rpcpb.AdminServiceClient, addresses []core.Address) {

	fromIndex := getAddrWithBalance(addresses)
	toIndex := rand.Intn(maxWallet)
	for toIndex == fromIndex {
		toIndex = rand.Intn(maxWallet)
	}
	sendAmount := calcSendAmount(addresses[fromIndex].String(), addresses[toIndex].String())
	err := sendTransaction(adminClient, addresses[fromIndex].String(), addresses[toIndex].String(), sendAmount)
	sendTXLogger := logger.WithFields(logger.Fields{
		"from":             addresses[fromIndex].String(),
		"to":               addresses[toIndex].String(),
		"amount":           sendAmount,
		"sender_balance":   currBalance[addresses[fromIndex].String()],
		"receiver_balance": currBalance[addresses[toIndex].String()],
	})
	if err != nil {
		sendTXLogger.WithError(err).Panic("Failed to send transaction!")
		return
	}
	sendTXLogger.Info("Transaction is sent!")
}

func calcSendAmount(from, to string) uint64 {
	fromBalance, _ := currBalance[from]
	toBalance, _ := currBalance[to]
	amount := uint64(0)
	if fromBalance < toBalance {
		amount = 1
	} else if fromBalance == toBalance {
		amount = fromBalance - 1
	} else {
		amount = (fromBalance - toBalance) / 3
	}

	if amount == 0 {
		amount = 1
	}
	return amount
}

func getAddrWithBalance(addresses []core.Address) int {
	fromIndex := rand.Intn(maxWallet)
	amount := currBalance[addresses[fromIndex].String()]
	//TODO: add time out to this loop
	for amount <= maxDefaultSendAmount+1 {
		fromIndex = rand.Intn(maxWallet)
		amount = currBalance[addresses[fromIndex].String()]
	}
	return fromIndex
}

func sendTransaction(adminClient rpcpb.AdminServiceClient, from, to string, amount uint64) error {
	_, err := adminClient.RpcSend(context.Background(), &rpcpb.SendRequest{
		From:       from,
		To:         to,
		Amount:     common.NewAmount(amount).Bytes(),
		Tip:        0,
		Walletpath: client.GetWalletFilePath(),
		Contract:   "",
	})
	if err != nil {
		return err
	}
	currBalance[from] -= amount
	currBalance[to] += amount
	return nil
}

func displayBalances(rpcClient rpcpb.RpcServiceClient, addresses []core.Address) {
	for _, addr := range addresses {
		amount, err := getBalance(rpcClient, addr.String())
		balanceLogger := logger.WithFields(logger.Fields{
			"address": addr.String(),
			"amount":  amount,
			"record":  currBalance[addr.String()],
		})
		if err != nil {
			balanceLogger.WithError(err).Warn("Failed to get wallet balance.")
		}
		balanceLogger.Info("Displaying wallet balance...")
		currBalance[addr.String()] = uint64(amount)
	}
}

func getBalance(rpcClient rpcpb.RpcServiceClient, address string) (int64, error) {
	getBalanceRequest := rpcpb.GetBalanceRequest{}
	getBalanceRequest.Name = "getBalance"
	getBalanceRequest.Address = address
	response, err := rpcClient.RpcGetBalance(context.Background(), &getBalanceRequest)
	return response.Amount, err
}

func getBlockHeight(rpcClient rpcpb.RpcServiceClient) uint64 {
	resp, err := rpcClient.RpcGetBlockchainInfo(
		context.Background(),
		&rpcpb.GetBlockchainInfoRequest{})
	if err != nil {
		logger.WithError(err).Panic("Cannot get block height.")
	}
	return resp.BlockHeight
}

func isBalanceSufficient(addr string, amount uint64) bool {
	return currBalance[addr] >= amount
}
