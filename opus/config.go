package opus

import (
	"time"

	codec "gopkg.in/hraban/opus.v2"
)

type Config struct {
	Latency  time.Duration
	BitRate  int
	Optimize codec.Application
}
