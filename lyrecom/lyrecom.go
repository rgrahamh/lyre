package lyrecom

var PAYLOAD_MAX = 65535

func Memset(buffer []byte, c byte, n int) {
	for i := range n {
		buffer[i] = c
	}
}
