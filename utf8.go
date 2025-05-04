// Portions of this file are derived from https://github.com/sugarme/tokenizer,
// licensed under the Apache License, Version 2.0.
// See https://www.apache.org/licenses/LICENSE-2.0

package mbpe

var BytesChar = GenerateBytesChar()

var CharBytes = func() map[string]uint8 {
	var bc = GenerateBytesChar()
	var cb = make(map[string]uint8)

	for b, c := range bc {
		cb[c] = b
	}

	return cb
}()

func encodeUTF8(r rune) []byte {
	const (
		// first byte of a 2-byte encoding starts 110 and carries 5 bits of data
		b2Lead = 0xC0 // 1100 0000
		b2Mask = 0x1F // 0001 1111

		// first byte of a 3-byte encoding starts 1110 and carries 4 bits of data
		b3Lead = 0xE0 // 1110 0000
		b3Mask = 0x0F // 0000 1111

		// first byte of a 4-byte encoding starts 11110 and carries 3 bits of data
		b4Lead = 0xF0 // 1111 0000
		b4Mask = 0x07 // 0000 0111

		// non-first bytes start 10 and carry 6 bits of data
		mbLead = 0x80 // 1000 0000
		mbMask = 0x3F // 0011 1111
	)

	switch i := uint32(r); {
	case i <= 1<<7-1: // single byte
		return []byte{byte(r)}
	case i <= 1<<11-1: // two bytes
		return []byte{
			b2Lead | byte(r>>6),
			mbLead | byte(r)&mbMask}
	case i <= 1<<16-1: // three bytes
		return []byte{
			b3Lead | byte(r>>12),
			mbLead | byte(r>>6)&mbMask,
			mbLead | byte(r)&mbMask}
	default:
		return []byte{
			b4Lead | byte(r>>18),
			mbLead | byte(r>>12)&mbMask,
			mbLead | byte(r>>6)&mbMask,
			mbLead | byte(r)&mbMask}
	}
}

func GenerateBytesChar() map[uint8]string {
	var bc map[uint8]string = make(map[uint8]string)

	n := 0

	// 0 ('Ā') - 32 ('Ġ') - control codes
	for i := 256; i <= 288; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 33 ('!') - 126 ('~')
	for i := 33; i <= 126; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 127 ('ġ') - 160 ('ł') - control codes
	for i := 289; i <= 322; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 161 ('¡') - 172 ('¬')
	for i := 161; i <= 172; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 173 - ('Ń') - control code
	if n == 173 {
		r := rune(323)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 174 ('®') - 255 ('ÿ')
	for i := 174; i <= 255; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	return bc
}

func Alphabet() []string {
	alphabet := make([]string, 256)

	for i, c := range BytesChar {
		alphabet[i] = c
	}

	return alphabet
}
