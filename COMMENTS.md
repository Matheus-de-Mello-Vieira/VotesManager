# Paredão BBB

## Tentativas

### Primeira tentativa

Quando recebi o problema, eu pensei o seguinte:

* Como o fluxo de votos é muito grande, e como no final de cada voto eu vou precisa dos resultados parciais, não é viável realizar esse cálculo múltiplas vezes. Além disso, a precisão dos dados não é um fator tão relevante assim.
* Os dados que vão para a produção são muito menos requisitados, de forma que poderia ter um custo computacional maior e se tornar mais preciso.

Ai eu pensei nas seguintes soluções e foi descartando elas:

1) ***\*armazenar o valor da soma total e atualizar toda vez que alguém votar\****: como seria poucos registros (um para cada participante) sendo atualizados pro milhares de agentes ao mesmo tempo, não seria algo viável devido a lentidão por locks.
2) **armazenar o valor da soma total e atualizar periodicamente**: isso separaria a leitura e escrita, de forma a evitar locks entre esses dois, e tornaria o sistema muito mais rápido.

Ai eu pensei em refinamentos

1. o sistema não precisa somar tudo, podemos simplesmente manter no banco de dados o registro da hora da última soma, e apenas completar a soma já existe com os registros feitos apartir dessa data 
   * O problema disso seriam dois:
     * **performance**: toda vez que essa rotina for executada, ela fará um query no banco de dados para filtros todos os votos que aconteceram depois de uma data especifica. Se o SGBD não suporta ter timestamp enquanto índice primário, então esse solução é inviável.
       * considerando que os índices do bancos de dados são feitos através de uma árvore B ou uma variação dessa, normalmente apenas o índice primário é ordenado (chave primária ou qualquer estrutura parecida)
         * o fato da chave de busca (no caso a data) ter um índice ordenado é importante pois apenas nesses casos existe um fila encadeada que bastaria encontrar o primeiro elemento e retornar todos os que o segue. Os bancos de dados fazem isso automaticamente quando existe esse índice. 
     * **consistência**: seria necessário uma transação dentro do banco dados para manter a consistência, que é algo que preferia evitar por problemas de performance
       * aqui são os cenários caso ele não exista:
         * usar a data antes da soma: os votos que acontecerão entre esse meio termo serão somadas duas vezes: uma na soma atual, e outra na próxima, pois eles terão uma data superior que data salva no banco da última soma
         * usar a data depois da soma: os votos que acontecerão entre esse meio termo nunca serão contabilizados, pois foram adicionados depois que a soma foi feita, e na próxima soma também não serão excluíveis por ter uma data inferior a data salva no banco da última soma.
       * se for usar transação, ele precisa ser *serializable*, pois qualquer nível de isolamento inferior causaria inconsistência
         * considerando o nível *repeatable read* que é um nível abaixo do *serializable* por permitir obtenção de resultados parciais de operações feitas por outros agentes ativos no momento que a transação foi iniciada (*Phanton*): se dois votos foram feitos ao mesmo tempo no inicio da transação, um pode entrar na soma e o outro não (que se perderá para sempre pelos motivos descritos acima).
2. A mesma solução da 1, mas ao invés de data, usar uma chave primária serializado, de forma que a rotina somadora precisa apenas armazenar o id do último evento somado
   * Apesar de ter superar a solução 1, ele pressupõe que seria viável um banco de dados estabelecer mais de 1000 conexões por segundos, principalmente porque eu pretendo testar o sistema na minha máquina local, então não teria tantos contêineres.
3. A mesma solução da 1, mas poderíamos pensar em armazenar os dados agrupados por tempo e por candidato, ao invés de apenas armazenar os dados simples.
  * esse agrupamento facilitaria também a parte do código para a produção do BBB, pois teria que somar muito menos dados
  * Se for armazenar dados ao longo prazo, poderia agrupar esses dados em partições cada vez maiores de acordo com a distância dos dados com o presente. Ou seja, quanto mais ao passado iria, menor seria a precisão dos dados, mas menor seria o espaço necessário para o armazenamento.

Dessa forma, eu pensei uma arquitetura os seguintes:

* `PostgreSQL`: 
  * armazenaria os dados brutos, os dados acumulados e os caches.
* `voters-frontend`:
  * interface do usuário comum
    * isso implica uma interface http que responda html
  * envia o registro de votação diretamente ao banco de dados
* `votes-aggregator`:
  * invocado pelo cron
  * periodicamente atualizar a soma que vai para os usuários
* `prodution-frontend`:
  * interface da produção
  * executar as somas refinadas sob demanda da produção

A separação desses componentes foi os motivos para cada um escalar:

* `voters-frontend` precisaria escalar de acordo com o seu consumo de CPU
* `votes-aggregator` precisaria escalar de maneira binária: ligando quando o cron o invoca, e desligando quando termina
* `prodution-frontend` não precisaria de escalonamento, pois é consumido por um número limitado de clientes (a produção).

