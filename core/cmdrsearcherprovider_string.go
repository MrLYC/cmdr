// Code generated by "stringer -type=CmdrSearcherProvider"; DO NOT EDIT.

package core

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CmdrSearcherProviderUnknown-0]
	_ = x[CmdrSearcherProviderDefault-1]
	_ = x[CmdrSearcherProviderApi-2]
	_ = x[CmdrSearcherProviderAtom-3]
}

const _CmdrSearcherProvider_name = "CmdrSearcherProviderUnknownCmdrSearcherProviderDefaultCmdrSearcherProviderApiCmdrSearcherProviderAtom"

var _CmdrSearcherProvider_index = [...]uint8{0, 27, 54, 77, 101}

func (i CmdrSearcherProvider) String() string {
	if i < 0 || i >= CmdrSearcherProvider(len(_CmdrSearcherProvider_index)-1) {
		return "CmdrSearcherProvider(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CmdrSearcherProvider_name[_CmdrSearcherProvider_index[i]:_CmdrSearcherProvider_index[i+1]]
}
