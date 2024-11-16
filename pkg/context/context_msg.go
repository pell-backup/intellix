package pkgcontext

const (
	ctxDvsRequestKey = "CTX_DVS_POST_RESPONSE"
)

func (c Context) WithDvsPostResponseData(postProcessResponseData []byte) Context {
	return c.WithValue(ctxDvsRequestKey, postProcessResponseData)
}

func (c Context) DvsPostResponseData() ([]byte, bool) {
	value := c.Value(ctxDvsRequestKey)
	val, ok := value.([]byte)
	return val, ok
}
