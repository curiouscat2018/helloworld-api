package vault

type mockVault struct {
}

func NewMockVault() Vault {
	return mockVault{}
}

func (v mockVault) GetSecret(url string) (string, error) {
	return "this is secret from mock vault", nil
}
