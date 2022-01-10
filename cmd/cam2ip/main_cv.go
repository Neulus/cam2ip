// +build cv2 cv4

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jamiealquiza/envy"

	"github.com/Neulus/cam2ip/camera"
	"github.com/Neulus/cam2ip/server"
	"github.com/Neulus/cam2ip/video"
)

const (
	name    = "cam2ip"
	version = "1.6"
)

func main() {
	srv := server.NewServer()

	flag.IntVar(&srv.Index, "index", 0, "Camera index")
	flag.IntVar(&srv.Delay, "delay", 10, "Delay between frames, in milliseconds")
	flag.Float64Var(&srv.FrameWidth, "width", 640, "Frame width")
	flag.Float64Var(&srv.FrameHeight, "height", 480, "Frame height")
	flag.IntVar(&srv.Rotate, "rotate", 0, "Rotate image, valid values are 90, 180, 270")
	flag.BoolVar(&srv.NoWebGL, "nowebgl", false, "Disable WebGL drawing of images (html handler)")
	flag.BoolVar(&srv.Timestamp, "timestamp", false, "Draws timestamp on images")
	flag.StringVar(&srv.Bind, "bind-addr", ":56000", "Bind address")
	flag.StringVar(&srv.Htpasswd, "htpasswd-file", "", "Path to htpasswd file, if empty auth is disabled")
	flag.StringVar(&srv.FileName, "video-file", "", "Use video file instead of camera")

	envy.Parse("CAM2IP")
	flag.Parse()

	srv.Name = name
	srv.Version = version

	var err error

	if srv.Htpasswd != "" {
		if _, err = os.Stat(srv.Htpasswd); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	}

	if srv.FileName != "" {
		if _, err = os.Stat(srv.FileName); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	}

	if srv.FileName != "" {
		vid, err := video.New(video.Options{srv.FileName, srv.Rotate})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		srv.Reader = vid
	} else {
		cam, err := camera.New(camera.Options{srv.Index, srv.Rotate, srv.FrameWidth, srv.FrameHeight, srv.Timestamp})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		srv.Reader = cam
	}

	defer srv.Reader.Close()

	fmt.Fprintf(os.Stderr, "Listening on %s\n", srv.Bind)

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
