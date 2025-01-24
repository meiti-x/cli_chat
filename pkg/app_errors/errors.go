package app_errors

import "errors"

// WebSocket errors
var (
	ErrSocketUpgradeFailed = errors.New("failed to upgrade to WebSocket connection")
	ErrSocketWriteFailed   = errors.New("failed to write WebSocket message")
	ErrSocketReadFailed    = errors.New("failed to read WebSocket message")
)

// NATS errors
var (
	ErrNatsInit               = errors.New("failed to initialize Nats connection")
	ErrNATSConnectionFailed   = errors.New("failed to connect to NATS server")
	ErrNATSSubscriptionFailed = errors.New("failed to subscribe to NATS channel")
	ErrNATSReceivedFailed     = errors.New("error receiving NATS message")
)

// Common errors
var (
	InvalidCommand        = errors.New("invalid command received")
	ErrSendWelcomeMessage = errors.New("error send welcome message")
	ErrSendJoinMessage    = errors.New("error send join message")
	ErrSendLeaveMessage   = errors.New("error send leave message")
	ErrSendOnlineUsers    = errors.New("error send leave message")
	ErrHttpStart          = errors.New("HTTP server error")
	ErrInitDB             = errors.New("error init db")
	ErrParseJSON          = errors.New("error parsing JSON message")
)
