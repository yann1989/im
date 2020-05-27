// Author       kevin
// Time         2019-08-27 14:58
// File Desc    封装公共的前后端交互的数据模型, 比如, 分页参数.

package view

const (
	MinPageNum      = 1  // 分页页数
	MinPageSize     = 1  // 分页大小
	DefaultPageSize = 10 // 默认分页大小
)

type PageVO struct {
	PageNum  int // 分页页数
	PageSize int // 分页大小
}

// 校验PageVO
func CheckPageVO(page PageVO) {

}
