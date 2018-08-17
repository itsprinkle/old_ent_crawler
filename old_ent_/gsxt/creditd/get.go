package creditd

import (
	"errors"
	"log"
	"time"

	"gsxt/credit"
)

var (
	InfoTimeout = time.Minute * 2
)

// 获取企业信息
func getInfo(credit credit.Credit, keyword, key string) (info credit.InfoV2, err error) {
	var errChan = make(chan error)
	var okChan = make(chan struct{})
	go func() {
		ent, err := credit.Enterprise(keyword, key)
		if err != nil {
			errChan <- err
			return
		}
		info.Enterprise = ent
		close(okChan)
	}()

	info.Business, err = credit.Business(keyword, key)
	if err != nil {
		return
	}

	select {
	case err = <-errChan:
	case <-okChan:
	case <-time.After(InfoTimeout):
		err = errors.New("超时没有返回")
	}
	return
}

func GetInfo(credit credit.Credit, keyword, key string, detailCount int) (info credit.InfoV2, err error) {
	for j := 0; j < detailCount; j++ {
		info, err = getInfo(credit, keyword, key)
		if err == nil {
			break
		}
	}
	return
}

// 获取企业信息
// keyword 搜索关键字
// searchCount 尝试搜索次数
// detailCount 获取详细信息的次数
// resultCount 返回结果条数
func Get(c credit.Credit, keyword string, searchCount, detailCount, resultCount int) (keys []string, infos []credit.InfoV2, err error) {
	if resultCount < 0 {
		resultCount = 10000
	}

	for i := 0; i < searchCount; i++ {
		keys, err = c.Search(keyword)
		if err != nil {
			if err == credit.ErrNotFound || err == credit.ErrOutOfLimit {
				return
			}
			continue
		}

		// 尝试再一次
		if len(keys) == 0 {
			continue
		}

		// 如果没有打开详情搜索, 直接返回
		if detailCount == 0 {
			return
		}

		for _, key := range keys {
			var info credit.InfoV2
			for j := 0; j < detailCount; j++ {
				info, err = getInfo(c, keyword, key)
				if err == nil {
					break
				}
				log.Printf("WARN: 搜索%s出错 - %v", keyword, err)
			}
			if err != nil {
				return
			}
			infos = append(infos, info)
			if len(infos) >= resultCount {
				keys = keys[:len(infos)]
				break
			}
		}
		break
	}
	return
}
