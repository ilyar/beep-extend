package opus

import (
	"github.com/pkg/errors"
	codec "gopkg.in/hraban/opus.v2"
)

func wrapError(err error) error {
	switch err.(type) {
	default:
		err = errors.Wrap(err, "opus")
	case codec.Error:
		err = codecError(err.(codec.Error))
	case nil:
	}

	return err
}

func codecError(err codec.Error) error {
	return errors.New(err.Error())
}

func flattenSamples(samples [][2]float64) []float32 {
	if len(samples) == 0 {
		return make([]float32, 0)
	}

	pcm := make([]float32, 2*len(samples))
	offset := 0
	for i, sample := range samples {
		pcm[i+offset] = float32(sample[0])
		offset++
		pcm[i+offset] = float32(sample[1])
	}

	return pcm
}
