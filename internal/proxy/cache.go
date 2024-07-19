package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

func getCacheDirPath() string {
	uhd, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(uhd, "wakili", "cache")
	}

	f, _ := os.Getwd()
	base := filepath.Base(f)

	return filepath.Join(base, "wakilicache")
}

func InitCache(baseUrl string) error {
	clearExpiredCache()

	baseDirPath := getBaseDirPath(baseUrl)
	if !isExist(baseDirPath) {
		if err := createDir(baseDirPath); err != nil {
			return err
		}
	}

	return nil
}

func getFromCache(baseUrl string, r *http.Request) (*proxyHandlerResponse, error) {
	if isNotCachable(r) {
		return nil, nil
	}

	filePath := getCachedFilePath(baseUrl, getRequestHash(r))
	bytes, err := os.ReadFile(filePath)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	byteStr := strings.Split(string(bytes), "\n")

	ttl, _ := strconv.Atoi(byteStr[0])
	if isCacheExpired(ttl) {
		return nil, nil
	}

	cacheValue := byteStr[1]
	if cacheValue != "" {
		res := new(proxyHandlerResponse)
		if err := json.Unmarshal([]byte(cacheValue), res); err != nil {
			return nil, nil
		}

		return res, nil
	}

	return nil, nil
}

func saveToCache(baseUrl string, r *http.Request, data string, ttl int) error {
	if isNotCachable(r) {
		return nil
	}

	filePath := getCachedFilePath(baseUrl, getRequestHash(r))
	content := fmt.Sprintf("%d\n%s", getExpiryTime(ttl), data)

	err := os.WriteFile(filePath, []byte(content), 0666)
	if err != nil {
		return err
	}

	return nil
}

func isExist(filePath string) bool {
	_, err := os.Stat(filePath)

	if err == nil {
		return true
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return false
}

func createDir(filePath string) error {
	err := os.MkdirAll(filePath, 0755)

	if err != nil {
		return err
	}

	return nil
}

func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func getBaseDirPath(baseUrl string) string {
	baseUrlDirName := getHashNumber(baseUrl)
	baseDirPath := filepath.Join(getCacheDirPath(), baseUrlDirName)

	return baseDirPath
}

func getCachedFilePath(baseUrl string, requestHash string) string {
	filename := fmt.Sprintf("%s.txt", requestHash)

	return filepath.Join(getBaseDirPath(baseUrl), filename)
}

func getHashNumber(data string) string {
	hUrl := fnv.New32()
	hUrl.Write([]byte(data))

	return strconv.Itoa(int(hUrl.Sum32()))
}

func getRequestHash(r *http.Request) string {
	data := fmt.Sprintf("%s%s", r.Method, strings.TrimSpace(r.URL.RequestURI()))

	return getHashNumber(data)
}

func isNotCachable(r *http.Request) bool {
	cachableMethods := []string{http.MethodGet}

	return !slices.Contains(cachableMethods, r.Method)
}

func getExpiryTime(ttl int) int64 {
	if ttl <= 0 {
		return 0
	}

	return time.Now().Add(time.Duration(ttl) * time.Second).Unix()
}

func isCacheExpired(ttl int) bool {
	// if ttl is 0, means data should be cached indefinetely
	if ttl <= 0 {
		return false
	}

	return time.Unix(int64(ttl), 0).Before(time.Now())
}

func clearExpiredCache() {
	cacheDirPath := getCacheDirPath()

	filepath.WalkDir(cacheDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && cacheDirPath != path {

			filepath.WalkDir(
				path,
				func(path string, file fs.DirEntry, err error) error {
					if !file.IsDir() {
						bytes, err := os.ReadFile(path)
						if err == nil {
							bytesArr := strings.Split(string(bytes), "\n")
							ttl, _ := strconv.Atoi(bytesArr[0])
							if isCacheExpired(ttl) {
								deleteFile(path)
							}
						}
					}
					return nil
				},
			)
		}

		return nil
	})
}

/*func purgeCache(baseUrl string) {
	if baseUrl == "" {
		deleteFile(getCacheDirPath())
	} else {
		deleteFile(getBaseDirPath(baseUrl))
	}
}*/
