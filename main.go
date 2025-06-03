package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// 版本信息，编译时通过 ldflags 注入
var version = "dev"

// Account 账号信息结构体
type Account struct {
	AuthToken string
	CT0       string
	Line      string // 原始行内容，用于输出
}

// Result 查询结果结构体
type Result struct {
	Account Account
	Status  string
}

func NewTwitterClient(authToken, ct0 string) *resty.Client {
	headers := map[string]string{
		"Accept-Language":    "en-US,en;q=0.8",
		"Authority":          "x.com",
		"Origin":             "https://x.com",
		"Referer":            "https://x.com/",
		"Sec-Ch-Ua":          `"Google Chrome";v="135", "Not;A=Brand";v="8", "Chromium";v="135"`,
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": "Windows",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"Sec-Gpc":            "1",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36",
		"Accept-Encoding":    "gzip, deflate, br",
		"authorization":      "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
	}

	cookies := []*http.Cookie{
		{Name: "auth_token", Value: authToken},
	}

	if ct0 != "" {
		headers["x-csrf-token"] = ct0
		cookies = append(cookies, &http.Cookie{Name: "ct0", Value: ct0})
	}

	client := resty.New().SetTimeout(60 * time.Second).SetHeaders(headers).SetCookies(cookies)
	return client
}

func check(account Account) string {
	error_codes := map[int]string{32: "BAD_TOKEN", 64: "SUSPENDED", 141: "SUSPENDED", 326: "LOCKED", 353: "LOCKED"}
	Twclient := NewTwitterClient(account.AuthToken, account.CT0)
	url := "https://x.com/i/api/1.1/jot/client_event.json"

	resp, err := Twclient.R().Post(url)
	if err != nil {
		fmt.Printf("Error checking account %s: %v\n", account.AuthToken[:10]+"...", err)
		return "ERROR"
	}
	if strings.Contains(resp.String(), "matching csrf cookie") {
		for _, c := range resp.Cookies() {
			if c.Name == "ct0" {
				Twclient.SetCookie(c)
				Twclient.SetHeader("x-csrf-token", c.Value)
				resp, err = Twclient.R().Post("https://x.com/i/api/1.1/jot/client_event.json")
				if err != nil {
					return "ERROR"
				}
				break
			}
		}
	}

	if resp.StatusCode() == 400 {
		return "GOOD"
	}

	code := gjson.Get(resp.String(), "errors.0.code")
	if code.Exists() {
		if status, ok := error_codes[int(code.Int())]; ok {
			return status
		}
	}
	return "ERROR"
}

// loadExistingResults 读取现有的输出文件，获取已经检查过的账号
func loadExistingResults() (map[string]bool, error) {
	checkedAccounts := make(map[string]bool)

	// 所有可能的输出文件
	outputFiles := []string{
		"good_accounts.txt",
		"bad_token_accounts.txt",
		"suspended_accounts.txt",
		"locked_accounts.txt",
		"error_accounts.txt",
	}

	for _, filename := range outputFiles {
		file, err := os.Open(filename)
		if err != nil {
			// 文件不存在是正常的，跳过
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// 解析行中的auth_token来标记为已检查
			parts := strings.Split(line, "----")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if len(part) == 40 && isLowerHex(part) {
					checkedAccounts[part] = true
					break
				}
			}
		}

		file.Close()

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("读取文件 %s 时出错: %v", filename, err)
		}
	}

	return checkedAccounts, nil
}

// readAccountsFromFile 从文件读取账号信息
func readAccountsFromFile(filename string) ([]Account, error) {
	// 首先加载已经检查过的账号
	checkedAccounts, err := loadExistingResults()
	if err != nil {
		return nil, fmt.Errorf("加载现有结果时出错: %v", err)
	}

	if len(checkedAccounts) > 0 {
		fmt.Printf("发现 %d 个已检查的账号，将跳过重复检查\n", len(checkedAccounts))
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 %s: %v", filename, err)
	}
	defer file.Close()

	var accounts []Account
	var skippedCount int
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "----")

		var authToken, ct0 string

		// 遍历所有部分，根据长度和格式判断类型
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			// 判断是否为auth_token（40位小写hex）
			if len(part) == 40 && isLowerHex(part) {
				if authToken == "" {
					authToken = part
				} else {
					fmt.Printf("警告: 第%d行发现多个可能的auth_token，使用第一个: %s\n", lineNum, line)
				}
			}

			// 判断是否为ct0（160位小写hex）
			if len(part) == 160 && isLowerHex(part) {
				if ct0 == "" {
					ct0 = part
				} else {
					fmt.Printf("警告: 第%d行发现多个可能的ct0，使用第一个: %s\n", lineNum, line)
				}
			}
		}

		// 验证必须有auth_token
		if authToken == "" {
			fmt.Printf("警告: 第%d行未找到有效的auth_token，跳过: %s\n", lineNum, line)
			continue
		}

		// 检查是否已经检查过此账号
		if checkedAccounts[authToken] {
			skippedCount++
			continue
		}

		// ct0可以为空
		if ct0 == "" {
			fmt.Printf("信息: 第%d行未找到ct0，将使用空值: %s\n", lineNum, authToken)
		}

		accounts = append(accounts, Account{
			AuthToken: authToken,
			CT0:       ct0,
			Line:      line,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件时出错: %v", err)
	}

	if skippedCount > 0 {
		fmt.Printf("跳过了 %d 个已经检查过的账号\n", skippedCount)
	}

	return accounts, nil
}

