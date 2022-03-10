//             ,%%%%%%%%,
//           ,%%/\%%%%/\%%
//          ,%%%\c "" J/%%%
// %.       %%%%/ o  o \%%%
// `%%.     %%%%    _  |%%%
//  `%%     `%%%%(__Y__)%%'
//  //       ;%%%%`\-/%%%'
// ((       /  `%%%%%%%'
//  \\    .'          |
//   \\  /       \  | |
//    \\/攻城狮保佑) | |
//     \         /_ | |__
//     (___________)))))))                   `\/'
/*
 * 修订记录:
 * long.qian 2022-01-17 15:16 创建
 */

/**
 * @author long.qian
 */

package file_watch

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	go_lazy_util "github.com/go-lazy-frame/go-lazy-util"
	"github.com/toolkits/file"
	"os"
	"strings"
	"sync"
	"time"
)

// 非 linux 下的处理
func (receiver *FileWatch) nonLinuxEventHandler(watchDirs []string) {
	if receiver.EnableFileCreateHandler {
		receiver.newFileCacheForNonLinux = new(sync.Map)
		// 文件创建缓存处理
		go func() {
			for {
				receiver.newFileCacheForNonLinux.Range(func(key, value interface{}) bool {
					filePath := key.(string)
					fileSize := value.(int64)
					size := go_lazy_util.FileUtil.FileSize(filePath)
					if size == fileSize {
						// 新文件
						receiver.newFileCacheForNonLinux.Delete(key)
						receiver.fileHandlerChannel <- &fileEvent{
							FilePath: filePath,
							Op:       Created,
						}
					} else {
						receiver.newFileCacheForNonLinux.Store(filePath, size)
					}
					return true
				})
				time.Sleep(time.Duration(time.Millisecond.Nanoseconds() * receiver.WriteTime))
			}
		}()

	}
	var err error
	receiver.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		_ = receiver.watcher.Close()
	}()

	go func() {
		for {
			select {
			case event, ok := <-receiver.watcher.Events:
				if !ok {
					fmt.Println("获取目录监听事件通道失败")
					time.Sleep(time.Second)
					continue
				}

				// 监听创建事件
				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Rename == fsnotify.Rename {
					if !file.IsFile(event.Name) {
						dir := event.Name
						err := receiver.watcher.Add(dir)
						if err != nil {
							fmt.Println("新目录：", dir, "，监听失败", err)
						} else {
							fmt.Println("成功监听目录：", dir)
						}
					} else {
						if !receiver.EnableFileCreateHandler {
							continue
						}
						receiver.newFileCacheForNonLinux.Store(event.Name, go_lazy_util.FileUtil.FileSize(event.Name))
					}
				}

				// 监听写操作
				if event.Op&fsnotify.Write == fsnotify.Write {
					if !receiver.EnableFileWriteHandler {
						continue
					}
					receiver.fileHandlerChannel <- &fileEvent{
						FilePath: event.Name,
						Op:       Write,
					}
				}

				// 监听删除操作
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if !receiver.EnableFileDelHandler {
						continue
					}
					receiver.fileHandlerChannel <- &fileEvent{
						FilePath: event.Name,
						Op:       Remove,
					}
				}
			case err, ok := <-receiver.watcher.Errors:
				fmt.Println("目录监听错误 ", err, ok)
				if !ok {
					fmt.Println("获取目录监听事件通道失败")
					time.Sleep(time.Second)
				}
			}
		}
	}()

	var done = make(chan bool)

	var listeningCount int32
	for _, dir := range watchDirs {
		dir := strings.TrimSpace(dir)
		if dir == "" {
			continue
		}
		homeDir, e := os.UserHomeDir()
		if e != nil {
			fmt.Println(e)
		}
		if strings.Contains(dir, "~") {
			dir = strings.Replace(dir, "~", homeDir, 1)
		}
		if file.IsExist(dir) {
			if file.IsFile(dir) {
				fmt.Println("路径：", dir, "，是一个文件，无法监听")
			} else {
				receiver.fsnotifyListening(dir, &listeningCount)
			}
		} else {
			fmt.Println("目录：", dir, "，不存在，监听失败")
		}

	}
	if listeningCount == 0 {
		done <- true
	} else {
		fmt.Println("共监听目录 ", listeningCount, " 个")
	}

	<-done
}
