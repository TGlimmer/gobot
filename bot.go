package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TGlimmer/gobot/models"
)

const (
	defaultPollTimeout      = time.Minute
	defaultUpdatesChanCap   = 1024
	defaultCheckInitTimeout = time.Second * 5
	defaultWorkers          = 1
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ErrorsHandler func(err error)
type DebugHandler func(format string, args ...any)
type Middleware func(next HandlerFunc) HandlerFunc
type HandlerFunc func(ctx context.Context, bot *Bot, update *models.Update)
type MatchFunc func(update *models.Update) bool

// Bot represents Telegram Bot main object
type Bot struct {
	lastUpdateID int64

	url                string
	token              string
	pollTimeout        time.Duration
	skipGetMe          bool
	webhookSecretToken string
	testEnvironment    bool
	workers            int
	notAsyncHandlers   bool

	defaultHandlerFunc HandlerFunc

	errorsHandler ErrorsHandler
	debugHandler  DebugHandler

	middlewares []Middleware

	handlersMx sync.RWMutex
	handlers   []handler

	client           HttpClient
	isDebug          bool
	checkInitTimeout time.Duration

	allowedUpdates AllowedUpdates

	updates chan *models.Update
}

// New creates new Bot instance
func New(token string, options ...Option) (*Bot, error) {
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("empty token")
	}

	b := &Bot{
		url:         "https://api.telegram.org",
		token:       token,
		pollTimeout: defaultPollTimeout,
		client: &http.Client{
			Timeout: defaultPollTimeout,
		},
		defaultHandlerFunc: defaultHandler,
		errorsHandler:      defaultErrorsHandler,
		debugHandler:       defaultDebugHandler,
		checkInitTimeout:   defaultCheckInitTimeout,
		workers:            defaultWorkers,

		updates: make(chan *models.Update, defaultUpdatesChanCap),
	}

	for _, o := range options {
		o(b)
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.checkInitTimeout)
	defer cancel()

	if !b.skipGetMe {
		_, err := b.GetMe(ctx)
		if err != nil {
			return nil, fmt.Errorf("error call getMe, %w", err)
		}
	}

	return b, nil
}

// ID returns the bot ID parsed from the token prefix. It returns 0 if the
// token is malformed; this preserves the original signature for callers that
// don't care about diagnostics. Use IDOrErr if you need to distinguish a real
// id from a parse failure (issue #277).
func (b *Bot) ID() int64 {
	id, _ := b.IDOrErr()
	return id
}

// IDOrErr returns the bot ID parsed from the token prefix together with a
// non-nil error if the token is empty or does not start with a numeric id
// followed by ':'. This is the variant to use when the caller needs to react
// to misconfiguration instead of silently treating the bot as id 0.
func (b *Bot) IDOrErr() (int64, error) {
	if b.token == "" {
		return 0, fmt.Errorf("bot token is empty")
	}
	prefix, _, ok := strings.Cut(b.token, ":")
	if !ok {
		return 0, fmt.Errorf("malformed bot token: missing ':' separator")
	}
	id, err := strconv.ParseInt(prefix, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("malformed bot token: id prefix is not numeric: %w", err)
	}
	return id, nil
}

// SetToken sets the bot token
func (b *Bot) SetToken(token string) {
	b.token = token
}

// Token returns the bot token
func (b *Bot) Token() string {
	return b.token
}

// StartWebhook starts the Bot with webhook mode
func (b *Bot) StartWebhook(ctx context.Context) {
	wg := sync.WaitGroup{}

	wg.Add(b.workers)
	for i := 0; i < b.workers; i++ {
		go b.waitUpdates(ctx, &wg)
	}

	wg.Wait()
}

// Start the bot
func (b *Bot) Start(ctx context.Context) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go b.getUpdates(ctx, &wg)

	wg.Add(b.workers)
	for i := 0; i < b.workers; i++ {
		go b.waitUpdates(ctx, &wg)
	}

	wg.Wait()
}

func defaultErrorsHandler(err error) {
	log.Printf("[TGBOT] [ERROR] %v", err)
}

func defaultDebugHandler(format string, args ...any) {
	log.Printf("[TGBOT] [DEBUG] "+format, args...)
}

func defaultHandler(_ context.Context, _ *Bot, update *models.Update) {
	log.Printf("[TGBOT] [UPDATE] %+v", update)
}

func (b *Bot) error(format string, args ...any) {
	b.errorsHandler(fmt.Errorf(format, args...))
}

// True and False returns the pointer to bool
func True() *bool {
	b := true
	return &b
}

// False and True returns the pointer to bool
func False() *bool {
	b := false
	return &b
}

// FileDownloadLink returns the file download link
func (b *Bot) FileDownloadLink(f *models.File) string {
	return fmt.Sprintf("%s/file/bot%s/%s", b.url, b.token, f.FilePath)
}
