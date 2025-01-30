package wordcloud

// Token は形態素解析結果のトークン
type Token struct {
	Surface  string `json:"surface"`   // 表層形
	BaseForm string `json:"base_form"` // 基本形
	POS      string `json:"pos"`       // 品詞
}

// WordCount は単語の出現回数情報
type WordCount struct {
	Text     string `json:"text"`     // 単語
	Count    int    `json:"count"`    // 出現回数
	FontSize int    `json:"fontSize"` // フォントサイズ
	Color    string `json:"color"`    // 色
}

// Config はワードクラウドの設定
type Config struct {
	MinCount    int    // 最小出現回数
	MaxWords    int    // 最大単語数
	MinFontSize int    // 最小フォントサイズ
	MaxFontSize int    // 最大フォントサイズ
	ColorScheme string // 色スキーム
	Width       int    // 画像の幅
	Height      int    // 画像の高さ
}

// defaultConfig はデフォルト設定を返す
func defaultConfig() Config {
	return Config{
		MinCount:    2,
		MaxWords:    100,
		MinFontSize: 12,
		MaxFontSize: 48,
		ColorScheme: "blue",
		Width:       800,
		Height:      600,
	}
}

// Option は設定オプション関数の型
type Option func(*Analyzer)

// WithStopWords はストップワードを設定するオプション
func WithStopWords(words []string) Option {
	return func(a *Analyzer) {
		for _, word := range words {
			a.stopWords[word] = true
		}
	}
}

// targetPOS は対象とする品詞
var targetPOS = map[string]bool{
	"名詞":  true,
	"動詞":  true,
	"形容詞": true,
}

// defaultStopWords はデフォルトのストップワード
func defaultStopWords() map[string]bool {
	return map[string]bool{
		"する": true, "ある": true, "いる": true,
		"なる": true, "できる": true, "思う": true,
		"の": true, "が": true, "を": true,
		"は": true, "に": true, "で": true,
		"と": true, "た": true, "です": true,
		"ます": true, "これ": true, "それ": true,
		"この": true, "その": true, "あの": true,
	}
}
