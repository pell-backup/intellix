package common

// ECCKeyPair is a struct that contains ECC private key and public key
type ECCKeyPair struct {
	ECCPrivateKeyPath string `json:"ecc_private_key_path" mapstructure:"ecc_private_key_path"`
	ECCPublicKey      string `json:"ecc_public_key" mapstructure:"ecc_public_key"`
}
