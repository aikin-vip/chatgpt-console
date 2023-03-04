# 使用文档

## 1. 项目介绍
    
    该项目是一个在Terminal中使用chatgpt api进行交互的一个工具。
    本项目代码全部是chatgpt生成的，如果你有其他需要也可以自己让chatgpt接着改进。
    
    
![image](https://github.com/aikin-vip/chatgpt-console/blob/main/preview.gif)

## 2. 项目使用

### 编译
    
    go build chatgpt.go

后续会写好github action，自动编译好可执行文件，然后发布到release中。

### 环境变量

    export OPENAI_API_KEY=你的OPENAI_API_KEY

### 运行

    将编译后的"chat"可执行文件放到你想放到的地方，然后运行即可。


### 快捷键
    
    Mac:

    Enter 输入时有多行时使用

    Control + C 退出

    Control + D 结束输入开始询问chatgpt

---

    Windows:(没有测试如果有问题请告诉我)

    Enter 输入时有多行使用回车换行继续输入

    Control + C 退出

    Control + D 结束输入开始询问chatgpt
