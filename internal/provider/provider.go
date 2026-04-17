package provider

type Provider interface {
	Get() ([]byte, error)

	Disconnect()
}
