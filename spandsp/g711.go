package spandsp

const (
	/*! The A-law alternate mark inversion mask */
	G711_ALAW_AMI_MASK = 0x55

	/* The usual values to use on idle channels, to emulate silence */
	/*! Idle value for A-law channels */
	G711_ALAW_IDLE_OCTET = (0x80 ^ G711_ALAW_AMI_MASK)
	/*! Idle value for u-law channels */
	G711_ULAW_IDLE_OCTET = 0xFF
)

const (
	G711_ALAW = 0
	G711_ULAW = 1
)

type g711_state_t = g711_state_s

const (

	/*
	 * Mu-law is basically as follows:
	 *
	 *      Biased Linear Input Code        Compressed Code
	 *      ------------------------        ---------------
	 *      00000001wxyza                   000wxyz
	 *      0000001wxyzab                   001wxyz
	 *      000001wxyzabc                   010wxyz
	 *      00001wxyzabcd                   011wxyz
	 *      0001wxyzabcde                   100wxyz
	 *      001wxyzabcdef                   101wxyz
	 *      01wxyzabcdefg                   110wxyz
	 *      1wxyzabcdefgh                   111wxyz
	 *
	 * Each biased linear code has a leading 1 which identifies the segment
	 * number. The value of the segment number is equal to 7 minus the number
	 * of leading 0's. The quantization interval is directly available as the
	 * four bits wxyz.  * The trailing bits (a - h) are ignored.
	 *
	 * Ordinarily the complement of the resulting code word is used for
	 * transmission, and so the code word is complemented before it is returned.
	 *
	 * For further information see John C. Bellamy's Digital Telephony, 1982,
	 * John Wiley & Sons, pps 98-111 and 472-476.
	 */

	/* Enable the trap as per the MIL-STD */
	//#define G711_ULAW_ZEROTRAP
	/*! Bias for u-law encoding from linear. */
	G711_ULAW_BIAS = 0x84
)

// linear_to_ulaw
// Encode a linear sample to u-law
// param linear The sample to encode.
// return The u-law value.
func linear_to_ulaw(linear int32_t) uint8_t {
	var u_val uint8_t
	var mask int_t
	var seg int_t

	/* Get the sign and the magnitude of the value. */
	if linear >= 0 {
		linear = G711_ULAW_BIAS + linear
		mask = 0xFF
	} else {
		linear = G711_ULAW_BIAS - linear
		mask = 0x7F
	}

	seg = top_bit(uint32_t(linear|0xFF)) - 7
	if seg >= 8 {
		u_val = (uint8_t)(0x7F ^ mask)
	} else {
		/* Combine the sign, segment, quantization bits, and complement the code word. */
		u_val = (uint8_t)(((seg << 4) | ((int_t(linear) >> (seg + 3)) & 0xF)) ^ mask)
	}

	return u_val
}

// ulaw_to_linear
// Decode an u-law sample to a linear value.
// param ulaw The u-law sample to decode.
// return The linear value.
func ulaw_to_linear(ulaw uint8_t) int16_t {
	var t int32_t

	/* Complement to obtain normal u-law value. */
	// ulaw = ~ulaw;
	ulaw = ^ulaw

	/*
	 * Extract and bias the quantization bits. Then
	 * shift up by the segment number and subtract out the bias.
	 */
	t = int32_t(((ulaw&0x0F)<<3)+G711_ULAW_BIAS) << ((int32_t(ulaw) & 0x70) >> 4)
	// return (int16_t) ((ulaw & 0x80)  ?  (G711_ULAW_BIAS - t)  :  (t - G711_ULAW_BIAS));
	if ulaw&0x80 != 0 {
		return int16_t(G711_ULAW_BIAS - t)
	} else {
		return int16_t(t - G711_ULAW_BIAS)
	}
}

/*
 * A-law is basically as follows:
 *
 *      Linear Input Code        Compressed Code
 *      -----------------        ---------------
 *      0000000wxyza             000wxyz
 *      0000001wxyza             001wxyz
 *      000001wxyzab             010wxyz
 *      00001wxyzabc             011wxyz
 *      0001wxyzabcd             100wxyz
 *      001wxyzabcde             101wxyz
 *      01wxyzabcdef             110wxyz
 *      1wxyzabcdefg             111wxyz
 *
 * For further information see John C. Bellamy's Digital Telephony, 1982,
 * John Wiley & Sons, pps 98-111 and 472-476.
 */

