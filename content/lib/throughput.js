		// create an SVG element inside the #graph div that fills 100% of the div
		var graph = d3.select("#throughput").append("svg:svg").attr("width", "100%").attr("height", "100%");

		// create a simple data array that we'll plot with a line (this array represents only the Y values, X will just be the index location)
		var download = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
		var upload = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];

		// X scale will fit values from 0-10 within pixels 0-100
		var x = d3.scale.linear().domain([0, 48]).range([-5, Math.floor(d3.select("#throughput").style('width').replace(/[^\d.-]/g,''))]); // starting point is -5 so the first value doesn't show and slides off the edge as part of the transition
		// Y scale will fit values from 0-10 within pixels 0-100
		var y = d3.scale.linear().domain([0, 4000]).range([200, 0]);

		// create a line object that represents the SVG line we're creating
		var line = d3.svg.line()
			.x(function(d,i) { return x(i); })
			.y(function(d) { return y(d); })
			.interpolate("basis")
	
			// display the line by appending an svg:path element with the data line we created above
			graph.append("svg:path").attr("d", line(download));
			graph.append("svg:path").attr("d", line(upload));
			// or it can be done like this
			//graph.selectAll("path").data([data]).enter().append("svg:path").attr("d", line);
			
			
			function redrawDownload() {
				// update with animation
				graph.selectAll("path")
					.data([download]) // set the new data
					.attr("transform", "translate(" + x(1) + ")") // set the transform to the right by x(1) pixels (6 for the scale we've set) to hide the new value
					.attr("d", line) // apply the new data values ... but the new value is hidden at this point off the right of the canvas
					.transition() // start a transition to bring the new value into view
					.ease("linear")
					.duration(1000) // for this demo we want a continual slide so set this to the same as the setInterval amount below
					.attr("transform", "translate(" + x(0) + ")"); // animate a slide to the left back to x(0) pixels to reveal the new value
			}
			function redrawUpload() {
				graph.selectAll("path")
					.data([download]) // set the new data
					.attr("transform", "translate(" + x(1) + ")") // set the transform to the right by x(1) pixels (6 for the scale we've set) to hide the new value
					.attr("d", line) // apply the new data values ... but the new value is hidden at this point off the right of the canvas
					.transition() // start a transition to bring the new value into view
					.ease("linear")
					.duration(1000) // for this demo we want a continual slide so set this to the same as the setInterval amount below
					.attr("transform", "translate(" + x(0) + ")"); // animate a slide to the left back to x(0) pixels to reveal the new value
			}
			
			function updateStats(stats) {
			   download.shift();
			   download.push(stats.DownloadRate / 1000);
			   upload.shift();
			   upload.push(stats.UploadRate / 1000);
			   redrawDownload();
			   redrawUpload();
			   console.log(download);
			   console.log(upload);
			}
