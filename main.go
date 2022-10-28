package shellops

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

// Performs copy operation between streams: os and tcp streams
func stream_copy(src io.Reader, dst io.Writer) <-chan int {
	log.Println("stream_copy")
	buf := make([]byte, 1024)
	sync_channel := make(chan int)
	go func() {
		defer func() {
			if con, ok := dst.(net.Conn); ok {
				con.Close()
				log.Printf("Connection from %v is closed\n", con.RemoteAddr())
			}
			sync_channel <- 0 // Notify that processing is finished
		}()
		for {

			var nBytes int
			var err error
			nBytes, err = src.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Read error: %s\n", err)
				}
				break
			}

			// _, err = dst.Write(buf[0:nBytes])
			// if err != nil {
			// 	log.Fatalf("Write error: %s\n", err)
			// }

			log.Println(buf[0:nBytes])

		}
	}()
	return sync_channel
}

func CmdBashSudo(cmdnya string) {
	if cmdnya != "" {
		if len(strings.TrimSpace(cmdnya)) > 0 {
			cmd := exec.Command("bash", "-c", "sudo "+cmdnya)

			cmd.Stdin = os.Stdin
			//cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			stdout, err := cmd.Output()
			if err != nil {
				//fmt.Println("Err", err)
				fmt.Fprintln(os.Stderr, err)
				return
			}
			fmt.Println(string(stdout))
		}	
	}

}

func AsyncCmdBashSudo(cmdnya string) {
	if cmdnya != "" {
		cmd := exec.Command("bash", "-c", "sudo "+cmdnya)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmdReader, _ := cmd.StdoutPipe()

		scanner := bufio.NewScanner(cmdReader)
		done := make(chan bool)
		go func() {
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			done <- true
		}()
		cmd.Start()
		<-done
		_ = cmd.Wait()
	}

}

func CmdBash(cmdnya string) {
	if cmdnya != "" {
		if len(strings.TrimSpace(cmdnya)) > 0 {
			cmd := exec.Command("bash", "-c", cmdnya)

			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println("Err", err)
				return
			}
			fmt.Println(string(stdout))
		}
	}
}
