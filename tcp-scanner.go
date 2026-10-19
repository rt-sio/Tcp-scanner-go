package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func worker(ports chan int, result chan int, target string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", target, p)
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err != nil {
			result <- 0
			continue
		}
		conn.Close()
		result <- p
	}
}

func main() {
	fmt.Println("Введите целевой хост или IP-адрес (например, google.com или 192.168.1.1):")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	target := strings.TrimSpace(scanner.Text())
	start := time.Now()
	ports := make(chan int, 100)
	result := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, result, target)
	}
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
		close(ports)
	}()

	for i := 0; i < 1024; i++ {
		ports := <-result
		if ports != 0 {
			openports = append(openports, ports)
		}
	}

	close(result)
	for _, ports := range openports {
		fmt.Printf("%d open\n", ports)
	}
	duration := time.Since(start)
	fmt.Println(duration)
}
