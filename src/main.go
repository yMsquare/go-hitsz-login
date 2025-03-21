package main

import (
	"context"
	"fmt"
	"os"

	chromedp "github.com/chromedp/chromedp"
	"github.com/manifoldco/promptui"
);

type AlreadyLoginError struct{}
func (e *AlreadyLoginError) Error() string {
    return i18n["✅ already_login"]
}

func main() {

	langPrompt := promptui.Select{
		Label: "Choose Language / 选择语言",
		Items: []string{"English", "中文"},
	}
	_, lang, err := langPrompt.Run()
	if err != nil {
		fmt.Println("Prompt failed")
		return
	}
	// 加载用户选择的语言包
	langCode := "en" // 默认语言
	if lang == "中文" {
		langCode = "zh"
	}

	err = loadLanguage(langCode)
	if err != nil {
		fmt.Println("Error loading language:", err)
		return
	}

	promptAccount := promptui.Prompt{
        Label: i18n["enter_account"],
    }

	accountID , err := promptAccount.Run()
	if err != nil {
		fmt.Println("Prompt failed")
		return
	}

	promptPassword := promptui.Prompt{
		Label: i18n["enter_password"],
		Mask: '*',
	}
	password , err := promptPassword.Run()

	if err != nil {
		fmt.Println("Prompt failed")
		return
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	err = chromedp.Run(ctx,
		chromedp.Navigate("http://10.248.98.2/srun_portal_pc?ac_id=1&theme=basic4"),
		fmt.Println(`now navigating to login page...`)
		// chromedp.ActionFunc(func(ctx context.Context) error {
			// return SaveScreenshot(ctx, "00-navigate.png")
		// }),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			var exists bool
            err := chromedp.Evaluate(`!!document.querySelector("#ipv4")`, &exists).Do(ctx)
            if err != nil {
                return err
            }
			if exists {
				return &AlreadyLoginError{}
			}
			return nil
		}),

		chromedp.WaitVisible(`button.btn.btn-account`, chromedp.ByQuery), // 等待按钮出现
		chromedp.Click(`button.btn.btn-account`, chromedp.ByQuery), 
		
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		fmt.Println(`now entering accountID and password...`)
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	return SaveScreenshot(ctx, "01-login-form.png")
		// }),
		chromedp.SendKeys(`#username`, accountID), 
		chromedp.SendKeys(`#password`,password),
		
		chromedp.Click(`#login-account`, chromedp.ByID),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	return SaveScreenshot(ctx, "02-fill.png")
		// }),
		chromedp.WaitVisible(`#ipv4`, chromedp.ByID),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	return SaveScreenshot(ctx, "03-login.png")
		// }),
	)	

	if err != nil{
		if _, ok := err.(*AlreadyLoginError); ok {
			fmt.Println(i18n["already_login"])
			return 
		} 
			fmt.Println(i18n["login_failed"], err)
			return
		
	}
	fmt.Println("✅ Successfuly Login! 🎉")
	
}

func SaveScreenshot(ctx context.Context, filename string) error {
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf)); err != nil {
		return err
	}
	return os.WriteFile(filename, buf, 0644)
}