// linear_to_alaw
// Encode a linear sample to A-law
// param linear The sample to encode.
// return The A-law value.
func linear_to_alaw(linear int32_t) uint8_t {
	var a_val uint8_t
	var mask int32_t
	var seg int32_t

	if linear >= 0 {
		/* Sign (bit 7) bit = 1 */
		mask = 0x80 | G711_ALAW_AMI_MASK
	} else {
		/* Sign (bit 7) bit = 0 */
		mask = G711_ALAW_AMI_MASK
		linear = -linear - 1
	}

	/* Convert the scaled magnitude to segment number. */
	seg = int32_t(top_bit(uint32_t(linear|0xFF)) - 7)
	if seg >= 8 {
		a_val = (uint8_t)(0x7F ^ mask)
	} else {
		/* Combine the sign, segment, and quantization bits. */
		// a_val = (uint8_t) (((seg << 4) | ((linear >> ((seg)  ?  (seg + 3)  :  4)) & 0x0F)) ^ mask);
		var value = seg
		if seg != 0 {
			value = seg + 3
		} else {
			value = 4
		}

		a_val = (uint8_t)(((seg << 4) | ((linear >> value) & 0x0F)) ^ mask)
	}
	return a_val
}

// alaw_to_linear Decode an A-law sample to a linear value.
// param alaw The A-law sample to decode.
// return The linear value.
func alaw_to_linear(alaw uint8_t) int16_t {
	var i int32_t
	var seg int32_t

	alaw ^= G711_ALAW_AMI_MASK
	i = (int32_t(alaw) & 0x0F) << 4
	seg = (int32_t(alaw) & 0x70) >> 4
	if seg != 0 {
		i = (i + 0x108) << (seg - 1)
	} else {
		i += 8
	}

	// return (int16_t)((alaw & 0x80)  ?  i: -i);
	if alaw&0x80 != 0 {
		return int16_t(i)
	} else {
		return int16_t(-i)
	}
}

/* Copied from the CCITT G.711 specification */
var ulaw_to_alaw_table = [256]uint8_t{
	42, 43, 40, 41, 46, 47, 44, 45, 34, 35, 32, 33, 38, 39, 36, 37,
	58, 59, 56, 57, 62, 63, 60, 61, 50, 51, 48, 49, 54, 55, 52, 53,
	10, 11, 8, 9, 14, 15, 12, 13, 2, 3, 0, 1, 6, 7, 4, 26,
	27, 24, 25, 30, 31, 28, 29, 18, 19, 16, 17, 22, 23, 20, 21, 106,
	104, 105, 110, 111, 108, 109, 98, 99, 96, 97, 102, 103, 100, 101, 122, 120,
	126, 127, 124, 125, 114, 115, 112, 113, 118, 119, 116, 117, 75, 73, 79, 77,
	66, 67, 64, 65, 70, 71, 68, 69, 90, 91, 88, 89, 94, 95, 92, 93,
	82, 82, 83, 83, 80, 80, 81, 81, 86, 86, 87, 87, 84, 84, 85, 85,
	170, 171, 168, 169, 174, 175, 172, 173, 162, 163, 160, 161, 166, 167, 164, 165,
	186, 187, 184, 185, 190, 191, 188, 189, 178, 179, 176, 177, 182, 183, 180, 181,
	138, 139, 136, 137, 142, 143, 140, 141, 130, 131, 128, 129, 134, 135, 132, 154,
	155, 152, 153, 158, 159, 156, 157, 146, 147, 144, 145, 150, 151, 148, 149, 234,
	232, 233, 238, 239, 236, 237, 226, 227, 224, 225, 230, 231, 228, 229, 250, 248,
	254, 255, 252, 253, 242, 243, 240, 241, 246, 247, 244, 245, 203, 201, 207, 205,
	194, 195, 192, 193, 198, 199, 196, 197, 218, 219, 216, 217, 222, 223, 220, 221,
	210, 210, 211, 211, 208, 208, 209, 209, 214, 214, 215, 215, 212, 212, 213, 213,
}

