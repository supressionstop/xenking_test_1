package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/supressionstop/xenking_test_1/internal/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	host := flag.String("host", "localhost", "gRPC server host")
	port := flag.String("port", "48002", "gRPC server port")

	conn, err := grpc.NewClient(
		*host+":"+*port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("could not create gRPC client: %w", err))
	}
	defer conn.Close()

	client := pb.NewLinesClient(conn)

	stream, err := client.SubscribeOnSportsLines(context.Background())
	if err != nil {
		log.Fatal(fmt.Errorf("could not subscribe to gRPC client: %w", err))
	}
	defer func(stream grpc.BidiStreamingClient[pb.Subscribe, pb.LinesData]) {
		err := stream.CloseSend()
		if err != nil {
			log.Fatal(fmt.Errorf("could not close stream: %w", err))
		}
	}(stream)

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Response: %s", in)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
	for {
		fmt.Print("$ ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		req, err := parseInput(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch input {
		case "quit":
			return
		default:
			err := stream.Send(req)
			if err != nil {
				log.Println(fmt.Errorf("could not send request: %w", err))
				continue
			}
		}
	}
}

func parseInput(s string) (*pb.Subscribe, error) {
	split := strings.Split(s, " ")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid string: %s, want: \"sport1,sport2,... 1\"", s)
	}

	sports := strings.Split(split[0], ",")
	interval, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("invalid interval: %s, must be integer", split[1])
	}

	return &pb.Subscribe{
		Sports:   sports,
		Interval: int32(interval),
	}, nil
}
