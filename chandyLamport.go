package main

import "time"

// Control Message is -1 

type Node struct {
    id int
    recorded bool
    inCh []*(chan int)
    outCh []*(chan int)
    marker map[*(chan int)]bool
    state map[*(chan int)][]int
}

func receiveMessage(node *Node, ch *(chan int)) {
    for {
        select {
        case m := <- *ch:
            if m != -1 { // receives a basic message
                println("Node", node.id, "received basic message", m)
//                 println("Node", node.id, "Recorded:", node.recorded, "Marker", node.marker[ch])
                if node.recorded && !node.marker[ch] {
                    node.state[ch] = append(node.state[ch], m)
                    println("Node", node.id, "updated the channel state")
                } else {
                    println("Node", node.id, "updated the channel state to ∅")
                }
            } else {
                println("Node", node.id, "received control message")
                TakeSnapshot(node)
                node.marker[ch] = true

                flag := 0
                for _, val := range node.marker {
                    if val { flag = 1 } 
                    else {
                        flag = 0
                    }
                }
                if flag == 1 {
                    println("Node", node.id, "successfully sent markers to all outgoing edges!")
//                     return // terminate
                }
            }
        }
    }
}

func Receive(node *Node) {
    for i := 0; i < len(node.inCh); i++ {
        go receiveMessage(node, node.inCh[i])
    }
}

func sendMarker(ch *(chan int)) {
    *ch <- -1
}

func TakeSnapshot(node *Node) {
    if !node.recorded {
        println("Node", node.id, "took local snapshot")
        node.recorded = true
        for i := 0; i < len(node.outCh); i++ {
            go sendMarker(node.outCh[i])           
        }
    }
}


func main(){
    // initialize nodes & channels
    chPQ := make(chan int, 10)
    chPR := make(chan int, 10)
    chRP := make(chan int, 10)
    chQR := make(chan int, 10)
                
    nP := Node{0, false,
               []*chan int{&chRP}, 
               []*chan int{&chPQ, &chPR},
               map[*(chan int)]bool{&chRP: false},
               map[*(chan int)][]int{&chRP: {}},
              }
    nQ := Node{1, false,
               []*chan int{&chPQ}, 
               []*chan int{&chQR},
               map[*(chan int)]bool{&chPQ: false},
               map[*(chan int)][]int{&chPQ: {}},
              }
    nR := Node{2, false,
               []*chan int{&chPR, &chQR}, 
               []*chan int{&chRP},
               map[*(chan int)]bool{&chPR: false, &chQR: false},
               map[*(chan int)][]int{&chPR: {}, &chQR: {}},
              }
    
    go Receive(&nP)
    go Receive(&nQ)
    go Receive(&nR)
    
    TakeSnapshot(&nP)
    chPQ <- 500        
    time.Sleep(3 * time.Second)
    chQR <- 200    
    time.Sleep(3 * time.Second)
    
}