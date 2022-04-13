package Controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"main/Model"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	clients     map[int64]*websocket.Conn
	clientGroup map[string]map[int64]struct{}
}

var myClient Client
var num int64

func init() {
	myClient = Client{
		clients:     make(map[int64]*websocket.Conn),
		clientGroup: make(map[string]map[int64]struct{}),
	}
	myClient.clientGroup["货代"] = make(map[int64]struct{})
	myClient.clientGroup["船代"] = make(map[int64]struct{})
	myClient.clientGroup["货商"] = make(map[int64]struct{})

	go func() {
		for {
			select {
			case <-time.Tick(time.Minute * 5):
				myClient.HeartBeat()
			}
		}
	}()
}

func BuildSocket(c *gin.Context) {
	id, _ := c.Get("companyId")
	upGrader := websocket.Upgrader{ // 将http协议升级成websocket
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil) // 获取链接
	if err != nil {                                         // 错误检测
		log.Println("websocket connect default ", err, c.Param("channel"))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Websocket建立错误"})
		return
	}
	myClient.Login(conn, id.(int64)) // 将链接信息交到管理中心
}

func (c *Client) Login(conn *websocket.Conn, id int64) {
	c.clients[num] = conn
	_, companyType := Model.GetCompanyBasicInfo(id)
	c.clientGroup[companyType][id] = struct{}{}
}

func (c *Client) SendMessageToId(message interface{}, id int64) {
	c.clients[id].WriteJSON(message)
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
	wg := sync.WaitGroup{}
	ch := make(chan struct{}, 20) // 同时校验上限
	for k, v := range c.clients {
		wg.Add(1)
		ch <- struct{}{} //写入信息到channel用于计数
		go func(id int64, conn *websocket.Conn) {
			ms := make(chan string, 1)         //用于接收数据协程与检测协程的同步
			conn.WriteMessage(1, []byte("hb")) // 发送验证消息hb
			go func() {                        // 创建协程接收信息
				_, message, _ := conn.ReadMessage() //仅需接收到消息即可
				ms <- string(message)               // 将消息写入channel可用于表明接收到消息，即链接联通
			}()
			select {
			case <-time.After(time.Second * 10): //十秒超时断开链接
				c.UnRegister(id) //超时断开连接
			case <-ms: //接收到消息，没啥问题
			}
			<-ch      //协程结束，计数--
			wg.Done() //协程结束
		}(k, v)
	}
	wg.Wait() //等待所有协程结束
}

func (c *Client) UnRegister(id int64) {
	companyName, companyType := Model.GetCompanyBasicInfo(id) //仅需获取type即可
	log.Println(companyName, " is Unconnected!")
	c.clients[id].Close() // 关闭websocket
	delete(c.clients, id) // 将信息从clients中删除
	delete(c.clientGroup[companyType], id)
}

func UseClient() *Client {
	return &myClient
}
