# booky-go-common
Booky common functions in Go

![Booky Logo](https://user-images.githubusercontent.com/96253/56195205-17e08980-6067-11e9-9488-d0dcd80b5ebf.png)

# Description
Contains common functions that can be used by other booky go modules

# [Module] encryption/
Data is encrypted using `AES` given a passphrase that is hashed using `sha256`

Usage:
- Encrypt(data, passphrase) - Returns utf-8 encrypted string
- EncryptB64(data, passphrase) - Returns base64 encrypted string
- Decrypt(ciphertext, passphrase) - Decrypts data that was encrypted in the same module in
- DecryptB64(ciphertext, passphrase) - Decrypts data encoded in b64

# [Module] config/
Used to load environment variables given a valid environment, config file, and target data

# [Module] photo/
Generates Booky's Image URL. Module is named `photo` as to not override go's `image` package

Usage: FormatImageURL(ID int, assetType string, filename string, extra ...string)
 	- ID - ID of entity.
	- assetType - Type of entity (e.g. 'offers' or 'brands')
	- filename - Image filename
	- extra - Accepts up to two optional parameters. Sets imageSize(default:`original`) and imageType(default:`photo`).
	- Sample Output: "https://assets1.phonebooky.com/brands/photos/000/000/020/original/sample.jpg"
