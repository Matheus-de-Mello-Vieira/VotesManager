document.addEventListener("DOMContentLoaded", async () => {
  await displayParticipantsButtons();
});

async function displayParticipantsButtons() {
  const buttons = document.getElementById("participants-buttons");
  const participants = await getParticipants();

  participants.forEach((participant) => {
    const result = document.createElement("button");
    result.value = participant.id;
    result.textContent = participant.name;
    result.className = "participant-button";
    result.onclick = () => onVote(participant.id);
    buttons.appendChild(result);
  });
}

async function getParticipants(retries = 3) {
  return await makeRequest(`api/participants`, retries);
}

async function onVote(participantId) {
  try {
    const response = await fetch(`/api/votes`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        participant_id: participantId,
        captcha_token: getCaptchaToken(),
      }),
    });

    if (!response.ok) {
      const content = await response.text();

      if (content.trim() == "Invalid Captcha") {
        alert("Você precisa fazer o CAPTCHA corretamente!");
        return;
      }

      throw new Error("Falha ao registrar o voto");
    }
  } catch (error) {
    console.error("Erro ao votar:", error);
    alert("Ocorreu um erro ao tentar votar. Por favor, tente novamente.");
    return;
  }

  window.location.replace(`after-vote`);
}

function getCaptchaToken() {
  const body = {
    success: checkCaptcha(),
    error_codes: [],
  };

  return btoa(JSON.stringify(body));
}

function checkCaptcha() {
  const captcha = document.getElementById("captcha");
  return captcha.checked;
}
