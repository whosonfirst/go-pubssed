window.addEventListener("load", function load(event){

    var root_url = location.protocol + "//" + location.host;
    var sse_url = root_url + "/sse";

    var ev = new EventSource(sse_url);

    ev.onopen = function(e){
	console.log("SSE connected");
    };

    ev.onerror = function(e){
	console.log("SSE error", e);
    };
    
    ev.onmessage = function(e) {
	var el = document.getElementById("message");
	el.innerText = e.data;
    }

});
