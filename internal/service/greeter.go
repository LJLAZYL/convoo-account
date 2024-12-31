package service

import (
	"context"

	v1 "convoo-accounts/api/helloworld/v1"
	"convoo-accounts/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello 接口测试-是否可以访问
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {

	return &v1.HelloReply{Message: "Hello ,welcome to " + in.Name}, nil
}

// Save 数据库相关-保存
func (s *GreeterService) Save(ctx context.Context, in *v1.SaveRequest) (*v1.SaveReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{UserName: in.UserName, Greeting: in.Greeting})
	if err != nil {
		return nil, err
	}

	return &v1.SaveReply{Id: g.ID}, nil
}

// UpdateByID 数据库相关-更新
func (s *GreeterService) UpdateByID(ctx context.Context, in *v1.UpdateByIDRequest) (*v1.UpdateByIDReply, error) {
	r, err := s.uc.UpdateByID(ctx, in.Id, in.UserName, in.Greeting)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateByIDReply{Result: r}, nil
}

// FindByID 数据库相关-通过id查询
func (s *GreeterService) FindByID(ctx context.Context, in *v1.FindByIDRequest) (*v1.FindByIDReply, error) {
	g, err := s.uc.FindByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &v1.FindByIDReply{Id: g.ID, UserName: g.UserName, Greeting: g.Greeting}, nil
}

// ListAll 数据库相关-查询所有
func (s *GreeterService) ListAll(ctx context.Context, in *v1.ListAllRequest) (*v1.ListAllReply, error) {
	greeters, err := s.uc.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	//将 biz.Greeter 转换为 v1.FindByIDReply 类型，假设 v1.FindByIDReply 是接口层返回的结构体
	var replyGreeters []*v1.FindByIDReply
	for _, greeter := range greeters {
		replyGreeters = append(replyGreeters, &v1.FindByIDReply{
			Id:       greeter.ID,
			UserName: greeter.UserName,
			Greeting: greeter.Greeting,
		})
	}

	return &v1.ListAllReply{Item: replyGreeters}, nil
}

// TestSetCache redis缓存-测试设置缓存并读取
func (s *GreeterService) TestSetCache(ctx context.Context, in *v1.TestSetCacheRequest) (*v1.TestSetCacheReply, error) {
	g, err := s.uc.TestSetCache(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &v1.TestSetCacheReply{Message: "success set redis,get value is:" + g}, nil
}
