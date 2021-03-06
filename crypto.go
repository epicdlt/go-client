package client

import (
	"context"
	"encoding/json"

	"github.com/epicdlt/go-client/models"
	"github.com/multiformats/go-multibase"
	"golang.org/x/crypto/ed25519"
	"gopkg.in/square/go-jose.v2"
)

func GenerateKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pubKey, pk, err := ed25519.GenerateKey(nil)
	return pubKey, pk, err
}

// GenerateKeyEncoded returns public key, private key, error
func GenerateKeyEncoded() (string, string, error) {
	pubKey, pk, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", err
	}
	return EncodeKey(pubKey), EncodeKey(pk), nil
}

func EncodeKey(k []byte) string {
	s, err := multibase.Encode(multibase.Base32, k)
	if err != nil {
		panic(err)
	}
	return s
}

func PrivateKeyFromHashed(ctx context.Context, privateKey string) (ed25519.PrivateKey, error) {
	_, pkBytes, err := multibase.Decode(privateKey)
	if err != nil {
		return nil, err
	}
	pk := ed25519.PrivateKey(pkBytes)
	return pk, nil
}

func PublicKeyFromPrivate(ctx context.Context, privateKey string) (string, error) {
	pk, err := PrivateKeyFromHashed(ctx, privateKey)
	if err != nil {
		return "", err
	}
	pubHashed, err := multibase.Encode(multibase.Base32, []byte(pk.Public().(ed25519.PublicKey)))
	if err != nil {
		return "", err
	}
	return pubHashed, nil
}

func Sign(ctx context.Context, privateKey string, toSign interface{}) (string, error) {
	pk, err := PrivateKeyFromHashed(ctx, privateKey)
	if err != nil {
		return "", err
	}
	return SignWithPK(ctx, pk, toSign)
}

func SignWithPK(ctx context.Context, pk ed25519.PrivateKey, toSign interface{}) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.EdDSA,
		Key:       pk,
	}, (&jose.SignerOptions{}).WithType("JWS"))
	if err != nil {
		return "", err
	}

	jsonValue, err := json.Marshal(toSign)
	if err != nil {
		return "", err
	}
	var payload = []byte(jsonValue)
	object, err := signer.Sign(payload)
	if err != nil {
		return "", err
	}

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()
	return serialized, nil
}

func VerifyTx(ctx context.Context, jws string) (*models.Transaction, error) {
	object, err := jose.ParseSigned(jws)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("jws parsed: %+v\n", object)

	unsafePayload := object.UnsafePayloadWithoutVerification()
	tx := &models.Transaction{}
	err = json.Unmarshal(unsafePayload, tx)
	if err != nil {
		return nil, err
	}

	// from is the publicKey
	from := tx.From

	// Now we can verify the signature on the payload. An error here would
	// indicate the the message failed to verify, e.g. because the signature was
	// broken or the message was tampered with.
	_, pubKeyBytes, err := multibase.Decode(from)
	if err != nil {
		return nil, err
	}
	pubKey := ed25519.PublicKey(pubKeyBytes)
	// todo: can use VerifyMulti here if we have multiple signatures
	payload, err := object.Verify(pubKey)
	if err != nil {
		return nil, err
	}
	tx = &models.Transaction{}
	err = json.Unmarshal(payload, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func Verify(ctx context.Context, publicKey string, jws string) ([]byte, error) {
	object, err := jose.ParseSigned(jws)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("jws parsed: %+v\n", object)

	// Now we can verify the signature on the payload. An error here would
	// indicate the the message failed to verify, e.g. because the signature was
	// broken or the message was tampered with.
	_, pubKeyBytes, err := multibase.Decode(publicKey)
	if err != nil {
		return nil, err
	}
	pubKey := ed25519.PublicKey(pubKeyBytes)
	payload, err := object.Verify(pubKey)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
