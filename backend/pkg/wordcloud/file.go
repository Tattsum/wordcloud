package wordcloud

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// FileProcessor はファイル処理を行う構造体
type FileProcessor struct {
	generator *Generator
	config    Config
}

// NewFileProcessor は新しいFileProcessorを作成
func NewFileProcessor(config Config) (*FileProcessor, error) {
	analyzer, err := NewAnalyzer()
	if err != nil {
		return nil, fmt.Errorf("アナライザーの初期化に失敗: %w", err)
	}

	generator := NewGenerator(config, analyzer)

	return &FileProcessor{
		generator: generator,
		config:    config,
	}, nil
}

// ProcessCSV はCSVファイルを処理してワードクラウドデータを生成
func (fp *FileProcessor) ProcessCSV(inputPath string, messageColumn int) ([]WordCount, error) {
	log.Printf("CSVファイル '%s' の処理を開始します...", inputPath)

	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("入力ファイルのオープンに失敗: %w", err)
	}
	defer file.Close()

	// 総行数を数える
	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("行数カウント中にエラー: %w", err)
	}

	// ファイルポインタを先頭に戻す
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("ファイルポインタのリセットに失敗: %w", err)
	}

	reader := csv.NewReader(file)

	// ヘッダーをスキップ
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("ヘッダーの読み込みに失敗: %w", err)
	}
	lineCount-- // ヘッダー行を除く

	// メッセージテキストを収集
	var texts []string
	processedLines := 0
	lastProgress := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("レコードの読み込みに失敗: %w", err)
		}

		processedLines++
		progress := (processedLines * 100) / lineCount

		// 10%単位で進捗を表示
		if progress/10 > lastProgress/10 {
			log.Printf("CSVファイルの処理中... %d%%完了", progress)
			lastProgress = progress
		}

		if len(record) > messageColumn {
			texts = append(texts, record[messageColumn])
		}
	}

	log.Printf("CSVファイルの読み込みが完了しました。%d行を処理しました。", processedLines)
	log.Print("ワードクラウドデータの生成を開始します...")

	// ワードクラウドデータの生成
	return fp.generator.Generate(texts)
}

// ExportJSON はワードクラウドデータをJSONファイルに出力
func (fp *FileProcessor) ExportJSON(data []WordCount, outputPath string) error {
	// 出力ディレクトリの作成
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("出力ディレクトリの作成に失敗: %w", err)
	}

	// JSONファイルに出力
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("出力ファイルの作成に失敗: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("JSONの書き込みに失敗: %w", err)
	}

	return nil
}
