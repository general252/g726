package spandsp

import (
	"fmt"
	"testing"
)

func Test_g726_init(t *testing.T) {
	var pcm_in = make([]int16, 480)
	for i := 0; i < len(pcm_in); i++ {
		pcm_in[i] = -0x7800
	}

	s_e, _ := G726_init(32000, G726_ENCODING_LINEAR, G726_PACKING_LEFT)
	s_d, _ := G726_init(32000, G726_ENCODING_LINEAR, G726_PACKING_LEFT)

	bitstream := s_e.Encode(pcm_in)
	pcm_out := s_d.Decode(bitstream)

	for i := 0; i < len(pcm_out); i++ {
		cin := int(pcm_in[i])
		out := int(pcm_out[i])
		diff := cin - out
		if diff < 0 {
			diff = -diff
		}

		fmt.Printf("%04d: [%04x - %04x: %d]\n", i, cin&0xFFFF, out&0xFFFF, diff&0xFFFF)
	}
	fmt.Println()
}
