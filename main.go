package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/timwhitez/Doge-PX/SigF"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)



func parseTagToBytes(tagStr string) []byte {
	tagStr = strings.Replace(tagStr, ` \x`, " ", -1)
	tagStr = strings.Replace(tagStr, `\x`, " ", -1)
	tagStr = strings.Replace(tagStr, `, `, " ", -1)
	tagStr = strings.TrimSpace(tagStr)
	tagSplit := strings.Split(tagStr, " ")
	data := make([]byte, len(tagSplit))
	for i := range tagSplit {
		bigint := new(big.Int)
		bigint.SetString(tagSplit[i], 16)
		data[i] = bigint.Bytes()[0]
	}
	return data
}

// 从test.txt中读取base64字符串，解码，然后生成文件
func base64ToFile(tByte string)[]byte {
	decodeData, err := base64.StdEncoding.DecodeString(tByte)
	if err != nil {
		panic(err)
	}
	tmp, _ := ReadZipFile(decodeData)
	return tmp
}

func ReadZipFile(zByte []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zByte), int64(len(zByte)))
	for _, zipFile := range zipReader.File {
		if strings.Contains(zipFile.Name,".exe"){
			f, err := zipFile.Open()
			if err != nil {
				return nil, err
			}
			defer f.Close()
			return ioutil.ReadAll(f)
		}
	}
	return nil, err
}


func zipData(Src []byte) []byte{
	// Create a buffer to write our archive to.
	fmt.Println("we are in the zipData function")
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)


	zipFile, err := zipWriter.Create("tmp.exe")
	if err != nil {
		fmt.Println(err)
	}

	_, err = zipFile.Write(Src)
	if err != nil {
		fmt.Println(err)
	}

	// Make sure to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		fmt.Println(err)
	}

	return buf.Bytes()

}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h"{
		fmt.Println("the Doge Packer for eXecutables")
		fmt.Println("Doge-PX usage:")
		fmt.Println("        ./DPX.exe target.exe\n")
		fmt.Println("sleep for evasion:")
		fmt.Println("        ./DPX.exe target.exe sleep\n")
		os.Exit(0)
	}
	if len(os.Args) == 3 && strings.ToLower(os.Args[2]) != "sleep"{
		fmt.Println("the Doge Packer for eXecutables")
		fmt.Println("Doge-PX usage:")
		fmt.Println("        ./DPX.exe target.exe\n")
		fmt.Println("sleep for evasion:")
		fmt.Println("        ./DPX.exe target.exe sleep\n")
		os.Exit(0)
	}

	//从tmp常量里读取并解压被写入文件
	f := bytes.NewReader(base64ToFile(tmpexe))

	if len(os.Args) == 3{
		//从tmp常量里读取并解压被写入文件
		f = bytes.NewReader(base64ToFile(sleepexe))
	}


	//读取原exe文件
	rawExeBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	//对原exe文件压缩
	rawExeBytes = zipData(rawExeBytes)

	//将压缩包加密后插入tmp.exe
	tempPEBytes, err := SigF.Inject(f, rawExeBytes, parseTagToBytes("11 ff a1 d3 11 ff a1 d3"), []byte("DPXpasswd"))
	if err != nil {
		fmt.Println(err)
		return
	}
	filename := strings.Trim(os.Args[1],".")
	filename = strings.Trim(filename,"\\")
	filename = strings.Trim(filename,"/")

	//生成加壳后的文件
	err = os.WriteFile("dpx_"+filename, tempPEBytes, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("generate %s successfully\n", "dpx_"+filename)
}
