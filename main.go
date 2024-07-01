package main

import (
    "fmt"
    "net"
)

type Server struct {
    listenAddr string
    ln         net.Listener
    quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
    return &Server{
        listenAddr: listenAddr,
        quitch:     make(chan struct{}),
    }
}

func (s *Server) Start() error {
    ln, err := net.Listen("tcp", s.listenAddr)
    if err != nil {
        return err
    }
    s.ln = ln
    defer s.ln.Close()

    go s.acceptLoop()

    <-s.quitch
    return nil
}

func (s *Server) acceptLoop() {
    connCount := 0
    for {
        conn, err := s.ln.Accept()
        if err != nil {
            fmt.Println("accept error", err)
            continue
        }
        connCount++
        go s.readLoop(conn, fmt.Sprintf("conn%d", connCount))
    }
}

func (s *Server) readLoop(conn net.Conn, connName string) {
    defer conn.Close()
    buf := make([]byte, 2048)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println("read error", err)
            continue
        }

        msg := buf[:n]
        fmt.Printf("%s: %s\n", connName, string(msg))
    }
}

func main() {
    server := NewServer(":3000")
    err := server.Start()
    if err != nil {
        fmt.Println("Server start error:", err)
    }
}