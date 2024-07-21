package main

import (
	"fmt"
	"log"
	"mario/emoji-cdn/cli"
	"mario/emoji-cdn/utils"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "go.uber.org/automaxprocs"
)

func main() {
	utils.LoadConfig(true)
	utils.InitHttpClient()

	args := os.Args[1:]

	if slices.Contains(args, "--update-db") {
		err := utils.InitEmojipediaNextjsBuildID()
		if err != nil {
			log.Fatal(err)
		}

		dbUpdaterThreads := 0

		for _, arg := range args {
			if strings.HasPrefix(arg, "--threads=") {
				split := strings.Split(arg, "--threads=")

				if len(split) == 2 {
					dbUpdaterThreads = utils.StrToI(split[1])
				}

				break
			}
		}

		if dbUpdaterThreads <= 0 {
			dbUpdaterThreads = 10
		}

		cli.UpdateDb(dbUpdaterThreads)

		return
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	if utils.Config.Get("enable_pprof").Bool() {
		pprof.Register(r)
	}
	if utils.Config.Get("corsAllowAll").Bool() {
		r.Use(cors.Default())
	}

	r.UseRawPath = true
	r.UnescapePathValues = false

	initErrors(r)
	initRoutes(r)

	// disable timeouts to prevent interruptions during large file uploads & downloads.
	// ... see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/ (the "About streaming" part)
	// ... and see https://ieftimov.com/posts/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/
	httpServer := &http.Server{
		Handler:           r,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		ReadHeaderTimeout: 30 * time.Second, // https://ieftimov.com/posts/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/#server-timeouts---first-principles
	}
	httpServer.SetKeepAlivesEnabled(true)

	LISTENER_TYPE := utils.Config.Get("listener.type").String()
	LISTENER_DATA := utils.Config.Get("listener.data").String()
	if LISTENER_TYPE != "tcp" && LISTENER_TYPE != "unix_socket" {
		log.Fatal("Invalid listener type. Please use either 'tcp' or 'unix_socket'")
	}
	fmt.Println("Listening at " + LISTENER_DATA + "...")
	if LISTENER_TYPE == "tcp" {
		httpServer.Addr = LISTENER_DATA
		if utils.Config.Get("ssl.enabled").Bool() {
			SSL_KEY_PATH := utils.Config.Get("ssl.key_path").String()
			SSL_CERT_PATH := utils.Config.Get("ssl.cert_path").String()
			err := httpServer.ListenAndServeTLS(SSL_CERT_PATH, SSL_KEY_PATH)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			err := httpServer.ListenAndServe()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		listener, err := net.Listen("unix", LISTENER_DATA)
		if err != nil {
			fmt.Println("failed to listen at the specified unix socket")
			return
		}
		httpServer.Serve(listener)
		listener.Close()
		os.Remove(LISTENER_DATA)
	}
}
