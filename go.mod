module github.com/scrambledeggs/booky-go-common

go 1.12

require (
	github.com/scrambledeggs/booky-go-common/config v0.0.0 // indirect
	github.com/scrambledeggs/booky-go-common/encryption v0.0.0-20190913074022-cfd0ef746a40
	github.com/scrambledeggs/booky-go-common/rds v0.0.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
)

replace github.com/scrambledeggs/booky-go-common/config v0.0.0 => ./config

replace github.com/scrambledeggs/booky-go-common/encryption v0.0.0 => ./encryption

replace github.com/scrambledeggs/booky-go-common/rds v0.0.0 => ./rds
