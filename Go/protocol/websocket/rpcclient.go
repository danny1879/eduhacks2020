package websocket

import (
	"context"
	"eduhacks2020/Go/pkg/setting"
	"eduhacks2020/Go/protocol/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
)

func grpcConn(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	return conn
}

// SendRPC2Client 发送到 RPC 客户端
func SendRPC2Client(addr string, messageID, sendUserID, clientID string, code int, message string, data *string) {
	conn := grpcConn(addr)
	defer conn.Close()

	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.Port,
		"add":      addr,
		"clientID": clientID,
		"msg":      data,
	}).Info("send to server")

	c := pb.NewCommonServiceClient(conn)
	_, err := c.Send2Client(context.Background(), &pb.Send2ClientReq{
		MessageID:  messageID,
		SendUserID: sendUserID,
		ClientID:   clientID,
		Code:       int32(code),
		Message:    message,
		Data:       *data,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

// CloseRPCClient 关闭 RPC 客户端
func CloseRPCClient(addr string, clientID, systemID string) {
	conn := grpcConn(addr)
	defer conn.Close()

	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.Port,
		"add":      addr,
		"clientID": clientID,
	}).Info("send close connection to server")

	c := pb.NewCommonServiceClient(conn)
	_, err := c.CloseClient(context.Background(), &pb.CloseClientReq{
		SystemID: systemID,
		ClientID: clientID,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

// SendRPCBindGroup 绑定分组
func SendRPCBindGroup(addr string, systemID string, groupName string, clientID string, userID string, extend string) {
	conn := grpcConn(addr)
	defer conn.Close()

	c := pb.NewCommonServiceClient(conn)
	_, err := c.BindGroup(context.Background(), &pb.BindGroupReq{
		SystemID:  systemID,
		GroupName: groupName,
		ClientID:  clientID,
		UserID:    userID,
		Extend:    extend,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

// SendGroupBroadcast 发送分组消息
func SendGroupBroadcast(systemID string, messageID, sendUserID, groupName string, code int, message string, data *string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	for _, addr := range setting.GlobalSetting.ServerList {
		conn := grpcConn(addr)
		defer conn.Close()

		c := pb.NewCommonServiceClient(conn)
		_, err := c.Send2Group(context.Background(), &pb.Send2GroupReq{
			SystemID:   systemID,
			MessageID:  messageID,
			SendUserID: sendUserID,
			GroupName:  groupName,
			Code:       int32(code),
			Message:    message,
			Data:       *data,
		})
		if err != nil {
			log.Errorf("failed to call: %v", err)
		}
	}
}

// SendSystemBroadcast 发送系统信息
func SendSystemBroadcast(systemID string, messageID, sendUserID string, code int, message string, data *string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	for _, addr := range setting.GlobalSetting.ServerList {
		conn := grpcConn(addr)
		defer conn.Close()

		c := pb.NewCommonServiceClient(conn)
		_, err := c.Send2System(context.Background(), &pb.Send2SystemReq{
			SystemID:   systemID,
			MessageID:  messageID,
			SendUserID: sendUserID,
			Code:       int32(code),
			Message:    message,
			Data:       *data,
		})
		if err != nil {
			log.Errorf("failed to call: %v", err)
		}
	}
}

// GetOnlineListBroadcast 获取在线列表广播
func GetOnlineListBroadcast(systemID *string, groupName *string) (clientIDList []string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()

	serverCount := len(setting.GlobalSetting.ServerList)

	onlineListChan := make(chan []string, serverCount)
	var wg sync.WaitGroup

	wg.Add(serverCount)
	for _, addr := range setting.GlobalSetting.ServerList {
		go func(addr string) {
			conn := grpcConn(addr)
			defer conn.Close()
			c := pb.NewCommonServiceClient(conn)
			response, err := c.GetGroupClients(context.Background(), &pb.GetGroupClientsReq{
				SystemID:  *systemID,
				GroupName: *groupName,
			})
			if err != nil {
				log.Errorf("failed to call: %v", err)
			} else {
				onlineListChan <- response.List
			}
			wg.Done()

		}(addr)
	}

	wg.Wait()

	for i := 1; i <= serverCount; i++ {
		list, ok := <-onlineListChan
		if ok {
			clientIDList = append(clientIDList, list...)
		} else {
			return
		}
	}
	close(onlineListChan)

	return
}
