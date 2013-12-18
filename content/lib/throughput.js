/*
var margin = {top: 20, right: 20, bottom: 20, left: 40},
	width = 300 - margin.left - margin.right,
	height = 200 - margin.top - margin.bottom;

var x = d3.scale.linear()
	.domain([0, d3.max(data, function(d) { return d[0] })])
	.range([0, width]);
 
var y = d3.scale.linear()
	.domain([0, d3.max(data, function(d) { return d[1] })])
	.range([height, 0]);

var line = d3.svg.line()
	.x(function(d, i) { return x(d[0]); })
	.y(function(d, i) { return y(d[1]); });

//var width = d3.select("#throughput").style("width");
//var height = d3.select("#throughput").style("height");

var svg = d3.select("#throughput").append("svg")
	.attr("width", width + margin.left + margin.right)
	.attr("height", height + margin.top + margin.bottom)
	.append("g")
		.attr("transform", "translate(" + margin.left + "," + margin.top + ")");

svg.append("defs").append("clipPath")
	.attr("id", "clip")
	.append("rect")
		.attr("width", width)
		.attr("height", height);

svg.append("g")
	.attr("class", "x axis")
	.attr("transform", "translate(0," + y(0) + ")")
	.call(d3.svg.axis().scale(x).orient("bottom"));

svg.append("g")
	.attr("class", "y axis")
	.call(d3.svg.axis().scale(y).orient("left"));

var path = svg.append("g")
	.attr("clip-path", "url(#clip)")
	.append("path")
		.datum(data)
			.attr("class", "line")
			.attr("d", line);

function updateStats(json) {
	if (!json.Stats) {
		return
	}
	data.push([json.Stats.DownloadRate, Date.now()]);
	if (data.length > 20) {
		data.pop()
	}
	// redraw the line, and slide it to the left
	path
		.attr("d", line)
		.attr("transform", null)
		.transition()
			.duration(500)
			.ease("linear")
			.attr("transform", "translate(" + x(-1) + ",0)")
			.each("end", updateStats)
		.call(d3.svg.axis().scale(x).orient("bottom"));
	console.log(data)
}
*/
