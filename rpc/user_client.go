package rpc

import (
	"llfile/rpc/user"

	"context"
	"fmt"
	"log"

	"github.com/llightos/efind"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn *grpc.ClientConn
}

func NewUser() *UserClient {
	client := new(UserClient)
	kv, err := efind.NewClient(efind.Config{
		EtcdAddr: "127.0.0.1:2379",
		TTL:      5,
	}).MatchAServer("user")

	fmt.Println("find a node", kv)
	if err != nil {
		log.Println(err)
	}
	client.conn, err = grpc.Dial(kv.Val, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	return client
}

func (c *UserClient) CallRegister(username, password string) (ok bool, id uint) {
	userClient := user.NewUserServerClient(c.conn)
	userRegisterRes, err := userClient.Register(context.Background(), &user.RegisterRequest{
		UserName: username,
		PassWord: password,
	})
	if err != nil {
		log.Println(err)
		return false, 0
	}
	return userRegisterRes.Ok, uint(userRegisterRes.Id)
}

func (c *UserClient) CallLogin(username, password string) (ok bool, token string, userId uint) {
	userClient := user.NewUserServerClient(c.conn)
	userLoginRes, err := userClient.Login(context.Background(), &user.LoginRequest{
		UserName: username,
		PassWord: password,
	})
	//fmt.Println("resss", userLoginRes.Token)
	if err != nil {
		log.Println(err)
		return false, "", 0
	}
	return true, userLoginRes.Token, uint(userLoginRes.Id)
}
