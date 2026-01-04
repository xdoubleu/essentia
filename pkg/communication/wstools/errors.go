package wstools

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/xdoubleu/essentia/pkg/contexttools"
	"github.com/xdoubleu/essentia/pkg/errortools"
	"github.com/xdoubleu/essentia/pkg/logging"
)

// ErrorResponse is used to handle any kind of error that occurred on a WebSocket.
func ErrorResponse(
	ctx context.Context,
	conn *websocket.Conn,
	status int,
	message any,
) {
	errorDto := errortools.NewErrorDto(status, message)
	err := wsjson.Write(ctx, conn, errorDto)
	if err != nil {
		contexttools.Logger(ctx).
			ErrorContext(ctx, "failed to write JSON", logging.ErrAttr(err))
	}
}

// ServerErrorResponse is used to handle
// internal server errors that occurred on a WebSocket.
func ServerErrorResponse(ctx context.Context, conn *websocket.Conn, err error) {
	if isCloseError(err) {
		return
	}

	contexttools.Logger(ctx).
		ErrorContext(ctx, "server error occurred", logging.ErrAttr(err))

	message := errortools.MessageInternalServerError
	if contexttools.ShowErrors(ctx) {
		message = err.Error()
	}

	ErrorResponse(ctx, conn, http.StatusInternalServerError, message)
}

// UpgradeErrorResponse is used to handle an error that
// occurred during the upgrade towards a WebSocket.
func UpgradeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	if !isWSProtocolViolation(err) {
		contexttools.Logger(r.Context()).
			ErrorContext(r.Context(), "WS upgrade error occurred", logging.ErrAttr(err))
	}

	w.WriteHeader(http.StatusInternalServerError)
}

// FailedValidationResponse is used to handle
// an error of a [validate.Validator] that occurred on a WebSocket.
func FailedValidationResponse(
	ctx context.Context,
	conn *websocket.Conn,
	errors map[string]string,
) {
	ErrorResponse(ctx, conn, http.StatusUnprocessableEntity, errors)
}

// ForbiddenResponse is used to handle an error when a user
// isn't authorized to access a certain resource.
func ForbiddenResponse(ctx context.Context, conn *websocket.Conn) {
	ErrorResponse(ctx, conn, http.StatusForbidden, errortools.MessageForbidden)
}

func isWSProtocolViolation(err error) bool {
	return strings.Contains(err.Error(), "WebSocket protocol violation")
}

func isCloseError(err error) bool {
	// EOF, close errors
	var closeError websocket.CloseError
	return errors.Is(err, io.EOF) ||
		errors.As(err, &closeError)
}
