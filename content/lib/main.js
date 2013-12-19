var throughputData = []
var socket = new WebSocket("ws://localhost:8080/ws", "protocolOne");

//initializeDiGraph();
//initThroughputChart();

socket.onopen = function() {
	console.log("SOCKET OPENED");
}

socket.onmessage = function(msg) {
	console.log(msg.data)
	var json = JSON.parse(msg.data);
	if (json.Stats) {
		console.log("Received Stats data");
		updateStats(json.Stats);
	} else if (json.DiGraph) {
		console.log("Received DiGraph data");
		//handleDiGraphMessage(json.DiGraph);
	} else if (json.Progress) {
		console.log("Received Progress data");
		handleProgressMessage(json.Progress);
	} else if (json.Piece) {
		console.log("Received Piece data");
		handlePieceMessage(json.Piece);
	} else {
		console.log("Received JSON of unknown type");
		console.log(json);
	}
};
socket.onclose = function() {
	console.log("SOCKET CLOSED");
};
