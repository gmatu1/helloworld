package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	flag "github.com/spf13/pflag"
)

const (
	minSpeedPercentage = 10
	maxSpeedPercentage = 500
)

var mp3File = "./audio/hallo_welt.mp3"
var speedPercentage = 100

func between(min, v, max int) int {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}

func play(audio *beep.Resampler) {
	done := make(chan bool)
	speaker.Play(
		beep.Seq(audio, beep.Callback(func() {
			done <- true
		})))
	<-done
}

func main() {
	flag.StringVarP(&mp3File, "file", "f", mp3File, "mp3 file to play")
	flag.IntVarP(&speedPercentage, "speed", "s", speedPercentage, "faster or slower (as percentage)")
	flag.Parse()

	speedPercentage = between(minSpeedPercentage, speedPercentage, maxSpeedPercentage)
	fmt.Printf("Play %s with %d%% speed\n", mp3File, speedPercentage)

	f, err := os.Open(mp3File)
	if err != nil {
		panic(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		panic(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))
	resampled := beep.ResampleRatio(6, float64(speedPercentage)/100, effects.Mono(streamer))

	fmt.Printf("Start playing...")
	play(resampled)
	fmt.Println("done")
}
