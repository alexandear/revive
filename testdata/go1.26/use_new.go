package pkg

// Reported: the variable exists only to have its address taken.

func assignToField(size *float32, values []float64, i int) {
	sizeVal := float64(*size) // MATCH /replace the temporary variable sizeVal with new(float64(*size))/
	snapshot.Size = &sizeVal

	nameVal := "a name" // MATCH /replace the temporary variable nameVal with new("a name")/
	snapshot.Name = &nameVal

	timeout := config.Timeout // MATCH /replace the temporary variable timeout with new(config.Timeout)/
	snapshot.Size = &timeout

	item := values[i] // MATCH /replace the temporary variable item with new(values[i])/
	snapshot.Size = &item
}

func definePointer() *int {
	count := 42 // MATCH /replace the temporary variable count with new(42)/
	countPtr := &count

	return countPtr
}

func returnAddress(b []byte) *string {
	s := string(b) // MATCH /replace the temporary variable s with new(string(b))/

	return &s
}

func passAddress(flag bool) {
	enabled := !flag // MATCH /replace the temporary variable enabled with new(!flag)/
	consume(&enabled)
}

func fillCompositeLiteral(b []byte) *snapshotType {
	name := string(b) // MATCH /replace the temporary variable name with new(string(b))/

	return &snapshotType{Name: &name}
}

func inClauses(ch chan int, i int) {
	switch i {
	case 0:
		inCase := int64(i) // MATCH /replace the temporary variable inCase with new(int64(i))/
		int64s = append(int64s, &inCase)
	}

	select {
	case v := <-ch:
		inComm := int64(v) // MATCH /replace the temporary variable inComm with new(int64(v))/
		int64s = append(int64s, &inComm)
	}
}

// Not reported: the variable has another use, or is not declared right before its address is taken.

func otherUses(flag bool, i int) {
	usedTwice := float64(i)
	snapshot.Size = &usedTwice
	_ = usedTwice

	var declared = float64(i)
	snapshot.Size = &declared

	multiValue, _ := lookup(i)
	snapshot.Size = &multiValue

	notAdjacent := float64(i)
	_ = flag
	snapshot.Size = &notAdjacent

	addressTakenTwice := float64(i)
	snapshot.Size = &addressTakenTwice
	snapshot.Size = &addressTakenTwice

	inFuncLit := !flag
	run(func() { consume(&inFuncLit) })
}

// Not reported: inlining the value would make the code harder to read.

func complexValues(name string) {
	compound := count(name) + 1
	snapshot.Size = &compound

	callResult := lookupSize(name)
	snapshot.Size = &callResult

	dereferenced := *newBuilder(name).Obj()
	consumeSnapshot(&dereferenced)

	literal := snapshotType{Name: snapshot.Name}
	consumeSnapshot(&literal)
}
