// Package communication 用于与服务端通信，暂时只支持 gRPC
package communication

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/jinlingan/gringotts/pkg/metadata/host"

	"github.com/jinlingan/gringotts/gringotts-agent/config"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
	"github.com/jinlingan/gringotts/pkg/log"
	"github.com/jinlingan/gringotts/pkg/message"
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

func NewClientForTools(serverAddr string) (*Client, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	instance := &Client{
		conn:   conn,
		client: message.NewGringottsClient(conn),
		logger: log.NewStdoutLogger(),
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
		AgentID:    agentID,
		ClientTime: time.Now().Unix(),
	}

	//TODO:带上 Job 信息
	return &req
}

//HeartBeat 发送心跳
func (c *Client) HeartBeat(ctx context.Context, agentID string) (*message.HeartBeatResponse, error) {
	hb := c.newHeartBeatRequest(agentID)
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
	defer func() {
		if err := tf.Close(); err != nil {
			c.logger.Infof("can not close tmp file %s: %v", tf.Name(), err)
		}
	}()
	defer func() {
		if err := os.Remove(tf.Name()); err != nil {
			c.logger.Infof("can not remove tmp file %s: %v", tf.Name(), err)
		}
	}()
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
			break
		}
		if err != nil {
			return errors.Wrap(err, "get unknown error from server")
		}
		if _, err := tf.Write(fc.GetData()); err != nil {
			return errors.Wrapf(err, "write to temp file %s fail", tf.Name())
		}
	}

	err = os.Rename(tf.Name(), destPath+string(os.PathSeparator)+filename)

	if err != nil {
		return errors.Wrapf(err, "mv file error")

	}

	return nil
}

//Register 注册 agent
func (c *Client) Register(agentID string, hostInfo *host.HostInfo) (*model.RegisterResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req := &message.RegisterRequest{
		AgentID:              agentID,
		HostName:             hostInfo.HostName,
		HostUUID:             hostInfo.HostUUID,
		Os:                   hostInfo.OS,
		Platform:             hostInfo.Platform,
		PlatformFamily:       hostInfo.PlatformFamily,
		PlatformVersion:      hostInfo.PlatformVersion,
		KernelVersion:        hostInfo.KernelVersion,
		VirtualizationSystem: hostInfo.VirtualizationSystem,
		VirtualizationRole:   hostInfo.VirtualizationRole,
	}

	for _, n := range hostInfo.Interfaces {
		nn := &message.RegisterRequest_Interface{
			Name:         n.Name,
			HardwareAddr: n.HardwareAddr,
		}

		for _, a := range n.Addrs {
			nn.IpAddrs = append(nn.IpAddrs, a)
		}
		req.Interfaces = append(req.Interfaces, nn)
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

//LoadConfig 从服务端加载配置
func (c *Client) LoadConfig(agentID string) (*message.GetJobsResponse, error) {
	req := &message.GetJobsRequest{
		AgentID: agentID,
	}
	return c.client.GetJobs(context.Background(), req)
}