O problema dessa solução que ela supõe que seria viável um banco de dados PostgreSQL sofrer escrita de milhares de agentes diferentes por segundo, que mesmo sem considerar os locks seria algo impossível.

### Segunda tentativa

Afim de minimizar o número de conexões sofridas pelo banco (por enquanto apenas na parte da escrita), eu pensei em tornar as votações assíncronas: o `voters-frontend` enviaria para uma fila, ai o `votes-register` consumia várias votações e salvaria no banco de dados com uma única conexão.

Eu pensei em calcular os dados acumulados dentro do `voters-frontend`, na seguinte maneira:

1) consumir `N` mensagens ou até um certo timeout

2) somar de acordo com a partição do acumulo e salvar em uma estrutura de dados na memória.

   * inicialmente estava pensando na partição de 30 segundos.

   * a estrutura de dados deveria ser algo hibrido entre:
     * um hash, devido o seu tempo esperado médio na busca ser $O(1)$ 
     * uma lista circular: como não faz sentido armazenar dados muitos antigos na memória, uma lista circular eliminaria os dados mais antigos de maneira automática.

   * Isso faz com que precise com que a memória interna do contêiner seja autoritativo a respeito das partições com as quais ele trabalhe. Ou seja, todas as mensagens de uma certa partição ir para apenas uma máquina especifica.

3) Salvar no banco de dados em uma transação só, pois assim seria aberta apenas uma conexão com o banco de dados.

4) Volte a etapa 1

Dessa forma, como eu precisava que o contêiner seja autoritativo a respeito das partições com as quais ele trabalha, eu pensei em usar o Kafka, pois:

* O Kafka particiona as mensagens de acordo com o atributo `partition key`, mensagens com esse atributo iguais sempre vão para a mesma partição.
* O Kafka garante que uma partição seja consumida apenas por um consumidor.
* Então se eu colocar um `partition key` relativo aos agrupamentos que quero implementar, eu terei tudo que preciso.

Ai eu pensei nos mesmos componentes da primeira tentativa, mas com dois adicionais:

* `kafka`
* `votes-register`

Entretanto, na realidade, ao invés de ter essa lógica de dados acumulados, eu usei apenas uma view materializada do SQL para armazenar as somas parciais. Essa solução é muito mais simples, mas tem o problema de sempre recalcular tudo novamente a cada refresh.

Com isso, eu percebi que mesmo com as somas armazenadas no banco de dados o sistema ficaria suficientemente eficiente para suportar a demanda, pois ainda assim teríamos 1000 conexões sendo feitas ao segundo, mesmo que apenas leitura.

### Terceira tentativa

Para tornar as leituras viáveis, eu pensei em ferramentas de cache, ai eu pensei sobre o Redis por:

* O Redis mantem todos os dados na memória, então é muito mais rápido que a view materializada que fica na memória estável do banco de dados.
* O Redis é uma ferramenta feita para lidar com múltiplos agentes requisitando os dados ao mesmo tempo.

O Redis faria o papel de armazenar as somas dos votos: toda vez que um voto fosse registrado, também seria somado 1 a soma toda de votos, a soma de votos do participante e a soma de votos da hora.

O beneficio disso é que não preciso mais `votes-aggregator`, pois o Redis vai armazenar as somas atualizadas. Dessa forma, os componentes que pensei foram os mesmos da tentativa anterior, exceto `votes-aggregator`.

## Estrutura do código

Por questão simplicidade, eu vou manter tudo no mesmo repositório, mas seria interessante imaginar isso como sendo um projeto grande, distribuído em nesses repositórios:

* `voting-commons`: códigos comuns entres os serviços, imagine isso enquanto uma biblioteca versionada via CI/CD
  
* `voters-frontend`: responsável pela interface web dos telespectadores
* `votes-register`: responsável por consumir a fila do Kafka e salvar os dados agrupados no banco de dados
* `prodution-frontend`: responsável pela interface dos telespectadores

A estrutura interna do código é inspirada na arquitura limpa, onde tem as camadas:
* `domain`: responsável pelas entidades, não conhece nenhuma outra camada
* `service`: responsável pelos casos de uso, conhece apenas a camada `domain`
* `controller`: responsável pela API, conhece as camadas `domain` e `service`
* `data-layer`: responsável pela comunicação com os bancos de dados e kafka, conhece todos as outras camadas

  

Vale destacar o papel que inversão de dependência fez, principalmente a camada de dados, que ao longo das tentativas eu só precisei adicionar novo código e não modicar o já existente, ou me deu muito menos trabalho.

## Coisas que implementaria se fosse um projeto real

* Separa os repositórios e os módulos
  * incluir o `voting-commons` em algum repositório de artefatos, de forma que cada componente pode puxar uma versão diferente dele.
* Proteção contra bots e DDOS, como rate limite
* Monitoramento de logs e de consumo de recursos
* Incluir imagens aos participantes
* Criação de um pipeline CI/CD
  * SonarQube
