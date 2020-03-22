package opus

import (
	"errors"
	"io"
	"math"

	"github.com/faiface/beep"
	codec "gopkg.in/hraban/opus.v2"
)

func EncodeBuffer(buffer *beep.Buffer, w io.Writer) (int, error) {
	buff := make([]byte, 512)
	encoder, err := NewEncoder(buffer, Config{})
	if err != nil {
		return 0, err
	}

	n := 0
	for {
		s, err := encoder.Read(buff)
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
		s, err = w.Write(buff[:s])
		if err != nil {
			return n, err
		}
		n += s
	}

	return n, nil
}

var supportedLatencies = []float64{
	2.5,
	5,
	10,
	20,
	40,
	60,
}

type encoder struct {
	buffer *beep.Buffer

	enc *codec.Encoder
	tmp [][2]float64
	pos int
}

func NewEncoder(buffer *beep.Buffer, config Config) (io.ReadCloser, error) {
	if buffer == nil {
		return nil, errors.New("buffer is nil")
	}

	if config.Latency == 0 {
		config.Latency = 2
	}

	if config.BitRate == 0 {
		config.BitRate = 32000
	}

	if config.Optimize == 0 {
		config.Optimize = codec.AppRestrictedLowdelay
	}

	// Select the nearest supported latency
	var targetLatency float64
	latencyInMS := float64(config.Latency.Milliseconds())
	nearestDist := math.Inf(+1)
	for _, latency := range supportedLatencies {
		dist := math.Abs(latency - latencyInMS)
		if dist >= nearestDist {
			break
		}

		nearestDist = dist
		targetLatency = latency
	}

	sampleRate := int(buffer.Format().SampleRate)
	bufferSamplesSize := int(targetLatency * float64(sampleRate) / 1000)
	bufferSamples := make([][2]float64, bufferSamplesSize)

	numChannels := buffer.Format().NumChannels
	enc, err := codec.NewEncoder(sampleRate, numChannels, config.Optimize)
	if err != nil {
		return nil, err
	}

	err = enc.SetBitrate(config.BitRate)
	if err != nil {
		return nil, err
	}

	return &encoder{
		buffer,
		enc,
		bufferSamples,
		0,
	}, nil
}

func (e *encoder) Read(p []byte) (n int, err error) {
	if e.pos >= e.buffer.Len() {
		return 0, io.EOF
	}

	to := e.pos + len(e.tmp)
	if to > e.buffer.Len() {
		to = e.buffer.Len()
	}
	s := e.buffer.Streamer(e.pos, to)
	n, ok := s.Stream(e.tmp)
	if !ok {
		return 0, wrapError(s.Err())
	}
	e.pos = to

	n, err = e.enc.EncodeFloat32(flattenSamples(e.tmp), p)
	if err != nil {
		return n, wrapError(err)
	}

	return n, nil
}

func (e *encoder) Close() error {
	return nil
}
