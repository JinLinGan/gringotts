package communication

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/jinlingan/gringotts/message"
)

var instance *Client
var mux sync.Mutex

// Client 用于表示服务，负责与服务器通信
type Client struct {
	conn   *grpc.ClientConn
	client message.GringottsClient
}

// NewClient 使用单例模式新建 Client
func NewClient(address string) (*Client, error) {

	// if instance == nil {
	// 	mux.Lock()
	// 	defer mux.Unlock()
	if instance == nil {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		instance = &Client{
			conn:   conn,
			client: message.NewGringottsClient(conn)}
	}

	// }
	return instance, nil
}

//Close 关闭连接
func (s *Client) Close() {
	s.conn.Close()
}

func newHeartBeatRequest(agentID string) *message.HeartBeatRequest {
	req := message.HeartBeatRequest{
		AgnetId: agentID,
		Time:    time.Now().UnixNano(),
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unkonw"
		log.Printf("get hostname with err: %s", err)
	}
	req.HostName = hostname
	return &req
}

//HeartBeat 发送心跳
func (s *Client) HeartBeat(ctx context.Context, agentID string) (*message.HeartBeatResponse, error) {
	hb := newHeartBeatRequest(agentID)
	return s.client.HeartBeat(ctx, hb)
}

//DownloadFile 下载文件
func (s *Client) DownloadFile(filename string, sha1 string, destPath string) error {
	tf, err := ioutil.TempFile("/Users/jinlin/code/golang/src/github.com/jinlingan/gringotts/tmp", "")

	if err != nil {
		return err
	}

	err = tf.Chmod(0755)
	if err != nil {
		return fmt.Errorf("can not change mod of temp file %s : %s", tf.Name(), err)
	}
	fcClient, err := s.client.DownloadFile(
		context.Background(),
		&message.File{
			FileName: filename,
			Sha1Hash: sha1,
		},
	)
	if err != nil {
		return err
	}
	// // tf, err := os.OpenFile(tempFile.Name(), os.O_WRONLY|os.O_CREATE, 0775)
	// if err != nil {
	// 	return err
	// }
	defer tf.Close()

	for {
		fc, err := fcClient.Recv()
		if err == io.EOF {
			log.Printf("reach EOF : %s ,%v", err, fc)
			break
		}
		if err != nil {
			log.Printf("get unknow error from server: %s", err)
			return err
		}
		if _, err := tf.Write(fc.GetData()); err != nil {
			log.Printf("fail to write file %s : %s", tf.Name(), err)
			break
		}
	}

	err = os.Rename(tf.Name(), destPath+"/"+filename)

	if err != nil {
		rmerr := os.Remove(tf.Name())
		if rmerr != nil {
			return fmt.Errorf("mv file error : %q and remove temp file error %q", err, rmerr)
		} else {
			return fmt.Errorf("mv file error : %q", err)
		}
	}
	return nil
}
