
function handlePieceMessage(pieceObj) {
	if (pieceObj.TotalPieces) {
		//console.log("The Piece object contains TotalPieces");
		initPiecesGrid(pieceObj.TotalPieces)
	} else if (pieceObj.FinishedPieces) {
		//console.log("The Piece object contains FinishedPieces");
		var finishedPieces = []
		for (var i = 0; i < pieceObj.FinishedPieces.length; i++) {
			//console.log("Marking " + pieceObj.FinishedPieces[i].PieceNum);
			//markFinishedPieces(pieceObj.FinishedPieces[i].PieceNum);
			finishedPieces.push(pieceObj.FinishedPieces[i].PieceNum);
		}
		markFinishedPieces(finishedPieces);
		handleFinishedPieces(finishedPieces);
	} else {
		console.log(pieceObj)
		throw "ERROR: Unknown pieceMessage type"
	}
}

function initPiecesGrid(numPieces) {
  var piecesDiv = d3.select("#pieces")
  var height = piecesDiv[0][0].clientHeight
  console.log("pieces grid heigh: " + height)
  var width = piecesDiv[0][0].clientWidth
  console.log("pieces grid width: " + width)
  var gridArray = createGridArray(width, height, numPieces);
  var sideLength = squareSideLength(width, height, numPieces);
  console.log(gridArray);

  var grid = piecesDiv.append("svg")
                .attr("width", width)
                .attr("height", height)
                .attr("class", "pieces");

  var row = grid.selectAll(".row")
              .data(gridArray)
            .enter().append("svg:g")
              .attr("class", "row");

  var col = row.selectAll(".piece")
               .data(function (d) { return d; })
              .enter().append("svg:rect")
               .attr("class", "piece")
               .attr("x", function(d) { return d.x; })
               .attr("y", function(d) { return d.y; })
               .attr("width", function(d) { return d.width; })
               .attr("height", function(d) { return d.height; })
               .style("fill", function(d) { return d.fill; })
               //.on('click', function() {
               //   console.log(d3.select(this));
               //   d3.select(this)
               //       .style('fill', '#0F0')
               //})
               //.on('dblclick', function() {
               //   d3.select(this)
               //       .style('fill', '#FFF')
               //})
               .style("stroke", '#FFF')
               // Make the border around the squares proportional to
               // their side. 
               .style("stroke-width", Math.floor(sideLength / 4));
} 


function squareSideLength(canvasWidth, canvasHeight, numPieces) {
  var canvasArea = canvasWidth * canvasHeight;
  //console.log("canvasArea: " + canvasArea);
  var squareMaxArea = canvasArea / numPieces;
  //console.log("squareMaxArea: " + squareMaxArea);
  // compute the max side length (an integer)
  var sideLength = Math.floor(Math.sqrt(squareMaxArea));
  //console.log("Initial side length: " + sideLength);

  var numRows
  var numCols
  var numSpaces
  while(true) {
    numRows = Math.floor(canvasHeight / sideLength);
    numCols = Math.floor(canvasWidth / sideLength);
    numSpaces = numRows * numCols;

    //console.log("Checking with side length: " + sideLength);
    //console.log("numRows: " + numRows + " numCols:" + numCols + " numSpaces:" + numSpaces);

    if (numSpaces >= numPieces) {
      // We've found a side length that would allow us to fit every square into the 
      // canvas
      return sideLength;
    } else {
      sideLength -= 1;
    }
  }
}



function createGridArray(width, height, numPieces) {
  var result = new Array();
  var gridSquareLength = squareSideLength(width, height, numPieces);
  var numRows = Math.floor(height / gridSquareLength);
  var numCols = Math.floor(width / gridSquareLength);
  var xpos = 0;
  var ypos = 0;
  var squareNum = 0;
  var fill;

  var rowEle;

  for (var row = 0; row < numRows; row++) {
    result.push(new Array());

    for (var col = 0; col < numCols; col++) {
      if (squareNum < numPieces) {
        // This square corresponds to an actual piece.
        fill = "lightgrey";
      } else {
        // this square does not correspond to a piece. It should be invisible. 
        fill = "#FFF";
      }

      result[row].push({ 
                    width: gridSquareLength,
                    height: gridSquareLength,
                    x: xpos,
                    y: ypos,
                    squareNum: squareNum,
                    fill: fill
                });

      squareNum++;
      xpos += gridSquareLength;
    }

    xpos = 0;
    ypos += gridSquareLength;
  }

  return result;
}


function markFinishedPieces(pieceNums) {
  var targetPieces = [];
  var svg = d3.select(".pieces");
  var canvasWidth = svg[0][0].scrollWidth;
  var canvasHeight = svg[0][0].scrollHeight;
  //console.log("height: " + canvasWidth);
  //console.log("width: " + canvasHeight);

  var gridSquareLength = svg.select(".piece")[0][0].width.animVal.value;
  //console.log("square length: " + gridSquareLength);

  var numCols = Math.floor(canvasWidth / gridSquareLength);
  //console.log("NUM COLS: " + numCols);

  for (var i = 0; i < pieceNums.length; i++) {
  	  var pieceNum = pieceNums[i];
	  var rowNum = Math.floor(pieceNum / numCols);
	  //console.log("ROW NUM: " + rowNum);
	  var colNum = pieceNum % numCols;
	  //console.log("COL NUM: " + colNum);

	  var xpos = colNum * gridSquareLength;
	  var ypos = rowNum * gridSquareLength;

	  //console.log("Filling in color for piece " + pieceNum + " at x:" + xpos + " y:" + ypos);

	  //console.log("ALL PIECES");
	  //console.log(svg.selectAll(".piece"));

	  var targetPiece = {
	    x: xpos,
	    y: ypos
	  };

	  targetPieces.push(targetPiece)
  }

  svg.selectAll(".piece")
      .data(targetPieces, function(d) {
        return "" + d.x + "-" + d.y;
      })
      .style("fill", "forestgreen");



}