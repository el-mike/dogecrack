package provider

// ProviderInstanceManager - entity responsible for managing actual, provider's instances.
type ProviderInstanceManager interface {
	// Sync - checks and synchronizes provider's instances.
	Sync() ([]ProviderInstance, error)

	// RunInstance - runs single instance.
	RunInstance(fileUrl, wallet string) (ProviderInstance, error)

	// GetInstance - returns active(rented) instance with given ID.
	GetInstance(instanceId int) (ProviderInstance, error)
}
