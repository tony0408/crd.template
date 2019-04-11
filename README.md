Data Collection

1.0 New Exchange Platform

    1.1 Develop the Basic Functions
    
        1.1.1 Duplicate Exchange Template
            1.1.1.1 Duplicate Folder [/data.binance/coin.realtime.data/exchange/blank]
            1.1.1.2 Rename blank folder -> ["exchange name"] and Put the folder under [/data.binance/coin.realtime.data/exchange/]
            1.1.1.3 Rename [blank.go] that is under [/data.binance/coin.realtime.data/exchange/"exchange name"] to ["exchange name".go]
            
        1.1.2 Duplicate Exchange Test Case Template
            1.1.2.1 Duplicate File [/data.binance/coin.realtime.data/test/blank_test.go]
            1.1.2.2 Rename [blank.go] -> ["exchange name"_test.go] and Put the file under [/data.binance/coin.realtime.data/test/]
            1.1.2.3 Replace [blank] to ["exchange name"]
            1.1.2.4 Replace [Blank] to ["Exchange Name"]
            1.1.2.5 Replace [BLANK] to ["EXCHANGE NAME"]
            
        1.1.3 Develop Basic Functions
            1.1.3.1 Follow the instruction on each file which is under [/data.binance/coin.realtime.data/exchange/"exchange name"]
            
        1.1.4 Test Basic Functions
            1.1.4.1 Run Each Test Case to Make Sure the function is working
            
    1.2 Get RealTime Data
    
        1.2.1 Init Exchange Config
            1.2.1.1 Modify [main.go]
                1.2.1.1.1 Add Function [init"ExchangeName"()]
                1.2.1.1.2 Modify Config Content
                1.2.1.1.3 Add [init"ExchangeName"()] in Init()
                
        1.2.2 Initial New Exchange Task
            1.2.2.1 Modify [init_task.go]
                1.2.2.1.1 Get Exchange Config [e"ExchangeName" := exMan.Get(exchange."EXCHANGENAME")]
                1.2.2.1.2 Call InitTask Function [m.InitTask(e"ExchangeName".GetPairs(), exchange."EXCHANGENAME", [pairs_amount])]
                
        1.2.3 Implement Gaining RealTime Data
            1.2.3.1 Modify [data.go]
                1.2.3.1.1 Get Exchange Config [e"ExchangeName" := exMan.Get(exchange."EXCHANGENAME")]
                1.2.3.1.2 Call InitTask Function [m.InitTask(e"ExchangeName".GetPairs(), exchange."EXCHANGENAME", [pairs_amount])]
                
        1.2.4 Deploy the program on Server