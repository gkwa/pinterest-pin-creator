package pinterest

import (
	"encoding/base64"
	"os"
)

func toBase64(imgPath string) string {
	bytes, err := os.ReadFile(imgPath)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}
