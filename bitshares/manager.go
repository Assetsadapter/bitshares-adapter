/*
 * Copyright 2018 The OpenWallet Authors
 * This file is part of the OpenWallet library.
 *
 * The OpenWallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OpenWallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package bitshares

import (
	"github.com/Assetsadapter/bitshares-adapter/addrdec"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/denkhaus/bitshares"
	"github.com/denkhaus/bitshares/config"
)

type WalletManager struct {
	openwallet.AssetsAdapterBase

	Api             *WalletClient                   // 节点客户端
	Config          *WalletConfig                   // 节点配置
	Decoder         openwallet.AddressDecoder       //地址编码器
	DecoderV2       openwallet.AddressDecoderV2     //地址编码器V2
	TxDecoder       openwallet.TransactionDecoder   //交易单编码器
	Log             *log.OWLogger                   //日志工具
	ContractDecoder openwallet.SmartContractDecoder //智能合约解析器
	Blockscanner    *BtsBlockScanner                //区块扫描器
	CacheManager    openwallet.ICacheManager        //缓存管理器
	WebsocketAPI    bitshares.WebsocketAPI          //bitshares WebsocketAPI
}

func NewWalletManager(cacheManager openwallet.ICacheManager) *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol)
	wm.Api = NewWalletClient(wm.Config.ServerAPI, wm.Config.WalletAPI, false)
	wm.Blockscanner = NewBlockScanner(&wm)
	wm.Decoder = NewAddressDecoder(&wm)
	wm.DecoderV2 = addrdec.NewAddressDecoderV2()
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	wm.CacheManager = cacheManager
	wm.ContractDecoder = NewContractDecoder(&wm)

	wm.WebsocketAPI = NewWebsocketAPI(wm.Config.ServerWS)
	//privateNetConfig :=config.ChainConfig{
	//	Name:      "private_net",
	//	CoreAsset: "CORE",
	//	Prefix:    "BTS",
	//	ID:"c3807ea4b7d28da76249c6be53dc6342a12888327aab9d66ac1d53e31e797ae3",
	//}
	//config.Add(privateNetConfig)
	//config.SetCurrent("c3807ea4b7d28da76249c6be53dc6342a12888327aab9d66ac1d53e31e797ae3")
	config.SetCurrent(config.ChainIDBTS)

	return &wm
}

func NewWebsocketAPI(api string) bitshares.WebsocketAPI {
	config.SetCurrent(config.ChainIDBTS)
	return bitshares.NewWebsocketAPI(api)
}
