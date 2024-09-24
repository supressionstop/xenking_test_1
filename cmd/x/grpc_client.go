package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/supressionstop/xenking_test_1/internal/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewLinesClient(conn)

	stream, err := client.SubscribeOnSportsLines(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer stream.CloseSend()

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Got message: %s", in)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		switch text {
		case "q":
			sub := &pb.Subscribe{
				Sports:   []string{"baseball"},
				Interval: 2,
			}
			_ = stream.Send(sub)
		case "w":
			sub := &pb.Subscribe{
				Sports:   []string{"baseball", "football", "soccer"},
				Interval: 3,
			}
			_ = stream.Send(sub)
		case "quit":
			return
		}
	}
}
