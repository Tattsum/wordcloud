package slack

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/slack-go/slack"
	"golang.org/x/sync/semaphore"
)

// SlackMessage はSlackのメッセージを表す構造体
type SlackMessage struct {
	ID        string         `json:"id"`
	Text      string         `json:"text"`
	UserID    string         `json:"user_id"`
	Timestamp string         `json:"timestamp"`
	Username  string         `json:"username"`
	ThreadTS  string         `json:"thread_ts"`
	Replies   []SlackMessage `json:"replies"`
}

// Client はSlack APIクライアントのラッパー構造体
type Client struct {
	api            *slack.Client
	rateLimit      time.Duration
	mutex          sync.Mutex
	lastCall       time.Time
	sem            *semaphore.Weighted
	maxConcurrency int
}

// ClientConfig はクライアントの設定オプション
type ClientConfig struct {
	Token          string
	RateLimit      time.Duration
	MaxConcurrency int
}

// NewClient は新しいSlackクライアントを作成
func NewClient(config ClientConfig) *Client {
	if config.RateLimit == 0 {
		config.RateLimit = time.Second
	}
	if config.MaxConcurrency <= 0 {
		config.MaxConcurrency = 5
	}

	return &Client{
		api:            slack.New(config.Token),
		rateLimit:      config.RateLimit,
		maxConcurrency: config.MaxConcurrency,
		sem:            semaphore.NewWeighted(int64(config.MaxConcurrency)),
	}
}

// waitForRateLimit はレートリミットを制御
func (c *Client) waitForRateLimit() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	if diff := now.Sub(c.lastCall); diff < c.rateLimit {
		time.Sleep(c.rateLimit - diff)
	}
	c.lastCall = time.Now()
}

func (c *Client) GetChannelMessages(channelID string, options ...MessageOption) ([]SlackMessage, error) {
	log.Printf("チャンネル %s のメッセージ取得を開始します", channelID)

	opts := defaultMessageOptions()
	for _, opt := range options {
		opt(opts)
	}

	var allMessages []SlackMessage
	cursor := ""

	for {
		// メッセージページを取得
		log.Printf("メッセージページを取得中...")
		params := &slack.GetConversationHistoryParameters{
			ChannelID: channelID,
			Cursor:    cursor,
			Limit:     100,
		}

		c.waitForRateLimit()
		history, err := c.api.GetConversationHistory(params)
		if err != nil {
			return nil, fmt.Errorf("メッセージの取得に失敗: %w", err)
		}

		log.Printf("%d 件のメッセージを取得しました。処理を開始します...", len(history.Messages))

		// メッセージを処理
		for _, msg := range history.Messages {
			message := SlackMessage{
				ID:        msg.Timestamp,
				Text:      msg.Text,
				UserID:    msg.User,
				Timestamp: msg.Timestamp,
				ThreadTS:  msg.ThreadTimestamp,
			}

			// スレッドの返信を取得
			if msg.ThreadTimestamp != "" && msg.ThreadTimestamp == msg.Timestamp {
				log.Printf("スレッド %s の返信を取得中...", msg.ThreadTimestamp)
				replies, err := c.getThreadReplies(channelID, msg.ThreadTimestamp)
				if err != nil {
					return nil, fmt.Errorf("スレッド返信の取得に失敗: %w", err)
				}
				message.Replies = replies
			}

			allMessages = append(allMessages, message)
		}

		// 次のページがなければ終了
		if !history.HasMore {
			break
		}
		cursor = history.ResponseMetaData.NextCursor
	}

	return allMessages, nil
}

func (c *Client) getThreadReplies(channelID, threadTS string) ([]SlackMessage, error) {
	var replies []SlackMessage
	cursor := ""

	for {
		c.waitForRateLimit()
		params := &slack.GetConversationRepliesParameters{
			ChannelID: channelID,
			Timestamp: threadTS,
			Cursor:    cursor,
		}

		messages, hasMore, nextCursor, err := c.api.GetConversationReplies(params)
		if err != nil {
			return nil, fmt.Errorf("スレッド返信の取得に失敗: %w", err)
		}

		// 最初のメッセージはスキップ（親メッセージのため）
		if len(messages) > 0 {
			messages = messages[1:]
		}

		for _, msg := range messages {
			reply := SlackMessage{
				ID:        msg.Timestamp,
				Text:      msg.Text,
				UserID:    msg.User,
				Timestamp: msg.Timestamp,
				ThreadTS:  msg.ThreadTimestamp,
			}
			replies = append(replies, reply)
		}

		if !hasMore {
			break
		}
		cursor = nextCursor
	}

	return replies, nil
}

// GetChannelInfo はチャンネル情報を取得
func (c *Client) GetChannelInfo(channelID string) (*Channel, error) {
	c.waitForRateLimit()

	info, err := c.api.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelID,
	})
	if err != nil {
		return nil, fmt.Errorf("チャンネル情報の取得に失敗: %w", err)
	}

	return &Channel{
		ID:          info.ID,
		Name:        info.Name,
		IsChannel:   info.IsChannel,
		IsPrivate:   info.IsPrivate,
		CreatorID:   info.Creator,
		Created:     time.Unix(int64(info.Created), 0),
		MemberCount: info.NumMembers,
	}, nil
}

// JoinChannel はBotをチャンネルに参加させる
func (c *Client) JoinChannel(channelID string) error {
	c.waitForRateLimit()

	channel, _, _, err := c.api.JoinConversation(channelID)
	if err != nil {
		return fmt.Errorf("チャンネルへの参加に失敗: %w", err)
	}

	if channel == nil {
		return fmt.Errorf("チャンネルの参加確認に失敗")
	}

	return nil
}

// Validate はトークンとBotの権限を検証
func (c *Client) Validate() error {
	c.waitForRateLimit()

	_, err := c.api.AuthTest()
	if err != nil {
		return fmt.Errorf("認証に失敗: %w", err)
	}
	return nil
}
