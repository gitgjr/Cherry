package bridge

//bridge package is for connect front-end and back-end

import (
	"encoding/json"
	"main/utils"
	"os"
)

var dataPath = utils.RootPath() + "/data"
var webPath = utils.RootPath() + "/web"

// GetMP4Files create a list of video file name for client
func GetMP4Files() {
	mp4Files, err := utils.FindFiles(dataPath, "", ".mp4")
	if err != nil {
		panic(err)
	}
	jsonData, err := json.Marshal(mp4Files)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(webPath + "/" + "MP4Files.json")

	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(jsonData)

}
