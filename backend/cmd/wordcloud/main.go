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
		outputFile  = flag.String("output", "", "Output JSON file path")
		minCount    = flag.Int("min-count", 2, "Minimum word count")
		maxWords    = flag.Int("max-words", 100, "Maximum number of words")
		colorScheme = flag.String("color", "blue", "Color scheme (blue/rainbow)")
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

	// JSONファイルに出力
	if err := processor.ExportJSON(wordCounts, *outputFile); err != nil {
		log.Fatalf("JSONファイルの出力に失敗: %v", err)
	}

	log.Printf("ワードクラウドデータの生成が完了しました: %s", *outputFile)
}
