package interceptor

import (
	"context"

	"github.com/bobacgo/kit/app/validator"
	"google.golang.org/grpc/metadata"
)

func GetLanguage(ctx context.Context) (lang string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(validator.GetLanguageCtxKey()); len(values) > 0 {
			lang = values[0]
		}
	}
	return
}
