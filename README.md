# SealDice

![Software MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
![SealDice](https://img.shields.io/badge/SealDice-blue)

海豹 TRPG 骰点核心，开源跑团辅助工具，支持 QQ/Kook/Discord 等。

轻量 · 易用 · 全能

## SealDice Project

本仓库为 SealDice 的主仓库，采用 monorepo 模式合并了以前的多个子仓库，便于管理代码。该仓库下每个文件夹为一个独立子仓库：

- `core`：对应 [原 core 仓库](https://github.com/sealdice/sealdice-core)，Go 后端代码，为海豹核心的主要逻辑部分；
- `ui`：对应 [原 UI 仓库](https://github.com/sealdice/sealdice-ui)，前端代码，基于 Vue3 + ElementPlus 开发；
- `android`：对应 [原 Android 仓库](https://github.com/sealdice/sealdice-android)，Android 应用源码；
- `builtins`: 对应 [原 builtins 仓库](https://github.com/sealdice/sealdice-builtins)，存储了海豹的内置牌堆与查询资料；
- `manual`：对应原 [手册仓库](https://github.com/sealdice/sealdice-manual-next)，官方手册源码，由 VitePress 驱动；
- `story-painter`：对应原 [染色器仓库](https://github.com/sealdice/story-painter)，跑团 Log 染色器前端源码；
- `verify`：对应 [原 verify 仓库](https://github.com/sealdice/sealdice-verify)，客户端验证工具前端源码；
- ……

你可以阅读各个子文件夹下的 `README.md` 文件，了解对应仓库的详情。

注：如无特殊说明，所有代码文件均遵循 MIT 开源协议。

## 文档

见 [使用手册](https://sealdice.github.io/sealdice-manual-next/)。

## 开发依赖

本项目提供了 [go-task](https://taskfile.dev/installation) 配置，方便执行预置好的任务。采用 [lefthook](https://github.com/evilmartians/lefthook#install) 配置预提交钩子，用于在提交前检测代码。在开发之前，**强烈推荐** 安装相应工具。

此外，对于不同的子项目，你可能需要配置相应的开发环境（如 `core` 需要配置 `golang` 环境，前端需要配置 `node` 等），请参考各语言相关自行配置。在配置时，一些子项目依赖额外的 Lint 工具，也强烈建议安装：

- 对于 Go 项目（如 `core`），请安装 [golangci-lint]((https://golangci-lint.run/welcome/install/#local-installation)) 进行静态分析。

```bash
# mac 下使用 brew 安装：
brew install go-task lefthook golintci-lint
```

更多信息请参考各子文件夹下的 `README.md` 文件。

## 贡献指南

见 [手册 - 参与项目](https://sealdice.github.io/sealdice-manual-next/about/develop.html)。
