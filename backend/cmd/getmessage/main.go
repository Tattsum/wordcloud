// cmd/getmessage/main.go
package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Tattsum/wordcloud/backend/pkg/slack"
)

func main() {
	token := flag.String("token", "", "Slack Bot User OAuth Token")
	channel := flag.String("channel", "", "Channel ID")
	output := flag.String("output", "data", "Output directory path")

	flag.Parse()

	// トークンとチャンネルIDの検証
	*token = strings.TrimSpace(*token)
	*channel = strings.TrimSpace(*channel)

	if *token == "" || *channel == "" {
		log.Println("Error: token and channel are required")
		flag.Usage()
		os.Exit(1)
	}

	// クライアントの初期化
	config := slack.ClientConfig{
		Token:          *token,
		RateLimit:      time.Second,
		MaxConcurrency: 5,
	}
	client := slack.NewClient(config)

	// トークンの検証
	if err := client.Validate(); err != nil {
		log.Fatalf("Slackトークンの検証に失敗: %v", err)
	}

	// メッセージの取得とCSV出力
	outputPath, err := client.ExportChannelMessages(
		*channel,
		slack.WithOutputDir(*output),
	)
	if err != nil {
		log.Fatalf("メッセージの出力に失敗: %v", err)
	}

	log.Printf("メッセージの出力が完了しました: %s", outputPath)
}
