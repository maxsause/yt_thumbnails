package client

import (
	"context"
	"fmt"
	"main/proto"
	"os"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client proto.ThumbnailServiceClient
}

func NewClient(serverAddr string) (*Client, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error connection to the gRPC server: %w", err)
	}

	client := proto.NewThumbnailServiceClient(conn)
	return &Client{conn: conn, client: client}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetThumbnail(videoURL string) (*ThumbnailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &proto.ThumbnailRequest{VideoUrl: videoURL}
	res, err := c.client.GetThumbnail(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error getting thumbnail: %w", err)
	}
	response := &ThumbnailResponse{
		Data:    res.GetData(),
		VideoID: res.GetVideoId(),
	}
	return response, nil
}

type ThumbnailResponse struct {
	Data    []byte
	VideoID string
}

func (c *Client) AsyncGetThumbnail(videoURLs []string) map[string]ThumbnailResult {
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	results := make(map[string]ThumbnailResult)

	for _, videoURL := range videoURLs {
		wg.Go(func() {
			response, err := c.GetThumbnail(videoURL)
			mu.Lock()
			results[videoURL] = ThumbnailResult{Response: response, Err: err}
			mu.Unlock()
		})
	}

	wg.Wait()
	return results
}

type ThumbnailResult struct {
	Response *ThumbnailResponse
	Err      error
}

func WriteFile(videoID string, data []byte, path string) error {
	name := videoID + ".jpg"
	fpath := filepath.Join(path, name)
	err := os.WriteFile(fpath, data, 0644)
	if err != nil {
		return fmt.Errorf("write thumbnail %q: %w", fpath, err)
	}
	return nil
}
