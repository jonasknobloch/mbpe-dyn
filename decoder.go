package mbpe

type Decoder interface {
	Decode([]string) string
}
