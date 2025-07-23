function require(modname) end -- 加载模块，返回模块值和加载数据

package = {}
package.cpath = ""                                    -- C加载器搜索路径
package.loaded = {}                                   -- 已加载模块表
package.path = ""                                     -- Lua加载器搜索路径
package.preload = {}                                  -- 特殊模块加载器表
package.config = ""                                   -- 包管理编译期配置信息
package.searchers = {}                                -- 模块搜索方式表（Lua 5.2+）

function package.loadlib(libname, funcname) end       -- 动态链接C库

function package.searchpath(name, path, sep, rep) end -- 在路径中搜索文件

---

package.loaders = {}                -- 模块加载方式表（Lua 5.1）
function package.seeall(module) end -- 设置模块继承全局环境（Lua 5.1）
