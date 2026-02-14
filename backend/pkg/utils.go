package pkg

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

// Check is number or not.
func IsNumber(ctx context.Context, numberString string, logger *log.Logger) bool {
	_, err := strconv.Atoi(numberString)
	if err != nil {
		logger.Debugf("This is not number current_value %s error: %s",
			numberString, err,
		)

		return false
	}

	return true
}

func StrToBool(logger log.FieldLogger, varString string) bool {
	boolValue, err := strconv.ParseBool(varString)
	if err != nil {
		logger.Debugf("Error convert to bool %s error: %s", varString, err)
		return false
	}
	return boolValue
}

func StrToInt(logger log.FieldLogger, varString string) int {
	intValue, err := strconv.Atoi(varString)
	if err != nil {
		logger.Debugf("Error convert to int %s error: %s", varString, err)
		return 0
	}
	return intValue
}

func PrettyfyPhone(str, replacement string) string {
	return string([]rune(str)[:0]) + replacement + string([]rune(str)[1:])
}

// UUIDToHex - возвращает UUID в виде hex.
func UUIDToHex(varUUID uuid.UUID) string {
	return strings.ReplaceAll(varUUID.String(), "-", "")
}
