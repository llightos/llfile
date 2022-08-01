package rpc

// 客户端代码是公用代码，所用端服务内的接口都可以调用，
// 用于颁发token和验证token是否过期，是否合法
import (
	"fmt"
	_ "llfile/config"
	auth "llfile/rpc/authentication"

	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/llightos/efind"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type AuthClient struct {
	conn *grpc.ClientConn
}

func NewAuth() *AuthClient {
	client := new(AuthClient)
	kv, err := efind.NewClient(efind.Config{
		EtcdAddr: "127.0.0.1:2379",
		TTL:      5,
	}).MatchAServer("auth")

	fmt.Println("find a node", kv)
	if err != nil {
		log.Println(err)
		g.Log("efind err").Error(context.TODO(), err)
	}
	client.conn, err = grpc.Dial(kv.Val, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		g.Log("err").Print(context.TODO(), err)
	}
	return client
}

func (c *AuthClient) CallCreateToken(userid uint) (*auth.CreateTokenResp, error) {
	authClient := auth.NewAuthClient(c.conn)
	authRes, err := authClient.CreateToken(context.TODO(), &auth.CreateTokenReq{Userid: uint32(userid)})
	if err != nil {
		g.Log("err").Print(context.TODO(), err)
		return &auth.CreateTokenResp{
			Ok:    false,
			Token: "",
			Data:  "err",
		}, err
	}
	return authRes, nil
}

func (c *AuthClient) CallAuthToken(token string, userid int) (*auth.AuthTokenResp, error) {
	authClient := auth.NewAuthClient(c.conn)
	authToken, err := authClient.AuthToken(context.TODO(), &auth.AuthTokenReq{
		Userid: uint32(userid),
		Token:  token,
		Data:   "",
	})
	if err != nil {
		g.Log("err").Print(context.TODO(), err)
		return &auth.AuthTokenResp{
			Ok:   false,
			Data: "err",
		}, err
	}
	return authToken, nil
}
