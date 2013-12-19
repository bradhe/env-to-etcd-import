package main

import (
	"github.com/coreos/go-etcd/etcd"

	"io"
	"strings"
	"bufio"
	"log"
	"os"
	"flag"
)

func main() {
	var node string
	var prefix string
	var envFile string

	flag.StringVar(&node, "node", "http://127.0.0.1:4001", "the node to connect to")
	flag.StringVar(&prefix, "key-prefix", "/", "the prefix for your key")
	flag.StringVar(&envFile, "env-file", "", "the .env file to import")
	flag.Parse()

	if envFile == "" {
		log.Fatal("Please specify a .env file with -env-file")
	}

	fd, err := os.Open(envFile)

	if err != nil {
		log.Fatalf("Failed to open %s for read. %s", envFile, err)
	}

	defer fd.Close()

	reader := bufio.NewReader(fd)

	// client for connecting to etcd
	client := etcd.NewClient(strings.Split(node, ","))

	// let's make sure the prefix is there.
	_, err = client.CreateDir(prefix, 0)

	if err != nil {
		switch err.(etcd.EtcdError).ErrorCode {
		case 105: // Do nothing.
		default: log.Panic(err)
		}
	}

	for {
		line, err := reader.ReadString('\n')

		// Did we get to the end of the line?
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panic(err)
		}

		// Skip lines that are invalid in some way.
		if !strings.Contains(line, ":") {
			continue
		}

		// Other
		toks := strings.SplitN(line, ":", 2)
		key, val := toks[0], toks[1]

		// Clean 'em both up a bit.
		key = strings.TrimSpace(strings.ToLower(key))
		val = strings.TrimSpace(val)

		// Now that we've got all that, let's write it to our dear friend etcd.
		key = strings.Join([]string { prefix, key }, "/")

		_, err = client.Create(key, val, 0)

		// how'd we do?
		if err == nil {
			log.Printf("Created key %s", key)
			continue
		} else if err.(etcd.EtcdError).ErrorCode == 105 {
			_, err = client.Update(key, val, 0)

			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Panic(err)
		}
	}
}
