/*
 * Copyright 2013-2019 The NATS Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

"use strict";

const NATS = require('nats');
const nats = NATS.connect('nats://172.31.46.240:4222');
const fs = require('fs');

///////////////////////////////////////
// Publish Performance
///////////////////////////////////////

const loop = 1*1000000;
const hash = loop*0.01;

console.log('Publish Performance Test');

var i =0;
var start = new Date();
var queries = hash;

function query(specific){
    
        if (queries == hash) start = new Date();
        if (queries-- > 0){
        let index  = Math.floor(Math.random()*loop);
        if (specific) index = specific;     
            
        if (queries % 1000 == 0 )process.stdout.write(".");
        nats.requestOne('state.query', 'state.'+index.toString(), {'max':1}, 3000, function(response) {

            if(response instanceof NATS.NatsError && response.code === NATS.REQ_TIMEOUT) {
                console.log(index, 'not found');

            }
            else if (index != response){
                console.log("Expected", index , "got", index);
            }
            query();
        });

        }
        else {
            const stop = new Date();
            const mps = parseInt(hash / ((stop - start) / 1000), 10);
            console.log('\nQueried at ' + mps + ' queries/sec');    
            console.log("\nFinished Queries successfully!");
        }

}

function publish(){

    if (i <= (loop - hash)){
        for (let j = 0; j< hash; j++) {
            nats.publish('state.'+i, i.toString());
            i++;
        }
        process.stdout.write('+');
        nats.flush(function() {
            setTimeout(publish, 0);
        });
    }
    else{
        nats.flush(function() {
            const stop = new Date();
            const mps = parseInt(loop / ((stop - start) / 1000), 10);
            console.log('\nPublished at ' + mps + ' msgs/sec');
            setTimeout(query, 10000);
        });    
        
    }
}

nats.on('connect', function() {
    publish();
    //query();
});