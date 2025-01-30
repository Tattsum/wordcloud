package slack

import "time"

// Message はSlackメッセージの構造体
type Message struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	Username  string `json:"username,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	Timestamp string `json:"timestamp"`
	ThreadTS  string `json:"thread_ts,omitempty"`
}

// Channel はSlackチャンネルの構造体
type Channel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	IsChannel   bool      `json:"is_channel"`
	IsPrivate   bool      `json:"is_private"`
	CreatorID   string    `json:"creator_id"`
	Created     time.Time `json:"created"`
	MemberCount int       `json:"member_count"`
}

// MessageOptions はメッセージ取得のオプション
type messageOptions struct {
	limit           int
	includeUserInfo bool
}

// MessageOption はメッセージ取得のオプション関数
type MessageOption func(*messageOptions)

// defaultMessageOptions はデフォルトのオプションを返す
func defaultMessageOptions() *messageOptions {
	return &messageOptions{
		limit:           0, // 0は制限なし
		includeUserInfo: false,
	}
}

// WithLimit は取得するメッセージ数を制限するオプション
func WithLimit(limit int) MessageOption {
	return func(opts *messageOptions) {
		opts.limit = limit
	}
}

// WithUserInfo はユーザー情報を含めるオプション
func WithUserInfo() MessageOption {
	return func(opts *messageOptions) {
		opts.includeUserInfo = true
	}
}
