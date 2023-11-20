package cache

type Result struct {
	data interface{}
	err  error
}

func (r *Result) Data() interface{} {
	return r.data
}

func (r *Result) Error() error {
	return r.err
}

func (r *Result) Int32() int32 {
	if _, ok := r.data.(int32); !ok {
		return 0
	}
	return r.data.(int32)
}

func (r *Result) Int64() int64 {
	if _, ok := r.data.(int64); !ok {
		return 0
	}
	return r.data.(int64)
}

func (r *Result) String() string {
	if _, ok := r.data.(string); !ok {
		return ""
	}
	return r.data.(string)
}
