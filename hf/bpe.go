package hf

import (
	"encoding/json"
)

type BPE struct {
	Type                    string          `json:"type"`
	Dropout                 json.RawMessage `json:"dropout"`
	UnkToken                json.RawMessage `json:"unk_token"`
	ContinuingSubwordPrefix json.RawMessage `json:"continuing_subword_prefix"`
	EndOfWordSuffix         json.RawMessage `json:"end_of_word_suffix"`
	FuseUnk                 bool            `json:"fuse_unk"`
	ByteFallback            bool            `json:"byte_fallback"`
	Vocab                   Vocab           `json:"vocab"`
	Merges                  Merges          `json:"merges"`
}
