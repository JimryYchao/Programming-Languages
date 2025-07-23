
##### source

- [CMath.c](./csrc/CMath.c)
- [test_CMath.lua](./test_CMath.lua)


##### build C module 

```powershell
$ gcc -shared -o 'CMath.dll' 'csrc\CMath.c' -I"include" -L"." -lLuaDll
$ lua test_CMath.lua
```