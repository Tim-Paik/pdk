package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/mholt/archiver/v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Install(packages []string, repoName string) (err error) {
	var (
		repo  Repositories
		index []int
	)
	if len(packages) == 0 {
		return fmt.Errorf("Error: no targets specified\n")
	}
	if repo, err = Read(RepoRoot + "/" + repoName + ".json"); err != nil {
		fmt.Println("Warning: Repo.json was not found, please use 'pdk update <URL>' to install repo")
		return err
	}
	if err := CheckArch(&repo); err != nil {
		return err
	}
	fmt.Println(Indent1 + "Searching for packages")
	if index, err = Search(&repo, packages); err != nil {
		return err
	}
	for _, i := range index {
		//PATH := PackageRoot + "/" + repo.Pkgs[i].Name + "-" + repo.Pkgs[i].Version + ".tar" + path.Ext(repo.Pkgs[i].URL)
		PATH := filepath.Base(repo.Pkgs[i].URL)
		if !IsExists(PATH) {
			fmt.Println(Indent2 + "Downloading " + repo.Pkgs[i].Name + "-" + repo.Pkgs[i].Version)
			if _, err := Download(repo.Pkgs[i].URL, PATH); err != nil {
				return err
			}
		} else {
			fmt.Println(Indent2 + "Find the local package: " + repo.Pkgs[i].Name + "-" + repo.Pkgs[i].Version)
		}
		var MD5 string
		if MD5, err = Md5Sum(PATH); err != nil {
			return err
		}
		if MD5 != repo.Pkgs[i].Md5 {
			fmt.Println(Indent2 + "Error: Md5 error")
			fmt.Println(Indent2 + "Want: " + repo.Pkgs[i].Md5)
			fmt.Println(Indent2 + "Get: " + MD5)
			fmt.Println(Indent2 + "Re-downloading " + repo.Pkgs[i].Name + "-" + repo.Pkgs[i].Version)
			if _, err := Download(repo.Pkgs[i].URL, PATH); err != nil {
				return err
			}
			if MD5, err := Md5Sum(PATH); err != nil {
				return err
			} else if MD5 != repo.Pkgs[i].Md5 {
				return fmt.Errorf("Error: Serious md5 error\n")
			}
		}
		fmt.Println(Indent2 + "Installing " + repo.Pkgs[i].Name + "-" + repo.Pkgs[i].Version)
		if err := UnpackAndCallback(PATH, repo.Pkgs[i].Name); err != nil {
			return err
		}
		fmt.Println(Indent2 + "Successful installation of " + repo.Pkgs[i].Name)
	}
	return nil
}

func Search(repo *Repositories, packages []string) (index []int, err error) {
	for _, packagesName := range packages {
		var pkgIndex []int
		for i, finish := 0, false; finish == false; i++ {
			if repo.Pkgs[i].Name == packagesName {
				pkgIndex = append(pkgIndex, i)
				finish = true
			} else if i == (len(repo.Pkgs) - 1) {
				return nil, fmt.Errorf("Error: Package %s not found\n", packagesName)
			}
		}
		index = append(index, pkgIndex...)
	}
	return index, nil
}

func Download(URL, PATH string) (written int64, err error) {
	var (
		resp *http.Response
		data *os.File
		//dataSize int64
	)

	if resp, err = http.Get(URL); err != nil {
		return 0, fmt.Errorf("Error: Sending request error\n")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	//if dataSize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64); err != nil {
	//return 0, err
	//}

	if data, err = os.Create(PATH); err != nil {
		return 0, err
	}
	defer func() {
		if err := data.Close(); err != nil {
			return
		}
	}()

	if written, err = io.Copy(data, resp.Body); err != nil {
		return 0, err
	}
	return written, nil
}

func Md5Sum(filePath string) (MD5string string, err error) {
	checkMD5 := md5.New()

	if file, err := os.Open(filePath); err != nil {
		return "", err
	} else {
		if _, err := io.Copy(checkMD5, file); err != nil {
			return "", err
		}
	}

	MD5string = hex.EncodeToString(checkMD5.Sum(nil))
	return MD5string, nil
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func CheckArch(repo *Repositories) (err error) {
	if repo.Arch == runtime.GOARCH && repo.OS == runtime.GOOS {
		return nil
	}
	if _, err := fmt.Printf("Warn: You are using %s instead of %s in %s\n", repo.OS+"/"+repo.Arch, runtime.GOOS+"/"+
		runtime.GOARCH, repo.Name); err != nil {
		return err
	}
	return nil
}

func UnpackAndCallback(PATH string, packageName string) (err error) {
	if err := archiver.Unarchive(PATH, AppRoot); err != nil {
		return err
	}
	//CALLBACK_SCRIPT
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command(AppData + "/" + packageName + "/install")
		out, _ := cmd.CombinedOutput()
		fmt.Println(string(out))
	}
	return nil
}
