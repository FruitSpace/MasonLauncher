package main

//var proxy = goproxy.NewProxyHttpServer()

//func (*App) Fetch(u string, o map[string]interface{}) map[string]interface{} {
//	var body io.Reader
//	method := http.MethodGet
//	header := make(http.Header)
//
//	if headers, ex := o["headers"]; ex {
//		if hmap, okh := headers.(map[string]interface{}); okh {
//			for key, value := range hmap {
//				v := value.([]interface{})
//				for _, item := range v {
//					var i []string
//					i = append(i, item.(string))
//					header[textproto.CanonicalMIMEHeaderKey(key)] = i
//				}
//			}
//		}
//	}
//
//	if b, ex := o["body"]; ex {
//		if bo, ok := b.(string); ok {
//			body = bytes.NewBufferString(bo)
//		}
//	}
//
//	if m, ex := o["method"]; ex {
//		if me, ok := m.(string); ok {
//			method = me
//		}
//	}
//
//	var toRet map[string]interface{}
//
//	res := httptest.NewRecorder()
//	req, err := http.NewRequest(method, u, body)
//	if err != nil {
//		toRet = map[string]interface{}{
//			"body":    fmt.Sprintf("Internal Server Error: %s", err.Error()),
//			"headers": make(map[string][]string),
//			"status":  http.StatusInternalServerError,
//			"method":  method,
//			"url":     u,
//		}
//	} else {
//		req.Header = header
//		proxy.ServeHTTP(res, req)
//		toRet = map[string]interface{}{
//			"body":    res.Body.String(),
//			"headers": map[string][]string(res.Header()),
//			"status":  res.Code,
//			"method":  method,
//			"url":     u,
//		}
//	}
//
//	return toRet
//}
