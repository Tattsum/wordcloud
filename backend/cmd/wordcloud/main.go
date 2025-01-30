package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/Tattsum/wordcloud/backend/pkg/wordcloud"
)

func main() {
	var (
		inputFile   = flag.String("input", "", "Input CSV file path")
		outputFile  = flag.String("output", "", "Output PNG file path")
		minCount    = flag.Int("min-count", 2, "Minimum word count")
		maxWords    = flag.Int("max-words", 100, "Maximum number of words")
		colorScheme = flag.String("color", "blue", "Color scheme (blue/rainbow)")
		width       = flag.Int("width", 800, "Image width in pixels")
		height      = flag.Int("height", 600, "Image height in pixels")
	)

	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 設定の初期化
	config := wordcloud.Config{
		MinCount:    *minCount,
		MaxWords:    *maxWords,
		MinFontSize: 12,
		MaxFontSize: 48,
		ColorScheme: *colorScheme,
		Width:       *width,
		Height:      *height,
	}

	// プロセッサーの初期化
	processor, err := wordcloud.NewFileProcessor(config)
	if err != nil {
		log.Fatalf("プロセッサーの初期化に失敗: %v", err)
	}

	// CSVファイルの処理
	wordCounts, err := processor.ProcessCSV(*inputFile, 3) // 3はメッセージカラムのインデックス
	if err != nil {
		log.Fatalf("CSVファイルの処理に失敗: %v", err)
	}

	// 出力ディレクトリの作成
	if err := os.MkdirAll(filepath.Dir(*outputFile), 0755); err != nil {
		log.Fatalf("出力ディレクトリの作成に失敗: %v", err)
	}

	// PNGファイルとして出力
	if err := processor.ExportPNG(wordCounts, *outputFile); err != nil {
		log.Fatalf("PNG画像の出力に失敗: %v", err)
	}

	log.Printf("ワードクラウド画像の生成が完了しました: %s", *outputFile)
}
