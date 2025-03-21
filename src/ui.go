package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

type PromptHandler struct{
	Label string 
	Items []string
}
// NewPromptHandler 创建一个新的 PromptHandler 实例
func NewPromptHandler(label string ,items []string ) *PromptHandler {
	return &PromptHandler{
		Label:  label,
		Items:  items,
	}
}

func (p *PromptHandler) Select() (string, error ){
	prompt := promptui.Select{
		Label: p.Label,
		Items: p.Items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("选择失败: %v", err)
	}

	return result, err

}

// Input 获取用户输入
func (p *PromptHandler) Input(validateFunc promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    p.Label,
		Validate: validateFunc,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("输入失败: %v", err)
	}

	return result, nil
}

// Confirm 获取用户确认
func (p *PromptHandler) Confirm() (bool, error) {
	prompt := promptui.Prompt{
		Label:     p.Label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("确认失败: %v", err)
	}

	return result == "y" || result == "Y", nil
}