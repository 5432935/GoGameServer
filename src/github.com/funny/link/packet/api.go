package packet

import (
	"encoding/json"
	"io"

	"github.com/funny/binary"
)

type InMessage interface {
	Unmarshal([]byte) error
}

type OutMessage interface {
	Marshal() ([]byte, error)
}

type FastInMessage interface {
	Unmarshal(r *io.LimitedReader) error
}

type FastOutMessage interface {
	MarshalSize() int
	Marshal(w *binary.Writer) error
}

type RAW []byte

func (msg RAW) MarshalSize() int {
	return len(msg)
}

func (msg RAW) Marshal(w *binary.Writer) error {
	_, err := w.Write(msg)
	return err
}

func (msg *RAW) Unmarshal(r *io.LimitedReader) error {
	if int64(cap(*msg)) >= r.N {
		*msg = (*msg)[0:r.N]
	} else {
		*msg = make([]byte, r.N)
	}
	_, err := io.ReadFull(r, *msg)
	return err
}

type JSON struct{ V interface{} }

func (msg JSON) Marshal() ([]byte, error) {
	return json.Marshal(msg.V)
}

func (msg JSON) Unmarshal(r *io.LimitedReader) error {
	return json.NewDecoder(r).Decode(msg.V)
}
