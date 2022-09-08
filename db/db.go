package db

import (
	"context"
	"os"
	"path/filepath"

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

	"github.com/allisterb/citizen5/crypto"
	"github.com/allisterb/citizen5/util"
)

var log = logging.Logger("citizen5/db")

func initIPFSRepo(ctx context.Context, privkey string, pubkey string) repo.Repo {
	pid := crypto.GetIdentity(pubkey)
	c := cfg.Config{}
	c.Pubsub.Enabled = cfg.True
	c.Bootstrap = []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	}
	c.Addresses.Swarm = []string{"/ip4/127.0.0.1/tcp/4001", "/ip4/127.0.0.1/udp/4001/quic"}
	c.Identity.PeerID = pid.Pretty()
	c.Identity.PrivKey = privkey

	return &repo.Mock{
		D: dsync.MutexWrap(ds.NewMapDatastore()),
		C: c,
	}
}

func InitIPFSApi(ctx context.Context, privkey string, pubkey string) (iface.CoreAPI, func(), error) {
	node, err := ipfsCore.NewNode(ctx, &ipfsCore.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo:    initIPFSRepo(ctx, privkey, pubkey),
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	log.Infof("IPFS Node created. We are: %s", node.Identity.Pretty())
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

func CreateDB(ctx context.Context, privkey string, pubkey string) error {
	h := util.GetUserHomeDir()
	datapath := filepath.Join(h, ".citizen5")
	dbpath := filepath.Join(datapath, "db")
	os.MkdirAll(datapath, os.ModePerm)
	ipfs, cleanup, err := InitIPFSApi(ctx, privkey, pubkey)
	if err != nil {
		return err
	}
	log.Infof("creating OrbitDB database at local path %s...", dbpath)
	d, e := orbitdb.NewOrbitDB(ctx, ipfs, &orbitdb.NewOrbitDBOptions{
		//DirectChannelFactory: directchannel.InitDirectChannelFactory(zap.NewNop(), node1.PeerHost),
		Directory: &dbpath,
		Logger:    log.Desugar(),
	})
	if e != nil {
		return err
	}
	docs, err := d.Docs(ctx, "reports", nil)
	if err != nil {
		log.Errorf("could not create OrbitDB document database: %v", err)
		return err
	} else {
		log.Infof("created OrbitDB document database 'reports' at IPFS address %s.", docs.Address().String())
	}
	d.Close()
	cleanup()
	return nil
}
