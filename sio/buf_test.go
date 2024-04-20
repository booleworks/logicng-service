package sio

import (
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	s := "((v0 | v1 | v2 | (v16 | v19 | v20 | v21 | v22 | v23) & ~v39 | (v24 | v25 | v26) & ~(v3 | v4)) & ~(v17 | v18) | (v17 | v18) & (v1 | v2 | v16 | v19 | v20 | v21 | v22 | v23 | (v24 | v25 | v26) & ~(v3 | v4))) & ~(v30 | v31 | v32 | v33 | v34 | v35 | v36 | v37 | v38 | v6 | v7 | v8 | v9 | v10 | v11 | v12 | v13 | v14 | v15 | v27 | v28 | v29) => v5"
	input := FormulaInput{s}
	bin, err := input.ProtoBuf()
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.Create("test.buf")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = file.Write(bin)
	if err != nil {
		fmt.Println(err)
	}
}
