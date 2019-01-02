package vault

type Vault interface {
	GetSecret(url string) (string, error)
}
