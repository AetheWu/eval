package eval

import (
	"fmt"
)

var errIFArgsNum = fmt.Errorf("op IF need 3 args")

func callNotProcedureError(n *node) error {
	return fmt.Errorf("%s is not a procedure;\nexpected a procedure that can be applied", n.tok.String())
}
