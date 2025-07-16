package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lkmio/g726"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	pcmIn, err := os.ReadFile("./example/audio-samples.pcm")
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i <= 3; i++ {
		var rate = g726.G726Rate(i)

		pcmOut, g726Data, err := encodeAndDecode(rate, pcmIn)
		if err != nil {
			log.Println(err)
			return
		} else {
			filenamePCM := fmt.Sprintf("./example/audio-samples-re-%vkbps.pcm", (i+2)*8)
			filenameG726 := fmt.Sprintf("./example/audio-samples-%vkbps.g726", (i+2)*8)
			log.Println(rate, filenamePCM, filenameG726)

			_ = os.WriteFile(filenamePCM, pcmOut, 0644)
			_ = os.WriteFile(filenameG726, g726Data, 0644)

		}
	}
}

func encodeAndDecode(rate g726.G726Rate, pcm []byte) (pcmOut, g726Data []byte, err error) {
	var s int
	switch rate {
	case g726.G726Rate16kbps:
		s = 4
	case g726.G726Rate24kbps:
		s = 8
	case g726.G726Rate32kbps:
		s = 2
	case g726.G726Rate40kbps:
		s = 8
	default:
		return nil, nil, fmt.Errorf("invalid rate")
	}

	s *= 2
	pcm = pcm[:len(pcm)/s*s]

	encoder := g726.G726_init_state(rate)
	g726Data, err = encoder.EncodeSimple(pcm)
	if err != nil {
		log.Println(err, rate, len(pcm))
		return nil, nil, err
	}

	decoder := g726.G726_init_state(rate)
	out, err := decoder.DecodeSimple(g726Data)
	if err != nil {
		log.Println(err, rate, len(g726Data))
		return nil, nil, err
	}

	return out, g726Data, nil
}
