package data

import (
	"context"
	"convoo-accounts/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"

	"convoo-accounts/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Save 数据库相关-保存
func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	greeter := &models.Greetings{
		UserName: g.UserName,
		Greeting: g.Greeting,
	}

	// 插入 Greeter 数据到数据库
	if err := r.data.DB.WithContext(ctx).Select([]string{"ID", "UserName", "Greeting", "ModifyTime"}).Create(&greeter).Error; err != nil {
		return nil, err
	}

	g.ID = greeter.ID
	return g, nil
}

// Update 数据库相关-更新
func (r *greeterRepo) Update(ctx context.Context, id int64, g *biz.Greeter) (bool, error) {
	greeter := &models.Greetings{}

	result := r.data.DB.WithContext(ctx).Model(&greeter).Where("id", id).Updates(models.Greetings{UserName: g.UserName, Greeting: g.Greeting})

	// 判断是否出错
	if result.Error != nil {
		return false, result.Error
	}

	// 判断是否有记录被更新
	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

// FindByID 数据库相关-通过id查询
func (r *greeterRepo) FindByID(ctx context.Context, id int64) (*biz.Greeter, error) {
	g := models.Greetings{}

	err := r.data.DB.WithContext(ctx).First(&g, "id=?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到记录，可以返回自定义的错误
			return nil, fmt.Errorf("Greeter with ID %d not found", id)
		}
		return nil, err
	}

	g1 := biz.Greeter{ID: g.ID, UserName: g.UserName, Greeting: g.Greeting}
	return &g1, nil
}

// ListAll 数据库相关-查询所有
func (r *greeterRepo) ListAll(ctx context.Context) ([]*biz.Greeter, error) {
	// 定义一个存储查询结果的切片
	var greetings []models.Greetings

	// 执行查询
	if err := r.data.DB.WithContext(ctx).Find(&greetings).Error; err != nil {
		return nil, err
	}

	// 将查询结果转换为 biz.Greeter 类型的切片
	var greeters []*biz.Greeter
	for _, greeting := range greetings {
		greeters = append(greeters, &biz.Greeter{
			ID:       greeting.ID,
			UserName: greeting.UserName,
			Greeting: greeting.Greeting,
		})
	}

	return greeters, nil
}

// TestSetCache 设置缓存
func (r *greeterRepo) TestSetCache(ctx context.Context, name string) (string, error) {

	key := "test:s:c"

	err := r.data.Redis.Set(ctx, key, name, 3600*time.Second).Err()
	if err != nil {

		return "", err
	}

	val := r.data.Redis.Get(ctx, key).Val()

	return val, nil
}
