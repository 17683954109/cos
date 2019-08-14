##### Cos 二次封装 (基于腾讯云对象存储 SDK v5)

> 下载, *github连不上, 不知道啥原因 ...*

```bash
    $ mkdir -p $GOPATH/src/github.com/17683954109
    $ cd $GOPATH/src/github.com/17683954109
    $ git clone -b master https://gitee.com/zjz17683954109/cos.git
```

> Api

```go
    package main

    import(
    	"fmt"
    	"github.com/17683954109/cos"
    )
    func main(){
    	
    	// 初始化
        client := &cos.Client{}
        
        // 设置认证信息
        client.NewAuthorization("SecretID", "SecretKey")
        
        // 获取存储桶列表, 这个不用预先调用 Init
        fmt.Println(client.GetBucketList())
        
        // 初始化一个 Bucket 访问
        client.Init("BucketUrl")
        
        // 读取路径下的对象列表, 结果数不超过 12
        res, err := client.GetObjList("doc/", 12)
        if err != nil {
            panic(err)
        }
        fmt.Println(res)
        
        // 下载文件到本地
        _ = client.Download("kele.mp3", "kele.mp3")
        
        // 不保存, 结果以字符串方式返回
        str, _ := client.ReadString("test/my.txt")
        fmt.Println(str)
        
        // 上传
        _ = client.Upload("kele.mp3", "test/kele.mp3")
        
        // 上传, 但写入为字符串
        _ = client.PutString("test.md", "这是内容")
        
        // 删除对象
        _ = client.Delete("kele.mp3")
        
        fmt.Println("Cos Client")
    }
``` 

> 简陋, 个人业余使用
