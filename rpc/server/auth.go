package main

// 注意，auth服务对所有信任的服务进行服务，不与底层数据库打交道，仅用于发布和鉴定token是否正确
import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/frame/g"
	_ "llfile/config"
	auth "llfile/rpc/authentication"
	"llfile/service"
	"llfile/util"
	"log"

	"context"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/llightos/efind"
	"google.golang.org/grpc"
)

var stru = service.JWTKey{Key: []byte(service.PrivateKey)}

type AuthServer struct {
	auth.UnimplementedAuthServer
}

func main() {
	node := "localhost:50001"
	session, err := efind.NewClient(efind.Config{
		EtcdAddr: "localhost:2379",
		TTL:      5,
	}).NewSession()
	if err != nil {
		log.Panicln(err)
	}
	elect := session.NewElect("auth", node)
	err = elect.Campaign()
	if err != nil {
		err.Error()
	}

	listen, err := net.Listen("tcp", node)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServer(grpcServer, &AuthServer{})

	g.Log("server start").Print(context.TODO(), fmt.Sprintf("server start at %s", node))
	if err = grpcServer.Serve(listen); err != nil {
		g.Log().Error(context.TODO(), err)
		log.Println(err)
		return
	}
}

func (a *AuthServer) CreateToken(ctx context.Context, CreateTokenReq *auth.CreateTokenReq) (CreateTokenResp *auth.CreateTokenResp, err error) {
	token, err := stru.CreateToken(service.TargetUser{
		UserID: int(CreateTokenReq.Userid),
		Mix:    gbase64.EncodeToString(util.Encrypt([]byte(strconv.Itoa(int(CreateTokenReq.Userid))), []byte(service.SecretKey))),
	})
	if err != nil {
		panic(err)
		return nil, err
	}
	return &auth.CreateTokenResp{
		Ok:    true,
		Token: token,
		Data:  "",
	}, nil
}

func (a *AuthServer) AuthToken(ctx context.Context, AuthTokenReq *auth.AuthTokenReq) (AuthTokenResp *auth.AuthTokenResp, err error) {
	//defer func() {
	//	recover()
	//}()
	target, err := stru.ParserToken(AuthTokenReq.Token)

	if err != nil {
		log.Println("stru.ParserToken err", err)
		return &auth.AuthTokenResp{
			Ok:   false,
			Data: err.Error(),
		}, nil
	}

	//fmt.Println("请求鉴权id，实际鉴权id", AuthTokenReq.Userid, target.UserID)
	if AuthTokenReq.Userid != uint32(target.UserID) {
		return &auth.AuthTokenResp{
			Ok:   false,
			Data: "鉴权id错误",
		}, nil
	}
	if err != nil {
		return nil, errors.New(err.Error() + "ParserToken err")
	}
	decodeString, err := gbase64.DecodeString(target.Mix)
	if err != nil {
		log.Println(err)
	}
	decrypt := util.Decrypt(decodeString, []byte(service.SecretKey))
	if string(decrypt) != strconv.Itoa(target.UserID) {
		return &auth.AuthTokenResp{
			Ok:   false,
			Data: "err token",
		}, err
	} else if target.ExpiresAt < time.Now().Unix() {
		return &auth.AuthTokenResp{
			Ok:   false,
			Data: "token time out",
		}, err
	}
	fmt.Println("什么情况")
	return &auth.AuthTokenResp{
		Ok:   true,
		Data: "",
	}, nil
}
