const API_URL = "http://localhost:8080";

async function makeRequest(url, retries = 3) {
    try {
      const response = await fetch(url);
  
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