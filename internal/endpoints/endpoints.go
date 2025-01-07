package endpoints

import (
	VideoDatabase "RolandQuest/internal/database"
	"RolandQuest/internal/localfiles"
	"RolandQuest/internal/views/pages"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func Home(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, pages.Home())
}

func VideoList(w http.ResponseWriter, r *http.Request) error {
	path, _ := strings.CutPrefix(r.URL.Path, "/videolist")
	return Render(w, r, pages.VideoList(path))
}

func Player(w http.ResponseWriter, r *http.Request) error {
	video_name, _ := strings.CutPrefix(r.URL.Path, "/player")
	
	return Render(w, r, pages.Player(video_name))
}

func VideoServe(w http.ResponseWriter, r *http.Request) error {
	title, _ := strings.CutPrefix(r.URL.Path, "/video//")
	file, err := localfiles.FindMovie(title)
	if err != nil {
		return err
	}
	fmt.Println(file)
	http.ServeFile(w, r, file)
	return nil
}

func HomebreweryServe(w http.ResponseWriter, r *http.Request) error {
	file, err := localfiles.FindImage(r.URL.Path)
	if err != nil {
		return err
	}
	fmt.Println(file)
	http.ServeFile(w, r, file)
	return nil
}

func Alt(w http.ResponseWriter, r *http.Request) error {
	
	pathSlice := strings.Split(r.URL.Path, "/")
	
	if len(pathSlice) == 2 {
		return Render(w, r, pages.AltMain(VideoDatabase.GetSeries()))
	}
	// if len(pathSlice) == 3 {
	// 	return Render(w, r, pages.AltSeries(pathSlice[2], VideoDatabase.GetSeasonsNames(pathSlice[2])))
	// }
	// if len(pathSlice) == 4 {
	// 	return Render(w, r, pages.AltSeason(pathSlice[2], pathSlice[3], VideoDatabase.GetEpisodenames(pathSlice[3])))
	// }
	
	return Render(w, r, pages.AltMain(VideoDatabase.GetSeries()))
}

func AltGet(w http.ResponseWriter, r *http.Request) error {
	file, err := localfiles.FindImage(r.URL.Path)
	if err != nil {
		return err
	}
	fmt.Println(file)
	http.ServeFile(w, r, file)
	return nil
}


func Test(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, pages.Test())
}

func TestUpload(w http.ResponseWriter, r *http.Request) error {
	
	r.ParseMultipartForm(0)
	defer r.MultipartForm.RemoveAll()
	
	fi, info, err := r.FormFile("video")
	if err != nil {
		fmt.Printf("Error: ")
		return nil
	}
	defer fi.Close()
	fmt.Printf("Recieved %v (%v)\n", info.Filename, info.Size)
	
	bytes_to_read := int64(1024 * 1024)
	steps := int(info.Size / bytes_to_read + 1)
	
	os.Create("whereami.avi");
	
	for i := 0; i < steps; i++ {
		data := make([]byte, bytes_to_read)
		fi.Read(data)
		fd, _ := os.OpenFile("whereami.avi", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644);
		fd.Write(data);
		if fd.Close() != nil {
			log.Fatal(err)
		}
	}
	
	
	return nil;
}