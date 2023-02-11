package mod

type Target interface {
	// GetName return secret name
	GetName() string
	// GetKey return mod key
	GetKey() string
	// GetSecret return secret
	GetSecret() string
	GetValidityPeriod() int64
	GetCreateTime() int64
	GetLastUsed() int64
	GetExpired() int64
}
