# Khicago | GOT | Preset

Preset 是一个用于读取静态配置表的库，应对与游戏等复杂系统中大量表格化的

##  Invalid Seal 和 Default-able Types

通过 `pseal.Invalid` 获取 Invalid Seal, 对一个 seal, 可以使用 seal.IsInvalid 来判断是否读取异常

如果在创建 seal 时解析失败, 这一格会被置为 `pseal.Invalid` 并返回异常, 这个异常会被 preset 逻辑捕获并以 Warning 形式打印, `pseal.Invalid` 会留在 preset 数据中

当填入的值为空时会特殊处理, 所有的 Default-able Types 会默认解析成这个类型的默认值, 而不会抛出或打印异常

所有 Seal 类型中, 除 PID 以外都是 Default-able Types, 意味着 PID 为空时, 解析在 seal 层会直接返回 Invalid, 其他类型会解析默认值

相关 constrains:
- `must_fill` 可以用来跳过 Default-able Types 默认行为, 指定这个 constrains, Default-able Types 在为空时也会置为 `pseal.Invalid` 并返回异常

## 各子模块定义以及相互依赖关系

- `preset` 用于读取静态配置表, 依赖 `table2d`
- `preset/parse` 用于解析静态配置表, 依赖 `preset`

