package main

    pro, _ := os.FindProcess(p.Pid) //查找进程
        fmt.Println(pro)                //&{5488 240 0}

	    err = p.Kill() 　　　　　//杀死进程但不释放进程相关资源
	        fmt.Println(err)

		    err = p.Release() 　　　//<span style="color:#FF0000;">释放进程相关资源，因为资源释放凋之后进程p就不能进行任何操作，此后进程Ｐ的任何操作都会被报错</span>
		        fmt.Println(err)
