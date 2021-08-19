import NonFungibleToken from "./NonFungibleToken.cdc"

// IdleMysitcHeros
//
pub contract IdleMysitcHeros: NonFungibleToken {

    // Events
    //
    pub event ContractInitialized()
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
    pub event Minted(id: UInt64)
    pub event ParentIDAddedToParent(parentA: UInt64, parentB: UInt64)

    // Named Paths
    //
    pub let CollectionStoragePath: StoragePath
    pub let CollectionPublicPath: PublicPath
    pub let MinterStoragePath: StoragePath
    pub let AdminStoragePath: StoragePath

    access(self) var infos: @{UInt64: HeroInfoData}

    pub var nextInfoID: UInt64
    // totalSupply
    // The total number of IdleMysitcHeros that have been minted
    //
    pub var totalSupply: UInt64

    pub resource Admin {

        pub fun createHeroInfoData(birth:String, batch:Int, first:Bool) {
            // Create the new Match
            var newHeroInfoData <- create HeroInfoData(birth: birth, inherit:false, 
                                                       parentAID:(0 as UInt64), parentBID:(0 as UInt64), parentA_dna:"", parentB_dna:"",
                                                       batch:batch, first:first)

            // Store it in the matchs mapping field
            IdleMysitcHeros.infos[newHeroInfoData.infoID] <-! newHeroInfoData
        }

        pub fun borrowHeroInfoData(infoID: UInt64): &HeroInfoData {
            pre {
                IdleMysitcHeros.infos[infoID] != nil: "Cannot borrow Match: The Match doesn't exist"
            }
            
            // Get a reference to the infos and return it
            // use `&` to indicate the reference to the object and type
            return &IdleMysitcHeros.infos[infoID] as &HeroInfoData
        } 

        pub fun inherit_hero(birth: String, parentAID: UInt64, parentBID:UInt64) {
            pre {
                IdleMysitcHeros.infos[parentAID] != nil: "Cannot borrow info: The info doesn't exist"
                IdleMysitcHeros.infos[parentBID] != nil: "Cannot borrow info: The info doesn't exist"
            }

            let pA = &IdleMysitcHeros.infos[parentAID] as &HeroInfoData
            let pB = &IdleMysitcHeros.infos[parentBID] as &HeroInfoData

            var newHeroInfoData <- create HeroInfoData(birth: birth, inherit:true, 
                                                      parentAID:pA.infoID, parentBID:pB.infoID,parentA_dna: pA.dna, parentB_dna: pB.dna,
                                                      batch:10, first:false)

            // Store it in the matchs mapping field
            IdleMysitcHeros.infos[newHeroInfoData.infoID] <-! newHeroInfoData
        } 
    }

    pub resource HeroInfoData {
        pub let infoID: UInt64

        pub let dna: String

        pub let birth: String

        pub var inheritCount: Int

        pub var parentA: UInt64

        pub var parentB: UInt64

        pub var childs: [UInt64]

        init(birth: String, inherit: Bool, parentAID:UInt64, parentBID:UInt64, 
             parentA_dna:String, parentB_dna:String, batch:Int, first:Bool) {
            self.infoID = IdleMysitcHeros.nextInfoID
            self.birth = birth
            self.inheritCount = 0
            self.childs = []

            fun get_rand_int(x: Int, step: Int): Int {
                if x == 0 {
                    return step
                }
                var v:Int = Int(unsafeRandom())
                var i = 0
                while v < 0 {
                    v = Int(unsafeRandom())
                    i = i+1
                    if i>10{
                        v = 0
                    }
                }
                let mold = Int(v)%x + step
                return Int(mold)
            }

            fun str_to_int(s:String):Int{
                switch s {
                    case "0":
                        return 0
                    case "1":
                        return 1
                    case "2":
                        return 2
                    case "3":
                        return 3
                    case "4":
                        return 4
                    case "5":
                        return 5
                    case "6":
                        return 6
                    case "7":
                        return 7
                    case "8":
                        return 8
                    case "9":
                        return 9
                }
                return 0
            }

            fun strs_to_int(str:String):Int{
                var a = 0
                var list: [Int] = []
                while a < str.length {
                    var s:String = str.slice(from: a, upTo: a+1)
                    a = a + 1
                    var v = str_to_int(s:s)
                    list.insert(at:0, v)
                }
                var res = 0
                var step = 1
                for i in list{
                res = res + i * step
                step = step*10
                }
                return res
            }

            if inherit == false {
                self.parentA = 0
                self.parentB = 0

                fun get_heroDNA():String {
                    fun get_hero_unit(batch:Int): String {
                        var unit:Int = batch + 90 + get_rand_int(x:3, step:0)
                        var res:String = unit.toString()
                        return res
                    }
                    
                    fun get_hero_camp(): String{
                        var rand:Int = get_rand_int(x:6, step:11)
                        var res:String = rand.toString()
                        return res
                    }

                    fun get_hero_attr(): String {
                        var health:Int = get_rand_int(x:17, step:27)
                        var speed:Int = get_rand_int(x:17, step:27)
                        var sum_avg = Int((140 - health - speed) / 2)
                        var skill:Int = 0
                        if sum_avg >= 35 {
                            var _scope = (43 - sum_avg) * 2
                            skill = 43 - get_rand_int(x:_scope, step:0)
                        }else{
                            var _scope = (sum_avg - 27) * 2
                            skill = 27 + get_rand_int(x:_scope, step:0)
                        }
                        var mood:Int = 140 - health - speed - skill
                        var total:Int = health+speed+skill+mood
                        if total != 140 {
                            panic("total error")
                        }
                        var res:String = health.toString()
                        res = res.concat(speed.toString())
                        res = res.concat(skill.toString())
                        res = res.concat(mood.toString())
                        return res
                    }

                    fun get_heroShow(first:Bool): String {

                        fun get_hero_head(first:Bool): Int {
                            let is_legend:Int = get_rand_int(x:100, step:1)
                            if is_legend <= 15 && first == true {
                                return 51
                            }
                            let res:Int =  10 + get_rand_int(x:4, step:1)
                            return res
                        }

                        fun get_hero_hand(first:Bool): Int {
                            let is_legend:Int = get_rand_int(x:100, step:1)
                            if is_legend <= 15 && first == true {
                                return 61
                            }
                            let res:Int =  20 + get_rand_int(x:6, step:1)
                            return res
                        }

                        fun get_hero_body(first:Bool): Int {
                            let is_legend:Int = get_rand_int(x:100, step:1)
                            if is_legend <= 15 && first == true {
                                return 71
                            }
                            let res:Int =  30 + get_rand_int(x:6, step:1)
                            return res
                        }

                        fun get_hero_weapon(first:Bool): Int {
                            let is_legend:Int = get_rand_int(x:100, step:1)
                            if is_legend <= 15 && first == true {
                                return 81
                            }
                            let res:Int =  40 + get_rand_int(x:6, step:1)
                            return res
                        }

                        fun get_hero_plat(): Int {
                            let res:Int = get_rand_int(x:6, step:1)
                            return res
                        }

                        fun get_hero_flag(): Int {
                            let res:Int = get_rand_int(x:6, step:1)
                            return res
                        }

                        var res_pre:String = get_hero_head(first:first).toString()
                        res_pre = res_pre.concat(get_hero_hand(first:first).toString()) 
                        res_pre = res_pre.concat(get_hero_body(first:first).toString()) 
                        res_pre = res_pre.concat(get_hero_weapon(first:first).toString()) 
                        let plat:String = "0".concat(get_hero_plat().toString())
                        let flag:String = "0".concat(get_hero_flag().toString())
                        let res_last:String = plat.concat(flag)
                        let res:String = res_pre.concat(res_last)
                        return res
                    }
                    
                    fun get_hero_skill(): String {
                        var s_head:Int = get_rand_int(x:6, step:1)
                        var s_hand:Int = get_rand_int(x:6, step:1)
                        var s_body:Int = get_rand_int(x:6, step:1)
                        var s_weapon:Int = get_rand_int(x:6, step:1)
                        var res:String = "0"
                        res = res.concat(s_head.toString())
                        res = res.concat("0")
                        res = res.concat(s_hand.toString())
                        res = res.concat("0")
                        res = res.concat(s_body.toString()) 
                        res = res.concat("0")
                        res = res.concat(s_weapon.toString()) 
                        return res
                    }

                    var showD:String = get_heroShow(first:first)
                    var skillD:String = get_hero_skill()
                    var showR1:String = get_heroShow(first:false)
                    var skillR1:String = get_hero_skill()
                    var showR2:String = get_heroShow(first:false)
                    var skillR2:String = get_hero_skill()

                    fun duplication(a:String, b:String, c:String): Bool {
                        let x = strs_to_int(str:a)
                        let y = strs_to_int(str:b)
                        let z = strs_to_int(str:c)
                        if x == y || x == z ||y == z {
                            return true
                        }
                        return  false
                    }

                    var showflag = duplication(a:showD, b:showR1, c:showR2)
                    var skillflag = duplication(a:skillD, b:skillR1, c:skillR2)
                    var i = 0
                    while showflag == true {
                        showD = get_heroShow(first:first)
                        showR1 = get_heroShow(first:false)
                        showR2 = get_heroShow(first:false)
                        showflag = duplication(a:showD, b:showR1, c:showR2)
                        if i > 26{
                            break
                        }
                        i = i+1
                    }
                    var j = 0
                    while skillflag == true {
                        skillD = get_hero_skill()
                        skillR1 = get_hero_skill()
                        skillR2 = get_hero_skill()
                        skillflag = duplication(a:skillD, b:skillR1, c:skillR2)
                        if j > 26{
                            break
                        }
                        j = j+1
                    }
                    
                    // var _dna:String = get_rand_int(x:3, step:0).toString()
                    var _dna:String = batch.toString()
                    _dna = _dna.concat(get_hero_unit(batch:batch)) 
                    _dna = _dna.concat(get_hero_camp())
                    _dna = _dna.concat(get_hero_attr())
                    _dna = _dna.concat(showD) 
                    _dna = _dna.concat(skillD) 
                    _dna = _dna.concat(showR1) 
                    _dna = _dna.concat(skillR1) 
                    _dna = _dna.concat(showR2) 
                    _dna = _dna.concat(skillR2)
                    return _dna
                }
                self.dna = get_heroDNA()
            } else {
                self.parentA = parentAID
                self.parentB = parentBID

                fun get_inherit_DNA(parentA_dna:String, parentB_dna:String):String {
                    let batchA = parentA_dna.slice(from: 0, upTo: 2)  // 2
                    let unitA = parentA_dna.slice(from: 2, upTo: 5)  //3
                    let campA = parentA_dna.slice(from: 5, upTo: 7)  //2
                    let attrA = parentA_dna.slice(from: 7, upTo: 15)  //8
                    let showA = parentA_dna.slice(from: 15, upTo: 27) //12
                    let skillA = parentA_dna.slice(from: 27, upTo: 35) //8
                    let showA_r1 = parentA_dna.slice(from: 35, upTo: 47) //12
                    let skillA_r1 = parentA_dna.slice(from: 47, upTo: 55) //8
                    let showA_r2 = parentA_dna.slice(from: 55, upTo: 67) // 12
                    let skillA_r2 = parentA_dna.slice(from: 67, upTo: 75) // 8

                    let batchB = parentB_dna.slice(from: 0, upTo: 2)  // 2
                    let unitB = parentB_dna.slice(from: 2, upTo: 5)  //3
                    let campB = parentB_dna.slice(from: 5, upTo: 7)  //2
                    let attrB = parentB_dna.slice(from: 7, upTo: 15)  //8
                    let showB = parentB_dna.slice(from: 15, upTo: 27) //12
                    let skillB = parentB_dna.slice(from: 27, upTo: 35) //8
                    let showB_r1 = parentB_dna.slice(from: 35, upTo: 47) //12
                    let skillB_r1 = parentB_dna.slice(from: 47, upTo: 55) //8
                    let showB_r2 = parentB_dna.slice(from: 55, upTo: 67) // 12
                    let skillB_r2 = parentB_dna.slice(from: 67, upTo: 75) // 8

                    let batchC:String = "10"
                    fun get_unit(unitA:String, unitB:String):String{
                        var r:Int = get_rand_int(x:100000, step:1)
                        if r <= 70000 {
                            return unitA
                        } else {
                            return unitB
                        }
                    }
                    let unitC:String = get_unit(unitA:unitA, unitB:unitB)

                    fun get_camp(campA:String, campB:String):String{
                        var r:Int = get_rand_int(x:100000, step:1)
                        if r <= 50000 {
                            return campA
                        } else {
                            return campB
                        }
                    }
                    let campC:String = get_camp(campA:campA, campB:campB)

                    fun get_attr(attrA:String, attrB:String):String{
                        let attrA_health = strs_to_int(str:attrA.slice(from: 0, upTo: 2))
                        let attrA_speed = strs_to_int(str:attrA.slice(from: 2, upTo: 4))
                        let attrA_skill = strs_to_int(str:attrA.slice(from: 4, upTo: 6))
                        let attrA_mood = strs_to_int(str:attrA.slice(from: 6, upTo: 8))
                        let attrB_health = strs_to_int(str:attrB.slice(from: 0, upTo: 2))
                        let attrB_speed = strs_to_int(str:attrB.slice(from: 2, upTo: 4))
                        let attrB_skill = strs_to_int(str:attrB.slice(from: 4, upTo: 6))
                        let attrB_mood = strs_to_int(str:attrB.slice(from: 6, upTo: 8))

                        var attrC_health = (attrA_health + attrB_health)/2
                        var attrC_speed = (attrA_speed + attrB_speed)/2
                        var attrC_skill = (attrA_skill + attrB_skill)/2
                        var attrC_mood = (attrA_mood + attrB_mood)/2
                        var attr_total = attrC_health + attrC_speed + attrC_skill + attrC_mood
                        var add = 140 - attr_total
                        while add > 0 {
                            add = add - 1
                            if attrC_health < 43 {
                                attrC_health = attrC_health + 1
                                continue
                            }
                            if attrC_speed < 43 {
                                attrC_speed = attrC_speed + 1
                                continue
                            }
                            if attrC_skill < 43 {
                                attrC_skill = attrC_skill + 1
                                continue
                            }          
                            if attrC_mood < 43 {
                                attrC_mood = attrC_mood + 1
                                continue
                            }
                            }
                            var res:String = attrC_health.toString()
                            res = res.concat(attrC_speed.toString())
                            res = res.concat(attrC_skill.toString())
                            res = res.concat(attrC_mood.toString())
                        return res
                    }
                    let attrC:String = get_attr(attrA:attrA, attrB:attrB)


                    fun legend_to_normal(s:String):String {
                        switch s {
                        case "51":
                            return "11"
                        case "61":
                            return "21"
                        case "71":
                            return "31"
                        case "81":
                            return "41"
                        default:
                            return s
                        }
                    }

                    fun get_show_skill(showA:String, skillA:String, showB:String, skillB:String, 
                                        showA_r1:String, skillA_r1:String, showA_r2:String, skillA_r2:String,
                                        showB_r1:String, skillB_r1:String, showB_r2:String, skillB_r2:String):[String]{
                        // show skill
                        var showA_head = showA.slice(from: 0, upTo: 2)
                        var showA_hand = showA.slice(from: 2, upTo: 4)
                        var showA_body = showA.slice(from: 4, upTo: 6)
                        var showA_weapon = showA.slice(from: 6, upTo: 8)
                        var showA_platform = showA.slice(from: 8, upTo: 10)
                        var showA_flag = showA.slice(from: 10, upTo: 12)
                        showA_head = legend_to_normal(s:showA_head)
                        showA_hand = legend_to_normal(s:showA_hand)
                        showA_body = legend_to_normal(s:showA_body)
                        showA_weapon = legend_to_normal(s:showA_weapon)
                        var skillA_head = skillA.slice(from: 0, upTo: 2)
                        var skillA_hand = skillA.slice(from: 2, upTo: 4)
                        var skillA_body = skillA.slice(from: 4, upTo: 6)
                        var skillA_weapon = skillA.slice(from: 6, upTo: 8)

                        var showB_head = showB.slice(from: 0, upTo: 2)
                        var showB_hand = showB.slice(from: 2, upTo: 4)
                        var showB_body = showB.slice(from: 4, upTo: 6)
                        var showB_weapon = showB.slice(from: 6, upTo: 8)
                        var showB_platform = showB.slice(from: 8, upTo: 10)
                        var showB_flag = showB.slice(from: 10, upTo: 12)
                        showB_head = legend_to_normal(s:showB_head)
                        showB_hand = legend_to_normal(s:showB_hand)
                        showB_body = legend_to_normal(s:showB_body)
                        showB_weapon = legend_to_normal(s:showB_weapon)
                        var skillB_head = skillB.slice(from: 0, upTo: 2)
                        var skillB_hand = skillB.slice(from: 2, upTo: 4)
                        var skillB_body = skillB.slice(from: 4, upTo: 6)
                        var skillB_weapon = skillB.slice(from: 6, upTo: 8)

                        // r1 show skill
                        var showA_r1_head = showA_r1.slice(from: 0, upTo: 2)
                        var showA_r1_hand = showA_r1.slice(from: 2, upTo: 4)
                        var showA_r1_body = showA_r1.slice(from: 4, upTo: 6)
                        var showA_r1_weapon = showA_r1.slice(from: 6, upTo: 8)
                        var showA_r1_platform = showA_r1.slice(from: 8, upTo: 10)
                        var showA_r1_flag = showA_r1.slice(from: 10, upTo: 12)
                        showA_r1_head = legend_to_normal(s:showA_r1_head)
                        showA_r1_hand = legend_to_normal(s:showA_r1_hand)
                        showA_r1_body = legend_to_normal(s:showA_r1_body)
                        showA_r1_weapon = legend_to_normal(s:showA_r1_weapon)
                        var skillA_r1_head = skillA_r1.slice(from: 0, upTo: 2)
                        var skillA_r1_hand = skillA_r1.slice(from: 2, upTo: 4)
                        var skillA_r1_body = skillA_r1.slice(from: 4, upTo: 6)
                        var skillA_r1_weapon = skillA_r1.slice(from: 6, upTo: 8)

                        var showB_r1_head = showB_r1.slice(from: 0, upTo: 2)
                        var showB_r1_hand = showB_r1.slice(from: 2, upTo: 4)
                        var showB_r1_body = showB_r1.slice(from: 4, upTo: 6)
                        var showB_r1_weapon = showB_r1.slice(from: 6, upTo: 8)
                        var showB_r1_platform = showB_r1.slice(from: 8, upTo: 10)
                        var showB_r1_flag = showB_r1.slice(from: 10, upTo: 12)
                        showB_r1_head = legend_to_normal(s:showB_r1_head)
                        showB_r1_hand = legend_to_normal(s:showB_r1_hand)
                        showB_r1_body = legend_to_normal(s:showB_r1_body)
                        showB_r1_weapon = legend_to_normal(s:showB_r1_weapon)
                        var skillB_r1_head = skillB_r1.slice(from: 0, upTo: 2)
                        var skillB_r1_hand = skillB_r1.slice(from: 2, upTo: 4)
                        var skillB_r1_body = skillB_r1.slice(from: 4, upTo: 6)
                        var skillB_r1_weapon = skillB_r1.slice(from: 6, upTo: 8)

                        // r2 show skill
                        var showA_r2_head = showA_r2.slice(from: 0, upTo: 2)
                        var showA_r2_hand = showA_r2.slice(from: 2, upTo: 4)
                        var showA_r2_body = showA_r2.slice(from: 4, upTo: 6)
                        var showA_r2_weapon = showA_r2.slice(from: 6, upTo: 8)
                        var showA_r2_platform = showA_r2.slice(from: 8, upTo: 10)
                        var showA_r2_flag = showA_r2.slice(from: 10, upTo: 12)
                        showA_r2_head = legend_to_normal(s:showA_r2_head)
                        showA_r2_hand = legend_to_normal(s:showA_r2_hand)
                        showA_r2_body = legend_to_normal(s:showA_r2_body)
                        showA_r2_weapon = legend_to_normal(s:showA_r2_weapon)
                        var skillA_r2_head = skillA_r2.slice(from: 0, upTo: 2)
                        var skillA_r2_hand = skillA_r2.slice(from: 2, upTo: 4)
                        var skillA_r2_body = skillA_r2.slice(from: 4, upTo: 6)
                        var skillA_r2_weapon = skillA_r2.slice(from: 6, upTo: 8)

                        var showB_r2_head = showB_r2.slice(from: 0, upTo: 2)
                        var showB_r2_hand = showB_r2.slice(from: 2, upTo: 4)
                        var showB_r2_body = showB_r2.slice(from: 4, upTo: 6)
                        var showB_r2_weapon = showB_r2.slice(from: 6, upTo: 8)
                        var showB_r2_platform = showB_r2.slice(from: 8, upTo: 10)
                        var showB_r2_flag = showB_r2.slice(from: 10, upTo: 12)
                        showB_r2_head = legend_to_normal(s:showB_r2_head)
                        showB_r2_hand = legend_to_normal(s:showB_r2_hand)
                        showB_r2_body = legend_to_normal(s:showB_r2_body)
                        showB_r2_weapon = legend_to_normal(s:showB_r2_weapon)
                        var skillB_r2_head = skillB_r2.slice(from: 0, upTo: 2)
                        var skillB_r2_hand = skillB_r2.slice(from: 2, upTo: 4)
                        var skillB_r2_body = skillB_r2.slice(from: 4, upTo: 6)
                        var skillB_r2_weapon = skillB_r2.slice(from: 6, upTo: 8)

                        fun get_showC_skill(list_skill:[String]):[String]{
                            // 1~37500 A;37501~75000 B;75001 ~ 84375 A_r1;84376~93750 B_r1;93751~96875 A_r2;96876~100000 B_r2
                            var showC_ = ""
                            var skillC_ = ""
                            var r = get_rand_int(x:100000, step:1)
                            if (r <= 37500){
                                showC_ = list_skill[0]
                                skillC_ = list_skill[1]
                            } else if (r>37501 && r<=75000){
                                showC_ = list_skill[2]
                                skillC_ = list_skill[3]
                            } else if (r>75001 && r<=84375){
                                showC_ = list_skill[4]
                                skillC_ = list_skill[5]
                            } else if (r>84376 && r<=93750){
                                showC_ = list_skill[6]
                                skillC_ = list_skill[7]
                            } else if (r>93751 && r<=96875){
                                showC_ = list_skill[8]
                                skillC_ = list_skill[9]
                            } else if (r>96876 && r<=100000){
                                showC_ = list_skill[10]
                                skillC_ = list_skill[11]
                            }
                            return [showC_, skillC_]
                        }
                        fun get_showC_platform():String{
                            // 1~37500 A;37501~75000 B;75001 ~ 84375 A_r1;84376~93750 B_r1;93751~96875 A_r2;96876~100000 B_r2
                            var showC_platform = ""
                            var r = get_rand_int(x:100000, step:1)
                            if (r <= 37500){
                                showC_platform = showA_weapon
                            } else if (r>37501 && r<=75000){
                                showC_platform = showB_weapon
                            } else if (r>75001 && r<=84375){
                                showC_platform = showA_r1_weapon
                            } else if (r>84376 && r<=93750){
                                showC_platform = showB_r1_body
                            } else if (r>93751 && r<=96875){
                                showC_platform = showA_r2_weapon
                            } else if (r>96876 && r<=100000){
                                showC_platform = showB_r2_weapon
                            }
                            return showC_platform
                        }
                        fun get_showC_flag():String{
                            // 1~37500 A;37501~75000 B;75001 ~ 84375 A_r1;84376~93750 B_r1;93751~96875 A_r2;96876~100000 B_r2
                            var showC_flag = ""
                            var r = get_rand_int(x:100000, step:1)
                            if (r <= 37500){
                                showC_flag = showA_weapon
                            } else if (r>37501 && r<=75000){
                                showC_flag = showB_weapon
                            } else if (r>75001 && r<=84375){
                                showC_flag = showA_r1_weapon
                            } else if (r>84376 && r<=93750){
                                showC_flag = showB_r1_body
                            } else if (r>93751 && r<=96875){
                                showC_flag = showA_r2_weapon
                            } else if (r>96876 && r<=100000){
                                showC_flag = showB_r2_weapon
                            }
                            return showC_flag
                        }
                        let showC_skill_head_list = [showA_head, skillA_head, showB_head, skillB_head, showA_r1_head, skillA_r1_head, 
                                                    showB_r1_head, skillB_r1_head, showA_r2_head, skillA_r2_head, showB_r2_head, skillB_r2_head]
                        let showC_skill_head_res = get_showC_skill(list_skill:showC_skill_head_list)
                        let showC_skill_hand_list = [showA_hand, skillA_hand, showB_hand, skillB_hand, showA_r1_hand, skillA_r1_hand, 
                                                    showB_r1_hand, skillB_r1_hand, showA_r2_hand, skillA_r2_hand, showB_r2_hand, skillB_r2_hand]
                        let showC_skill_hand_res = get_showC_skill(list_skill:showC_skill_hand_list)
                        let showC_skill_body_list = [showA_body, skillA_body, showB_body, skillB_body, showA_r1_body, skillA_r1_body, 
                                                    showB_r1_body, skillB_r1_body, showA_r2_body, skillA_r2_body, showB_r2_body, skillB_r2_body]
                        let showC_skill_body_res = get_showC_skill(list_skill:showC_skill_body_list) 
                        let showC_skill_weapon_list = [showA_weapon, skillA_weapon, showB_weapon, skillB_weapon, showA_r1_weapon, skillA_r1_weapon, 
                                                    showB_r1_weapon, skillB_r1_weapon, showA_r2_weapon, skillA_r2_weapon, showB_r2_weapon, skillB_r2_weapon]
                        let showC_skill_weapon_res = get_showC_skill(list_skill:showC_skill_weapon_list)                       
                        let showC_head = showC_skill_head_res[0]
                        let showC_hand = showC_skill_hand_res[0]
                        let showC_body = showC_skill_body_res[0]
                        let showC_weapon = showC_skill_weapon_res[0]
                        let showC_platform = get_showC_platform()
                        let showC_flag = get_showC_flag()

                        let skillC_head = showC_skill_head_res[1]
                        let skillC_hand = showC_skill_hand_res[1]
                        let skillC_body = showC_skill_body_res[1]
                        let skillC_weapon = showC_skill_weapon_res[1]
                        
                        fun get_C_r1_r2(list_show:[String], list_skill:[String], showC_check:String):[String] {
                                
                            var a = strs_to_int(str:showC_check)
                            
                            var showC_part_r1 = ""
                            var showC_part_r2 = ""
                            var skillC_part_r1 = ""
                            var skillC_part_r2 = ""
                            
                            var i = 0
                            var j = 0
                            var b = strs_to_int(str:list_show[i])
                            while a == b && i < list_show.length - 1{
                                i = i + 1
                                b = strs_to_int(str:list_show[i])
                            }
                            var c = strs_to_int(str:list_show[j])
                            while (a == c || b == c) && (j < list_show.length - 1){
                                j = j + 1
                                c = strs_to_int(str:list_show[j])
                            }
                            showC_part_r1 = list_show[i]
                            skillC_part_r1 = list_skill[i]
                            showC_part_r2 = list_show[j]
                            skillC_part_r2 = list_skill[j]
                            return [showC_part_r1, skillC_part_r1, showC_part_r2, skillC_part_r2]
                        }
                        
                        var list_show_head = [showA_head, showB_head, showA_r1_head, showB_r1_head, showA_r2_head, showB_r2_head]
                        var list_skill_head = [skillA_head, skillB_head, skillA_r1_head, skillB_r1_head, skillA_r2_head, skillB_r2_head]
                        let showC_head_res = get_C_r1_r2(list_show:list_show_head, list_skill:list_skill_head, showC_check:showC_head)
                        
                        let showC_head_r1 = showC_head_res[0]
                        let skillC_head_r1 = showC_head_res[1]
                        let showC_head_r2 = showC_head_res[2]
                        let skillC_head_r2 = showC_head_res[3]
                        
                        var list_show_hand = [showA_hand, showB_hand, showA_r1_hand, showB_r1_hand, showA_r2_hand, showB_r2_hand]
                        var list_skill_hand = [skillA_hand, skillB_hand, skillA_r1_hand, skillB_r1_hand, skillA_r2_hand, skillB_r2_hand]
                        let showC_hand_res = get_C_r1_r2(list_show:list_show_hand, list_skill:list_skill_hand, showC_check:showC_hand)
                        let showC_hand_r1 = showC_hand_res[0]
                        let skillC_hand_r1 = showC_hand_res[1]
                        let showC_hand_r2 = showC_hand_res[2]
                        let skillC_hand_r2 = showC_hand_res[3]

                        var list_show_body = [showA_body, showB_body, showA_r1_body, showB_r1_body, showA_r2_body, showB_r2_body]
                        var list_skill_body = [skillA_body, skillB_body, skillA_r1_body, skillB_r1_body, skillA_r2_body, skillB_r2_body]
                        let showC_body_res = get_C_r1_r2(list_show:list_show_body, list_skill:list_skill_body, showC_check:showC_body)
                        let showC_body_r1 = showC_body_res[0]
                        let skillC_body_r1 = showC_body_res[1]
                        let showC_body_r2 = showC_body_res[2]
                        let skillC_body_r2 = showC_body_res[3]  

                        var list_show_weapon = [showA_weapon, showB_weapon, showA_r1_weapon, showB_r1_weapon, showA_r2_weapon, showB_r2_weapon]
                        var list_skill_weapon = [skillA_weapon, skillB_weapon, skillA_r1_weapon, skillB_r1_weapon, skillA_r2_weapon, skillB_r2_weapon]
                        let showC_weapon_res = get_C_r1_r2(list_show:list_show_weapon, list_skill:list_skill_weapon, showC_check:showC_weapon)
                        let showC_weapon_r1 = showC_weapon_res[0]
                        let skillC_weapon_r1 = showC_weapon_res[1]
                        let showC_weapon_r2 = showC_weapon_res[2]
                        let skillC_weapon_r2 = showC_weapon_res[3]

                        fun get_showC_r1_r2_platform():[String] {
                            var list_show_platform = [showA_platform, showB_platform, showA_r1_platform, showB_r1_platform, showA_r2_platform, showB_r2_platform]
                                                        
                            var a = strs_to_int(str:showC_platform)
                            
                            var showC_part_r1 = ""
                            var showC_part_r2 = ""
                            var skillC_part_r1 = ""
                            var skillC_part_r2 = ""
                            
                            var i = 0
                            var j = 0
                            var b = strs_to_int(str:list_show_platform[i])
                            while a == b && i < list_show_platform.length - 1{
                                i = i + 1
                                b = strs_to_int(str:list_show_platform[i])
                            }
                            var c = strs_to_int(str:list_show_platform[j])
                            while (a == c || b == c) && (j < list_show_platform.length - 1){
                                j = j + 1
                                c = strs_to_int(str:list_show_platform[j])
                            }
                            var showC_platform_r1 = list_show_platform[i]
                            var showC_platform_r2 = list_show_platform[j]
                            return [showC_platform_r1, showC_platform_r2]
                        }
                        let showC_platform_res = get_showC_r1_r2_platform()
                        let showC_platform_r1_res = showC_platform_res[0]
                        let showC_platform_r2_res = showC_platform_res[1]

                        fun get_showC_r1_r2_flag():[String] {
                            var list_show_flag = [showA_flag, showB_flag, showA_r1_flag, showB_r1_flag, showA_r2_flag, showB_r2_flag]
                            
                            var a = strs_to_int(str:showC_flag)
                            
                            var showC_part_r1 = ""
                            var showC_part_r2 = ""
                            var skillC_part_r1 = ""
                            var skillC_part_r2 = ""
                            
                            var i = 0
                            var j = 0
                            var b = strs_to_int(str:list_show_flag[i])
                            while a == b && i < list_show_flag.length - 1{
                                i = i + 1
                                b = strs_to_int(str:list_show_flag[i])
                            }
                            var c = strs_to_int(str:list_show_flag[j])
                            while (a == c || b == c) && (j < list_show_flag.length - 1){
                                j = j + 1
                                c = strs_to_int(str:list_show_flag[j])
                            }
                            var showC_flag_r1 = list_show_flag[i]
                            var showC_flag_r2 = list_show_flag[j]
                            return [showC_flag_r1, showC_flag_r2]
                        }
                        let showC_flag_res = get_showC_r1_r2_flag()
                        let showC_flag_r1_res = showC_flag_res[0]
                        let showC_flag_r2_res = showC_flag_res[1]

                        var showC = showC_head.concat(showC_hand)
                        showC = showC.concat(showC_body)
                        showC = showC.concat(showC_weapon)
                        showC = showC.concat(showC_platform)
                        showC = showC.concat(showC_flag)

                        var skillC = skillC_head.concat(skillC_hand)
                        skillC = skillC.concat(skillC_body)
                        skillC = skillC.concat(skillC_weapon)
                    
                        
                        var showC_r1 = showC_head_r1.concat(showC_hand_r1)
                        showC_r1 = showC_r1.concat(showC_body_r1)
                        showC_r1 = showC_r1.concat(showC_weapon_r1)
                        showC_r1 = showC_r1.concat(showC_platform_r1_res)
                        showC_r1 = showC_r1.concat(showC_flag_r1_res)

                        var skillC_r1 = skillC_head_r1.concat(skillC_hand_r1)
                        skillC_r1 = skillC_r1.concat(skillC_body_r1)
                        skillC_r1 = skillC_r1.concat(skillC_weapon_r1)

                        var showC_r2 = showC_head_r2.concat(showC_hand_r2)
                        showC_r2 = showC_r2.concat(showC_body_r2)
                        showC_r2 = showC_r2.concat(showC_weapon_r2)
                        showC_r2 = showC_r2.concat(showC_platform_r2_res)
                        showC_r2 = showC_r2.concat(showC_flag_r2_res)

                        var skillC_r2 = skillC_head_r2.concat(skillC_hand_r2)
                        skillC_r2 = skillC_r2.concat(skillC_body_r2)
                        skillC_r2 = skillC_r2.concat(skillC_weapon_r2)
                        
                        return [showC, skillC, showC_r1, skillC_r1, showC_r2, skillC_r2]
                    }
                    
                    let show_skill = get_show_skill(showA:showA, skillA:skillA, showB:showB, skillB:skillB, 
                                    showA_r1:showA_r1, skillA_r1:skillA_r1, showA_r2:showA_r2, skillA_r2:skillA_r2,
                                    showB_r1:showB_r1, skillB_r1:skillB_r1, showB_r2:showB_r2, skillB_r2:skillB_r2
                                    )
                    let showC = show_skill[0]
                    let skillC = show_skill[1]
                    let showC_r1 = show_skill[2]
                    let skillC_r1 = show_skill[3]
                    let showC_r2 = show_skill[4]
                    let skillC_r2 = show_skill[5]
                    
                    //C = [batchC, unitC, campC, attrC, showC, skillC, showC_r1, skillC_r1, showC_r2, skillC_r2]
                    var dna_res = batchC
                    dna_res = dna_res.concat(unitC)
                    dna_res = dna_res.concat(campC)
                    dna_res = dna_res.concat(attrC)
                    dna_res = dna_res.concat(showC)
                    dna_res = dna_res.concat(skillC)
                    dna_res = dna_res.concat(showC_r1)
                    dna_res = dna_res.concat(skillC_r1)
                    dna_res = dna_res.concat(showC_r2)
                    dna_res = dna_res.concat(skillC_r2)
                    return dna_res
                }
                self.dna = get_inherit_DNA(parentA_dna:parentA_dna, parentB_dna:parentB_dna)
            }
        }

        pub fun checkInheritHero(parentA: UInt64, parentB: UInt64) {
            pre {
                self.parentA == self.parentB: "parentA and parentB must different!"
            }
            emit ParentIDAddedToParent(parentA: parentA, parentB: parentB)
        }

        pub fun riseInheritCount() {
            pre {
                self.inheritCount < 8: "inherit count must less than 7!"
            }
            self.inheritCount = self.inheritCount + 1
        }

        pub fun addChild(childID: UInt64) {
            pre {
                self.inheritCount < 8: "inherit count max cant add child!"
            }
            self.childs.append(childID)
        }
    }

    // NFT
    // A Hero Item as an NFT
    //
    pub resource NFT: NonFungibleToken.INFT {
        // The token's ID
        pub let id: UInt64

        pub let infoID: UInt64

        // initializer
        //
        init(initID: UInt64) {
            self.id = IdleMysitcHeros.nextInfoID
            self.infoID = IdleMysitcHeros.nextInfoID
            IdleMysitcHeros.nextInfoID = IdleMysitcHeros.nextInfoID +  (1 as UInt64)
        }
    }

    // This is the interface that users can cast their IdleMysitcHeros Collection as
    // to allow others to deposit IdleMysitcHeros into their Collection. It also allows for reading
    // the details of IdleMysitcHeros in the Collection.
    pub resource interface HerosCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowHero(id: UInt64): &IdleMysitcHeros.NFT? {
            // If the result isn't nil, the id of the returned reference
            // should be the same as the argument to the function
            post {
                (result == nil) || (result?.id == id):
                    "Cannot borrow Hero reference: The ID of the returned reference is incorrect"
            }
        }
        pub fun borrowHeroInfoData(infoID: UInt64): &HeroInfoData?{
            post {
                (result == nil) || (result?.infoID == infoID):
                    "Cannot borrow Hero Inf oData reference: The ID of the returned reference is incorrect"
            }
        }
    }

    // Collection
    // A collection of Hero NFTs owned by an account
    //
    pub resource Collection: HerosCollectionPublic, NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic {
        // dictionary of NFT conforming tokens
        // NFT is a resource type with an `UInt64` ID field
        //
        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}

        // withdraw
        // Removes an NFT from the collection and moves it to the caller
        //
        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")

            emit Withdraw(id: token.id, from: self.owner?.address)

            return <-token
        }

        // deposit
        // Takes a NFT and adds it to the collections dictionary
        // and adds the ID to the id array
        //
        pub fun deposit(token: @NonFungibleToken.NFT) {
            let token <- token as! @IdleMysitcHeros.NFT

            let id: UInt64 = token.id

            // add the new token to the dictionary which removes the old one
            let oldToken <- self.ownedNFTs[id] <- token

            emit Deposit(id: id, to: self.owner?.address)

            destroy oldToken
        }

        // getIDs
        // Returns an array of the IDs that are in the collection
        //
        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        // borrowNFT
        // Gets a reference to an NFT in the collection
        // so that the caller can read its metadata and call its methods
        //
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            return &self.ownedNFTs[id] as &NonFungibleToken.NFT
        }

        // borrowHero
        // Gets a reference to an NFT in the collection as a Hero,
        // exposing all of its fields (including the typeID).
        // This is safe as there are no functions that can be called on the Hero.
        //
        pub fun borrowHero(id: UInt64): &IdleMysitcHeros.NFT? {
            if self.ownedNFTs[id] != nil {
                let ref = &self.ownedNFTs[id] as auth &NonFungibleToken.NFT
                return ref as! &IdleMysitcHeros.NFT
            } else {
                return nil
            }
        }

        pub fun borrowHeroInfoData(infoID: UInt64): &HeroInfoData? {
            pre {
                IdleMysitcHeros.infos[infoID] != nil: "Cannot borrow InfoData: The InfoData doesn't exist"
            }
            
            // Get a reference to the Match and return it
            // use `&` to indicate the reference to the object and type
            return &IdleMysitcHeros.infos[infoID] as &HeroInfoData
        } 
        

        // destructor
        destroy() {
            destroy self.ownedNFTs
        }

        // initializer
        //
        init () {
            self.ownedNFTs <- {}
        }
    }

    // createEmptyCollection
    // public function that anyone can call to create a new empty collection
    //
    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create Collection()
    }
    
    // NFTMinter
    // Resource that an admin or something similar would own to be
    // able to mint new NFTs
    //
	pub resource NFTMinter {

		// mintNFT
        // Mints a new NFT with a new ID
		// and deposit it in the recipients collection using their collection reference
        //
		pub fun mintNFT(recipient: &{NonFungibleToken.CollectionPublic}) {
            emit Minted(id: IdleMysitcHeros.totalSupply)

			// deposit it in the recipient's account using their reference
			recipient.deposit(token: <-create IdleMysitcHeros.NFT(initID: IdleMysitcHeros.totalSupply))

            IdleMysitcHeros.totalSupply = IdleMysitcHeros.totalSupply + (1 as UInt64)
		}
	}

    // fetch
    // Get a reference to a Hero from an account's Collection, if available.
    // If an account does not have a IdleMysitcHeros.Collection, panic.
    // If it has a collection but does not contain the itemID, return nil.
    // If it has a collection and that collection contains the itemID, return a reference to that.
    //
    pub fun fetch(_ from: Address, itemID: UInt64): &IdleMysitcHeros.NFT? {
        let collection = getAccount(from)
            .getCapability(IdleMysitcHeros.CollectionPublicPath)!
            .borrow<&IdleMysitcHeros.Collection{IdleMysitcHeros.HerosCollectionPublic}>()
            ?? panic("Couldn't get collection")
        // We trust IdleMysitcHeros.Collection.borowHero to get the correct itemID
        // (it checks it before returning it).
        return collection.borrowHero(id: itemID)
    }

    // initializer
    //
	init() {
        self.infos <- {}
        // Set our named paths
        //FIXME: REMOVE SUFFIX BEFORE RELEASE
        self.CollectionStoragePath = /storage/herosCollection7
        self.CollectionPublicPath = /public/herosCollection7
        self.MinterStoragePath = /storage/herosMinter7
        self.AdminStoragePath = /storage/herosAdmin7
        
        // Initialize the total supply
        self.totalSupply = 0
        self.nextInfoID = 1

        // Create a Minter resource and save it to storage
        let minter <- create NFTMinter()
        self.account.save(<-minter, to: self.MinterStoragePath)
        self.account.save<@Admin>(<- create Admin(), to: self.AdminStoragePath)

        emit ContractInitialized()
	}
}
