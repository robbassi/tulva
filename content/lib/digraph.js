var force = d3.layout.force();
var nodeMap = {};
var linkMap = {};
var smallCircleRadius = 7;
var largeCircleRadius = 16;

function initDiGraph() {
    var digraph = d3.select(".digraph")
    //console.log(digraph)
    var height = digraph[0][0].clientHeight
    //console.log("digraph heigh: " + height)
    var width = digraph[0][0].clientWidth
    //console.log("digraph width: " + width)

    digraph.append("svg")
    	.attr("width", width)
    	.attr("height", height);

    // build the arrow.
	digraph.select("svg").append("svg:defs").selectAll("marker")
	    .data(["end"])      // Different link/path types can be defined here
	  .enter().append("svg:marker")    // This section adds in the arrows
	    .attr("id", String)
	    .attr("viewBox", "0 -5 10 10")
	    .attr("refX", 15)
	    .attr("refY", -1.5)
	    .attr("markerWidth", 6)
	    .attr("markerHeight", 6)
	    .attr("orient", "auto")
	  .append("svg:path")
	    .attr("d", "M0,-5L10,0L0,5");


	// Create a group element that will hold path elements (edges)
	digraph.select("svg").append("g")


	force
	    // The total width/height of the SVG 'canvas'
	    .size([width, height])

	    // The length of the links that connect two nodes
	    .linkDistance(60)

	    // the amount of attraction or repulsion between nodes. A positive value 
	    // results in attraction, while a negative value results in repulsion. As
	    // mentioned in the documentation, "For graph layout, negative values should be used"
	    .charge(-300)

	    // As the graph is 'converging', this function seems to be called roughly 30 times 
	    // per second. 
	    .on("tick", digraphTick)

	    // Start the simulation
	    .start();


}


function handleDiGraphMessage(graphMessage) {
	if (graphMessage.AddNodes) {
		for (var i = 0; i < graphMessage.AddNodes.length; i++) {
			addNode(graphMessage.AddNodes[i].NodeName, graphMessage.AddNodes[i].NodeID)
		}
	} else if (graphMessage.RemoveNodes) {
		for (var i = 0; i < graphMessage.RemoveNodes.length; i++) {
			removeNode(graphMessage.RemoveNodes[i])
		}
	} else if (graphMessage.AddEdges) {
		for (var i = 0; i < graphMessage.AddEdges.length; i++) {
			//console.log("SourceNode: " + graphMessage.AddEdges[i].SourceNode);
			//console.log("TargetNode: " + graphMessage.AddEdges[i].TargetNode);
			//console.log("Name: " + graphMessage.AddEdges[i].Name);
			//console.log("Intensity: " + graphMessage.AddEdges[i].Intensity);

			addLink(
				graphMessage.AddEdges[i].SourceNode,
				graphMessage.AddEdges[i].TargetNode,
				graphMessage.AddEdges[i].Name,
				(graphMessage.AddEdges[i].Intensity / 100));
		}

	} else if (graphMessage.UpdateEdges) {

	} else if (graphMessage.RemoveEdges) {

	} else {
		throw "ERROR: Unknown graphMessage type"
	}
}


function digraphTick() {
    //console.log("TICK!");
    var svg = d3.select(".digraph").select("svg");
    var path = svg.select("g").selectAll("path");
    var node = svg.selectAll(".node");
    path.attr("d", function(d) {
        var dx = d.target.x - d.source.x,
            dy = d.target.y - d.source.y,
            dr = Math.sqrt(dx * dx + dy * dy);
        return "M" + 
            d.source.x + "," + 
            d.source.y + "A" + 
            dr + "," + dr + " 0 0,1 " + 
            d.target.x + "," + 
            d.target.y;
    });

    node
        .attr("transform", function(d) { 
        return "translate(" + d.x + "," + d.y + ")"; });
}


// action to take on mouse click
function digraphClick() {
    d3.select(this).select("text").transition()
        .duration(750)
        .attr("x", 0)
        .style("fill", "steelblue")
        .style("stroke", "lightsteelblue")
        .style("stroke-width", ".5px")
        .style("font", "20px sans-serif");
    d3.select(this).select("circle").transition()
        .duration(750)
        .attr("r", largeCircleRadius)
        .style("fill", "lightsteelblue");
}

// action to take on mouse double click
function digraphDoubleClick() {
    d3.select(this).select("circle").transition()
        .duration(750)
        .attr("r", 6)
        .style("fill", "#ccc");
    d3.select(this).select("text").transition()
        .duration(750)
        //.attr("x", 12)
        .style("stroke", "none")
        .style("fill", "black")
        .style("stroke", "none")
        .style("font", "10px sans-serif");
}

function getNode(nodeName) {
  //console.log("nodeName: " + nodeName)
  //console.log("NODEMAP")
  //console.log(nodeMap)

  var node = nodeMap[nodeName];
  if (!node) {
  	throw "Tried removing node with name " + node.name + ", but it doesn't exist in the graph.";
  }
  return node
}

