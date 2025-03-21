# Go-hitsz-login

一个基于 chronium + chromedp 的 HITsz 校园网自动登录命令行工具。

## 使用方式

### 源码运行方式

首先需要确保本机上安装了 chronium 内核的浏览器。

```shell
git clone https://github.com/yMsquare/go-hitsz-login.git
cd go-hitsz-login
go mod tidy
go run ./src
```

根据命令行提示，输入校园网账号和密码即可完成登录。

![image-20250321200528643](https://s1.vika.cn/space/2025/03/21/7cd66b44c8b84c23ad3f0acdef10bfc1)

## 原理

chromedp 能够以 headless mode 运行一个 chroniuim 内核的浏览器并在其上进行页面操作，所以需要安装 chrome 或 chronium 内核的浏览器。

#### （那么为什么不直接手动点击浏览器进行操作呢？）

主要是自己在宿舍有一台 ubuntu 服务器，平常在教室上课时通过 ssh 连接到服务器进行操作。由于泥深校园网具有一个账号同时只能登录 3 台设备的限制，经常发现服务器没连上校园网需要重新登录，但常规 ssh 连接方式并不支持图形化界面，而自己比较懒不想配置 vnc 等方式（其实是配置了但失败了或效果不好），所以希望提供一种可以通过命令行交互完成的登录方式。
