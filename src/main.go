package main

import (
	"context"
	"fmt"
	"os"
	"strings"
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

	lang, err := NewPromptHandler("Choose Language / 选择语言", []string{"English", "中文"}).Select()

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

	// 使用 chromedp 开启浏览器
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	err = chromedp.Run(ctx,
		chromedp.Navigate("http://10.248.98.2/srun_portal_pc?ac_id=1&theme=basic4"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println(i18n["navigate_page"])
			// return SaveScreenshot(ctx, "00-navigate.png")
			return nil
		}),
		chromedp.WaitVisible(`h3.title`, chromedp.ByQuery),
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
				fmt.Println(`login user id:`, loginUserID)
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
			}

			input, err := logoutPrompt.Run()
			if err != nil {
				fmt.Println("Prompt failed")
				return
			}

			input = strings.ToLower(strings.TrimSpace(input))
			if input == "yes" {
				// 执行注销逻辑
				fmt.Println(i18n["check_logout"])
				chromedp.Run(ctx,
					chromedp.Click(`#logout`, chromedp.ByID),
					chromedp.WaitVisible(`button.btn-confirm`, chromedp.ByQuery), // 等待按钮出现
					chromedp.Sleep(1*time.Second),                                // 等待弹窗出现
					chromedp.ActionFunc(func(ctx context.Context) error {
						return SaveScreenshot(ctx, "cdp.png")
					}),
					chromedp.Click(`button.btn-confirm`, chromedp.ByQuery), // 点击按钮
					chromedp.Sleep(3*time.Second),                          // 等待弹窗出现
				)
				var successVisible bool
				err := chromedp.Run(ctx,
					chromedp.Evaluate(`!!document.querySelector(".alert.alert-success")`, &successVisible),
				)

				if err != nil {
					fmt.Println(i18n["logout_failed"], err)
				} else if successVisible {
					fmt.Println(i18n["logout_success"])
				} else {
					fmt.Println(i18n["logout_unknown"])
				}
			} else if input == "no" || input == "n" {
				fmt.Println(i18n["no_logout"])
			}
			fmt.Println(`Bye`)
			return
		}
	}

	loginMethod, err := NewPromptHandler(i18n["login_method"], []string{i18n["local"], i18n["unified"]}).Select()
	// 检查用户选择的登录方式

	// username
	// password
	// login_submit

	accountID, err := NewPromptHandler(i18n["enter_account"], []string{}).Input()
	password, err := NewPromptHandler(i18n["enter_password"], []string{}, WithMask('*')).Input()

	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			if loginMethod == i18n["unified"] {
				err = chromedp.Run(ctx,
					chromedp.WaitVisible(`button.btn.btn-sso`, chromedp.ByQuery), // 等待按钮出现
					chromedp.Click(`button.btn.btn-sso`, chromedp.ByQuery),       // 等待按钮出现
					//https://ids.hit.edu.cn/authserver/login?service=http%3A%2F%2F10.248.98.2%2Fsrun_portal_sso
					chromedp.Sleep(1*time.Second),
				)
			} else if loginMethod == i18n["local"] {
				err = chromedp.Run(ctx,
					chromedp.WaitVisible(`button.btn.btn-account`, chromedp.ByQuery), // 等待按钮出现
					chromedp.Click(`button.btn.btn-account`, chromedp.ByQuery),
					chromedp.Sleep(1*time.Second),
				)
			}
			return err
		}),
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println(`now entering accountID and password...`)
			// return SaveScreenshot(ctx, "01-login-form.png")
			return nil
		}),
		chromedp.SendKeys(`#username`, accountID),
		chromedp.SendKeys(`#password`, password),

		chromedp.ActionFunc(func(ctx context.Context) error {
			if loginMethod == i18n["unified"] {
				err = chromedp.Run(ctx,
					chromedp.WaitVisible(`#login_submit`, chromedp.ByID), // 等待按钮出现
					chromedp.Click(`login_submit`, chromedp.ByID),        // 等待按钮出现
					//https://ids.hit.edu.cn/authserver/login?service=http%3A%2F%2F10.248.98.2%2Fsrun_portal_sso
					chromedp.Sleep(2*time.Second),
				)
			} else if loginMethod == i18n["local"] {
				err = chromedp.Run(ctx,
					chromedp.WaitVisible(`#login-account`, chromedp.ByID), // 等待按钮出现
					chromedp.Click(`#login-account`, chromedp.ByID),
					chromedp.Sleep(2*time.Second),
				)
			}
			return err
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
			// return SaveScreenshot(ctx, "03-login.png")
			return nil
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
