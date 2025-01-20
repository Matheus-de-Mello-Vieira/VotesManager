document.addEventListener("DOMContentLoaded", async () => {
  const totalsByParticipant = await getTotalByParticipant();
  displayChartTotalsByParticipant(totalsByParticipant);
});

async function getTotalByParticipant(retries = 3) {
  return await makeRequest(`/api/votes/totals/rough`, retries);
}

function displayChartTotalsByParticipant(totalByParticipant) {
  const labels = Object.keys(totalByParticipant);
  const totals = Object.values(totalByParticipant);

  const subtotal = totals.reduce((a, b) => a + b, 0);
  const percentages = totals.map((a) => a / subtotal);

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
          type: "linear",
          min: 0,
          ticks: {
            format: {
              style: "percent",
            },
          },
        },
        x: {
          type: "category",
        },
      },
    },
  };

  // 4. Create and render the chart
  const ctx = document.getElementById("resultChart").getContext("2d");
  const resultChart = new Chart(ctx, config);
}
