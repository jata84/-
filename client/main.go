package client

func main() {
	client := NewClient()
	client.Run()
	defer client.Close()
}
