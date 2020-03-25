module github.com/scrambledeggs/booky-go-common

go 1.12

require (
	github.com/scrambledeggs/booky-go-common/config v0.0.0-20191003085113-ddb19ea6d89c // indirect
	github.com/scrambledeggs/booky-go-common/encryption v0.0.0-20191003071505-64d725be3c0
	github.com/scrambledeggs/booky-go-common/rds v0.0.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
)

replace github.com/scrambledeggs/booky-go-common/config v0.0.0 => ./config

replace github.com/scrambledeggs/booky-go-common/encryption v0.0.0 => ./encryption

replace github.com/scrambledeggs/booky-go-common/rds v0.0.0 => ./rds
