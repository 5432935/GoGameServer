package gateway

import (
	"github.com/funny/binary"
	"github.com/funny/link/packet"
)

//import (
//	. "tools"
//	"strconv"
//)

const (
	CMD_NEW_1 = 1
	CMD_NEW_2 = 2
	CMD_NEW_3 = 3
	CMD_DEL   = 4
	CMD_MSG   = 5
	CMD_BRD   = 6
	CMD_PING  = 7
	CMD_PONG  = 8
)

type gatewayMsg struct {
	Command   uint8
	ClientId  uint64
	ClientIds []uint64
	Data      []byte
	Message   interface{}
}

func (msg *gatewayMsg) Unmarshal(r *binary.Reader) error {
	msg.Command = r.ReadUint8()
	//	DEBUG("收到:" + strconv.Itoa(int(msg.Command)))
	switch msg.Command {
	case CMD_NEW_1:
		msg.ClientId = r.ReadUint64BE()
		msg.Data = r.ReadPacket(binary.SplitByUint8)
	case CMD_NEW_2:
		msg.ClientId = r.ReadUint64BE()
		msg.ClientIds = []uint64{r.ReadUint64BE()}
	case CMD_NEW_3:
		msg.ClientId = r.ReadUint64BE()
	case CMD_DEL:
		msg.ClientId = r.ReadUint64BE()
	case CMD_MSG:
		msg.ClientId = r.ReadUint64BE()
		msg.Data = r.ReadPacket(binary.SplitByUvarint)
	case CMD_BRD:
		num := int(r.ReadUvarint())
		msg.ClientIds = make([]uint64, num)
		for i := 0; i < num; i++ {
			msg.ClientIds[i] = r.ReadUvarint()
		}
		msg.Data = r.ReadPacket(binary.SplitByUvarint)
	}
	return nil
}

func (msg *gatewayMsg) Marshal(w *binary.Writer) error {
	w.WriteUint8(msg.Command)
	//	DEBUG("发送:" + strconv.Itoa(int(msg.Command)))
	switch msg.Command {
	case CMD_NEW_1:
		w.WriteUint64BE(msg.ClientId)
		w.WritePacket(msg.Data, binary.SplitByUint8)
	case CMD_NEW_2:
		w.WriteUint64BE(msg.ClientId)
		w.WriteUint64BE(msg.ClientIds[0])
	case CMD_NEW_3:
		w.WriteUint64BE(msg.ClientId)
	case CMD_DEL:
		w.WriteUint64BE(msg.ClientId)
	case CMD_MSG:
		w.WriteUint64BE(msg.ClientId)
		goto ENCODE
	case CMD_BRD:
		w.WriteUvarint(uint64(len(msg.ClientIds)))
		for i := 0; i < len(msg.ClientIds); i++ {
			w.WriteUvarint(msg.ClientIds[i])
		}
		goto ENCODE
	}
	return w.Flush()
ENCODE:
	if fast, ok := msg.Message.(packet.FastOutMessage); ok {
		binary.SplitByUvarint.WriteHead(w, fast.MarshalSize())
		if err := fast.Marshal(w); err != nil {
			return err
		}
	} else {
		data, err := msg.Message.(packet.OutMessage).Marshal()
		if err != nil {
			return err
		}
		w.WritePacket(data, binary.SplitByUvarint)
	}
	return w.Flush()
}
