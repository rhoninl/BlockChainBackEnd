package Model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"runtime"
	"sync"
	"time"
)

type Client struct {
	clients     map[int64]*websocket.Conn     //公司Id与链接的map
	clientGroup map[string]map[int64]struct{} // 通过公司类型对链接进行分类保存
	clientNum   int                           //记录当前用户总数，便于维护心跳检测
}

var myClient Client

func init() {
	myClient = Client{
		clients:     make(map[int64]*websocket.Conn),
		clientGroup: make(map[string]map[int64]struct{}),
		clientNum:   0,
	}
	myClient.clientGroup["未选择"] = make(map[int64]struct{})
	myClient.clientGroup["船代"] = make(map[int64]struct{})
	myClient.clientGroup["货代"] = make(map[int64]struct{})
	myClient.clientGroup["陆运公司"] = make(map[int64]struct{})
	go myClient.HeartBeat()
	fmt.Println("[Websocket]init Success!")
}

func (c *Client) Login(conn *websocket.Conn, id int64) {
	_, exists := c.clients[id]
	if exists {
		c.SendMessageToId(gin.H{"message": "您的帐号在另一个地方登陆，寄"}, id)
		c.UnRegister(id)
	}
	c.clients[id] = conn
	_, companyType := GetCompanyBasicInfo(id)
	c.clientGroup[companyType][id] = struct{}{}
	c.clientNum++
	num, _ := GetUnReadNum(id)
	c.SendMessageToId(gin.H{
		"UnReadMessageNum": num,
		"OnlineNum":        c.clientNum,
	}, id)
}

func (c *Client) SendMessageToId(message interface{}, id int64) {
	con, exists := c.clients[id]
	if exists {
		con.WriteJSON(message)
	}
}

func (c *Client) SendMessageToGroup(message interface{}, group string) {
	wg := sync.WaitGroup{}
	for item := range c.clientGroup[group] {
		wg.Add(1)
		go func(id int64) {
			c.SendMessageToId(message, id)
			wg.Done()
		}(item)
	}
	wg.Wait()
}

func (c *Client) HeartBeat() {
	ch := make(chan struct{}, 20) // 同时校验上限
	wg := sync.WaitGroup{}
	for {
		select {
		case <-time.After(time.Minute * 5):
			for k, v := range c.clients {
				wg.Add(1)
				ch <- struct{}{} //写入信息到channel用于计数
				go func(id int64, conn *websocket.Conn) {
					defer func() {
						wg.Done()
						<-ch
					}()
					err := conn.WriteJSON(gin.H{"OnlineNum": c.clientNum})
					if err != nil {
						c.UnRegister(id)
					}
				}(k, v)
			}
		}
		wg.Wait() //等待所有协程结束
		runtime.GC()
	}
	close(ch)
}

func (c *Client) UnRegister(id int64) {
	companyName, companyType := GetCompanyBasicInfo(id) //仅需获取type即可
	log.Println(companyName, " is Unconnected!")
	if _, e := c.clients[id]; e {
		c.clients[id].Close() // 关闭websocket
		delete(c.clients, id) // 将信息从clients中删除
		delete(c.clientGroup[companyType], id)
		c.clientNum--
	}
}

func UseClient() *Client {
	return &myClient
}
