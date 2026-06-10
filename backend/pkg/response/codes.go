package response

const (
	CodeInvalidParams       = "INVALID_PARAMS"
	CodeInvalidRecordID     = "INVALID_RECORD_ID"
	CodeUsernameExists      = "USERNAME_EXISTS"
	CodeInvalidCredentials  = "INVALID_CREDENTIALS"
	CodeTokenRequired       = "AUTH_TOKEN_REQUIRED"
	CodeTokenMalformed      = "AUTH_TOKEN_MALFORMED"
	CodeTokenInvalid        = "AUTH_TOKEN_INVALID"
	CodeSessionExpired      = "AUTH_SESSION_EXPIRED"
	CodeRateLimited         = "RATE_LIMITED"
	CodeUserNotFound        = "USER_NOT_FOUND"
	CodeHistoryNotFound     = "HISTORY_NOT_FOUND"
	CodeHistoryListFailed   = "HISTORY_LIST_FAILED"
	CodeHistorySaveFailed   = "HISTORY_SAVE_FAILED"
	CodeHistoryDeleteFailed = "HISTORY_DELETE_FAILED"
	CodeConvertFailed       = "MARKDOWN_CONVERT_FAILED"
	CodeInternalError       = "INTERNAL_ERROR"
)
