package main

import (
	"fmt"
	"log"
	"os"

	"github.com/general252/g726"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	pcmIn, err := os.ReadFile("audio-samples.pcm")
	if err != nil {
		fmt.Println(err)
		return
	}

	pcmOut, err := encodeAndDecode(g726.G726Rate32kbps, pcmIn)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		_ = os.WriteFile("audio-samples-re-32kbps.pcm", pcmOut, 0644)
	}

	pcmOut, err = encodeAndDecode(g726.G726Rate16kbps, pcmIn)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		_ = os.WriteFile("audio-samples-re-16kbps.pcm", pcmOut, 0644)
	}
}

func encodeAndDecode(rate g726.G726Rate, pcm []byte) ([]byte, error) {
	encoder := g726.G726_init_state(rate)
	g726Data, err := encoder.EncodeSimple(pcm)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	decoder := g726.G726_init_state(rate)
	out, err := decoder.DecodeSimple(g726Data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return out, nil
}
