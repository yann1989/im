// Author       kevin
// Time         2019-08-08 20:16
// File Desc    服务接口(interface.go)实现类

package chat

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

type YannChat struct {
	lock         *sync.RWMutex
	serviceFuncs []ServiceConstructor //服务构造回调
	services     map[string]Service   //目前正在运行子服务(一个子服务分配一个goroutine)
}

func NewYannChat() *YannChat {
	return &YannChat{
		services:     make(map[string]Service),
		serviceFuncs: make([]ServiceConstructor, 1),
		lock:         new(sync.RWMutex),
	}
}

// 根据服务对象类型获取服务对象
func (c *YannChat) Service(serviceType string) (result interface{}, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if element, ok := c.services[serviceType]; ok {
		return element, nil
	}
	return nil, errors.New("unknown service object")
}

// 注册服务至Taurus中
func (c *YannChat) Register(constructor ServiceConstructor) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if constructor == nil {
		return errors.New("constructor cannot be nil")
	}
	c.serviceFuncs = append(c.serviceFuncs, constructor)
	return nil
}

// 构建Taurus
func (c *YannChat) Build(ctx context.Context) error {
	for _, constructor := range c.serviceFuncs[1:] {
		service, err := constructor(nil)
		if err != nil {
			return err
		}
		kind := reflect.TypeOf(service).String()
		if _, exists := c.services[kind]; exists {
			return errors.New("service and existence")
		}
		c.services[kind] = service
	}
	return nil
}

// 启动已经注册后的服务
func (c *YannChat) Start(ctx context.Context) error {
	for _, s := range c.services {
		if err := s.Start(ctx, c); err != nil {
			return err
		}
	}
	return nil
}

// 停止所有的服务
func (c *YannChat) Stop(ctx context.Context) error {
	for _, s := range c.services {
		if err := s.Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}
