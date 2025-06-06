package mbpe

type ByteLevel struct {
	addPrefixSpace bool
	fsa            *FSA
}

type PreTokenizer interface {
	PreTokenize(string) []string
}

func NewByteLevel(addPrefixSpace bool) *ByteLevel {
	return &ByteLevel{
		addPrefixSpace: addPrefixSpace,
		fsa:            NewFSA(),
	}
}

func (p *ByteLevel) PreTokenize(phrase string) []string {
	if phrase == "" {
		return []string{}
	}

	if p.addPrefixSpace && phrase[0] != ' ' {
		phrase = " " + phrase
	}

	compounds := p.fsa.FindAll(phrase)

	for i, compound := range compounds {
		r := ""

		for _, b := range []byte(compound) {
			r += BytesChar[b]
		}

		compounds[i] = r
	}

	return compounds
}

func (p *ByteLevel) Decode(tokens []string) string {
	phrase := ""

	for _, token := range tokens {
		r := make([]byte, 0)

		for _, c := range token {
			r = append(r, CharBytes[string(c)])
		}

		phrase += string(r)
	}

	return phrase
}
