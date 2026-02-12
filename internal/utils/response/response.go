// Package response
package response

import "context"

type JSON = map[string]any

func Error(ctx context.Context, code string, msg string) JSON {
	return JSON{
		"error": JSON{
			"code":    code,
			"message": msg,
		},
	}
}

func Success(ctx context.Context, data any) JSON {
	return JSON{
		"data": data,
	}
}
