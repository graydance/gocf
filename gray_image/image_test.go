package gray_image

import (
	"testing"
)

func Test_ResizeProportional(t *testing.T) {
	rawFile := "./testdata/big.jpg"
	dstFile := "./testdata/out.jpg"
	err := ResizeProportional(rawFile, dstFile, 500, 100)
	if err != nil {
		t.Errorf("resize error: %v", err)
	}
}
