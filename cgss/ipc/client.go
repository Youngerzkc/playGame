package ipc
import (
	"encoding/json"
)
type IpcClient struct{
	coon chan string
}
func NewIpcClient(server *IpcServer) *IpcClient {
	c:=server.Connect()
	return &IpcClient{c}
}
func (client *IpcClient)Call(method ,params string)(resp *Response,err error)  {
	req:=&Request{method,params}
	var b []byte 
	b,err =json.Marshal(req)
	if err!=nil{
	return 
}
	client.coon<-string(b)
	str:=<-client.coon
	var respl Response
	err =json.Unmarshal([]byte(str), &respl)
	return
}
func (client *IpcClient)Close () {
	client.coon<-"CLOSE"
}