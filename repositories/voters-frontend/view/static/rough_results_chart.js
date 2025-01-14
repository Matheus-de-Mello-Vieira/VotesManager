const API_URL = "http://localhost:8080";

document.addEventListener("DOMContentLoaded", async () => {
  const totalsByParticipant = await getTotalByParticipant();
  displayChartTotalsByParticipant(totalsByParticipant);
});

async function getTotalByParticipant(retries = 3) {
  try {
    const response = await fetch(`${API_URL}/votes/totals/rough`);

    if (!response.ok) {
      throw new Error("Falha ao obter resultados");
    }
    return await response.json();
  } catch (error) {
    console.error("Erro ao obter resultados:", error);
    if (retries > 0) {
      console.log(`Tentando novamente. Tentativas restantes: ${retries - 1}`);
      setTimeout(() => getResults(retries - 1), 2000);
    } else {
      alert(
        "Ocorreu um erro ao tentar obter os resultados. Por favor, tente novamente mais tarde."
      );
    }
  }
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
}
