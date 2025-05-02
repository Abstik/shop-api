package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"go.uber.org/zap"

	"mxshop-api/oss-web/global"
)

type PolicyToken struct {
	Host             string `json:"host"`
	Dir              string `json:"dir"`
	Policy           string `json:"policy"`
	Signature        string `json:"signature"`
	SignatureVersion string `json:"x_oss_signature_version"`
	Credential       string `json:"x_oss_credential"`
	Date             string `json:"x_oss_date"`
	SecurityToken    string `json:"security_token"`
	Callback         string `json:"callback"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

func GetPolicyToken() string {
	region := global.ServerConfig.OssInfo.Region
	bucketName := global.ServerConfig.OssInfo.BucketName
	product := "oss"

	// 如果配置中已有完整host，则使用它，否则构建标准OSS地址
	host := global.ServerConfig.OssInfo.Host
	if host == "" {
		host = fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", bucketName, region)
	}

	// 使用配置中的上传目录
	dir := global.ServerConfig.OssInfo.UploadDir

	// 使用配置中的回调URL
	callbackUrl := global.ServerConfig.OssInfo.CallBackUrl

	// 创建凭证配置
	config := new(credentials.Config).
		SetType("ram_role_arn").
		SetAccessKeyId(global.ServerConfig.OssInfo.ApiKey).
		SetAccessKeySecret(global.ServerConfig.OssInfo.ApiSecret).
		SetRoleArn(global.ServerConfig.OssInfo.RoleArn).
		SetRoleSessionName("Role_Session_Name").
		SetPolicy("").
		SetRoleSessionExpiration(3600)

	// 根据配置创建凭证提供器
	provider, err := credentials.NewCredential(config)
	if err != nil {
		zap.S().Errorf("创建凭证提供器失败: %v", err)
	}

	// 从凭证提供器获取凭证
	cred, err := provider.GetCredential()
	if err != nil {
		zap.S().Errorf("获取凭证失败: %v", err)
	}

	// 构建policy
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	expiration := utcTime.Add(1 * time.Hour)
	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": bucketName},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, region, product)},
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			map[string]string{"x-oss-security-token": *cred.SecurityToken},
			[]string{"starts-with", "$key", dir},
		},
	}

	// 将policy转换为JSON格式
	policy, err := json.Marshal(policyMap)
	if err != nil {
		zap.S().Errorf("序列化Policy失败: %v", err)
	}

	// 构造待签名字符串(StringToSign)
	stringToSign := base64.StdEncoding.EncodeToString(policy)

	// 构建签名
	signature, err := generateSignature(stringToSign, *cred.AccessKeySecret, date, region, product)
	if err != nil {
		zap.S().Errorf("生成签名失败: %v", err)
	}

	// 构建回调参数
	callbackParam := CallbackParam{
		CallbackUrl:      callbackUrl,
		CallbackBody:     "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}",
		CallbackBodyType: "application/x-www-form-urlencoded",
	}

	callbackStr, err := json.Marshal(callbackParam)
	if err != nil {
		zap.S().Errorf("序列化回调参数失败: %v", err)
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)

	// 构建返回给前端的表单
	policyToken := PolicyToken{
		Policy:           stringToSign,
		SecurityToken:    *cred.SecurityToken,
		SignatureVersion: "OSS4-HMAC-SHA256",
		Credential:       fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, region, product),
		Date:             utcTime.Format("20060102T150405Z"),
		Signature:        signature,
		Host:             host,
		Dir:              dir,
		Callback:         callbackBase64,
	}

	response, err := json.Marshal(policyToken)
	if err != nil {
		zap.S().Errorf("序列化响应失败: %v", err)
	}

	return string(response)
}

// 生成OSS签名
func generateSignature(stringToSign, secretKey, date, region, product string) (string, error) {
	hmacHash := func() hash.Hash { return sha256.New() }

	// 构建signing key
	signingKey := "aliyun_v4" + secretKey
	h1 := hmac.New(hmacHash, []byte(signingKey))
	io.WriteString(h1, date)
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	io.WriteString(h2, region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	io.WriteString(h3, product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	io.WriteString(h4, "aliyun_v4_request")
	h4Key := h4.Sum(nil)

	// 生成签名
	h := hmac.New(hmacHash, h4Key)
	io.WriteString(h, stringToSign)
	signature := hex.EncodeToString(h.Sum(nil))

	return signature, nil
}
