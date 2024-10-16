package rocketmq

import (
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
	"github.com/atompi/aliyunbot/pkg/utils"
	"go.uber.org/zap"
)

func createApiClient(opts options.AliyunOptions) (*openapi.Client, error) {
	config := utils.CreateClientConfig(
		tea.String(opts.AccessKeyId),
		tea.String(opts.AccessKeySecret),
		tea.String(opts.RegionId),
		tea.String(opts.Endpoint),
	)

	return openapi.NewClient(config)
}

func createApiInfo(action string, pathName string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String("2022-08-01"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		Pathname:    tea.String(pathName),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
}

func callApi(t options.TaskOptions, action string, pathName string, body map[string]interface{}) error {
	c, err := createApiClient(t.Aliyun)
	if err != nil {
		zap.S().Errorf("create api client failed: %v", err)
		return err
	}

	params := createApiInfo(action, pathName)

	runtime := &teautil.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Body: body,
	}

	_, err = c.CallApi(params, request, runtime)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return err
	}

	return nil
}

func CreateTopic(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, instanceId string, topicName string, messageType string, remark string) {
	defer func() { wg.Done(); <-ch }()

	action := "CreateTopic"
	pathName := "/instances/" + instanceId + "/topics/" + topicName
	body := map[string]interface{}{
		"messageType": messageType,
		"remark":      remark,
	}

	err := callApi(t, action, pathName, body)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}

func CreateConsumerGroup(ch chan int, wg *sync.WaitGroup, t options.TaskOptions, instanceId string, consumerGroupId string, remark string) {
	defer func() { wg.Done(); <-ch }()

	action := "CreateConsumerGroup"
	pathName := "/instances/" + instanceId + "/consumerGroups/" + consumerGroupId
	body := map[string]interface{}{
		"remark": remark,
	}

	err := callApi(t, action, pathName, body)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return
	}
}
