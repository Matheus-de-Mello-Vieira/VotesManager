import { check } from "k6";
import http from "k6/http";

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: "constant-arrival-rate",
      rate: 1000, // Number of iterations (connections) per second
      timeUnit: "1s", // Per second
      duration: "30s", // Total duration of the test
      preAllocatedVUs: 100, // Initial number of virtual users
      maxVUs: 1000, // Maximum number of virtual users
    },
  },
};

const url = __ENV.URL;

const participantsPDF = [
  { id: 1, probability: 0.5 },
  { id: 2, probability: 0.4 },
  { id: 3, probability: 0.1 },
];

export default function () {
  const tester = new Tester(url, participantsPDF);
  tester.test();
}

class Tester {
  constructor(url, participantsPDF) {
    this.url = url;
    this.participantsCDF = this.calcParticipantsCDF(participantsPDF);
  }

  calcParticipantsCDF(participantsPDF) {
    let cumulativeProbability = 0;
    return participantsPDF.map((participant) => {
      cumulativeProbability += participant.probability;
      return {
        ...participant,
        cumulativeProbability: cumulativeProbability,
      };
    });
  }

  selectParticipant() {
    const r = Math.random();
    for (const participantCDF of this.participantsCDF) {
      if (r <= participantCDF.cumulativeProbability) {
        return participantCDF.id;
      }
    }
    return this.participantsCDF[this.participantsCDF.length - 1].id;
  }

  test() {
    this.testGet("");
    this.testGet("api/participants");

    this.testVote();

    this.testGet("after-vote");
    this.testGet("api/votes/totals/rough");
  }

  testGet(route) {
    const response = http.get(`${this.url}/${route}`);

    check(response, {
      "status is 200": (r) => r.status === 200,
    });
  }

  testVote() {
    const participantId = this.selectParticipant();
    const payload = JSON.stringify({
      participant_id: participantId,
    });

    const headers = { "Content-Type": "application/json" };
    const response = http.post(`${this.url}/api/votes`, payload, {
      headers: headers,
    });

    check(response, {
      "status is 201": (r) => r.status === 201,
    });
  }
}
