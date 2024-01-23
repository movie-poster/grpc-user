package utils

import (
	"github.com/google/uuid"
)

func RandomNameObjectUUID() string {
	objID, _ := uuid.NewRandom()
	return objID.String()
}
