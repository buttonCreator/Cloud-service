package logger

import "context"

type constKey int

const ctxFieldsKey constKey = 1

func contextFields(ctx context.Context) []any {
	fieldsValueVal := ctx.Value(ctxFieldsKey)

	if fieldsValueVal == nil {
		return nil
	}

	fieldsValueSlice, ok := fieldsValueVal.([]any)
	if !ok {
		return nil
	}

	return fieldsValueSlice
}
