package db

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	orbitdb "berty.tech/go-orbit-db"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	cfg "github.com/ipfs/go-ipfs-config"
	repo "github.com/ipfs/go-ipfs/repo"
	iface "github.com/ipfs/interface-go-ipfs-core"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	ipfsCore "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
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

func InitIPFSNode(ctx context.Context) (*ipfsCore.IpfsNode, func()) {
	core, err := ipfsCore.NewNode(ctx, &ipfsCore.BuildCfg{
		Online: true,
		Repo:   initIPFSRepo(ctx),
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})
	if err != nil {
		panic(err)
	}
	cleanup := func() { core.Close() }
	return core, cleanup
}

func InitIPFSApi(ctx context.Context) (iface.CoreAPI, func(), error) {
	core, err := ipfsCore.NewNode(ctx, &ipfsCore.BuildCfg{
		Online: true,
		Repo:   initIPFSRepo(ctx),
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	c, e := coreapi.NewCoreAPI(core)
	if e != nil {
		return nil, nil, e
	} else {
		clean := func() {
			core.Close()
		}
		return c, clean, e
	}
}
func initIPFSAPIs(ctx context.Context, count int) ([]iface.CoreAPI, func()) {
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
		if err != nil {
			panic(err)
		}
		coreAPIs[i], err = coreapi.NewCoreAPI(core)
		if err != nil {
			panic(err)
		}
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

func CreateDB(ctx context.Context, name *string) (orbitdb.OrbitDB, func(), error) {
	ipfs, cleanup, err := InitIPFSApi(ctx)
	if err != nil {
		return nil, nil, err
	}
	d, e := orbitdb.NewOrbitDB(ctx, ipfs, &orbitdb.NewOrbitDBOptions{
		// DirectChannelFactory: directchannel.InitDirectChannelFactory(zap.NewNop(), node1.PeerHost),
		Directory: name,
	})
	return d, cleanup, e
	//orbitdb1.Identity()
}
