// Author       Yann
// Time         2020-05-23 08:16
// File Desc    服务接口

package chat

import "context"

type ServiceConstructor func(ctx *context.Context) (Service, error)

type Service interface {
	Start(ctx context.Context, yannChat *YannChat) error
	Stop(ctx context.Context) error
}
