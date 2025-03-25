package common

// ECCKeyPair is a struct that contains ECC private key and public key
type ECCKeyPair struct {
	ECCPrivateKey string `json:"ecc_private_key" mapstructure:"ecc_private_key"`
	ECCPublicKey  string `json:"ecc_public_key" mapstructure:"ecc_public_key"`
}
