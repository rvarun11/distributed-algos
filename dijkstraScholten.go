package main

import (
    "time"
//     "os"
)

/*
NOTE: In the implementation below, process is made passive manually.
In a distributed network however, the process becomes automatically and only becomes active again
when a basic message is received.
*/

// For simplicity a map is used instead of array for Tree
var treeMp map[int]bool

type Message struct {
    from int
    to int
    val int // val = -1 for control message; -111 to make process passive
}

type Node struct {
    id int
    parent int // parent id
    isInitiator bool
    cc int
    outCh map[int]*(chan Message)
    inCh map[int]*(chan Message)
}

func receiveMessage(node *Node, ch *(chan Message)) {
    for {
        select {
            case m := <- *ch:
            if m.val == -1 { // control message
                println("Node", node.id, "received a control message from", m.from)
                node.cc--
            } else if treeMp[node.id] {
                println("Node", node.id, "exists in Tree. Sending control message to parent node", node.parent)
                *node.outCh[node.parent] <- Message{node.id, node.parent, -1}
            } else {
                println("Adding node", node.id, "to Tree with parent", m.from)
                treeMp[node.id] = true
                node.parent = m.from
                node.cc = 0
            }
            println("Node", node.id, "cc:", node.cc)
        }   
    }
}

func Receive(node *Node) {
    for _, ch := range node.inCh {
        go receiveMessage(node, ch)
    }
}

func sendMessage(node *Node, ch *(chan Message), m Message) {
    println("Node", node.id, "sends basic message to", m.to)
    node.cc++
    *ch <- m
} 

func processPassive(node *Node) {
     if node.cc == 0 {
        println("Node", node.id, "is passive")

        if node.isInitiator {
            println("Node", node.id, "calls Announce!")
        } else if node.cc == 0 {
            println("Node", node.id, "sends an ACK to", node.parent)
            treeMp[node.id] = false
            *node.outCh[node.parent] <- Message{node.id, node.parent, -1}
        }
    } else {
        println("Node", node.id, "remains in tree")
    }    
}

func main() {

    chPQ, chPR := make(chan Message), make(chan Message)
    chQP, chQR := make(chan Message), make(chan Message)
    chRP, chRQ := make(chan Message), make(chan Message)    
        
    nP := Node{0, -1, true, 0,
               map[int]*chan Message{1: &chPQ, 2: &chPR}, 
               map[int]*chan Message{1: &chQP, 2: &chRP}}
    nQ := Node{1, -1, false, 0, 
               map[int]*chan Message{0: &chQP, 2: &chQR}, 
               map[int]*chan Message{0: &chPQ, 2: &chRQ}}
    nR := Node{2, -1, false, 0,
               map[int]*chan Message{0: &chRP, 1: &chRQ}, 
               map[int]*chan Message{0: &chPR, 1: &chQR}}
    
    treeMp = map[int]bool{0:true, 1:false, 2:false}

    go Receive(&nP)
    go Receive(&nQ)
    go Receive(&nR)
    time.Sleep(3 * time.Second)
    
    sendMessage(&nP, &chPQ, Message{0, 1, 300}) // P sends message to Q
    sendMessage(&nP, &chPR, Message{0, 2, 300}) // P sends message to R
    time.Sleep(2 * time.Second)
    
    processPassive(&nP) // P becomes passive
    processPassive(&nR) // R becomes passive
    time.Sleep(2 * time.Second)
    
    sendMessage(&nQ, &chQR, Message{1, 2, 300}) // Q sends message to R
    time.Sleep(2 * time.Second)
    processPassive(&nQ) // Q becomes passive
    time.Sleep(2 * time.Second)
    processPassive(&nR) // R becomes passive
    time.Sleep(2 * time.Second)
    processPassive(&nQ) // Q stays passive
    time.Sleep(2 * time.Second)
    processPassive(&nP) // P stays passive & exits
        
}