time = require("time")
process = require("process")
cmd = require("cmd")
util = require("util")
md5 = require("md5")

while (true) do
    remote_version, download_url, strMd5 = util.get_remote_version_info()
    if download_url ~= "" then
        if remote_version > util.get_current_version() then
            calcMd5, err = md5.CalcFileMd5("ytfs-node")
            if err ~= nil then
                print("update calc md5 err:", err)
                err = util.download_ytfs_node(remote_version,download_url)
                if err == nil then
                    process.pkill("ytfs-node")
                    mv("ytfs-node","ytfs-node.old")
                    print("update 更新完成")
                    break
                else
                    print("update 更新失败",err)
                end
            elseif strMd5 ~= calcMd5 then
                err = util.download_ytfs_node(remote_version,download_url)
                if err == nil then
                    process.pkill("ytfs-node")
                    mv("ytfs-node","ytfs-node.old")
                    print("md5 inconsistent update 更新完成")
                    break
                else
                    print("md5 inconsistent update 更新失败",err)
                end
            else
                print("update ytfs-node local md5:", calcMd5)
                print("update ytfs-node remote md5:", strMd5)
                print("update md5 consistent 不需要更新矿机")
            end
        end
    end

    time.sleep(time.minute * 10)
end

