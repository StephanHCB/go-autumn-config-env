package auconfigenv

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func tstRequireErrorMsg(t *testing.T, expectedMsg string, actual error) {
	require.NotNil(t, actual)
	require.Contains(t, actual.Error(), expectedMsg)
}

func TestObtainPatternValidator(t *testing.T) {
	configValues = make(map[string]string, 0)
	configValues["pattern1"] = "abc"
	configValues["pattern2"] = "abbbbbbbc"
	configValues["pattern3_fail"] = "qabcd"

	cut := ObtainPatternValidator("^ab+c$")

	require.Nil(t, cut("pattern1"))
	require.Nil(t, cut("pattern2"))
	tstRequireErrorMsg(t, "must match ^ab+c$", cut("pattern3_fail"))
}

func TestObtainNotEmptyValidator(t *testing.T) {
	configValues = make(map[string]string, 0)
	configValues["notempty1"] = "abc"
	configValues["notempty2_fail"] = ""

	cut := ObtainNotEmptyValidator()

	require.Nil(t, cut("notempty1"))
	tstRequireErrorMsg(t, "must not be empty", cut("notempty2_fail"))
}

func TestObtainUintRangeValidator(t *testing.T) {
	configValues = make(map[string]string, 0)
	configValues["uint1"] = "4832"
	configValues["uint2_range"] = "9201"
	configValues["uint3_parse"] = "hallo"
	configValues["uint4_parse"] = "hallo 1234"
	configValues["uint5_neg"] = "-1"

	cut := ObtainUintRangeValidator(400, 9200)

	require.Nil(t, cut("uint1"))
	tstRequireErrorMsg(t, "value 9201 is out of range [400..9200]", cut("uint2_range"))
	tstRequireErrorMsg(t, "value hallo is not a valid integer: ", cut("uint3_parse"))
	tstRequireErrorMsg(t, "value hallo 1234 is not a valid integer: ", cut("uint4_parse"))
	tstRequireErrorMsg(t, "value -1 is negative", cut("uint5_neg"))
}

func TestObtainIntRangeValidator(t *testing.T) {
	configValues = make(map[string]string, 0)
	configValues["int1"] = "4832"
	configValues["int2_range"] = "-401"
	configValues["int3_parse"] = "hallo"
	configValues["int4_parse"] = "hallo 1234"
	configValues["int5"] = "-1"

	cut := ObtainIntRangeValidator(-400, 9200)

	require.Nil(t, cut("int1"))
	tstRequireErrorMsg(t, "value -401 is out of range [-400..9200]", cut("int2_range"))
	tstRequireErrorMsg(t, "value hallo is not a valid integer: ", cut("int3_parse"))
	tstRequireErrorMsg(t, "value hallo 1234 is not a valid integer: ", cut("int4_parse"))
	require.Nil(t, cut("int5"))
}

func TestObtainIsBooleanValidator(t *testing.T) {
	configValues = make(map[string]string, 0)
	configValues["bool1"] = "true"
	configValues["bool2"] = "false"
	configValues["bool3"] = "TRUE"
	configValues["bool4_fail"] = "hallo"

	cut := ObtainIsBooleanValidator()

	require.Nil(t, cut("bool1"))
	require.Nil(t, cut("bool2"))
	require.Nil(t, cut("bool3"))
	tstRequireErrorMsg(t, "value hallo is not a valid boolean value", cut("bool4_fail"))
}

func TestObtainIsRegexValidator(t *testing.T) {
	configValues = make(map[string]string, 0)
	configValues["regexp1"] = "^hello$"
	configValues["regexp2_fail"] = "^(hello$"

	cut := ObtainIsRegexValidator()

	require.Nil(t, cut("regexp1"))
	tstRequireErrorMsg(t, "value ^(hello$ is not a valid regex pattern", cut("regexp2_fail"))
}

func TestObtainSingleCharacterValidator(t *testing.T) {
	configValues = make(map[string]string, 0)
	configValues["char1"] = "a"
	configValues["char2_fail"] = ""
	configValues["char3_fail"] = "aa"

	cut := ObtainSingleCharacterValidator()

	require.Nil(t, cut("char1"))
	tstRequireErrorMsg(t, "cannot be empty", cut("char2_fail"))
	tstRequireErrorMsg(t, "cannot consist of multiple characters", cut("char3_fail"))
}