// isLowerHex 检查字符串是否为小写十六进制
func isLowerHex(s string) bool {
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}

// writeResultsToFiles 将结果写入不同的文件
func writeResultsToFiles(results []Result) error {
	// 按状态分组
	statusGroups := make(map[string][]Result)
	for _, result := range results {
		statusGroups[result.Status] = append(statusGroups[result.Status], result)
	}

	// 为每个状态创建文件
	for status, group := range statusGroups {
		filename := fmt.Sprintf("%s_accounts.txt", strings.ToLower(status))
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("无法创建文件 %s: %v", filename, err)
		}

		for _, result := range group {
			_, err := file.WriteString(result.Account.Line + "\n")
			if err != nil {
				file.Close()
				return fmt.Errorf("写入文件 %s 时出错: %v", filename, err)
			}
		}

		file.Close()
		fmt.Printf("已将 %d 个 %s 状态账号写入 %s\n", len(group), status, filename)
	}

	return nil
}

// worker 工作协程
func worker(id int, accounts <-chan Account, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for account := range accounts {
		status := check(account)
		results <- Result{
			Account: account,
			Status:  status,
		}
		fmt.Printf("Worker %d 完成检查: %s... -> %s\n", id, account.AuthToken[:10], status)
	}
}

func main() {
	// 检查版本参数
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("Twitter 账号状态检查工具 v%s\n", version)
		fmt.Println("项目地址: https://github.com/your-username/twitter-checker")
		return
	}

	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Printf("Twitter 账号状态检查工具 v%s\n\n", version)
		fmt.Println("使用方法: go run main.go <accounts_file.txt>")
		fmt.Println("文件格式: 每行按----分割，至少包含auth_token")
		fmt.Println("---------------------------------------------------")
		fmt.Println("good_accounts.txt， 正常账号")
		fmt.Println("bad_token_accounts.txt， 无效token")
		fmt.Println("suspended_accounts.txt， 封禁账号")
		fmt.Println("locked_accounts.txt， 锁定账号")
		fmt.Println("error_accounts.txt， 错误账号，请重命名后重新查询")
		fmt.Println("")
		fmt.Println("参数:")
		fmt.Println("  --version, -v    显示版本信息")

		return
	}

	filename := os.Args[1]

	fmt.Printf("Twitter 账号状态检查工具 v%s\n", version)
	fmt.Printf("正在从文件 %s 读取账号信息...\n", filename)

	// 读取账号信息
	accounts, err := readAccountsFromFile(filename)
	if err != nil {
		fmt.Printf("读取账号文件失败: %v\n", err)
		return
	}

	if len(accounts) == 0 {
		fmt.Println("没有找到有效的账号信息")
		return
	}

	fmt.Printf("成功读取 %d 个账号，开始并发检查...\n", len(accounts))

	// 创建通道
	accountChan := make(chan Account, len(accounts))
	resultChan := make(chan Result, len(accounts))

	// 启动10个工作协程
	var wg sync.WaitGroup
	const workerCount = 10

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i+1, accountChan, resultChan, &wg)
	}

	// 发送账号到通道
	for _, account := range accounts {
		accountChan <- account
	}
	close(accountChan)

	// 等待所有工作完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	var results []Result
	for result := range resultChan {
		results = append(results, result)
	}

	fmt.Printf("\n检查完成，共处理 %d 个账号\n", len(results))

	// 统计结果
	statusCount := make(map[string]int)
	for _, result := range results {
		statusCount[result.Status]++
	}

	fmt.Println("\n状态统计:")
	for status, count := range statusCount {
		fmt.Printf("%s: %d 个账号\n", status, count)
	}

	// 写入结果文件
	fmt.Println("\n正在将结果写入文件...")
	err = writeResultsToFiles(results)
	if err != nil {
		fmt.Printf("写入结果文件失败: %v\n", err)
		return
	}

	fmt.Println("\n任务完成！")
}
