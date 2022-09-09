package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

var log = logging.Logger("citizen5/crypto")

//decodePEMFile reads and decodes generic PEM files.
func decodePEMFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	p, _ := pem.Decode(buf)
	if p == nil {
		return nil, fmt.Errorf("no pem block found")
	}
	return p.Bytes, nil
}

func GetIdentity(pubkey string) peer.ID {
	pubb, err := base64.StdEncoding.DecodeString(pubkey)
	if err != nil {
		panic(err)
	}
	pub, err := crypto.UnmarshalPublicKey(pubb)
	if err != nil {
		panic(err)
	}
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(err)
	}
	return id
}

func GenerateIdentity() (string, string) {
	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		panic(err)
	}
	privkeyb, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		panic(err)
	}
	pubkeyb, err := crypto.MarshalPublicKey(pub)
	if err != nil {
		panic(err)
	}
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(err)
	}
	log.Infof("generated identity %s", id.Pretty())
	return base64.StdEncoding.EncodeToString(privkeyb), base64.StdEncoding.EncodeToString(pubkeyb)
}

func GenerateED25519Key(private string, public string) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	//Create file and write public key
	pubOut, err := os.OpenFile(public, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("failed to create %s file: %s", private, err)
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Fatalf("unable to marshal public key: %v", err)
	}
	//Encode public key using PEM format
	if err := pem.Encode(pubOut, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}); err != nil {
		log.Fatalf("failed to write data to %s file: %s", public, err)
	}
	if err := pubOut.Close(); err != nil {
		log.Fatalf("error closing %s file: %s", public, err)
	}

	//Create file and write private key
	keyOut, err := os.OpenFile(private, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("failed to create %s file: %s", private, err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		log.Fatalf("unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}); err != nil {
		log.Fatalf("Failed to write data to %s file: %s", private, err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("error closing %s file: %s", private, err)
	}
}

//GetPrivateKey reads the private key from input file and
//returns the initialized PrivateKey.
func GetED25519PrivateKey(privateKey string) (ed25519.PrivateKey, error) {
	p, _ := decodePEMFile(privateKey)
	key, err := x509.ParsePKCS8PrivateKey(p)
	if err != nil {
		return nil, err
	}
	edKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not ed25519 key")
	}
	return ed25519.PrivateKey(edKey), nil
}

//GetPublicKey reads the public key from input file and
//returns the initialized PublicKey.
func GetED25519PublicKey(publicKey string) (ed25519.PublicKey, error) {
	p, _ := decodePEMFile(publicKey)
	key, err := x509.ParsePKIXPublicKey(p)
	if err != nil {
		return nil, err
	}
	edKey, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not ed25519 key")
	}
	return ed25519.PublicKey(edKey), nil
}

func ED25519Sign(p ed25519.PrivateKey, data []byte) (string, error) {
	signature := ed25519.Sign(ed25519.PrivateKey(p), data)
	return hex.EncodeToString(signature), nil
}

//Sign reads the input file and compute the ED25519 signature
//using the private key.
func ED25519SignFile(p ed25519.PrivateKey, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return ED25519Sign(p, buf)
}

//Verify checks that input signature is valid. That is, if
//input file was signed by private key corresponding to input
//public key.
func ED25519Verify(signature string, p ed25519.PublicKey, data []byte) (bool, error) {
	byteSign, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}
	ok := ed25519.Verify(p, data, byteSign)
	return ok, nil
}

func ED25519VerifyFile(signature string, p ed25519.PublicKey, file string) (bool, error) {
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return false, err
	}
	return ED25519Verify(signature, p, buf)
}
