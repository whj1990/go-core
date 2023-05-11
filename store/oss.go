package store

import (
	"bytes"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	uuid "github.com/satori/go.uuid"
	"github.com/whj1990/go-core/config"
	"github.com/whj1990/go-core/handler"
	"io/ioutil"
	"mime/multipart"
	"sort"
	"strings"
)

func NewOSSConfig() *OSSConfig {
	return &OSSConfig{
		EndPoint:        config.GetNaCosString("oss.endPoint", ""),
		AccessKeyId:     config.GetNaCosString("oss.accessKeyId", ""),
		AccessKeySecret: config.GetNaCosString("oss.accessKeySecret", ""),
		Bucket:          config.GetNaCosString("oss.bucket", ""),
		Host:            config.GetNaCosString("oss.host", ""),
		RegionId:        config.GetNaCosString("oss.regionId", ""),
		RoleArn:         config.GetNaCosString("oss.roleArn", ""),
	}
}

type OSSConfig struct {
	EndPoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Host            string
	RegionId        string
	RoleArn         string
}

func getOSSObjectKey(dir, fileName string) string {
	filenameArray := strings.Split(fileName, ".")
	newFileName := uuid.NewV4().String()
	if len(filenameArray) > 1 {
		newFileName = fmt.Sprintf("%s.%s", uuid.NewV4().String(), filenameArray[len(filenameArray)-1])
	}
	return fmt.Sprintf("%s/%s", dir, newFileName)
}

func UploadOSSFile(dir, fileName string, file *multipart.File, config *OSSConfig, isPublic bool) (string, error) {
	client, err := oss.New(config.EndPoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		return "", err
	}
	objectKey := getOSSObjectKey(dir, fileName)
	// 上传文件流。
	err = bucket.PutObject(objectKey, *file)
	if err != nil {
		return "", err
	}
	// 设置文件的访问权限。
	objectACL := oss.ACLPrivate
	if isPublic {
		objectACL = oss.ACLPublicRead
	}
	err = bucket.SetObjectACL(objectKey, objectACL)
	if err != nil {
		return "", err
	}
	if isPublic {
		return fmt.Sprintf("%s/%s", config.Host, objectKey), nil
	}
	return objectKey, nil
}

func UploadOSSFileBytes(dir, fileName string, file []byte, config *OSSConfig, isPublic bool) (string, error) {
	client, err := oss.New(config.EndPoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		return "", err
	}
	objectKey := getOSSObjectKey(dir, fileName)
	// 上传文件流。
	err = bucket.PutObject(objectKey, bytes.NewReader(file))
	if err != nil {
		return "", err
	}
	// 设置文件的访问权限。
	objectACL := oss.ACLPrivate
	if isPublic {
		objectACL = oss.ACLPublicRead
	}
	err = bucket.SetObjectACL(objectKey, objectACL)
	if err != nil {
		return "", err
	}
	if isPublic {
		return fmt.Sprintf("%s/%s", config.Host, objectKey), nil
	}
	return objectKey, nil
}

func UploadDingtalkOSSFile(endpoint, accessKeyId, accessKeySecret, accessToken, fileName, bucketName string, file []byte) error {
	client, err := oss.New("https://"+endpoint, accessKeyId, accessKeySecret, oss.SecurityToken(accessToken))
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	// 上传文件流。
	err = bucket.PutObject(fileName, bytes.NewReader(file))
	if err != nil {
		return err
	}
	return nil
}

func GetOSSPrivateFile(path string, config *OSSConfig) (string, error) {
	client, err := sts.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = config.RoleArn
	request.RoleSessionName = uuid.NewV4().String()
	request.DurationSeconds = "900"

	response, err := client.AssumeRole(request)
	if err != nil {
		return "", err
	}
	// 获取STS临时凭证后，您可以通过其中的安全令牌（SecurityToken）和临时访问密钥（AccessKeyId和AccessKeySecret）生成OSSClient。
	ossClient, err := oss.New(config.EndPoint, response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, oss.SecurityToken(response.Credentials.SecurityToken))
	if err != nil {
		return "", err
	}
	bucket, err := ossClient.Bucket(config.Bucket)
	if err != nil {
		return "", err
	}
	signedURL, err := bucket.SignURL(path, oss.HTTPGet, 600)
	if err != nil {
		return "", err
	}
	return signedURL, nil
}

func GetOSSWriteToken(ossConfig *OSSConfig) (string, string, string, error) {
	roleArn := config.GetNaCosString("oss.roleArn", "")
	client, err := sts.NewClientWithAccessKey(ossConfig.RegionId, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = roleArn
	request.RoleSessionName = uuid.NewV4().String()
	request.DurationSeconds = "900"

	response, err := client.AssumeRole(request)
	if err != nil {
		return "", "", "", err
	}
	return response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, response.Credentials.SecurityToken, nil
}

func ListOSSFile(config *OSSConfig, path string) ([]oss.ObjectProperties, error) {
	client, err := oss.New(config.EndPoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	bucket, err := client.Bucket(config.Bucket)
	marker := oss.Marker("")
	prefix := oss.Prefix(path)
	var result []oss.ObjectProperties
	for {
		lor, err := bucket.ListObjects(marker, prefix)
		if err != nil {
			return nil, handler.HandleError(err)
		}
		result = append(result, lor.Objects...)
		if lor.IsTruncated {
			prefix = oss.Prefix(lor.Prefix)
			marker = oss.Marker(lor.NextMarker)
		} else {
			break
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastModified.Before(result[j].LastModified)
	})
	return result, nil
}

func GetOSSFileBytes(config *OSSConfig, key string) ([]byte, error) {
	client, err := oss.New(config.EndPoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		fmt.Println(err)
	}
	bucket, err := client.Bucket(config.Bucket)
	body, err := bucket.GetObject(key)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, handler.HandleError(err)
	}
	return data, nil
}
