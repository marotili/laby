package server

func handleConnection(conn net.Conn) {
	player := &Player{Conn: conn}
	defer func() {
		log.Printf("Player %q left (%s)\n", player.Nick, conn.RemoteAddr())
		mu.Lock()
		game.Leave(player)
		mu.Unlock()
	}()

	var nickLength uint32
	binary.Read(conn, binary.BigEndian, &nickLength)
	nickBytes := make([]byte, nickLength)
	if _, err := io.ReadFull(conn, nickBytes); err != nil {
		log.Println(err)
		return
	}
	player.Nick = string(nickBytes)

	mu.Lock()
	game.Join(player, 200, 200)
	mu.Unlock()

	log.Printf("Player %q joined (%s)\n", player.Nick, conn.RemoteAddr())

	buf := make([]byte, 12)
	for {
		// handle player input and send updates
	}
}

var (
	game Game
	mu   sync.Mutex
)

func main() {
	var err error

	if game, err = NewGame(); err != nil {
		log.Fatal(err)
	}

	// go func() {
	listen, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
	// }()
}
