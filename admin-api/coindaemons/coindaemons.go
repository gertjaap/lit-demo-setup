package coindaemons

import "github.com/mit-dci/lit/btcutil/chaincfg"

type CoinDaemon struct {
	ImageID                string
	ImageName              string
	DataFolderInContainer  string
	DataSubFolderOnHost    string
	ConfigName             string
	P2PPort                uint
	RPCPort                uint
	ContainerName          string
	LitConfigPrefix        string
	Command                []string
	CoinParams             chaincfg.Params
	LitCoinType            uint32
	InitialFunding         uint32
	InitialBlocks          uint32
	NodeChannelCapacity    int64
	NodeChannelInitialSend int64
}

var CoinDaemons = []CoinDaemon{
	{
		ImageID:                "",
		ImageName:              "bitcoind",
		DataFolderInContainer:  "/bitcoin/.bitcoin",
		DataSubFolderOnHost:    "bitcoind",
		ConfigName:             "bitcoin.conf",
		P2PPort:                18444,
		RPCPort:                18443,
		ContainerName:          "litdemobtcregtest",
		LitConfigPrefix:        "reg",
		LitCoinType:            257,
		InitialFunding:         10000,
		InitialBlocks:          3000,
		NodeChannelCapacity:    int64(100000000),
		NodeChannelInitialSend: int64(400000),
		CoinParams: chaincfg.Params{
			// Address encoding magics
			PubKeyHashAddrID: 0x6f, // starts with m or n
			ScriptHashAddrID: 0xc4, // starts with 2
			PrivateKeyID:     0xef, // starts with 9 (uncompressed) or c (compressed)
			Bech32Prefix:     "bcrt",

			// BIP32 hierarchical deterministic extended key magics
			HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
			HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
		},
	},
	{
		ImageID:                "",
		ImageName:              "litecoind",
		DataFolderInContainer:  "/home/litecoin/.litecoin",
		DataSubFolderOnHost:    "litecoind",
		ConfigName:             "litecoin.conf",
		P2PPort:                19444,
		RPCPort:                19443,
		LitCoinType:            258,
		ContainerName:          "litdemoltcregtest",
		LitConfigPrefix:        "litereg",
		InitialFunding:         10000,
		InitialBlocks:          3000,
		NodeChannelCapacity:    int64(500000000),
		NodeChannelInitialSend: int64(50000000),
		CoinParams: chaincfg.Params{
			// Address encoding magics
			PubKeyHashAddrID: 0x6f, // starts with m or n
			ScriptHashAddrID: 0xc4, // starts with 2
			Bech32Prefix:     "rltc",
			PrivateKeyID:     0xef, // starts with 9 7(uncompressed) or c (compressed)

			// BIP32 hierarchical deterministic extended key magics
			HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
			HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
		},
	},
	/*{
		ImageID:               "sha256:32f8620c9d3e9d20ae6fe6d19806b6acbb6fc37929fc33d9a926d10194b81af5",
		DataFolderInContainer: "/data",
		DataSubFolderOnHost:   "vertcoind",
		ConfigName:            "vertcoin.conf",
		P2PPort:               18444,
		RPCPort:               18443,
		ContainerName:         "litdemovtcregtest",
		LitConfigPrefix:       "rtvtc",
		LitCoinType:           261,
		CoinParams: chaincfg.Params{
			PubKeyHashAddrID: 0x6f,
			ScriptHashAddrID: 0xc4,
			Bech32Prefix:     "rvtc",
			PrivateKeyID:     0xef,

			// BIP32 hierarchical deterministic extended key magics
			HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
			HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

		},
	},*/
	{
		ImageID:                "",
		ImageName:              "dummyusdd",
		DataFolderInContainer:  "/bitcoin/.bitcoin",
		DataSubFolderOnHost:    "dummyusdd",
		ConfigName:             "bitcoin.conf",
		P2PPort:                26999,
		RPCPort:                18443,
		ContainerName:          "litdemousdregtest",
		LitConfigPrefix:        "dusd",
		LitCoinType:            262,
		InitialFunding:         200000,
		InitialBlocks:          10000,
		NodeChannelCapacity:    int64(10000000000),
		NodeChannelInitialSend: int64(1000000000),
		CoinParams: chaincfg.Params{
			PubKeyHashAddrID: 0x1e, // starts with D
			ScriptHashAddrID: 0x5a, // starts with d
			PrivateKeyID:     0x83, // starts with u
			Bech32Prefix:     "dusd",

			// BIP32 hierarchical deterministic extended key magics
			HDPrivateKeyID: [4]byte{0x04, 0xA5, 0xB3, 0xF4}, // starts with tprv
			HDPublicKeyID:  [4]byte{0x04, 0xA5, 0xB7, 0x8F}, // starts with tpub},
		},
	},
}
