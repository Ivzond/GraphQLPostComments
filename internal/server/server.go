package server

func NewServer() *handler.Server {
	server := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}}))
	return server
}
