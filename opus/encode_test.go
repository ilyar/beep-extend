package opus_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ilyar/beep-extend/mp3"
	"github.com/ilyar/beep-extend/opus"
)

func TestEncodeBuffer(t *testing.T) {
	tmp, err := mp3.NewBufferFromMP3("../testdata/simple.mp3")

	assert.Nil(t, err)

	var out bytes.Buffer
	n, err := opus.EncodeBuffer(tmp, &out)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, n, 1770)

	//err = ioutil.WriteFile("../testdata/simple.opus.dat", out.Bytes(), 0644)
	//assert.Nil(t, err)
}
