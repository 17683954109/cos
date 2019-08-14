package cos

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Cos 客户端
type Client struct {
	url    string
	client *cos.Client
	auth   auth
}

// 凭证接口, 需要实现获取客户端和获取访问链接方法
type auth interface {
	GetClient() *cos.Client
	GetBaseUrl(string)
	GetClientWithoutBucket() *cos.Client
}

// 永久凭证
type authorization struct {
	SecretID  string
	SecretKey string
	BaseUrl   *cos.BaseURL
}

// 永久凭证获取 Client 方法
func (a *authorization) GetClient() *cos.Client {
	if a.BaseUrl == nil {
		return nil
	}
	return cos.NewClient(a.BaseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  a.SecretID,
			SecretKey: a.SecretKey,
		},
	})
}

// 获取存储桶列表
func (a *authorization) GetClientWithoutBucket() *cos.Client {
	return cos.NewClient(nil, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  a.SecretID,
			SecretKey: a.SecretKey,
		},
	})
}

// 永久凭证获取访问地址的方法
func (a *authorization) GetBaseUrl(bucketUrl string) {
	u, _ := url.Parse(bucketUrl)
	a.BaseUrl = &cos.BaseURL{BucketURL: u}
}

// 临时凭证
type tmpAuthorization struct {
	SecretID     string
	SecretKey    string
	SessionToken string
	BaseUrl      *cos.BaseURL
}

// 临时凭证获取 Client 方法
func (t *tmpAuthorization) GetClient() *cos.Client {
	if t.BaseUrl == nil {
		return nil
	}
	return cos.NewClient(t.BaseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     t.SecretID,
			SecretKey:    t.SecretKey,
			SessionToken: t.SessionToken,
		},
	})
}

// 临时凭证获取访问链接的方法
func (t *tmpAuthorization) GetBaseUrl(bucketUrl string) {
	u, _ := url.Parse(bucketUrl)
	t.BaseUrl = &cos.BaseURL{BucketURL: u}
}

// 获取存储桶列表
func (t *tmpAuthorization) GetClientWithoutBucket() *cos.Client {
	return cos.NewClient(nil, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     t.SecretID,
			SecretKey:    t.SecretKey,
			SessionToken: t.SessionToken,
		},
	})
}

// 设置一个临时凭证
func (c *Client) NewTempAuthorization(secretId string, secretKey string, sessionToken string) {
	c.auth = (auth)(&tmpAuthorization{
		SecretID:     secretId,
		SecretKey:    secretKey,
		SessionToken: sessionToken,
	})
}

// 设置一个永久凭证
func (c *Client) NewAuthorization(secretId string, secretKey string) {
	c.auth = (auth)(&authorization{
		SecretID:  secretId,
		SecretKey: secretKey,
	})
}

// 设置存储桶地址
func (c *Client) Init(bucketUrl string) {
	c.auth.GetBaseUrl(bucketUrl)
	c.client = c.auth.GetClient()
}

// 获取存储桶列表
func (c *Client) GetBucketList() ([]cos.Bucket, error){
	c.client = c.auth.GetClientWithoutBucket()
	s, _, err := c.client.Service.Get(context.Background())
	if err != nil {
		return nil, err
	}
	return s.Buckets, err
}

// 上传对象(本地文件)
func (c *Client) Upload(name string, filePath string) error {
	_, err := c.client.Object.PutFromFile(context.Background(), name, filePath, nil)
	return err
}

// 上传文件内容
func (c *Client) PutString(name string, content string) error {
	f := strings.NewReader(content)
	_, err := c.client.Object.Put(context.Background(), name, f, nil)
	return err
}

// 读取内容列表
func (c *Client) GetObjList(path string, maxList int) ([]cos.Object, error) {
	opt := &cos.BucketGetOptions{
		Prefix:  path,
		MaxKeys: maxList,
	}
	v, _, err := c.client.Bucket.Get(context.Background(), opt)
	if err != nil {
		return nil, err
	}
	return v.Contents, nil
}

// 下载对象
func (c *Client) Download(name string, savePath string) error {
	_, err := c.client.Object.GetToFile(context.Background(), name, savePath, nil)
	return err
}

// 读取内容, 而不保存
func (c *Client) ReadString(name string) (string, error) {
	resp, err := c.client.Object.Get(context.Background(), name, nil)
	if err != nil {
		return "", err
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return string(bs), nil
}

// 删除对象
func (c *Client) Delete(name string) error {
	_, err := c.client.Object.Delete(context.Background(), name)
	return err
}
