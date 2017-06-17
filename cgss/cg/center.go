package cg

import (
	"sync"
	"errors"
	"~/gopath/src/cgss/ipc"
	"encoding/json"
)
type Message struct {
	From string `json:"from"`
	To string `json:"to"`
	Content string `json:"content"`
}
type CenterServer struct{
	servers map[string] ipc.Server 
	players []*Player
	rooms  []*Room 
	mutex  sync.RWMutex
}
func NewCenterServer() *CenterServer{
	servers:=make(map[string] ipc.Server)
	players:=make([]*Player,0)
	return&CenterServer{servers,players}
}
func (server *CenterServer)addPlayer(params string)error  {
	player:=NewPlayer()
	err:=json.Unmarshal([]byte(params),&player)
	if err !=nil{
		return err
	}
	server.mutex.Lock()
	defer server.mutex.Unlock()
	// 登录检查
	server.players=append(server.players,player)
	return nil
}
func (server *CenterServer)removePlayer(params string)error{
	server.mutex.Lock()
	defer server.mutex.Unlock()
	for i,v:=range server.players{
		if v.Name ==params{
			if len(server.players)==1{
				server.players=make([]*Player, 0)
			} else if i==len(server.players)-1{
				server.players=server.players[:1]
			}else if i==0{
				server.players=server.players[1:]
			}else{
				server.players=append(server.players[:i-1],server.players[:i+1] ...)
			}
			return nil
		}
	}
	return errors.New("Player not found.")
}
func (server *CenterServer)listPlayer(params string)(players string,err error)  {
	server.mutex.RLock()
	defer server.mutex.RUnlock()
	if len(server.players)>0{
		b,_:=json.Marshal(server.players)
		players=string(b)
	}else{
		err:=errors.New("No player online")
	}
	return
}
func (server *CenterServer)broadcast(params string) error  {
	var message Message
	err:=json.Marshal([]byte(params),&message)
	if err !=nil{
		return err
	}
	server.mutex.Lock()
	defer server.mutex.Unlock()
	if len(server.players)>0{
		for _,player:=range server.players{
			player.mq<-&message
		}
	}else{
		err =errors.New("No play online")
	}
	return err
}
func (server *CenterServer)Handle(method,params string)*ipc.Response {
	switch method {
	case "addplayer":
		err :=server.addPlayer(params)
		if err!=nil{
		// 
			return
		}
		
	case "removeplayer":
		err:=server.removePlayer(params)
		if err !=nil{
			return
		}
	case "listplayer":
		err:=server.listPlayer(params)
		if err !=nil{
			return
		}
		return &ipc.Response{"200",players}
	case "broadcast":
		err:=server.broadcast(params)
		if err !=nil{
			return
		}
		return &ipc.Response{Code:"200"}
	default:
		return &ipc.Response{Code:"404"}	
	}
	return &ipc.Response{Code:"200"}
}
func (server*CenterServer)Name()string{
	return "CenterServer"
}