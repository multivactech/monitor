package errorInfo

type ErrorInfo struct {
	Host         string
	KeyWord      string
	ErrorContent string
	Node         string
}

var AllError chan *ErrorInfo

//TODO: sortite errorInfo to more type, and add whitelist...........some error we can ignore
func FindErrorFromString(content string) {
	// 1. 按行切分， 以 [ERR] 所在行为界，进行上下判断
	
	// 2. 对每一段进行选取error

	// 3. 判断是否为白名单的错误, 向AllError发送错误。
		
}
