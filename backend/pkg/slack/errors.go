package slack

import "errors"

var (
	// ErrInvalidToken は無効なトークンエラー
	ErrInvalidToken = errors.New("無効なSlackトークン")

	// ErrChannelNotFound はチャンネルが見つからないエラー
	ErrChannelNotFound = errors.New("チャンネルが見つかりません")

	// ErrNotChannel は指定されたIDがチャンネルでないエラー
	ErrNotChannel = errors.New("指定されたIDはチャンネルではありません")

	// ErrBotNotInChannel はBotがチャンネルに参加していないエラー
	ErrBotNotInChannel = errors.New("Botがチャンネルに参加していません")

	// ErrRateLimitExceeded はレートリミット超過エラー
	ErrRateLimitExceeded = errors.New("APIレートリミットを超過しました")
)

// IsNotFoundError はチャンネルが見つからないエラーかを判定
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrChannelNotFound)
}

// IsRateLimitError はレートリミットエラーかを判定
func IsRateLimitError(err error) bool {
	return errors.Is(err, ErrRateLimitExceeded)
}
