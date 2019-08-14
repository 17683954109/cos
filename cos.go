package main

import(
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

// Cos 客户端
type CosClient struct{
	url     string
	client  cos.Client
	auth    auth
}

// 凭证接口, 需要实现获取客户端和获取访问链接方法
type auth interface {
	GetClient() *cos.Client
	GetBaseUrl(string)
}

// 永久凭证
type authgrition struct{
	SecretID  string
	SecretKey string
	BaseUrl   *cos.BaseURL
}

// 永久凭证获取 Client 方法
func (a *authgrition) GetClient() *cos.Client{
	if a.BaseUrl == nil{
		return nil
	}
	return cos.NewClient(a.BaseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID: a.SecretID,
			SecretKey: a.SecretKey,                      
		},           
	})
}

// 永久凭证获取访问地址的方法
func (a *authgrition) GetBaseUrl(bucketUrl string){
	u, _ := url.Parse(bucketUrl)
	a.BaseUrl = &cos.BaseURL{BucketURL: u}
}

// 临时凭证
type tmpAuthgrition struct{
	SecretID string
	SecretKey string   
	SessionToken string
	BaseUrl  *cos.BaseURL
}

// 临时凭证获取 Client 方法
func (t *tmpAuthgrition) GetClient() *cos.Client{
	if t.BaseUrl == nil{
		return nil
	}
	return cos.NewClient(t.BaseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID: t.SecretID,
			SecretKey: t.SecretKey,    
			SessionToken: t.SessionToken,
		},           
	});
}

// 临时凭证获取访问链接的方法
func (t *tmpAuthgrition) GetBaseUrl(bucketUrl string){
	u, _ := url.Parse(bucketUrl)
	t.BaseUrl = &cos.BaseURL{BucketURL: u}
}

// 设置一个临时凭证
func (c *CosClient) NewTmpAuthgrition(secretId string, secretKey string, sessionToken string){
	c.auth = (auth)(&tmpAuthgrition{
		SecretID: secretId,
		SecretKey: secretKey,   
		SessionToken: sessionToken,
	})
}

// 设置一个永久凭证
func (c *CosClient) NewAuthgrition (secretId string, secretKey string){
	c.auth = (auth)(&authgrition{
		SecretID: secretId,
		SecretKey: secretKey,
	})
}

// 设置存储桶地址
func (c *CosClient) Init(bucketUrl string) {
	c.auth.GetBaseUrl(bucketUrl)
}

// TODO: 完成 COS 客户端 CURD 操作