// These transcoding tables are copied from the CCITT G.711 specification. To achieve
// optimal results, do not change them.
var alaw_to_ulaw_table = [256]uint8_t{
	42, 43, 40, 41, 46, 47, 44, 45, 34, 35, 32, 33, 38, 39, 36, 37,
	57, 58, 55, 56, 61, 62, 59, 60, 49, 50, 47, 48, 53, 54, 51, 52,
	10, 11, 8, 9, 14, 15, 12, 13, 2, 3, 0, 1, 6, 7, 4, 5,
	26, 27, 24, 25, 30, 31, 28, 29, 18, 19, 16, 17, 22, 23, 20, 21,
	98, 99, 96, 97, 102, 103, 100, 101, 93, 93, 92, 92, 95, 95, 94, 94,
	116, 118, 112, 114, 124, 126, 120, 122, 106, 107, 104, 105, 110, 111, 108, 109,
	72, 73, 70, 71, 76, 77, 74, 75, 64, 65, 63, 63, 68, 69, 66, 67,
	86, 87, 84, 85, 90, 91, 88, 89, 79, 79, 78, 78, 82, 83, 80, 81,
	170, 171, 168, 169, 174, 175, 172, 173, 162, 163, 160, 161, 166, 167, 164, 165,
	185, 186, 183, 184, 189, 190, 187, 188, 177, 178, 175, 176, 181, 182, 179, 180,
	138, 139, 136, 137, 142, 143, 140, 141, 130, 131, 128, 129, 134, 135, 132, 133,
	154, 155, 152, 153, 158, 159, 156, 157, 146, 147, 144, 145, 150, 151, 148, 149,
	226, 227, 224, 225, 230, 231, 228, 229, 221, 221, 220, 220, 223, 223, 222, 222,
	244, 246, 240, 242, 252, 254, 248, 250, 234, 235, 232, 233, 238, 239, 236, 237,
	200, 201, 198, 199, 204, 205, 202, 203, 192, 193, 191, 191, 196, 197, 194, 195,
	214, 215, 212, 213, 218, 219, 216, 217, 207, 207, 206, 206, 210, 211, 208, 209,
}

func alaw_to_ulaw(alaw uint8_t) uint8_t {
	return alaw_to_ulaw_table[alaw]
}

func ulaw_to_alaw(ulaw uint8_t) uint8_t {
	return ulaw_to_alaw_table[ulaw]
}

// G.711 state
type g711_state_s struct {
	/*! One of the G.711_xxx options */
	mode int
}

func g711_init(mode int) *g711_state_s {
	return &g711_state_s{
		mode: mode,
	}
}

func (s *g711_state_s) g711_decode(g711_data []uint8_t) (amp []int16_t) {
	amp = make([]int16_t, len(g711_data))
	var i int
	switch s.mode {
	case G711_ALAW:
		for i = 0; i < len(g711_data); i++ {
			amp[i] = alaw_to_linear(g711_data[i])
		}
	case G711_ULAW:
		for i = 0; i < len(g711_data); i++ {
			amp[i] = ulaw_to_linear(g711_data[i])
		}
	}

	return amp
}

func (s *g711_state_s) g711_encode(amp []int16_t) (g711_data []uint8_t) {
	g711_data = make([]uint8_t, len(amp))
	var i int
	switch s.mode {
	case G711_ALAW:
		for i = 0; i < len(amp); i++ {
			g711_data[i] = linear_to_alaw(int32_t(amp[i]))
		}
	case G711_ULAW:
		for i = 0; i < len(amp); i++ {
			g711_data[i] = linear_to_ulaw(int32_t(amp[i]))
		}
	}

	return g711_data
}

func (s *g711_state_s) g711_transcode(g711_in []uint8_t) (g711_out []uint8_t) {
	var i int

	switch s.mode {
	case G711_ALAW:
		for i = 0; i < len(g711_in); i++ {
			g711_out[i] = alaw_to_ulaw_table[g711_in[i]]
		}
	case G711_ULAW:
		for i = 0; i < len(g711_in); i++ {
			g711_out[i] = ulaw_to_alaw_table[g711_in[i]]
		}
	}

	return g711_out
}
