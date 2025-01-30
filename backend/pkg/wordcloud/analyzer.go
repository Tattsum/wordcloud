package wordcloud

import (
	"strings"
	"sync"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// Analyzer は形態素解析を行う構造体
type Analyzer struct {
	tokenizer *tokenizer.Tokenizer
	stopWords map[string]bool
	mu        sync.Mutex
}

// NewAnalyzer は新しいAnalyzerを作成
func NewAnalyzer(options ...Option) (*Analyzer, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}

	a := &Analyzer{
		tokenizer: t,
		stopWords: defaultStopWords(),
	}

	// オプションを適用
	for _, opt := range options {
		opt(a)
	}

	return a, nil
}

// Analyze はテキストを解析して単語のスライスを返す
func (a *Analyzer) Analyze(text string) []Token {
	a.mu.Lock()
	tokens := a.tokenizer.Tokenize(text)
	a.mu.Unlock()

	var results []Token
	for _, t := range tokens {
		features := t.Features()
		if len(features) < 7 {
			continue
		}

		pos := features[0]      // 品詞
		baseForm := features[6] // 基本形

		// 絵文字やスラックの特殊表記を処理
		if strings.HasPrefix(baseForm, ":") && strings.HasSuffix(baseForm, ":") {
			results = append(results, Token{
				Surface:  baseForm,
				BaseForm: baseForm,
				POS:      "絵文字",
			})
			continue
		}

		// URLは除外
		if strings.HasPrefix(baseForm, "http") {
			continue
		}

		// ユーザーメンション (@User) は匿名化して保持
		if strings.HasPrefix(baseForm, "<@") {
			results = append(results, Token{
				Surface:  "某メンバー",
				BaseForm: "某メンバー",
				POS:      "固有名詞",
			})
			continue
		}

		if !a.isStopWord(baseForm) && a.isTargetPOS(pos) {
			results = append(results, Token{
				Surface:  t.Surface,
				BaseForm: baseForm,
				POS:      pos,
			})
		}
	}

	return results
}

// isStopWord は単語がストップワードかどうかを判定
func (a *Analyzer) isStopWord(word string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.stopWords[word]
}

// isTargetPOS は品詞が対象かどうかを判定
func (a *Analyzer) isTargetPOS(pos string) bool {
	return targetPOS[pos]
}

// AddStopWords はストップワードを追加
func (a *Analyzer) AddStopWords(words ...string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, word := range words {
		a.stopWords[word] = true
	}
}
