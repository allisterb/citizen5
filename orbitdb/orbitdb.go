package orbitdb

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	iface "github.com/ipfs/interface-go-ipfs-core"
	cfg "github.com/ipfs/kubo/config"
	ipfsCore "github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/repo"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func initIPFSRepo(ctx context.Context) repo.Repo {
	c := cfg.Config{}
	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		panic(err)
	}

	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(err)
	}

	privkeyb, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		panic(err)
	}

	c.Pubsub.Enabled = cfg.True
	c.Bootstrap = []string{}
	c.Addresses.Swarm = []string{"/ip4/127.0.0.1/tcp/4001", "/ip4/127.0.0.1/udp/4001/quic"}
	c.Identity.PeerID = pid.Pretty()
	c.Identity.PrivKey = base64.StdEncoding.EncodeToString(privkeyb)

	return &repo.Mock{
		D: dsync.MutexWrap(ds.NewMapDatastore()),
		C: c,
	}
}

func createIPFSAPIs(ctx context.Context, count int) ([]iface.CoreAPI, func()) {
	coreAPIs := make([]iface.CoreAPI, count)
	cleans := make([]func(), count)

	for i := 0; i < count; i++ {
		core, err := ipfsCore.NewNode(ctx, &ipfsCore.BuildCfg{
			Online: true,
			Repo:   initIPFSRepo(ctx),
			ExtraOpts: map[string]bool{
				"pubsub": true,
			},
		})
		//require.NoError(t, err)

		//coreAPIs[i] = testingCoreAPI(t, core)
		cleans[i] = func() {
			core.Close()
		}
	}

	return coreAPIs, func() {
		for i := 0; i < count; i++ {
			cleans[i]()
		}
	}
}
