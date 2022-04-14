%%writefile echo.go

package main

import (
    "time"
    "os"
)


type Node struct {
    id int
    isInitiator bool
    parent int
    received int
    outCh map[int]*(chan int)
    inCh map[int]*(chan int)
}

func receiveWave(node *Node, ch *(chan int)) {
    for {
        select {
            case from := <- *ch:
            println(node.id, "received wave from", from)
            node.received++
            if node.parent == -1 && !node.isInitiator { // no parent
                node.parent = from
                if (len(node.outCh) > 1) {
                    for key, c := range node.outCh {
                        if key != node.parent {
                
                            *c <- node.id // send wave from current node
                        }
                    }
                } else {
                    *node.outCh[from] <- node.id
                }
            } else if node.received == len(node.inCh) {
                if node.parent != -1 {
                    *node.outCh[node.parent] <- node.id
                } else {
                    println(node.id, "will decide")
                    os.Exit(1)
                }
            }
        }   
    }
}

func Receive(node *Node) {
    println(node.id, "listening for waves")
    for _, ch := range node.inCh {
        go receiveWave(node, ch)
    }
}

// node.id is sent as a wave
func sendWave(id int, to int, ch *(chan int)) {
    println(id, "sending wave to", to)
    *ch <- id
}

func Echo(node *Node) {
    for k, ch := range node.outCh {
        go sendWave(node.id, k, ch)
    }
}

func main() {

    chPQ, chPR := make(chan int), make(chan int)
    chQP, chQR := make(chan int), make(chan int)
    chRP, chRQ := make(chan int), make(chan int)    
        
    nP := Node{0, true, -1, 0,
               map[int]*chan int{1: &chPQ, 2: &chPR}, 
               map[int]*chan int{1: &chQP, 2: &chRP}}
    nQ := Node{1, false, -1, 0, 
               map[int]*chan int{0: &chQP, 2: &chQR}, 
               map[int]*chan int{0: &chPQ, 2: &chRQ}}
    nR := Node{2, false, -1, 0,
               map[int]*chan int{0: &chRP, 1: &chRQ}, 
               map[int]*chan int{0: &chPR, 1: &chQR}}

    go Receive(&nP)
    go Receive(&nQ)
    go Receive(&nR)
    time.Sleep(3 * time.Second)
    
    Echo(&nP)
    
    time.Sleep(3 * time.Second)

}