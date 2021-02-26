package aliyun

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dawei101/gor/rconfig"
	"github.com/dawei101/gor/rlog"
)

type BucketConfig struct {
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	AccessKeyID     string `json:"accessKeyID" yaml:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`
	BucketName      string `json:"bucketName" yaml:"bucketName"`
	UrlPrefix       string `json:"urlPrefix" yaml:"urlPrefix"`
}

// https://github.com/aliyun/aliyun-oss-go-sdk
type Bucket struct {
	*oss.Bucket
	Config *BucketConfig
	Once   sync.Once
}

func (b *Bucket) init() error {
	var err error = nil
	b.Once.Do(func() {
		var client *oss.Client
		c := b.Config
		client, err = oss.New(c.Endpoint, c.AccessKeyID, c.AccessKeySecret)
		if err != nil {
			panic(err.Error())
		}
		bucket, err := client.Bucket(c.BucketName)
		if err != nil {
			panic(err.Error())
		}
		b.Bucket = bucket
	})
	return err
}

func (b *Bucket) FileUrl(objectKey string) string {
	u, _ := url.Parse(b.Config.UrlPrefix)
	u.Path = path.Join(u.Path, objectKey)
	return u.String()
}

func (b *Bucket) PutObjectWithBytes(objectKey string, byts []byte) (string, error) {
	r := strings.NewReader(string(byts))
	err := b.PutObject(objectKey, r)
	return b.FileUrl(objectKey), err
}

var buckets = make(map[string]*Bucket)
var default_config_once sync.Once

// 上传示例：
// b, _ := GetOssBucket("config")
// h.PutObject(objkey, f)

func GetOssBucket(configname string) (*Bucket, error) {
	default_config_once.Do(func() {
		configs := make(map[string]*BucketConfig)
		rconfig.DefConf().ValueMustAssignTo("aliyun.oss", &configs)
		rlog.Info(context.Background(), "load aliyun.oss config:", configs)
		for name, config := range configs {
			buckets[name] = &Bucket{Config: config}
		}
	})
	bucket, ok := buckets[configname]
	if !ok {
		return nil, errors.New("no oss config named:" + configname)
	}
	if err := bucket.init(); err != nil {
		return nil, err
	}
	return bucket, nil
}

func DefaultOssBucket() (*Bucket, error) {
	return GetOssBucket("default")
}
