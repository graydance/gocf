package gray_image

import "testing"

func TestResizeHeif(t *testing.T) {
	rawFile := "./testdata/raw.HEIC"
	dstFile := "./testdata/heic_out.jpg"
	err := ResizeHeif(rawFile, dstFile, 500, 100)
	if err != nil {
		t.Errorf("resize error: %v", err)
	}
}
