package main

import (
	"llfile/model"
	"llfile/rpc"
	"llfile/rpc/user"

	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/llightos/efind"
	"google.golang.org/grpc"
)

type UserServer struct {
	db *model.DB
	user.UnimplementedUserServerServer
}

func main() {
	node := "localhost:51001"
	session, err := efind.NewClient(efind.Config{
		EtcdAddr: "localhost:2379",
		TTL:      5,
	}).NewSession()
	if err != nil {
		log.Panicln(err)
	}
	elect := session.NewElect("user", node)
	err = elect.Campaign()
	if err != nil {
		err.Error()
	}

	listen, err := net.Listen("tcp", node)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	user.RegisterUserServerServer(grpcServer, &UserServer{})

	g.Log("user server start").Print(context.TODO(), fmt.Sprintf("user server start at %s", node))
	if err = grpcServer.Serve(listen); err != nil {
		g.Log().Error(context.TODO(), err)
		log.Println(err)
		return
	}
}

func (u *UserServer) Register(ctx context.Context, userRegisterReq *user.RegisterRequest) (userRegisterRes *user.RegisterResponse, err error) {
	u.db = model.NewModelDB()

	id, err := u.db.AddUser(userRegisterReq.UserName, userRegisterReq.PassWord)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResponse{
		Ok: true,
		Id: uint32(id),
	}, nil
}

func (u *UserServer) Login(ctx context.Context, userLoginRequest *user.LoginRequest) (res *user.LoginResponse, err error) {
	u.db = model.NewModelDB()

	userId, err := u.db.GetUser(userLoginRequest.GetUserName(), userLoginRequest.GetPassWord())

	fmt.Println("userID", userId)

	if err != nil {
		return nil, err
	}

	token, err := rpc.NewAuth().CallCreateToken(userId)
	if err != nil {
		return nil, errors.New("create token err")
	}

	return &user.LoginResponse{
		Ok:    true,
		Token: token.Token,
		Id:    uint32(userId),
	}, nil
}
