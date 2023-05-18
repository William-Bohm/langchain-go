package documentLoaders

import (
	"github.com/William-Bohm/langchain-go/langchain-go/documentStore/documentSchema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type S3DirectoryLoader struct {
	bucket string
	prefix string
}

type S3FileLoader struct {
	bucket string
	key    string
}

func NewS3DirectoryLoader(bucket string, prefix string) *S3DirectoryLoader {
	return &S3DirectoryLoader{
		bucket: bucket,
		prefix: prefix,
	}
}

func (s *S3DirectoryLoader) Load() []documentSchema.Document {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))

	svc := s3.New(sess)
	var docs []documentSchema.Document

	err := svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(s.prefix),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			loader := NewS3FileLoader(s.bucket, *obj.Key)
			docs = append(docs, loader.Load()...)
		}
		return true
	})

	if err != nil {
		panic(err)
	}

	return docs
}

func NewS3FileLoader(bucket string, key string) *S3FileLoader {
	return &S3FileLoader{
		bucket: bucket,
		key:    key,
	}
}

func (s *S3FileLoader) Load() []documentSchema.Document {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))

	svc := s3.New(sess)

	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, s.key)
	_ = os.MkdirAll(filepath.Dir(filePath), 0755)

	_, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
	})

	if err != nil {
		panic(err)
	}

	loader := NewUnstructuredFileLoader(filePath, "", map[string]interface{}{})

	docs, err := loader.Load()
	return docs
}
