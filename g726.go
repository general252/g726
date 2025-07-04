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

func (state_ptr *G726_state) Encode(pcm []int16) ([]byte, error) {
	switch state_ptr.rate {
	case G726Rate16kbps:
		input_len := len(pcm)
		var out = make([]byte, 0, input_len/4)
		for i := 0; i < input_len/4; i++ {
			a := state_ptr.g726_16_encoder(int(pcm[i*4]))
			b := state_ptr.g726_16_encoder(int(pcm[i*4+1]))
			c := state_ptr.g726_16_encoder(int(pcm[i*4+2]))
			d := state_ptr.g726_16_encoder(int(pcm[i*4+3]))

			v := byte((a << 6) | (b << 4) | (c << 2) | d)
			out = append(out, v)
		}
		return out, nil
	case G726Rate24kbps:
	case G726Rate32kbps:
		input_len := len(pcm)
		var out = make([]byte, 0, input_len/2)
		for i := 0; i < input_len/2; i++ {
			a := state_ptr.g726_32_encoder(int(pcm[i*2]))
			b := state_ptr.g726_32_encoder(int(pcm[i*2+1]))

			out = append(out, byte((a<<4)|b))
		}
		return out, nil
	case G726Rate40kbps:
	default:
		return nil, fmt.Errorf("invalid rate")
	}

	return nil, fmt.Errorf("not implemented")
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
	default:
		return nil, fmt.Errorf("invalid rate")
	}

	return nil, fmt.Errorf("not implemented")
}

func (state_ptr *G726_state) EncodeSimple(pcm []byte) ([]byte, error) {
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
