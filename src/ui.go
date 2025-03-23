package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

type PromptHandler struct{
	Label string 
	Items []string
	Mask rune
}

type Option func(*PromptHandler)

// NewPromptHandler 创建一个新的 PromptHandler 实例
func NewPromptHandler(label string ,items []string, option...Option ) *PromptHandler {
	handler := &PromptHandler{
		Label:  label,
		Items:  items,
		Mask:  0,
	}
	for _, opt := range option { 
		opt(handler) 
	}
	return handler
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
func (p *PromptHandler) Input() (string, error) {
	prompt := promptui.Prompt{
		Label:    p.Label,
		Mask:	  p.Mask,
		// Validate: validateFunc,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("输入失败: %v", err)
	}

	return result, nil
}

func WithMask(mask rune) Option {
	return func(s *PromptHandler){
		s.Mask = mask
	}
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