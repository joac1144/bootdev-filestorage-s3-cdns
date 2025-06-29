package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}
	videoUrlSlice := strings.Split(*video.VideoURL, ",")
	if len(videoUrlSlice) != 2 {
		return video, nil
	}
	bucket, key := videoUrlSlice[0], videoUrlSlice[1]
	presignedUrl, err := generatePresignedURL(cfg.s3Client, bucket, key, 3*time.Minute)
	if err != nil {
		return video, fmt.Errorf("failed to generate presigned URL for video %s: %w", video.ID, err)
	}
	video.VideoURL = &presignedUrl
	return video, nil
}

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)
	presignedHttpReq, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedHttpReq.URL, nil
}
