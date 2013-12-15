var socket = new WebSocket("ws://localhost:8080/ws", "protocolOne");
socket.onmessage = function(msg) {
	console.log(msg.data);
};

var w = 800;
var h = 600;
var force = d3.layout.force()
	.nodes(dataset.nodes)
	.links(dataset.edges)
	.size([w, h])
	.start();

var svg = d3.select("body").append("svg")
	.attr("width", w)
    .attr("height", h);


/*
var svg = d3.select("body")
			.append("svg")
			.attr("width", 600)
			.attr("height", 400);

var n = 100
var max = Math.ceil(Math.sqrt(n))
for (var row = 0; row < max; row++) {
	for (var col = 0; col < max; col++) {
		var rect = svg.append("rect")
			.attr("x", col * 20 + 5)
			.attr("y", row * 20 + 5)
			.attr("height", 10)
			.attr("width", 10)
			.attr("fill", "black");
	}
}
*/
