package oss

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var defaultOSS *minio.Client

func InitDefaultOSS() error {
	ossCnf := GetDefaultOSSConfig()
	if ossCnf == nil {
		return nil
	}
	minioClient, err := minio.New(ossCnf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ossCnf.AccessKeyID, ossCnf.SecretAccessKey, ossCnf.Token),
		Secure: ossCnf.UseSSL,
	})
	if err != nil {
		return err
	}

	ctx := context.Background()

	// 检查默认桶是否存在，不存在则创建
	exists, err := minioClient.BucketExists(ctx, ossCnf.BucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, ossCnf.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}

		// 设置桶的访问策略为公开读
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + ossCnf.BucketName + `/*"]
				}
			]
		}`
		err = minioClient.SetBucketPolicy(ctx, ossCnf.BucketName, policy)
		if err != nil {
			return err
		}
	}

	defaultOSS = minioClient

	return nil
}
func GetDefaultOSS() *minio.Client {
	return defaultOSS
}
