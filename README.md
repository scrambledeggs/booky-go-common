# booky-go-common
![Booky Logo](https://user-images.githubusercontent.com/96253/56195205-17e08980-6067-11e9-9488-d0dcd80b5ebf.png)

# Description
Contains common functions that can be used by other booky go modules

# [Module] encryption/
Data is encrypted using `AES` given a passphrase that is hashed using `sha256`

Usage:
- `Encrypt(data, passphrase)` - Returns utf-8 encrypted string
- `EncryptB64(data, passphrase)` - Returns base64 encrypted string
- `Decrypt(ciphertext, passphrase)` - Decrypts data in byte mode
- `DecryptB64(ciphertext, passphrase)` - Decrypts data encoded in b64

Encrypting a string:
- Run `go run runnables/encryption/main.go <string> <key>` in order to get encrypted key. Key should be at least 12 in length.

# [Module] config/
Used to load environment variables given a valid environment, config file, and target data

Usage:
- `New(application_environment, config_map)` - *Required.* Reads through config file.
	- Parameters:
		- **application_environment** - Defaults to **development**. Accepts **development**, **staging**, **test**, and **production** values only.
		- **config_map**- (map[string]string) String array which contains config keys to set that will be set to become ENV variables. Stored in the format of **[ENV_VARIABLE_NAME] = [config-key]** (key-value). Encrypted data should be enclosed in **ENC()** (E.g. **config-key: ENC(encrypted-data)**) to flag it for decryption.
- `ApplyEnvConfig()`- *Required.* Sets environment variables based on config read.
- `SetConfig(configFile, configFilePath, configFileType)` - *Optional.* Sets config file name and path. Defaults to **./config.yml**
- `SetCipherPass(passphrase)` - *Required when config contains encrypted keys.* Provide passphrase for the encrypted data

# [Module] photo/
Generates Booky's Image URL. Module is named `photo` as to not override go's `image` package

Usage: FormatImageURL(ID int, assetType string, filename string, extra ...string)
 	- ID - ID of entity.
	- assetType - Type of entity (e.g. 'offers' or 'brands')
	- filename - Image filename
	- extra - Accepts up to two optional parameters. Sets imageSize(default:`original`) and imageType(default:`photo`).
	- Sample Output: "https://assets1.phonebooky.com/brands/photos/000/000/020/original/sample.jpg"

# [Module] marshalling/
Customized MarshalMap for dynamodb attributes. This will keep empty values (empty strings, zeros, empty structs/objects) and persist them in dynamodb tables.

Usage: 
+ CustomMarshalMap(in interface{})
	- in - interface{} of value struct such as `{Renamed Brand 1 0xc0000874b0  inactive 0xc0000908c0}`

	- Sample Output of type `map[string]*dynamodb.AttributeValue`:

	```go
	map[brand_name:{
		S: "Renamed Brand 1"
	} brand_status:{
		S: "inactive"
	} description:{
		NULL: true
	} offer:{
		NULL: true
	} offer_limit:{
		N: "0"
	}]
	```