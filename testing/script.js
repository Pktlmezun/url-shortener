import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js'; // Library to generate unique IDs

export const options = {
    vus: 500, // Number of concurrent users
    duration: '30s', // Test duration
};

export default function () {
    const url = 'https://url-shortener-zskc.onrender.com/signup'; // Replace with your actual signup endpoint

    // Generate a unique username and email for each request
    const uniqueID = uuidv4(); // Generates a UUID
    const username = `testuser_${uniqueID}`;
    const email = `testuser_${uniqueID}@example.com`;

    const payload = JSON.stringify({
        username: username,
        email: email,
        password: 'testpassword123', // Use the same password or customize as needed
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    http.post(url, payload, params);
    sleep(1); // Wait 1 second between requests
}
