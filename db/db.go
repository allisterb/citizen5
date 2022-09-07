package db

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	logging "github.com/ipfs/go-log/v2"

	orbitdb "berty.tech/go-orbit-db"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	cfg "github.com/ipfs/go-ipfs/config"
	ipfsCore "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	repo "github.com/ipfs/go-ipfs/repo"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

var log = logging.Logger("db")

func initIPFSRepo(ctx context.Context) repo.Repo {
	c := cfg.Config{}
	//dc, _ := cfg.Init(io.Discard, 2048)
	//log.Info(dc.API)
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
	c.Bootstrap = []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	}
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
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo:    initIPFSRepo(ctx),
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
	node, err := ipfsCore.NewNode(ctx, &ipfsCore.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo:    initIPFSRepo(ctx),
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	log.Infof("IPFS Node created. We are: %s", node.PeerHost.ID())
	c, e := coreapi.NewCoreAPI(node)
	if e != nil {
		return nil, nil, e
	} else {
		clean := func() {
			node.Close()
		}
		return c, clean, e
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
}
