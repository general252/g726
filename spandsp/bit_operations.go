package spandsp

// top_bit Find the bit position of the highest set bit in a word
// param bits The word to be searched
// return The bit number of the highest set bit, or -1 if the word is zero.
func top_bit(bits uint32_t) int_t {
	/* Visual Studio x86_64 */
	/* TODO: Need the appropriate x86_64 code */
	var res int_t

	if bits == 0 {
		return -1
	}
	res = 0
	if (bits & 0xFFFF0000) != 0 {
		bits &= 0xFFFF0000
		res += 16
	}
	if (bits & 0xFF00FF00) != 0 {
		bits &= 0xFF00FF00
		res += 8
	}
	if (bits & 0xF0F0F0F0) != 0 {
		bits &= 0xF0F0F0F0
		res += 4
	}
	if (bits & 0xCCCCCCCC) != 0 {
		bits &= 0xCCCCCCCC
		res += 2
	}
	if (bits & 0xAAAAAAAA) != 0 {
		bits &= 0xAAAAAAAA
		res += 1
	}

	return res
}
