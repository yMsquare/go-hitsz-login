package main

import (
	"context"
	"fmt"
	"os"
	"time"

	chromedp "github.com/chromedp/chromedp"
	"github.com/manifoldco/promptui"
)

type AlreadyLoginError struct {
	LoginUserID string
}

func (e *AlreadyLoginError) Error() string {
	return fmt.Sprintf(e.LoginUserID)
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

	ctx, cancel := chromedp.NewContext(context.Background())
	err = chromedp.Run(ctx,
		chromedp.Navigate("http://10.248.98.2/srun_portal_pc?ac_id=1&theme=basic4"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println(`now navigating to login page...`)
			return SaveScreenshot(ctx, "00-navigate.png")
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var exists bool
			err := chromedp.Evaluate(`!!document.querySelector("#ipv4")`, &exists).Do(ctx)
			if err != nil {
				return err
			}
			if exists {
				SaveScreenshot(ctx, "debug-before-error.png")
				var loginUserID string
				err := chromedp.Evaluate(`
				(function() {
					let usernameInput = document.querySelector("#username");
					return usernameInput ? usernameInput.textContent: "unknown_user";
				})();
			`, &loginUserID).Do(ctx)

				if err != nil {
					fmt.Println(`error getting login user id:`, err)
					return err
				}
				return &AlreadyLoginError{LoginUserID: loginUserID}
			}
			return nil
		}),
	)
	if err != nil {
		if _, ok := err.(*AlreadyLoginError); ok {
			fmt.Println(i18n["already_login"])
			fmt.Println(`login user:`, err.(*AlreadyLoginError).LoginUserID)

			logoutPrompt := promptui.Prompt{
				Label: i18n["logout_confirm"],
				Validate: func(input string) error {
					if input == "yes" {

						chromedp.Run(ctx,
							chromedp.Click(`#logout`, chromedp.ByID),
							chromedp.WaitVisible(`button.btn-confirm`, chromedp.ByQuery), // 等待按钮出现
							chromedp.Click(`button.btn-confirm`, chromedp.ByQuery),       // 点击按钮
							chromedp.Sleep(4*time.Second), // 等待弹窗出现
						)
						var successVisible bool
						err := chromedp.Run(ctx,
							chromedp.Evaluate(`!!document.querySelector(".alert.alert-success")`, &successVisible),
						)
						
						if err != nil {
							fmt.Println("检查注销状态失败:", err)
						} else if successVisible {
							fmt.Println("✅ 注销成功！")
						} else {
							fmt.Println("❌ 未检测到注销成功提示！")
						}
						return nil
					}
					return fmt.Errorf(i18n["logout_confirm_error"])
				},
			}

			_, err := logoutPrompt.Run()
			if err != nil {
				fmt.Println("Prompt failed")
				return
			}
			fmt.Println(i18n["login_failed"], err)
			return
		}
	}

	promptAccount := promptui.Prompt{
		Label: i18n["enter_account"],
	}

	accountID, err := promptAccount.Run()
	if err != nil {
		fmt.Println("Prompt failed")
		return
	}

	promptPassword := promptui.Prompt{
		Label: i18n["enter_password"],
		Mask:  '*',
	}
	password, err := promptPassword.Run()

	if err != nil {
		fmt.Println("Prompt failed")
		return
	}

	defer cancel()
	err = chromedp.Run(ctx,
		chromedp.Navigate("http://10.248.98.2/srun_portal_pc?ac_id=1&theme=basic4"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println(`now navigating to login page...`)
			return SaveScreenshot(ctx, "00-navigate.png")
		}),

		chromedp.WaitVisible(`button.btn.btn-account`, chromedp.ByQuery), // 等待按钮出现
		chromedp.Click(`button.btn.btn-account`, chromedp.ByQuery),

		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println(`now entering accountID and password...`)
			return SaveScreenshot(ctx, "01-login-form.png")
		}),
		chromedp.SendKeys(`#username`, accountID),
		chromedp.SendKeys(`#password`, password),

		chromedp.Click(`#login-account`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return SaveScreenshot(ctx, "02-fill.png")
		}),

		// **检测是否弹出错误提示框**
		chromedp.Sleep(1*time.Second), // 等待弹窗出现
		chromedp.ActionFunc(func(ctx context.Context) error {
			var dialogVisible bool
			err := chromedp.Evaluate(`!!document.querySelector("div.component.dialog.confirm.active")`, &dialogVisible).Do(ctx)
			if err != nil {
				return err
			}
			if dialogVisible {
				// 读取错误信息
				var errorMsg string
				chromedp.Run(ctx, chromedp.Text(`div.component.dialog.confirm.active div.section`, &errorMsg, chromedp.ByQuery))
				fmt.Println(i18n["login_failed_password"], errorMsg)

				// 截图保存错误提示
				SaveScreenshot(ctx, "error-dialog.png")

				// 关闭错误弹窗
				chromedp.Run(ctx, chromedp.Click(`div.component.dialog.confirm.active button.btn-confirm`, chromedp.ByQuery))

				return fmt.Errorf("Login failed: %s", errorMsg)
			}
			return nil
		}),

		chromedp.WaitVisible(`#ipv4`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return SaveScreenshot(ctx, "03-login.png")
		}),
	)
	if err != nil {
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
