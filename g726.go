package g726

import (
	"encoding/binary"
	"fmt"
)

type G726Rate int

const (
	G726Rate16kbps G726Rate = 0
	G726Rate24kbps G726Rate = 1
	G726Rate32kbps G726Rate = 2
	G726Rate40kbps G726Rate = 3
)

func (r G726Rate) String() string {
	switch r {
	case G726Rate16kbps:
		return "16kbps"
	case G726Rate24kbps:
		return "24kbps"
	case G726Rate32kbps:
		return "32kbps"
	case G726Rate40kbps:
		return "40kbps"
	default:
		return ""
	}
}

func (state_ptr *G726_state) Encode(pcm []int16) ([]byte, error) {
	switch state_ptr.rate {
	case G726Rate16kbps:
		input_len := len(pcm)
		if input_len%4 != 0 {
			return nil, fmt.Errorf("input length must be a multiple of 4 for 16kbps encoding")
		}

		var out = make([]byte, 0, input_len/4)
		for i := 0; i < input_len; i += 4 {
			a := state_ptr.g726_16_encoder(int(pcm[i+0]))
			b := state_ptr.g726_16_encoder(int(pcm[i+1]))
			c := state_ptr.g726_16_encoder(int(pcm[i+2]))
			d := state_ptr.g726_16_encoder(int(pcm[i+3]))

			// 4b -> 1b
			v := byte((a << 6) | (b << 4) | (c << 2) | d)
			out = append(out, v)
		}
		return out, nil
	case G726Rate24kbps:
		input_len := len(pcm)
		if input_len%8 != 0 {
			return nil, fmt.Errorf("input length must be multiple of 8 for 24kbps encoding")
		}

		out := make([]byte, 0, input_len/8*3)

		for i := 0; i < input_len; i += 8 {
			s0 := state_ptr.g726_24_encoder(int(pcm[i]))
			s1 := state_ptr.g726_24_encoder(int(pcm[i+1]))
			s2 := state_ptr.g726_24_encoder(int(pcm[i+2]))
			s3 := state_ptr.g726_24_encoder(int(pcm[i+3]))
			s4 := state_ptr.g726_24_encoder(int(pcm[i+4]))
			s5 := state_ptr.g726_24_encoder(int(pcm[i+5]))
			s6 := state_ptr.g726_24_encoder(int(pcm[i+6]))
			s7 := state_ptr.g726_24_encoder(int(pcm[i+7]))

			// 打包8个3位值到3个字节
			// 8b -> 3b
			b0 := byte(s0<<5) | byte(s1<<2) | byte(s2>>1)
			b1 := byte((s2&1)<<7) | byte(s3<<4) | byte(s4<<1) | byte(s5>>2)
			b2 := byte((s5&3)<<6) | byte(s6<<3) | byte(s7)

			out = append(out, b0, b1, b2)
		}
		return out, nil
	case G726Rate32kbps:
		input_len := len(pcm)
		if input_len%2 != 0 {
			return nil, fmt.Errorf("input length must be a multiple of 2 for 32kbps encoding")
		}

		var out = make([]byte, 0, input_len/2)
		for i := 0; i < input_len; i += 2 {
			a := state_ptr.g726_32_encoder(int(pcm[i+0]))
			b := state_ptr.g726_32_encoder(int(pcm[i+1]))

			// 2b -> 1b
			out = append(out, byte((a<<4)|b))
		}
		return out, nil
	case G726Rate40kbps:
		input_len := len(pcm)
		if input_len%8 != 0 {
			return nil, fmt.Errorf("input length must be multiple of 8 for 40kbps encoding")
		}

		out_len := input_len * 5 / 8
		out := make([]byte, 0, out_len)

		for i := 0; i < input_len; i += 8 {
			s0 := state_ptr.g726_40_encoder(int(pcm[i]))
			s1 := state_ptr.g726_40_encoder(int(pcm[i+1]))
			s2 := state_ptr.g726_40_encoder(int(pcm[i+2]))
			s3 := state_ptr.g726_40_encoder(int(pcm[i+3]))
			s4 := state_ptr.g726_40_encoder(int(pcm[i+4]))
			s5 := state_ptr.g726_40_encoder(int(pcm[i+5]))
			s6 := state_ptr.g726_40_encoder(int(pcm[i+6]))
			s7 := state_ptr.g726_40_encoder(int(pcm[i+7]))

			// 将8个5位值打包成5个字节
			b0 := byte((s0 << 3) | (s1 >> 2))
			b1 := byte(((s1 & 0x03) << 6) | (s2 << 1) | (s3 >> 4))
			b2 := byte(((s3 & 0x0F) << 4) | (s4 >> 1))
			b3 := byte(((s4 & 0x01) << 7) | (s5 << 2) | (s6 >> 3))
			b4 := byte(((s6 & 0x07) << 5) | s7)

			out = append(out, b0, b1, b2, b3, b4)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("invalid rate")
	}

}

func (state_ptr *G726_state) Decode(bitstream []byte) ([]int16, error) {
	switch state_ptr.rate {
	case G726Rate16kbps:
		input_len := len(bitstream)

		var out = make([]int16, 0, input_len*4)
		for i := 0; i < input_len; i++ {
			a := (bitstream[i] & byte(192)) >> 6
			b := (bitstream[i] & byte(48)) >> 4
			c := (bitstream[i] & byte(12)) >> 2
			d := (bitstream[i] & byte(3)) >> 0

			out = append(out, int16(state_ptr.g726_16_decoder(int(a))))
			out = append(out, int16(state_ptr.g726_16_decoder(int(b))))
			out = append(out, int16(state_ptr.g726_16_decoder(int(c))))
			out = append(out, int16(state_ptr.g726_16_decoder(int(d))))
		}
		return out, nil
	case G726Rate24kbps:
		input_len := len(bitstream)
		if input_len%3 != 0 {
			return nil, fmt.Errorf("input length must be multiple of 3 for 24kbps decoding")
		}

		out_len := input_len * 8 / 3
		out := make([]int16, 0, out_len)

		for i := 0; i < input_len; i += 3 {
			b0 := bitstream[i]
			b1 := bitstream[i+1]
			b2 := bitstream[i+2]

			s0 := (b0 & 0xE0) >> 5
			s1 := (b0 & 0x1C) >> 2
			s2 := ((b0 & 0x03) << 1) | ((b1 & 0x80) >> 7)
			s3 := (b1 & 0x70) >> 4
			s4 := (b1 & 0x0E) >> 1
			s5 := ((b1 & 0x01) << 2) | ((b2 & 0xC0) >> 6)
			s6 := (b2 & 0x38) >> 3
			s7 := (b2 & 0x07) >> 0

			out = append(out, int16(state_ptr.g726_24_decoder(int(s0))))
			out = append(out, int16(state_ptr.g726_24_decoder(int(s1))))
			out = append(out, int16(state_ptr.g726_24_decoder(int(s2))))
			out = append(out, int16(state_ptr.g726_24_decoder(int(s3))))
			out = append(out, int16(state_ptr.g726_24_decoder(int(s4))))
			out = append(out, int16(state_ptr.g726_24_decoder(int(s5))))
			out = append(out, int16(state_ptr.g726_24_decoder(int(s6))))
			out = append(out, int16(state_ptr.g726_24_decoder(int(s7))))
		}
		return out, nil
	case G726Rate32kbps:
		input_len := len(bitstream)

		var out = make([]int16, 0, input_len*2)
		for i := 0; i < input_len; i++ {
			a := (bitstream[i] & byte(240)) >> 4
			b := (bitstream[i] & byte(15)) >> 0

			out = append(out, int16(state_ptr.g726_32_decoder(int(a))))
			out = append(out, int16(state_ptr.g726_32_decoder(int(b))))
		}
		return out, nil
	case G726Rate40kbps:
		input_len := len(bitstream)
		if input_len%5 != 0 {
			return nil, fmt.Errorf("input length must be multiple of 5 for 40kbps decoding")
		}

		out_len := input_len * 8 / 5
		out := make([]int16, 0, out_len)

		for i := 0; i < input_len; i += 5 {
			b0 := bitstream[i]
			b1 := bitstream[i+1]
			b2 := bitstream[i+2]
			b3 := bitstream[i+3]
			b4 := bitstream[i+4]

			// 解包5个字节到8个5位值
			s0 := (b0 & 0xF8) >> 3
			s1 := ((b0 & 0x07) << 2) | ((b1 & 0xC0) >> 6)
			s2 := (b1 & 0x3E) >> 1
			s3 := ((b1 & 0x01) << 4) | ((b2 & 0xF0) >> 4)
			s4 := ((b2 & 0x0F) << 1) | ((b3 & 0x80) >> 7)
			s5 := (b3 & 0x7C) >> 2
			s6 := ((b3 & 0x03) << 3) | ((b4 & 0xE0) >> 5)
			s7 := (b4 & 0x1F) >> 0

			out = append(out, int16(state_ptr.g726_40_decoder(int(s0))))
			out = append(out, int16(state_ptr.g726_40_decoder(int(s1))))
			out = append(out, int16(state_ptr.g726_40_decoder(int(s2))))
			out = append(out, int16(state_ptr.g726_40_decoder(int(s3))))
			out = append(out, int16(state_ptr.g726_40_decoder(int(s4))))
			out = append(out, int16(state_ptr.g726_40_decoder(int(s5))))
			out = append(out, int16(state_ptr.g726_40_decoder(int(s6))))
			out = append(out, int16(state_ptr.g726_40_decoder(int(s7))))
		}
		return out, nil
	default:
		return nil, fmt.Errorf("invalid rate")
	}
}

func (state_ptr *G726_state) EncodeSimple(pcm []byte) ([]byte, error) {
	if len(pcm)%2 != 0 {
		return nil, fmt.Errorf("pcm length must be even")
	}

	pcm_in := make([]int16, len(pcm)/2)
	for i := 0; i < len(pcm_in); i++ {
		// 每2字节组合为一个int16
		pcm_in[i] = int16(binary.LittleEndian.Uint16(pcm[2*i : 2*i+2]))
	}

	return state_ptr.Encode(pcm_in)
}

func (state_ptr *G726_state) DecodeSimple(bitstream []byte) ([]byte, error) {
	pcm_out, err := state_ptr.Decode(bitstream)
	if err != nil {
		return nil, err
	}

	pcm := make([]byte, len(pcm_out)*2)
	for i := 0; i < len(pcm_out); i++ {
		binary.LittleEndian.PutUint16(pcm[2*i:2*i+2], uint16(pcm_out[i]))
	}
	return pcm, nil
}
