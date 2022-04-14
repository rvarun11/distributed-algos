package main

import (
    "time"
    "os"
)

/*
Control message: value = -1, counter & tag = nil
Basic message will have value, tag & counter = nil
*/

type Message struct {
    value int
    tag bool 
    counter int
}

type Node struct {
    id int
    recorded bool
    inCh map[int]*(chan Message)
    outCh map[int]*(chan Message)
    counter map[int]int
    state map[int][]int
}

func receiveMessage(node *Node, ch *(chan Message), idx int) {
    for {
        select {
        case m := <- *ch:
            if m.value != -1 { // receives a basic message
                println("Node", node.id, "received basic message", m.value, "with tag", m.tag)
                if m.tag {
                    TakeSnapshot(node)
                } else if node.recorded {
                    node.state[idx] = append(node.state[idx], m.value)
                    canTerminate(node)
                }
            } else { 
                println("Node", node.id, "received control message")
                node.counter[idx] = m.counter
                TakeSnapshot(node)
                canTerminate(node)
            }
        }
    }
}

func Receive(node *Node) {
    for i, ch := range node.inCh {
        go receiveMessage(node, ch, i)
    }
}

func sendMarker(node *Node, idx int) {
    c := node.counter[idx] + 1
    *node.outCh[idx] <- Message{-1, false, c}
}

func sendBasicMessage(node *Node, val int, idx int) {
    m := Message{val, node.recorded, -1}
    *node.outCh[idx] <- m
    if !node.recorded {
        node.counter[idx]++
    }
}

func canTerminate(node *Node) {
    for i, v := range node.state {
        if len(v) + 1 !=  node.counter[i] {
            return
        }       
    }
    os.Exit(1)   
}

func TakeSnapshot(node *Node) {
    if !node.recorded {
        println("Node", node.id, "took local snapshot")
        node.recorded = true
        for i, _ := range node.outCh {
            go sendMarker(node, i)           
        }
    }
}

func main(){
    // initialize nodes & channels
    chPQ := make(chan Message)
    chQP := make(chan Message)
                
    nP := Node{0, false,
               map[int]*chan Message{1: &chQP}, 
               map[int]*chan Message{1: &chPQ},
               map[int]int{1:0},
               map[int][]int{1:{}},
              }
    nQ := Node{1, false,
               map[int]*chan Message{0: &chPQ}, 
               map[int]*chan Message{0: &chQP},
               map[int]int{0:0},
               map[int][]int{0: {}},
              }
    
    go Receive(&nP); go Receive(&nQ)
    
    sendBasicMessage(&nP, 100, 1)
    sendBasicMessage(&nP, 200, 1)
    time.Sleep(2 * time.Second)
    TakeSnapshot(&nP)
    sendBasicMessage(&nP, 300, 1)
    time.Sleep(2 * time.Second)
    sendBasicMessage(&nP, 400, 1)
    time.Sleep(2 * time.Second)
    
}