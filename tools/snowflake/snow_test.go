/**
 * @Author: yann
 * @Date: 2019/9/14 下午2:17
 */
package snowflake

import (
	"fmt"
	"testing"
	"time"
)

//理论每秒产生4095000个id 这里测试包含了其他判断,赋值等操作所以2s
func TestIdWorker_NextId(t *testing.T) {
	err := IDWORKER.InitIdWorker(1, 1)
	if err != nil {
		panic(err)
	}
	ids := make(map[int64]int64)
	unix := time.Now().Unix()
	for i := 0; i < 4095000; i++ {
		id := IDWORKER.NextId()
		if _, ok := ids[id]; ok {
			fmt.Println("已经包含这个id:", ids[id])
			break
		}
		ids[id] = id
		//fmt.Println(id)
	}
	fmt.Printf("4095000个id用时%d 秒\n:", time.Now().Unix()-unix)
}

//获取下一个时间戳
func TestNextMillis(t *testing.T) {
	err := IDWORKER.InitIdWorker(1, 1)
	if err != nil {
		fmt.Println("初始化失败")
	}
	for i := 0; i < 10; i++ {
		IDWORKER.lastTimestamp = IDWORKER.nextMillis()
		fmt.Println(IDWORKER.lastTimestamp)
	}
}

//获取当前时间戳
func TestTimeGen(t *testing.T) {
	err := IDWORKER.InitIdWorker(1, 1)
	if err != nil {
		fmt.Println("初始化失败")
	}
	for i := 0; i < 10; i++ {
		fmt.Println(IDWORKER.timeGen())
	}
}

func TestMath(t *testing.T) {
	i := -1 << 12
	fmt.Printf("%b\n", i)
	i = -1 ^ i
	fmt.Printf("%b\n", i)

	for j := 1; j <= 4096; j++ {
		fmt.Println(j & i)
	}
}

func TestGoRouting(t *testing.T) {
	err := IDWORKER.InitIdWorker(1, 1)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(IDWORKER.NextId())
		}()
	}
}
