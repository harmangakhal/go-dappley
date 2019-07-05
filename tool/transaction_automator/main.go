package main

import (
	"io/ioutil"
	_ "net/http/pprof"
	"os"

	"github.com/dappley/go-dappley/common"
	"github.com/dappley/go-dappley/config"
	"github.com/dappley/go-dappley/sdk"
	tool "github.com/dappley/go-dappley/tool/tool_util"
	tx_automator_configpb "github.com/dappley/go-dappley/tool/transaction_automator/pb"
	"github.com/dappley/go-dappley/tool/transaction_automator/util"
	logger "github.com/sirupsen/logrus"
)

const (
	contractAddrFilePath = "contract/contractAddr"
	contractFilePath     = "contract/test_contract.js"
	configFilePath       = "default.conf"
)

func main() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})

	logger.Info("*************************************")
	logger.Info("**Transaction automator tool starts**")
	logger.Info("*************************************")

	dappSdk, wallet, toolConfigs := initial_setup()

	nextBlockTicker := tool.NewNextBlockTicker(dappSdk)
	nextBlockTicker.Run()

	logger.Info("Start funding...")

	waitTillBlockHeightTwo(nextBlockTicker, dappSdk)
	fund(dappSdk, wallet, toolConfigs.GetInitialAmount())

	isScDeployed, scAddr := deploySmartContract(dappSdk, wallet)

	sender := util.NewBatchTxSender(toolConfigs.GetTps(), wallet, dappSdk, toolConfigs.GetScFreq(), scAddr)
	if isScDeployed {
		sender.EnableSmartContract()
	}
	sender.Run()

	for {
		select {
		case <-nextBlockTicker.GetTickerChan():
			sender.EnableSmartContract()
			wallet.DisplayBalances()
		}
	}
}

func initial_setup() (*sdk.DappSdk, *sdk.DappSdkWallet, *tx_automator_configpb.Config) {
	toolConfigs := &tx_automator_configpb.Config{}
	config.LoadConfig(configFilePath, toolConfigs)

	grpcClient := sdk.NewDappSdkGrpcClient(toolConfigs.GetPort())
	dappSdk := sdk.NewDappSdk(grpcClient)
	wallet := sdk.NewDappSdkWallet(
		toolConfigs.GetMaxWallet(),
		toolConfigs.GetPassword(),
		dappSdk,
	)

	return dappSdk, wallet, toolConfigs
}

func waitTillBlockHeightTwo(ticker *tool.NextBlockTicker, dappSdk *sdk.DappSdk) {
	logger.Info("Waiting till the next block is mined...")
	for {
		select {
		case <-ticker.GetTickerChan():
			blkHeight, _ := dappSdk.GetBlockHeight()
			if blkHeight > 1 {
				return
			}
		}
	}
}

func fund(dappSdk *sdk.DappSdk, wallet *sdk.DappSdkWallet, initialAmount uint64) {
	fundAddr := getFundAddr(wallet)
	fundRequest := tool.NewFundRequest(dappSdk)
	fundRequest.Fund(fundAddr, common.NewAmount(initialAmount))

	logger.WithFields(logger.Fields{
		"initial_total_amount": initialAmount,
	}).Info("Funding is completed. Script starts.")

	wallet.Update()
	wallet.DisplayBalances()
}

func getFundAddr(wallet *sdk.DappSdkWallet) string {
	return wallet.GetAddrs()[0].String()
}

func deploySmartContract(dappSdk *sdk.DappSdk, wallet *sdk.DappSdkWallet) (bool, string) {

	from := getFundAddr(wallet)
	smartContractAddr := getSmartContractAddr()
	if smartContractAddr != "" {
		logger.WithFields(logger.Fields{
			"contractAddr": smartContractAddr,
		}).Info("Smart contract has already been deployed. If you are sure it is not deployed, empty the file:", contractAddrFilePath)
		return true, smartContractAddr
	}

	data, err := ioutil.ReadFile(contractFilePath)
	if err != nil {
		logger.WithError(err).WithFields(logger.Fields{
			"file_path": contractFilePath,
		}).Panic("Unable to read smart contract file!")
	}

	contract := string(data)
	resp, err := dappSdk.Send(from, "", 1, contract)
	if err != nil {
		logger.WithError(err).WithFields(logger.Fields{
			"file_path":     contractFilePath,
			"contract_addr": smartContractAddr,
		}).Panic("Deploy smart contract failed!")
	}
	smartContractAddr = resp.ContractAddress

	recordSmartContractAddr(smartContractAddr)

	logger.WithFields(logger.Fields{
		"contract_addr": smartContractAddr,
	}).Info("Smart contract has been deployed")

	wallet.Update()

	return false, smartContractAddr
}

func getSmartContractAddr() string {
	bytes, err := ioutil.ReadFile(contractAddrFilePath)
	if err != nil {
		logger.WithError(err).WithFields(logger.Fields{
			"file_path": contractAddrFilePath,
		}).Panic("Unable to read file!")
	}
	return string(bytes)
}

func recordSmartContractAddr(addr string) {
	err := ioutil.WriteFile(contractAddrFilePath, []byte(addr), os.FileMode(777))
	if err != nil {
		logger.WithError(err).WithFields(logger.Fields{
			"file_path":     contractAddrFilePath,
			"contract_addr": addr,
		}).Panic("Unable to record smart contract address!")
	}
}
