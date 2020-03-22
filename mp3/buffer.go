package mp3

import (
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
)

func NewBufferFromMP3(filePath string) (*beep.Buffer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return nil, err
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	err = streamer.Close()
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
