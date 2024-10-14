package internal

//import (
//	_ "github.com/3d0c/gmf"
//)

//func getVideoDuration(filePath string) (float64, error) {
//	//if err := gmf.Init(); err != nil {
//	//	return 0, err
//	//}
//	//defer gmf.Release()
//	//
//	//type probeData struct {
//	//	Streams []struct {
//	//		CodecType string `json:"codec_type"`
//	//		Duration  string `json:"duration"`
//	//	} `json:"streams"`
//	//}
//	//var pd probeData
//	//if err := json.Unmarshal([]byte(data), &pd); err != nil {
//	//	return 0, err
//	//}
//	//
//	//for _, stream := range pd.Streams {
//	//	if stream.CodecType == "video" {
//	//		duration, err := strconv.ParseFloat(stream.Duration, 64)
//	//		if err != nil {
//	//			return 0, err
//	//		}
//	//		return duration, nil
//	//	}
//	//}
//	//return 0, fmt.Errorf("no video stream found")
//}
