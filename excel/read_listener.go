package excel

type ReadListener interface {
	// interface is option's dest without ptr
	ReadCompleteTrigger(rCtx ReadSheetContext, dest interface{}) error
	// interface is option's dest elem with ptr
	ReadCellCompleteTrigger(rCtx ReadCellContext, destElem interface{}, err error) error
}
