# go-spawn

简易的 spawn 工具，主要为了解决 spawn 命令在切换用户后就无法使用的问题。此工具的全部命令都会在一个 shell 中执行。

对 spawn 命令语法稍作改变。



## usage

config：

```json
{
  "addr": "aaa.bbb.ccc.ddd:22",
  "user": "XXXXX",
  "password": "",
  "private_key": "",
  "ssh_dir_path": "",
  "timeout": -1,
  "command_chain_path": "config/commands.txt"
}
```

commands.txt：

```
expect "$" { cd .. }
expect "$" { su your-account }
expect "Password:" { your-password }
expect "$" { cd XXXX/ }
expect "$" { ./YYYY.sh args }
expect "$" { su super-user }
expect "#" { docker ps -a }
expect "#" { docker exec -it 92a2024fa59f bash }
expect "#" { whoami }
```

