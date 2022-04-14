%%writefile tarry.go

package main

import (
    "time"
    "os"
)

const N = 3

type Token struct {
    from int
    to int
    visited []*(chan Token)
    spanning [N][N]int
}

type Node struct {
    id int
    parent int // parent = node.id; -1 if unvisited
    outCh map[int]*(chan Token) // key represents the node.id to which Token is being sent
    inCh map[int]*(chan Token) // key represents the node.id from where Token is coming
}

// receiveToken on incoming channel ch of a node 
func receiveToken(node *Node, ch *(chan Token)) {
    select {
    case token := <- *ch:
        if node.parent == -1 {
            node.parent = token.from
            println("Token visited node", node.id, "for the first time")
            token.spanning[token.from][token.to] = 1
            println("Spanning tree updated")
            printGraph(token.spanning)
        } 
        token.visited = append(token.visited, ch)
        println("Channel added to visited array")
        
        flag := 0 // number of availabe channels to send except parent
        for k, val := range node.outCh {
            // id should not be of parent & ch should not be visited
            if (k != token.from) && !find(val, token.visited) {
                flag++
            }
        }
                
        // Forwarding token to a neighbour
        token.from = node.id // new parent is the current node
        
        if flag == 0 { // send back to parent
            for _, c := range node.outCh {
                if find(c, token.visited) {
                println("All channels visited! Final spanning tree: ")
                printGraph(token.spanning)
                os.Exit(1)                    
                }
            } 
            token.to = node.parent
            println("Sending token back to parent", token.to)
        } else {
            for key, _ := range node.outCh {
                if key != node.parent {
                    token.to = key
                    break
                }
            }
            println("Sending token back to neighbour", token.to)   
        }

        println("Token was received by: ", token.to, " from ", token.from)
        
        *node.outCh[token.to] <- token
    }

}

func Receive(node *Node) {
    for _, ch := range node.inCh {
        go receiveToken(node, ch)
    }
}

func find(ch *(chan Token), arr []*(chan Token)) bool {
    for i := 0; i < len(arr); i++ {
        if ch == arr[i] {
            return true
        }
    }
    return false
}

func printGraph(g [N][N]int) {
    for i := 0; i < N; i++ {
        for j := 0; j < N; j++ {
            print(g[i][j], " ")
        }
        println()
    }
}

func main(){
    // initialize nodes & channels
    chPQ, chPR := make(chan Token), make(chan Token)
    chQP, chQR := make(chan Token), make(chan Token)
    chRP, chRQ := make(chan Token), make(chan Token)    
        
    nP := Node{0, 0, 
               map[int]*chan Token{1: &chPQ, 2: &chPR}, 
               map[int]*chan Token{1: &chQP, 2: &chRP}}
    nQ := Node{1, -1, 
               map[int]*chan Token{0: &chQP, 2: &chQR}, 
               map[int]*chan Token{0: &chPQ, 2: &chRQ}}
    nR := Node{2, -1, 
               map[int]*chan Token{0: &chRP, 1: &chRQ}, 
               map[int]*chan Token{0: &chPR, 1: &chQR}}

    go Receive(&nP)
    go Receive(&nQ)
    go Receive(&nR)
    
    token := Token{0, 1, []*(chan Token){}, [N][N]int{}}
    chPQ <- token
    
    time.Sleep(3 * time.Second)
}