package mod

type Target interface {
	GetName() string
	GetKey() string
	GetSecret() string
	GetValidityPeriod() int64
	GetCreateTime() int64
	GetLastUsed() int64
	GetExpired() int64
}
