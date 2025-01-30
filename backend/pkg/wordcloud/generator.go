package wordcloud

import (
	"fmt"
	"log"
	"math"
	"sort"
)

// Generator はワードクラウドのデータを生成する構造体
type Generator struct {
	config   Config
	analyzer *Analyzer
}

// NewGenerator は新しいGeneratorを作成
func NewGenerator(config Config, analyzer *Analyzer) *Generator {
	if analyzer == nil {
		var err error
		analyzer, err = NewAnalyzer()
		if err != nil {
			panic(err) // 実際のアプリケーションではエラーハンドリングを適切に行う
		}
	}

	return &Generator{
		config:   config,
		analyzer: analyzer,
	}
}

// Generate はテキストからワードクラウドデータを生成
func (g *Generator) Generate(texts []string) ([]WordCount, error) {
	log.Printf("テキスト解析を開始します（%d件）...", len(texts))

	// 単語のカウント
	wordCounts := make(map[string]int)
	for i, text := range texts {
		tokens := g.analyzer.Analyze(text)
		for _, token := range tokens {
			wordCounts[token.BaseForm]++
		}

		// 1000件ごとに進捗を表示
		if (i+1)%1000 == 0 {
			log.Printf("テキスト解析中... %d/%d件完了", i+1, len(texts))
		}
	}

	log.Printf("単語の出現回数集計が完了しました。%d個の一意な単語が見つかりました。", len(wordCounts))

	// WordCountのスライスに変換
	var counts []WordCount
	for word, count := range wordCounts {
		if count >= g.config.MinCount {
			counts = append(counts, WordCount{
				Text:  word,
				Count: count,
			})
		}
	}

	// 出現回数でソート
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	// 最大単語数に制限
	if len(counts) > g.config.MaxWords {
		counts = counts[:g.config.MaxWords]
	}

	// フォントサイズと色を計算
	maxCount := counts[0].Count
	for i := range counts {
		counts[i].FontSize = g.calculateFontSize(counts[i].Count, maxCount)
		counts[i].Color = g.calculateColor(counts[i].Count, maxCount)
	}

	return counts, nil
}

// calculateFontSize は出現回数からフォントサイズを計算
func (g *Generator) calculateFontSize(count, maxCount int) int {
	ratio := math.Log(float64(count)) / math.Log(float64(maxCount))
	size := g.config.MinFontSize + int(ratio*float64(g.config.MaxFontSize-g.config.MinFontSize))

	if size < g.config.MinFontSize {
		return g.config.MinFontSize
	}
	if size > g.config.MaxFontSize {
		return g.config.MaxFontSize
	}
	return size
}

// calculateColor は出現回数から色を計算
func (g *Generator) calculateColor(count, maxCount int) string {
	ratio := float64(count) / float64(maxCount)

	switch g.config.ColorScheme {
	case "blue":
		intensity := uint8(150 + int(105*ratio))
		return fmt.Sprintf("#0000%02x", intensity)
	case "rainbow":
		hue := int(240 * ratio)
		return fmt.Sprintf("hsl(%d, 70%%, 50%%)", hue)
	default:
		return "#000000"
	}
}
