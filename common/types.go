package common

type ECCKeyPair struct {
	ECCPrivateKey string `json:"ecc_private_key"`
	ECCPublicKey  string `json:"ecc_public_key"`
}
