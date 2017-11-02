package main;
import (
    "fmt"
    "io"
  
    "os"
    "path"
    "net/http"
    "errors"
	"strings"
	"image"
	"image/draw"
    "image/gif"
    "image/jpeg"
	"image/png"
	
	// "github.com/nfnt/resize"
)
const html = `<html>
<head></head>
<body>
	<form method="post" enctype="multipart/form-data">
		<input type="file" name="image" />
		<input type="submit" />
	</form>
</body>
</html>`
func main() {
    http.HandleFunc("/upload/", uploadHandle) // 上传
    http.HandleFunc("/uploaded/", showPicHandle)  //显示图片
    err := http.ListenAndServe(":443", nil)
    fmt.Println(err)
}
 
// 上传图像接口
func uploadHandle (w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/html")
 
    req.ParseForm()
    if req.Method != "POST" {
        w.Write([]byte(html))
    } else {
        // 接收图片
        uploadFile, handle, err := req.FormFile("image")
        errorHandle(err, w)
 
        // 检查图片后缀
        ext := strings.ToLower(path.Ext(handle.Filename))
        if ext != ".jpg" && ext != ".png" {
            errorHandle(errors.New("只支持jpg/png图片上传"), w);
            return
            //defer os.Exit(2)
        }
 
        // 保存图片
        os.Mkdir("./uploaded/", 0777)
        saveFile, err := os.OpenFile("./uploaded/" + handle.Filename, os.O_WRONLY|os.O_CREATE, 0666);
        errorHandle(err, w)
        io.Copy(saveFile, uploadFile);
 
        defer uploadFile.Close()
        defer saveFile.Close()
        // 上传图片成功
        w.Write([]byte("success: <a target='_blank' href='/uploaded/" + handle.Filename + "'>" + handle.Filename + "</a>"));
    }
}
 
// 显示图片接口
func showPicHandle( w http.ResponseWriter, req *http.Request ) {
    file, err := os.Open("." + req.URL.Path + "13654118.jpg")
	errorHandle(err, w);
	// fmt.Printf(req.URL.Path)
	defer file.Close()
	
	src, err := GetImageObj("." + req.URL.Path + "1.png")
    if err != nil {
		 errorHandle(err, w);
		 fmt.Printf("Error1")
	}
	srcB := src.Bounds().Max

	src1, err := GetImageObj("." + req.URL.Path + "2.png")
    if err != nil {
		errorHandle(err, w);
		fmt.Printf("Error2")
    }
	src1B := src.Bounds().Max

	src2, err := GetImageObj("." + req.URL.Path + "3.png")
    if err != nil {
		 errorHandle(err, w);
		 fmt.Printf("Error1")
	}
	src2B := src.Bounds().Max

	src3, err := GetImageObj("." + req.URL.Path + "4.png")
    if err != nil {
		 errorHandle(err, w);
		 fmt.Printf("Error1")
	}
	src3B := src.Bounds().Max

	src4, err := GetImageObj("." + req.URL.Path + "5.png")
    if err != nil {
		 errorHandle(err, w);
		 fmt.Printf("Error1")
	}
	src4B := src.Bounds().Max

	src5, err := GetImageObj("." + req.URL.Path + "6.png")
    if err != nil {
		 errorHandle(err, w);
		 fmt.Printf("Error1")
	}
	src5B := src.Bounds().Max
	
	newWidth := srcB.X + src1B.X + src2B.X + src3B.X + src4B.X + src5B.X
    newHeight := srcB.Y
    if src1B.Y > newHeight {
        newHeight = src1B.Y
	}
	
	des := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板
	srcWidth := srcB.X
	draw.Draw(des, des.Bounds(), src, src.Bounds().Min, draw.Over)                      //首先将一个图片信息存入jpg
	draw.Draw(des, image.Rect(srcWidth, 0, newWidth, src1B.Y), src1, image.ZP, draw.Over) //将另外一张图片信息存入jpg
	draw.Draw(des, image.Rect(srcWidth * 2, 0, newWidth, src1B.Y), src2, image.ZP, draw.Over) //将另外一张图片信息存入jpg
	draw.Draw(des, image.Rect(srcWidth * 3, 0, newWidth, src1B.Y), src3, image.ZP, draw.Over) //将另外一张图片信息存入jpg
	draw.Draw(des, image.Rect(srcWidth * 4, 0, newWidth, src1B.Y), src4, image.ZP, draw.Over) //将另外一张图片信息存入jpg
	draw.Draw(des, image.Rect(srcWidth * 5, 0, newWidth, src1B.Y), src5, image.ZP, draw.Over) //将另外一张图片信息存入jpg
	
	fSave, err := os.Create("." + req.URL.Path + "7.png")
	if err != nil {
		errorHandle(err, w);
		fmt.Printf("Error3")
	}

	defer fSave.Close()

	var opt jpeg.Options
	opt.Quality = 80

	// newImage := resize.Resize(1024, 0, des, resize.Lanczos3)

	err = jpeg.Encode(fSave, des, &opt) // put quality to 80%
	if err != nil {
		errorHandle(err, w);
		fmt.Printf("Error4")
	}
	// buff, err := ioutil.ReadAll(src)
	
    // errorHandle(err, w);
    // w.Write(buff)
}
 
// 统一错误输出接口
func errorHandle(err error, w http.ResponseWriter) {
    if  err != nil {
		w.Write([]byte(err.Error()))
		fmt.Printf("Error")
    }
}

func GetImageObj(filePath string) (img image.Image, err error) {
    f1Src, err := os.Open(filePath)

    if err != nil {
        return nil, err
    }
    defer f1Src.Close()

    buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
    _, err = f1Src.Read(buff)

    if err != nil {
        return nil, err
    }

    filetype := http.DetectContentType(buff)

    fmt.Println(filetype)

    fSrc, err := os.Open(filePath)
    defer fSrc.Close()

    switch filetype {
    case "image/jpeg", "image/jpg":
        img, err = jpeg.Decode(fSrc)
        if err != nil {
            fmt.Println("jpeg error")
            return nil, err
        }

    case "image/gif":
        img, err = gif.Decode(fSrc)
        if err != nil {
            return nil, err
        }

    case "image/png":
        img, err = png.Decode(fSrc)
        if err != nil {
            return nil, err
        }
    default:
        return nil, err
    }
    return img, nil
}
