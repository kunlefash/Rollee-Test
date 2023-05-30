package main

import (
    "net/http"
    "net"
    "fmt"
    "log"
    "time"
    "os"
    "regexp"
    "errors"
    "sync/atomic"
)

type httpServer struct {
    server   *http.Server
    mux      *http.ServeMux
    listener  net.Listener // non-nil when server is running

    // These are set by setListenAddr.
	endpoint string
	host     string
	port     int

    storageRequest chan string
    storageResponce chan struct{}

    retrievalRequest chan struct{}
    retrievalResponce chan string

    logger *log.Logger
}


// setListenAddr configure the listening of the server.
// the address can only be set while the server isn't running
func (h *httpServer) setListenAddr(host string, port int) error {
    if h.listener != nil && (host != h.host || port != h.port) {
		return fmt.Errorf("HTTP server already running on %s", h.endpoint)
	}

    h.host, h.port = host, port
   
    h.endpoint = fmt.Sprintf("%s:%d", host, port)
    
    return nil
}

type httpService struct {
    http *httpServer

    httpStore storageService
}

type storageService struct {
    count int32

    prevWord string

    cache chan string

    store map[string]int32
}

func (s *httpService) startServer() {
    s.http.logger =  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

    host := ""
    port := 8545

    if err := s.http.setListenAddr(host, port); err != nil {
        return
    }
     
    if s.http.endpoint == "" || s.http.listener != nil { return }

    s.http.mux = http.NewServeMux()
    s.http.mux.HandleFunc("/service/word", s.http.handleWord) 
    s.http.mux.HandleFunc("/service/prefix", s.http.handlePrefix) 

    s.http.server = &http.Server { Handler: s.http.mux }

    // timeout config
    s.http.server.ReadTimeout = 30 * time.Second
    s.http.server.ReadHeaderTimeout = 30 * time.Second
    s.http.server.WriteTimeout = 30 * time.Second
    s.http.server.IdleTimeout = 30 * time.Second

    s.http.storageRequest  = make(chan string)
    s.http.storageResponce = make(chan struct{})

    s.http.retrievalRequest  = make(chan struct{})
    s.http.retrievalResponce = make(chan string) 

    s.httpStore.store = make(map[string]int32)
    // cache 
    s.httpStore.cache = make(chan string, 1)

    // start the server
    listener, err := net.Listen("tcp", s.http.endpoint)
    if err != nil { return }


    // storage service started
    go s.storageTask()



    s.http.listener = listener

    s.http.logger.Println("Http Service started")

    s.http.server.Serve(listener)
}



func (s *httpService) storageTask() {
    for {
        select {
            // storage
            case newWord := <- s.http.storageRequest:
                s.httpStore.store[newWord] = atomic.AddInt32(&s.httpStore.count, s.httpStore.count + 1)
                s.httpStore.prevWord = newWord
                s.http.storageResponce <- struct{}{}
            // retrieval
            case <-s.http.retrievalRequest:
                s.http.retrievalResponce  <-s.httpStore.prevWord 

        }
    }
}


func (h *httpServer) handleWord(w http.ResponseWriter, r *http.Request) {
    if ! (r.Method == http.MethodPost)  {
        http.Error(w, errors.New("method not allowed").Error(), http.StatusMethodNotAllowed)
        return 
	}
   
    word := r.FormValue("word")
    isValid, err := regexp.MatchString(`[a-zZ-Z]+`, word)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotAcceptable )
        return 
    }

    if !isValid {
        http.Error(w, errors.New("Not the correct word format [a-zZ-Z]+").Error(), http.StatusNotAcceptable )
        return 
    }

    
    select {
        case h.storageRequest <-word: 
            h.logger.Printf("The word %s has been sent to store \r\n", word)
    }

    <- h.storageResponce

    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)

    fmt.Fprint(w, "Stored \r\n")
}












func (h *httpServer)  handlePrefix(w http.ResponseWriter, r *http.Request) {
    if !( r.Method == http.MethodGet ) {
        http.Error(w, errors.New("method not allowed").Error(), http.StatusMethodNotAllowed)
        return 
	}
   // validate the word
    word := r.FormValue("prefix")
    isValid, err := regexp.MatchString(`[a-zZ-Z]+`, word)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotAcceptable )
        return 
    }

    if !isValid {
        http.Error(w, errors.New("Not the correct word format [a-zZ-Z]+").Error(), http.StatusNotAcceptable )
        return 
    }

    // algorithm to validate prefix
    // we assume the algorithm is true
    // var validPrefix = true


    select {
        case h.retrievalRequest <-struct{}{}: 
            h.logger.Printf("Processing last word stored from storage service \r\n")
            fmt.Fprint(w, "Processing last word stored from storage service \r\n")
    }

    

    w.Header().Set("Content-Type", "text/plain")
    

    fmt.Fprintf(w, "%s \r\n", <-h.retrievalResponce)
}




func main() {
    server := httpServer{}
    h := httpService{ http: &server }
    h.startServer()
}



    


 