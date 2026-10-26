import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    vus: 100,
    duration: "20s",
};

const DEFAULT_AI_URL = "http://127.0.0.1:8000";
const DEFAULT_AI_KEY = "1yOojwxW0POZMQxQ";

const AI_URL = __ENV.AI_URL || DEFAULT_AI_URL;
const AI_API_KEY = __ENV.AI_API_KEY || DEFAULT_AI_KEY;

const CONCERNS = [
    "tulong naman po, emergency suicidal thoughts",
    "how to get student id and campus map saan po",
    "delayed po ang graduation ko may grades issue",
    "ordinary routine check for library hours please",
    "sobrang stress at depress na ako sa grades ko par",
    "di ko na kaya ang pressure sa thesis namin",
    "inquire lang po regarding shifting procedure",
];

export default function () {
    const randomIndex = Math.floor(Math.random() * CONCERNS.length);
    const randomText = CONCERNS[randomIndex];
    const payload = JSON.stringify({ text: randomText });

    const params = {
        headers: {
            "Content-Type": "application/json",
            "X-API-Key": AI_API_KEY,
        },
    };

    const res = http.post(`${AI_URL}/api/v1/classify`, payload, params);

    check(res, {
        "AI: status is 200": (r) => r.status === 200,
        "AI: has level": (r) => {
            try {
                const body = r.json();
                return body && body.level !== undefined;
            } catch (e) {
                return false;
            }
        },
    });

    sleep(0.5);
}
