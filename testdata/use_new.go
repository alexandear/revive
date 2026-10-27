package fixtures

type useNewSnapshot struct {
	Size *float64
	Name *string
}

type useNewAPISnapshot struct{}

func (useNewAPISnapshot) GetSizeOk() (*float32, bool) { return nil, true }

func useNewConsume(*bool) {}

func useNewRun(func()) {}

func useNewLookup(int) (float64, bool) { return 0, false }

var useNewInt64s []*int64

func useNewAssign(api useNewAPISnapshot, snapshot *useNewSnapshot) {
	if size, ok := api.GetSizeOk(); ok && size != nil {
		sizeVal := float64(*size)
		snapshot.Size = &sizeVal
	}

	nameVal := "a name"
	snapshot.Name = &nameVal
}

func useNewDefine() *int {
	count := 42
	countPtr := &count

	return countPtr
}

func useNewReturn(b []byte) *string {
	s := string(b)

	return &s
}

func useNewCallArgument(flag bool) {
	enabled := !flag
	useNewConsume(&enabled)
}

func useNewCompositeLiteral(b []byte) *useNewSnapshot {
	name := string(b)

	return &useNewSnapshot{Name: &name}
}

func useNewField(config useNewConfigType, snapshot *useNewSnapshot) {
	timeout := config.Timeout
	snapshot.Size = &timeout
}

func useNewClauses(ch chan int, i int) {
	switch i {
	case 0:
		inCase := int64(i)
		useNewInt64s = append(useNewInt64s, &inCase)
	}

	select {
	case v := <-ch:
		inComm := int64(v)
		useNewInt64s = append(useNewInt64s, &inComm)
	}
}

func useNewValid(flag bool, i int, snapshot *useNewSnapshot) {
	usedTwice := float64(i)
	snapshot.Size = &usedTwice
	_ = usedTwice

	var declared = float64(i)
	snapshot.Size = &declared

	multiValue, _ := useNewLookup(i)
	snapshot.Size = &multiValue

	notAdjacent := float64(i)
	_ = flag
	snapshot.Size = &notAdjacent

	inFuncLit := !flag
	useNewRun(func() { useNewConsume(&inFuncLit) })

	addressTakenTwice := float64(i)
	snapshot.Size = &addressTakenTwice
	snapshot.Size = &addressTakenTwice

	literal := useNewSnapshot{Name: snapshot.Name}
	useNewConsumeSnapshot(&literal)

	computed := func() float64 { return 1 }()
	snapshot.Size = &computed
}

// The value of a reported variable has to be simple enough to read inside the call to new,
// and the resulting statement has to stay short.
func useNewComplexValue(state *useNewRequeueState, snapshot *useNewSnapshot, name string) {
	compound := useNewDeref(state.Count, 0) + 1
	state.Count = &compound

	notAConversion := useNewAmountQuantity(name, 1)
	state.Quantity = &notAConversion

	chainedCalls := useNewJobType{}.GetObjectKind().GroupVersionKind()
	useNewDiff(useNewWantGVK, &chainedCalls)

	dereferencedBuilder := *useNewMakePodSet(name, 1).Obj()
	useNewRequests = []useNewPodSetRequest{{PodSet: &dereferencedBuilder}}

	longValue := float64(useNewConfig.Server.Timeouts.ReadHeaderTimeoutSeconds)
	snapshot.Size = &longValue

	lostInStatement := int32(2)
	useNewReportPending(name, useNewQuantityToFloat(&lostInStatement), useNewCustomLabels, useNewRoleTracker)
}

type useNewConfigType struct {
	Timeout float64
}

type useNewQuantity struct{}

type useNewRequeueState struct {
	Count    *int32
	Quantity *useNewQuantity
}

type useNewGVK struct{}

type useNewObjectKind struct{}

func (useNewObjectKind) GroupVersionKind() useNewGVK { return useNewGVK{} }

type useNewJobType struct{}

func (useNewJobType) GetObjectKind() useNewObjectKind { return useNewObjectKind{} }

type useNewPodSet struct{}

type useNewPodSetBuilder struct{}

func (*useNewPodSetBuilder) Obj() *useNewPodSet { return nil }

type useNewPodSetRequest struct {
	PodSet *useNewPodSet
}

func useNewConsumeSnapshot(*useNewSnapshot) {}

func useNewDeref(*int32, int32) int32 { return 0 }

func useNewAmountQuantity(string, int) useNewQuantity { return useNewQuantity{} }

func useNewDiff(useNewGVK, *useNewGVK) string { return "" }

func useNewMakePodSet(string, int) *useNewPodSetBuilder { return nil }

func useNewQuantityToFloat(*int32) float64 { return 0 }

func useNewReportPending(string, float64, []string, string) {}

var (
	useNewWantGVK      useNewGVK
	useNewRequests     []useNewPodSetRequest
	useNewCustomLabels []string
	useNewRoleTracker  string
	useNewConfig       struct {
		Server struct {
			Timeouts struct {
				ReadHeaderTimeoutSeconds int
			}
		}
	}
)