function addNode(nodeName, nodeID) {

  var node = { name: nodeName };

  if (force.nodes().indexOf(node) >= 0) {
    // This node already exists in the graph. 
    throw "Tried adding duplicate node with name " + node.name;
  } else {
    console.log("Adding node. name:" + node.name);
  }

  force.nodes().push(node);
  nodeMap[node.name] = node;

  var svg = d3.select(".digraph").select("svg");
  // define the nodes
  var nodeElement = svg
      .selectAll(".node")
      .data(force.nodes())
      .enter()
      .append("g")
      .attr("class", "node")
      .on("click", digraphClick)
      .on("dblclick", digraphDoubleClick)
      .call(force.drag);

  // add the nodes
  nodeElement.append("circle")
      .attr("r", smallCircleRadius);

  // add the text 
  
  
  nodeElement.append("text")
      .attr("dx", -4)
      .attr("dy", 3)
      //.attr("x", 12)
      //.attr("dy", ".35em")
      .text(nodeID);

  nodeElement.append("title")
      .text(node.name);

  force.start();
}

function removeNode(nodeName) {
  //console.log("REMOVE NODE: " + nodeName);
  var node = getNode(nodeName);
  //console.log(node)

  // If this node has any links, remove them first. 
  removeAttachedLinks(node);

  var indexToRemove = force.nodes().indexOf(node);

  if (indexToRemove === -1) {
    throw "Tried removing node with name " + node.name + ", but it doesn't exist in the graph.";
  } else {
    console.log("Removing node. name:" + node.name);
  }
  
  force.nodes().splice(indexToRemove, 1);


  var svg = d3.select(".digraph").select("svg");

  var join = svg.selectAll(".node")
      .data(force.nodes(), function(d) {
        return d.name
      });

  join.exit().remove();

  force.start();
}


function addLink(sourceName, targetName, name, intensity) {

  var source = getNode(sourceName);
  var target = getNode(targetName);

  if (!source) {
  	throw "node with name " + sourceName + " not found in graph";
  } else if (!target) {
  	throw "node with name " + targetName + " not found in graph";
  }

  var linkId = source.name + "-" + target.name + "-" + name;
  //console.log("LINK ID: " + linkId);
  var link = {
  	source: source,
  	target: target,
  	name: name,
  	opacity: intensity
  };
  linkMap[linkId] = link

  var existingLink
  for (var i = 0; i < force.links().length; i++) {
    existingLink = force.links()[i]
    if (existingLink.source === link.source &&
      existingLink.target === link.target &&
      existingLink.name === link.name) {
      throw "Tried adding duplicate link with source:" + link.source.name + " target:" + link.target.name + " name:" + link.name + " opacity:" + link.opacity;
    }
  }

  console.log("Adding link. source:" + link.source.name + " target:" + link.target.name + " name:" + link.name + " opacity:" + link.opacity);

  force.links().push(link);

  var svg = d3.select(".digraph").select("svg").select("g");

  var path = svg
    .selectAll("path")
    .data(force.links())
    .enter()
    .append("svg:path")
    .style("opacity", link.opacity)
    .attr("class", "link")
    .attr("marker-end", "url(#end)");

  force.start();

}

function updateLinkIntensity(source, target, name, intensity) {

  var existingLink = null
  for (var i = 0; i < force.links().length; i++) {
    if (force.links()[i].source === link.source &&
      force.links()[i].target === link.target &&
      force.links()[i].name === link.name) {
      existingLink = force.links()[i]
      break;
    }
  }

  if (existingLink === null) {
    throw "Unable to update link as it wasn't found. source:" + link.source.name + " target:" + link.target.name + " name:" + link.name;
  } 

  console.log("Updating link opacity from " + existingLink.opacity + " to " + link.opacity + " source:" + link.source.name + " target:" + link.target.name + " name:" + link.name);
  
  existingLink.opacity = link.opacity;

  var svg = d3.select(".digraph").select("svg").select("g");

  var path = svg
    .selectAll("path")
    .data([existingLink], function(d) {
      return d.source.name + d.target.name + d.name;
    })
    .style("opacity", link.opacity);

  force.start();
  

}

function removeAttachedLinks(node) {

  // Remove all links attached to this node from force.links()
  var link
  for (var i = 0; i < force.links().length; i++) {
    link = force.links()[i]
    if ((link.source.name === node.name) ||
      (link.target.name === node.name)) {
        console.log("Removing link. source:" + link.source.name + " target:" + link.target.name + " name:" + link.name + " opacity:" + link.opacity);
        force.links().splice(i, 1);
        i -= 1;
      }
  }


  var svg = d3.select(".digraph").select("svg").select("g");

  var join = svg.selectAll("path")
          .data(force.links(), function(d) {
            return d.source.name + d.target.name + d.name;
          });

  join.exit().remove();

  force.start();

}
