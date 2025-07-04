package main

import (
	"fmt"
	"log"

	"github.com/general252/g726"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var pcm_in = make([]int16, 480)
	for i := 0; i < len(pcm_in); i++ {
		pcm_in[i] = -0x7800
	}

	var rate = g726.G726Rate32kbps

	encoder := g726.G726_init_state(rate)
	bitstream, err := encoder.Encode(pcm_in)
	if err != nil {
		fmt.Println(err)
		return
	}

	decoder := g726.G726_init_state(rate)
	pcm_out, err := decoder.Decode(bitstream)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(pcm_out); i++ {
		cin := int(pcm_in[i])
		out := int(pcm_out[i])
		diff := g726.ABS(cin - out)

		fmt.Printf("%04d: [%04x - %04x: %d]\n", i, cin&0xFFFF, out&0xFFFF, diff&0xFFFF)
	}
}
