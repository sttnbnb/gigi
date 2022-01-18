# gigi
Discord BOT "gigi"

## spec

- golang 17
- [discordgo](https://github.com/bwmarrin/discordgo)

## build and run

### build

```bash
go build
```

### run

```bash
./gigi -token <token> [-guild <guild_id>]
```
If you want to enable Guild Command, you need -guild option.

## 参考
https://github.com/bwmarrin/discordgo/blob/master/examples/slash_commands/main.go
https://github.com/bwmarrin/discordgo/blob/master/examples/components/main.go
