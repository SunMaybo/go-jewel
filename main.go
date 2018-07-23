package main

import (
	"os"
	"strings"
	"github.com/alecthomas/kingpin"
	"github.com/gin-gonic/gin"
)

var (
	app          = kingpin.UsageTemplate(kingpin.DefaultUsageTemplate)
	debug        = app.Flag("debug", "Enable debug mode.").Bool()
	serverIP     = app.Flag("server", "Server address.").Default("127.0.0.1").IP()
	register     = app.Command("register", "Register a new user.")
	registerNick = register.Arg("nick", "Nickname for user.").Required().String()
	registerName = register.Arg("name", "Name of user.").Required().String()
	post         = app.Command("post", "Post a message to a channel.")
	postImage    = post.Flag("image", "Image to post.").File()
	postChannel  = post.Arg("channel", "Channel to post to.").Required().String()
	postText     = post.Arg("text", "Text to post.").Strings()
	//get          = post.Command("get", "get a message channel.")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Register user
	case register.FullCommand():
		println(*registerNick)

		// Post message
	case post.FullCommand():
		if *postImage != nil {
		}
		text := strings.Join(*postText, " ")
		println("Post:", text)
	}
}
