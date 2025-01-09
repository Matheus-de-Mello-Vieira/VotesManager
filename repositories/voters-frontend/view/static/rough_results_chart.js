// results.js: Creates a bar chart using Chart.js

// 1. Read the JSON data from the hidden <script> tag
const dataEl = document.getElementById("chartData");
const jsonData = JSON.parse(dataEl.textContent);

// Extract labels and votes
const labels = jsonData.labels;
const votes = jsonData.votes;

// 2. Prepare the Chart.js data object
const data = {
  labels: labels,
  datasets: [
    {
      label: "Número de Votos",
      data: votes,
      backgroundColor: ["rgba(54, 162, 235, 0.2)"],
      borderColor: ["rgba(54, 162, 235, 1)"],
      borderWidth: 1,
    },
  ],
};

// 3. Chart config
const config = {
  type: "bar",
  data: data,
  options: {
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  },
};

// 4. Create and render the chart
const ctx = document.getElementById("resultChart").getContext("2d");
const resultChart = new Chart(ctx, config);
