package wordcloud

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
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

// Rectangle は単語の配置領域を表す
type Rectangle struct {
	X, Y, W, H float64
}

// Overlaps は2つの矩形が重なっているかを判定
func (r Rectangle) Overlaps(other Rectangle) bool {
	return !(r.X+r.W < other.X ||
		other.X+other.W < r.X ||
		r.Y+r.H < other.Y ||
		other.Y+other.H < r.Y)
}

// ExportPNG はワードクラウドデータをPNG画像として出力
func (fp *FileProcessor) ExportPNG(data []WordCount, outputPath string) error {
	// デバッグ用のログ追加
	log.Printf("処理開始: 単語数=%d", len(data))

	dc := gg.NewContext(fp.config.Width, fp.config.Height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// フォント設定
	fontPath := "/Library/Fonts/Arial Unicode.ttf"
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return fmt.Errorf("フォントファイルの読み込みに失敗: %w", err)
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return fmt.Errorf("フォントのパースに失敗: %w", err)
	}

	// 頻出度の最大値と最小値を取得（'*'を除外）
	maxCount := 0
	minCount := math.MaxInt32
	for _, word := range data {
		if word.Text == "*" {
			continue // '*'は除外
		}
		if word.Count > maxCount {
			maxCount = word.Count
		}
		if word.Count < minCount {
			minCount = word.Count
		}
	}

	log.Printf("頻出度範囲: 最小=%d, 最大=%d", minCount, maxCount)

	// 色を決定する関数
	getColorHex := func(count int, text string) string {
		if text == "*" {
			return "#66CCFF" // '*'は最も薄い青に固定
		}

		// 頻出度に基づいて0-1の値を計算
		ratio := float64(count-minCount) / float64(maxCount-minCount)

		// 頻出度に応じて5段階の色を返す
		switch {
		case ratio >= 0.8:
			return "#FF0000" // 赤（最も頻出）
		case ratio >= 0.6:
			return "#FF6600" // オレンジ
		case ratio >= 0.4:
			return "#0066FF" // 青
		case ratio >= 0.2:
			return "#3399FF" // 明るい青
		default:
			return "#66CCFF" // 最も薄い青
		}
	}

	// 配置済みの単語の領域を管理するスライスを初期化
	occupied := make([]Rectangle, 0)

	// 配置パラメータの調整
	centerX := float64(fp.config.Width) / 2
	centerY := float64(fp.config.Height) / 2
	maxRadius := math.Min(float64(fp.config.Width), float64(fp.config.Height)) * 0.4
	spiralDelta := 0.1 // スパイラルの増加率を小さくする

	// 単語を描画
	for _, word := range data {
		fontSize := float64(word.FontSize)
		face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
		dc.SetFontFace(face)
		dc.SetHexColor(getColorHex(word.Count, word.Text))

		// 単語の大きさを計算
		w, h := dc.MeasureString(word.Text)

		// スパイラル状に配置を試行
		placed := false
		for radius := float64(0); radius < maxRadius; radius += spiralDelta {
			for angle := float64(0); angle < 2*math.Pi*radius; angle += spiralDelta {
				x := centerX + math.Cos(angle)*radius - w/2
				y := centerY + math.Sin(angle)*radius + h/2

				// 画像の範囲内かチェック
				if x < 0 || x+w > float64(fp.config.Width) ||
					y-h < 0 || y > float64(fp.config.Height) {
					continue
				}

				// 配置領域の作成
				newRect := Rectangle{
					X: x - 2,     // マージンを追加
					Y: y - h - 2, // マージンを追加
					W: w + 4,     // マージンを追加
					H: h + 4,     // マージンを追加
				}

				// 重なりチェック
				overlap := false
				for _, rect := range occupied {
					if newRect.Overlaps(rect) {
						overlap = true
						break
					}
				}

				if !overlap {
					dc.DrawString(word.Text, x, y)
					occupied = append(occupied, newRect)
					placed = true
					break
				}
			}
			if placed {
				break
			}
		}

		if !placed {
			log.Printf("警告: '%s' の配置に失敗しました", word.Text)
		}
	}

	// PNG画像として保存
	if err := dc.SavePNG(outputPath); err != nil {
		return fmt.Errorf("PNG画像の保存に失敗: %w", err)
	}

	log.Printf("処理完了: %s", outputPath)
	return nil
}
