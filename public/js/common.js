
function delayMs (ms) {
    
    return new Promise(function(res) {
        setTimeout(res, ms)
    })
}

function toggleDisplayById(tagId){

    if(document.getElementById(tagId).style.display == "none"){

        document.getElementById(tagId).style.display = "block"
        
    } else {

        document.getElementById(tagId).style.display = "none"
    }


}

function getNewDateSortedList(flag, fieldName, orgList){


    let newList = []

    if(flag != "asc" && flag != "desc"){
        return orgList
    }



    for(let i = 0; i < orgList.length; i ++){


        let d1 = Date.parse(orgList[i][fieldName])

        if(newList.length == 0){

            newList.push(orgList[i])
        
            continue
        } 

        for(let j = 0; j < newList.length; j ++){


            let d2 = Date.parse(newList[j][fieldName])
            

            if(flag == "asc"){

                if(d1 < d2){

                    newList.splice(j, 0, orgList[i])

                    break
                }

                if (d1 >= d2 && (j != newList.length - 1)){

                    continue

                } else if (d1 >= d2 && (j == newList.length - 1)){

                    newList.push(orgList[i])

                    break
                }


            } else if(flag == "desc"){


                if(d1 > d2){

                    newList.splice(j, 0, orgList[i])

                    break
                }

                if (d1 <= d2 && (j != newList.length - 1)){

                    continue

                } else if (d1 <= d2 && (j == newList.length - 1)){

                    newList.push(orgList[i])

                    break
                }


            }




        }

    }


    return newList


  


}