package http

// func (r Request) String() (out string) {
// 	out = "Request:\n"
// 	out += "  Param:\n"
// 	for k, v := range r.Param {
// 		out += fmt.Sprintf("    %s: %s\n", k, v)
// 	}

// 	out += "  Header:\n"
// 	for k, v := range r.Header {
// 		out += fmt.Sprintf("    %s: %s\n", k, v)
// 	}

// 	out += "  Cookie:\n"
// 	for k, v := range r.Cookie {
// 		out += fmt.Sprintf("    %s: %s\n", k, v)
// 	}

// 	if r.Body == nil {
// 		out += "  Body: "
// 	} else {
// 		out += fmt.Sprintf("  Body: %v", r.Body)
// 	}

// 	return out
// }

// func (r Response) String() (out string) {
// 	out = "Response:\n"
// 	out += fmt.Sprintf("  Status: %v\n  Header:\n", r.StatusCode)
// 	for k, v := range r.Header {
// 		out += fmt.Sprintf("    %s: %s\n", k, v)
// 	}

// 	out += "  Cookie:\n"
// 	for k, v := range r.Cookie {
// 		out += fmt.Sprintf("    %s: %s\n", k, v)
// 	}

// 	if r.Body == nil {
// 		out += "  Body: "
// 	} else {
// 		out += fmt.Sprintf("  Body: %s", r.Body)
// 	}

// 	return out
// }

// func (h ApiHttp) String() string {
// 	out := "Name: " + h.Name + "\n"
// 	out += "URL: " + h.URL + "\n"
// 	out += "Method: " + h.Method + "\n"
// 	out += "Variables:\n"
// 	for k, v := range h.Variables {
// 		out += fmt.Sprintf("  %s: %s\n", k, v)
// 	}
// 	out += h.Request.String() + "\n"
// 	out += h.Response.String()
// 	return out
// }
