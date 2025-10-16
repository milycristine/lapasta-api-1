package s3client

import (
	"bytes"
	"fmt"
	"lapasta/config"
	utils "lapasta/internal/Utils"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	s3Session *s3.S3
	bucket    string
}

func gerarNomeImagem(extensao string) string {
	chaveAleatoria := utils.GerarStringAleatoria(12)
	return fmt.Sprintf("%s%s", chaveAleatoria, extensao)
}

func NovoS3Client(region, bucket string) *S3Client {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			config.Yml.AWS.AccessKeyID,
			config.Yml.AWS.SecretAccessKey,
			""),
	}))

	return &S3Client{
		s3Session: s3.New(sess),
		bucket:    bucket,
	}
}

func (c *S3Client) UploadBase64File(fileBytes []byte) (string, error) {

	extensao := filepath.Ext(".png")
	fileName := gerarNomeImagem(extensao)

	fileExt := filepath.Ext(fileName)
	if fileExt == "" {
		return "", fmt.Errorf("extensão de arquivo inválida")
	}

	_, err := c.s3Session.PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(c.bucket),
		Key:                aws.String(fileName),
		Body:               bytes.NewReader(fileBytes),
		ContentLength:      aws.Int64(int64(len(fileBytes))),
		ContentType:        aws.String(http.DetectContentType(fileBytes)),
		ContentDisposition: aws.String("inline"),
	})
	if err != nil {
		return "", fmt.Errorf("erro ao enviar o arquivo ao S3: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.bucket, fileName)

	return url, nil
}
