package main

type transaction struct {
	ProtocolVersion int                 `json:"protocolVersion"`
	Inputs          []transactionInput  `json:"inputs"`
	Outputs         []transactionOutput `json:"outputs"`
}

type transactionInput struct {
	PreviousTransactionHash  string `json:"previousTransactionHash"`
	PreviousTransactionIndex int    `json:"previousTransactionIndex"`
	SenderSignature          string `json:"senderSignature"`
}

type transactionOutput struct {
	ReceiverAddress string  `json:"receiverAddress"`
	Amount          float64 `json:"amount"`
}
