package g726

type g726_param_32 struct {
	qtab_721 [7]int
	_dqlntab [16]int
	_witab   [16]int
	_fitab   [16]int
}

var p32 g726_param_32

func init() {
	p32.qtab_721 = [7]int{-124, 80, 178, 246, 300, 349, 400}
	/*
	 * Maps G.721 code word to reconstructed scale factor normalized log
	 * magnitude values.
	 */
	p32._dqlntab = [16]int{-2048, 4, 135, 213, 273, 323, 373, 425,
		425, 373, 323, 273, 213, 135, 4, -2048}

	/* Maps G.721 code word to log of scale factor multiplier. */
	p32._witab = [16]int{-12, 18, 41, 64, 112, 198, 355, 1122,
		1122, 355, 198, 112, 64, 41, 18, -12}
	/*
	 * Maps G.721 code words to a set of values whose long and short
	 * term averages are computed and then compared to give an indication
	 * how stationary (steady state) the signal is.
	 */
	p32._fitab = [16]int{0, 0, 0, 0x200, 0x200, 0x200, 0x600, 0xE00,
		0xE00, 0x600, 0x200, 0x200, 0x200, 0, 0, 0}
}

func (state_ptr *G726_state) g726_32_encoder(sl int) int {
	var (
		sezi  int
		sez   int
		se    int
		d     int
		y     int
		i     int
		dq    int
		sr    int
		dqsez int
	)

	sl >>= 2 /* 14-bit dynamic range */

	sezi = state_ptr.predictor_zero()
	sez = sezi >> 1
	se = (sezi + state_ptr.predictor_pole()) >> 1 /* estimated signal */

	d = sl - se /* estimation difference */

	/* quantize the prediction difference */
	y = state_ptr.step_size()           /* quantizer step size */
	i = quantize(d, y, p32.qtab_721[:]) /* i = ADPCM code */

	dq = reconstruct(i&8, p32._dqlntab[i], y) /* quantized est diff */

	sr = IfElse[int](dq < 0, se-(dq&0x3FFF), se+dq) /* reconst. signal */

	dqsez = sr + sez - se /* pole prediction diff. */

	t0 := int(p32._witab[i])
	t1 := t0 << 5
	t2 := int(t1)
	_ = t2
	state_ptr.update(4, y, p32._witab[i]<<5, p32._fitab[i], dq, sr, dqsez)

	return i
}

func (state_ptr *G726_state) g726_32_decoder(i int) int {
	var (
		sezi  int
		sez   int
		sei   int
		se    int
		y     int
		dq    int
		sr    int
		dqsez int
		lino  int
	)

	i &= 0x0f /* mask to get proper bits */
	sezi = state_ptr.predictor_zero()
	sez = sezi >> 1
	sei = sezi + state_ptr.predictor_pole()
	se = sei >> 1 /* se = estimated signal */

	y = state_ptr.step_size() /* dynamic quantizer step size */

	dq = reconstruct(i&0x08, p32._dqlntab[i], y) /* quantized diff. */

	sr = IfElse[int](dq < 0, se-(dq&0x3FFF), se+dq) /* reconst. signal */

	dqsez = sr - se + sez /* pole prediction diff. */

	state_ptr.update(4, y, p32._witab[i]<<5, p32._fitab[i], dq, sr, dqsez)

	lino = sr << 2 /* this seems to overflow a short*/
	lino = IfElse[int](lino > 32767, 32767, lino)
	lino = IfElse[int](lino < -32768, -32768, lino)

	return lino //(sr << 2);	/* sr was 14-bit dynamic range */
}
