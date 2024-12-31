package biz

import (
	"context"

	v1 "convoo-accounts/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter 是业务逻辑中的领域模型，业务层的核心对象
type Greeter struct {
	ID       int64
	UserName string
	Greeting string
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Save(context.Context, *Greeter) (*Greeter, error)
	Update(context.Context, int64, *Greeter) (bool, error)
	FindByID(context.Context, int64) (*Greeter, error)
	ListAll(context.Context) ([]*Greeter, error)
	TestSetCache(context.Context, string) (string, error)
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter 数据库相关-保存
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	g, err := uc.repo.Save(ctx, g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// UpdateByID 数据库相关-更新
func (uc *GreeterUsecase) UpdateByID(ctx context.Context, id int64, userName, greeting string) (bool, error) {
	g := Greeter{
		UserName: userName,
		Greeting: greeting,
	}

	r, err := uc.repo.Update(ctx, id, &g)
	if err != nil {
		return false, err
	}

	return r, nil
}

// FindByID 数据库相关-通过id查询
func (uc *GreeterUsecase) FindByID(ctx context.Context, id int64) (*Greeter, error) {
	g, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// ListAll 数据库相关-查询所有
func (uc *GreeterUsecase) ListAll(ctx context.Context) ([]*Greeter, error) {
	greeters, err := uc.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return greeters, nil
}

// TestSetCache redis缓存-测试设置缓存并读取
func (uc *GreeterUsecase) TestSetCache(ctx context.Context, name string) (string, error) {
	uc.log.WithContext(ctx).Infof("redis缓存-测试设置缓存并读取 name: %v", name)

	val, err := uc.repo.TestSetCache(ctx, name)
	if err != nil {
		return "", err
	}

	return val, nil
}
