package requests

type TextRequestsWrapper struct {
	headers           map[string]string
	extra             string
	arbitrary_allowed bool
}

func NewTextRequestsWrapper() TextRequestsWrapper {
	return TextRequestsWrapper{
		extra:             "forbid",
		arbitrary_allowed: true,
	}
}

func (wrapper *TextRequestsWrapper) requests() *Requests {
	return NewRequests(wrapper.headers)
}

func (wrapper *TextRequestsWrapper) Get(url string) (string, error) {
	return wrapper.requests().get(url)
}

func (wrapper *TextRequestsWrapper) Post(url string, data map[string]interface{}) (string, error) {
	return wrapper.requests().post(url, data)
}

func (wrapper *TextRequestsWrapper) Patch(url string, data map[string]interface{}) (string, error) {
	return wrapper.requests().patch(url, data)
}

func (wrapper *TextRequestsWrapper) Put(url string, data map[string]interface{}) (string, error) {
	return wrapper.requests().put(url, data)
}

func (wrapper *TextRequestsWrapper) Delete(url string) (string, error) {
	return wrapper.requests().delete(url)
}
