package translated

import "github.com/stephens2424/php/passes/togo/internal/phpctx"

func Echo(ctx phpctx.PHPContext) {
	ctx.Echo.Write("test")
}
