package utils

import (
	"demo/chatRoom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//这里将这些方法关联到结构体中
type Transfer struct {
	//分析他应该有什么字段
	Conn net.Conn
	Buf  [8096]byte //这是传输时使用的缓冲
}

//写包
func (this *Transfer) WritePkg(data []byte) (err error) {

	//先把data转成[]byte
	var pkgLen uint32
	//获得data长度
	pkgLen = uint32(len(data))
	//将pkhLen转成[]byte
	binary.BigEndian.PutUint32(this.Buf[:4], pkgLen)
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	n, err = this.Conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	return
}

//读包
func (this *Transfer) ReadPkg() (mes message.Message, err error) {

	n, err := this.Conn.Read(this.Buf[:4])
	if n == 0 {
		fmt.Println("服务器检测到客户端已关闭，断开链接！")
		return
	}
	if err != nil {
		fmt.Println("ReadPkg() Error = ", err)
		return
	}
	//根据buf[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4]) //pkgLen变成了之前的接受的数据的长度

	//根据pkgLen读取消息内容
	n, err = this.Conn.Read(this.Buf[:pkgLen]) //把pkgLen长度的消息读入buf中
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Read(buf[:pkgLen]) err =", err)
	}
	//把buf反序列化
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err =", err)
		return
	}
	//把mes反序列化，得到LoginMes结构体

	return
}
