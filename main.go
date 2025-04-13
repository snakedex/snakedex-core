package main

import (
	"context"
)

func main() {
	ctx := context.Background()
	docker(ctx)
	test := TestData{"onPrem", []string{"data", "block", "test"}}
	template_engine(test)

}
