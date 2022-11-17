package svg

import (
	"fmt"
)

func UrlFor(ident string) string {
	return fmt.Sprintf("url(#%s", ident)
}