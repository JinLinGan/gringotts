// Package communication 用于与服务端通信，暂时只支持 GRPC
package communication

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/jinlingan/gringotts/gringotts-agent/config"
	"github.com/jinlingan/gringotts/gringotts-agent/log"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
	"github.com/jinlingan/gringotts/message"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Client 用于表示服务，负责与服务器通信
type Client struct {
	conn   *grpc.ClientConn
	client message.GringottsClient
	logger log.Logger
}

// NewClient 使用单例模式新建 Client
func NewClient(cfg *config.AgentConfig, logger log.Logger) (*Client, error) {

	conn, err := grpc.Dial(cfg.GetServerAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	instance := &Client{
		conn:   conn,
		client: message.NewGringottsClient(conn),
		logger: logger,
	}

	return instance, nil
}

//Close 关闭连接
func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		c.logger.Infof("close server conn error: %v", err)
	}
}

func (c *Client) newHeartBeatRequest(agentID string) *message.HeartBeatRequest {
	req := message.HeartBeatRequest{
		AgentId: agentID,
		Time:    time.Now().UnixNano(),
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
		c.logger.Infof("get hostname with err: %s", err)
	}
	req.HostName = hostname
	return &req
}

//HeartBeat 发送心跳
func (c *Client) HeartBeat(ctx context.Context, agentID string) (*message.HeartBeatResponse, error) {
	hb := c.newHeartBeatRequest(agentID)
	fmt.Println(c.client)
	return c.client.HeartBeat(ctx, hb)
}

//DownloadFile 下载文件
func (c *Client) DownloadFile(filename string, sha1 string, destPath string, tempPath string) error {

	fcClient, err := c.client.DownloadFile(
		context.Background(),
		&message.File{
			FileName: filename,
			Sha1Hash: sha1,
		},
	)
	if err != nil {
		return err
	}

	tf, err := ioutil.TempFile(tempPath, "")
	if err != nil {
		return errors.Wrapf(err, "can not create temp file in %s", tempPath)
	}
	err = tf.Chmod(0755)
	if err != nil {
		return errors.Wrapf(err, "can not change mod of temp file %s", tf.Name())
	}

	for {
		fc, err := fcClient.Recv()
		if err == io.EOF {
			c.logger.Infof("reach EOF : %s ,%v", err, fc)
			break
		}
		if err != nil {
			c.logger.Infof("get unknown error from server: %s", err)
			return err
		}
		if _, err := tf.Write(fc.GetData()); err != nil {
			c.logger.Infof("fail to write file %s : %s", tf.Name(), err)
			break
		}
	}

	err = os.Rename(tf.Name(), destPath+string(os.PathSeparator)+filename)

	if err != nil {
		return errors.Wrapf(err, "mv file error")

	}

	if err := tf.Close(); err != nil {
		c.logger.Infof("can not close tmp file %s: %v", tf.Name(), err)
	}
	if err := os.Remove(tf.Name()); err != nil {
		c.logger.Infof("can not remove tmp file %s: %v", tf.Name(), err)
	}

	return nil
}

//Register 注册 agent
func (c *Client) Register(hostName string, netInfos *model.NetInfos) (*model.RegisterResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req := &message.RegisterRequest{
		HostName: hostName,
	}
	req.NetInfo = make([]*message.RegisterRequest_NetInfo, len(*netInfos))

	index := 0
	for _, v := range *netInfos {
		n := &message.RegisterRequest_NetInfo{
			IpAddress:  v.IPAddress,
			MacAddress: v.MacAddress,
		}
		req.NetInfo[index] = n
		index++
	}

	resp, err := c.client.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	ret := &model.RegisterResp{
		AgentID:       resp.AgentId,
		ConfigVersion: resp.ConfigVersion,
	}
	return ret, nil
}
