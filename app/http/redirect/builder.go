// app/http/redirect/builder.go
package redirect

import (
	"github.com/goravel/framework/contracts/http"
)

type Builder struct {
	ctx      http.Context
	location string
	code     int
}

func New(ctx http.Context) *Builder {
	return &Builder{
		ctx:  ctx,
		code: 302, // Default redirect code
	}
}

func (b *Builder) To(location string) *Builder {
	b.location = location
	return b
}

func (b *Builder) Back() *Builder {
	b.location = b.ctx.Request().Header("Referer", "/")
	return b
}

func (b *Builder) WithStatus(code int) *Builder {
	b.code = code
	return b
}

func (b *Builder) With(key string, value any) *Builder {
	b.ctx.Request().Session().Flash(key, value)
	return b
}

func (b *Builder) WithErrors(errors any) *Builder {
	b.ctx.Request().Session().Flash("errors", errors)
	return b
}

func (b *Builder) WithInput() *Builder {
	input := b.ctx.Request().All()
	delete(input, "password")
	delete(input, "password_confirmation")
	b.ctx.Request().Session().Flash("_old_input", input)
	return b
}

func (b *Builder) Go() http.Response {
	return b.ctx.Response().Redirect(b.code, b.location)
}
