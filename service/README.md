
https://www.jianshu.com/p/eee8a7de179c


检查plist语法是否正确

``
plutil /Library/LaunchDaemons/io.airdb.noah.plist
``

添加自启动项

``
//launchctl工具提供了一系列接口方便使用launchd程序
launchctl load  /Library/LaunchDaemons/io.airdb.noah.plist
``

启动自启动项
``
launchctl start aria2
``

删除自启动项
`
launchctl unload ~/Library/LaunchAgents/aria2.plist
`

查看当前所有自启动项
`
launchctl list
`
