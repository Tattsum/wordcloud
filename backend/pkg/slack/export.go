package slack

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ExportOptions はエクスポートのオプション
type ExportOptions struct {
	OutputDir     string
	IncludeThread bool
	TimeLocation  *time.Location
}

// defaultExportOptions はデフォルトのエクスポートオプションを返す
func defaultExportOptions() *ExportOptions {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return &ExportOptions{
		OutputDir:     "data",
		IncludeThread: true,
		TimeLocation:  jst,
	}
}

// ExportOption はエクスポートのオプション関数の型
type ExportOption func(*ExportOptions)

// WithOutputDir は出力ディレクトリを指定するオプション
func WithOutputDir(dir string) ExportOption {
	return func(opts *ExportOptions) {
		opts.OutputDir = dir
	}
}

// WithoutThread はスレッド情報を除外するオプション
func WithoutThread() ExportOption {
	return func(opts *ExportOptions) {
		opts.IncludeThread = false
	}
}

// ExportChannelMessages はチャンネルのメッセージをCSVに出力
func (c *Client) ExportChannelMessages(channelID string, options ...ExportOption) (string, error) {
	log.Println("メッセージのエクスポートを開始します")

	// オプションの設定
	opts := defaultExportOptions()
	for _, opt := range options {
		opt(opts)
	}

	// チャンネル情報の取得
	log.Println("チャンネル情報を取得中...")
	channel, err := c.GetChannelInfo(channelID)
	if err != nil {
		return "", fmt.Errorf("チャンネル情報の取得に失敗: %w", err)
	}
	log.Printf("チャンネル情報を取得しました: %s", channel.Name)

	// 出力ディレクトリの作成
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗: %w", err)
	}

	// CSVファイルの作成
	filename := fmt.Sprintf("messages_%s_%s.csv",
		channel.Name,
		time.Now().In(opts.TimeLocation).Format("20060102_150405"),
	)
	filepath := filepath.Join(opts.OutputDir, filename)

	log.Printf("CSVファイルを作成します: %s", filepath)
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("CSVファイルの作成に失敗: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーの書き込み
	headers := []string{"Timestamp", "UserID", "Username", "Message"}
	if opts.IncludeThread {
		headers = append(headers, "ThreadTS")
	}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("ヘッダーの書き込みに失敗: %w", err)
	}

	// メッセージの取得と書き込み
	log.Println("メッセージを取得中...")
	messages, err := c.GetChannelMessages(channelID, WithUserInfo())
	if err != nil {
		return "", fmt.Errorf("メッセージの取得に失敗: %w", err)
	}

	log.Printf("CSVファイルに %d 件のメッセージを書き込み中...", len(messages))
	for i, msg := range messages {
		if i > 0 && i%1000 == 0 {
			log.Printf("進捗: %d/%d 件処理完了", i, len(messages))
		}

		record := []string{
			msg.Timestamp,
			msg.UserID,
			msg.Username,
			msg.Text,
		}
		if opts.IncludeThread {
			record = append(record, msg.ThreadTS)
		}

		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("レコードの書き込みに失敗: %w", err)
		}
	}

	log.Printf("メッセージのエクスポートが完了しました: %s", filepath)
	return filepath, nil
}
