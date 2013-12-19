
var files = [];
var filesToNum = {};

function handleProgressMessage(progressObj) {
	if (progressObj.AddFiles) {
		for (var i = 0; i < progressObj.AddFiles.length; i++) {
			progressObj.AddFiles[i].row = i;
			progressObj.AddFiles[i].totalPieces = progressObj.AddFiles[i].LastPiece - progressObj.AddFiles[i].FirstPiece + 1;
			progressObj.AddFiles[i].finishedPieces = 0;
		}
		addFiles(progressObj.AddFiles)
	} else {
		console.log(progressObj)
		throw "ERROR: Unknown pieceMessage type"
	}
}




function addFiles(fileInfo) {
  files = fileInfo
  //console.log("--FILES--");
  //console.log(files)

  var div = d3.select("#progress").selectAll("div")
    .data(files)
    .enter()
    .append("div");

  div.append("span")
      //.style("width", "74px")
      .style("width", "300px")
      .style("text-align", "left")
      .style("font", "8px sans-serif")
      .style("vertical-align", "middle")
      .style("line-height", "20px")
      .text(function(d) {
      	return d.FileName
      });

  div.append("span")
      .attr("class", "progress-outline")
      .style("width", "240px")
      .style("left", "120px")
    .append("span")
      .attr("class", "value")
      .attr("row", function(d) {
      	d.row
      })
      .style("text-align", "center")
      .text("0%")
      .style("width", "0%");

  div.transition()
      .style("height", "20px")
}

function handleFinishedPieces(finishedPieces) {
  for (var fileIndex = 0; fileIndex < files.length; fileIndex++) {
  	var file = files[fileIndex];
  	for (var pieceIndex = 0; pieceIndex < finishedPieces.length; pieceIndex++) {
	  if (finishedPieces[pieceIndex] >= file.FirstPiece && finishedPieces[pieceIndex] <= file.LastPiece) {
	  	// This file contains this piece
	  	file.finishedPieces++;
	  }
  	}
  }
  updateProgress();
}

function updateProgress() {
  //console.log("Updating progress for files.");

  d3.select("#progress").selectAll(".value")
      .data(files, function(d) {
        return d.row;
      })
      .text(function(d) {
      	//console.log("UPDATING TEXT FOR BAR");
      	//console.log(d);
      	//console.log("RETURN: " + ((d.finishedPieces / d.totalPieces) * 100).toFixed(2) + "%");
      	return ((d.finishedPieces / d.totalPieces) * 100).toFixed(2) + "%";
      })
      .style("width", function(d) {
      	//console.log("UPDATING STYLE FOR BAR")
      	//console.log(d)
      	//console.log("RETURN: " + ((d.finishedPieces / d.totalPieces) * 100) + "%");
      	return ((d.finishedPieces / d.totalPieces) * 100) + "%";
      });
      //.text("" + (percentage * 100).toFixed(2) + "%")
      //.style("width", (percentage * 100) + "%");

}
