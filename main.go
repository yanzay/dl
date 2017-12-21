package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/yanzay/log"
)

const (
	youtubeCommand = "youtube-dl"
	bestFormat     = "bestvideo[height<=?1080]+bestaudio/best"
)

var (
	addr = flag.String("addr", ":8080", "Address to listen")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", downloadHandler)
	log.Infof("Starting server at %s", *addr)
	http.ListenAndServe(*addr, nil)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "invalid video id", http.StatusBadRequest)
		return
	}

	log.Infof("Getting filename for video %s", id)
	filename, err := getFilename(id)
	if err != nil {
		http.Error(w, "unable to fetch video name", http.StatusUnprocessableEntity)
		return
	}

	log.Infof("Downloading video to %s", filename)
	err = downloadFile(id)
	if err != nil {
		http.Error(w, "unable to download video", http.StatusUnprocessableEntity)
		return
	}

	log.Infof("Serving file %s", filename)
	f, err := os.Open(filename)
	if err != nil {
		http.Error(w, "unable to open file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename='%s'", filename))
	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "unable to send file", http.StatusUnprocessableEntity)
		return
	}
	f.Close()

	log.Infof("Removing file %s", filename)
	err = os.Remove(filename)
	if err != nil {
		log.Errorf("can't remove file: %v", err)
		return
	}
}

func getFilename(id string) (string, error) {
	cmd := exec.Command(youtubeCommand, id, "-f", bestFormat, "--get-filename")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func downloadFile(id string) error {
	cmd := exec.Command(youtubeCommand, id, "-f", bestFormat)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
