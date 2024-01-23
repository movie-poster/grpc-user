package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
)

func DecodeBase64ToBytes(strBase64 string) (*bytes.Buffer, float64, error) {
	data, err := base64.StdEncoding.DecodeString(strBase64)
	if err != nil {
		return nil, 0, errors.New("no es posible decodificar el archivo debido a errores existentes")
	}

	buf := bytes.NewBuffer(data)
	fileSizeInBytes := buf.Len()
	fileSizeInMB := float64(fileSizeInBytes) / (1024 * 1024)

	return buf, fileSizeInMB, nil
}

func ConvertStringToBuffer(content string) (*bytes.Buffer, float64, error) {
	buf := bytes.NewBuffer([]byte(content))

	fileSizeInBytes := buf.Len()
	fileSizeInMB := float64(fileSizeInBytes) / (1024 * 1024)
	return buf, fileSizeInMB, nil
}
