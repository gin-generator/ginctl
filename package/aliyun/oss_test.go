package aliyun

import (
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/gin-generator/ginctl/package/get"
	"os"
	"testing"
)

func TestOss(t *testing.T) {

	get.NewViper("env.yaml", "../../app/http/admin/etc")

	client := NewOssClient()
	reader, err := os.Open("./WechatIMG339.jpg")
	if err != nil {
		t.Error(err)
		return
	}
	defer reader.Close()

	result, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr("week-eight"),
		Key:    oss.Ptr("WechatIMG339.jpg"),
		Body:   reader,
		ProgressFn: func(increment, transferred, total int64) {
			fmt.Printf("increment:%v, transferred:%v, total:%v\n", increment, transferred, total)
		},
	})

	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
