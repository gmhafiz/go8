// Install k6 at https://k6.io/
// Run with:
//      k6 run k6.js

import http from 'k6/http';

export const options = {
   stages: [
       { duration: '1m', target: 2 }, // traffic ramp-up from 1 to N users over X minutes.
       { duration: '1m', target: 4 }, // stay at N users for X minute
       { duration: '1m', target: 1 }, // ramp-down to N users
   ],
};

export default function () {
    http.get('http://localhost:3080/api/v1/author');
}
/*
rate(node_cpu_seconds_total{mode="system"}[1m])
node_filesystem_avail_bytes
rate(node_network_receive_bytes_total[1m])

(sum by(instance) (irate(node_cpu_seconds_total{instance="$node",job="$job", mode!="idle"}[$__rate_interval])) / on(instance) group_left sum by (instance)((irate(node_cpu_seconds_total{instance="$node",job="$job"}[$__rate_interval])))) * 100
 */