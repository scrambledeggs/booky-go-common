package photo

import (
	"strconv"
	"strings"
)

const (
	ASSET_URL string = "https://assets1.phonebooky.com"
)

// Format image filename
// Accept asset id, asset type, image filename (Optional: image size, image type)
// Image URL Format : ASSET_URL/ASSET_TYPE/IMAGE_TYPE/ID_URL/IMAGE_SIZE/FILENAME
func FormatImageURL(id int, assetType string, filename string, extra ...string) string {
	if filename == "" {
		return ""
	}

	// default values
	imageSize := "original"
	imageType := "photos"
	if len(extra) > 0 {
		if extra[0] != "" {
			imageSize = extra[0]
		}

		if len(extra) > 1 && extra[1] != "" {
			imageType = extra[1]
		}
	}

	return strings.Join([]string{ASSET_URL, assetType, imageType, FormatIDUrl(id), imageSize, filename}, "/")
}

// Format ID for URL
func FormatIDUrl(id int) string {
	standardLength := 9

	// Pad with 0 up to 9 digits
	idURL := strings.Repeat("0", standardLength) + strconv.Itoa(id)
	idURL = idURL[(len(idURL) - standardLength):]
	// Add / for every len 3 e.g. 000/009/012
	idURL = idURL[:3] + "/" + idURL[3:6] + "/" + idURL[6:9]

	return idURL
}
