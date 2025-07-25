---
title: "为什么Athens重要？"
date: 2018-11-06T13:58:58-07:00
weight: 4

---

### 不变性

Go社区的许多问题都是由于库（library）的消失或者在没有告警的情况下突然变化所引起的。上游的软件包维护人员很容易对他们的代码进行更改，但这可能会破坏您的代码，而大多数情况下这是一个意外！
 如果您的软件使用的某个依赖项执行下列操作，该软件的构建是否会中断？

- 提交（Commit） `abdef` 被删除了
- 标签（Tag） `v0.1.0` 被强制推送（push）
- 源码库被完全删除

由于应用程序的依赖项直接来自VCS（版本控制系统，如Github），因此上述情况都可能发生在您身上，并且当它们发生时，您的软件的构建过程可能会中断-哦，不！
Athens通过将代码从VCS复制到_不可变_存储中来解决这些问题。

在这种方式下，您就不需要手动将任何内容上传到Athens后端存储。Go第一次向Athens请求依赖包时，Athens会从VCS（Github、Bitbucket等）获取。但一旦检索到该模块，它将永远保存在Athens的后端存储中，并且代理将不再返回到VCS中获取同一版本的依赖包。这就是雅典如何实现模块不变性。需要注意的是，后端存储掌握在您的手中。


### 逻辑 

go命令行现在可以ping _您自己的服务器_ 来下载依赖项，这意味着您可以编写任何需要的逻辑来提供这种依赖项。包括访问控制（下面将讨论）、添加自定义版本、自定义分支和自定义包等。例如，Athens提供了一个[验证钩子（hook）](https://github.com/leimeng-go/athens/blob/main/config.dev.toml#L127)，每个模块下载时都会调用它来确定是否应该下载此模块。因此，您可以用自己的逻辑扩展athens，比如扫描模块路径或代码以查找标红代码等。


### 性能 

从Athens下载存储的依赖关系_显著_比从版本控制系统下载依赖更快。这是因为`go get`默认情况下使用VCS的下载模块，例如`git clone`。而`go get`启用GOPROXY时，将使用HTTP下载zip压缩包。因此，根据您的计算机和网络连接速度，从GitHub下载CockroachDB源代码zip文件只需要10秒，但git clone需要将近4分钟。


### 访问控制 

比软件包消失更糟糕的是，软件包可能是恶意的。为了确保您的团队或公司不会安装此类恶意软件包，当go命令行请求一个被排除的模块（恶意软件）时，您可以让代理服务器返回500。这将导致构建失败，因为Go需要200 HTTP响应码。使用Athens，您可以通过过滤器（filter）文件实现此目的。


### Vendor 目录成为可选

有了不变的、高性能和高鲁棒性的代理服务器，用户不再需要在每个库中都将vendor目录纳入其版本控制。`go.sum`文件确保在第一次安装之后不会处理任何包。此外，您的CI/CD在每次构建安装所有依赖项时都只需要很短的时间。
