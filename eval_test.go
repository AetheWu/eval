package eval

import "testing"

var filename = "D:/code/Golang/balance.rkt"
var input = FileLoad(filename)

func BenchmarkEval(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RunSourceFile(filename)
	}
}
