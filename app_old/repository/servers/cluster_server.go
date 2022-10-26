package servers

// Define the Cluster Rest server:
// - issue logs for FSM
// - cluster membership changes
// - authentication
// - cluster admin

type Cluster struct {
	//router *http.ServeMux
}

// type repo struct {
// 	r *Runner
// 	*clients.Runner
// }

// // take repo and service as argument
// func NewServer(runner *Runner, client *clients.Runner) *Cluster {
// 	router := http.NewServeMux()
// 	// create repo
// 	// inject repo in service
// 	// inject service in makeHandler
// 	// repo := &repo{r: runner, client}
// 	// service := cluster.NewService(repo)
// 	return &Cluster{router: router}
// }

// func (s *Cluster) makeHandler() {}

// func (s *Cluster) Start() {
// 	log.Println("START CLUSTER SERVER")
// 	log.Fatal(http.ListenAndServe(":8080", s.router))
// }
