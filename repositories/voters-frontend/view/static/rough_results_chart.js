const dataEl = document.getElementById("chartData");
const jsonData = JSON.parse(dataEl.textContent);

const labels = jsonData.labels;
const votes = jsonData.votes;

const total = votes.reduce((a, b) => a + b, 0)
const percentages = votes.map((a) => a / total)

const data = {
  labels: labels,
  datasets: [
    {
      label: "Porcentagem dos Votos",
      data: percentages,
      backgroundColor: ["rgba(54, 162, 235, 0.2)"],
      borderColor: ["rgba(54, 162, 235, 1)"],
      borderWidth: 1,
    },
  ],
};

const config = {
  type: "bar",
  data: data,
  options: {
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          format: {
            style: "percent",
          },
        },
      },
    },
  },
};

// 4. Create and render the chart
const ctx = document.getElementById("resultChart").getContext("2d");
const resultChart = new Chart(ctx, config);
