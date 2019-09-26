import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
  vus: 100,
  duration: "60s"
};
export default function() {

const counter = 10000;

for (var id = 1; id <= counter; id++) {

      var url = `http://172.31.10.75:8080/api/${__VU}.${__ITER}.`+id; 
      var payload = JSON.stringify({ user: __VU, iteration: __ITER, counter: id });
      var params =  { headers: { "Content-Type": "application/json" } }
      let res =  http.post(url, payload, params);
      check(res, {
        "status was 200": (r) => r.status == 200
      });
}

for (var id = 1; id <= counter; id++) {
      var url =  `http://172.31.10.75:8080/api/${__VU}.${__ITER}.`+id; 
      let res =  http.get(url);
      
      if (res.status == 200)
      {
          try{
            var response = JSON.parse(res.body)
             check(res, {
                "status was 200": (r) => r.status == 200,
                "Counter is correct": () => response.counter == id,
                "Counter is correct": () => response.user == __VU,
                "Counter is correct": () => response.iteration == __ITER
              });
    
          }
          catch(e){
              console.log("Invalid json ", res.body);
          }
      }
}



};