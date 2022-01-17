http = require("http")
time = require("time")
config = require("config")
cmd = require("cmd")
util = require("util")

function init()
    version,err = cmd.exec("./ytfs-node --version")
    if err ~= nil then
        print(err)

        print("下载最新ytfs-node...")
        remote_version,download_url = util.get_remote_version_info()
        util.download_ytfs_node(remote_version,download_url)
        print("下载完成")
    else
        print("当前版本:",version)
    end
end

init()

print("启动ytfs-node")
mv("ytfs-node.new","ytfs-node")
cmd.command("./ytfs-node daemon").run()
print("ytfs-node 运行结束")
exit()