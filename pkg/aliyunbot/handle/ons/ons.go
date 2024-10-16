package ons

import (
	"sync"

	"github.com/atompi/aliyunbot/pkg/aliyun/ons"
	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
	"github.com/atompi/aliyunbot/pkg/dataio"
	"github.com/atompi/aliyunbot/pkg/utils"
	"go.uber.org/zap"
)

func CreateTopicHandler(t options.TaskOptions) error {
	res, err := dataio.InputCSV(t.Input)
	if err != nil {
		zap.S().Errorf("input error: %v", err)
		return err
	}

	data, err := utils.DataToMap(&res)
	if err != nil {
		zap.S().Errorf("data convert error: %v", err)
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan int, t.Threads)

	for _, row := range *data {
		wg.Add(1)
		ch <- 1
		go ons.CreateTopic(ch, &wg, t, row["instanceId"], row["topicName"], row["messageType"], row["remark"])
	}

	wg.Wait()
	return nil
}

func CreateConsumerGroupHandler(t options.TaskOptions) error {
	res, err := dataio.InputCSV(t.Input)
	if err != nil {
		zap.S().Errorf("input error: %v", err)
		return err
	}

	data, err := utils.DataToMap(&res)
	if err != nil {
		zap.S().Errorf("data convert error: %v", err)
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan int, t.Threads)

	for _, row := range *data {
		wg.Add(1)
		ch <- 1
		go ons.CreateConsumerGroup(ch, &wg, t, row["instanceId"], row["consumerGroupId"], row["remark"])
	}

	wg.Wait()
	return nil
}