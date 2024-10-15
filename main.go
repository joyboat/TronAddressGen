package main

import (
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"os"
	"strings"
	"sync"
)

// GenerateKey 生成一个 Tron 地址和对应的私钥（WIF 格式）。
func GenerateKey() (address string, wif string) {
	for {
		pri, err := btcec.NewPrivateKey() // 创建新的私钥
		if err != nil {
			continue // 如果出错，重试
		}
		address = addr.PubkeyToAddress(pri.ToECDSA().PublicKey).String() // 根据公钥生成地址
		wif = pri.Key.String()                                           // 获取私钥字符串
		return
	}
}

// isValidSuffix 判断末尾相同字符
func isValidSuffix(address string, endRepeatTimes int) bool {
	lowerAddress := strings.ToLower(address)
	suffix := lowerAddress[len(lowerAddress)-endRepeatTimes:]
	return strings.Count(suffix, string(suffix[0])) == endRepeatTimes
}

// generateBeginAndEndRepeatAccount 生成具有指定开始和结束重复字符模式的地址。
func generateBeginAndEndRepeatAccount(endRepeatTimes, numWorker int) (string, string, error) {
	if endRepeatTimes == 0 {
		return "", "", fmt.Errorf("endRepeatTimes 不能为0")
	}
	var wg sync.WaitGroup
	resultChan := make(chan []string, 1) // 用于接收结果的通道
	stopChan := make(chan struct{})      // 用于通知协程停止的通道
	worker := func() {
		defer wg.Done() // 确保工作结束时减少计数
		for {
			select {
			case <-stopChan:
				return // 如果接收到停止信号，退出协程
			default:
				address, privateKey := GenerateKey()
				if isValidSuffix(address, endRepeatTimes) {
					select {
					case resultChan <- []string{address, privateKey}: // 找到符合条件的地址后发送到通道
						close(stopChan) // 关闭通道，通知其他协程停止
					default:
					}
					return
				}
			}
		}
	}
	// 启动多个并发工作协程
	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		go worker()
	}
	// 等待第一个协程找到符合条件的结果
	result := <-resultChan
	wg.Wait() // 等待所有协程完成
	return result[0], result[1], nil
}

// Product 根据指定的模式生成多个 Tron 地址。
func Product(endTimes int, numAddr, numWorker int) [][]string {
	fileName := fmt.Sprintf("addr_%v.txt", endTimes)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("打开文件出错: %v\n", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	for i := 0; i < numAddr; i++ {
		var tronAddress, privateKey string
		var err error
		if endTimes != 0 {
			// 根据重复字符模式生成地址
			tronAddress, privateKey, err = generateBeginAndEndRepeatAccount(endTimes, numWorker)
		} else {
			fmt.Println("无效的参数。")
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("生成账号出现错误: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n>>> 钱包地址: %v, >>> 私钥: %v\n", tronAddress, privateKey)
		_, err = file.WriteString(fmt.Sprintf("地址: %v >>>> 私钥: %v\n", tronAddress, privateKey))
		if err != nil {
			fmt.Printf("写入文件出错: %v\n", err)
			os.Exit(1)
		}
	}
	return nil
}

func main() {
	endTimes := flag.Int("endTimes", 0, "结束重复次数")
	numAddr := flag.Int("numAddr", 1, "生成账号数量")
	numWorker := flag.Int("numWorker", 8, "并发数")
	flag.Parse()
	fmt.Printf("开始生成靓号 --> endTimes: %v\n", *endTimes)
	Product(*endTimes, *numAddr, *numWorker)
}
