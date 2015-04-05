package translated

import "github.com/stephens2424/php/passes/togo/internal/phpctx"

func Shell(ctx phpctx.PHPContext) {
	ctx.Shell(`ls -al`)
}
