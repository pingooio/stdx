package httpx

import (
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

func ServerInternalErrorPlaintext(res http.ResponseWriter, message *string) {
	var data string

	if message == nil {
		data = "Internal Error\n"
	} else {
		data = *message
	}

	ServeError(res, data, http.StatusInternalServerError)
}

// TODO: Do we force close the TCP connection?
// https://stackoverflow.com/questions/72368886/is-there-a-way-to-drop-an-http-connection-in-golang-without-sending-anything-to
// https://www.bentasker.co.uk/posts/blog/software-development/golang-net-http-net-http-2-does-not-reliably-close-failed-connections-allowing-attempted-reuse.html
func ServerForbiddenError(res http.ResponseWriter) {
	sleepForMs := rand.Int64N(1000) + 1500
	time.Sleep(time.Duration(sleepForMs) * time.Millisecond)
	res.Header().Set(HeaderConnection, "close")
	ServeError(res, "Forbidden\n", http.StatusForbidden)
}

func ServeError(res http.ResponseWriter, message string, statusCode int) {
	res.Header().Del(HeaderETag)
	res.Header().Set(HeaderCacheControl, CacheControlNoCache)
	// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set(HeaderContentType, MediaTypeHtmlUtf8)
	res.Header().Set(HeaderContentLength, strconv.FormatInt(int64(len(message)), 10))
	res.WriteHeader(statusCode)
	res.Write([]byte(message))
}
