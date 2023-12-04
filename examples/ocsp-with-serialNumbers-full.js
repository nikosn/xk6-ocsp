import http from 'k6/http';
import {
    check,
    fail,
    sleep
} from 'k6';
import { SharedArray } from 'k6/data';
import {
    Counter,
    Rate
} from 'k6/metrics';
import { describe } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import ocspmodule from 'k6/x/ocsp';

let ErrorCount = new Counter("errors");
let ErrorRate = new Rate("error_rate");

export let options = {
    scenarios: {
        constantocsp: {
            executor: 'constant-arrival-rate',

            // 50 iterations per second
            rate: 50,
            timeUnit: '1s',

            // for how long do you want the scenario to run
            duration: '5m',

            // this number doesn't really matter, as long as it's high enough
            // that there is always a free VU to run an iteration on
            preAllocatedVUs: 250,
        },
    },
    summaryTrendStats: ["min", "med", "avg", "max", "p(90)", "p(95)", "p(99)"],
    thresholds: {
        'fatal_errors': [{
            threshold: 'count<1',
            abortOnFail: true
        }]
    }
};
const abortMetric = new Counter('fatal_errors');

// load the hex serialNumbers from file
const serialNumbers = new SharedArray('serialNumbers', function () {
    const f = open(`${__ENV.HEX_SERIALNUMBERS_FILE_PATH}`).split(/\r?\n/);
    f.pop();
    return f;
});

export default function () {
    const endpointURL = `${__ENV.ENDPOINT_URL}`;
    const issuerCertPath = `${__ENV.ISSUER_CERT_PATH}`;
    const hashAlgorithm = `${__ENV.HASH_ALGORITHM}`;

    check(endpointURL, {
        ['OCSP responder: ' + endpointURL]: (endpointURL) => endpointURL != null
    });
    check(issuerCertPath, {
        ['issuer certificate: ' + issuerCertPath]: (issuerCertPath) => issuerCertPath != null
    });
    check(hashAlgorithm, {
        ['hashAlgorithm: ' + hashAlgorithm]: (hashAlgorithm) => hashAlgorithm != null
    });
    const serialNumber = randomItem(serialNumbers);
    console.log(serialNumber);

    describe('Creating OCSP request and calling OCSP responder', () => {
        if (endpointURL === "undefined" || endpointURL === null || endpointURL === "") {
            abortMetric.add(1);
            sleep(1);
            fail(`ENDPOINT_URL has to be specified. (Currently set to '${endpointURL}')`)
        }
        const ocspRequest = ocspmodule.createRequest(serialNumber, issuerCertPath, hashAlgorithm);

        const params = {
            headers: {
                'Content-Type': 'application/ocsp-request',
                'Accept': 'application/ocsp-response'
            },
            responseType: "binary"
        };
        const response = http.post(endpointURL, ocspRequest, params);

        const success = check(response, {
            'response time <= 1000ms': (r) => r.timings.duration <= 1000,
            'status code = http 200': (r) => r.status === 200,
        });
        if (!success) {
            ErrorCount.add(1);
            ErrorRate.add(true);
        }
        check(response, {
            'Content-Type is application/ocsp-response': (response) => response.headers['Content-Type'] === 'application/ocsp-response',
        });

        describe('Checking OCSP response', () => {
            const ocspresponse = ocspmodule.checkResponse(response.body, true);
            console.log(ocspresponse);
            const validOcsp = check(ocspresponse, {
                'OCSP status is Good, Revoked or Unknown': (value) => value === 'Good' || value === 'Revoked' || value === 'Unknown',
            });
            if (!validOcsp) {
                ErrorCount.add(1);
                ErrorRate.add(true);
            }
        })
    })
}
