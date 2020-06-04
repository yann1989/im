// Author: yann
// Date: 2019/9/22 上午11:52

package snowflake

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
	chat "yann-chat"
)

type SnowConfig struct {
	WorkerId     int64
	DatacenterId int64
}

type IdWorker struct {
	startTime             int64       //开始时间设置一个固定比现在早的毫秒时间戳
	workerIdBits          uint        //机器id位数
	datacenterIdBits      uint        //业务id位数
	maxWorkerId           int64       //机器id最大值
	maxDatacenterId       int64       //业务id最大值
	sequenceBits          uint        //毫秒内自增位
	workerIdLeftShift     uint        //机器id左偏移12位
	datacenterIdLeftShift uint        //业务id做偏移17位
	timestampLeftShift    uint        //时间毫秒左移22位
	sequenceMask          int64       //并发掩码
	workerId              int64       //机器id
	datacenterId          int64       //业务id
	sequence              int64       //并发控制
	lastTimestamp         int64       //上次生产id时间戳
	idLock                *sync.Mutex //并发控制
}

var (
	idWordker *IdWorker
	once      sync.Once
)

func (s *SnowConfig) Start(ctx context.Context, yannChat *chat.YannChat) error {
	once.Do(func() {
		idWordker = new(IdWorker)
		if err := idWordker.initIdWorker(s.WorkerId, s.DatacenterId); err != nil {
			logrus.Errorf("雪花id初始化失败:%s", err.Error())
		}
		logrus.Infof("snowflake[%d:%d] 初始化完成", s.WorkerId, s.DatacenterId)
	})
	return nil
}

func (s *SnowConfig) Stop(ctx context.Context) error {
	logrus.Infof("snowflake 已关闭")
	return nil
}

func NextId() int64 {
	return idWordker.nextId()
}

//workerId  datacenterId  传入每个服务各自的机器id和业务id
func (this *IdWorker) initIdWorker(workerId, datacenterId int64) error {
	var baseValue int64 = -1
	this.startTime = 1463834116272
	this.workerIdBits = 5
	this.datacenterIdBits = 5
	this.maxWorkerId = baseValue ^ (baseValue << this.workerIdBits)
	this.maxDatacenterId = baseValue ^ (baseValue << this.datacenterIdBits)
	this.sequenceBits = 12
	this.workerIdLeftShift = this.sequenceBits
	this.datacenterIdLeftShift = this.workerIdBits + this.workerIdLeftShift
	this.timestampLeftShift = this.datacenterIdBits + this.datacenterIdLeftShift
	this.sequenceMask = baseValue ^ (baseValue << this.sequenceBits) //1111 1111 1111
	this.sequence = 0
	this.lastTimestamp = -1
	this.idLock = &sync.Mutex{}

	if this.workerId < 0 || this.workerId > this.maxWorkerId {
		return errors.New(fmt.Sprintf("workerId[%v] 必须大于0并且小于 maxWorkerId[%v].", workerId, datacenterId))
	}
	if this.datacenterId < 0 || this.datacenterId > this.maxDatacenterId {
		return errors.New(fmt.Sprintf("datacenterId[%d] 必须大于0并且小于 maxDatacenterId[%d].", workerId, datacenterId))
	}
	this.workerId = workerId
	this.datacenterId = datacenterId
	return nil
}

//获取雪花id
func (this *IdWorker) nextId() int64 {
	this.idLock.Lock()
	//获取当前毫秒时间戳
	timestamp := this.timeGen()
	//如果最后时间戳和当前时间戳相同 则通过并发计算,获取下一个毫秒值
	if timestamp == this.lastTimestamp {
		this.sequence = (this.sequence + 1) & this.sequenceMask
		//如果毫秒内序列溢出(sequence = 1000000000000 & 111111111111) 则阻塞到下一个毫秒获取新的时间戳
		//相当于每一毫秒可以获得4095个id
		if this.sequence == 0 {
			timestamp = this.nextMillis()
		}
	} else { //时间戳改变,毫秒内序列重置
		this.sequence = 0
	}
	this.lastTimestamp = timestamp

	//拼接id
	id := ((timestamp - this.startTime) << this.timestampLeftShift) |
		(this.datacenterId << this.datacenterIdLeftShift) |
		(this.workerId << this.workerIdLeftShift) |
		this.sequence

	this.idLock.Unlock()
	if id < 0 {
		id = -id
	}
	return id
}

//阻塞到获取下一个毫秒值
func (this *IdWorker) nextMillis() int64 {
	timestamp := this.timeGen()
	for timestamp <= this.lastTimestamp {
		timestamp = this.timeGen()
	}
	return timestamp
}

//获取当前毫秒值
func (this *IdWorker) timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
