package spandsp

import (
	"fmt"
	"math/bits"
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

func Test_alaw_to_linear(t *testing.T) {
	var alaw2lpcm = make([]int16, 0, 256)
	for i := 0; i <= 255; i++ {
		var linear = alaw_to_linear(uint8_t(i))
		alaw2lpcm = append(alaw2lpcm, linear)
	}

	t.Log(alaw2lpcm)

	var count int
	for i := -32768; i <= 32767; i++ {
		v := int16(i)
		a := encodeAlawFrame(v)
		b := linear_to_alaw(v)

		if a != b {
			t.Logf("%v %v %v", v, a, b)
		} else {
			count++
		}
	}
	t.Logf("success count: %v", count)
}

func encodeAlawFrame(frame int16) uint8 {
	var compressedByte, seg, sign int16
	sign = ((^frame) >> 8) & 0x80
	if sign == 0 {
		frame = ^frame
	}
	compressedByte = frame >> 4
	if compressedByte > 15 {
		seg = int16(12 - bits.LeadingZeros16(uint16(compressedByte)))
		compressedByte >>= seg - 1
		compressedByte -= 16
		compressedByte += seg << 4
	}
	return uint8((sign | compressedByte) ^ 0x0055)
}
