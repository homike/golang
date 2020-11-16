package inject

import (
	"fmt"

	"github.com/facebookarchive/inject"
)

type Conf struct {
	A int
}

type DB struct {
	D string
}

type UserController struct {
	UserService *UserService `inject:""`
	Conf        *Conf        `inject:""`
}

type PostController struct {
	UserService *UserService `inject:""`
	PostService *PostService `inject:""`
	Conf        *Conf
}

type UserService struct {
	Db   *DB   `inject:""`
	Conf *Conf `inject:""`
}

type PostService struct {
	Db *DB `inject:""`
}

type Server struct {
	UserApi *UserController `inject:""`
	PostApi *PostController `inject:""`
}

func (self *Server) Run() {
	fmt.Print(self.UserApi.Conf.A, " \n [db] ", self.PostApi.PostService.Db.D)
}

func Inject() {
	conf := &Conf{1}
	db := &DB{"DDD"}

	server := Server{}

	graph := inject.Graph{}

	if err := graph.Provide(
		&inject.Object{
			Value: &server,
		},
		&inject.Object{
			Value: conf,
		},
		&inject.Object{
			Value: db,
		},
	); err != nil {
		panic(err)
	}

	if err := graph.Populate(); err != nil {
		panic(err)
	}

	server.Run()
}

func Normal() {
	conf := &Conf{1}
	db := &DB{"DDD"}

	userService := &UserService{
		Db:   db,
		Conf: conf,
	}

	postService := &PostService{
		Db: db,
	}

	userHandler := &UserController{
		UserService: userService,
		Conf:        conf,
	}

	postHandler := &PostController{
		UserService: userService,
		PostService: postService,
		Conf:        conf,
	}

	server := &Server{
		UserApi: userHandler,
		PostApi: postHandler,
	}

	server.Run()
}
