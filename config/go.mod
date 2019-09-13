module github.com/scrambledeggs/booky-go-common/config

go 1.12

require (
	github.com/scrambledeggs/booky-go-common/encryption v0.0.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
)

replace github.com/scrambledeggs/booky-go-common/encryption v0.0.0 => ../encryption
