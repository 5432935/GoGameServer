package worldProxy

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"global"
	"protos"
	"protos/gameProto"
	"protos/systemProto"
	. "tools"
)

var (
	session *link.Session
)

//初始化
func InitClient(ip string, port string) error {
	addr := ip + ":" + port
	client, err := link.Connect("tcp", addr, packet.New(binary.SplitByUint32BE, 1024, 1024, 1024))
	if err != nil {
		return err
	}

	session = client
	go dealReceiveMsg()
	ConnectWorldServer()

	return nil
}

//处理从TransferServer发回的消息
func dealReceiveMsg() {
	for {
		var msg packet.RAW
		if err := session.Receive(&msg); err != nil {
			break
		}
		dealReceiveMsgS2C(msg)
	}
}

//处理接收到的系统消息
func dealReceiveSystemMsgS2C(msg packet.RAW) {
	protoMsg := systemProto.UnmarshalProtoMsg(msg)
	if protoMsg == systemProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case systemProto.ID_System_ConnectWorldServerS2C:
		connectWorldServerCallBack(protoMsg)
	}
}

//处理接收到的消息
func dealReceiveMsgS2C(msg packet.RAW) {
	if len(msg) < 2 {
		return
	}

	msgID := binary.GetUint16LE(msg[:2])
	//	DEBUG(global.ServerName, msgID)
	if systemProto.IsValidID(msgID) {
		dealReceiveSystemMsgS2C(msg)
	} else if gameProto.IsValidID(msgID) {
		if msgID%2 == 0 {
			//S2C消息，直接发送到用户客户端
			msgIdentification := binary.GetUint64LE(msg[2:10])
			userSession := global.GetSession(msgIdentification)
			if userSession == nil {
				return
			}
			protos.Send(msg, userSession)
		}
	} else {
		ERR(global.ServerName, "收到未处理消息")
	}
}

//发送系统消息到WorldServer
func SendSystemMsgToServer(msg []byte) {
	if session == nil {
		return
	}
	protos.Send(msg, session)
}

//发送游戏消息到WorldServer
func SendGameMsgToServer(msg []byte) {
	if session == nil {
		return
	}
	protos.Send(msg, session)
}

//发送连接WorldServer
func ConnectWorldServer() {
	INFO(global.ServerName + " Connect WorldServer ...")
	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectWorldServerC2S{
		ServerName: protos.String(global.ServerName),
		ServerID:   protos.Uint32(global.ServerID),
	})
	SendSystemMsgToServer(send_msg)
}

//连接Transfer服务器返回
func connectWorldServerCallBack(protoMsg systemProto.ProtoMsg) {
	//	rev_msg := protoMsg.Body.(*systemProto.System_ConnectWorldServerS2C)
	INFO(global.ServerName + " Connect WorldServer Success")
}